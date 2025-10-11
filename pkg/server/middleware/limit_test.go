/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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
