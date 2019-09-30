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

package operations

import (
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// CreateNote creates a note with the next usn and updates the user's max_usn.
// It returns the created note.
func CreateNote(user database.User, clock clock.Clock, bookUUID, content string, addedOn *int64, editedOn *int64, public bool) (database.Note, error) {
	db := database.DBConn
	tx := db.Begin()

	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return database.Note{}, errors.Wrap(err, "incrementing user max_usn")
	}

	var noteAddedOn int64
	if addedOn == nil {
		noteAddedOn = clock.Now().UnixNano()
	} else {
		noteAddedOn = *addedOn
	}

	var noteEditedOn int64
	if editedOn == nil {
		noteEditedOn = 0
	} else {
		noteEditedOn = *editedOn
	}

	note := database.Note{
		UUID:      helpers.GenUUID(),
		BookUUID:  bookUUID,
		UserID:    user.ID,
		AddedOn:   noteAddedOn,
		EditedOn:  noteEditedOn,
		USN:       nextUSN,
		Body:      content,
		Public:    public,
		Encrypted: false,
	}
	if err := tx.Create(&note).Error; err != nil {
		tx.Rollback()
		return note, errors.Wrap(err, "inserting note")
	}

	tx.Commit()

	return note, nil
}

// UpdateNote creates a note with the next usn and updates the user's max_usn
func UpdateNote(tx *gorm.DB, user database.User, clock clock.Clock, note database.Note, bookUUID, content *string) (database.Note, error) {
	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		return note, errors.Wrap(err, "incrementing user max_usn")
	}

	if bookUUID != nil {
		note.BookUUID = *bookUUID
	}
	if content != nil {
		note.Body = *content
	}

	note.USN = nextUSN
	note.EditedOn = clock.Now().UnixNano()
	note.Deleted = false
	// TODO: remove after all users are migrated
	note.Encrypted = false

	if err := tx.Save(&note).Error; err != nil {
		return note, errors.Wrap(err, "editing note")
	}

	return note, nil
}

// DeleteNote marks a note deleted with the next usn and updates the user's max_usn
func DeleteNote(tx *gorm.DB, user database.User, note database.Note) (database.Note, error) {
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
