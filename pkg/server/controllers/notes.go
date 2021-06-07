package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/operations"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/dnote/dnote/pkg/server/views"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// NewNotes creates a new Notes controller.
// It panics if the necessary templates are not parsed.
func NewNotes(app *app.App) *Notes {
	return &Notes{
		IndexView: views.NewView(app.Config.PageTemplateDir, views.Config{Title: "", Layout: "base", HeaderTemplate: "navbar"}, "notes/index"),
		ShowView:  views.NewView(app.Config.PageTemplateDir, views.Config{Title: "", Layout: "base", HeaderTemplate: "navbar"}, "notes/show"),
		app:       app,
	}
}

// Notes is a user controller.
type Notes struct {
	IndexView *views.View
	ShowView  *views.View
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

// GetNotesResponse is a reponse by getNotesHandler
type GetNotesResponse struct {
	Notes []presenters.Note `json:"notes"`
	Total int               `json:"total"`
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

func (n *Notes) getNote(r *http.Request) (database.Note, error) {
	user := context.User(r.Context())

	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	note, ok, err := operations.GetNote(n.app.DB, noteUUID, user)
	if !ok {
		return database.Note{}, app.ErrNotFound
	}
	if err != nil {
		return database.Note{}, errors.Wrap(err, "finding note")
	}

	return note, nil
}

// Show shows note
func (n *Notes) Show(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	note, err := n.getNote(r)
	if err != nil {
		handleHTMLError(w, r, err, "getting notes", n.IndexView, vd)
		return
	}

	vd.Yield = struct {
		Note database.Note
	}{
		Note: note,
	}

	n.ShowView.Render(w, r, vd)
}

// V3Show is api for show
func (n *Notes) V3Show(w http.ResponseWriter, r *http.Request) {
	note, err := n.getNote(r)
	if err != nil {
		handleJSONError(w, err, "getting note")
		return
	}

	respondJSON(w, http.StatusOK, presenters.PresentNote(note))
}

type createNotePayload struct {
	BookUUID string `schema:"book_uuid" json:"book_uuid"`
	Content  string `schema:"content" json:"content"`
	AddedOn  *int64 `schema:"added_on" json:"added_on"`
	EditedOn *int64 `schema:"edited_on" json:"edited_on"`
}

func validateCreateNotePayload(p createNotePayload) error {
	if p.BookUUID == "" {
		return app.ErrBookUUIDRequired
	}

	return nil
}

func (n *Notes) create(r *http.Request) (database.Note, error) {
	user := context.User(r.Context())
	if user == nil {
		return database.Note{}, app.ErrLoginRequired
	}

	var params createNotePayload
	if err := parseRequestData(r, &params); err != nil {
		return database.Note{}, errors.Wrap(err, "parsing request payload")
	}

	if err := validateCreateNotePayload(params); err != nil {
		return database.Note{}, err
	}

	var book database.Book
	if err := n.app.DB.Where("uuid = ? AND user_id = ?", params.BookUUID, user.ID).First(&book).Error; err != nil {
		return database.Note{}, errors.Wrap(err, "finding book")
	}

	client := getClientType(r)
	note, err := n.app.CreateNote(*user, params.BookUUID, params.Content, params.AddedOn, params.EditedOn, false, client)
	if err != nil {
		return database.Note{}, errors.Wrap(err, "creating note")
	}

	// preload associations
	note.User = *user
	note.Book = book

	return note, nil
}

func (n *Notes) del(r *http.Request) (database.Note, error) {
	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	user := context.User(r.Context())
	if user == nil {
		return database.Note{}, app.ErrLoginRequired
	}

	var note database.Note
	if err := n.app.DB.Where("uuid = ? AND user_id = ?", noteUUID, user.ID).Preload("Book").First(&note).Error; err != nil {
		return database.Note{}, errors.Wrap(err, "finding note")
	}

	tx := n.app.DB.Begin()

	note, err := n.app.DeleteNote(tx, *user, note)
	if err != nil {
		tx.Rollback()
		return database.Note{}, errors.Wrap(err, "deleting note")
	}

	tx.Commit()

	return note, nil
}

// CreateNoteResp is a response for creating a note
type CreateNoteResp struct {
	Result presenters.Note `json:"result"`
}

// Create creates note
func (n *Notes) Create(w http.ResponseWriter, r *http.Request) {
	note, err := n.create(r)
	if err != nil {
		handleJSONError(w, err, "creating note")
		return
	}

	respondJSON(w, http.StatusCreated, CreateNoteResp{
		Result: presenters.PresentNote(note),
	})
}

// V3Create creates note
func (n *Notes) V3Create(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	note, err := n.create(r)
	if err != nil {
		handleHTMLError(w, r, err, "creating note", n.IndexView, vd)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/notes/%s", note.UUID), http.StatusCreated)
}

// Delete shows note
func (n *Notes) Delete(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	_, err := n.del(r)
	if err != nil {
		handleHTMLError(w, r, err, "getting notes", n.IndexView, vd)
		return
	}

	http.Redirect(w, r, "/notes", http.StatusOK)
}

type DeleteNoteResp struct {
	Status int             `json:"status"`
	Result presenters.Note `json:"result"`
}

// V3Delete deletes note
func (n *Notes) V3Delete(w http.ResponseWriter, r *http.Request) {
	note, err := n.del(r)
	if err != nil {
		handleJSONError(w, err, "deleting note")
		return
	}

	respondJSON(w, http.StatusOK, DeleteNoteResp{
		Status: http.StatusNoContent,
		Result: presenters.PresentNote(note),
	})
}

type updateNotePayload struct {
	BookUUID *string `schema:"book_uuid" json:"book_uuid"`
	Content  *string `schema:"content" json:"content"`
	Public   *bool   `schema:"public" json:"public"`
}

func validateUpdateNotePayload(p updateNotePayload) error {
	if p.BookUUID == nil && p.Content == nil && p.Public == nil {
		return app.ErrEmptyUpdate
	}

	return nil
}

func (n *Notes) update(r *http.Request) (database.Note, error) {
	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	user := context.User(r.Context())
	if user == nil {
		return database.Note{}, app.ErrLoginRequired
	}

	var params updateNotePayload
	err := parseRequestData(r, &params)
	if err != nil {
		return database.Note{}, errors.Wrap(err, "decoding params")
	}

	if err := validateUpdateNotePayload(params); err != nil {
		return database.Note{}, err
	}

	var note database.Note
	if err := n.app.DB.Where("uuid = ? AND user_id = ?", noteUUID, user.ID).First(&note).Error; err != nil {
		return database.Note{}, errors.Wrap(err, "finding note")
	}

	tx := n.app.DB.Begin()

	note, err = n.app.UpdateNote(tx, *user, note, &app.UpdateNoteParams{
		BookUUID: params.BookUUID,
		Content:  params.Content,
		Public:   params.Public,
	})
	if err != nil {
		tx.Rollback()
		return database.Note{}, errors.Wrap(err, "updating note")
	}

	var book database.Book
	if err := tx.Where("uuid = ? AND user_id = ?", note.BookUUID, user.ID).First(&book).Error; err != nil {
		tx.Rollback()
		return database.Note{}, errors.Wrapf(err, "finding book %s to preload", note.BookUUID)
	}

	tx.Commit()

	// preload associations
	note.User = *user
	note.Book = book

	return note, nil
}

type updateNoteResp struct {
	Status int             `json:"status"`
	Result presenters.Note `json:"result"`
}

// V3Update updates a note
func (n *Notes) V3Update(w http.ResponseWriter, r *http.Request) {
	note, err := n.update(r)
	if err != nil {
		handleJSONError(w, err, "updating note")
		return
	}

	respondJSON(w, http.StatusOK, updateNoteResp{
		Status: http.StatusOK,
		Result: presenters.PresentNote(note),
	})
}

// Update updates a note
func (n *Notes) Update(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	note, err := n.update(r)
	if err != nil {
		handleHTMLError(w, r, err, "updating note", n.IndexView, vd)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/notes/%s", note.UUID), http.StatusOK)
}
