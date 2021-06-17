package middleware

import (
	"net/http"
	"net/url"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/gorilla/schema"
)

// Middleware is a middleware for request handlers
type Middleware func(h http.Handler, app *app.App, rateLimit bool) http.Handler

type payload struct {
	Method string `schema:"_method"`
}

func parseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder()

	// Ignore CSRF token field
	dec.IgnoreUnknownKeys(true)

	if err := dec.Decode(dst, values); err != nil {
		return err
	}

	return nil
}

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

	ret = logging(ret)
	ret = methodOverride(ret)

	return ret
}
