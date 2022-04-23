package controllers

import (
	"fmt"
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/dnote/dnote/pkg/server/views"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// NewBooks creates a new Books controller.
// It panics if the necessary templates are not parsed.
func NewBooks(app *app.App) *Books {
	return &Books{
		IndexView: views.NewView(app, views.Config{Title: "", Layout: "base", HeaderTemplate: "navbar"}, "books/index"),
		ShowView:  views.NewView(app, views.Config{Title: "", Layout: "base", HeaderTemplate: "navbar"}, "books/show"),
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

	vd.Yield = map[string]interface{}{
		"Books": result,
	}

	b.IndexView.Render(w, r, &vd, http.StatusOK)
}

// V3Index gets books
func (b *Books) V3Index(w http.ResponseWriter, r *http.Request) {
	result, err := b.getBooks(r)
	if err != nil {
		handleJSONError(w, err, "getting books")
		return
	}

	respondJSON(w, http.StatusOK, presenters.PresentBooks(result))
}

// V3Show gets a book
func (b *Books) V3Show(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user == nil {
		handleJSONError(w, app.ErrLoginRequired, "login required")
		return
	}

	vars := mux.Vars(r)
	bookUUID := vars["bookUUID"]

	if !helpers.ValidateUUID(bookUUID) {
		handleJSONError(w, app.ErrInvalidUUID, "login required")
		return
	}

	var book database.Book
	conn := b.app.DB.Where("uuid = ? AND user_id = ?", bookUUID, user.ID).First(&book)

	if conn.RecordNotFound() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := conn.Error; err != nil {
		handleJSONError(w, err, "finding the book")
		return
	}

	respondJSON(w, http.StatusOK, presenters.PresentBook(book))
}

type createBookPayload struct {
	Name string `schema:"name" json:"name"`
}

func validateCreateBookPayload(p createBookPayload) error {
	if p.Name == "" {
		return app.ErrBookNameRequired
	}

	return nil
}

func (b *Books) create(r *http.Request) (database.Book, error) {
	user := context.User(r.Context())
	if user == nil {
		return database.Book{}, app.ErrLoginRequired
	}

	var params createBookPayload
	if err := parseRequestData(r, &params); err != nil {
		return database.Book{}, errors.Wrap(err, "parsing request payload")
	}

	if err := validateCreateBookPayload(params); err != nil {
		return database.Book{}, errors.Wrap(err, "validating payload")
	}

	var bookCount int
	err := b.app.DB.Model(database.Book{}).
		Where("user_id = ? AND label = ?", user.ID, params.Name).
		Count(&bookCount).Error
	if err != nil {
		return database.Book{}, errors.Wrap(err, "checking duplicate")
	}
	if bookCount > 0 {
		return database.Book{}, app.ErrDuplicateBook
	}

	book, err := b.app.CreateBook(*user, params.Name)
	if err != nil {
		return database.Book{}, errors.Wrap(err, "inserting a book")
	}

	return book, nil
}

// Create creates a book
func (b *Books) Create(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	_, err := b.create(r)
	if err != nil {
		handleHTMLError(w, r, err, "creating a books", b.IndexView, vd)
		return
	}

	http.Redirect(w, r, "/books", http.StatusCreated)
}

// createBookResp is the response from create book api
type createBookResp struct {
	Book presenters.Book `json:"book"`
}

// V3Create creates a book
func (b *Books) V3Create(w http.ResponseWriter, r *http.Request) {
	result, err := b.create(r)
	if err != nil {
		handleJSONError(w, err, "creating a book")
		return
	}

	resp := createBookResp{
		Book: presenters.PresentBook(result),
	}
	respondJSON(w, http.StatusCreated, resp)
}

type updateBookPayload struct {
	Name *string `schema:"name" json:"name"`
}

// UpdateBookResp is the response from create book api
type UpdateBookResp struct {
	Book presenters.Book `json:"book"`
}

func (b *Books) update(r *http.Request) (database.Book, error) {
	user := context.User(r.Context())
	if user == nil {
		return database.Book{}, app.ErrLoginRequired
	}

	vars := mux.Vars(r)
	uuid := vars["bookUUID"]

	if !helpers.ValidateUUID(uuid) {
		return database.Book{}, app.ErrInvalidUUID
	}

	tx := b.app.DB.Begin()

	var book database.Book
	if err := tx.Where("user_id = ? AND uuid = ?", user.ID, uuid).First(&book).Error; err != nil {
		return database.Book{}, errors.Wrap(err, "finding book")
	}

	var params updateBookPayload
	if err := parseRequestData(r, &params); err != nil {
		return database.Book{}, errors.Wrap(err, "decoding payload")
	}

	book, err := b.app.UpdateBook(tx, *user, book, params.Name)
	if err != nil {
		tx.Rollback()
		return database.Book{}, errors.Wrap(err, "updating a book")
	}

	tx.Commit()

	return book, nil
}

// Update updates a book
func (b *Books) Update(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	_, err := b.update(r)
	if err != nil {
		handleHTMLError(w, r, err, "creating a books", b.IndexView, vd)
		return
	}

	http.Redirect(w, r, "/books", http.StatusOK)
}

// V3Update updates a book
func (b *Books) V3Update(w http.ResponseWriter, r *http.Request) {
	book, err := b.update(r)
	if err != nil {
		handleJSONError(w, err, "updating a book")
		return
	}

	resp := UpdateBookResp{
		Book: presenters.PresentBook(book),
	}
	respondJSON(w, http.StatusOK, resp)
}

func (b *Books) del(r *http.Request) (database.Book, error) {
	user := context.User(r.Context())
	if user == nil {
		return database.Book{}, app.ErrLoginRequired
	}

	vars := mux.Vars(r)
	uuid := vars["bookUUID"]

	if !helpers.ValidateUUID(uuid) {
		return database.Book{}, app.ErrInvalidUUID
	}

	tx := b.app.DB.Begin()

	var book database.Book
	if err := tx.Where("user_id = ? AND uuid = ?", user.ID, uuid).First(&book).Error; err != nil {
		return database.Book{}, errors.Wrap(err, "finding a book")
	}

	var notes []database.Note
	if err := tx.Where("book_uuid = ? AND NOT deleted", uuid).Order("usn ASC").Find(&notes).Error; err != nil {
		return database.Book{}, errors.Wrap(err, "finding notes for the book")
	}

	for _, note := range notes {
		if _, err := b.app.DeleteNote(tx, *user, note); err != nil {
			tx.Rollback()
			return database.Book{}, errors.Wrap(err, "deleting a note in the book")
		}
	}

	book, err := b.app.DeleteBook(tx, *user, book)
	if err != nil {
		return database.Book{}, errors.Wrap(err, "deleting the book")
	}

	tx.Commit()

	return book, nil
}

// Delete updates a book
func (b *Books) Delete(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	_, err := b.del(r)
	if err != nil {
		handleHTMLError(w, r, err, "creating a books", b.IndexView, vd)
		return
	}

	http.Redirect(w, r, "/", http.StatusOK)
}

// deleteBookResp is the response from create book api
type deleteBookResp struct {
	Status int             `json:"status"`
	Book   presenters.Book `json:"book"`
}

// Delete updates a book
func (b *Books) V3Delete(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	book, err := b.del(r)
	if err != nil {
		handleHTMLError(w, r, err, "creating a books", b.IndexView, vd)
		return
	}

	resp := deleteBookResp{
		Status: http.StatusOK,
		Book:   presenters.PresentBook(book),
	}
	respondJSON(w, http.StatusOK, resp)
}

// IndexOptions is a handler for OPTIONS endpoint for notes
func (b *Books) IndexOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Version")
}
