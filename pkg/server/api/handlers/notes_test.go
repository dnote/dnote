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
	"time"

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

func TestGetNotes(t *testing.T) {
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
		Label:  "css",
	}
	testutils.MustExec(t, db.Save(&b3), "preparing b3")

	n1 := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n1 content",
		USN:      11,
		Deleted:  false,
		AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n1), "preparing n1")
	n2 := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n2 content",
		USN:      14,
		Deleted:  false,
		AddedOn:  time.Date(2018, time.August, 11, 22, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n2), "preparing n2")
	n3 := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n3 content",
		USN:      17,
		Deleted:  false,
		AddedOn:  time.Date(2017, time.January, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n3), "preparing n3")
	n4 := database.Note{
		UserID:   user.ID,
		BookUUID: b2.UUID,
		Body:     "n4 content",
		USN:      18,
		Deleted:  false,
		AddedOn:  time.Date(2018, time.September, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n4), "preparing n4")
	n5 := database.Note{
		UserID:   anotherUser.ID,
		BookUUID: b3.UUID,
		Body:     "n5 content",
		USN:      19,
		Deleted:  false,
		AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n5), "preparing n5")
	n6 := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "",
		USN:      11,
		Deleted:  true,
		AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n6), "preparing n6")

	// Execute
	req := testutils.MakeReq(server, "GET", "/notes?year=2018&month=8", "")
	res := testutils.HTTPAuthDo(t, req, user)

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
			{
				CreatedAt: n2Record.CreatedAt,
				UpdatedAt: n2Record.UpdatedAt,
				UUID:      n2Record.UUID,
				Body:      n2Record.Body,
				AddedOn:   n2Record.AddedOn,
				USN:       n2Record.USN,
				Book: presenters.NoteBook{
					UUID:  b1.UUID,
					Label: b1.Label,
				},
				User: presenters.NoteUser{
					Name: user.Name,
					UUID: user.UUID,
				},
			},
			{
				CreatedAt: n1Record.CreatedAt,
				UpdatedAt: n1Record.UpdatedAt,
				UUID:      n1Record.UUID,
				Body:      n1Record.Body,
				AddedOn:   n1Record.AddedOn,
				USN:       n1Record.USN,
				Book: presenters.NoteBook{
					UUID:  b1.UUID,
					Label: b1.Label,
				},
				User: presenters.NoteUser{
					Name: user.Name,
					UUID: user.UUID,
				},
			},
		},
		Total: 2,
	}

	if ok := reflect.DeepEqual(payload, expected); !ok {
		t.Errorf("Payload does not match.\nActual:   %+v\nExpected: %+v", payload, expected)
	}
}

func TestGetNote(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	user := testutils.SetupUserData()

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

	n1 := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n1 content",
		USN:      1123,
		AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
		Public:   true,
	}
	testutils.MustExec(t, db.Save(&n1), "preparing n1")
	n2 := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n2 content",
		USN:      1888,
		AddedOn:  time.Date(2018, time.August, 11, 22, 0, 0, 0, time.UTC).UnixNano(),
	}
	testutils.MustExec(t, db.Save(&n2), "preparing n2")

	// Execute
	url := fmt.Sprintf("/notes/%s", n1.UUID)
	req := testutils.MakeReq(server, "GET", url, "")
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var payload presenters.Note
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	var n1Record database.Note
	testutils.MustExec(t, db.Where("uuid = ?", n1.UUID).First(&n1Record), "finding n1Record")

	expected := presenters.Note{
		UUID:      n1Record.UUID,
		CreatedAt: n1Record.CreatedAt,
		UpdatedAt: n1Record.UpdatedAt,
		Body:      n1Record.Body,
		AddedOn:   n1Record.AddedOn,
		Public:    n1Record.Public,
		USN:       n1Record.USN,
		Book: presenters.NoteBook{
			UUID:  b1.UUID,
			Label: b1.Label,
		},
		User: presenters.NoteUser{
			Name: user.Name,
			UUID: user.UUID,
		},
	}

	if ok := reflect.DeepEqual(payload, expected); !ok {
		t.Errorf("Payload does not match.\nActual:   %+v\nExpected: %+v", payload, expected)
	}
}

// TODO: finish the test after implementing note sharing
// func TestGetNote_guestAccessPrivate(t *testing.T) {
// 	defer testutils.ClearData()
// 	db := database.DBConn
//
// 	// Setup
// 	server := httptest.NewServer(NewRouter(&App{
// 		Clock: clock.NewMock(),
// 	}))
// 	defer server.Close()
//
// 	user := testutils.SetupUserData()
//
// 	b1 := database.Book{
// 		UUID:   "b1-uuid",
// 		UserID: user.ID,
// 		Label:  "js",
// 	}
// 	testutils.MustExec(t, db.Save(&b1), "preparing b1")
//
// 	n1 := database.Note{
// 		UserID:   user.ID,
// 		BookUUID: b1.UUID,
// 		Body:     "n1 content",
// 		AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
// 		Public:   false,
// 	}
// 	testutils.MustExec(t, db.Save(&n1), "preparing n1")
//
// 	// Execute
// 	url := fmt.Sprintf("/notes/%s", n1.UUID)
// 	req := testutils.MakeReq(server, "GET", url, "")
//
// 	res := testutils.HTTPDo(t, req)
//
// 	// Test
// 	assert.StatusCodeEquals(t, res, http.StatusNotFound, "")
// }

// func TestGetNote_nonOwnerAccessPrivate(t *testing.T) {
// 	defer testutils.ClearData()
// 	db := database.DBConn
//
// 	// Setup
// 	server := httptest.NewServer(NewRouter(&App{
// 		Clock: clock.NewMock(),
// 	}))
// 	defer server.Close()
//
// 	owner := testutils.SetupUserData()
//
// 	nonOwner := testutils.SetupUserData()
// 	testutils.MustExec(t, db.Model(&nonOwner).Update("api_key", "non-owner-api-key"), "preparing user max_usn")
//
// 	b1 := database.Book{
// 		UUID:   "b1-uuid",
// 		UserID: owner.ID,
// 		Label:  "js",
// 	}
// 	testutils.MustExec(t, db.Save(&b1), "preparing b1")
//
// 	n1 := database.Note{
// 		UserID:   owner.ID,
// 		BookUUID: b1.UUID,
// 		Body:     "n1 content",
// 		AddedOn:  time.Date(2018, time.August, 10, 23, 0, 0, 0, time.UTC).UnixNano(),
// 		Public:   false,
// 	}
// 	testutils.MustExec(t, db.Save(&n1), "preparing n1")
//
// 	// Execute
// 	url := fmt.Sprintf("/notes/%s", n1.UUID)
// 	req := testutils.MakeReq(server, "GET", url, "")
// 	res := testutils.HTTPAuthDo(t, req, nonOwner)
//
// 	// Test
// 	assert.StatusCodeEquals(t, res, http.StatusNotFound, "")
// }
