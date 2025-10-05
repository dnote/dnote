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

func getExpectedNotePayload(n database.Note, b database.Book, u database.User) presenters.Note {
	return presenters.Note{
		UUID:      n.UUID,
		CreatedAt: truncateMicro(n.CreatedAt),
		UpdatedAt: truncateMicro(n.UpdatedAt),
		Body:      n.Body,
		AddedOn:   n.AddedOn,
		Public:    n.Public,
		USN:       n.USN,
		Book: presenters.NoteBook{
			UUID:  b.UUID,
			Label: b.Label,
		},
		User: presenters.NoteUser{
			UUID: u.UUID,
		},
	}
}

func TestGetNotes(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	// Setup
	a := app.NewTest()
	a.DB = db
	a.Clock = clock.NewMock()
	server := MustNewServer(t, &a)
	defer server.Close()

	user := testutils.SetupUserData(db)
	testutils.SetupAccountData(db, user, "alice@test.com", "pass1234")
	anotherUser := testutils.SetupUserData(db)
	testutils.SetupAccountData(db, anotherUser, "bob@test.com", "pass1234")

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
		Label:  "css",
	}
	testutils.MustExec(t, db.Save(&b3), "preparing b3")

	n1 := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n1 content",
		USN:      11,
		Deleted:  false,
		AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n1), "preparing n1")
	n2 := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n2 content",
		USN:      14,
		Deleted:  false,
		AddedOn:  time.Date(2018, time.August, 11, 22, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n2), "preparing n2")
	n3 := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n3 content",
		USN:      17,
		Deleted:  false,
		AddedOn:  time.Date(2017, time.January, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n3), "preparing n3")
	n4 := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   user.ID,
		BookUUID: b2.UUID,
		Body:     "n4 content",
		USN:      18,
		Deleted:  false,
		AddedOn:  time.Date(2018, time.September, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n4), "preparing n4")
	n5 := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   anotherUser.ID,
		BookUUID: b3.UUID,
		Body:     "n5 content",
		USN:      19,
		Deleted:  false,
		AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n5), "preparing n5")
	n6 := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "",
		USN:      11,
		Deleted:  true,
		AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n6), "preparing n6")

	// Execute
	endpoint := "/api/v3/notes"

	req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("%s?year=2018&month=8", endpoint), "")
	res := testutils.HTTPAuthDo(t, db, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var payload GetNotesResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	var n2Record, n1Record database.Note
	testutils.MustExec(t, db.Where("uuid = ?", n2.UUID).First(&n2Record), "finding n2Record")
	testutils.MustExec(t, db.Where("uuid = ?", n1.UUID).First(&n1Record), "finding n1Record")

	expected := GetNotesResponse{
		Notes: []presenters.Note{
			getExpectedNotePayload(n2Record, b1, user),
			getExpectedNotePayload(n1Record, b1, user),
		},
		Total: 2,
	}

	assert.DeepEqual(t, payload, expected, "payload mismatch")
}

