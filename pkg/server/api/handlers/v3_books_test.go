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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/api/presenters"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func init() {
	testutils.InitTestDB()
}

func TestGetBooks(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	user := testutils.SetupUserData()
	anotherUser := testutils.SetupUserData()

	b1 := database.Book{
		UserID:  user.ID,
		Label:   "js",
		USN:     1123,
		Deleted: false,
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")
	b2 := database.Book{
		UserID:  user.ID,
		Label:   "css",
		USN:     1125,
		Deleted: false,
	}
	testutils.MustExec(t, db.Save(&b2), "preparing b2")
	b3 := database.Book{
		UserID:  anotherUser.ID,
		Label:   "css",
		USN:     1128,
		Deleted: false,
	}
	testutils.MustExec(t, db.Save(&b3), "preparing b3")
	b4 := database.Book{
		UserID:  user.ID,
		Label:   "",
		USN:     1129,
		Deleted: true,
	}
	testutils.MustExec(t, db.Save(&b4), "preparing b4")

	// Execute
	req := testutils.MakeReq(server, "GET", "/v3/books", "")
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var payload []presenters.Book
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	var b1Record, b2Record database.Book
	testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&b1Record), "finding b1")
	testutils.MustExec(t, db.Where("id = ?", b2.ID).First(&b2Record), "finding b2")
	testutils.MustExec(t, db.Where("id = ?", b2.ID).First(&b2Record), "finding b2")

	expected := []presenters.Book{
		{
			UUID:      b2Record.UUID,
			CreatedAt: b2Record.CreatedAt,
			UpdatedAt: b2Record.UpdatedAt,
			Label:     b2Record.Label,
			USN:       b2Record.USN,
		},
		{
			UUID:      b1Record.UUID,
			CreatedAt: b1Record.CreatedAt,
			UpdatedAt: b1Record.UpdatedAt,
			Label:     b1Record.Label,
			USN:       b1Record.USN,
		},
	}

	assert.DeepEqual(t, payload, expected, "payload mismatch")
}

func TestGetBooksByName(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	user := testutils.SetupUserData()
	anotherUser := testutils.SetupUserData()
	req := testutils.MakeReq(server, "GET", "/v3/books?name=js", "")

	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")
	b2 := database.Book{
		UserID: user.ID,
		Label:  "css",
	}
	testutils.MustExec(t, db.Save(&b2), "preparing b2")
	b3 := database.Book{
		UserID: anotherUser.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b3), "preparing b3")

	// Execute
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var payload []presenters.Book
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	var b1Record database.Book
	testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&b1Record), "finding b1")

	expected := []presenters.Book{
		{
			UUID:      b1Record.UUID,
			CreatedAt: b1Record.CreatedAt,
			UpdatedAt: b1Record.UpdatedAt,
			Label:     b1Record.Label,
			USN:       b1Record.USN,
		},
	}

	assert.DeepEqual(t, payload, expected, "payload mismatch")
}

