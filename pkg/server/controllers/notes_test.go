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

package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func getExpectedNotePayload(n database.Note, b database.Book, u database.User) presenters.Note {
	return presenters.Note{
		UUID:      n.UUID,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
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
	testutils.RunForWebAndAPI(t, "get notes", func(t *testing.T, target testutils.EndpointType) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
		defer server.Close()

		user := testutils.SetupUserData()
		anotherUser := testutils.SetupUserData()

		b1 := database.Book{
			UserID: user.ID,
			Label:  "js",
		}
		testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
		b2 := database.Book{
			UserID: user.ID,
			Label:  "css",
		}
		testutils.MustExec(t, testutils.DB.Save(&b2), "preparing b2")
		b3 := database.Book{
			UserID: anotherUser.ID,
			Label:  "css",
		}
		testutils.MustExec(t, testutils.DB.Save(&b3), "preparing b3")

		n1 := database.Note{
			UserID:   user.ID,
			BookUUID: b1.UUID,
			Body:     "n1 content",
			USN:      11,
			Deleted:  false,
			AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
		}
		testutils.MustExec(t, testutils.DB.Save(&n1), "preparing n1")
		n2 := database.Note{
			UserID:   user.ID,
			BookUUID: b1.UUID,
			Body:     "n2 content",
			USN:      14,
			Deleted:  false,
			AddedOn:  time.Date(2018, time.August, 11, 22, 0, 0, 0, time.UTC).UnixNano(),
		}
		testutils.MustExec(t, testutils.DB.Save(&n2), "preparing n2")
		n3 := database.Note{
			UserID:   user.ID,
			BookUUID: b1.UUID,
			Body:     "n3 content",
			USN:      17,
			Deleted:  false,
			AddedOn:  time.Date(2017, time.January, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
		}
		testutils.MustExec(t, testutils.DB.Save(&n3), "preparing n3")
		n4 := database.Note{
			UserID:   user.ID,
			BookUUID: b2.UUID,
			Body:     "n4 content",
			USN:      18,
			Deleted:  false,
			AddedOn:  time.Date(2018, time.September, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
		}
		testutils.MustExec(t, testutils.DB.Save(&n4), "preparing n4")
		n5 := database.Note{
			UserID:   anotherUser.ID,
			BookUUID: b3.UUID,
			Body:     "n5 content",
			USN:      19,
			Deleted:  false,
			AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
		}
		testutils.MustExec(t, testutils.DB.Save(&n5), "preparing n5")
		n6 := database.Note{
			UserID:   user.ID,
			BookUUID: b1.UUID,
			Body:     "",
			USN:      11,
			Deleted:  true,
			AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
		}
		testutils.MustExec(t, testutils.DB.Save(&n6), "preparing n6")

		// Execute
		var endpoint string
		if target == testutils.EndpointWeb {
			endpoint = "/"
		} else {
			endpoint = "/api/v3/notes"
		}

		req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("%s?year=2018&month=8", endpoint), "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		if target == testutils.EndpointAPI {
			var payload GetNotesResponse
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			var n2Record, n1Record database.Note
			testutils.MustExec(t, testutils.DB.Where("uuid = ?", n2.UUID).First(&n2Record), "finding n2Record")
			testutils.MustExec(t, testutils.DB.Where("uuid = ?", n1.UUID).First(&n1Record), "finding n1Record")

			expected := GetNotesResponse{
				Notes: []presenters.Note{
					getExpectedNotePayload(n2Record, b1, user),
					getExpectedNotePayload(n1Record, b1, user),
				},
				Total: 2,
			}

			assert.DeepEqual(t, payload, expected, "payload mismatch")
		}
	})
}

func TestGetNote(t *testing.T) {
	defer testutils.ClearData(testutils.DB)

	// Setup
	server := MustNewServer(t, &app.App{
		Clock: clock.NewMock(),
		Config: config.Config{
			PageTemplateDir: "../views",
		},
	})
	defer server.Close()

	user := testutils.SetupUserData()
	anotherUser := testutils.SetupUserData()

	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")

	privateNote := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "privateNote content",
		Public:   false,
	}
	testutils.MustExec(t, testutils.DB.Save(&privateNote), "preparing privateNote")
	publicNote := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "publicNote content",
		Public:   true,
	}
	testutils.MustExec(t, testutils.DB.Save(&publicNote), "preparing publicNote")
	deletedNote := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Deleted:  true,
	}
	testutils.MustExec(t, testutils.DB.Save(&deletedNote), "preparing deletedNote")

	getURL := func(noteUUID string, target testutils.EndpointType) string {
		if target == testutils.EndpointWeb {
			return fmt.Sprintf("/notes/%s", noteUUID)
		}

		return fmt.Sprintf("/api/v3/notes/%s", noteUUID)
	}

	testutils.RunForWebAndAPI(t, "owner accessing private note", func(t *testing.T, target testutils.EndpointType) {
		// Execute
		url := getURL(publicNote.UUID, target)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		if target == testutils.EndpointAPI {
			var payload presenters.Note
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			var n2Record database.Note
			testutils.MustExec(t, testutils.DB.Where("uuid = ?", publicNote.UUID).First(&n2Record), "finding n2Record")

			expected := getExpectedNotePayload(n2Record, b1, user)
			assert.DeepEqual(t, payload, expected, "payload mismatch")
		}
	})

	testutils.RunForWebAndAPI(t, "owner accessing public note", func(t *testing.T, target testutils.EndpointType) {
		// Execute
		url := getURL(publicNote.UUID, target)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		if target == testutils.EndpointAPI {
			var payload presenters.Note
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			var n2Record database.Note
			testutils.MustExec(t, testutils.DB.Where("uuid = ?", publicNote.UUID).First(&n2Record), "finding n2Record")

			expected := getExpectedNotePayload(n2Record, b1, user)
			assert.DeepEqual(t, payload, expected, "payload mismatch")
		}
	})

	testutils.RunForWebAndAPI(t, "non-owner accessing public note", func(t *testing.T, target testutils.EndpointType) {
		// Execute
		url := getURL(publicNote.UUID, target)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, req, anotherUser)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		if target == testutils.EndpointAPI {
			var payload presenters.Note
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			var n2Record database.Note
			testutils.MustExec(t, testutils.DB.Where("uuid = ?", publicNote.UUID).First(&n2Record), "finding n2Record")

			expected := getExpectedNotePayload(n2Record, b1, user)
			assert.DeepEqual(t, payload, expected, "payload mismatch")
		}
	})

	testutils.RunForWebAndAPI(t, "non-owner accessing private note", func(t *testing.T, target testutils.EndpointType) {
		// Execute
		url := getURL(privateNote.UUID, target)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, req, anotherUser)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

		if target == testutils.EndpointAPI {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(errors.Wrap(err, "reading body"))
			}

			assert.DeepEqual(t, string(body), "not found\n", "payload mismatch")
		}
	})

	testutils.RunForWebAndAPI(t, "guest accessing public note", func(t *testing.T, target testutils.EndpointType) {
		// Execute
		url := getURL(publicNote.UUID, target)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		if target == testutils.EndpointAPI {
			var payload presenters.Note
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			var n2Record database.Note
			testutils.MustExec(t, testutils.DB.Where("uuid = ?", publicNote.UUID).First(&n2Record), "finding n2Record")

			expected := getExpectedNotePayload(n2Record, b1, user)
			assert.DeepEqual(t, payload, expected, "payload mismatch")
		}
	})

	testutils.RunForWebAndAPI(t, "guest accessing private note", func(t *testing.T, target testutils.EndpointType) {
		// Execute
		url := getURL(privateNote.UUID, target)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

		if target == testutils.EndpointAPI {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(errors.Wrap(err, "reading body"))
			}

			assert.DeepEqual(t, string(body), "not found\n", "payload mismatch")
		}
	})

	testutils.RunForWebAndAPI(t, "nonexistent", func(t *testing.T, target testutils.EndpointType) {
		// Execute
		url := getURL("somerandomstring", target)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

		if target == testutils.EndpointAPI {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(errors.Wrap(err, "reading body"))
			}

			assert.DeepEqual(t, string(body), "not found\n", "payload mismatch")
		}
	})

	testutils.RunForWebAndAPI(t, "deleted", func(t *testing.T, target testutils.EndpointType) {
		// Execute
		url := getURL(deletedNote.UUID, target)
		req := testutils.MakeReq(server.URL, "GET", url, "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

		if target == testutils.EndpointAPI {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(errors.Wrap(err, "reading body"))
			}

			assert.DeepEqual(t, string(body), "not found\n", "payload mismatch")
		}
	})
}

