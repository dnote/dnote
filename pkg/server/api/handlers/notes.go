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
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/api/presenters"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func respondWithNote(w http.ResponseWriter, note database.Note) {
	presentedNote := presenters.PresentNote(note)

	respondJSON(w, presentedNote)
}

func preloadNote(conn *gorm.DB) *gorm.DB {
	return conn.Preload("Book").Preload("User")
}

func (a *App) getDemoNote(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn
	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	demoUserID, err := helpers.GetDemoUserID()
	if err != nil {
		handleError(w, "finding demo user", err, http.StatusInternalServerError)
		return
	}

	var note database.Note
	conn := db.Where("uuid = ? AND user_id = ?", noteUUID, demoUserID)
	conn = preloadNote(conn)
	conn.Find(&note)

	if conn.RecordNotFound() {
		handleError(w, "not found", nil, http.StatusNotFound)
		return
	} else if err := conn.Error; err != nil {
		handleError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	respondWithNote(w, note)
}

func (a *App) getNote(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	db := database.DBConn
	vars := mux.Vars(r)
	noteUUID := vars["noteUUID"]

	var note database.Note
	conn := db.Where("uuid = ? AND user_id = ?", noteUUID, user.ID)
	conn = preloadNote(conn)
	conn.Find(&note)

	if conn.RecordNotFound() {
		handleError(w, "not found", nil, http.StatusNotFound)
		return
	} else if err := conn.Error; err != nil {
		handleError(w, "finding note", err, http.StatusInternalServerError)
		return
	}

	respondWithNote(w, note)
}

/**** getNotesHandler */

// GetNotesResponse is a reponse by getNotesHandler
type GetNotesResponse struct {
	Notes    []presenters.Note `json:"notes"`
	Total    int               `json:"total"`
	PrevDate *int64            `json:"prev_date"`
}

func (a *App) getDemoNotes(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetDemoUserID()
	if err != nil {
		handleError(w, "finding demo user id", err, http.StatusInternalServerError)
		return
	}
	query := r.URL.Query()

	respondGetNotes(userID, query, w)
}

func (a *App) getNotes(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}
	query := r.URL.Query()

	respondGetNotes(user.ID, query, w)
}

func respondGetNotes(userID int, query url.Values, w http.ResponseWriter) {
	err := validateGetNotesQuery(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q, err := parseGetNotesQuery(query)
	if err != nil {
		handleError(w, "parsing query parameters", err, http.StatusBadRequest)
		return
	}

	dateLowerbound, dateUpperbound := getDateBounds(q.Year, q.Month)

	baseConn := getNotesBaseQuery(userID, q)
	conn := baseConn.Where("notes.added_on >= ? AND notes.added_on < ?", dateLowerbound, dateUpperbound)

	var total int
	err = conn.Model(database.Note{}).Count(&total).Error
	if err != nil {
		handleError(w, "counting total", err, http.StatusInternalServerError)
		return
	}

	notes := []database.Note{}
	if total != 0 {
		conn = orderGetNotes(conn)
		conn = preloadNote(conn)
		conn = paginate(conn, q.Page)

		err = conn.Find(&notes).Error
		if err != nil {
			handleError(w, "finding notes", err, http.StatusInternalServerError)
			return
		}
	}

	// peek the prev date
	var prevDateUpperbound int64
	if len(notes) > 0 {
		lastNote := notes[len(notes)-1]
		prevDateUpperbound = lastNote.AddedOn
	} else {
		prevDateUpperbound = dateLowerbound
	}

	prevDate, err := getPrevDate(baseConn, prevDateUpperbound)
	if err != nil {
		handleError(w, "getting prevDate", err, http.StatusInternalServerError)
		return
	}

	presentedNotes := presenters.PresentNotes(notes)

	response := GetNotesResponse{
		Notes:    presentedNotes,
		Total:    total,
		PrevDate: prevDate,
	}
	respondJSON(w, response)
}

func getPrevDate(baseConn *gorm.DB, dateUpperbound int64) (*int64, error) {
	var prevNote database.Note

	conn := baseConn.
		Select("notes.added_on").
		Where("notes.added_on < ?", dateUpperbound).
		Order("notes.added_on DESC")

	if conn.First(&prevNote).RecordNotFound() {
		return nil, nil
	}

	if err := conn.Error; err != nil {
		return nil, errors.Wrap(err, "querying previous note")
	}

	return &prevNote.AddedOn, nil
}

func validateGetNotesQuery(q url.Values) error {
	if q.Get("year") == "" {
		return errors.New("'year' is required")
	}
	if q.Get("month") == "" {
		return errors.New("'month' is required")
	}

	return nil
}

type getNotesQuery struct {
	Year      int
	Month     int
	Page      int
	BookUUID  string
	Encrypted *bool
}

func parseGetNotesQuery(q url.Values) (getNotesQuery, error) {
	yearStr := q.Get("year")
	monthStr := q.Get("month")
	bookStr := q.Get("book")
	pageStr := q.Get("page")
	encryptedStr := q.Get("encrypted")

	var page int
	if len(pageStr) > 0 {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return getNotesQuery{}, errors.Wrap(err, "parsing page")
		}

		page = p
	} else {
		page = 1
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return getNotesQuery{}, errors.Wrapf(err, "invalid year %s", yearStr)
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return getNotesQuery{}, errors.Wrapf(err, "invalid month %s", monthStr)
	}
	if month < 1 || month > 12 {
		return getNotesQuery{}, errors.Errorf("Invalid month %s", monthStr)
	}

	var encrypted *bool
	if strings.ToLower(encryptedStr) == "true" {
		*encrypted = true
	} else if strings.ToLower(encryptedStr) == "false" {
		*encrypted = false
	}

	ret := getNotesQuery{
		Year:      year,
		Month:     month,
		Page:      page,
		BookUUID:  bookStr,
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

func getNotesBaseQuery(userID int, q getNotesQuery) *gorm.DB {
	db := database.DBConn

	conn := db.Where("notes.user_id = ? AND notes.deleted = ?", userID, false)

	if len(q.BookUUID) > 0 {
		conn = conn.Joins("INNER JOIN books ON books.uuid = notes.book_uuid").
			Where("books.uuid = ?", q.BookUUID)
	}
	if q.Encrypted != nil {
		conn = conn.Where("notes.encrypted = ?", *q.Encrypted)
	}

	return conn
}

func orderGetNotes(conn *gorm.DB) *gorm.DB {
	return conn.Order("notes.added_on DESC, notes.id DESC")
}

func (a *App) legacyGetNotes(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var notes []database.Note
	db := database.DBConn
	if err := db.Where("user_id = ? AND encrypted = false", user.ID).Find(&notes).Error; err != nil {
		handleError(w, "finding notes", err, http.StatusInternalServerError)
		return
	}

	presented := presenters.PresentNotes(notes)
	respondJSON(w, presented)
}