func TestDeleteBook(t *testing.T) {
	testCases := []struct {
		label          string
		deleted        bool
		expectedB2USN  int
		expectedMaxUSN int
		expectedN2USN  int
		expectedN3USN  int
	}{
		{
			label:          "n1 content",
			deleted:        false,
			expectedMaxUSN: 61,
			expectedB2USN:  61,
			expectedN2USN:  59,
			expectedN3USN:  60,
		},
		{
			label:          "",
			deleted:        true,
			expectedMaxUSN: 59,
			expectedB2USN:  59,
			expectedN2USN:  5,
			expectedN3USN:  6,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("originally deleted %t", tc.deleted), func(t *testing.T) {
			defer testutils.ClearData()
			db := database.DBConn

			// Setup
			server := httptest.NewServer(NewRouter(&App{
				Clock: clock.NewMock(),
			}))
			defer server.Close()

			user := testutils.SetupUserData()
			testutils.MustExec(t, db.Model(&user).Update("max_usn", 58), "preparing user max_usn")
			anotherUser := testutils.SetupUserData()
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 109), "preparing another user max_usn")

			b1 := database.Book{
				UserID: user.ID,
				Label:  "js",
				USN:    1,
			}
			testutils.MustExec(t, db.Save(&b1), "preparing a book data")
			b2 := database.Book{
				UserID:  user.ID,
				Label:   tc.label,
				USN:     2,
				Deleted: tc.deleted,
			}
			testutils.MustExec(t, db.Save(&b2), "preparing a book data")
			b3 := database.Book{
				UserID: anotherUser.ID,
				Label:  "linux",
				USN:    3,
			}
			testutils.MustExec(t, db.Save(&b3), "preparing a book data")

			var n2Body string
			if !tc.deleted {
				n2Body = "n2 content"
			}
			var n3Body string
			if !tc.deleted {
				n3Body = "n3 content"
			}

			n1 := database.Note{
				UserID:   user.ID,
				BookUUID: b1.UUID,
				Body:     "n1 content",
				USN:      4,
			}
			testutils.MustExec(t, db.Save(&n1), "preparing a note data")
			n2 := database.Note{
				UserID:   user.ID,
				BookUUID: b2.UUID,
				Body:     n2Body,
				USN:      5,
				Deleted:  tc.deleted,
			}
			testutils.MustExec(t, db.Save(&n2), "preparing a note data")
			n3 := database.Note{
				UserID:   user.ID,
				BookUUID: b2.UUID,
				Body:     n3Body,
				USN:      6,
				Deleted:  tc.deleted,
			}
			testutils.MustExec(t, db.Save(&n3), "preparing a note data")
			n4 := database.Note{
				UserID:   user.ID,
				BookUUID: b2.UUID,
				Body:     "",
				USN:      7,
				Deleted:  true,
			}
			testutils.MustExec(t, db.Save(&n4), "preparing a note data")
			n5 := database.Note{
				UserID:   anotherUser.ID,
				BookUUID: b3.UUID,
				Body:     "n5 content",
				USN:      8,
			}
			testutils.MustExec(t, db.Save(&n5), "preparing a note data")

			endpoint := fmt.Sprintf("/v3/books/%s", b2.UUID)
			req := testutils.MakeReq(server, "DELETE", endpoint, "")
			req.Header.Set("Version", "0.1.1")
			req.Header.Set("Origin", "chrome-extension://iaolnfnipkoinabdbbakcmkkdignedce")

			// Execute
			res := testutils.HTTPAuthDo(t, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, "")

			var b1Record, b2Record, b3Record database.Book
			var n1Record, n2Record, n3Record, n4Record, n5Record database.Note
			var userRecord database.User
			var bookCount, noteCount int

			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
			testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
			testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&b1Record), "finding b1")
			testutils.MustExec(t, db.Where("id = ?", b2.ID).First(&b2Record), "finding b2")
			testutils.MustExec(t, db.Where("id = ?", b3.ID).First(&b3Record), "finding b3")
			testutils.MustExec(t, db.Where("id = ?", n1.ID).First(&n1Record), "finding n1")
			testutils.MustExec(t, db.Where("id = ?", n2.ID).First(&n2Record), "finding n2")
			testutils.MustExec(t, db.Where("id = ?", n3.ID).First(&n3Record), "finding n3")
			testutils.MustExec(t, db.Where("id = ?", n4.ID).First(&n4Record), "finding n4")
			testutils.MustExec(t, db.Where("id = ?", n5.ID).First(&n5Record), "finding n5")
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

			assert.Equal(t, bookCount, 3, "book count mismatch")
			assert.Equal(t, noteCount, 5, "note count mismatch")

			assert.Equal(t, userRecord.MaxUSN, tc.expectedMaxUSN, "user max_usn mismatch")

			assert.Equal(t, b1Record.Deleted, false, "b1 deleted mismatch")
			assert.Equal(t, b1Record.Label, b1.Label, "b1 content mismatch")
			assert.Equal(t, b1Record.USN, b1.USN, "b1 usn mismatch")
			assert.Equal(t, b2Record.Deleted, true, "b2 deleted mismatch")
			assert.Equal(t, b2Record.Label, "", "b2 content mismatch")
			assert.Equal(t, b2Record.USN, tc.expectedB2USN, "b2 usn mismatch")
			assert.Equal(t, b3Record.Deleted, false, "b3 deleted mismatch")
			assert.Equal(t, b3Record.Label, b3.Label, "b3 content mismatch")
			assert.Equal(t, b3Record.USN, b3.USN, "b3 usn mismatch")

			assert.Equal(t, n1Record.USN, n1.USN, "n1 usn mismatch")
			assert.Equal(t, n1Record.Deleted, false, "n1 deleted mismatch")
			assert.Equal(t, n1Record.Body, n1.Body, "n1 content mismatch")

			assert.Equal(t, n2Record.USN, tc.expectedN2USN, "n2 usn mismatch")
			assert.Equal(t, n2Record.Deleted, true, "n2 deleted mismatch")
			assert.Equal(t, n2Record.Body, "", "n2 content mismatch")

			assert.Equal(t, n3Record.USN, tc.expectedN3USN, "n3 usn mismatch")
			assert.Equal(t, n3Record.Deleted, true, "n3 deleted mismatch")
			assert.Equal(t, n3Record.Body, "", "n3 content mismatch")

			// if already deleted, usn should remain the same and hence should not contribute to bumping the max_usn
			assert.Equal(t, n4Record.USN, n4.USN, "n4 usn mismatch")
			assert.Equal(t, n4Record.Deleted, true, "n4 deleted mismatch")
			assert.Equal(t, n4Record.Body, "", "n4 content mismatch")

			assert.Equal(t, n5Record.USN, n5.USN, "n5 usn mismatch")
			assert.Equal(t, n5Record.Deleted, false, "n5 deleted mismatch")
			assert.Equal(t, n5Record.Body, n5.Body, "n5 content mismatch")
		})
	}
}