func TestCreateNote(t *testing.T) {
	testutils.RunForWebAndAPI(t, "success", func(t *testing.T, target testutils.EndpointType) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
		defer server.Close()

		user := testutils.SetupUserData()
		testutils.MustExec(t, testutils.DB.Model(&user).Update("max_usn", 101), "preparing user max_usn")

		b1 := database.Book{
			UserID: user.ID,
			Label:  "js",
			USN:    58,
		}
		testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")

		// Execute

		var req *http.Request
		if target == testutils.EndpointAPI {
			dat := fmt.Sprintf(`{"book_uuid": "%s", "content": "note content"}`, b1.UUID)
			req = testutils.MakeReq(server.URL, "POST", "/api/v3/notes", dat)
		} else {
			dat := url.Values{}
			dat.Set("book_uuid", b1.UUID)
			dat.Set("content", "note content")
			req = testutils.MakeFormReq(server.URL, "POST", "/notes", dat)
		}
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusCreated, "")

		var noteRecord database.Note
		var bookRecord database.Book
		var userRecord database.User
		var bookCount, noteCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), "counting books")
		testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), "counting notes")
		testutils.MustExec(t, testutils.DB.First(&noteRecord), "finding note")
		testutils.MustExec(t, testutils.DB.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
		testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), "finding user record")

		assert.Equalf(t, bookCount, 1, "book count mismatch")
		assert.Equalf(t, noteCount, 1, "note count mismatch")

		assert.Equal(t, bookRecord.Label, b1.Label, "book name mismatch")
		assert.Equal(t, bookRecord.UUID, b1.UUID, "book uuid mismatch")
		assert.Equal(t, bookRecord.UserID, b1.UserID, "book user_id mismatch")
		assert.Equal(t, bookRecord.USN, 58, "book usn mismatch")

		assert.NotEqual(t, noteRecord.UUID, "", "note uuid should have been generated")
		assert.Equal(t, noteRecord.BookUUID, b1.UUID, "note book_uuid mismatch")
		assert.Equal(t, noteRecord.Body, "note content", "note content mismatch")
		assert.Equal(t, noteRecord.USN, 102, "note usn mismatch")
	})
}

