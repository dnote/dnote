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

package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
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

	d, err := a.App.GetUserDigestByUUID(user.ID, digestUUID)
	if d == nil {
		RespondNotFound(w)
		return
	}
	if err != nil {
		HandleError(w, "finding digest", err, http.StatusInternalServerError)
		return
	}

	digest, err := a.App.PreloadDigest(*d)
	if err != nil {
		HandleError(w, "finding digest", err, http.StatusInternalServerError)
		return
	}

	// mark as read
	if _, err := a.App.MarkDigestRead(digest, user); err != nil {
		log.ErrorWrap(err, fmt.Sprintf("marking digest as read for %s", digest.UUID))
	}

	presented := presenters.PresentDigest(digest)
	respondJSON(w, http.StatusOK, presented)
}

// DigestsResponse is a response for getting digests
type DigestsResponse struct {
	Total int                 `json:"total"`
	Items []presenters.Digest `json:"items"`
}

type getDigestsParams struct {
	page   int
	status string
}

func parseGetDigestsParams(r *http.Request) (getDigestsParams, error) {
	var page int
	var err error

	q := r.URL.Query()

	pageStr := q.Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return getDigestsParams{}, errors.Wrap(err, "parsing page")
		}
	} else {
		page = 1
	}

	status := q.Get("status")

	return getDigestsParams{
		page:   page,
		status: status,
	}, nil
}

func (a *API) getDigests(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	params, err := parseGetDigestsParams(r)
	if err != nil {
		HandleError(w, "parsing params", err, http.StatusBadRequest)
		return
	}

	perPage := 30
	offset := (params.page - 1) * perPage
	p := app.GetDigestsParam{
		UserID:  user.ID,
		Offset:  offset,
		PerPage: perPage,
		Status:  params.status,
		Order:   "created_at DESC",
	}

	digests, err := a.App.GetDigests(p)
	if err != nil {
		HandleError(w, "querying digests", err, http.StatusInternalServerError)
		return
	}

	total, err := a.App.CountDigests(p)
	if err != nil {
		HandleError(w, "counting digests", err, http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, DigestsResponse{
		Total: total,
		Items: presenters.PresentDigests(digests),
	})
}