func TestGetNote(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	// Setup
	a := app.NewTest()
	a.DB = db
	a.Clock = clock.NewMock()
	server := MustNewServer(t, &a)
	defer server.Close()

	user := testutils.SetupUserData(db)
	anotherUser := testutils.SetupUserData(db)

	b1 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")

	privateNote := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "privateNote content",
		Public:   false,
	}
	testutils.MustExec(t, db.Save(&privateNote), "preparing privateNote")
	publicNote := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "publicNote content",
		Public:   true,
	}
	testutils.MustExec(t, db.Save(&publicNote), "preparing publicNote")
	deletedNote := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Deleted:  true,
	}
	testutils.MustExec(t, db.Save(&deletedNote), "preparing deletedNote")

	getURL := func(noteUUID string) string {
		return fmt.Sprintf("/api/v3/notes/%s", noteUUID)
	}

	t.Run("owner accessing private note", func(t *testing.T) {
		// Execute
		url := getURL(publicNote.UUID)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var payload presenters.Note
		if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
			t.Fatal(errors.Wrap(err, "decoding payload"))
		}

		var n2Record database.Note
		testutils.MustExec(t, db.Where("uuid = ?", publicNote.UUID).First(&n2Record), "finding n2Record")

		expected := getExpectedNotePayload(n2Record, b1, user)
		assert.DeepEqual(t, payload, expected, "payload mismatch")
	})

	t.Run("owner accessing public note", func(t *testing.T) {
		// Execute
		url := getURL(publicNote.UUID)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var payload presenters.Note
		if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
			t.Fatal(errors.Wrap(err, "decoding payload"))
		}

		var n2Record database.Note
		testutils.MustExec(t, db.Where("uuid = ?", publicNote.UUID).First(&n2Record), "finding n2Record")

		expected := getExpectedNotePayload(n2Record, b1, user)
		assert.DeepEqual(t, payload, expected, "payload mismatch")
	})

	t.Run("non-owner accessing public note", func(t *testing.T) {
		// Execute
		url := getURL(publicNote.UUID)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, db, req, anotherUser)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var payload presenters.Note
		if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
			t.Fatal(errors.Wrap(err, "decoding payload"))
		}

		var n2Record database.Note
		testutils.MustExec(t, db.Where("uuid = ?", publicNote.UUID).First(&n2Record), "finding n2Record")

		expected := getExpectedNotePayload(n2Record, b1, user)
		assert.DeepEqual(t, payload, expected, "payload mismatch")
	})

	t.Run("non-owner accessing private note", func(t *testing.T) {
		// Execute
		url := getURL(privateNote.UUID)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, db, req, anotherUser)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(errors.Wrap(err, "reading body"))
		}

		assert.DeepEqual(t, string(body), "not found\n", "payload mismatch")
	})

	t.Run("guest accessing public note", func(t *testing.T) {
		// Execute
		url := getURL(publicNote.UUID)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var payload presenters.Note
		if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
			t.Fatal(errors.Wrap(err, "decoding payload"))
		}

		var n2Record database.Note
		testutils.MustExec(t, db.Where("uuid = ?", publicNote.UUID).First(&n2Record), "finding n2Record")

		expected := getExpectedNotePayload(n2Record, b1, user)
		assert.DeepEqual(t, payload, expected, "payload mismatch")
	})

	t.Run("guest accessing private note", func(t *testing.T) {
		// Execute
		url := getURL(privateNote.UUID)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(errors.Wrap(err, "reading body"))
		}

		assert.DeepEqual(t, string(body), "not found\n", "payload mismatch")
	})

	t.Run("nonexistent", func(t *testing.T) {
		// Execute
		url := getURL("somerandomstring")
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(errors.Wrap(err, "reading body"))
		}

		assert.DeepEqual(t, string(body), "not found\n", "payload mismatch")
	})

	t.Run("deleted", func(t *testing.T) {
		// Execute
		url := getURL(deletedNote.UUID)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(errors.Wrap(err, "reading body"))
		}

		assert.DeepEqual(t, string(body), "not found\n", "payload mismatch")
	})
}

func TestCreateNote(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	// Setup
	a := app.NewTest()
	a.DB = db
	a.Clock = clock.NewMock()
	server := MustNewServer(t, &a)
	defer server.Close()

	user := testutils.SetupUserData(db)
	testutils.SetupAccountData(db, user, "alice@test.com", "pass1234")
	testutils.MustExec(t, db.Model(&user).Update("max_usn", 101), "preparing user max_usn")

	b1 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: user.ID,
		Label:  "js",
		USN:    58,
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")

	// Execute

	dat := fmt.Sprintf(`{"book_uuid": "%s", "content": "note content"}`, b1.UUID)
	req := testutils.MakeReq(server.URL, "POST", "/api/v3/notes", dat)
	res := testutils.HTTPAuthDo(t, db, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusCreated, "")

	var noteRecord database.Note
	var bookRecord database.Book
	var userRecord database.User
	var bookCount, noteCount int64
	testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
	testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
	testutils.MustExec(t, db.First(&noteRecord), "finding note")
	testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
	testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

	assert.Equalf(t, bookCount, int64(1), "book count mismatch")
	assert.Equalf(t, noteCount, int64(1), "note count mismatch")

	assert.Equal(t, bookRecord.Label, b1.Label, "book name mismatch")
	assert.Equal(t, bookRecord.UUID, b1.UUID, "book uuid mismatch")
	assert.Equal(t, bookRecord.UserID, b1.UserID, "book user_id mismatch")
	assert.Equal(t, bookRecord.USN, 58, "book usn mismatch")

	assert.NotEqual(t, noteRecord.UUID, "", "note uuid should have been generated")
	assert.Equal(t, noteRecord.BookUUID, b1.UUID, "note book_uuid mismatch")
	assert.Equal(t, noteRecord.Body, "note content", "note content mismatch")
	assert.Equal(t, noteRecord.USN, 102, "note usn mismatch")
}

