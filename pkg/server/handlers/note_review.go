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
	"net/http"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/gorilla/mux"
	// "github.com/dnote/dnote/pkg/server/operations"
	// 	"github.com/dnote/dnote/pkg/server/presenters"
	// "github.com/jinzhu/gorm"
	// "github.com/pkg/errors"
)

type createNoteReviewParams struct {
	DigestUUID string `json:"digest_uuid"`
	NoteUUID   string `json:"note_uuid"`
}

func (a *API) createNoteReview(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var params createNoteReviewParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, "decoding params", err, http.StatusInternalServerError)
		return
	}

	db := a.App.DB

	var digest database.Digest
	conn := db.Where("uuid = ? AND user_id = ?", params.DigestUUID, user.ID).First(&digest)
	if conn.RecordNotFound() {
		http.Error(w, "digest not found for the given uuid", http.StatusBadRequest)
		return
	} else if err := conn.Error; err != nil {
		HandleError(w, "finding digest", err, http.StatusInternalServerError)
		return
	}

	var note database.Note
	conn2 := db.Where("uuid = ? AND user_id = ?", params.NoteUUID, user.ID).First(&note)
	if conn2.RecordNotFound() {
		http.Error(w, "note not found for the given uuid", http.StatusBadRequest)
		return
	} else if err := conn.Error; err != nil {
		HandleError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	var nr database.NoteReview
	if err := db.FirstOrCreate(&nr, database.NoteReview{
		UserID:   user.ID,
		DigestID: digest.ID,
		NoteID:   note.ID,
	}).Error; err != nil {
		HandleError(w, "saving note review", err, http.StatusInternalServerError)
		return
	}
}

func (a *API) deleteNoteReview(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	noteReviewUUID := vars["noteReviewUUID"]

	db := a.App.DB

	var note database.Note
	conn2 := db.Where("uuid = ? AND user_id = ?", noteReviewUUID, user.ID).First(&note)
	if conn2.RecordNotFound() {
		http.Error(w, "note not found for the given uuid", http.StatusBadRequest)
		return
	} else if err := conn2.Error; err != nil {
		HandleError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	var nr database.NoteReview
	conn := db.Where("note_id = ? AND user_id = ?", note.ID, user.ID).First(&nr)
	if conn.RecordNotFound() {
		http.Error(w, "no record found", http.StatusBadRequest)
		return
	} else if err := conn.Error; err != nil {
		HandleError(w, "finding record", err, http.StatusInternalServerError)
		return
	}

	if err := db.Delete(&nr).Error; err != nil {
		HandleError(w, "deleting record", err, http.StatusInternalServerError)
		return
	}
}
