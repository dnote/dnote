/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/handlers"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type updateNotePayload struct {
	BookUUID *string `json:"book_uuid"`
	Content  *string `json:"content"`
	Public   *bool   `json:"public"`
}

type updateNoteResp struct {
	Status int             `json:"status"`
	Result presenters.Note `json:"result"`
}

func validateUpdateNotePayload(p updateNotePayload) bool {
	return p.BookUUID != nil || p.Content != nil || p.Public != nil
}

// UpdateNote updates note
func (a *API) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handlers.DoError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var params updateNotePayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		handlers.DoError(w, "decoding params", err, http.StatusInternalServerError)
		return
	}

	if ok := validateUpdateNotePayload(params); !ok {
		handlers.DoError(w, "Invalid payload", nil, http.StatusBadRequest)
		return
	}

	var note database.Note
	if err := a.App.DB.Where("uuid = ? AND user_id = ?", noteUUID, user.ID).First(&note).Error; err != nil {
		handlers.DoError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	tx := a.App.DB.Begin()

	note, err = a.App.UpdateNote(tx, user, note, &app.UpdateNoteParams{
		BookUUID: params.BookUUID,
		Content:  params.Content,
		Public:   params.Public,
	})
	if err != nil {
		tx.Rollback()
		handlers.DoError(w, "updating note", err, http.StatusInternalServerError)
		return
	}

	var book database.Book
	if err := tx.Where("uuid = ? AND user_id = ?", note.BookUUID, user.ID).First(&book).Error; err != nil {
		tx.Rollback()
		handlers.DoError(w, fmt.Sprintf("finding book %s to preload", note.BookUUID), err, http.StatusInternalServerError)
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
	handlers.RespondJSON(w, http.StatusOK, resp)
}

type deleteNoteResp struct {
	Status int             `json:"status"`
	Result presenters.Note `json:"result"`
}

// DeleteNote removes note
func (a *API) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handlers.DoError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var note database.Note
	if err := a.App.DB.Where("uuid = ? AND user_id = ?", noteUUID, user.ID).Preload("Book").First(&note).Error; err != nil {
		handlers.DoError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	tx := a.App.DB.Begin()

	n, err := a.App.DeleteNote(tx, user, note)
	if err != nil {
		tx.Rollback()
		handlers.DoError(w, "deleting note", err, http.StatusInternalServerError)
		return
	}

	tx.Commit()

	resp := deleteNoteResp{
		Status: http.StatusNoContent,
		Result: presenters.PresentNote(n),
	}
	handlers.RespondJSON(w, http.StatusOK, resp)
}

type createNotePayload struct {
	BookUUID string `json:"book_uuid"`
	Content  string `json:"content"`
	AddedOn  *int64 `json:"added_on"`
	EditedOn *int64 `json:"edited_on"`
}

func validateCreateNotePayload(p createNotePayload) error {
	if p.BookUUID == "" {
		return errors.New("bookUUID is required")
	}

	return nil
}

// CreateNoteResp is a response for creating a note
type CreateNoteResp struct {
	Result presenters.Note `json:"result"`
}

// CreateNote creates a note
func (a *API) CreateNote(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handlers.DoError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var params createNotePayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		handlers.DoError(w, "decoding payload", err, http.StatusInternalServerError)
		return
	}

	err = validateCreateNotePayload(params)
	if err != nil {
		handlers.DoError(w, "validating payload", err, http.StatusBadRequest)
		return
	}

	var book database.Book
	if err := a.App.DB.Where("uuid = ? AND user_id = ?", params.BookUUID, user.ID).First(&book).Error; err != nil {
		handlers.DoError(w, "finding book", err, http.StatusInternalServerError)
		return
	}

	client := getClientType(r)
	note, err := a.App.CreateNote(user, params.BookUUID, params.Content, params.AddedOn, params.EditedOn, false, client)
	if err != nil {
		handlers.DoError(w, "creating note", err, http.StatusInternalServerError)
		return
	}

	// preload associations
	note.User = user
	note.Book = book

	resp := CreateNoteResp{
		Result: presenters.PresentNote(note),
	}
	handlers.RespondJSON(w, http.StatusCreated, resp)
}

// NotesOptions is a handler for OPTIONS endpoint for notes
func (a *API) NotesOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Version")
}
