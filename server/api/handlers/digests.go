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
	"strconv"

	"github.com/dnote/dnote/server/api/helpers"
	"github.com/dnote/dnote/server/api/logger"
	"github.com/dnote/dnote/server/api/presenters"
	"github.com/dnote/dnote/server/database"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (a App) getDigestNotes(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	digestUUID := vars["digestUUID"]

	db := database.DBConn

	var digest database.Digest
	conn := db.Debug().Where("user_id = ? AND uuid = ? ", user.ID, digestUUID).
		First(&digest)
	if conn.RecordNotFound() {
		http.Error(w, "finding digest", http.StatusNotFound)
		return
	} else if err := conn.Error; err != nil {
		logger.Err("finding digest %s", err.Error())
		http.Error(w, "finding digest", http.StatusInternalServerError)
		return
	}

	var notes []database.Note
	conn2 := db.Model(&database.Note{}).
		Joins("INNER JOIN digest_notes ON digest_notes.note_uuid = notes.uuid").
		Where("digest_notes.digest_uuid = ?", digest.UUID).
		Preload("Book").
		Order("notes.created_at DESC").
		Find(&notes)
	if conn2.RecordNotFound() {
		http.Error(w, "finding digest", http.StatusNotFound)
		return
	} else if err := conn2.Error; err != nil {
		logger.Err("finding digest notes %s", err.Error())
		http.Error(w, "finding digest", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	presented := presenters.PresentNotes(notes)
	if err := json.NewEncoder(w).Encode(presented); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseGetDigestsParams(r *http.Request) (int, error) {
	var page int
	var err error

	q := r.URL.Query()
	pageStr := q.Get("page")

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return 0, errors.Wrap(err, "parsing page")
		}

	}

	return page, nil
}

func (a *App) getDigests(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	db := database.DBConn

	page, err := parseGetDigestsParams(r)
	if err != nil {
		http.Error(w, "parsing params", http.StatusBadRequest)
		return
	}
	offset := page - 1
	perPage := 25

	var digests []database.Digest
	conn := db.Where("user_id = ?", user.ID).Offset(offset).Limit(perPage)

	if err := conn.Find(&digests).Error; err != nil {
		logger.Err("finding digests %s", err.Error())
		http.Error(w, "finding digests", http.StatusInternalServerError)
		return
	}

	presented := presenters.PresentDigests(digests)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(presented); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
