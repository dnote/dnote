package net

import (
	"github.com/dnote/dnote/pkg/server/log"
	"net/http"
)

// LifecycleWriter  wraps http.ResponseWriter to track state of the http response.
// The optional interfaces of http.ResponseWriter are lost because of the wrapping, and
// such interfaces should be implemented if needed. (i.e. http.Pusher, http.Flusher, etc.)
type LifecycleWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader wraps the WriteHeader call and marks the response state as done.
func (w *LifecycleWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// IsHeaderWritten returns true if a response has been written.
func IsHeaderWritten(w http.ResponseWriter) bool {
	if lw, ok := w.(*LifecycleWriter); ok {
		return lw.StatusCode != 0
	}

	// the response writer must have been wrapped in the middleware chain.
	log.Error("unable to log because writer is not a LifecycleWriter")
	return false
}
