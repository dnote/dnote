package controllers

import (
	"net/http"
	"strings"

	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/views"
)

// NewStatic creates a new Static controller.
func NewStatic(cfg config.Config) *Static {
	return &Static{
		NotFoundView: views.NewView(cfg.PageTemplateDir, views.Config{Title: "Not Found", Layout: "base"}, "static/not_found"),
	}
}

// Static is a static controller
type Static struct {
	NotFoundView *views.View
}

// NotFound is a catch-all handler for requests with no matching handler
func (s *Static) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "text/html") {
		s.NotFoundView.Render(w, r, nil)
	} else {
		statusText := http.StatusText(http.StatusNotFound)
		w.Write([]byte(statusText))
	}
}
