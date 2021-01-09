package handlers

import (
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
)

// Middleware is a middleware for request handlers
type Middleware func(h http.Handler, app *app.App, rateLimit bool) http.Handler

// WebMw is the middleware for the web
func WebMw(h http.Handler, app *app.App, rateLimit bool) http.Handler {
	ret := h
	return ret
}

// APIMw is the middleware for the API
func APIMw(h http.Handler, app *app.App, rateLimit bool) http.Handler {
	return h
}
