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

package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

// truncateMicro rounds time to microsecond precision to match SQLite storage
func truncateMicro(t time.Time) time.Time {
	return t.Round(time.Microsecond)
}

func TestGetBooks(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	// Setup
	a := app.NewTest()
	a.DB = db
	a.Clock = clock.NewMock()
	server := MustNewServer(t, &a)
	defer server.Close()

	user := testutils.SetupUserData(db, "alice@test.com", "pass1234")
	anotherUser := testutils.SetupUserData(db, "bob@test.com", "pass1234")

	b1 := database.Book{
		UUID:    testutils.MustUUID(t),
		UserID:  user.ID,
		Label:   "js",
		USN:     1123,
		Deleted: false,
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")
	b2 := database.Book{
		UUID:    testutils.MustUUID(t),
		UserID:  user.ID,
		Label:   "css",
		USN:     1125,
		Deleted: false,
	}
	testutils.MustExec(t, db.Save(&b2), "preparing b2")
	b3 := database.Book{
		UUID:    testutils.MustUUID(t),
		UserID:  anotherUser.ID,
		Label:   "css",
		USN:     1128,
		Deleted: false,
	}
	testutils.MustExec(t, db.Save(&b3), "preparing b3")
	b4 := database.Book{
		UUID:    testutils.MustUUID(t),
		UserID:  user.ID,
		Label:   "",
		USN:     1129,
		Deleted: true,
	}
	testutils.MustExec(t, db.Save(&b4), "preparing b4")

	// Execute
	endpoint := "/api/v3/books"

	req := testutils.MakeReq(server.URL, "GET", endpoint, "")
	res := testutils.HTTPAuthDo(t, db, req, user)

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
			CreatedAt: truncateMicro(b2Record.CreatedAt),
			UpdatedAt: truncateMicro(b2Record.UpdatedAt),
			Label:     b2Record.Label,
			USN:       b2Record.USN,
		},
		{
			UUID:      b1Record.UUID,
			CreatedAt: truncateMicro(b1Record.CreatedAt),
			UpdatedAt: truncateMicro(b1Record.UpdatedAt),
			Label:     b1Record.Label,
			USN:       b1Record.USN,
		},
	}

	// Truncate payload timestamps to match SQLite precision
	for i := range payload {
		payload[i].CreatedAt = truncateMicro(payload[i].CreatedAt)
		payload[i].UpdatedAt = truncateMicro(payload[i].UpdatedAt)
	}

	assert.DeepEqual(t, payload, expected, "payload mismatch")
}

func TestGetBooksByName(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	// Setup
	a := app.NewTest()
	a.DB = db
	a.Clock = clock.NewMock()
	server := MustNewServer(t, &a)
	defer server.Close()

	user := testutils.SetupUserData(db, "alice@test.com", "pass1234")
	anotherUser := testutils.SetupUserData(db, "bob@test.com", "pass1234")

	b1 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")
	b2 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: user.ID,
		Label:  "css",
	}
	testutils.MustExec(t, db.Save(&b2), "preparing b2")
	b3 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: anotherUser.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b3), "preparing b3")

	// Execute
	endpoint := "/api/v3/books?name=js"

	req := testutils.MakeReq(server.URL, "GET", endpoint, "")
	res := testutils.HTTPAuthDo(t, db, req, user)

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
			CreatedAt: truncateMicro(b1Record.CreatedAt),
			UpdatedAt: truncateMicro(b1Record.UpdatedAt),
			Label:     b1Record.Label,
			USN:       b1Record.USN,
		},
	}

	for i := range payload {
		payload[i].CreatedAt = truncateMicro(payload[i].CreatedAt)
		payload[i].UpdatedAt = truncateMicro(payload[i].UpdatedAt)
	}

	assert.DeepEqual(t, payload, expected, "payload mismatch")
}

func TestGetBook(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	// Setup
	a := app.NewTest()
	a.DB = db
	a.Clock = clock.NewMock()
	server := MustNewServer(t, &a)
	defer server.Close()

	user := testutils.SetupUserData(db, "alice@test.com", "pass1234")
	anotherUser := testutils.SetupUserData(db, "bob@test.com", "pass1234")

	b1 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")
	b2 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: user.ID,
		Label:  "css",
	}
	testutils.MustExec(t, db.Save(&b2), "preparing b2")
	b3 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: anotherUser.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b3), "preparing b3")

	// Execute
	endpoint := fmt.Sprintf("/api/v3/books/%s", b1.UUID)
	req := testutils.MakeReq(server.URL, "GET", endpoint, "")
	res := testutils.HTTPAuthDo(t, db, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var payload presenters.Book
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	var b1Record database.Book
	testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&b1Record), "finding b1")

	expected := presenters.Book{
		UUID:      b1Record.UUID,
		CreatedAt: truncateMicro(b1Record.CreatedAt),
		UpdatedAt: truncateMicro(b1Record.UpdatedAt),
		Label:     b1Record.Label,
		USN:       b1Record.USN,
	}

	payload.CreatedAt = truncateMicro(payload.CreatedAt)
	payload.UpdatedAt = truncateMicro(payload.UpdatedAt)

	assert.DeepEqual(t, payload, expected, "payload mismatch")
}

