package controllers

import (
	"fmt"
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/dnote/dnote/pkg/server/views"
)

// NewBooks creates a new Books controller.
// It panics if the necessary templates are not parsed.
func NewBooks(app *app.App) *Books {
	return &Books{
		IndexView: views.NewView(app.Config.PageTemplateDir, views.Config{Title: "", Layout: "base", HeaderTemplate: "navbar"}, "books/index"),
		ShowView:  views.NewView(app.Config.PageTemplateDir, views.Config{Title: "", Layout: "base", HeaderTemplate: "navbar"}, "books/show"),
		app:       app,
	}
}

// Books is a user controller.
type Books struct {
	IndexView *views.View
	ShowView  *views.View
	app       *app.App
}

func (b *Books) getBooks(r *http.Request) ([]database.Book, error) {
	user := context.User(r.Context())
	if user == nil {
		return []database.Book{}, app.ErrLoginRequired
	}

	conn := b.app.DB.Where("user_id = ? AND NOT deleted", user.ID).Order("label ASC")

	query := r.URL.Query()
	name := query.Get("name")
	encryptedStr := query.Get("encrypted")

	if name != "" {
		part := fmt.Sprintf("%%%s%%", name)
		conn = conn.Where("LOWER(label) LIKE ?", part)
	}
	if encryptedStr != "" {
		var encrypted bool
		if encryptedStr == "true" {
			encrypted = true
		} else {
			encrypted = false
		}

		conn = conn.Where("encrypted = ?", encrypted)
	}

	var books []database.Book
	if err := conn.Find(&books).Error; err != nil {
		return []database.Book{}, nil
	}

	return books, nil
}

// Index handles GET /
func (b *Books) Index(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	result, err := b.getBooks(r)
	if err != nil {
		handleHTMLError(w, r, err, "getting books", b.IndexView, vd)
		return
	}

	vd.Yield = struct {
		Books []database.Book
	}{
		Books: result,
	}

	b.IndexView.Render(w, r, vd)
}

// Index gets books
func (b *Books) V3Index(w http.ResponseWriter, r *http.Request) {
	result, err := b.getBooks(r)
	if err != nil {
		handleJSONError(w, err, "getting books")
		return
	}

	respondJSON(w, http.StatusOK, presenters.PresentBooks(result))
}
