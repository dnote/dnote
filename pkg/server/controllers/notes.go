package controllers

import (
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/models"
	"github.com/dnote/dnote/pkg/server/views"
	// "github.com/nadproject/nad/pkg/server/context"
)

// NewNotes creates a new Notes controller.
// It panics if the necessary templates are not parsed.
func NewNotes(cfg config.Config, app *app.App) *Notes {
	return &Notes{
		IndexView: views.NewView(cfg.PageTemplateDir, views.Config{Title: "", Layout: "base", HeaderTemplate: "navbar"}, "notes/index"),
		app:       app,
	}
}

// Notes is a user controller.
type Notes struct {
	IndexView *views.View
	app       *app.App
}

// Index handles GET /
func (n *Notes) Index(w http.ResponseWriter, r *http.Request) {
	// user := context.User(r.Context())

	var vd views.Data
	vd.Yield = struct {
		Notes []models.Note
	}{
		Notes: nil,
	}

	n.IndexView.Render(w, r, vd)
}
