package controllers

import (
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
)

// NewHealth creates a new Health controller.
// It panics if the necessary templates are not parsed.
func NewHealth(app *app.App) *Health {
	return &Health{}
}

// Health is a health controller.
type Health struct {
}

// Index handles GET /
func (n *Health) Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
