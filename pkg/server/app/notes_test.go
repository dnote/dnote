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
	"fmt"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestCreateNote(t *testing.T) {
	serverTime := time.Date(2017, time.March, 14, 21, 15, 0, 0, time.UTC)
	mockClock := clock.NewMock()
	mockClock.SetNow(serverTime)

	ts1 := time.Date(2018, time.November, 12, 10, 11, 0, 0, time.UTC).UnixNano()
	ts2 := time.Date(2018, time.November, 15, 0, 1, 10, 0, time.UTC).UnixNano()

	testCases := []struct {
		userUSN          int
		addedOn          *int64
		editedOn         *int64
		expectedUSN      int
		expectedAddedOn  int64
		expectedEditedOn int64
	}{
		{
			userUSN:          8,
			addedOn:          nil,
			editedOn:         nil,
			expectedUSN:      9,
			expectedAddedOn:  serverTime.UnixNano(),
			expectedEditedOn: 0,
		},
		{
			userUSN:          102229,
			addedOn:          &ts1,
			editedOn:         nil,
			expectedUSN:      102230,
			expectedAddedOn:  ts1,
			expectedEditedOn: 0,
		},
		{
			userUSN:          8099,
			addedOn:          &ts1,
			editedOn:         &ts2,
			expectedUSN:      8100,
			expectedAddedOn:  ts1,
			expectedEditedOn: ts2,
		},
	}

	for idx, tc := range testCases {
		func() {
			defer testutils.ClearData(testutils.DB)

			user := testutils.SetupUserData()
			testutils.MustExec(t, testutils.DB.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			anotherUser := testutils.SetupUserData()
			testutils.MustExec(t, testutils.DB.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			b1 := database.Book{UserID: user.ID, Label: "js", Deleted: false}
			testutils.MustExec(t, testutils.DB.Save(&b1), fmt.Sprintf("preparing b1 for test case %d", idx))

			a := NewTest(&App{
				Clock: mockClock,
			})

			tx := testutils.DB.Begin()
			if _, err := a.CreateNote(user, b1.UUID, "note content", tc.addedOn, tc.editedOn, false, ""); err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "deleting note"))
			}
			tx.Commit()

			var bookCount, noteCount int
			var noteRecord database.Note
			var userRecord database.User

			testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), fmt.Sprintf("counting book for test case %d", idx))
			testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), fmt.Sprintf("counting notes for test case %d", idx))
			testutils.MustExec(t, testutils.DB.First(&noteRecord), fmt.Sprintf("finding note for test case %d", idx))
			testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), fmt.Sprintf("finding user for test case %d", idx))

			assert.Equal(t, bookCount, 1, "book count mismatch")
			assert.Equal(t, noteCount, 1, "note count mismatch")
			assert.NotEqual(t, noteRecord.UUID, "", "note UUID should have been generated")
			assert.Equal(t, noteRecord.UserID, user.ID, "note UserID mismatch")
			assert.Equal(t, noteRecord.Body, "note content", "note Body mismatch")
			assert.Equal(t, noteRecord.Deleted, false, "note Deleted mismatch")
			assert.Equal(t, noteRecord.USN, tc.expectedUSN, "note Label mismatch")
			assert.Equal(t, noteRecord.AddedOn, tc.expectedAddedOn, "note AddedOn mismatch")
			assert.Equal(t, noteRecord.EditedOn, tc.expectedEditedOn, "note EditedOn mismatch")

			assert.Equal(t, userRecord.MaxUSN, tc.expectedUSN, "user max_usn mismatch")
		}()
	}
}

