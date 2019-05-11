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

func respondWithDigest(w http.ResponseWriter, userID int, digestUUID string) {
	db := database.DBConn

	var digest database.Digest
	conn := db.Preload("Notes.Book").Where("user_id = ? AND uuid = ? ", userID, digestUUID).First(&digest)
	if conn.RecordNotFound() {
		http.Error(w, "finding digest", http.StatusNotFound)
		return
	} else if err := conn.Error; err != nil {
		logger.Err("finding digest %s", err.Error())
		http.Error(w, "finding digest", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	presented := presenters.PresentDigest(digest)
	if err := json.NewEncoder(w).Encode(presented); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a App) getDigest(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	digestUUID := vars["digestUUID"]

	respondWithDigest(w, user.ID, digestUUID)
}

func (a App) getDemoDigest(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetDemoUserID()
	if err != nil {
		http.Error(w, errors.Wrap(err, "finding demo user").Error(), http.StatusInternalServerError)
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
		http.Error(w, "parsing params", http.StatusBadRequest)
		return
	}
	perPage := 25
	offset := (page - 1) * perPage

	var digests []database.Digest
	conn := db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(perPage)
	if err := conn.Find(&digests).Error; err != nil {
		logger.Err("finding digests %s", err.Error())
		http.Error(w, "finding digests", http.StatusInternalServerError)
		return
	}

	var total int
	if err := db.Model(database.Digest{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		logger.Err("counting digests %s", err.Error())
		http.Error(w, "finding digests", http.StatusInternalServerError)
		return
	}

	res := DigestsResponse{
		Total:   total,
		Digests: presenters.PresentDigests(digests),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getDigests(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	respondWithDigests(w, r, user.ID)
}

func (a *App) getDemoDigests(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetDemoUserID()
	if err != nil {
		http.Error(w, errors.Wrap(err, "finding demo user").Error(), http.StatusInternalServerError)
		return
	}

	respondWithDigests(w, r, userID)
}
