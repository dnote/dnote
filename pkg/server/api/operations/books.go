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

// CreateBook creates a book with the next usn and updates the user's max_usn
func CreateBook(user database.User, clock clock.Clock, name string) (database.Book, error) {
	db := database.DBConn
	tx := db.Begin()

	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return database.Book{}, errors.Wrap(err, "incrementing user max_usn")
	}

	book := database.Book{
		UUID:      helpers.GenUUID(),
		UserID:    user.ID,
		Label:     name,
		AddedOn:   clock.Now().UnixNano(),
		USN:       nextUSN,
		Encrypted: false,
	}
	if err := tx.Create(&book).Error; err != nil {
		tx.Rollback()
		return book, errors.Wrap(err, "inserting book")
	}

	tx.Commit()

	return book, nil
}

// DeleteBook marks a book deleted with the next usn and updates the user's max_usn
func DeleteBook(tx *gorm.DB, user database.User, book database.Book) (database.Book, error) {
	if user.ID != book.UserID {
		return book, errors.New("Not allowed")
	}

	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		return book, errors.Wrap(err, "incrementing user max_usn")
	}

	if err := tx.Model(&book).
		Update(map[string]interface{}{
			"usn":     nextUSN,
			"deleted": true,
			"label":   "",
		}).Error; err != nil {
		return book, errors.Wrap(err, "deleting book")
	}

	return book, nil
}

// UpdateBook updaates the book, the usn and the user's max_usn
func UpdateBook(tx *gorm.DB, c clock.Clock, user database.User, book database.Book, label *string) (database.Book, error) {
	if user.ID != book.UserID {
		return book, errors.New("Not allowed")
	}

	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		return book, errors.Wrap(err, "incrementing user max_usn")
	}

	if label != nil {
		book.Label = *label
	}

	book.USN = nextUSN
	book.EditedOn = c.Now().UnixNano()
	book.Deleted = false
	// TODO: remove after all users have been migrated
	book.Encrypted = false

	if err := tx.Save(&book).Error; err != nil {
		return book, errors.Wrap(err, "updating the book")
	}

	return book, nil
}
