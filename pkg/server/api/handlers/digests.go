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
	"net/http"
	"strconv"

	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/api/presenters"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func respondWithDigest(w http.ResponseWriter, userID int, digestUUID string) {
	db := database.DBConn

	var digest database.Digest
	conn := db.Preload("Notes.Book").Where("user_id = ? AND uuid = ? ", userID, digestUUID).First(&digest)
	if conn.RecordNotFound() {
		handleError(w, "digest not found", nil, http.StatusNotFound)
		return
	} else if err := conn.Error; err != nil {

		handleError(w, "finding digest", err, http.StatusInternalServerError)
		return
	}

	presented := presenters.PresentDigest(digest)
	respondJSON(w, presented)
}

func (a App) getDigest(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	digestUUID := vars["digestUUID"]

	respondWithDigest(w, user.ID, digestUUID)
}

func (a App) getDemoDigest(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetDemoUserID()
	if err != nil {
		handleError(w, "finding demo user", err, http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	digestUUID := vars["digestUUID"]

	respondWithDigest(w, userID, digestUUID)
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

// DigestsResponse is a response for getting digests
type DigestsResponse struct {
	Total   int                 `json:"total"`
	Digests []presenters.Digest `json:"digests"`
}

func respondWithDigests(w http.ResponseWriter, r *http.Request, userID int) {
	db := database.DBConn

	page, err := parseGetDigestsParams(r)
	if err != nil {
		handleError(w, "parsing params", err, http.StatusBadRequest)
		return
	}
	perPage := 25
	offset := (page - 1) * perPage

	var digests []database.Digest
	conn := db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(perPage)
	if err := conn.Find(&digests).Error; err != nil {
		handleError(w, "finding digests", err, http.StatusInternalServerError)
		return
	}

	var total int
	if err := db.Model(database.Digest{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		handleError(w, "counting digests", err, http.StatusInternalServerError)
		return
	}

	res := DigestsResponse{
		Total:   total,
		Digests: presenters.PresentDigests(digests),
	}
	respondJSON(w, res)
}

func (a *App) getDigests(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	respondWithDigests(w, r, user.ID)
}

func (a *App) getDemoDigests(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetDemoUserID()
	if err != nil {
		handleError(w, "finding demo user", err, http.StatusInternalServerError)
		return
	}

	respondWithDigests(w, r, userID)
}
