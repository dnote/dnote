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
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func init() {
	testutils.InitTestDB()
}

func TestCreateBook(t *testing.T) {
	testCases := []struct {
		userUSN     int
		expectedUSN int
		label       string
	}{
		{
			userUSN:     0,
			expectedUSN: 1,
			label:       "js",
		},
		{
			userUSN:     3,
			expectedUSN: 4,
			label:       "js",
		},
		{
			userUSN:     15,
			expectedUSN: 16,
			label:       "css",
		},
	}

	for idx, tc := range testCases {
		func() {
			defer testutils.ClearData()
			db := database.DBConn

			user := testutils.SetupUserData()
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			anotherUser := testutils.SetupUserData()
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			c := clock.NewMock()

			book, err := CreateBook(user, c, tc.label)
			if err != nil {
				t.Fatal(errors.Wrap(err, "creating book"))
			}

			var bookCount int
			var bookRecord database.Book
			var userRecord database.User

			if err := db.Model(&database.Book{}).Count(&bookCount).Error; err != nil {
				t.Fatal(errors.Wrap(err, "counting books"))
			}
			if err := db.First(&bookRecord).Error; err != nil {
				t.Fatal(errors.Wrap(err, "finding book"))
			}
			if err := db.Where("id = ?", user.ID).First(&userRecord).Error; err != nil {
				t.Fatal(errors.Wrap(err, "finding user"))
			}

			assert.Equal(t, bookCount, 1, "book count mismatch")
			assert.Equal(t, bookRecord.UserID, user.ID, "book user_id mismatch")
			assert.Equal(t, bookRecord.Label, tc.label, "book label mismatch")
			assert.Equal(t, bookRecord.USN, tc.expectedUSN, "book label mismatch")

			assert.NotEqual(t, book.UUID, "", "book uuid should have been generated")
			assert.Equal(t, book.UserID, user.ID, "returned book user_id mismatch")
			assert.Equal(t, book.Label, tc.label, "returned book label mismatch")
			assert.Equal(t, book.USN, tc.expectedUSN, "returned book usn mismatch")
			assert.Equal(t, userRecord.MaxUSN, tc.expectedUSN, "user max_usn mismatch")
		}()
	}
}

func TestDeleteBook(t *testing.T) {
	testCases := []struct {
		userUSN     int
		expectedUSN int
	}{
		{
			userUSN:     3,
			expectedUSN: 4,
		},
		{
			userUSN:     9787,
			expectedUSN: 9788,
		},
		{
			userUSN:     787,
			expectedUSN: 788,
		},
	}

	for idx, tc := range testCases {
		func() {
			defer testutils.ClearData()
			db := database.DBConn

			user := testutils.SetupUserData()
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			anotherUser := testutils.SetupUserData()
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			book := database.Book{UserID: user.ID, Label: "js", Deleted: false}
			testutils.MustExec(t, db.Save(&book), fmt.Sprintf("preparing book for test case %d", idx))

			tx := db.Begin()
			ret, err := DeleteBook(tx, user, book)
			if err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "deleting book"))
			}
			tx.Commit()

			var bookCount int
			var bookRecord database.Book
			var userRecord database.User

			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), fmt.Sprintf("counting books for test case %d", idx))
			testutils.MustExec(t, db.First(&bookRecord), fmt.Sprintf("finding book for test case %d", idx))
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), fmt.Sprintf("finding user for test case %d", idx))

			assert.Equal(t, bookCount, 1, "book count mismatch")
			assert.Equal(t, bookRecord.UserID, user.ID, "book user_id mismatch")
			assert.Equal(t, bookRecord.Label, "", "book label mismatch")
			assert.Equal(t, bookRecord.Deleted, true, "book deleted flag mismatch")
			assert.Equal(t, bookRecord.USN, tc.expectedUSN, "book label mismatch")

			assert.Equal(t, ret.UserID, user.ID, "returned book user_id mismatch")
			assert.Equal(t, ret.Label, "", "returned book label mismatch")
			assert.Equal(t, ret.Deleted, true, "returned book deleted flag mismatch")
			assert.Equal(t, ret.USN, tc.expectedUSN, "returned book label mismatch")

			assert.Equal(t, userRecord.MaxUSN, tc.expectedUSN, "user max_usn mismatch")
		}()
	}
}

func TestUpdateBook(t *testing.T) {
	js := "js"

	testCases := []struct {
		usn             int
		userUSN         int
		label           string
		payloadLabel    *string
		expectedUSN     int
		expectedUserUSN int
		expectedLabel   string
	}{
		{
			userUSN:         1,
			usn:             1,
			label:           "js",
			payloadLabel:    nil,
			expectedUSN:     2,
			expectedUserUSN: 2,
			expectedLabel:   "js",
		},
		{
			userUSN:         8,
			usn:             3,
			label:           "css",
			payloadLabel:    &js,
			expectedUSN:     9,
			expectedUserUSN: 9,
			expectedLabel:   "js",
		},
	}

	for idx, tc := range testCases {
		func() {
			defer testutils.ClearData()
			db := database.DBConn

			user := testutils.SetupUserData()
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			anotherUser := testutils.SetupUserData()
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			c := clock.NewMock()

			b := database.Book{UserID: user.ID, Deleted: false, Label: tc.expectedLabel}
			testutils.MustExec(t, db.Save(&b), fmt.Sprintf("preparing book for test case %d", idx))

			tx := db.Begin()

			book, err := UpdateBook(tx, c, user, b, tc.payloadLabel)
			if err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "updating book"))
			}

			tx.Commit()

			var bookCount int
			var bookRecord database.Book
			var userRecord database.User
			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), fmt.Sprintf("counting books for test case %d", idx))
			testutils.MustExec(t, db.First(&bookRecord), fmt.Sprintf("finding book for test case %d", idx))
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), fmt.Sprintf("finding user for test case %d", idx))

			assert.Equal(t, bookCount, 1, "book count mismatch")

			assert.Equal(t, bookRecord.UserID, user.ID, "book user_id mismatch")
			assert.Equal(t, bookRecord.Label, tc.expectedLabel, "book label mismatch")
			assert.Equal(t, bookRecord.USN, tc.expectedUSN, "book label mismatch")
			assert.Equal(t, bookRecord.EditedOn, c.Now().UnixNano(), "book edited_on mismatch")
			assert.Equal(t, book.UserID, user.ID, "returned book user_id mismatch")
			assert.Equal(t, book.Label, tc.expectedLabel, "returned book label mismatch")
			assert.Equal(t, book.USN, tc.expectedUSN, "returned book usn mismatch")
			assert.Equal(t, book.EditedOn, c.Now().UnixNano(), "returned book edited_on mismatch")

			assert.Equal(t, userRecord.MaxUSN, tc.expectedUserUSN, "user max_usn mismatch")
		}()
	}
}