func TestDeleteNote(t *testing.T) {
	b1UUID := "37868a8e-a844-4265-9a4f-0be598084733"

	testCases := []struct {
		content        string
		deleted        bool
		originalUSN    int
		expectedUSN    int
		expectedMaxUSN int
	}{
		{
			content:        "n1 content",
			deleted:        false,
			originalUSN:    12,
			expectedUSN:    982,
			expectedMaxUSN: 982,
		},
		{
			content:        "",
			deleted:        true,
			originalUSN:    12,
			expectedUSN:    982,
			expectedMaxUSN: 982,
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

			user := testutils.SetupUserData(db)
			testutils.SetupAccountData(db, user, "alice@test.com", "pass1234")
			testutils.MustExec(t, db.Model(&user).Update("max_usn", 981), "preparing user max_usn")

			b1 := database.Book{
				UUID:   b1UUID,
				UserID: user.ID,
				Label:  "js",
			}
			testutils.MustExec(t, db.Save(&b1), "preparing b1")
			note := database.Note{
				UUID:     testutils.MustUUID(t),
				UserID:   user.ID,
				BookUUID: b1.UUID,
				Body:     tc.content,
				Deleted:  tc.deleted,
				USN:      tc.originalUSN,
			}
			testutils.MustExec(t, db.Save(&note), "preparing note")

			// Execute
			endpoint := fmt.Sprintf("/api/v3/notes/%s", note.UUID)
			req := testutils.MakeReq(server.URL, "DELETE", endpoint, "")
			res := testutils.HTTPAuthDo(t, db, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, "")

			var bookRecord database.Book
			var noteRecord database.Note
			var userRecord database.User
			var bookCount, noteCount int64
			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
			testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
			testutils.MustExec(t, db.Where("uuid = ?", note.UUID).First(&noteRecord), "finding note")
			testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

			assert.Equalf(t, bookCount, int64(1), "book count mismatch")
			assert.Equalf(t, noteCount, int64(1), "note count mismatch")

			assert.Equal(t, noteRecord.UUID, note.UUID, "note uuid mismatch for test case")
			assert.Equal(t, noteRecord.Body, "", "note content mismatch for test case")
			assert.Equal(t, noteRecord.Deleted, true, "note deleted mismatch for test case")
			assert.Equal(t, noteRecord.BookUUID, note.BookUUID, "note book_uuid mismatch for test case")
			assert.Equal(t, noteRecord.UserID, note.UserID, "note user_id mismatch for test case")
			assert.Equal(t, noteRecord.USN, tc.expectedUSN, "note usn mismatch for test case")

			assert.Equal(t, userRecord.MaxUSN, tc.expectedMaxUSN, "user max_usn mismatch for test case")
		})
	}
}

