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
	// "time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestGetBooks(t *testing.T) {
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
			UserID:  user.ID,
			Label:   "js",
			USN:     1123,
			Deleted: false,
		}
		testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
		b2 := database.Book{
			UserID:  user.ID,
			Label:   "css",
			USN:     1125,
			Deleted: false,
		}
		testutils.MustExec(t, testutils.DB.Save(&b2), "preparing b2")
		b3 := database.Book{
			UserID:  anotherUser.ID,
			Label:   "css",
			USN:     1128,
			Deleted: false,
		}
		testutils.MustExec(t, testutils.DB.Save(&b3), "preparing b3")
		b4 := database.Book{
			UserID:  user.ID,
			Label:   "",
			USN:     1129,
			Deleted: true,
		}
		testutils.MustExec(t, testutils.DB.Save(&b4), "preparing b4")

		// Execute
		var endpoint string
		if target == testutils.EndpointWeb {
			endpoint = "/books"
		} else {
			endpoint = "/api/v3/books"
		}

		req := testutils.MakeReq(server.URL, "GET", endpoint, "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		if target == testutils.EndpointAPI {
			var payload []presenters.Book
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			var b1Record, b2Record database.Book
			testutils.MustExec(t, testutils.DB.Where("id = ?", b1.ID).First(&b1Record), "finding b1")
			testutils.MustExec(t, testutils.DB.Where("id = ?", b2.ID).First(&b2Record), "finding b2")
			testutils.MustExec(t, testutils.DB.Where("id = ?", b2.ID).First(&b2Record), "finding b2")

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
	})
}

func TestGetBooksByName(t *testing.T) {
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
			Label:  "js",
		}
		testutils.MustExec(t, testutils.DB.Save(&b3), "preparing b3")

		// Execute
		var endpoint string
		if target == testutils.EndpointWeb {
			endpoint = "/books?name=js"
		} else {
			endpoint = "/api/v3/books?name=js"
		}

		req := testutils.MakeReq(server.URL, "GET", endpoint, "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		if target == testutils.EndpointAPI {
			var payload []presenters.Book
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			var b1Record database.Book
			testutils.MustExec(t, testutils.DB.Where("id = ?", b1.ID).First(&b1Record), "finding b1")

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
	})
}

func TestGetBook(t *testing.T) {
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
		Label:  "js",
	}
	testutils.MustExec(t, testutils.DB.Save(&b3), "preparing b3")

	// Execute
	endpoint := fmt.Sprintf("/api/v3/books/%s", b1.UUID)
	req := testutils.MakeReq(server.URL, "GET", endpoint, "")
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var payload presenters.Book
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	var b1Record database.Book
	testutils.MustExec(t, testutils.DB.Where("id = ?", b1.ID).First(&b1Record), "finding b1")

	expected := presenters.Book{
		UUID:      b1Record.UUID,
		CreatedAt: b1Record.CreatedAt,
		UpdatedAt: b1Record.UpdatedAt,
		Label:     b1Record.Label,
		USN:       b1Record.USN,
	}

	assert.DeepEqual(t, payload, expected, "payload mismatch")
}

func TestGetBookNonOwner(t *testing.T) {
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
	nonOwner := testutils.SetupUserData()

	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")

	// Execute
	endpoint := fmt.Sprintf("/api/v3/books/%s", b1.UUID)
	req := testutils.MakeReq(server.URL, "GET", endpoint, "")
	res := testutils.HTTPAuthDo(t, req, nonOwner)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusNotFound, "")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(errors.Wrap(err, "reading body"))
	}
	assert.DeepEqual(t, string(body), "", "payload mismatch")
}

