/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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

package app

import (
	"errors"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateNote creates a note with the next usn and updates the user's max_usn.
// It returns the created note.
func (a *App) CreateNote(user database.User, bookUUID, content string, addedOn *int64, editedOn *int64, client string) (database.Note, error) {
	tx := a.DB.Begin()

	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return database.Note{}, pkgErrors.Wrap(err, "incrementing user max_usn")
	}

	var noteAddedOn int64
	if addedOn == nil {
		noteAddedOn = a.Clock.Now().UnixNano()
	} else {
		noteAddedOn = *addedOn
	}

	var noteEditedOn int64
	if editedOn == nil {
		noteEditedOn = 0
	} else {
		noteEditedOn = *editedOn
	}

	uuid, err := helpers.GenUUID()
	if err != nil {
		tx.Rollback()
		return database.Note{}, err
	}

	note := database.Note{
		UUID:     uuid,
		BookUUID: bookUUID,
		UserID:   user.ID,
		AddedOn:  noteAddedOn,
		EditedOn: noteEditedOn,
		USN:      nextUSN,
		Body:     content,
		Client:   client,
	}
	if err := tx.Create(&note).Error; err != nil {
		tx.Rollback()
		return note, pkgErrors.Wrap(err, "inserting note")
	}

	tx.Commit()

	return note, nil
}

// UpdateNoteParams is the parameters for updating a note
type UpdateNoteParams struct {
	BookUUID *string
	Content  *string
}

// GetBookUUID gets the bookUUID from the UpdateNoteParams
func (r UpdateNoteParams) GetBookUUID() string {
	if r.BookUUID == nil {
		return ""
	}

	return *r.BookUUID
}

// GetContent gets the content from the UpdateNoteParams
func (r UpdateNoteParams) GetContent() string {
	if r.Content == nil {
		return ""
	}

	return *r.Content
}

// UpdateNote creates a note with the next usn and updates the user's max_usn
func (a *App) UpdateNote(tx *gorm.DB, user database.User, note database.Note, p *UpdateNoteParams) (database.Note, error) {
	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		return note, pkgErrors.Wrap(err, "incrementing user max_usn")
	}

	if p.BookUUID != nil {
		note.BookUUID = p.GetBookUUID()
	}
	if p.Content != nil {
		note.Body = p.GetContent()
	}

	note.USN = nextUSN
	note.EditedOn = a.Clock.Now().UnixNano()
	note.Deleted = false

	if err := tx.Save(&note).Error; err != nil {
		return note, pkgErrors.Wrap(err, "editing note")
	}

	return note, nil
}

// DeleteNote marks a note deleted with the next usn and updates the user's max_usn
func (a *App) DeleteNote(tx *gorm.DB, user database.User, note database.Note) (database.Note, error) {
	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		return note, pkgErrors.Wrap(err, "incrementing user max_usn")
	}

	if err := tx.Model(&note).
		Updates(map[string]interface{}{
			"usn":     nextUSN,
			"deleted": true,
			"body":    "",
		}).Error; err != nil {
		return note, pkgErrors.Wrap(err, "deleting note")
	}

	return note, nil
}

// GetUserNoteByUUID retrives a digest by the uuid for the given user
func (a *App) GetUserNoteByUUID(userID int, uuid string) (*database.Note, error) {
	var ret database.Note
	err := a.DB.Where("user_id = ? AND uuid = ?", userID, uuid).First(&ret).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, pkgErrors.Wrap(err, "finding digest")
	}

	return &ret, nil
}

// GetNotesParams is params for finding notes
type GetNotesParams struct {
	Year    int
	Month   int
	Page    int
	Books   []string
	Search  string
	PerPage int
}

type ftsParams struct {
	HighlightAll bool
}

func getFTSBodyExpression(params *ftsParams) string {
	if params != nil && params.HighlightAll {
		return "highlight(notes_fts, 0, '<dnotehl>', '</dnotehl>') AS body"
	}

	return "snippet(notes_fts, 0, '<dnotehl>', '</dnotehl>', '...', 50) AS body"
}

func selectFTSFields(conn *gorm.DB, params *ftsParams) *gorm.DB {
	bodyExpr := getFTSBodyExpression(params)

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
` + bodyExpr)
}

func getNotesBaseQuery(db *gorm.DB, userID int, q GetNotesParams) *gorm.DB {
	conn := db.Where(
		"notes.user_id = ? AND notes.deleted = ?",
		userID, false,
	)

	if q.Search != "" {
		conn = selectFTSFields(conn, nil)
		conn = conn.Joins("INNER JOIN notes_fts ON notes_fts.rowid = notes.id")
		conn = conn.Where("notes_fts MATCH ?", q.Search)
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

func orderGetNotes(conn *gorm.DB) *gorm.DB {
	return conn.Order("notes.updated_at DESC, notes.id DESC")
}

func paginate(conn *gorm.DB, page, perPage int) *gorm.DB {
	// Paginate
	if page > 0 {
		offset := perPage * (page - 1)
		conn = conn.Offset(offset)
	}

	conn = conn.Limit(perPage)

	return conn
}

// GetNotesResult is the result of getting notes
type GetNotesResult struct {
	Notes []database.Note
	Total int64
}

// GetNotes returns a list of matching notes
func (a *App) GetNotes(userID int, params GetNotesParams) (GetNotesResult, error) {
	conn := getNotesBaseQuery(a.DB, userID, params)

	var total int64
	if err := conn.Model(database.Note{}).Count(&total).Error; err != nil {
		return GetNotesResult{}, pkgErrors.Wrap(err, "counting total")
	}

	notes := []database.Note{}
	if total != 0 {
		conn = orderGetNotes(conn)
		conn = database.PreloadNote(conn)
		conn = paginate(conn, params.Page, params.PerPage)

		if err := conn.Find(&notes).Error; err != nil {
			return GetNotesResult{}, pkgErrors.Wrap(err, "finding notes")
		}
	}

	res := GetNotesResult{
		Notes: notes,
		Total: total,
	}

	return res, nil
}
