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
