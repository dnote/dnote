/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	pkgErrors "github.com/pkg/errors"
)

// NewBooks creates a new Books controller.
// It panics if the necessary templates are not parsed.
func NewBooks(app *app.App) *Books {
	return &Books{
		app: app,
	}
}

// Books is a user controller.
type Books struct {
	app *app.App
}

func (b *Books) getBooks(r *http.Request) ([]database.Book, error) {
	user := context.User(r.Context())
	if user == nil {
		return []database.Book{}, app.ErrLoginRequired
	}

	conn := b.app.DB.Where("user_id = ? AND NOT deleted", user.ID).Order("label ASC")

	query := r.URL.Query()
	name := query.Get("name")

	if name != "" {
		part := fmt.Sprintf("%%%s%%", name)
		conn = conn.Where("LOWER(label) LIKE ?", part)
	}

	var books []database.Book
	if err := conn.Find(&books).Error; err != nil {
		return []database.Book{}, nil
	}

	return books, nil
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
	err := b.app.DB.Where("uuid = ? AND user_id = ?", bookUUID, user.ID).First(&book).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
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
		return database.Book{}, pkgErrors.Wrap(err, "parsing request payload")
	}

	if err := validateCreateBookPayload(params); err != nil {
		return database.Book{}, pkgErrors.Wrap(err, "validating payload")
	}

	var bookCount int64
	err := b.app.DB.Model(database.Book{}).
		Where("user_id = ? AND label = ?", user.ID, params.Name).
		Count(&bookCount).Error
	if err != nil {
		return database.Book{}, pkgErrors.Wrap(err, "checking duplicate")
	}
	if bookCount > 0 {
		return database.Book{}, app.ErrDuplicateBook
	}

	book, err := b.app.CreateBook(*user, params.Name)
	if err != nil {
		return database.Book{}, pkgErrors.Wrap(err, "inserting a book")
	}

	return book, nil
}

// CreateBookResp is the response from create book api
type CreateBookResp struct {
	Book presenters.Book `json:"book"`
}

// V3Create creates a book
func (b *Books) V3Create(w http.ResponseWriter, r *http.Request) {
	result, err := b.create(r)
	if err != nil {
		handleJSONError(w, err, "creating a book")
		return
	}

	resp := CreateBookResp{
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
		return database.Book{}, pkgErrors.Wrap(err, "finding book")
	}

	var params updateBookPayload
	if err := parseRequestData(r, &params); err != nil {
		return database.Book{}, pkgErrors.Wrap(err, "decoding payload")
	}

	book, err := b.app.UpdateBook(tx, *user, book, params.Name)
	if err != nil {
		tx.Rollback()
		return database.Book{}, pkgErrors.Wrap(err, "updating a book")
	}

	tx.Commit()

	return book, nil
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
		return database.Book{}, pkgErrors.Wrap(err, "finding a book")
	}

	var notes []database.Note
	if err := tx.Where("book_uuid = ? AND NOT deleted", uuid).Order("usn ASC").Find(&notes).Error; err != nil {
		return database.Book{}, pkgErrors.Wrap(err, "finding notes for the book")
	}

	for _, note := range notes {
		if _, err := b.app.DeleteNote(tx, *user, note); err != nil {
			tx.Rollback()
			return database.Book{}, pkgErrors.Wrap(err, "deleting a note in the book")
		}
	}

	book, err := b.app.DeleteBook(tx, *user, book)
	if err != nil {
		return database.Book{}, pkgErrors.Wrap(err, "deleting the book")
	}

	tx.Commit()

	return book, nil
}

// deleteBookResp is the response from create book api
type deleteBookResp struct {
	Status int             `json:"status"`
	Book   presenters.Book `json:"book"`
}

// Delete updates a book
func (b *Books) V3Delete(w http.ResponseWriter, r *http.Request) {
	book, err := b.del(r)
	if err != nil {
		handleJSONError(w, err, "creating a books")
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
