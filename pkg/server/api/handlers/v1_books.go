/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/api/operations"
	"github.com/dnote/dnote/pkg/server/api/presenters"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func respondWithBooks(userID int, query url.Values, w http.ResponseWriter) {
	db := database.DBConn

	var books []database.Book
	conn := db.Where("user_id = ? AND NOT deleted", userID).Order("label ASC")
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

	if err := conn.Find(&books).Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding books").Error(), http.StatusInternalServerError)
		return
	}

	presentedBooks := presenters.PresentBooks(books)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(presentedBooks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetDemoBooks returns books for demo
func (a *App) GetDemoBooks(w http.ResponseWriter, r *http.Request) {
	demoUserID, err := helpers.GetDemoUserID()
	if err != nil {
		http.Error(w, errors.Wrap(err, "finding demo user").Error(), http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()

	respondWithBooks(demoUserID, query, w)
}

// GetBooks returns books for the user
func (a *App) GetBooks(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		return
	}

	query := r.URL.Query()

	respondWithBooks(user.ID, query, w)
}

// GetBook returns a book for the user
func (a *App) GetBook(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		return
	}

	db := database.DBConn

	vars := mux.Vars(r)
	bookUUID := vars["bookUUID"]

	var book database.Book
	conn := db.Where("uuid = ? AND user_id = ?", bookUUID, user.ID).First(&book)

	if conn.RecordNotFound() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := conn.Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding book").Error(), http.StatusInternalServerError)
		return
	}

	p := presenters.PresentBook(book)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type createBookPayload struct {
	Name string `json:"name"`
}

// CreateBookResp is the response from create book api
type CreateBookResp struct {
	Book presenters.Book `json:"book"`
}

func validateCreateBookPayload(p createBookPayload) error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

// CreateBook creates a new book
func (a *App) CreateBook(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not supported. Please upgrade your client.", http.StatusGone)
	return
}

type updateBookPayload struct {
	Name *string `json:"name"`
}

// UpdateBookResp is the response from create book api
type UpdateBookResp struct {
	Book presenters.Book `json:"book"`
}

// UpdateBook updates a book
func (a *App) UpdateBook(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	uuid := vars["bookUUID"]

	db := database.DBConn
	tx := db.Begin()

	var book database.Book
	if err := tx.Where("user_id = ? AND uuid = ?", user.ID, uuid).First(&book).Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding book").Error(), http.StatusInternalServerError)
		return
	}

	var params updateBookPayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, errors.Wrap(err, "decoding payload").Error(), http.StatusInternalServerError)
		return
	}

	book, err = operations.UpdateBook(tx, a.Clock, user, book, params.Name)
	if err != nil {
		tx.Rollback()
		http.Error(w, errors.Wrap(err, "updating a book").Error(), http.StatusInternalServerError)
	}

	tx.Commit()

	resp := UpdateBookResp{
		Book: presenters.PresentBook(book),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteBookResp is the response from create book api
type DeleteBookResp struct {
	Status int             `json:"status"`
	Book   presenters.Book `json:"book"`
}

// DeleteBook removes a book
func (a *App) DeleteBook(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	uuid := vars["bookUUID"]

	db := database.DBConn
	tx := db.Begin()

	var book database.Book
	if err := tx.Where("user_id = ? AND uuid = ?", user.ID, uuid).First(&book).Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding book").Error(), http.StatusInternalServerError)
		return
	}

	var notes []database.Note
	if err := tx.Where("book_uuid = ? AND NOT deleted", uuid).Order("usn ASC").Find(&notes).Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding notes").Error(), http.StatusInternalServerError)
		return
	}

	for _, note := range notes {
		if _, err := operations.DeleteNote(tx, user, note); err != nil {
			http.Error(w, errors.Wrap(err, "deleting a note").Error(), http.StatusInternalServerError)
			return
		}
	}
	b, err := operations.DeleteBook(tx, user, book)
	if err != nil {
		http.Error(w, errors.Wrap(err, "deleting book").Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit()

	resp := DeleteBookResp{
		Status: http.StatusOK,
		Book:   presenters.PresentBook(b),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// BooksOptions handles OPTIONS endpoint
func (a *App) BooksOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Version")
}
