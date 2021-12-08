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

package app

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// CreateNote creates a note with the next usn and updates the user's max_usn.
// It returns the created note.
func (a *App) CreateNote(user database.User, bookUUID, content string, addedOn *int64, editedOn *int64, public bool, client string) (database.Note, error) {
	tx := a.DB.Begin()

	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return database.Note{}, errors.Wrap(err, "incrementing user max_usn")
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
		return database.Note{}, err
	}

	note := database.Note{
		UUID:      uuid,
		BookUUID:  bookUUID,
		UserID:    user.ID,
		AddedOn:   noteAddedOn,
		EditedOn:  noteEditedOn,
		USN:       nextUSN,
		Body:      content,
		Public:    public,
		Encrypted: false,
		Client:    client,
	}
	if err := tx.Create(&note).Error; err != nil {
		tx.Rollback()
		return note, errors.Wrap(err, "inserting note")
	}

	tx.Commit()

	return note, nil
}

// UpdateNoteParams is the parameters for updating a note
type UpdateNoteParams struct {
	BookUUID *string
	Content  *string
	Public   *bool
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

// GetPublic gets the public field from the UpdateNoteParams
func (r UpdateNoteParams) GetPublic() bool {
	if r.Public == nil {
		return false
	}

	return *r.Public
}

// UpdateNote creates a note with the next usn and updates the user's max_usn
func (a *App) UpdateNote(tx *gorm.DB, user database.User, note database.Note, p *UpdateNoteParams) (database.Note, error) {
	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		return note, errors.Wrap(err, "incrementing user max_usn")
	}

	if p.BookUUID != nil {
		note.BookUUID = p.GetBookUUID()
	}
	if p.Content != nil {
		note.Body = p.GetContent()
	}
	if p.Public != nil {
		note.Public = p.GetPublic()
	}

	note.USN = nextUSN
	note.EditedOn = a.Clock.Now().UnixNano()
	note.Deleted = false
	// TODO: remove after all users are migrated
	note.Encrypted = false

	if err := tx.Save(&note).Error; err != nil {
		return note, errors.Wrap(err, "editing note")
	}

	return note, nil
}

// DeleteNote marks a note deleted with the next usn and updates the user's max_usn
func (a *App) DeleteNote(tx *gorm.DB, user database.User, note database.Note) (database.Note, error) {
	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		return note, errors.Wrap(err, "incrementing user max_usn")
	}

	if err := tx.Model(&note).
		Update(map[string]interface{}{
			"usn":     nextUSN,
			"deleted": true,
			"body":    "",
		}).Error; err != nil {
		return note, errors.Wrap(err, "deleting note")
	}

	return note, nil
}

// GetUserNoteByUUID retrives a digest by the uuid for the given user
func (a *App) GetUserNoteByUUID(userID int, uuid string) (*database.Note, error) {
	var ret database.Note
	conn := a.DB.Where("user_id = ? AND uuid = ?", userID, uuid).First(&ret)

	if conn.RecordNotFound() {
		return nil, nil
	}
	if err := conn.Error; err != nil {
		return nil, errors.Wrap(err, "finding digest")
	}

	return &ret, nil
}