func TestDeleteNote(t *testing.T) {
	testutils.RunForWebAndAPI(t, "success", func(t *testing.T, target testutils.EndpointType) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
		defer server.Close()

		user := testutils.SetupUserData()
		testutils.MustExec(t, testutils.DB.Model(&user).Update("max_usn", 101), "preparing user max_usn")

		b1 := database.Book{
			UserID: user.ID,
			Label:  "js",
			USN:    58,
		}
		testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
		n1 := database.Note{
			UserID:   user.ID,
			BookUUID: b1.UUID,
			Body:     "n1 content",
			USN:      11,
			Deleted:  false,
			AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
		}
		testutils.MustExec(t, testutils.DB.Save(&n1), "preparing n1")

		fmt.Println(n1.UUID)

		// Execute
		var req *http.Request
		if target == testutils.EndpointAPI {
			endpoint := fmt.Sprintf("/api/v3/notes/%s", n1.UUID)
			req = testutils.MakeReq(server.URL, "DELETE", endpoint, "")
		} else {
			endpoint := fmt.Sprintf("/notes/%s", n1.UUID)
			req = testutils.MakeFormReq(server.URL, "DELETE", endpoint, nil)
		}
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var noteRecord database.Note
		var bookRecord database.Book
		var userRecord database.User
		var bookCount, noteCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), "counting books")
		testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), "counting notes")
		testutils.MustExec(t, testutils.DB.First(&noteRecord), "finding note")
		testutils.MustExec(t, testutils.DB.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
		testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), "finding user record")

		assert.Equalf(t, bookCount, 1, "book count mismatch")
		assert.Equalf(t, noteCount, 1, "note count mismatch")

		assert.Equal(t, bookRecord.Label, b1.Label, "book name mismatch")
		assert.Equal(t, bookRecord.UUID, b1.UUID, "book uuid mismatch")
		assert.Equal(t, bookRecord.UserID, b1.UserID, "book user_id mismatch")
		assert.Equal(t, bookRecord.USN, 58, "book usn mismatch")

		assert.NotEqual(t, noteRecord.UUID, "", "note uuid should have been generated")
		assert.Equal(t, noteRecord.BookUUID, b1.UUID, "note book_uuid mismatch")
		assert.Equal(t, noteRecord.Body, "", "note content mismatch")
		assert.Equal(t, noteRecord.USN, 102, "note usn mismatch")
		assert.Equal(t, noteRecord.Deleted, true, "note usn mismatch")
	})
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
		testutils.RunForWebAndAPI(t, fmt.Sprintf("test case %d", idx), func(t *testing.T, target testutils.EndpointType) {
			defer testutils.ClearData(testutils.DB)

			// Setup
			server := MustNewServer(t, &app.App{
				Clock: clock.NewMock(),
				Config: config.Config{
					PageTemplateDir: "../views",
				},
			})
			defer server.Close()

			user := testutils.SetupUserData()
			testutils.MustExec(t, testutils.DB.Model(&user).Update("max_usn", 101), "preparing user max_usn")

			b1 := database.Book{
				UUID:   b1UUID,
				UserID: user.ID,
				Label:  "css",
			}
			testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
			b2 := database.Book{
				UUID:   b2UUID,
				UserID: user.ID,
				Label:  "js",
			}
			testutils.MustExec(t, testutils.DB.Save(&b2), "preparing b2")

			note := database.Note{
				UserID:   user.ID,
				UUID:     tc.noteUUID,
				BookUUID: tc.noteBookUUID,
				Body:     tc.noteBody,
				Deleted:  tc.noteDeleted,
				Public:   tc.notePublic,
			}
			testutils.MustExec(t, testutils.DB.Save(&note), "preparing note")

			// Execute
			var req *http.Request

			fmt.Println("URLVALUES")
			fmt.Println(tc.payload.ToURLValues().Get("book_uuid"))
			fmt.Println("JSONVALUES")
			fmt.Println(tc.payload.ToJSON(t))
			if target == testutils.EndpointWeb {
				endpoint := fmt.Sprintf("/notes/%s", note.UUID)
				req = testutils.MakeFormReq(server.URL, "PATCH", endpoint, tc.payload.ToURLValues())
			} else {
				endpoint := fmt.Sprintf("/api/v3/notes/%s", note.UUID)
				req = testutils.MakeReq(server.URL, "PATCH", endpoint, tc.payload.ToJSON(t))
			}

			res := testutils.HTTPAuthDo(t, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, "status code mismatch for test case")

			var bookRecord database.Book
			var noteRecord database.Note
			var userRecord database.User
			var noteCount, bookCount int
			testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), "counting books")
			testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), "counting notes")
			testutils.MustExec(t, testutils.DB.Where("uuid = ?", note.UUID).First(&noteRecord), "finding note")
			testutils.MustExec(t, testutils.DB.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
			testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), "finding user record")

			assert.Equalf(t, bookCount, 2, "book count mismatch")
			assert.Equalf(t, noteCount, 1, "note count mismatch")

			assert.Equal(t, noteRecord.UUID, tc.noteUUID, "note uuid mismatch for test case")
			assert.Equal(t, noteRecord.Body, tc.expectedNoteBody, "note content mismatch for test case")
			assert.Equal(t, noteRecord.BookUUID, tc.expectedNoteBookUUID, "note book_uuid mismatch for test case")
			assert.Equal(t, noteRecord.Public, tc.expectedNotePublic, "note public mismatch for test case")
			assert.Equal(t, noteRecord.USN, 102, "note usn mismatch for test case")

			assert.Equal(t, userRecord.MaxUSN, 102, "user max_usn mismatch for test case")
		})
	}
}
