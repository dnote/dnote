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
	"net/http/httptest"
	"testing"
)

func TestLimit(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limiter := NewRateLimiter()
	middleware := limiter.Limit(handler)

	// Make burst + 5 requests from same IP
	numRequests := serverRateLimitBurst + 5
	blockedCount := 0

	for range numRequests {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code == http.StatusTooManyRequests {
			blockedCount++
		}
	}

	// At least some requests after burst should be blocked
	if blockedCount == 0 {
		t.Error("Expected some requests to be rate limited after burst")
	}
}

func TestLimit_DifferentIPs(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limiter := NewRateLimiter()
	middleware := limiter.Limit(handler)

	// Exhaust rate limit for first IP
	for range serverRateLimitBurst + 5 {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		middleware.ServeHTTP(w, req)
	}

	// Request from different IP should still succeed
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:5678"
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Request from different IP should succeed, got status %d", w.Code)
	}
}