func TestCreateBook(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	user := testutils.SetupUserData()
	testutils.MustExec(t, db.Model(&user).Update("max_usn", 101), "preparing user max_usn")

	req := testutils.MakeReq(server, "POST", "/v3/books", `{"name": "js"}`)
	req.Header.Set("Version", "0.1.1")
	req.Header.Set("Origin", "chrome-extension://iaolnfnipkoinabdbbakcmkkdignedce")

	// Execute
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var bookRecord database.Book
	var userRecord database.User
	var bookCount, noteCount int
	testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
	testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
	testutils.MustExec(t, db.First(&bookRecord), "finding book")
	testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

	maxUSN := 102

	assert.Equalf(t, bookCount, 1, "book count mismatch")
	assert.Equalf(t, noteCount, 0, "note count mismatch")

	assert.NotEqual(t, bookRecord.UUID, "", "book uuid should have been generated")
	assert.Equal(t, bookRecord.Label, "js", "book name mismatch")
	assert.Equal(t, bookRecord.UserID, user.ID, "book user_id mismatch")
	assert.Equal(t, bookRecord.USN, maxUSN, "book user_id mismatch")
	assert.Equal(t, userRecord.MaxUSN, maxUSN, "user max_usn mismatch")

	var got CreateBookResp
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(errors.Wrap(err, "decoding got"))
	}
	expected := CreateBookResp{
		Book: presenters.Book{
			UUID:      bookRecord.UUID,
			USN:       bookRecord.USN,
			CreatedAt: bookRecord.CreatedAt,
			UpdatedAt: bookRecord.UpdatedAt,
			Label:     "js",
		},
	}

	if ok := reflect.DeepEqual(got, expected); !ok {
		t.Errorf("Payload does not match.\nActual:   %+v\nExpected: %+v", got, expected)
	}
}