func TestCreateBook(t *testing.T) {
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

		var req *http.Request
		if target == testutils.EndpointWeb {
			dat := url.Values{}
			dat.Set("name", "js")
			req = testutils.MakeFormReq(server.URL, "POST", "/books", dat)
		} else {
			req = testutils.MakeReq(server.URL, "POST", "/api/v3/books", `{"name": "js"}`)
		}

		// Execute
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusCreated, "")

		var bookRecord database.Book
		var userRecord database.User
		var bookCount, noteCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), "counting books")
		testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), "counting notes")
		testutils.MustExec(t, testutils.DB.First(&bookRecord), "finding book")
		testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), "finding user record")

		maxUSN := 102

		assert.Equalf(t, bookCount, 1, "book count mismatch")
		assert.Equalf(t, noteCount, 0, "note count mismatch")

		assert.NotEqual(t, bookRecord.UUID, "", "book uuid should have been generated")
		assert.Equal(t, bookRecord.Label, "js", "book name mismatch")
		assert.Equal(t, bookRecord.UserID, user.ID, "book user_id mismatch")
		assert.Equal(t, bookRecord.USN, maxUSN, "book user_id mismatch")
		assert.Equal(t, userRecord.MaxUSN, maxUSN, "user max_usn mismatch")

		if target == testutils.EndpointAPI {
			var got createBookResp
			if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
				t.Fatal(errors.Wrap(err, "decoding"))
			}
			expected := createBookResp{
				Book: presenters.Book{
					UUID:      bookRecord.UUID,
					USN:       bookRecord.USN,
					CreatedAt: bookRecord.CreatedAt,
					UpdatedAt: bookRecord.UpdatedAt,
					Label:     "js",
				},
			}

			assert.DeepEqual(t, got, expected, "payload mismatch")
		}
	})

	testutils.RunForWebAndAPI(t, "duplicate", func(t *testing.T, target testutils.EndpointType) {
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
		testutils.MustExec(t, testutils.DB.Save(&b1), "preparing book data")

		// Execute
		var req *http.Request
		if target == testutils.EndpointWeb {
			dat := url.Values{}
			dat.Set("name", "js")
			req = testutils.MakeFormReq(server.URL, "POST", "/books", dat)
		} else {
			req = testutils.MakeReq(server.URL, "POST", "/api/v3/books", `{"name": "js"}`)
		}
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusConflict, "")

		var bookRecord database.Book
		var bookCount, noteCount int
		var userRecord database.User
		testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), "counting books")
		testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), "counting notes")
		testutils.MustExec(t, testutils.DB.First(&bookRecord), "finding book")
		testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), "finding user record")

		assert.Equalf(t, bookCount, 1, "book count mismatch")
		assert.Equalf(t, noteCount, 0, "note count mismatch")

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
				UUID:    tc.bookUUID,
				UserID:  user.ID,
				Label:   tc.bookLabel,
				Deleted: tc.bookDeleted,
			}
			testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
			b2 := database.Book{
				UUID:   b2UUID,
				UserID: user.ID,
				Label:  "js",
			}
			testutils.MustExec(t, testutils.DB.Save(&b2), "preparing b2")

			// Execute
			var req *http.Request
			if target == testutils.EndpointWeb {
				endpoint := fmt.Sprintf("/books/%s", tc.bookUUID)
				req = testutils.MakeFormReq(server.URL, "PATCH", endpoint, tc.payload.ToURLValues())
			} else {
				endpoint := fmt.Sprintf("/api/v3/books/%s", tc.bookUUID)
				req = testutils.MakeReq(server.URL, "PATCH", endpoint, tc.payload.ToJSON(t))
			}
			res := testutils.HTTPAuthDo(t, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, fmt.Sprintf("status code mismatch for test case %d", idx))

			var bookRecord database.Book
			var userRecord database.User
			var noteCount, bookCount int
			testutils.MustExec(t, testutils.DB.Model(&database.Book{}).Count(&bookCount), "counting books")
			testutils.MustExec(t, testutils.DB.Model(&database.Note{}).Count(&noteCount), "counting notes")
			testutils.MustExec(t, testutils.DB.Where("id = ?", b1.ID).First(&bookRecord), "finding book")
			testutils.MustExec(t, testutils.DB.Where("id = ?", user.ID).First(&userRecord), "finding user record")

			assert.Equalf(t, bookCount, 2, "book count mismatch")
			assert.Equalf(t, noteCount, 0, "note count mismatch")

			assert.Equalf(t, bookRecord.UUID, tc.bookUUID, "book uuid mismatch")
			assert.Equalf(t, bookRecord.Label, tc.expectedBookLabel, "book label mismatch")
			assert.Equalf(t, bookRecord.USN, 102, "book usn mismatch")
			assert.Equalf(t, bookRecord.Deleted, false, "book Deleted mismatch")

			assert.Equal(t, userRecord.MaxUSN, 102, fmt.Sprintf("user max_usn mismatch for test case %d", idx))
		})
	}
}