func TestUpdateNote(t *testing.T) {
	testCases := []struct {
		userUSN int
	}{
		{
			userUSN: 8,
		},
		{
			userUSN: 102229,
		},
		{
			userUSN: 8099,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			defer testutils.ClearData(testutils.DB)

			user := testutils.SetupUserData()
			testutils.MustExec(t, testutils.DB.Model(&user).Update("max_usn", tc.userUSN), "preparing user max_usn for test case")

			anotherUser := testutils.SetupUserData()
			testutils.MustExec(t, testutils.DB.Model(&anotherUser).Update("max_usn", 55), "preparing user max_usn for test case")

			b1 := database.Book{UserID: user.ID, Label: "js", Deleted: false}
			testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1 for test case")

			note := database.Note{UserID: user.ID, Deleted: false, Body: "test content", BookUUID: b1.UUID}
			testutils.MustExec(t, testutils.DB.Save(&note), "preparing note for test case")

			c := clock.NewMock()
			content := "updated test content"
			public := true

			a := NewTest(&App{
				Clock: c,
			})

			tx := testutils.DB.Begin()
			if _, err := a.UpdateNote(tx, user, note, &UpdateNoteParams{
				Content: &content,
				Public:  &public,
			}); err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "deleting note"))
			}
			tx.Commit()

			var bookCount, noteCount int
			var noteRecord database.Note
			var userRecord database.User

			testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), "counting book for test case")
			testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), "counting notes for test case")
			testutils.MustExec(t, testutils.DB.First(&noteRecord), "finding note for test case")
			testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), "finding user for test case")

			expectedUSN := tc.userUSN + 1
			assert.Equal(t, bookCount, 1, "book count mismatch")
			assert.Equal(t, noteCount, 1, "note count mismatch")
			assert.Equal(t, noteRecord.UserID, user.ID, "note UserID mismatch")
			assert.Equal(t, noteRecord.Body, content, "note Body mismatch")
			assert.Equal(t, noteRecord.Public, public, "note Public mismatch")
			assert.Equal(t, noteRecord.Deleted, false, "note Deleted mismatch")
			assert.Equal(t, noteRecord.USN, expectedUSN, "note USN mismatch")
			assert.Equal(t, userRecord.MaxUSN, expectedUSN, "user MaxUSN mismatch")
		})
	}
}

func TestDeleteNote(t *testing.T) {
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
			defer testutils.ClearData(testutils.DB)

			user := testutils.SetupUserData()
			testutils.MustExec(t, testutils.DB.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			anotherUser := testutils.SetupUserData()
			testutils.MustExec(t, testutils.DB.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			b1 := database.Book{UserID: user.ID, Label: "testBook"}
			testutils.MustExec(t, testutils.DB.Save(&b1), fmt.Sprintf("preparing b1 for test case %d", idx))

			note := database.Note{UserID: user.ID, Deleted: false, Body: "test content", BookUUID: b1.UUID}
			testutils.MustExec(t, testutils.DB.Save(&note), fmt.Sprintf("preparing note for test case %d", idx))

			a := NewTest(nil)

			tx := testutils.DB.Begin()
			ret, err := a.DeleteNote(tx, user, note)
			if err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "deleting note"))
			}
			tx.Commit()

			var noteCount int
			var noteRecord database.Note
			var userRecord database.User

			testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), fmt.Sprintf("counting notes for test case %d", idx))
			testutils.MustExec(t, testutils.DB.First(&noteRecord), fmt.Sprintf("finding note for test case %d", idx))
			testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), fmt.Sprintf("finding user for test case %d", idx))

			assert.Equal(t, noteCount, 1, "note count mismatch")

			assert.Equal(t, noteRecord.UserID, user.ID, "note user_id mismatch")
			assert.Equal(t, noteRecord.Body, "", "note content mismatch")
			assert.Equal(t, noteRecord.Deleted, true, "note deleted flag mismatch")
			assert.Equal(t, noteRecord.USN, tc.expectedUSN, "note label mismatch")
			assert.Equal(t, userRecord.MaxUSN, tc.expectedUSN, "user max_usn mismatch")

			assert.Equal(t, ret.UserID, user.ID, "note user_id mismatch")
			assert.Equal(t, ret.Body, "", "note content mismatch")
			assert.Equal(t, ret.Deleted, true, "note deleted flag mismatch")
			assert.Equal(t, ret.USN, tc.expectedUSN, "note label mismatch")
		}()
	}
}