func TestCreateBookDuplicate(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	user := testutils.SetupUserData()
	testutils.MustExec(t, db.Model(&user).Update("max_usn", 101), "preparing user max_usn")

	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
		USN:    58,
	}
	testutils.MustExec(t, db.Save(&b1), "preparing book data")

	// Execute
	req := testutils.MakeReq(server, "POST", "/v3/books", `{"name": "js"}`)
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusConflict, "")

	var bookRecord database.Book
	var bookCount, noteCount int
	var userRecord database.User
	testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
	testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
	testutils.MustExec(t, db.First(&bookRecord), "finding book")
	testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

	assert.Equalf(t, bookCount, 1, "book count mismatch")
	assert.Equalf(t, noteCount, 0, "note count mismatch")

	assert.Equal(t, bookRecord.Label, "js", "book name mismatch")
	assert.Equal(t, bookRecord.UserID, user.ID, "book user_id mismatch")
	assert.Equal(t, bookRecord.USN, b1.USN, "book usn mismatch")
	assert.Equal(t, userRecord.MaxUSN, 101, "user max_usn mismatch")
}

func TestUpdateBook(t *testing.T) {
	updatedLabel := "updated-label"

	b1UUID := "ead8790f-aff9-4bdf-8eec-f734ccd29202"
	b2UUID := "0ecaac96-8d72-4e04-8925-5a21b79a16da"

	testCases := []struct {
		payload           string
		bookUUID          string
		bookDeleted       bool
		bookLabel         string
		expectedBookLabel string
	}{
		{
			payload: fmt.Sprintf(`{
				"name": "%s"
			}`, updatedLabel),
			bookUUID:          b1UUID,
			bookDeleted:       false,
			bookLabel:         "original-label",
			expectedBookLabel: updatedLabel,
		},
		// if a deleted book is updated, it should be un-deleted
		{
			payload: fmt.Sprintf(`{
				"name": "%s"
			}`, updatedLabel),
			bookUUID:          b1UUID,
			bookDeleted:       true,
			bookLabel:         "",
			expectedBookLabel: updatedLabel,
		},
	}

	for idx, tc := range testCases {
		func() {
			defer testutils.ClearData()
			db := database.DBConn

			// Setup
			server := httptest.NewServer(NewRouter(&App{
				Clock: clock.NewMock(),
			}))
			defer server.Close()

			user := testutils.SetupUserData()
			testutils.MustExec(t, db.Model(&user).Update("max_usn", 101), "preparing user max_usn")

			b1 := database.Book{
				UUID:    tc.bookUUID,
				UserID:  user.ID,
				Label:   tc.bookLabel,
				Deleted: tc.bookDeleted,
			}
			testutils.MustExec(t, db.Save(&b1), "preparing b1")
			b2 := database.Book{
				UUID:   b2UUID,
				UserID: user.ID,
				Label:  "js",
			}
			testutils.MustExec(t, db.Save(&b2), "preparing b2")

			// Execute
			endpoint := fmt.Sprintf("/v3/books/%s", tc.bookUUID)
			req := testutils.MakeReq(server, "PATCH", endpoint, tc.payload)
			res := testutils.HTTPAuthDo(t, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, fmt.Sprintf("status code mismatch for test case %d", idx))

			var bookRecord database.Book
			var userRecord database.User
			var noteCount, bookCount int
			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
			testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
			testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

			assert.Equalf(t, bookCount, 2, "book count mismatch")
			assert.Equalf(t, noteCount, 0, "note count mismatch")

			assert.Equalf(t, bookRecord.UUID, tc.bookUUID, "book uuid mismatch")
			assert.Equalf(t, bookRecord.Label, tc.expectedBookLabel, "book label mismatch")
			assert.Equalf(t, bookRecord.USN, 102, "book usn mismatch")
			assert.Equalf(t, bookRecord.Deleted, false, "book Deleted mismatch")

			assert.Equal(t, userRecord.MaxUSN, 102, fmt.Sprintf("user max_usn mismatch for test case %d", idx))
		}()
	}
}
