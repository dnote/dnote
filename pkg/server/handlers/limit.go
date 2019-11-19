/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package handlers

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dnote/dnote/pkg/server/log"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mtx sync.RWMutex

func init() {
	go cleanupVisitors()
}

// addVisitor adds a new visitor to the map and returns a limiter for the visitor
func addVisitor(identifier string) *rate.Limiter {
	// initialize a token bucket
	limiter := rate.NewLimiter(rate.Every(1*time.Second), 60)

	mtx.Lock()
	visitors[identifier] = &visitor{
		limiter:  limiter,
		lastSeen: time.Now()}
	mtx.Unlock()

	return limiter
}

// getVisitor returns a limiter for a visitor with the given identifier. It
// adds the visitor to the map if not seen before.
func getVisitor(identifier string) *rate.Limiter {
	mtx.RLock()
	v, exists := visitors[identifier]

	if !exists {
		mtx.RUnlock()
		return addVisitor(identifier)
	}

	v.lastSeen = time.Now()
	mtx.RUnlock()

	return v.limiter
}

// cleanupVisitors deletes visitors that has not been seen in a while from the
// map of visitors
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mtx.Lock()

		for identifier, v := range visitors {
			if time.Now().Sub(v.lastSeen) > 3*time.Minute {
				delete(visitors, identifier)
			}
		}

		mtx.Unlock()
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

// limit is a middleware to rate limit the handler
func limit(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identifier := lookupIP(r)
		limiter := getVisitor(identifier)

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