func TestGetBookNonOwner(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	// Setup
	a := app.NewTest()
	a.DB = db
	a.Clock = clock.NewMock()
	server := MustNewServer(t, &a)
	defer server.Close()

	user := testutils.SetupUserData(db, "alice@test.com", "pass1234")
	nonOwner := testutils.SetupUserData(db, "bob@test.com", "pass1234")

	b1 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")

	// Execute
	endpoint := fmt.Sprintf("/api/v3/books/%s", b1.UUID)
	req := testutils.MakeReq(server.URL, "GET", endpoint, "")
	res := testutils.HTTPAuthDo(t, db, req, nonOwner)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(errors.Wrap(err, "reading body"))
	}
	assert.DeepEqual(t, string(body), "", "payload mismatch")
}

func TestCreateBook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.DB = db
		a.Clock = clock.NewMock()
		server := MustNewServer(t, &a)
		defer server.Close()

		user := testutils.SetupUserData(db, "alice@test.com", "pass1234")
		testutils.MustExec(t, db.Model(&user).Update("max_usn", 101), "preparing user max_usn")

		req := testutils.MakeReq(server.URL, "POST", "/api/v3/books", `{"name": "js"}`)

		// Execute
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusCreated, "")

		var bookRecord database.Book
		var userRecord database.User
		var bookCount, noteCount int64
		testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
		testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
		testutils.MustExec(t, db.First(&bookRecord), "finding book")
		testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

		maxUSN := 102

		assert.Equalf(t, bookCount, int64(1), "book count mismatch")
		assert.Equalf(t, noteCount, int64(0), "note count mismatch")

		assert.NotEqual(t, bookRecord.UUID, "", "book uuid should have been generated")
		assert.Equal(t, bookRecord.Label, "js", "book name mismatch")
		assert.Equal(t, bookRecord.UserID, user.ID, "book user_id mismatch")
		assert.Equal(t, bookRecord.USN, maxUSN, "book user_id mismatch")
		assert.Equal(t, userRecord.MaxUSN, maxUSN, "user max_usn mismatch")

		var got CreateBookResp
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatal(errors.Wrap(err, "decoding"))
		}
		expected := CreateBookResp{
			Book: presenters.Book{
				UUID:      bookRecord.UUID,
				USN:       bookRecord.USN,
				CreatedAt: truncateMicro(bookRecord.CreatedAt),
				UpdatedAt: truncateMicro(bookRecord.UpdatedAt),
				Label:     "js",
			},
		}

		got.Book.CreatedAt = truncateMicro(got.Book.CreatedAt)
		got.Book.UpdatedAt = truncateMicro(got.Book.UpdatedAt)

		assert.DeepEqual(t, got, expected, "payload mismatch")
	})

	t.Run("duplicate", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.DB = db
		a.Clock = clock.NewMock()
		server := MustNewServer(t, &a)
		defer server.Close()

		user := testutils.SetupUserData(db, "alice@test.com", "pass1234")
		testutils.MustExec(t, db.Model(&user).Update("max_usn", 101), "preparing user max_usn")

		b1 := database.Book{
			UUID:   testutils.MustUUID(t),
			UserID: user.ID,
			Label:  "js",
			USN:    58,
		}
		testutils.MustExec(t, db.Save(&b1), "preparing book data")

		// Execute
		req := testutils.MakeReq(server.URL, "POST", "/api/v3/books", `{"name": "js"}`)
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusConflict, "")

		var bookRecord database.Book
		var bookCount, noteCount int64
		var userRecord database.User
		testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
		testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
		testutils.MustExec(t, db.First(&bookRecord), "finding book")
		testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

		assert.Equalf(t, bookCount, int64(1), "book count mismatch")
		assert.Equalf(t, noteCount, int64(0), "note count mismatch")

		assert.Equal(t, bookRecord.Label, "js", "book name mismatch")
		assert.Equal(t, bookRecord.UserID, user.ID, "book user_id mismatch")
		assert.Equal(t, bookRecord.USN, b1.USN, "book usn mismatch")
		assert.Equal(t, userRecord.MaxUSN, 101, "user max_usn mismatch")
	})
}

