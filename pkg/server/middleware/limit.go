/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dnote/dnote/pkg/server/log"
	"golang.org/x/time/rate"
)

const (
	// serverRateLimitPerSecond is the max requests per second the server will accept per IP
	serverRateLimitPerSecond = 50
	// serverRateLimitBurst is the burst capacity for rate limiting
	serverRateLimitBurst = 100
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter holds the rate limiting state for visitors
type RateLimiter struct {
	visitors map[string]*visitor
	mtx      sync.RWMutex
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
	}
	go rl.cleanupVisitors()
	return rl
}

var defaultLimiter = NewRateLimiter()

// addVisitor adds a new visitor to the map and returns a limiter for the visitor
func (rl *RateLimiter) addVisitor(identifier string) *rate.Limiter {
	// Calculate interval from rate: 1 second / requests per second
	interval := time.Second / time.Duration(serverRateLimitPerSecond)
	limiter := rate.NewLimiter(rate.Every(interval), serverRateLimitBurst)

	rl.mtx.Lock()
	rl.visitors[identifier] = &visitor{
		limiter:  limiter,
		lastSeen: time.Now()}
	rl.mtx.Unlock()

	return limiter
}

// getVisitor returns a limiter for a visitor with the given identifier. It
// adds the visitor to the map if not seen before.
func (rl *RateLimiter) getVisitor(identifier string) *rate.Limiter {
	rl.mtx.RLock()
	v, exists := rl.visitors[identifier]

	if !exists {
		rl.mtx.RUnlock()
		return rl.addVisitor(identifier)
	}

	v.lastSeen = time.Now()
	rl.mtx.RUnlock()

	return v.limiter
}

// cleanupVisitors deletes visitors that has not been seen in a while from the
// map of visitors
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mtx.Lock()

		for identifier, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, identifier)
			}
		}

		rl.mtx.Unlock()
	}
}

// lookupIP returns the request's IP
func lookupIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP")
	forwardedFor := r.Header.Get("X-Forwarded-For")

	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		return parts[0]
	}

	if realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}

// Limit is a middleware to rate limit the handler
func (rl *RateLimiter) Limit(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identifier := lookupIP(r)
		limiter := rl.getVisitor(identifier)

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			log.WithFields(log.Fields{
				"ip": identifier,
			}).Warn("Too many requests")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ApplyLimit applies rate limit conditionally using the global limiter
func ApplyLimit(h http.HandlerFunc, rateLimit bool) http.Handler {
	ret := h

	if rateLimit {
		ret = defaultLimiter.Limit(ret)
	}

	return ret
}
