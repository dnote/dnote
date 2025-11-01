/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"gorm.io/gorm"
	"github.com/pkg/errors"
)

// CreateBook creates a book with the next usn and updates the user's max_usn
func (a *App) CreateBook(user database.User, name string) (database.Book, error) {
	tx := a.DB.Begin()

	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return database.Book{}, errors.Wrap(err, "incrementing user max_usn")
	}

	uuid, err := helpers.GenUUID()
	if err != nil {
		tx.Rollback()
		return database.Book{}, err
	}

	book := database.Book{
		UUID:    uuid,
		UserID:  user.ID,
		Label:   name,
		AddedOn: a.Clock.Now().UnixNano(),
		USN:     nextUSN,
	}
	if err := tx.Create(&book).Error; err != nil {
		tx.Rollback()
		return book, errors.Wrap(err, "inserting book")
	}

	tx.Commit()

	return book, nil
}

// DeleteBook marks a book deleted with the next usn and updates the user's max_usn
func (a *App) DeleteBook(tx *gorm.DB, user database.User, book database.Book) (database.Book, error) {
	if user.ID != book.UserID {
		return book, errors.New("Not allowed")
	}

	nextUSN, err := incrementUserUSN(tx, user.ID)
	if err != nil {
		return book, errors.Wrap(err, "incrementing user max_usn")
	}

	if err := tx.Model(&book).
		Updates(map[string]interface{}{
			"usn":     nextUSN,
			"deleted": true,
			"label":   "",
		}).Error; err != nil {
		return book, errors.Wrap(err, "deleting book")
	}

	return book, nil
}

// UpdateBook updaates the book, the usn and the user's max_usn
func (a *App) UpdateBook(tx *gorm.DB, user database.User, book database.Book, label *string) (database.Book, error) {
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
	book.EditedOn = a.Clock.Now().UnixNano()
	book.Deleted = false

	if err := tx.Save(&book).Error; err != nil {
		return book, errors.Wrap(err, "updating the book")
	}

	return book, nil
}
