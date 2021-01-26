package controllers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/dnote/dnote/pkg/server/views"
	"github.com/pkg/errors"
)

// NewNotes creates a new Notes controller.
// It panics if the necessary templates are not parsed.
func NewNotes(app *app.App) *Notes {
	return &Notes{
		IndexView: views.NewView(app.Config.PageTemplateDir, views.Config{Title: "", Layout: "base", HeaderTemplate: "navbar"}, "notes/index"),
		app:       app,
	}
}

// Notes is a user controller.
type Notes struct {
	IndexView *views.View
	app       *app.App
}

// escapeSearchQuery escapes the query for full text search
func escapeSearchQuery(searchQuery string) string {
	return strings.Join(strings.Fields(searchQuery), "&")
}

func parseSearchQuery(q url.Values) string {
	searchStr := q.Get("q")

	return escapeSearchQuery(searchStr)
}

func parseGetNotesQuery(q url.Values) (app.GetNotesParams, error) {
	yearStr := q.Get("year")
	monthStr := q.Get("month")
	books := q["book"]
	pageStr := q.Get("page")
	encryptedStr := q.Get("encrypted")

	var page int
	if len(pageStr) > 0 {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return app.GetNotesParams{}, errors.Errorf("invalid page %s", pageStr)
		}
		if p < 1 {
			return app.GetNotesParams{}, errors.Errorf("invalid page %s", pageStr)
		}

		page = p
	} else {
		page = 1
	}

	var year int
	if len(yearStr) > 0 {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			return app.GetNotesParams{}, errors.Errorf("invalid year %s", yearStr)
		}

		year = y
	}

	var month int
	if len(monthStr) > 0 {
		m, err := strconv.Atoi(monthStr)
		if err != nil {
			return app.GetNotesParams{}, errors.Errorf("invalid month %s", monthStr)
		}
		if m < 1 || m > 12 {
			return app.GetNotesParams{}, errors.Errorf("invalid month %s", monthStr)
		}

		month = m
	}

	var encrypted bool
	if strings.ToLower(encryptedStr) == "true" {
		encrypted = true
	} else {
		encrypted = false
	}

	ret := app.GetNotesParams{
		Year:      year,
		Month:     month,
		Page:      page,
		Search:    parseSearchQuery(q),
		Books:     books,
		Encrypted: encrypted,
	}

	return ret, nil
}

func (n *Notes) getNotes(r *http.Request) (app.GetNotesResult, error) {
	user := context.User(r.Context())
	if user == nil {
		return app.GetNotesResult{}, app.ErrLoginRequired
	}

	query := r.URL.Query()
	p, err := parseGetNotesQuery(query)
	if err != nil {
		return app.GetNotesResult{}, errors.Wrap(err, "parsing query")
	}

	res, err := n.app.GetNotes(user.ID, p)
	if err != nil {
		return app.GetNotesResult{}, errors.Wrap(err, "getting notes")
	}

	return res, nil
}

// Index handles GET /
func (n *Notes) Index(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	result, err := n.getNotes(r)
	if err != nil {
		handleHTMLError(w, r, err, "getting notes", n.IndexView, vd)
		return
	}

	vd.Yield = struct {
		Notes []database.Note
		Total int
	}{
		Notes: result.Notes,
		Total: result.Total,
	}

	n.IndexView.Render(w, r, vd)
}

// V3Index is a v3 handler for getting notes
func (n *Notes) V3Index(w http.ResponseWriter, r *http.Request) {
	result, err := n.getNotes(r)
	if err != nil {
		handleJSONError(w, err, "getting notes")
		return
	}

	respondJSON(w, http.StatusOK, GetNotesResponse{
		Notes: presenters.PresentNotes(result.Notes),
		Total: result.Total,
	})
}

// GetNotesResponse is a reponse by getNotesHandler
type GetNotesResponse struct {
	Notes []presenters.Note `json:"notes"`
	Total int               `json:"total"`
}
