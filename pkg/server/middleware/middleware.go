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
