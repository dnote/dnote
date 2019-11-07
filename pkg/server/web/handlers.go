// Package web provides handlers for the web application
package web

import (
	"net/http"

	"github.com/dnote/dnote/pkg/server/handlers"
	"github.com/dnote/dnote/pkg/server/tmpl"
	"github.com/pkg/errors"
)

// Context contains contents of web assets
type Context struct {
	IndexHTML        []byte
	RobotsTxt        []byte
	ServiceWorkerJs  []byte
	StaticFileSystem http.FileSystem
}

// GetRootHandler returns an HTTP handler that serves the app shell
func GetRootHandler(b []byte) http.HandlerFunc {
	appShell, err := tmpl.NewAppShell(b)
	if err != nil {
		panic(errors.Wrap(err, "initializing app shell"))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// index.html must not be cached
		w.Header().Set("Cache-Control", "no-cache")

		buf, err := appShell.Execute(r)
		if err != nil {
			if errors.Cause(err) == tmpl.ErrNotFound {
				handlers.RespondNotFound(w)
			} else {
				handlers.HandleError(w, "executing app shell", err, http.StatusInternalServerError)
			}
			return
		}

		w.Write(buf)
	}
}

// GetRobotsHandler returns an HTTP handler that serves robots.txt
func GetRobotsHandler(b []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(b)
	}
}

// GetSWHandler returns an HTTP handler that serves service worker
func GetSWHandler(b []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(b)
	}
}

// GetStaticHandler returns an HTTP handler that serves static files from a filesystem
func GetStaticHandler(root http.FileSystem) http.Handler {
	return http.StripPrefix("/static/", http.FileServer(root))
}
