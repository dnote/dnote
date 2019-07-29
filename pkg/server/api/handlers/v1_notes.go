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

	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/api/operations"
	"github.com/dnote/dnote/pkg/server/api/presenters"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/gorilla/mux"
)

// CreateNote creates a note by generating an action and feeding it to the reducer
func (a *App) CreateNote(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not supported. Please upgrade your client.", http.StatusGone)
	return
}

// NotesOptions is a handler for OPTIONS endpoint for notes
func (a *App) NotesOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Version")
}

type updateNotePayload struct {
	BookUUID *string `json:"book_uuid"`
	Content  *string `json:"content"`
}

type updateNoteResp struct {
	Status int             `json:"status"`
	Result presenters.Note `json:"result"`
}

func validateUpdateNotePayload(p updateNotePayload) bool {
	return p.BookUUID != nil || p.Content != nil
}

// UpdateNote updates note
func (a *App) UpdateNote(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn
	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var params updateNotePayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		handleError(w, "decoding params", err, http.StatusInternalServerError)
		return
	}

	if ok := validateUpdateNotePayload(params); !ok {
		handleError(w, "Invalid payload", nil, http.StatusBadRequest)
		return
	}

	var note database.Note
	if err := db.Where("uuid = ? AND user_id = ?", noteUUID, user.ID).First(&note).Error; err != nil {
		handleError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	tx := db.Begin()

	note, err = operations.UpdateNote(tx, user, a.Clock, note, params.BookUUID, params.Content)
	if err != nil {
		tx.Rollback()
		handleError(w, "updating note", err, http.StatusInternalServerError)
		return
	}

	var book database.Book
	if err := tx.Where("uuid = ? AND user_id = ?", note.BookUUID, user.ID).First(&book).Error; err != nil {
		tx.Rollback()
		handleError(w, fmt.Sprintf("finding book %s to preload", note.BookUUID), err, http.StatusInternalServerError)
		return
	}

	tx.Commit()

	// preload associations
	note.User = user
	note.Book = book

	resp := updateNoteResp{
		Status: http.StatusOK,
		Result: presenters.PresentNote(note),
	}
	respondJSON(w, resp)
}

type deleteNoteResp struct {
	Status int             `json:"status"`
	Result presenters.Note `json:"result"`
}

// DeleteNote removes note
func (a *App) DeleteNote(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var note database.Note
	if err := db.Where("uuid = ? AND user_id = ?", noteUUID, user.ID).Preload("Book").First(&note).Error; err != nil {
		handleError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	tx := db.Begin()

	n, err := operations.DeleteNote(tx, user, note)
	if err != nil {
		tx.Rollback()
		handleError(w, "deleting note", err, http.StatusInternalServerError)
		return
	}

	tx.Commit()

	resp := deleteNoteResp{
		Status: http.StatusNoContent,
		Result: presenters.PresentNote(n),
	}
	respondJSON(w, resp)
}