func TestUpdateBook(t *testing.T) {
	updatedLabel := "updated-label"

	b1UUID := "ead8790f-aff9-4bdf-8eec-f734ccd29202"
	b2UUID := "0ecaac96-8d72-4e04-8925-5a21b79a16da"

	type payloadData struct {
		Name *string `schema:"name" json:"name,omitempty"`
	}

	testCases := []struct {
		payload           testutils.PayloadWrapper
		bookUUID          string
		bookDeleted       bool
		bookLabel         string
		expectedBookLabel string
	}{
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					Name: &updatedLabel,
				},
			},
			bookUUID:          b1UUID,
			bookDeleted:       false,
			bookLabel:         "original-label",
			expectedBookLabel: updatedLabel,
		},
		// if a deleted book is updated, it should be un-deleted
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					Name: &updatedLabel,
				},
			},
			bookUUID:          b1UUID,
			bookDeleted:       true,
			bookLabel:         "",
			expectedBookLabel: updatedLabel,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			db := testutils.InitMemoryDB(t)

			// Setup
			a := app.NewTest()
			a.DB = db
			a.Clock = clock.NewMock()
			server := MustNewServer(t, &a)
			defer server.Close()

			user := testutils.SetupUserData(db, "alice@test.com", "pass1234")
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
			endpoint := fmt.Sprintf("/api/v3/books/%s", tc.bookUUID)
			req := testutils.MakeReq(server.URL, "PATCH", endpoint, tc.payload.ToJSON(t))
			res := testutils.HTTPAuthDo(t, db, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, fmt.Sprintf("status code mismatch for test case %d", idx))

			var bookRecord database.Book
			var userRecord database.User
			var noteCount, bookCount int64
			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
			testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
			testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

			assert.Equalf(t, bookCount, int64(2), "book count mismatch")
			assert.Equalf(t, noteCount, int64(0), "note count mismatch")

			assert.Equalf(t, bookRecord.UUID, tc.bookUUID, "book uuid mismatch")
			assert.Equalf(t, bookRecord.Label, tc.expectedBookLabel, "book label mismatch")
			assert.Equalf(t, bookRecord.USN, 102, "book usn mismatch")
			assert.Equalf(t, bookRecord.Deleted, false, "book Deleted mismatch")

			assert.Equal(t, userRecord.MaxUSN, 102, fmt.Sprintf("user max_usn mismatch for test case %d", idx))
		})
	}
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
			db := testutils.InitMemoryDB(t)

			// Setup
			a := app.NewTest()
			a.DB = db
			a.Clock = clock.NewMock()
			server := MustNewServer(t, &a)
			defer server.Close()

			user := testutils.SetupUserData(db, "alice@test.com", "pass1234")
			testutils.MustExec(t, db.Model(&user).Update("max_usn", 58), "preparing user max_usn")
			anotherUser := testutils.SetupUserData(db, "bob@test.com", "pass1234")
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 109), "preparing another user max_usn")

			b1 := database.Book{
				UUID:   testutils.MustUUID(t),
				UserID: user.ID,
				Label:  "js",
				USN:    1,
			}
			testutils.MustExec(t, db.Save(&b1), "preparing a book data")
			b2 := database.Book{
				UUID:    testutils.MustUUID(t),
				UserID:  user.ID,
				Label:   tc.label,
				USN:     2,
				Deleted: tc.deleted,
			}
			testutils.MustExec(t, db.Save(&b2), "preparing a book data")
			b3 := database.Book{
				UUID:   testutils.MustUUID(t),
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
				UUID:     testutils.MustUUID(t),
				UserID:   user.ID,
				BookUUID: b1.UUID,
				Body:     "n1 content",
				USN:      4,
			}
			testutils.MustExec(t, db.Save(&n1), "preparing a note data")
			n2 := database.Note{
				UUID:     testutils.MustUUID(t),
				UserID:   user.ID,
				BookUUID: b2.UUID,
				Body:     n2Body,
				USN:      5,
				Deleted:  tc.deleted,
			}
			testutils.MustExec(t, db.Save(&n2), "preparing a note data")
			n3 := database.Note{
				UUID:     testutils.MustUUID(t),
				UserID:   user.ID,
				BookUUID: b2.UUID,
				Body:     n3Body,
				USN:      6,
				Deleted:  tc.deleted,
			}
			testutils.MustExec(t, db.Save(&n3), "preparing a note data")
			n4 := database.Note{
				UUID:     testutils.MustUUID(t),
				UserID:   user.ID,
				BookUUID: b2.UUID,
				Body:     "",
				USN:      7,
				Deleted:  true,
			}
			testutils.MustExec(t, db.Save(&n4), "preparing a note data")
			n5 := database.Note{
				UUID:     testutils.MustUUID(t),
				UserID:   anotherUser.ID,
				BookUUID: b3.UUID,
				Body:     "n5 content",
				USN:      8,
			}
			testutils.MustExec(t, db.Save(&n5), "preparing a note data")

			// Execute
			endpoint := fmt.Sprintf("/api/v3/books/%s", b2.UUID)

			req := testutils.MakeReq(server.URL, "DELETE", endpoint, "")
			res := testutils.HTTPAuthDo(t, db, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, "")

			var b1Record, b2Record, b3Record database.Book
			var n1Record, n2Record, n3Record, n4Record, n5Record database.Note
			var userRecord database.User
			var bookCount, noteCount int64

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

			assert.Equal(t, bookCount, int64(3), "book count mismatch")
			assert.Equal(t, noteCount, int64(5), "note count mismatch")

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
