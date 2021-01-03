package routes

import (
	"net/http"
	"strings"

	"github.com/dnote/dnote/pkg/server/config"
)

type middleware func(h http.Handler, c config.Config, rateLimit bool) http.Handler

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

// WebMw is the middleware for the web
func WebMw(h http.Handler, c config.Config, rateLimit bool) http.Handler {
	return h
}

// APIMw is the middleware for the API
func APIMw(h http.Handler, c config.Config, rateLimit bool) http.Handler {
	return h
}
