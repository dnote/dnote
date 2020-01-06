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
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type createNoteReviewParams struct {
	DigestUUID string `json:"digest_uuid"`
	NoteUUID   string `json:"note_uuid"`
}

func getDigestByUUID(db *gorm.DB, uuid string) (*database.Digest, error) {
	var ret database.Digest
	conn := db.Where("uuid = ?", uuid).First(&ret)

	if conn.RecordNotFound() {
		return nil, nil
	}
	if err := conn.Error; err != nil {
		return nil, errors.Wrap(err, "finding digest")
	}

	return &ret, nil
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

	digest, err := a.App.GetUserDigestByUUID(user.ID, params.DigestUUID)
	if digest == nil {
		http.Error(w, "digest not found for the given uuid", http.StatusBadRequest)
		return
	}
	if err != nil {
		HandleError(w, "finding digest", err, http.StatusInternalServerError)
		return
	}

	note, err := a.App.GetUserNoteByUUID(user.ID, params.NoteUUID)
	if note == nil {
		http.Error(w, "note not found for the given uuid", http.StatusBadRequest)
		return
	}
	if err != nil {
		HandleError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	var nr database.NoteReview
	if err := a.App.DB.Debug().FirstOrCreate(&nr, database.NoteReview{
		UserID:   user.ID,
		DigestID: digest.ID,
		NoteID:   note.ID,
	}).Error; err != nil {
		HandleError(w, "saving note review", err, http.StatusInternalServerError)
		return
	}
}

type deleteNoteReviewParams struct {
	DigestUUID string `json:"digest_uuid"`
	NoteUUID   string `json:"note_uuid"`
}

func (a *API) deleteNoteReview(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var params deleteNoteReviewParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		HandleError(w, "decoding params", err, http.StatusInternalServerError)
		return
	}

	db := a.App.DB

	note, err := a.App.GetUserNoteByUUID(user.ID, params.NoteUUID)
	if note == nil {
		http.Error(w, "note not found for the given uuid", http.StatusBadRequest)
		return
	}
	if err != nil {
		HandleError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	digest, err := a.App.GetUserDigestByUUID(user.ID, params.DigestUUID)
	if digest == nil {
		http.Error(w, "digest not found for the given uuid", http.StatusBadRequest)
		return
	}
	if err != nil {
		HandleError(w, "finding digest", err, http.StatusInternalServerError)
		return
	}

	var nr database.NoteReview
	conn := db.Where("note_id = ? AND digest_id = ? AND user_id = ?", note.ID, digest.ID, user.ID).First(&nr)
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
