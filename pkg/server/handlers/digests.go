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

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (a *API) getDigest(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	digestUUID := vars["digestUUID"]

	db := a.App.DB

	var digest database.Digest
	conn := db.Preload("Notes.Book").Where("user_id = ? AND uuid = ? ", user.ID, digestUUID).First(&digest)
	if conn.RecordNotFound() {
		HandleError(w, "digest not found", nil, http.StatusNotFound)
		return
	} else if err := conn.Error; err != nil {

		HandleError(w, "finding digest", err, http.StatusInternalServerError)
		return
	}

	presented := presenters.PresentDigest(digest)
	respondJSON(w, http.StatusOK, presented)
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
	Total int                 `json:"total"`
	Items []presenters.Digest `json:"digests"`
}

func (a *API) getDigests(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	db := a.App.DB

	page, err := parseGetDigestsParams(r)
	if err != nil {
		HandleError(w, "parsing params", err, http.StatusBadRequest)
		return
	}

	perPage := 25
	offset := (page - 1) * perPage

	var digests []database.Digest
	conn := db.Where("user_id = ?", user.ID).Order("created_at DESC").Offset(offset).Limit(perPage)
	if err := conn.Find(&digests).Error; err != nil {
		HandleError(w, "finding digests", err, http.StatusInternalServerError)
		return
	}

	var total int
	if err := db.Model(database.Digest{}).Where("user_id = ?", user.ID).Count(&total).Error; err != nil {
		HandleError(w, "counting digests", err, http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, DigestsResponse{
		Total: total,
		Items: presenters.PresentDigests(digests),
	})
}
