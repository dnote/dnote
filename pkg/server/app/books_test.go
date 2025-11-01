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
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

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
			db := testutils.InitMemoryDB(t)

			user := testutils.SetupUserData(db, "user@test.com", "password123")
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			anotherUser := testutils.SetupUserData(db, "another@test.com", "password123")
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			a := NewTest()
			a.DB = db
			a.Clock = clock.NewMock()

			book, err := a.CreateBook(user, tc.label)
			if err != nil {
				t.Fatal(errors.Wrap(err, "creating book"))
			}

			var bookCount int64
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

			assert.Equal(t, bookCount, int64(1), "book count mismatch")
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
			db := testutils.InitMemoryDB(t)

			user := testutils.SetupUserData(db, "user@test.com", "password123")
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			anotherUser := testutils.SetupUserData(db, "another@test.com", "password123")
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			book := database.Book{UserID: user.ID, Label: "js", Deleted: false}
			testutils.MustExec(t, db.Save(&book), fmt.Sprintf("preparing book for test case %d", idx))

			tx := db.Begin()
			a := NewTest()
			a.DB = db
			ret, err := a.DeleteBook(tx, user, book)
			if err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "deleting book"))
			}
			tx.Commit()

			var bookCount int64
			var bookRecord database.Book
			var userRecord database.User

			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), fmt.Sprintf("counting books for test case %d", idx))
			testutils.MustExec(t, db.First(&bookRecord), fmt.Sprintf("finding book for test case %d", idx))
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), fmt.Sprintf("finding user for test case %d", idx))

			assert.Equal(t, bookCount, int64(1), "book count mismatch")
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
			db := testutils.InitMemoryDB(t)

			user := testutils.SetupUserData(db, "user@test.com", "password123")
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			anotherUser := testutils.SetupUserData(db, "another@test.com", "password123")
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			b := database.Book{UserID: user.ID, Deleted: false, Label: tc.expectedLabel}
			testutils.MustExec(t, db.Save(&b), fmt.Sprintf("preparing book for test case %d", idx))

			c := clock.NewMock()
			a := NewTest()
			a.DB = db
			a.Clock = c

			tx := db.Begin()
			book, err := a.UpdateBook(tx, user, b, tc.payloadLabel)
			if err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "updating book"))
			}

			tx.Commit()

			var bookCount int64
			var bookRecord database.Book
			var userRecord database.User
			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), fmt.Sprintf("counting books for test case %d", idx))
			testutils.MustExec(t, db.First(&bookRecord), fmt.Sprintf("finding book for test case %d", idx))
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), fmt.Sprintf("finding user for test case %d", idx))

			assert.Equal(t, bookCount, int64(1), "book count mismatch")

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
