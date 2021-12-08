/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/handlers"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/operations"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type ftsParams struct {
	HighlightAll bool
}

func getHeadlineOptions(params *ftsParams) string {
	headlineOptions := []string{
		"StartSel=<dnotehl>",
		"StopSel=</dnotehl>",
		"ShortWord=0",
	}

	if params != nil && params.HighlightAll {
		headlineOptions = append(headlineOptions, "HighlightAll=true")
	} else {
		headlineOptions = append(headlineOptions, "MaxFragments=3, MaxWords=50, MinWords=10")
	}

	return strings.Join(headlineOptions, ",")
}

func selectFTSFields(conn *gorm.DB, search string, params *ftsParams) *gorm.DB {
	headlineOpts := getHeadlineOptions(params)

	return conn.Select(` 
notes.id,
notes.uuid,
notes.created_at,
notes.updated_at,
notes.book_uuid,
notes.user_id,
notes.added_on,
notes.edited_on,
notes.usn,
notes.deleted,
notes.encrypted,
ts_headline('english_nostop', notes.body, plainto_tsquery('english_nostop', ?), ?) AS body
	`, search, headlineOpts)
}

func respondWithNote(w http.ResponseWriter, note database.Note) {
	presentedNote := presenters.PresentNote(note)

	handlers.RespondJSON(w, http.StatusOK, presentedNote)
}

func parseSearchQuery(q url.Values) string {
	searchStr := q.Get("q")

	return escapeSearchQuery(searchStr)
}

func getNoteBaseQuery(db *gorm.DB, noteUUID string, search string) *gorm.DB {
	var conn *gorm.DB
	if search != "" {
		conn = selectFTSFields(db, search, &ftsParams{HighlightAll: true})
	} else {
		conn = db
	}

	conn = conn.Where("notes.uuid = ? AND deleted = ?", noteUUID, false)

	return conn
}

func (a *API) getNote(w http.ResponseWriter, r *http.Request) {
	user, _, err := handlers.AuthWithSession(a.App.DB, r, nil)
	if err != nil {
		handlers.DoError(w, "authenticating", err, http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	note, ok, err := operations.GetNote(a.App.DB, noteUUID, user)
	if !ok {
		handlers.RespondNotFound(w)
		return
	}
	if err != nil {
		handlers.DoError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	respondWithNote(w, note)
}

/**** getNotesHandler */

// GetNotesResponse is a reponse by getNotesHandler
type GetNotesResponse struct {
	Notes []presenters.Note `json:"notes"`
	Total int               `json:"total"`
}

type dateRange struct {
	lower int64
	upper int64
}

func (a *API) getNotes(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handlers.DoError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}
	query := r.URL.Query()

	respondGetNotes(a.App.DB, user.ID, query, w)
}

func respondGetNotes(db *gorm.DB, userID int, query url.Values, w http.ResponseWriter) {
	q, err := parseGetNotesQuery(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn := getNotesBaseQuery(db, userID, q)

	var total int
	if err := conn.Model(database.Note{}).Count(&total).Error; err != nil {
		handlers.DoError(w, "counting total", err, http.StatusInternalServerError)
		return
	}

	notes := []database.Note{}
	if total != 0 {
		conn = orderGetNotes(conn)
		conn = database.PreloadNote(conn)
		conn = paginate(conn, q.Page)

		if err := conn.Find(&notes).Error; err != nil {
			handlers.DoError(w, "finding notes", err, http.StatusInternalServerError)
			return
		}
	}

	response := GetNotesResponse{
		Notes: presenters.PresentNotes(notes),
		Total: total,
	}
	handlers.RespondJSON(w, http.StatusOK, response)
}

type getNotesQuery struct {
	Year      int
	Month     int
	Page      int
	Books     []string
	Search    string
	Encrypted bool
}

func parseGetNotesQuery(q url.Values) (getNotesQuery, error) {
	yearStr := q.Get("year")
	monthStr := q.Get("month")
	books := q["book"]
	pageStr := q.Get("page")
	encryptedStr := q.Get("encrypted")

	fmt.Println("books", books)

	var page int
	if len(pageStr) > 0 {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return getNotesQuery{}, errors.Errorf("invalid page %s", pageStr)
		}
		if p < 1 {
			return getNotesQuery{}, errors.Errorf("invalid page %s", pageStr)
		}

		page = p
	} else {
		page = 1
	}

	var year int
	if len(yearStr) > 0 {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			return getNotesQuery{}, errors.Errorf("invalid year %s", yearStr)
		}

		year = y
	}

	var month int
	if len(monthStr) > 0 {
		m, err := strconv.Atoi(monthStr)
		if err != nil {
			return getNotesQuery{}, errors.Errorf("invalid month %s", monthStr)
		}
		if m < 1 || m > 12 {
			return getNotesQuery{}, errors.Errorf("invalid month %s", monthStr)
		}

		month = m
	}

	var encrypted bool
	if strings.ToLower(encryptedStr) == "true" {
		encrypted = true
	} else {
		encrypted = false
	}

	ret := getNotesQuery{
		Year:      year,
		Month:     month,
		Page:      page,
		Search:    parseSearchQuery(q),
		Books:     books,
		Encrypted: encrypted,
	}

	return ret, nil
}

func getDateBounds(year, month int) (int64, int64) {
	var yearUpperbound, monthUpperbound int

	if month == 12 {
		monthUpperbound = 1
		yearUpperbound = year + 1
	} else {
		monthUpperbound = month + 1
		yearUpperbound = year
	}

	lower := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).UnixNano()
	upper := time.Date(yearUpperbound, time.Month(monthUpperbound), 1, 0, 0, 0, 0, time.UTC).UnixNano()

	return lower, upper
}

func getNotesBaseQuery(db *gorm.DB, userID int, q getNotesQuery) *gorm.DB {
	conn := db.Where(
		"notes.user_id = ? AND notes.deleted = ? AND notes.encrypted = ?",
		userID, false, q.Encrypted,
	)

	if q.Search != "" {
		conn = selectFTSFields(conn, q.Search, nil)
		conn = conn.Where("tsv @@ plainto_tsquery('english_nostop', ?)", q.Search)
	}

	if len(q.Books) > 0 {
		conn = conn.Joins("INNER JOIN books ON books.uuid = notes.book_uuid").
			Where("books.label in (?)", q.Books)
	}

	if q.Year != 0 || q.Month != 0 {
		dateLowerbound, dateUpperbound := getDateBounds(q.Year, q.Month)
		conn = conn.Where("notes.added_on >= ? AND notes.added_on < ?", dateLowerbound, dateUpperbound)
	}

	return conn
}

func orderGetNotes(conn *gorm.DB) *gorm.DB {
	return conn.Order("notes.updated_at DESC, notes.id DESC")
}

// escapeSearchQuery escapes the query for full text search
func escapeSearchQuery(searchQuery string) string {
	return strings.Join(strings.Fields(searchQuery), "&")
}

func (a *API) legacyGetNotes(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handlers.DoError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var notes []database.Note
	if err := a.App.DB.Where("user_id = ? AND encrypted = true", user.ID).Find(&notes).Error; err != nil {
		handlers.DoError(w, "finding notes", err, http.StatusInternalServerError)
		return
	}

	presented := presenters.PresentNotes(notes)
	handlers.RespondJSON(w, http.StatusOK, presented)
}
