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
	"fmt"
	"net/http"
	"strconv"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func preloadDigest(conn *gorm.DB, userID int) *gorm.DB {
	return conn.
		Preload("Notes", func(db *gorm.DB) *gorm.DB {
			return db.Order("notes.created_at DESC")
		}).
		Preload("Notes.Book").
		Preload("Notes.NoteReview").
		Preload("Rule").
		Preload("Receipts", func(db *gorm.DB) *gorm.DB {
			return db.Where("digest_receipts.user_id = ?", userID)
		})
}

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
	conn := db.Where("user_id = ? AND uuid = ? ", user.ID, digestUUID)
	conn = preloadDigest(conn, user.ID)
	conn = conn.First(&digest)

	if conn.RecordNotFound() {
		RespondNotFound(w)
		return
	} else if err := conn.Error; err != nil {
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

func queryDigestIDs(db *gorm.DB, p getDigestsParams, userID, offset, perPage int) ([]int, error) {
	var whereClause string
	if p.status == "unread" {
		whereClause = "WHERE t1.receipt_count = 0"
	} else if p.status == "read" {
		whereClause = "WHERE t1.receipt_count > 0"
	}

	query := fmt.Sprintf(`
SELECT t1.digest_id FROM
(
	SELECT
		digests.id AS digest_id,
		COUNT(digest_receipts.id) AS receipt_count
	FROM digests
	LEFT JOIN digest_receipts ON digest_receipts.digest_id = digests.id
	WHERE digests.user_id = ?
	GROUP BY digests.id, digests.created_at
	ORDER BY digests.created_at DESC
) AS t1
%s
OFFSET ?
LIMIT ?;
`, whereClause)

	ret := []int{}
	rows, err := db.Debug().Raw(query, userID, offset, perPage).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "getting rows")
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return []int{}, errors.Wrap(err, "scanning row")
		}

		ret = append(ret, id)
	}

	return ret, nil

}

func (a *API) getDigests(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	db := a.App.DB

	params, err := parseGetDigestsParams(r)
	if err != nil {
		HandleError(w, "parsing params", err, http.StatusBadRequest)
		return
	}

	perPage := 30
	offset := (params.page - 1) * perPage

	IDs, err := queryDigestIDs(db, params, user.ID, offset, perPage)
	if err != nil {
		HandleError(w, "querying digest IDs", err, http.StatusInternalServerError)
		return
	}

	var digests []database.Digest
	conn := db.Debug().
		Where("id IN (?)", IDs).
		Order("created_at DESC").
		Preload("Rule").Preload("Receipts").
		Find(&digests)
	if err := conn.Error; err != nil && !conn.RecordNotFound() {
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
