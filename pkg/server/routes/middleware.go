package routes

import (
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
)

type middleware func(h http.Handler, app *app.App, rateLimit bool) http.Handler

// WebMw is the middleware for the web
func WebMw(h http.Handler, app *app.App, rateLimit bool) http.Handler {
	ret := h
	return ret
}

// APIMw is the middleware for the API
func APIMw(h http.Handler, app *app.App, rateLimit bool) http.Handler {
	return h
}