func TestUpdateNote(t *testing.T) {
	updatedBody := "some updated content"

	b1UUID := "37868a8e-a844-4265-9a4f-0be598084733"
	b2UUID := "8f3bd424-6aa5-4ed5-910d-e5b38ab09f8c"

	type payloadData struct {
		Content  *string `schema:"content" json:"content,omitempty"`
		BookUUID *string `schema:"book_uuid" json:"book_uuid,omitempty"`
		Public   *bool   `schema:"public" json:"public,omitempty"`
	}

	testCases := []struct {
		payload              testutils.PayloadWrapper
		noteUUID             string
		noteBookUUID         string
		noteBody             string
		notePublic           bool
		noteDeleted          bool
		expectedNoteBody     string
		expectedNoteBookName string
		expectedNoteBookUUID string
		expectedNotePublic   bool
	}{
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					Content: &updatedBody,
				},
			},
			noteUUID:             "ab50aa32-b232-40d8-b10f-10a7f9134053",
			noteBookUUID:         b1UUID,
			notePublic:           false,
			noteBody:             "original content",
			noteDeleted:          false,
			expectedNoteBookUUID: b1UUID,
			expectedNoteBody:     "some updated content",
			expectedNoteBookName: "css",
			expectedNotePublic:   false,
		},
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					BookUUID: &b1UUID,
				},
			},
			noteUUID:             "ab50aa32-b232-40d8-b10f-10a7f9134053",
			noteBookUUID:         b1UUID,
			notePublic:           false,
			noteBody:             "original content",
			noteDeleted:          false,
			expectedNoteBookUUID: b1UUID,
			expectedNoteBody:     "original content",
			expectedNoteBookName: "css",
			expectedNotePublic:   false,
		},
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					BookUUID: &b2UUID,
				},
			},
			noteUUID:             "ab50aa32-b232-40d8-b10f-10a7f9134053",
			noteBookUUID:         b1UUID,
			notePublic:           false,
			noteBody:             "original content",
			noteDeleted:          false,
			expectedNoteBookUUID: b2UUID,
			expectedNoteBody:     "original content",
			expectedNoteBookName: "js",
			expectedNotePublic:   false,
		},
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					BookUUID: &b2UUID,
					Content:  &updatedBody,
				},
			},
			noteUUID:             "ab50aa32-b232-40d8-b10f-10a7f9134053",
			noteBookUUID:         b1UUID,
			notePublic:           false,
			noteBody:             "original content",
			noteDeleted:          false,
			expectedNoteBookUUID: b2UUID,
			expectedNoteBody:     "some updated content",
			expectedNoteBookName: "js",
			expectedNotePublic:   false,
		},
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					BookUUID: &b1UUID,
					Content:  &updatedBody,
				},
			},
			noteUUID:             "ab50aa32-b232-40d8-b10f-10a7f9134053",
			noteBookUUID:         b1UUID,
			notePublic:           false,
			noteBody:             "",
			noteDeleted:          true,
			expectedNoteBookUUID: b1UUID,
			expectedNoteBody:     updatedBody,
			expectedNoteBookName: "js",
			expectedNotePublic:   false,
		},
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					Public: &testutils.TrueVal,
				},
			},
			noteUUID:             "ab50aa32-b232-40d8-b10f-10a7f9134053",
			noteBookUUID:         b1UUID,
			notePublic:           false,
			noteBody:             "original content",
			noteDeleted:          false,
			expectedNoteBookUUID: b1UUID,
			expectedNoteBody:     "original content",
			expectedNoteBookName: "css",
			expectedNotePublic:   true,
		},
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					Public: &testutils.FalseVal,
				},
			},
			noteUUID:             "ab50aa32-b232-40d8-b10f-10a7f9134053",
			noteBookUUID:         b1UUID,
			notePublic:           true,
			noteBody:             "original content",
			noteDeleted:          false,
			expectedNoteBookUUID: b1UUID,
			expectedNoteBody:     "original content",
			expectedNoteBookName: "css",
			expectedNotePublic:   false,
		},
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					Content: &updatedBody,
					Public:  &testutils.FalseVal,
				},
			},
			noteUUID:             "ab50aa32-b232-40d8-b10f-10a7f9134053",
			noteBookUUID:         b1UUID,
			notePublic:           true,
			noteBody:             "original content",
			noteDeleted:          false,
			expectedNoteBookUUID: b1UUID,
			expectedNoteBody:     updatedBody,
			expectedNoteBookName: "css",
			expectedNotePublic:   false,
		},
		{
			payload: testutils.PayloadWrapper{
				Data: payloadData{
					BookUUID: &b2UUID,
					Content:  &updatedBody,
					Public:   &testutils.TrueVal,
				},
			},
			noteUUID:             "ab50aa32-b232-40d8-b10f-10a7f9134053",
			noteBookUUID:         b1UUID,
			notePublic:           false,
			noteBody:             "original content",
			noteDeleted:          false,
			expectedNoteBookUUID: b2UUID,
			expectedNoteBody:     updatedBody,
			expectedNoteBookName: "js",
			expectedNotePublic:   true,
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

			user := testutils.SetupUserData(db)
			testutils.SetupAccountData(db, user, "alice@test.com", "pass1234")

			testutils.MustExec(t, db.Model(&user).Update("max_usn", 101), "preparing user max_usn")

			b1 := database.Book{
				UUID:   b1UUID,
				UserID: user.ID,
				Label:  "css",
			}
			testutils.MustExec(t, db.Save(&b1), "preparing b1")
			b2 := database.Book{
				UUID:   b2UUID,
				UserID: user.ID,
				Label:  "js",
			}
			testutils.MustExec(t, db.Save(&b2), "preparing b2")

			note := database.Note{
				UUID:     tc.noteUUID,
				UserID:   user.ID,
				BookUUID: tc.noteBookUUID,
				Body:     tc.noteBody,
				Deleted:  tc.noteDeleted,
				Public:   tc.notePublic,
			}
			testutils.MustExec(t, db.Save(&note), "preparing note")

			// Execute
			var req *http.Request

			endpoint := fmt.Sprintf("/api/v3/notes/%s", note.UUID)
			req = testutils.MakeReq(server.URL, "PATCH", endpoint, tc.payload.ToJSON(t))

			res := testutils.HTTPAuthDo(t, db, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, "status code mismatch for test case")

			var bookRecord database.Book
			var noteRecord database.Note
			var userRecord database.User
			var noteCount, bookCount int64
			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting books")
			testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes")
			testutils.MustExec(t, db.Where("uuid = ?", note.UUID).First(&noteRecord), "finding note")
			testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user record")

			assert.Equalf(t, bookCount, int64(2), "book count mismatch")
			assert.Equalf(t, noteCount, int64(1), "note count mismatch")

			assert.Equal(t, noteRecord.UUID, tc.noteUUID, "note uuid mismatch for test case")
			assert.Equal(t, noteRecord.Body, tc.expectedNoteBody, "note content mismatch for test case")
			assert.Equal(t, noteRecord.BookUUID, tc.expectedNoteBookUUID, "note book_uuid mismatch for test case")
			assert.Equal(t, noteRecord.Public, tc.expectedNotePublic, "note public mismatch for test case")
			assert.Equal(t, noteRecord.USN, 102, "note usn mismatch for test case")

			assert.Equal(t, userRecord.MaxUSN, 102, "user max_usn mismatch for test case")
		})
	}
}
