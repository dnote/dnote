/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

// import (
// 	"fmt"
// 	"net/http"
// 	"testing"
//
// 	"github.com/dnote/dnote/pkg/assert"
// 	"github.com/dnote/dnote/pkg/clock"
// 	"github.com/dnote/dnote/pkg/server/app"
// 	"github.com/dnote/dnote/pkg/server/database"
// 	"github.com/dnote/dnote/pkg/server/testutils"
// )
//
// func TestCreateNote(t *testing.T) {
//
// 	defer testutils.ClearData(testutils.DB)
//
// 	// Setup
// 	server := MustNewServer(t, &app.App{
//
// 		Clock: clock.NewMock(),
// 	})
// 	defer server.Close()
//
// 	user := testutils.SetupUserData()
// 	testutils.MustExec(t, testutils.DB.Model(&user).Update("max_usn", 101), "preparing user max_usn")
//
// 	b1 := database.Book{
// 		UserID: user.ID,
// 		Label:  "js",
// 		USN:    58,
// 	}
// 	testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
//
// 	// Execute
// 	dat := fmt.Sprintf(`{"book_uuid": "%s", "content": "note content"}`, b1.UUID)
// 	req := testutils.MakeReq(server.URL, "POST", "/v3/notes", dat)
// 	res := testutils.HTTPAuthDo(t, req, user)
//
// 	// Test
// 	assert.StatusCodeEquals(t, res, http.StatusCreated, "")
//
// 	var noteRecord database.Note
// 	var bookRecord database.Book
// 	var userRecord database.User
// 	var bookCount, noteCount int
// 	testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), "counting books")
// 	testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), "counting notes")
// 	testutils.MustExec(t, testutils.DB.First(&noteRecord), "finding note")
// 	testutils.MustExec(t, testutils.DB.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
// 	testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), "finding user record")
//
// 	assert.Equalf(t, bookCount, 1, "book count mismatch")
// 	assert.Equalf(t, noteCount, 1, "note count mismatch")
//
// 	assert.Equal(t, bookRecord.Label, b1.Label, "book name mismatch")
// 	assert.Equal(t, bookRecord.UUID, b1.UUID, "book uuid mismatch")
// 	assert.Equal(t, bookRecord.UserID, b1.UserID, "book user_id mismatch")
// 	assert.Equal(t, bookRecord.USN, 58, "book usn mismatch")
//
// 	assert.NotEqual(t, noteRecord.UUID, "", "note uuid should have been generated")
// 	assert.Equal(t, noteRecord.BookUUID, b1.UUID, "note book_uuid mismatch")
// 	assert.Equal(t, noteRecord.Body, "note content", "note content mismatch")
// 	assert.Equal(t, noteRecord.USN, 102, "note usn mismatch")
// }
//
// func TestDeleteNote(t *testing.T) {
// 	b1UUID := "37868a8e-a844-4265-9a4f-0be598084733"
//
// 	testCases := []struct {
// 		content        string
// 		deleted        bool
// 		originalUSN    int
// 		expectedUSN    int
// 		expectedMaxUSN int
// 	}{
// 		{
// 			content:        "n1 content",
// 			deleted:        false,
// 			originalUSN:    12,
// 			expectedUSN:    982,
// 			expectedMaxUSN: 982,
// 		},
// 		{
// 			content:        "",
// 			deleted:        true,
// 			originalUSN:    12,
// 			expectedUSN:    982,
// 			expectedMaxUSN: 982,
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(fmt.Sprintf("originally deleted %t", tc.deleted), func(t *testing.T) {
//
// 			defer testutils.ClearData(testutils.DB)
//
// 			// Setup
// 			server := MustNewServer(t, &app.App{
//
// 				Clock: clock.NewMock(),
// 			})
// 			defer server.Close()
//
// 			user := testutils.SetupUserData()
// 			testutils.MustExec(t, testutils.DB.Model(&user).Update("max_usn", 981), "preparing user max_usn")
//
// 			b1 := database.Book{
// 				UUID:   b1UUID,
// 				UserID: user.ID,
// 				Label:  "js",
// 			}
// 			testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
// 			note := database.Note{
// 				UserID:   user.ID,
// 				BookUUID: b1.UUID,
// 				Body:     tc.content,
// 				Deleted:  tc.deleted,
// 				USN:      tc.originalUSN,
// 			}
// 			testutils.MustExec(t, testutils.DB.Save(&note), "preparing note")
//
// 			// Execute
// 			endpoint := fmt.Sprintf("/v3/notes/%s", note.UUID)
// 			req := testutils.MakeReq(server.URL, "DELETE", endpoint, "")
// 			res := testutils.HTTPAuthDo(t, req, user)
//
// 			// Test
// 			assert.StatusCodeEquals(t, res, http.StatusOK, "")
//
// 			var bookRecord database.Book
// 			var noteRecord database.Note
// 			var userRecord database.User
// 			var bookCount, noteCount int
// 			testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), "counting books")
// 			testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), "counting notes")
// 			testutils.MustExec(t, testutils.DB.Where("uuid = ?", note.UUID).First(&noteRecord), "finding note")
// 			testutils.MustExec(t, testutils.DB.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
// 			testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), "finding user record")
//
// 			assert.Equalf(t, bookCount, 1, "book count mismatch")
// 			assert.Equalf(t, noteCount, 1, "note count mismatch")
//
// 			assert.Equal(t, noteRecord.UUID, note.UUID, "note uuid mismatch for test case")
// 			assert.Equal(t, noteRecord.Body, "", "note content mismatch for test case")
// 			assert.Equal(t, noteRecord.Deleted, true, "note deleted mismatch for test case")
// 			assert.Equal(t, noteRecord.BookUUID, note.BookUUID, "note book_uuid mismatch for test case")
// 			assert.Equal(t, noteRecord.UserID, note.UserID, "note user_id mismatch for test case")
// 			assert.Equal(t, noteRecord.USN, tc.expectedUSN, "note usn mismatch for test case")
//
// 			assert.Equal(t, userRecord.MaxUSN, tc.expectedMaxUSN, "user max_usn mismatch for test case")
// 		})
// 	}
// }
