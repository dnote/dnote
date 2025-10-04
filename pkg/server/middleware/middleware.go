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

	"github.com/dnote/dnote/pkg/server/app"
)

// Middleware is a middleware for request handlers
type Middleware func(h http.Handler, app *app.App, rateLimit bool) http.Handler

// methodOverrideKey is the form key for overriding the method
var methodOverrideKey = "_method"

// methodOverride overrides the request's method to simulate form actions that
// are not natively supported by web browsers
func methodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			method := r.PostFormValue(methodOverrideKey)

			if method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete {
				r.Method = method
			}
		}

		next.ServeHTTP(w, r)
	})
}

// WebMw is the middleware for the web
func WebMw(h http.Handler, app *app.App, rateLimit bool) http.Handler {
	ret := h

	ret = ApplyLimit(ret.ServeHTTP, rateLimit)

	return ret
}

// APIMw is the middleware for the API
func APIMw(h http.Handler, app *app.App, rateLimit bool) http.Handler {
	ret := h

	ret = ApplyLimit(ret.ServeHTTP, rateLimit)

	return ret
}

// Global is the middleware for all routes
func Global(h http.Handler) http.Handler {
	ret := h

	ret = Logging(ret)
	ret = methodOverride(ret)

	return ret
}
