/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package sync

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/client"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/testutils"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
)

var dbPath = "../../tmp/.dnote.db"

func TestProcessFragments(t *testing.T) {
	fragments := []client.SyncFragment{
		{
			FragMaxUSN:  10,
			UserMaxUSN:  10,
			CurrentTime: 1550436136,
			Notes: []client.SyncFragNote{
				{
					UUID: "45546de0-40ed-45cf-9bfc-62ce729a7d3d",
					Body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\n Donec ac libero efficitur, posuere dui non, egestas lectus.\n Aliquam urna ligula, sagittis eu volutpat vel, consequat et augue.\n\n Ut mi urna, dignissim a ex eget, venenatis accumsan sem. Praesent facilisis, ligula hendrerit auctor varius, mauris metus hendrerit dolor, sit amet pulvinar.",
				},
				{
					UUID: "a25a5336-afe9-46c4-b881-acab911c0bc3",
					Body: "foo bar baz quz\nqux",
				},
			},
			Books: []client.SyncFragBook{
				{
					UUID:  "e8ac6f25-d95b-435a-9fae-094f7506a5ac",
					Label: "foo",
				},
				{
					UUID:  "05fd8b95-ddcd-4071-9380-4358ffb8a436",
					Label: "foo-bar-baz-1000",
				},
			},
			ExpungedNotes: []string{},
			ExpungedBooks: []string{},
		},
	}

	// exec
	sl, err := processFragments(fragments)
	if err != nil {
		t.Fatalf(errors.Wrap(err, "executing").Error())
	}

	expected := syncList{
		Notes: map[string]client.SyncFragNote{
			"45546de0-40ed-45cf-9bfc-62ce729a7d3d": {
				UUID: "45546de0-40ed-45cf-9bfc-62ce729a7d3d",
				Body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\n Donec ac libero efficitur, posuere dui non, egestas lectus.\n Aliquam urna ligula, sagittis eu volutpat vel, consequat et augue.\n\n Ut mi urna, dignissim a ex eget, venenatis accumsan sem. Praesent facilisis, ligula hendrerit auctor varius, mauris metus hendrerit dolor, sit amet pulvinar.",
			},
			"a25a5336-afe9-46c4-b881-acab911c0bc3": {
				UUID: "a25a5336-afe9-46c4-b881-acab911c0bc3",
				Body: "foo bar baz quz\nqux",
			},
		},
		Books: map[string]client.SyncFragBook{
			"e8ac6f25-d95b-435a-9fae-094f7506a5ac": {
				UUID:  "e8ac6f25-d95b-435a-9fae-094f7506a5ac",
				Label: "foo",
			},
			"05fd8b95-ddcd-4071-9380-4358ffb8a436": {
				UUID:  "05fd8b95-ddcd-4071-9380-4358ffb8a436",
				Label: "foo-bar-baz-1000",
			},
		},
		ExpungedNotes:  map[string]bool{},
		ExpungedBooks:  map[string]bool{},
		MaxUSN:         10,
		MaxCurrentTime: 1550436136,
	}

	// test
	assert.DeepEqual(t, sl, expected, "syncList mismatch")
}

func TestGetLastSyncAt(t *testing.T) {
	// set up
	db := database.InitTestDB(t, "../../tmp/.dnote", nil)
	defer database.CloseTestDB(t, db)
	database.MustExec(t, "setting up last_sync_at", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastSyncAt, 1541108743)

	// exec
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	got, err := getLastSyncAt(tx)
	if err != nil {
		t.Fatalf(errors.Wrap(err, "getting last_sync_at").Error())
	}

	tx.Commit()

	// test
	assert.Equal(t, got, 1541108743, "last_sync_at mismatch")
}

func TestGetLastMaxUSN(t *testing.T) {
	// set up
	db := database.InitTestDB(t, "../../tmp/.dnote", nil)
	defer database.CloseTestDB(t, db)
	database.MustExec(t, "setting up last_max_usn", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, 20001)

	// exec
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	got, err := getLastMaxUSN(tx)
	if err != nil {
		t.Fatalf(errors.Wrap(err, "getting last_max_usn").Error())
	}

	tx.Commit()

	// test
	assert.Equal(t, got, 20001, "last_max_usn mismatch")
}

func TestResolveLabel(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "js",
			expected: "js_2",
		},
		{
			input:    "css",
			expected: "css_3",
		},
		{
			input:    "linux",
			expected: "linux_4",
		},
		{
			input:    "cool_ideas",
			expected: "cool_ideas_2",
		},
	}

	for idx, tc := range testCases {
		func() {
			// set up
			db := database.InitTestDB(t, "../../tmp/.dnote", nil)
			defer database.CloseTestDB(t, db)

			database.MustExec(t, fmt.Sprintf("inserting book for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", "b1-uuid", "js")
			database.MustExec(t, fmt.Sprintf("inserting book for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", "b2-uuid", "css_2")
			database.MustExec(t, fmt.Sprintf("inserting book for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", "b3-uuid", "linux_(1)")
			database.MustExec(t, fmt.Sprintf("inserting book for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", "b4-uuid", "linux_2")
			database.MustExec(t, fmt.Sprintf("inserting book for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", "b5-uuid", "linux_3")
			database.MustExec(t, fmt.Sprintf("inserting book for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", "b6-uuid", "cool_ideas")

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
			}

			got, err := resolveLabel(tx, tc.input)
			if err != nil {
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
			}
			tx.Rollback()

			assert.Equal(t, got, tc.expected, fmt.Sprintf("output mismatch for test case %d", idx))
		}()
	}
}

func TestSyncDeleteNote(t *testing.T) {
	t.Run("exists on server only", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		if err := syncDeleteNote(tx, "nonexistent-note-uuid"); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 0, "book count mismatch")
	})

	t.Run("local copy is dirty", func(t *testing.T) {
		b1UUID := utils.GenerateUUID()

		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		database.MustExec(t, "inserting b1 for test case %d", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1-label")
		database.MustExec(t, "inserting n1 for test case %d", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", b1UUID, 10, "n1 body", 1541108743, false, true)
		database.MustExec(t, "inserting n2 for test case %d", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n2-uuid", b1UUID, 11, "n2 body", 1541108743, false, true)

		var n1 database.Note
		database.MustScan(t, "getting n1 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", "n1-uuid"),
			&n1.UUID, &n1.BookUUID, &n1.USN, &n1.AddedOn, &n1.EditedOn, &n1.Body, &n1.Deleted, &n1.Dirty)
		var n2 database.Note
		database.MustScan(t, "getting n2 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body, deleted, dirty FROM notes WHERE uuid = ?", "n2-uuid"),
			&n2.UUID, &n2.BookUUID, &n2.USN, &n2.AddedOn, &n2.EditedOn, &n2.Body, &n2.Deleted, &n2.Dirty)

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction for test case").Error())
		}

		if err := syncDeleteNote(tx, "n1-uuid"); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes for test case", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books for test case", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		// do not delete note if local copy is dirty
		assert.Equalf(t, noteCount, 2, "note count mismatch for test case")
		assert.Equalf(t, bookCount, 1, "book count mismatch for test case")

		var n1Record database.Note
		database.MustScan(t, "getting n1 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", n1.UUID),
			&n1Record.UUID, &n1Record.BookUUID, &n1Record.USN, &n1Record.AddedOn, &n1Record.EditedOn, &n1Record.Body, &n1Record.Deleted, &n1Record.Dirty)
		var n2Record database.Note
		database.MustScan(t, "getting n2 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", n2.UUID),
			&n2Record.UUID, &n2Record.BookUUID, &n2Record.USN, &n2Record.AddedOn, &n2Record.EditedOn, &n2Record.Body, &n2Record.Deleted, &n2Record.Dirty)

		assert.Equal(t, n1Record.UUID, n1.UUID, "n1 UUID mismatch for test case")
		assert.Equal(t, n1Record.BookUUID, n1.BookUUID, "n1 BookUUID mismatch for test case")
		assert.Equal(t, n1Record.USN, n1.USN, "n1 USN mismatch for test case")
		assert.Equal(t, n1Record.AddedOn, n1.AddedOn, "n1 AddedOn mismatch for test case")
		assert.Equal(t, n1Record.EditedOn, n1.EditedOn, "n1 EditedOn mismatch for test case")
		assert.Equal(t, n1Record.Body, n1.Body, "n1 Body mismatch for test case")
		assert.Equal(t, n1Record.Deleted, n1.Deleted, "n1 Deleted mismatch for test case")
		assert.Equal(t, n1Record.Dirty, n1.Dirty, "n1 Dirty mismatch for test case")

		assert.Equal(t, n2Record.UUID, n2.UUID, "n2 UUID mismatch for test case")
		assert.Equal(t, n2Record.BookUUID, n2.BookUUID, "n2 BookUUID mismatch for test case")
		assert.Equal(t, n2Record.USN, n2.USN, "n2 USN mismatch for test case")
		assert.Equal(t, n2Record.AddedOn, n2.AddedOn, "n2 AddedOn mismatch for test case")
		assert.Equal(t, n2Record.EditedOn, n2.EditedOn, "n2 EditedOn mismatch for test case")
		assert.Equal(t, n2Record.Body, n2.Body, "n2 Body mismatch for test case")
		assert.Equal(t, n2Record.Deleted, n2.Deleted, "n2 Deleted mismatch for test case")
		assert.Equal(t, n2Record.Dirty, n2.Dirty, "n2 Dirty mismatch for test case")
	})

	t.Run("local copy is not dirty", func(t *testing.T) {
		b1UUID := utils.GenerateUUID()

		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		database.MustExec(t, "inserting b1 for test case %d", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1-label")
		database.MustExec(t, "inserting n1 for test case %d", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", b1UUID, 10, "n1 body", 1541108743, false, false)
		database.MustExec(t, "inserting n2 for test case %d", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n2-uuid", b1UUID, 11, "n2 body", 1541108743, false, false)

		var n1 database.Note
		database.MustScan(t, "getting n1 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", "n1-uuid"),
			&n1.UUID, &n1.BookUUID, &n1.USN, &n1.AddedOn, &n1.EditedOn, &n1.Body, &n1.Deleted, &n1.Dirty)
		var n2 database.Note
		database.MustScan(t, "getting n2 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", "n2-uuid"),
			&n2.UUID, &n2.BookUUID, &n2.USN, &n2.AddedOn, &n2.EditedOn, &n2.Body, &n2.Deleted, &n2.Dirty)

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction for test case").Error())
		}

		if err := syncDeleteNote(tx, "n1-uuid"); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes for test case", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books for test case", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 1, "note count mismatch for test case")
		assert.Equalf(t, bookCount, 1, "book count mismatch for test case")

		var n2Record database.Note
		database.MustScan(t, "getting n2 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", n2.UUID),
			&n2Record.UUID, &n2Record.BookUUID, &n2Record.USN, &n2Record.AddedOn, &n2Record.EditedOn, &n2Record.Body, &n2Record.Deleted, &n2Record.Dirty)

		assert.Equal(t, n2Record.UUID, n2.UUID, "n2 UUID mismatch for test case")
		assert.Equal(t, n2Record.BookUUID, n2.BookUUID, "n2 BookUUID mismatch for test case")
		assert.Equal(t, n2Record.USN, n2.USN, "n2 USN mismatch for test case")
		assert.Equal(t, n2Record.AddedOn, n2.AddedOn, "n2 AddedOn mismatch for test case")
		assert.Equal(t, n2Record.EditedOn, n2.EditedOn, "n2 EditedOn mismatch for test case")
		assert.Equal(t, n2Record.Body, n2.Body, "n2 Body mismatch for test case")
		assert.Equal(t, n2Record.Deleted, n2.Deleted, "n2 Deleted mismatch for test case")
		assert.Equal(t, n2Record.Dirty, n2.Dirty, "n2 Dirty mismatch for test case")
	})
}

func TestSyncDeleteBook(t *testing.T) {
	t.Run("exists on server only", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)
		database.MustExec(t, "inserting b1 for test case %d", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", "b1-uuid", "b1-label")

		var b1 database.Book
		database.MustScan(t, "getting b1 for test case",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b1-uuid"),
			&b1.UUID, &b1.Label, &b1.USN, &b1.Dirty)

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		if err := syncDeleteBook(tx, "nonexistent-book-uuid"); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 1, "book count mismatch")

		var b1Record database.Book
		database.MustScan(t, "getting b1 for test case",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b1-uuid"),
			&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)

		assert.Equal(t, b1Record.UUID, b1.UUID, "b1 UUID mismatch for test case")
		assert.Equal(t, b1Record.Label, b1.Label, "b1 Label mismatch for test case")
		assert.Equal(t, b1Record.USN, b1.USN, "b1 USN mismatch for test case")
		assert.Equal(t, b1Record.Dirty, b1.Dirty, "b1 Dirty mismatch for test case")
	})

	t.Run("local copy is dirty", func(t *testing.T) {
		b1UUID := utils.GenerateUUID()

		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		database.MustExec(t, "inserting b1 for test case %d", db, "INSERT INTO books (uuid, label, usn, dirty) VALUES (?, ?, ?, ?)", b1UUID, "b1-label", 12, true)
		database.MustExec(t, "inserting n1 for test case %d", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", b1UUID, 10, "n1 body", 1541108743, false, true)

		var b1 database.Book
		database.MustScan(t, "getting b1 for test case",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", b1UUID),
			&b1.UUID, &b1.Label, &b1.USN, &b1.Dirty)
		var n1 database.Note
		database.MustScan(t, "getting n1 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", "n1-uuid"),
			&n1.UUID, &n1.BookUUID, &n1.USN, &n1.AddedOn, &n1.EditedOn, &n1.Body, &n1.Deleted, &n1.Dirty)

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction for test case").Error())
		}

		if err := syncDeleteBook(tx, b1UUID); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes for test case", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books for test case", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		// do not delete note if local copy is dirty
		assert.Equalf(t, noteCount, 1, "note count mismatch for test case")
		assert.Equalf(t, bookCount, 1, "book count mismatch for test case")

		var b1Record database.Book
		database.MustScan(t, "getting b1Record for test case",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", b1UUID),
			&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)
		var n1Record database.Note
		database.MustScan(t, "getting n1 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body, deleted, dirty FROM notes WHERE uuid = ?", n1.UUID),
			&n1Record.UUID, &n1Record.BookUUID, &n1Record.USN, &n1Record.AddedOn, &n1Record.EditedOn, &n1Record.Body, &n1Record.Deleted, &n1Record.Dirty)

		assert.Equal(t, b1Record.UUID, b1.UUID, "b1 UUID mismatch for test case")
		assert.Equal(t, b1Record.Label, b1.Label, "b1 Label mismatch for test case")
		assert.Equal(t, b1Record.USN, b1.USN, "b1 USN mismatch for test case")
		assert.Equal(t, b1Record.Dirty, b1.Dirty, "b1 Dirty mismatch for test case")

		assert.Equal(t, n1Record.UUID, n1.UUID, "n1 UUID mismatch for test case")
		assert.Equal(t, n1Record.BookUUID, n1.BookUUID, "n1 BookUUID mismatch for test case")
		assert.Equal(t, n1Record.USN, n1.USN, "n1 USN mismatch for test case")
		assert.Equal(t, n1Record.AddedOn, n1.AddedOn, "n1 AddedOn mismatch for test case")
		assert.Equal(t, n1Record.EditedOn, n1.EditedOn, "n1 EditedOn mismatch for test case")
		assert.Equal(t, n1Record.Body, n1.Body, "n1 Body mismatch for test case")
		assert.Equal(t, n1Record.Deleted, n1.Deleted, "n1 Deleted mismatch for test case")
		assert.Equal(t, n1Record.Dirty, n1.Dirty, "n1 Dirty mismatch for test case")
	})

	t.Run("local copy is not dirty", func(t *testing.T) {
		b1UUID := utils.GenerateUUID()
		b2UUID := utils.GenerateUUID()

		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		database.MustExec(t, "inserting b1 for test case %d", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1-label")
		database.MustExec(t, "inserting n1 for test case %d", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", b1UUID, 10, "n1 body", 1541108743, false, false)
		database.MustExec(t, "inserting b2 for test case %d", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "b2-label")
		database.MustExec(t, "inserting n2 for test case %d", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n2-uuid", b2UUID, 11, "n2 body", 1541108743, false, false)

		var b2 database.Book
		database.MustScan(t, "getting b2 for test case",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", b2UUID),
			&b2.UUID, &b2.Label, &b2.USN, &b2.Dirty)
		var n2 database.Note
		database.MustScan(t, "getting n2 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", "n2-uuid"),
			&n2.UUID, &n2.BookUUID, &n2.USN, &n2.AddedOn, &n2.EditedOn, &n2.Body, &n2.Deleted, &n2.Dirty)

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction for test case").Error())
		}

		if err := syncDeleteBook(tx, b1UUID); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes for test case", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books for test case", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 1, "note count mismatch for test case")
		assert.Equalf(t, bookCount, 1, "book count mismatch for test case")

		var b2Record database.Book
		database.MustScan(t, "getting b2 for test case",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", b2UUID),
			&b2Record.UUID, &b2Record.Label, &b2Record.USN, &b2Record.Dirty)
		var n2Record database.Note
		database.MustScan(t, "getting n2 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", n2.UUID),
			&n2Record.UUID, &n2Record.BookUUID, &n2Record.USN, &n2Record.AddedOn, &n2Record.EditedOn, &n2Record.Body, &n2Record.Deleted, &n2Record.Dirty)

		assert.Equal(t, b2Record.UUID, b2.UUID, "b2 UUID mismatch for test case")
		assert.Equal(t, b2Record.Label, b2.Label, "b2 Label mismatch for test case")
		assert.Equal(t, b2Record.USN, b2.USN, "b2 USN mismatch for test case")
		assert.Equal(t, b2Record.Dirty, b2.Dirty, "b2 Dirty mismatch for test case")

		assert.Equal(t, n2Record.UUID, n2.UUID, "n2 UUID mismatch for test case")
		assert.Equal(t, n2Record.BookUUID, n2.BookUUID, "n2 BookUUID mismatch for test case")
		assert.Equal(t, n2Record.USN, n2.USN, "n2 USN mismatch for test case")
		assert.Equal(t, n2Record.AddedOn, n2.AddedOn, "n2 AddedOn mismatch for test case")
		assert.Equal(t, n2Record.EditedOn, n2.EditedOn, "n2 EditedOn mismatch for test case")
		assert.Equal(t, n2Record.Body, n2.Body, "n2 Body mismatch for test case")
		assert.Equal(t, n2Record.Deleted, n2.Deleted, "n2 Deleted mismatch for test case")
		assert.Equal(t, n2Record.Dirty, n2.Dirty, "n2 Dirty mismatch for test case")
	})

	t.Run("local copy has at least one note that is dirty", func(t *testing.T) {
		b1UUID := utils.GenerateUUID()

		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		database.MustExec(t, "inserting b1 for test case %d", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1-label")
		database.MustExec(t, "inserting n1 for test case %d", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", b1UUID, 10, "n1 body", 1541108743, false, true)

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction for test case").Error())
		}

		if err := syncDeleteBook(tx, b1UUID); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes for test case", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books for test case", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
		assert.Equalf(t, noteCount, 1, "note count mismatch for test case")
		assert.Equalf(t, bookCount, 1, "book count mismatch for test case")

		var b1Record database.Book
		database.MustScan(t, "getting b1 for test case",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", b1UUID),
			&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)
		var n1Record database.Note
		database.MustScan(t, "getting n1 for test case",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, body,deleted, dirty FROM notes WHERE uuid = ?", "n1-uuid"),
			&n1Record.UUID, &n1Record.BookUUID, &n1Record.USN, &n1Record.AddedOn, &n1Record.Body, &n1Record.Deleted, &n1Record.Dirty)

		assert.Equal(t, b1Record.UUID, b1UUID, "b1 UUID mismatch for test case")
		assert.Equal(t, b1Record.Label, "b1-label", "b1 Label mismatch for test case")
		assert.Equal(t, b1Record.Dirty, true, "b1 Dirty mismatch for test case")

		assert.Equal(t, n1Record.UUID, "n1-uuid", "n1 UUID mismatch for test case")
		assert.Equal(t, n1Record.BookUUID, b1UUID, "n1 BookUUID mismatch for test case")
		assert.Equal(t, n1Record.USN, 10, "n1 USN mismatch for test case")
		assert.Equal(t, n1Record.AddedOn, int64(1541108743), "n1 AddedOn mismatch for test case")
		assert.Equal(t, n1Record.Body, "n1 body", "n1 Body mismatch for test case")
		assert.Equal(t, n1Record.Deleted, false, "n1 Deleted mismatch for test case")
		assert.Equal(t, n1Record.Dirty, true, "n1 Dirty mismatch for test case")
	})
}

func TestFullSyncNote(t *testing.T) {
	t.Run("exists on server only", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		b1UUID := utils.GenerateUUID()
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1-label")

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		n := client.SyncFragNote{
			UUID:     "n1-uuid",
			BookUUID: b1UUID,
			USN:      128,
			AddedOn:  1541232118,
			EditedOn: 1541219321,
			Body:     "n1-body",
			Deleted:  false,
		}

		if err := fullSyncNote(tx, n); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 1, "note count mismatch")
		assert.Equalf(t, bookCount, 1, "book count mismatch")

		var n1 database.Note
		database.MustScan(t, "getting n1",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", n.UUID),
			&n1.UUID, &n1.BookUUID, &n1.USN, &n1.AddedOn, &n1.EditedOn, &n1.Body, &n1.Deleted, &n1.Dirty)

		assert.Equal(t, n1.UUID, n.UUID, "n1 UUID mismatch")
		assert.Equal(t, n1.BookUUID, n.BookUUID, "n1 BookUUID mismatch")
		assert.Equal(t, n1.USN, n.USN, "n1 USN mismatch")
		assert.Equal(t, n1.AddedOn, n.AddedOn, "n1 AddedOn mismatch")
		assert.Equal(t, n1.EditedOn, n.EditedOn, "n1 EditedOn mismatch")
		assert.Equal(t, n1.Body, n.Body, "n1 Body mismatch")
		assert.Equal(t, n1.Deleted, n.Deleted, "n1 Deleted mismatch")
		assert.Equal(t, n1.Dirty, false, "n1 Dirty mismatch")
	})

	t.Run("exists on server and client", func(t *testing.T) {
		b1UUID := utils.GenerateUUID()
		b2UUID := utils.GenerateUUID()
		conflictBookUUID := utils.GenerateUUID()

		testCases := []struct {
			addedOn          int64
			clientUSN        int
			clientEditedOn   int64
			clientBody       string
			clientDeleted    bool
			clientBookUUID   string
			clientDirty      bool
			serverUSN        int
			serverEditedOn   int64
			serverBody       string
			serverDeleted    bool
			serverBookUUID   string
			expectedUSN      int
			expectedAddedOn  int64
			expectedEditedOn int64
			expectedBody     string
			expectedDeleted  bool
			expectedBookUUID string
			expectedDirty    bool
		}{
			// server has higher usn and client is dirty
			{
				clientDirty:      true,
				clientUSN:        1,
				clientEditedOn:   0,
				clientBody:       "n1 body",
				clientDeleted:    false,
				clientBookUUID:   b1UUID,
				addedOn:          1541232118,
				serverUSN:        21,
				serverEditedOn:   1541219321,
				serverBody:       "n1 body edited",
				serverDeleted:    false,
				serverBookUUID:   b2UUID,
				expectedUSN:      21,
				expectedAddedOn:  1541232118,
				expectedEditedOn: 1541219321,
				expectedBody: `<<<<<<< Local
Moved to the book b1-label
=======
Moved to the book b2-label
>>>>>>> Server

<<<<<<< Local
n1 body
=======
n1 body edited
>>>>>>> Server
`,
				expectedDeleted:  false,
				expectedBookUUID: conflictBookUUID,
				expectedDirty:    true,
			},
			{
				clientDirty:      true,
				clientUSN:        1,
				clientEditedOn:   0,
				clientBody:       "n1 body",
				clientDeleted:    false,
				clientBookUUID:   b1UUID,
				addedOn:          1541232118,
				serverUSN:        21,
				serverEditedOn:   1541219321,
				serverBody:       "n1 body edited",
				serverDeleted:    false,
				serverBookUUID:   b1UUID,
				expectedUSN:      21,
				expectedAddedOn:  1541232118,
				expectedEditedOn: 1541219321,
				expectedBody: `<<<<<<< Local
n1 body
=======
n1 body edited
>>>>>>> Server
`,
				expectedDeleted:  false,
				expectedBookUUID: b1UUID,
				expectedDirty:    true,
			},
			// server has higher usn and client deleted locally
			{
				clientDirty:      true,
				clientUSN:        1,
				clientEditedOn:   0,
				clientBody:       "",
				clientDeleted:    true,
				clientBookUUID:   b1UUID,
				addedOn:          1541232118,
				serverUSN:        21,
				serverEditedOn:   1541219321,
				serverBody:       "n1 body server",
				serverDeleted:    false,
				serverBookUUID:   b2UUID,
				expectedUSN:      21,
				expectedAddedOn:  1541232118,
				expectedEditedOn: 1541219321,
				expectedBody:     "n1 body server",
				expectedDeleted:  false,
				expectedBookUUID: b2UUID,
				expectedDirty:    false,
			},
			// server has higher usn and client is not dirty
			{
				clientDirty:      false,
				clientUSN:        1,
				clientEditedOn:   0,
				clientBody:       "n1 body",
				clientDeleted:    false,
				clientBookUUID:   b1UUID,
				addedOn:          1541232118,
				serverUSN:        21,
				serverEditedOn:   1541219321,
				serverBody:       "n1 body edited",
				serverDeleted:    false,
				serverBookUUID:   b2UUID,
				expectedUSN:      21,
				expectedAddedOn:  1541232118,
				expectedEditedOn: 1541219321,
				expectedBody:     "n1 body edited",
				expectedDeleted:  false,
				expectedBookUUID: b2UUID,
				expectedDirty:    false,
			},
			// they're in sync
			{
				clientDirty:      true,
				clientUSN:        21,
				clientEditedOn:   1541219321,
				clientBody:       "n1 body",
				clientDeleted:    false,
				clientBookUUID:   b2UUID,
				addedOn:          1541232118,
				serverUSN:        21,
				serverEditedOn:   1541219321,
				serverBody:       "n1 body",
				serverDeleted:    false,
				serverBookUUID:   b2UUID,
				expectedUSN:      21,
				expectedAddedOn:  1541232118,
				expectedEditedOn: 1541219321,
				expectedBody:     "n1 body",
				expectedDeleted:  false,
				expectedBookUUID: b2UUID,
				expectedDirty:    true,
			},
			// they have the same usn but client is dirty
			// not sure if this is a possible scenario but if it happens, the local copy will
			// be uploaded to the server anyway.
			{
				clientDirty:      true,
				clientUSN:        21,
				clientEditedOn:   1541219320,
				clientBody:       "n1 body client",
				clientDeleted:    false,
				clientBookUUID:   b2UUID,
				addedOn:          1541232118,
				serverUSN:        21,
				serverEditedOn:   1541219321,
				serverBody:       "n1 body server",
				serverDeleted:    false,
				serverBookUUID:   b2UUID,
				expectedUSN:      21,
				expectedAddedOn:  1541232118,
				expectedEditedOn: 1541219320,
				expectedBody:     "n1 body client",
				expectedDeleted:  false,
				expectedBookUUID: b2UUID,
				expectedDirty:    true,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				db := database.InitTestDB(t, dbPath, nil)
				defer database.CloseTestDB(t, db)

				database.MustExec(t, fmt.Sprintf("inserting b1 for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1-label")
				database.MustExec(t, fmt.Sprintf("inserting b2 for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "b2-label")
				database.MustExec(t, fmt.Sprintf("inserting conflitcs book for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", conflictBookUUID, "conflicts")
				n1UUID := utils.GenerateUUID()
				database.MustExec(t, fmt.Sprintf("inserting n1 for test case %d", idx), db, "INSERT INTO notes (uuid, book_uuid, usn, added_on, edited_on, body, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", n1UUID, tc.clientBookUUID, tc.clientUSN, tc.addedOn, tc.clientEditedOn, tc.clientBody, tc.clientDeleted, tc.clientDirty)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				// update all fields but uuid and bump usn
				n := client.SyncFragNote{
					UUID:     n1UUID,
					BookUUID: tc.serverBookUUID,
					USN:      tc.serverUSN,
					AddedOn:  tc.addedOn,
					EditedOn: tc.serverEditedOn,
					Body:     tc.serverBody,
					Deleted:  tc.serverDeleted,
				}

				if err := fullSyncNote(tx, n); err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				var noteCount, bookCount int
				database.MustScan(t, fmt.Sprintf("counting notes for test case %d", idx), db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
				database.MustScan(t, fmt.Sprintf("counting books for test case %d", idx), db.QueryRow("SELECT count(*) FROM books"), &bookCount)

				assert.Equalf(t, noteCount, 1, fmt.Sprintf("note count mismatch for test case %d", idx))
				assert.Equalf(t, bookCount, 3, fmt.Sprintf("book count mismatch for test case %d", idx))

				var n1 database.Note
				database.MustScan(t, fmt.Sprintf("getting n1 for test case %d", idx),
					db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body, deleted, dirty FROM notes WHERE uuid = ?", n.UUID),
					&n1.UUID, &n1.BookUUID, &n1.USN, &n1.AddedOn, &n1.EditedOn, &n1.Body, &n1.Deleted, &n1.Dirty)

				assert.Equal(t, n1.UUID, n.UUID, fmt.Sprintf("n1 UUID mismatch for test case %d", idx))
				assert.Equal(t, n1.BookUUID, tc.expectedBookUUID, fmt.Sprintf("n1 BookUUID mismatch for test case %d", idx))
				assert.Equal(t, n1.USN, tc.expectedUSN, fmt.Sprintf("n1 USN mismatch for test case %d", idx))
				assert.Equal(t, n1.AddedOn, tc.expectedAddedOn, fmt.Sprintf("n1 AddedOn mismatch for test case %d", idx))
				assert.Equal(t, n1.EditedOn, tc.expectedEditedOn, fmt.Sprintf("n1 EditedOn mismatch for test case %d", idx))
				assert.Equal(t, n1.Body, tc.expectedBody, fmt.Sprintf("n1 Body mismatch for test case %d", idx))
				assert.Equal(t, n1.Deleted, tc.expectedDeleted, fmt.Sprintf("n1 Deleted mismatch for test case %d", idx))
				assert.Equal(t, n1.Dirty, tc.expectedDirty, fmt.Sprintf("n1 Dirty mismatch for test case %d", idx))
			}()
		}
	})
}

func TestFullSyncBook(t *testing.T) {
	t.Run("exists on server only", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		b1UUID := utils.GenerateUUID()
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", b1UUID, 555, "b1-label", true, false)

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		b2UUID := utils.GenerateUUID()
		b := client.SyncFragBook{
			UUID:    b2UUID,
			USN:     1,
			AddedOn: 1541108743,
			Label:   "b2-label",
			Deleted: false,
		}

		if err := fullSyncBook(tx, b); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 2, "book count mismatch")

		var b1, b2 database.Book
		database.MustScan(t, "getting b1",
			db.QueryRow("SELECT uuid, usn, label, dirty, deleted FROM books WHERE uuid = ?", b1UUID),
			&b1.UUID, &b1.USN, &b1.Label, &b1.Dirty, &b1.Deleted)
		database.MustScan(t, "getting b2",
			db.QueryRow("SELECT uuid, usn, label, dirty, deleted FROM books WHERE uuid = ?", b2UUID),
			&b2.UUID, &b2.USN, &b2.Label, &b2.Dirty, &b2.Deleted)

		assert.Equal(t, b1.UUID, b1UUID, "b1 UUID mismatch")
		assert.Equal(t, b1.USN, 555, "b1 USN mismatch")
		assert.Equal(t, b1.Label, "b1-label", "b1 Label mismatch")
		assert.Equal(t, b1.Dirty, true, "b1 Dirty mismatch")
		assert.Equal(t, b1.Deleted, false, "b1 Deleted mismatch")

		assert.Equal(t, b2.UUID, b2UUID, "b2 UUID mismatch")
		assert.Equal(t, b2.USN, b.USN, "b2 USN mismatch")
		assert.Equal(t, b2.Label, b.Label, "b2 Label mismatch")
		assert.Equal(t, b2.Dirty, false, "b2 Dirty mismatch")
		assert.Equal(t, b2.Deleted, b.Deleted, "b2 Deleted mismatch")
	})

	t.Run("exists on server and client", func(t *testing.T) {
		testCases := []struct {
			clientDirty     bool
			clientUSN       int
			clientLabel     string
			clientDeleted   bool
			serverUSN       int
			serverLabel     string
			serverDeleted   bool
			expectedUSN     int
			expectedLabel   string
			expectedDeleted bool
		}{
			// server has higher usn and client is dirty
			{
				clientDirty:     true,
				clientUSN:       1,
				clientLabel:     "b2-label",
				clientDeleted:   false,
				serverUSN:       3,
				serverLabel:     "b2-label-updated",
				serverDeleted:   false,
				expectedUSN:     3,
				expectedLabel:   "b2-label-updated",
				expectedDeleted: false,
			},
			{
				clientDirty:     true,
				clientUSN:       1,
				clientLabel:     "b2-label",
				clientDeleted:   false,
				serverUSN:       3,
				serverLabel:     "",
				serverDeleted:   true,
				expectedUSN:     3,
				expectedLabel:   "",
				expectedDeleted: true,
			},
			// server has higher usn and client is not dirty
			{
				clientDirty:     false,
				clientUSN:       1,
				clientLabel:     "b2-label",
				clientDeleted:   false,
				serverUSN:       3,
				serverLabel:     "b2-label-updated",
				serverDeleted:   false,
				expectedUSN:     3,
				expectedLabel:   "b2-label-updated",
				expectedDeleted: false,
			},
			// they are in sync
			{
				clientDirty:     false,
				clientUSN:       3,
				clientLabel:     "b2-label",
				clientDeleted:   false,
				serverUSN:       3,
				serverLabel:     "b2-label",
				serverDeleted:   false,
				expectedUSN:     3,
				expectedLabel:   "b2-label",
				expectedDeleted: false,
			},
			// they have the same usn but client is dirty
			{
				clientDirty:     true,
				clientUSN:       3,
				clientLabel:     "b2-label-client",
				clientDeleted:   false,
				serverUSN:       3,
				serverLabel:     "b2-label",
				serverDeleted:   false,
				expectedUSN:     3,
				expectedLabel:   "b2-label-client",
				expectedDeleted: false,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				db := database.InitTestDB(t, dbPath, nil)
				defer database.CloseTestDB(t, db)

				b1UUID := utils.GenerateUUID()
				database.MustExec(t, fmt.Sprintf("inserting book for test case %d", idx), db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", b1UUID, tc.clientUSN, tc.clientLabel, tc.clientDirty, tc.clientDeleted)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				// update all fields but uuid and bump usn
				b := client.SyncFragBook{
					UUID:    b1UUID,
					USN:     tc.serverUSN,
					Label:   tc.serverLabel,
					Deleted: tc.serverDeleted,
				}

				if err := fullSyncBook(tx, b); err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				var noteCount, bookCount int
				database.MustScan(t, fmt.Sprintf("counting notes for test case %d", idx), db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
				database.MustScan(t, fmt.Sprintf("counting books for test case %d", idx), db.QueryRow("SELECT count(*) FROM books"), &bookCount)

				assert.Equalf(t, noteCount, 0, fmt.Sprintf("note count mismatch for test case %d", idx))
				assert.Equalf(t, bookCount, 1, fmt.Sprintf("book count mismatch for test case %d", idx))

				var b1 database.Book
				database.MustScan(t, "getting b1",
					db.QueryRow("SELECT uuid, usn, label, dirty, deleted FROM books WHERE uuid = ?", b1UUID),
					&b1.UUID, &b1.USN, &b1.Label, &b1.Dirty, &b1.Deleted)

				assert.Equal(t, b1.UUID, b1UUID, fmt.Sprintf("b1 UUID mismatch for idx %d", idx))
				assert.Equal(t, b1.USN, tc.expectedUSN, fmt.Sprintf("b1 USN mismatch for test case %d", idx))
				assert.Equal(t, b1.Label, tc.expectedLabel, fmt.Sprintf("b1 Label mismatch for test case %d", idx))
				assert.Equal(t, b1.Dirty, tc.clientDirty, fmt.Sprintf("b1 Dirty mismatch for test case %d", idx))
				assert.Equal(t, b1.Deleted, tc.expectedDeleted, fmt.Sprintf("b1 Deleted mismatch for test case %d", idx))
			}()
		}
	})
}

func TestStepSyncNote(t *testing.T) {
	t.Run("exists on server only", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		b1UUID := utils.GenerateUUID()
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1-label")

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		n := client.SyncFragNote{
			UUID:     "n1-uuid",
			BookUUID: b1UUID,
			USN:      128,
			AddedOn:  1541232118,
			EditedOn: 1541219321,
			Body:     "n1-body",
			Deleted:  false,
		}

		if err := stepSyncNote(tx, n); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 1, "note count mismatch")
		assert.Equalf(t, bookCount, 1, "book count mismatch")

		var n1 database.Note
		database.MustScan(t, "getting n1",
			db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body, deleted, dirty FROM notes WHERE uuid = ?", n.UUID),
			&n1.UUID, &n1.BookUUID, &n1.USN, &n1.AddedOn, &n1.EditedOn, &n1.Body, &n1.Deleted, &n1.Dirty)

		assert.Equal(t, n1.UUID, n.UUID, "n1 UUID mismatch")
		assert.Equal(t, n1.BookUUID, n.BookUUID, "n1 BookUUID mismatch")
		assert.Equal(t, n1.USN, n.USN, "n1 USN mismatch")
		assert.Equal(t, n1.AddedOn, n.AddedOn, "n1 AddedOn mismatch")
		assert.Equal(t, n1.EditedOn, n.EditedOn, "n1 EditedOn mismatch")
		assert.Equal(t, n1.Body, n.Body, "n1 Body mismatch")
		assert.Equal(t, n1.Deleted, n.Deleted, "n1 Deleted mismatch")
		assert.Equal(t, n1.Dirty, false, "n1 Dirty mismatch")
	})

	t.Run("exists on server and client", func(t *testing.T) {
		b1UUID := utils.GenerateUUID()
		b2UUID := utils.GenerateUUID()
		conflictBookUUID := utils.GenerateUUID()

		testCases := []struct {
			addedOn          int64
			clientUSN        int
			clientEditedOn   int64
			clientBody       string
			clientDeleted    bool
			clientBookUUID   string
			clientDirty      bool
			serverUSN        int
			serverEditedOn   int64
			serverBody       string
			serverDeleted    bool
			serverBookUUID   string
			expectedUSN      int
			expectedAddedOn  int64
			expectedEditedOn int64
			expectedBody     string
			expectedDeleted  bool
			expectedBookUUID string
			expectedDirty    bool
		}{
			{
				clientDirty:      true,
				clientUSN:        1,
				clientEditedOn:   0,
				clientBody:       "n1 body",
				clientDeleted:    false,
				clientBookUUID:   b1UUID,
				addedOn:          1541232118,
				serverUSN:        21,
				serverEditedOn:   1541219321,
				serverBody:       "n1 body edited",
				serverDeleted:    false,
				serverBookUUID:   b2UUID,
				expectedUSN:      21,
				expectedAddedOn:  1541232118,
				expectedEditedOn: 1541219321,
				expectedBody: `<<<<<<< Local
Moved to the book b1-label
=======
Moved to the book b2-label
>>>>>>> Server

<<<<<<< Local
n1 body
=======
n1 body edited
>>>>>>> Server
`,
				expectedDeleted:  false,
				expectedBookUUID: conflictBookUUID,
				expectedDirty:    true,
			},
			// if deleted locally, resurrect it
			{
				clientDirty:      true,
				clientUSN:        1,
				clientEditedOn:   1541219321,
				clientBody:       "",
				clientDeleted:    true,
				clientBookUUID:   b1UUID,
				addedOn:          1541232118,
				serverUSN:        21,
				serverEditedOn:   1541219321,
				serverBody:       "n1 body edited",
				serverDeleted:    false,
				serverBookUUID:   b2UUID,
				expectedUSN:      21,
				expectedAddedOn:  1541232118,
				expectedEditedOn: 1541219321,
				expectedBody:     "n1 body edited",
				expectedDeleted:  false,
				expectedBookUUID: b2UUID,
				expectedDirty:    false,
			},
			{
				clientDirty:      false,
				clientUSN:        1,
				clientEditedOn:   0,
				clientBody:       "n1 body",
				clientDeleted:    false,
				clientBookUUID:   b1UUID,
				addedOn:          1541232118,
				serverUSN:        21,
				serverEditedOn:   1541219321,
				serverBody:       "n1 body edited",
				serverDeleted:    false,
				serverBookUUID:   b2UUID,
				expectedUSN:      21,
				expectedAddedOn:  1541232118,
				expectedEditedOn: 1541219321,
				expectedBody:     "n1 body edited",
				expectedDeleted:  false,
				expectedBookUUID: b2UUID,
				expectedDirty:    false,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				db := database.InitTestDB(t, dbPath, nil)
				defer database.CloseTestDB(t, db)

				database.MustExec(t, fmt.Sprintf("inserting b1 for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1-label")
				database.MustExec(t, fmt.Sprintf("inserting b2 for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "b2-label")
				database.MustExec(t, fmt.Sprintf("inserting conflitcs book for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", conflictBookUUID, "conflicts")
				n1UUID := utils.GenerateUUID()
				database.MustExec(t, fmt.Sprintf("inserting n1 for test case %d", idx), db, "INSERT INTO notes (uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", n1UUID, tc.clientBookUUID, tc.clientUSN, tc.addedOn, tc.clientEditedOn, tc.clientBody, tc.clientDeleted, tc.clientDirty)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				// update all fields but uuid and bump usn
				n := client.SyncFragNote{
					UUID:     n1UUID,
					BookUUID: tc.serverBookUUID,
					USN:      tc.serverUSN,
					AddedOn:  tc.addedOn,
					EditedOn: tc.serverEditedOn,
					Body:     tc.serverBody,
					Deleted:  tc.serverDeleted,
				}

				if err := stepSyncNote(tx, n); err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				var noteCount, bookCount int
				database.MustScan(t, fmt.Sprintf("counting notes for test case %d", idx), db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
				database.MustScan(t, fmt.Sprintf("counting books for test case %d", idx), db.QueryRow("SELECT count(*) FROM books"), &bookCount)

				assert.Equalf(t, noteCount, 1, fmt.Sprintf("note count mismatch for test case %d", idx))
				assert.Equalf(t, bookCount, 3, fmt.Sprintf("book count mismatch for test case %d", idx))

				var n1 database.Note
				database.MustScan(t, fmt.Sprintf("getting n1 for test case %d", idx),
					db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body,  deleted, dirty FROM notes WHERE uuid = ?", n.UUID),
					&n1.UUID, &n1.BookUUID, &n1.USN, &n1.AddedOn, &n1.EditedOn, &n1.Body, &n1.Deleted, &n1.Dirty)

				assert.Equal(t, n1.UUID, n.UUID, fmt.Sprintf("n1 UUID mismatch for test case %d", idx))
				assert.Equal(t, n1.BookUUID, tc.expectedBookUUID, fmt.Sprintf("n1 BookUUID mismatch for test case %d", idx))
				assert.Equal(t, n1.USN, tc.expectedUSN, fmt.Sprintf("n1 USN mismatch for test case %d", idx))
				assert.Equal(t, n1.AddedOn, tc.expectedAddedOn, fmt.Sprintf("n1 AddedOn mismatch for test case %d", idx))
				assert.Equal(t, n1.EditedOn, tc.expectedEditedOn, fmt.Sprintf("n1 EditedOn mismatch for test case %d", idx))
				assert.Equal(t, n1.Body, tc.expectedBody, fmt.Sprintf("n1 Body mismatch for test case %d", idx))
				assert.Equal(t, n1.Deleted, tc.expectedDeleted, fmt.Sprintf("n1 Deleted mismatch for test case %d", idx))
				assert.Equal(t, n1.Dirty, tc.expectedDirty, fmt.Sprintf("n1 Dirty mismatch for test case %d", idx))
			}()
		}
	})
}

func TestStepSyncBook(t *testing.T) {
	t.Run("exists on server only", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		b1UUID := utils.GenerateUUID()
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", b1UUID, 555, "b1-label", true, false)

		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		b2UUID := utils.GenerateUUID()
		b := client.SyncFragBook{
			UUID:    b2UUID,
			USN:     1,
			AddedOn: 1541108743,
			Label:   "b2-label",
			Deleted: false,
		}

		if err := stepSyncBook(tx, b); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 2, "book count mismatch")

		var b1, b2 database.Book
		database.MustScan(t, "getting b1",
			db.QueryRow("SELECT uuid, usn, label, dirty, deleted FROM books WHERE uuid = ?", b1UUID),
			&b1.UUID, &b1.USN, &b1.Label, &b1.Dirty, &b1.Deleted)
		database.MustScan(t, "getting b2",
			db.QueryRow("SELECT uuid, usn, label, dirty, deleted FROM books WHERE uuid = ?", b2UUID),
			&b2.UUID, &b2.USN, &b2.Label, &b2.Dirty, &b2.Deleted)

		assert.Equal(t, b1.UUID, b1UUID, "b1 UUID mismatch")
		assert.Equal(t, b1.USN, 555, "b1 USN mismatch")
		assert.Equal(t, b1.Label, "b1-label", "b1 Label mismatch")
		assert.Equal(t, b1.Dirty, true, "b1 Dirty mismatch")
		assert.Equal(t, b1.Deleted, false, "b1 Deleted mismatch")

		assert.Equal(t, b2.UUID, b2UUID, "b2 UUID mismatch")
		assert.Equal(t, b2.USN, b.USN, "b2 USN mismatch")
		assert.Equal(t, b2.Label, b.Label, "b2 Label mismatch")
		assert.Equal(t, b2.Dirty, false, "b2 Dirty mismatch")
		assert.Equal(t, b2.Deleted, b.Deleted, "b2 Deleted mismatch")
	})

	t.Run("exists on server and client", func(t *testing.T) {
		testCases := []struct {
			clientDirty              bool
			clientUSN                int
			clientLabel              string
			clientDeleted            bool
			serverUSN                int
			serverLabel              string
			serverDeleted            bool
			expectedUSN              int
			expectedLabel            string
			expectedDeleted          bool
			anotherBookLabel         string
			expectedAnotherBookLabel string
			expectedAnotherBookDirty bool
		}{
			{
				clientDirty:              true,
				clientUSN:                1,
				clientLabel:              "b2-label",
				clientDeleted:            false,
				serverUSN:                3,
				serverLabel:              "b2-label-updated",
				serverDeleted:            false,
				expectedUSN:              3,
				expectedLabel:            "b2-label-updated",
				expectedDeleted:          false,
				anotherBookLabel:         "foo",
				expectedAnotherBookLabel: "foo",
				expectedAnotherBookDirty: false,
			},
			{
				clientDirty:              false,
				clientUSN:                1,
				clientLabel:              "b2-label",
				clientDeleted:            false,
				serverUSN:                3,
				serverLabel:              "b2-label-updated",
				serverDeleted:            false,
				expectedUSN:              3,
				expectedLabel:            "b2-label-updated",
				expectedDeleted:          false,
				anotherBookLabel:         "foo",
				expectedAnotherBookLabel: "foo",
				expectedAnotherBookDirty: false,
			},
			{
				clientDirty:              false,
				clientUSN:                1,
				clientLabel:              "b2-label",
				clientDeleted:            false,
				serverUSN:                3,
				serverLabel:              "foo",
				serverDeleted:            false,
				expectedUSN:              3,
				expectedLabel:            "foo",
				expectedDeleted:          false,
				anotherBookLabel:         "foo",
				expectedAnotherBookLabel: "foo_2",
				expectedAnotherBookDirty: true,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				db := database.InitTestDB(t, dbPath, nil)
				defer database.CloseTestDB(t, db)

				b1UUID := utils.GenerateUUID()
				database.MustExec(t, fmt.Sprintf("inserting book for test case %d", idx), db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", b1UUID, tc.clientUSN, tc.clientLabel, tc.clientDirty, tc.clientDeleted)
				b2UUID := utils.GenerateUUID()
				database.MustExec(t, fmt.Sprintf("inserting book for test case %d", idx), db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", b2UUID, 2, tc.anotherBookLabel, false, false)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				// update all fields but uuid and bump usn
				b := client.SyncFragBook{
					UUID:    b1UUID,
					USN:     tc.serverUSN,
					Label:   tc.serverLabel,
					Deleted: tc.serverDeleted,
				}

				if err := fullSyncBook(tx, b); err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				var noteCount, bookCount int
				database.MustScan(t, fmt.Sprintf("counting notes for test case %d", idx), db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
				database.MustScan(t, fmt.Sprintf("counting books for test case %d", idx), db.QueryRow("SELECT count(*) FROM books"), &bookCount)

				assert.Equalf(t, noteCount, 0, fmt.Sprintf("note count mismatch for test case %d", idx))
				assert.Equalf(t, bookCount, 2, fmt.Sprintf("book count mismatch for test case %d", idx))

				var b1Record, b2Record database.Book
				database.MustScan(t, "getting b1Record",
					db.QueryRow("SELECT uuid, usn, label, dirty, deleted FROM books WHERE uuid = ?", b1UUID),
					&b1Record.UUID, &b1Record.USN, &b1Record.Label, &b1Record.Dirty, &b1Record.Deleted)
				database.MustScan(t, "getting b2Record",
					db.QueryRow("SELECT uuid, usn, label, dirty, deleted FROM books WHERE uuid = ?", b2UUID),
					&b2Record.UUID, &b2Record.USN, &b2Record.Label, &b2Record.Dirty, &b2Record.Deleted)

				assert.Equal(t, b1Record.UUID, b1UUID, fmt.Sprintf("b1Record UUID mismatch for idx %d", idx))
				assert.Equal(t, b1Record.USN, tc.expectedUSN, fmt.Sprintf("b1Record USN mismatch for test case %d", idx))
				assert.Equal(t, b1Record.Label, tc.expectedLabel, fmt.Sprintf("b1Record Label mismatch for test case %d", idx))
				assert.Equal(t, b1Record.Dirty, tc.clientDirty, fmt.Sprintf("b1Record Dirty mismatch for test case %d", idx))
				assert.Equal(t, b1Record.Deleted, tc.expectedDeleted, fmt.Sprintf("b1Record Deleted mismatch for test case %d", idx))

				assert.Equal(t, b2Record.UUID, b2UUID, fmt.Sprintf("b2Record UUID mismatch for idx %d", idx))
				assert.Equal(t, b2Record.USN, 2, fmt.Sprintf("b2Record USN mismatch for test case %d", idx))
				assert.Equal(t, b2Record.Label, tc.expectedAnotherBookLabel, fmt.Sprintf("b2Record Label mismatch for test case %d", idx))
				assert.Equal(t, b2Record.Dirty, tc.expectedAnotherBookDirty, fmt.Sprintf("b2Record Dirty mismatch for test case %d", idx))
				assert.Equal(t, b2Record.Deleted, false, fmt.Sprintf("b2Record Deleted mismatch for test case %d", idx))
			}()
		}
	})
}

func TestMergeBook(t *testing.T) {
	t.Run("insert, no duplicates", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		// test
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		b1 := client.SyncFragBook{
			UUID:    "b1-uuid",
			USN:     12,
			AddedOn: 1541108743,
			Label:   "b1-label",
			Deleted: false,
		}

		if err := mergeBook(tx, b1, modeInsert); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// execute
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 1, "book count mismatch")

		var b1Record database.Book
		database.MustScan(t, "getting b1",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b1-uuid"),
			&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)

		assert.Equal(t, b1Record.UUID, b1.UUID, "b1 UUID mismatch")
		assert.Equal(t, b1Record.Label, b1.Label, "b1 Label mismatch")
		assert.Equal(t, b1Record.USN, b1.USN, "b1 USN mismatch")
	})

	t.Run("insert, 1 duplicate", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b1-uuid", 1, "foo", false, false)

		// test
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		b := client.SyncFragBook{
			UUID:    "b2-uuid",
			USN:     12,
			AddedOn: 1541108743,
			Label:   "foo",
			Deleted: false,
		}

		if err := mergeBook(tx, b, modeInsert); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// execute
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 2, "book count mismatch")

		var b1Record, b2Record database.Book
		database.MustScan(t, "getting b1",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b1-uuid"),
			&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)
		database.MustScan(t, "getting b2",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b2-uuid"),
			&b2Record.UUID, &b2Record.Label, &b2Record.USN, &b2Record.Dirty)

		assert.Equal(t, b1Record.Label, "foo_2", "b1 Label mismatch")
		assert.Equal(t, b1Record.USN, 1, "b1 USN mismatch")
		assert.Equal(t, b1Record.Dirty, true, "b1 should have been marked dirty")

		assert.Equal(t, b2Record.Label, "foo", "b2 Label mismatch")
		assert.Equal(t, b2Record.USN, 12, "b2 USN mismatch")
		assert.Equal(t, b2Record.Dirty, false, "b2 Dirty mismatch")
	})

	t.Run("insert, 3 duplicates", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b1-uuid", 1, "foo", false, false)
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b2-uuid", 2, "foo_2", true, false)
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b3-uuid", 3, "foo_3", false, false)

		// test
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		b := client.SyncFragBook{
			UUID:    "b4-uuid",
			USN:     12,
			AddedOn: 1541108743,
			Label:   "foo",
			Deleted: false,
		}

		if err := mergeBook(tx, b, modeInsert); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// execute
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 4, "book count mismatch")

		var b1Record, b2Record, b3Record, b4Record database.Book
		database.MustScan(t, "getting b1",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b1-uuid"),
			&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)
		database.MustScan(t, "getting b2",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b2-uuid"),
			&b2Record.UUID, &b2Record.Label, &b2Record.USN, &b2Record.Dirty)
		database.MustScan(t, "getting b3",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b3-uuid"),
			&b3Record.UUID, &b3Record.Label, &b3Record.USN, &b3Record.Dirty)
		database.MustScan(t, "getting b4",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b4-uuid"),
			&b4Record.UUID, &b4Record.Label, &b4Record.USN, &b4Record.Dirty)

		assert.Equal(t, b1Record.Label, "foo_4", "b1 Label mismatch")
		assert.Equal(t, b1Record.USN, 1, "b1 USN mismatch")
		assert.Equal(t, b1Record.Dirty, true, "b1 Dirty mismatch")

		assert.Equal(t, b2Record.Label, "foo_2", "b2 Label mismatch")
		assert.Equal(t, b2Record.USN, 2, "b2 USN mismatch")
		assert.Equal(t, b2Record.Dirty, true, "b2 Dirty mismatch")

		assert.Equal(t, b3Record.Label, "foo_3", "b3 Label mismatch")
		assert.Equal(t, b3Record.USN, 3, "b3 USN mismatch")
		assert.Equal(t, b3Record.Dirty, false, "b3 Dirty mismatch")

		assert.Equal(t, b4Record.Label, "foo", "b4 Label mismatch")
		assert.Equal(t, b4Record.USN, 12, "b4 USN mismatch")
		assert.Equal(t, b4Record.Dirty, false, "b4 Dirty mismatch")
	})

	t.Run("update, no duplicates", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		// test
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		b1UUID := utils.GenerateUUID()
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", b1UUID, 1, "b1-label", false, false)

		b1 := client.SyncFragBook{
			UUID:    b1UUID,
			USN:     12,
			AddedOn: 1541108743,
			Label:   "b1-label-edited",
			Deleted: false,
		}

		if err := mergeBook(tx, b1, modeUpdate); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// execute
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 1, "book count mismatch")

		var b1Record database.Book
		database.MustScan(t, "getting b1",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", b1UUID),
			&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)

		assert.Equal(t, b1Record.UUID, b1UUID, "b1 UUID mismatch")
		assert.Equal(t, b1Record.Label, "b1-label-edited", "b1 Label mismatch")
		assert.Equal(t, b1Record.USN, 12, "b1 USN mismatch")
	})

	t.Run("update, 1 duplicate", func(t *testing.T) {
		// set up
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b1-uuid", 1, "foo", false, false)
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b2-uuid", 2, "bar", false, false)

		// test
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		b := client.SyncFragBook{
			UUID:    "b1-uuid",
			USN:     12,
			AddedOn: 1541108743,
			Label:   "bar",
			Deleted: false,
		}

		if err := mergeBook(tx, b, modeUpdate); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// execute
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 2, "book count mismatch")

		var b1Record, b2Record database.Book
		database.MustScan(t, "getting b1",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b1-uuid"),
			&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)
		database.MustScan(t, "getting b2",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b2-uuid"),
			&b2Record.UUID, &b2Record.Label, &b2Record.USN, &b2Record.Dirty)

		assert.Equal(t, b1Record.Label, "bar", "b1 Label mismatch")
		assert.Equal(t, b1Record.USN, 12, "b1 USN mismatch")
		assert.Equal(t, b1Record.Dirty, false, "b1 Dirty mismatch")

		assert.Equal(t, b2Record.Label, "bar_2", "b2 Label mismatch")
		assert.Equal(t, b2Record.USN, 2, "b2 USN mismatch")
		assert.Equal(t, b2Record.Dirty, true, "b2 Dirty mismatch")
	})

	t.Run("update, 3 duplicate", func(t *testing.T) {
		// set uj
		db := database.InitTestDB(t, dbPath, nil)
		defer database.CloseTestDB(t, db)

		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b1-uuid", 1, "foo", false, false)
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b2-uuid", 2, "bar", false, false)
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b3-uuid", 3, "bar_2", true, false)
		database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, usn, label, dirty, deleted) VALUES (?, ?, ?, ?, ?)", "b4-uuid", 4, "bar_3", false, false)

		// test
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}

		b := client.SyncFragBook{
			UUID:    "b1-uuid",
			USN:     12,
			AddedOn: 1541108743,
			Label:   "bar",
			Deleted: false,
		}

		if err := mergeBook(tx, b, modeUpdate); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// execute
		var noteCount, bookCount int
		database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
		database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

		assert.Equalf(t, noteCount, 0, "note count mismatch")
		assert.Equalf(t, bookCount, 4, "book count mismatch")

		var b1Record, b2Record, b3Record, b4Record database.Book
		database.MustScan(t, "getting b1",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b1-uuid"),
			&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)
		database.MustScan(t, "getting b2",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b2-uuid"),
			&b2Record.UUID, &b2Record.Label, &b2Record.USN, &b2Record.Dirty)
		database.MustScan(t, "getting b3",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b3-uuid"),
			&b3Record.UUID, &b3Record.Label, &b3Record.USN, &b3Record.Dirty)
		database.MustScan(t, "getting b4",
			db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", "b4-uuid"),
			&b4Record.UUID, &b4Record.Label, &b4Record.USN, &b4Record.Dirty)

		assert.Equal(t, b1Record.Label, "bar", "b1 Label mismatch")
		assert.Equal(t, b1Record.USN, 12, "b1 USN mismatch")
		assert.Equal(t, b1Record.Dirty, false, "b1 Dirty mismatch")

		assert.Equal(t, b2Record.Label, "bar_4", "b2 Label mismatch")
		assert.Equal(t, b2Record.USN, 2, "b2 USN mismatch")
		assert.Equal(t, b2Record.Dirty, true, "b2 Dirty mismatch")

		assert.Equal(t, b3Record.Label, "bar_2", "b3 Label mismatch")
		assert.Equal(t, b3Record.USN, 3, "b3 USN mismatch")
		assert.Equal(t, b3Record.Dirty, true, "b3 Dirty mismatch")

		assert.Equal(t, b4Record.Label, "bar_3", "b4 Label mismatch")
		assert.Equal(t, b4Record.USN, 4, "b4 USN mismatch")
		assert.Equal(t, b4Record.Dirty, false, "b4 Dirty mismatch")
	})
}

func TestSaveServerState(t *testing.T) {
	// set up
	ctx := context.InitTestCtx(t, "../../tmp", nil)
	defer context.TeardownTestCtx(t, ctx)
	testutils.Login(t, &ctx)

	db := ctx.DB

	database.MustExec(t, "inserting last synced at", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastSyncAt, int64(1231108742))
	database.MustExec(t, "inserting last max usn", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, 8)

	// execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	serverTime := int64(1541108743)
	serverMaxUSN := 100

	err = saveSyncState(tx, serverTime, serverMaxUSN)
	if err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "executing").Error())
	}

	tx.Commit()

	// test
	var lastSyncedAt int64
	var lastMaxUSN int

	database.MustScan(t, "getting system value",
		db.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastSyncAt), &lastSyncedAt)
	database.MustScan(t, "getting system value",
		db.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastMaxUSN), &lastMaxUSN)

	assert.Equal(t, lastSyncedAt, serverTime, "last synced at mismatch")
	assert.Equal(t, lastMaxUSN, serverMaxUSN, "last max usn mismatch")
}

// TestSendBooks tests that books are put to correct 'buckets' by running a test server and recording the
// uuid from the incoming data. It also tests that the uuid of the created books and book_uuids of their notes
// are updated accordingly based on the server response.
func TestSendBooks(t *testing.T) {
	// set up
	ctx := context.InitTestCtx(t, "../../tmp", nil)
	defer context.TeardownTestCtx(t, ctx)
	testutils.Login(t, &ctx)

	db := ctx.DB

	database.MustExec(t, "inserting last max usn", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, 0)

	// should be ignored
	database.MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b1-uuid", "b1-label", 1, false, false)
	database.MustExec(t, "inserting b2", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b2-uuid", "b2-label", 2, false, false)
	// should be created
	database.MustExec(t, "inserting b3", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b3-uuid", "b3-label", 0, false, true)
	database.MustExec(t, "inserting b4", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b4-uuid", "b4-label", 0, false, true)
	// should be only expunged locally without syncing to server
	database.MustExec(t, "inserting b5", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b5-uuid", "b5-label", 0, true, true)
	// should be deleted
	database.MustExec(t, "inserting b6", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b6-uuid", "b6-label", 10, true, true)
	// should be updated
	database.MustExec(t, "inserting b7", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b7-uuid", "b7-label", 11, false, true)
	database.MustExec(t, "inserting b8", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b8-uuid", "b8-label", 18, false, true)

	// some random notes
	database.MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", "b1-uuid", 10, "n1 body", 1541108743, false, false)
	database.MustExec(t, "inserting n2", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n2-uuid", "b5-uuid", 10, "n2 body", 1541108743, false, false)
	database.MustExec(t, "inserting n3", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n3-uuid", "b6-uuid", 10, "n3 body", 1541108743, false, false)
	database.MustExec(t, "inserting n4", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n4-uuid", "b7-uuid", 10, "n4 body", 1541108743, false, false)
	// notes that belong to the created book. Their book_uuid should be updated.
	database.MustExec(t, "inserting n5", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n5-uuid", "b3-uuid", 10, "n5 body", 1541108743, false, false)
	database.MustExec(t, "inserting n6", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n6-uuid", "b3-uuid", 10, "n6 body", 1541108743, false, false)
	database.MustExec(t, "inserting n7", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n7-uuid", "b4-uuid", 10, "n7 body", 1541108743, false, false)

	var createdLabels []string
	var updatesUUIDs []string
	var deletedUUIDs []string

	// fire up a test server. It decrypts the payload for test purposes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/v3/books" && r.Method == "POST" {
			var payload client.CreateBookPayload

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				t.Fatalf(errors.Wrap(err, "decoding payload in the test server").Error())
				return
			}

			createdLabels = append(createdLabels, payload.Name)

			resp := client.CreateBookResp{
				Book: client.RespBook{
					UUID: fmt.Sprintf("server-%s-uuid", payload.Name),
				},
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		p := strings.Split(r.URL.Path, "/")
		if len(p) == 4 && p[0] == "" && p[1] == "v3" && p[2] == "books" {
			if r.Method == "PATCH" {
				uuid := p[3]
				updatesUUIDs = append(updatesUUIDs, uuid)

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{}"))
				return
			} else if r.Method == "DELETE" {
				uuid := p[3]
				deletedUUIDs = append(deletedUUIDs, uuid)

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{}"))
				return
			}
		}

		t.Fatalf("unrecognized endpoint reached Method: %s Path: %s", r.Method, r.URL.Path)
	}))
	defer ts.Close()

	ctx.APIEndpoint = ts.URL

	// execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if _, err := sendBooks(ctx, tx); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "executing").Error())
	}

	tx.Commit()

	// test

	// First, decrypt data so that they can be asserted
	sort.SliceStable(createdLabels, func(i, j int) bool {
		return strings.Compare(createdLabels[i], createdLabels[j]) < 0
	})

	assert.DeepEqual(t, createdLabels, []string{"b3-label", "b4-label"}, "createdLabels mismatch")
	assert.DeepEqual(t, updatesUUIDs, []string{"b7-uuid", "b8-uuid"}, "updatesUUIDs mismatch")
	assert.DeepEqual(t, deletedUUIDs, []string{"b6-uuid"}, "deletedUUIDs mismatch")

	var b1, b2, b3, b4, b7, b8 database.Book
	database.MustScan(t, "getting b1", db.QueryRow("SELECT uuid, dirty FROM books WHERE label = ?", "b1-label"), &b1.UUID, &b1.Dirty)
	database.MustScan(t, "getting b2", db.QueryRow("SELECT uuid, dirty FROM books WHERE label = ?", "b2-label"), &b2.UUID, &b2.Dirty)
	database.MustScan(t, "getting b3", db.QueryRow("SELECT uuid, dirty FROM books WHERE label = ?", "b3-label"), &b3.UUID, &b3.Dirty)
	database.MustScan(t, "getting b4", db.QueryRow("SELECT uuid, dirty FROM books WHERE label = ?", "b4-label"), &b4.UUID, &b4.Dirty)
	database.MustScan(t, "getting b7", db.QueryRow("SELECT uuid, dirty FROM books WHERE label = ?", "b7-label"), &b7.UUID, &b7.Dirty)
	database.MustScan(t, "getting b8", db.QueryRow("SELECT uuid, dirty FROM books WHERE label = ?", "b8-label"), &b8.UUID, &b8.Dirty)

	var bookCount int
	database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	assert.Equalf(t, bookCount, 6, "book count mismatch")

	assert.Equal(t, b1.Dirty, false, "b1 Dirty mismatch")
	assert.Equal(t, b2.Dirty, false, "b2 Dirty mismatch")
	assert.Equal(t, b3.Dirty, false, "b3 Dirty mismatch")
	assert.Equal(t, b4.Dirty, false, "b4 Dirty mismatch")
	assert.Equal(t, b7.Dirty, false, "b7 Dirty mismatch")
	assert.Equal(t, b8.Dirty, false, "b8 Dirty mismatch")
	assert.Equal(t, b1.UUID, "b1-uuid", "b1 UUID mismatch")
	assert.Equal(t, b2.UUID, "b2-uuid", "b2 UUID mismatch")
	// uuids of created books should have been updated
	assert.Equal(t, b3.UUID, "server-b3-label-uuid", "b3 UUID mismatch")
	assert.Equal(t, b4.UUID, "server-b4-label-uuid", "b4 UUID mismatch")
	assert.Equal(t, b7.UUID, "b7-uuid", "b7 UUID mismatch")
	assert.Equal(t, b8.UUID, "b8-uuid", "b8 UUID mismatch")

	var n1, n2, n3, n4, n5, n6, n7 database.Note
	database.MustScan(t, "getting n1", db.QueryRow("SELECT book_uuid FROM notes WHERE body = ?", "n1 body"), &n1.BookUUID)
	database.MustScan(t, "getting n2", db.QueryRow("SELECT book_uuid FROM notes WHERE body = ?", "n2 body"), &n2.BookUUID)
	database.MustScan(t, "getting n3", db.QueryRow("SELECT book_uuid FROM notes WHERE body = ?", "n3 body"), &n3.BookUUID)
	database.MustScan(t, "getting n4", db.QueryRow("SELECT book_uuid FROM notes WHERE body = ?", "n4 body"), &n4.BookUUID)
	database.MustScan(t, "getting n5", db.QueryRow("SELECT book_uuid FROM notes WHERE body = ?", "n5 body"), &n5.BookUUID)
	database.MustScan(t, "getting n6", db.QueryRow("SELECT book_uuid FROM notes WHERE body = ?", "n6 body"), &n6.BookUUID)
	database.MustScan(t, "getting n7", db.QueryRow("SELECT book_uuid FROM notes WHERE body = ?", "n7 body"), &n7.BookUUID)
	assert.Equal(t, n1.BookUUID, "b1-uuid", "n1 bookUUID mismatch")
	assert.Equal(t, n2.BookUUID, "b5-uuid", "n2 bookUUID mismatch")
	assert.Equal(t, n3.BookUUID, "b6-uuid", "n3 bookUUID mismatch")
	assert.Equal(t, n4.BookUUID, "b7-uuid", "n4 bookUUID mismatch")
	assert.Equal(t, n5.BookUUID, "server-b3-label-uuid", "n5 bookUUID mismatch")
	assert.Equal(t, n6.BookUUID, "server-b3-label-uuid", "n6 bookUUID mismatch")
	assert.Equal(t, n7.BookUUID, "server-b4-label-uuid", "n7 bookUUID mismatch")
}

func TestSendBooks_isBehind(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/v3/books" && r.Method == "POST" {
			var payload client.CreateBookPayload

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				t.Fatalf(errors.Wrap(err, "decoding payload in the test server").Error())
				return
			}

			resp := client.CreateBookResp{
				Book: client.RespBook{
					USN: 11,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		p := strings.Split(r.URL.Path, "/")
		if len(p) == 4 && p[0] == "" && p[1] == "v3" && p[2] == "books" {
			if r.Method == "PATCH" {
				resp := client.UpdateBookResp{
					Book: client.RespBook{
						USN: 11,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			} else if r.Method == "DELETE" {
				resp := client.DeleteBookResp{
					Book: client.RespBook{
						USN: 11,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}

		t.Fatalf("unrecognized endpoint reached Method: %s Path: %s", r.Method, r.URL.Path)
	}))
	defer ts.Close()

	t.Run("create book", func(t *testing.T) {
		testCases := []struct {
			systemLastMaxUSN int
			expectedIsBehind bool
		}{
			{
				systemLastMaxUSN: 10,
				expectedIsBehind: false,
			},
			{
				systemLastMaxUSN: 9,
				expectedIsBehind: true,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				ctx := context.InitTestCtx(t, "../../tmp", nil)
				ctx.APIEndpoint = ts.URL
				defer context.TeardownTestCtx(t, ctx)
				testutils.Login(t, &ctx)

				db := ctx.DB

				database.MustExec(t, fmt.Sprintf("inserting last max usn for test case %d", idx), db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, tc.systemLastMaxUSN)
				database.MustExec(t, fmt.Sprintf("inserting b1 for test case %d", idx), db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b1-uuid", "b1-label", 0, false, true)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				isBehind, err := sendBooks(ctx, tx)
				if err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				assert.Equal(t, isBehind, tc.expectedIsBehind, fmt.Sprintf("isBehind mismatch for test case %d", idx))
			}()
		}
	})

	t.Run("delete book", func(t *testing.T) {
		testCases := []struct {
			systemLastMaxUSN int
			expectedIsBehind bool
		}{
			{
				systemLastMaxUSN: 10,
				expectedIsBehind: false,
			},
			{
				systemLastMaxUSN: 9,
				expectedIsBehind: true,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				ctx := context.InitTestCtx(t, "../../tmp", nil)
				ctx.APIEndpoint = ts.URL
				defer context.TeardownTestCtx(t, ctx)
				testutils.Login(t, &ctx)

				db := ctx.DB

				database.MustExec(t, fmt.Sprintf("inserting last max usn for test case %d", idx), db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, tc.systemLastMaxUSN)
				database.MustExec(t, fmt.Sprintf("inserting b1 for test case %d", idx), db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b1-uuid", "b1-label", 1, true, true)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				isBehind, err := sendBooks(ctx, tx)
				if err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				assert.Equal(t, isBehind, tc.expectedIsBehind, fmt.Sprintf("isBehind mismatch for test case %d", idx))
			}()
		}
	})

	t.Run("update book", func(t *testing.T) {
		testCases := []struct {
			systemLastMaxUSN int
			expectedIsBehind bool
		}{
			{
				systemLastMaxUSN: 10,
				expectedIsBehind: false,
			},
			{
				systemLastMaxUSN: 9,
				expectedIsBehind: true,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				ctx := context.InitTestCtx(t, "../../tmp", nil)
				ctx.APIEndpoint = ts.URL
				defer context.TeardownTestCtx(t, ctx)
				testutils.Login(t, &ctx)

				db := ctx.DB

				database.MustExec(t, fmt.Sprintf("inserting last max usn for test case %d", idx), db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, tc.systemLastMaxUSN)
				database.MustExec(t, fmt.Sprintf("inserting b1 for test case %d", idx), db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b1-uuid", "b1-label", 11, false, true)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				isBehind, err := sendBooks(ctx, tx)
				if err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				assert.Equal(t, isBehind, tc.expectedIsBehind, fmt.Sprintf("isBehind mismatch for test case %d", idx))
			}()
		}
	})
}

// TestSendNotes tests that notes are put to correct 'buckets' by running a test server and recording the
// uuid from the incoming data.
func TestSendNotes(t *testing.T) {
	// set up
	ctx := context.InitTestCtx(t, "../../tmp", nil)
	defer context.TeardownTestCtx(t, ctx)
	testutils.Login(t, &ctx)

	db := ctx.DB

	database.MustExec(t, "inserting last max usn", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, 0)

	b1UUID := "b1-uuid"
	database.MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b1UUID, "b1-label", 1, false, false)

	// should be ignored
	database.MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", b1UUID, 10, "n1-body", 1541108743, false, false)
	// should be created
	database.MustExec(t, "inserting n2", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n2-uuid", b1UUID, 0, "n2-body", 1541108743, false, true)
	// should be updated
	database.MustExec(t, "inserting n3", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n3-uuid", b1UUID, 11, "n3-body", 1541108743, false, true)
	// should be only expunged locally without syncing to server
	database.MustExec(t, "inserting n4", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n4-uuid", b1UUID, 0, "n4-body", 1541108743, true, true)
	// should be deleted
	database.MustExec(t, "inserting n5", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n5-uuid", b1UUID, 17, "n5-body", 1541108743, true, true)
	// should be created
	database.MustExec(t, "inserting n6", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n6-uuid", b1UUID, 0, "n6-body", 1541108743, false, true)
	// should be ignored
	database.MustExec(t, "inserting n7", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n7-uuid", b1UUID, 12, "n7-body", 1541108743, false, false)
	// should be updated
	database.MustExec(t, "inserting n8", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n8-uuid", b1UUID, 17, "n8-body", 1541108743, false, true)
	// should be deleted
	database.MustExec(t, "inserting n9", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n9-uuid", b1UUID, 17, "n9-body", 1541108743, true, true)
	// should be created
	database.MustExec(t, "inserting n10", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n10-uuid", b1UUID, 0, "n10-body", 1541108743, false, true)

	var createdBodys []string
	var updatedUUIDs []string
	var deletedUUIDs []string

	// fire up a test server. It decrypts the payload for test purposes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/v3/notes" && r.Method == "POST" {
			var payload client.CreateNotePayload

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				t.Fatalf(errors.Wrap(err, "decoding payload in the test server").Error())
				return
			}

			createdBodys = append(createdBodys, payload.Body)

			resp := client.CreateNoteResp{
				Result: client.RespNote{
					UUID: fmt.Sprintf("server-%s-uuid", payload.Body),
				},
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		p := strings.Split(r.URL.Path, "/")
		if len(p) == 4 && p[0] == "" && p[1] == "v3" && p[2] == "notes" {
			if r.Method == "PATCH" {
				uuid := p[3]
				updatedUUIDs = append(updatedUUIDs, uuid)

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{}"))
				return
			} else if r.Method == "DELETE" {
				uuid := p[3]
				deletedUUIDs = append(deletedUUIDs, uuid)

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{}"))
				return
			}
		}

		t.Fatalf("unrecognized endpoint reached Method: %s Path: %s", r.Method, r.URL.Path)
	}))
	defer ts.Close()

	ctx.APIEndpoint = ts.URL

	// execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if _, err := sendNotes(ctx, tx); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "executing").Error())
	}

	tx.Commit()

	// test
	sort.SliceStable(createdBodys, func(i, j int) bool {
		return strings.Compare(createdBodys[i], createdBodys[j]) < 0
	})

	assert.DeepEqual(t, createdBodys, []string{"n10-body", "n2-body", "n6-body"}, "createdBodys mismatch")
	assert.DeepEqual(t, updatedUUIDs, []string{"n3-uuid", "n8-uuid"}, "updatedUUIDs mismatch")
	assert.DeepEqual(t, deletedUUIDs, []string{"n5-uuid", "n9-uuid"}, "deletedUUIDs mismatch")

	var noteCount int
	database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	assert.Equalf(t, noteCount, 7, "note count mismatch")

	var n1, n2, n3, n6, n7, n8, n10 database.Note
	database.MustScan(t, "getting n1", db.QueryRow("SELECT uuid, added_on, dirty FROM notes WHERE body = ?", "n1-body"), &n1.UUID, &n1.AddedOn, &n1.Dirty)
	database.MustScan(t, "getting n2", db.QueryRow("SELECT uuid, added_on, dirty FROM notes WHERE body = ?", "n2-body"), &n2.UUID, &n2.AddedOn, &n2.Dirty)
	database.MustScan(t, "getting n3", db.QueryRow("SELECT uuid, added_on, dirty FROM notes WHERE body = ?", "n3-body"), &n3.UUID, &n3.AddedOn, &n3.Dirty)
	database.MustScan(t, "getting n6", db.QueryRow("SELECT uuid, added_on, dirty FROM notes WHERE body = ?", "n6-body"), &n6.UUID, &n6.AddedOn, &n6.Dirty)
	database.MustScan(t, "getting n7", db.QueryRow("SELECT uuid, added_on, dirty FROM notes WHERE body = ?", "n7-body"), &n7.UUID, &n7.AddedOn, &n7.Dirty)
	database.MustScan(t, "getting n8", db.QueryRow("SELECT uuid, added_on, dirty FROM notes WHERE body = ?", "n8-body"), &n8.UUID, &n8.AddedOn, &n8.Dirty)
	database.MustScan(t, "getting n10", db.QueryRow("SELECT uuid, added_on, dirty FROM notes WHERE body = ?", "n10-body"), &n10.UUID, &n10.AddedOn, &n10.Dirty)

	assert.Equalf(t, noteCount, 7, "note count mismatch")

	assert.Equal(t, n1.Dirty, false, "n1 Dirty mismatch")
	assert.Equal(t, n2.Dirty, false, "n2 Dirty mismatch")
	assert.Equal(t, n3.Dirty, false, "n3 Dirty mismatch")
	assert.Equal(t, n6.Dirty, false, "n6 Dirty mismatch")
	assert.Equal(t, n7.Dirty, false, "n7 Dirty mismatch")
	assert.Equal(t, n8.Dirty, false, "n8 Dirty mismatch")
	assert.Equal(t, n10.Dirty, false, "n10 Dirty mismatch")

	assert.Equal(t, n1.AddedOn, int64(1541108743), "n1 AddedOn mismatch")
	assert.Equal(t, n2.AddedOn, int64(1541108743), "n2 AddedOn mismatch")
	assert.Equal(t, n3.AddedOn, int64(1541108743), "n3 AddedOn mismatch")
	assert.Equal(t, n6.AddedOn, int64(1541108743), "n6 AddedOn mismatch")
	assert.Equal(t, n7.AddedOn, int64(1541108743), "n7 AddedOn mismatch")
	assert.Equal(t, n8.AddedOn, int64(1541108743), "n8 AddedOn mismatch")
	assert.Equal(t, n10.AddedOn, int64(1541108743), "n10 AddedOn mismatch")

	// UUIDs of created notes should have been updated with those from the server response
	assert.Equal(t, n1.UUID, "n1-uuid", "n1 UUID mismatch")
	assert.Equal(t, n2.UUID, "server-n2-body-uuid", "n2 UUID mismatch")
	assert.Equal(t, n3.UUID, "n3-uuid", "n3 UUID mismatch")
	assert.Equal(t, n6.UUID, "server-n6-body-uuid", "n6 UUID mismatch")
	assert.Equal(t, n7.UUID, "n7-uuid", "n7 UUID mismatch")
	assert.Equal(t, n8.UUID, "n8-uuid", "n8 UUID mismatch")
	assert.Equal(t, n10.UUID, "server-n10-body-uuid", "n10 UUID mismatch")
}

func TestSendNotes_addedOn(t *testing.T) {
	// set up
	ctx := context.InitTestCtx(t, "../../tmp", nil)
	defer context.TeardownTestCtx(t, ctx)
	testutils.Login(t, &ctx)

	db := ctx.DB

	database.MustExec(t, "inserting last max usn", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, 0)

	// should be created
	b1UUID := "b1-uuid"
	database.MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", b1UUID, 0, "n1-body", 1541108743, false, true)

	// fire up a test server. It decrypts the payload for test purposes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/v3/notes" && r.Method == "POST" {
			resp := client.CreateNoteResp{
				Result: client.RespNote{
					UUID: utils.GenerateUUID(),
				},
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		t.Fatalf("unrecognized endpoint reached Method: %s Path: %s", r.Method, r.URL.Path)
	}))
	defer ts.Close()

	ctx.APIEndpoint = ts.URL

	// execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if _, err := sendNotes(ctx, tx); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "executing").Error())
	}

	tx.Commit()

	// test
	var n1 database.Note
	database.MustScan(t, "getting n1", db.QueryRow("SELECT uuid, added_on, dirty FROM notes WHERE body = ?", "n1-body"), &n1.UUID, &n1.AddedOn, &n1.Dirty)
	assert.Equal(t, n1.AddedOn, int64(1541108743), "n1 AddedOn mismatch")
}

func TestSendNotes_isBehind(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/v3/notes" && r.Method == "POST" {
			var payload client.CreateBookPayload

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				t.Fatalf(errors.Wrap(err, "decoding payload in the test server").Error())
				return
			}

			resp := client.CreateNoteResp{
				Result: client.RespNote{
					USN: 11,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		p := strings.Split(r.URL.Path, "/")
		if len(p) == 4 && p[0] == "" && p[1] == "v3" && p[2] == "notes" {
			if r.Method == "PATCH" {
				resp := client.UpdateNoteResp{
					Result: client.RespNote{
						USN: 11,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			} else if r.Method == "DELETE" {
				resp := client.DeleteNoteResp{
					Result: client.RespNote{
						USN: 11,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}

		t.Fatalf("unrecognized endpoint reached Method: %s Path: %s", r.Method, r.URL.Path)
	}))
	defer ts.Close()

	t.Run("create note", func(t *testing.T) {
		testCases := []struct {
			systemLastMaxUSN int
			expectedIsBehind bool
		}{
			{
				systemLastMaxUSN: 10,
				expectedIsBehind: false,
			},
			{
				systemLastMaxUSN: 9,
				expectedIsBehind: true,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				ctx := context.InitTestCtx(t, "../../tmp", nil)
				defer context.TeardownTestCtx(t, ctx)
				testutils.Login(t, &ctx)
				ctx.APIEndpoint = ts.URL

				db := ctx.DB

				database.MustExec(t, fmt.Sprintf("inserting last max usn for test case %d", idx), db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, tc.systemLastMaxUSN)
				database.MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b1-uuid", "b1-label", 1, false, false)
				database.MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", "b1-uuid", 1, "n1 body", 1541108743, false, true)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				isBehind, err := sendNotes(ctx, tx)
				if err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				assert.Equal(t, isBehind, tc.expectedIsBehind, fmt.Sprintf("isBehind mismatch for test case %d", idx))
			}()
		}
	})

	t.Run("delete note", func(t *testing.T) {
		testCases := []struct {
			systemLastMaxUSN int
			expectedIsBehind bool
		}{
			{
				systemLastMaxUSN: 10,
				expectedIsBehind: false,
			},
			{
				systemLastMaxUSN: 9,
				expectedIsBehind: true,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				ctx := context.InitTestCtx(t, "../../tmp", nil)
				defer context.TeardownTestCtx(t, ctx)
				testutils.Login(t, &ctx)
				ctx.APIEndpoint = ts.URL

				db := ctx.DB

				database.MustExec(t, fmt.Sprintf("inserting last max usn for test case %d", idx), db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, tc.systemLastMaxUSN)
				database.MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b1-uuid", "b1-label", 1, false, false)
				database.MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", "b1-uuid", 2, "n1 body", 1541108743, true, true)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				isBehind, err := sendNotes(ctx, tx)
				if err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				assert.Equal(t, isBehind, tc.expectedIsBehind, fmt.Sprintf("isBehind mismatch for test case %d", idx))
			}()
		}
	})

	t.Run("update note", func(t *testing.T) {
		testCases := []struct {
			systemLastMaxUSN int
			expectedIsBehind bool
		}{
			{
				systemLastMaxUSN: 10,
				expectedIsBehind: false,
			},
			{
				systemLastMaxUSN: 9,
				expectedIsBehind: true,
			},
		}

		for idx, tc := range testCases {
			func() {
				// set up
				ctx := context.InitTestCtx(t, "../../tmp", nil)
				defer context.TeardownTestCtx(t, ctx)
				testutils.Login(t, &ctx)
				ctx.APIEndpoint = ts.URL

				db := ctx.DB

				database.MustExec(t, fmt.Sprintf("inserting last max usn for test case %d", idx), db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemLastMaxUSN, tc.systemLastMaxUSN)
				database.MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b1-uuid", "b1-label", 1, false, false)
				database.MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", "b1-uuid", 8, "n1 body", 1541108743, false, true)

				// execute
				tx, err := db.Begin()
				if err != nil {
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
				}

				isBehind, err := sendNotes(ctx, tx)
				if err != nil {
					tx.Rollback()
					t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
				}

				tx.Commit()

				// test
				assert.Equal(t, isBehind, tc.expectedIsBehind, fmt.Sprintf("isBehind mismatch for test case %d", idx))
			}()
		}
	})
}

func TestMergeNote(t *testing.T) {
	b1UUID := "b1-uuid"
	b2UUID := "b2-uuid"
	conflictBookUUID := utils.GenerateUUID()

	testCases := []struct {
		addedOn          int64
		clientUSN        int
		clientEditedOn   int64
		clientBody       string
		clientDeleted    bool
		clientBookUUID   string
		clientDirty      bool
		serverUSN        int
		serverEditedOn   int64
		serverBody       string
		serverDeleted    bool
		serverBookUUID   string
		expectedUSN      int
		expectedAddedOn  int64
		expectedEditedOn int64
		expectedBody     string
		expectedDeleted  bool
		expectedBookUUID string
		expectedDirty    bool
	}{
		// local copy is not dirty
		{
			clientDirty:      false,
			clientUSN:        1,
			clientEditedOn:   0,
			clientBody:       "n1 body",
			clientDeleted:    false,
			clientBookUUID:   b1UUID,
			addedOn:          1541232118,
			serverUSN:        21,
			serverEditedOn:   1541219321,
			serverBody:       "n1 body edited",
			serverDeleted:    false,
			serverBookUUID:   b1UUID,
			expectedUSN:      21,
			expectedAddedOn:  1541232118,
			expectedEditedOn: 1541219321,
			expectedBody:     "n1 body edited",
			expectedDeleted:  false,
			expectedBookUUID: b1UUID,
			expectedDirty:    false,
		},
		// local copy is dirty and needs conflict resolution
		{
			clientDirty:      true,
			clientUSN:        1,
			clientEditedOn:   1541219320,
			clientBody:       "n1 body",
			clientDeleted:    false,
			clientBookUUID:   b1UUID,
			addedOn:          1541232118,
			serverUSN:        21,
			serverEditedOn:   1541219321,
			serverBody:       "n1 body edited",
			serverDeleted:    false,
			serverBookUUID:   b1UUID,
			expectedUSN:      21,
			expectedAddedOn:  1541232118,
			expectedEditedOn: 1541219321,
			expectedBody: `<<<<<<< Local
n1 body
=======
n1 body edited
>>>>>>> Server
`,
			expectedDeleted:  false,
			expectedBookUUID: b1UUID,
			expectedDirty:    true,
		},
		{
			clientDirty:      true,
			clientUSN:        1,
			clientEditedOn:   1541219319,
			clientBody:       "n1 body",
			clientDeleted:    false,
			clientBookUUID:   b1UUID,
			addedOn:          1541232118,
			serverUSN:        21,
			serverEditedOn:   1541219321,
			serverBody:       "n1 body edited",
			serverDeleted:    false,
			serverBookUUID:   b2UUID,
			expectedUSN:      21,
			expectedAddedOn:  1541232118,
			expectedEditedOn: 1541219321,
			expectedBody: `<<<<<<< Local
Moved to the book b1-label
=======
Moved to the book b2-label
>>>>>>> Server

<<<<<<< Local
n1 body
=======
n1 body edited
>>>>>>> Server
`,
			expectedDeleted:  false,
			expectedBookUUID: conflictBookUUID,
			expectedDirty:    true,
		},
		// deleted locally and edited on server
		{
			clientDirty:      true,
			clientUSN:        1,
			clientEditedOn:   1541219321,
			clientBody:       "",
			clientDeleted:    true,
			clientBookUUID:   b1UUID,
			addedOn:          1541232118,
			serverUSN:        21,
			serverEditedOn:   1541219321,
			serverBody:       "n1 body edited",
			serverDeleted:    false,
			serverBookUUID:   b2UUID,
			expectedUSN:      21,
			expectedAddedOn:  1541232118,
			expectedEditedOn: 1541219321,
			expectedBody:     "n1 body edited",
			expectedDeleted:  false,
			expectedBookUUID: b2UUID,
			expectedDirty:    false,
		},
	}

	for idx, tc := range testCases {
		func() {
			// set up
			db := database.InitTestDB(t, "../../tmp/.dnote", nil)
			defer database.CloseTestDB(t, db)

			database.MustExec(t, fmt.Sprintf("inserting b1 for test case %d", idx), db, "INSERT INTO books (uuid, label, usn, dirty) VALUES (?, ?, ?, ?)", b1UUID, "b1-label", 5, false)
			database.MustExec(t, fmt.Sprintf("inserting b2 for test case %d", idx), db, "INSERT INTO books (uuid, label, usn, dirty) VALUES (?, ?, ?, ?)", b2UUID, "b2-label", 6, false)
			database.MustExec(t, fmt.Sprintf("inserting conflitcs book for test case %d", idx), db, "INSERT INTO books (uuid, label) VALUES (?, ?)", conflictBookUUID, "conflicts")
			n1UUID := utils.GenerateUUID()
			database.MustExec(t, fmt.Sprintf("inserting n1 for test case %d", idx), db, "INSERT INTO notes (uuid, book_uuid, usn, added_on, edited_on, body, deleted, dirty) VALUES (?, ?, ?,  ?, ?, ?, ?, ?)", n1UUID, b1UUID, tc.clientUSN, tc.addedOn, tc.clientEditedOn, tc.clientBody, tc.clientDeleted, tc.clientDirty)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
			}

			// update all fields but uuid and bump usn
			fragNote := client.SyncFragNote{
				UUID:     n1UUID,
				BookUUID: tc.serverBookUUID,
				USN:      tc.serverUSN,
				AddedOn:  tc.addedOn,
				EditedOn: tc.serverEditedOn,
				Body:     tc.serverBody,
				Deleted:  tc.serverDeleted,
			}
			var localNote database.Note
			database.MustScan(t, fmt.Sprintf("getting localNote for test case %d", idx),
				db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body, deleted, dirty FROM notes WHERE uuid = ?", n1UUID),
				&localNote.UUID, &localNote.BookUUID, &localNote.USN, &localNote.AddedOn, &localNote.EditedOn, &localNote.Body, &localNote.Deleted, &localNote.Dirty)

			if err := mergeNote(tx, fragNote, localNote); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
			}

			tx.Commit()

			// test
			var noteCount, bookCount int
			database.MustScan(t, fmt.Sprintf("counting notes for test case %d", idx), db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
			database.MustScan(t, fmt.Sprintf("counting books for test case %d", idx), db.QueryRow("SELECT count(*) FROM books"), &bookCount)

			assert.Equalf(t, noteCount, 1, fmt.Sprintf("note count mismatch for test case %d", idx))
			assert.Equalf(t, bookCount, 3, fmt.Sprintf("book count mismatch for test case %d", idx))

			var n1Record database.Note
			database.MustScan(t, fmt.Sprintf("getting n1Record for test case %d", idx),
				db.QueryRow("SELECT uuid, book_uuid, usn, added_on, edited_on, body, deleted, dirty FROM notes WHERE uuid = ?", n1UUID),
				&n1Record.UUID, &n1Record.BookUUID, &n1Record.USN, &n1Record.AddedOn, &n1Record.EditedOn, &n1Record.Body, &n1Record.Deleted, &n1Record.Dirty)
			var b1Record database.Book
			database.MustScan(t, "getting b1Record for test case",
				db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", b1UUID),
				&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Dirty)
			var b2Record database.Book
			database.MustScan(t, "getting b2Record for test case",
				db.QueryRow("SELECT uuid, label, usn, dirty FROM books WHERE uuid = ?", b2UUID),
				&b2Record.UUID, &b2Record.Label, &b2Record.USN, &b2Record.Dirty)

			assert.Equal(t, b1Record.UUID, b1UUID, fmt.Sprintf("b1Record UUID mismatch for test case %d", idx))
			assert.Equal(t, b1Record.Label, "b1-label", fmt.Sprintf("b1Record Label mismatch for test case %d", idx))
			assert.Equal(t, b1Record.USN, 5, fmt.Sprintf("b1Record USN mismatch for test case %d", idx))
			assert.Equal(t, b1Record.Dirty, false, fmt.Sprintf("b1Record Dirty mismatch for test case %d", idx))

			assert.Equal(t, b2Record.UUID, b2UUID, fmt.Sprintf("b2Record UUID mismatch for test case %d", idx))
			assert.Equal(t, b2Record.Label, "b2-label", fmt.Sprintf("b2Record Label mismatch for test case %d", idx))
			assert.Equal(t, b2Record.USN, 6, fmt.Sprintf("b2Record USN mismatch for test case %d", idx))
			assert.Equal(t, b2Record.Dirty, false, fmt.Sprintf("b2Record Dirty mismatch for test case %d", idx))

			assert.Equal(t, n1Record.UUID, n1UUID, fmt.Sprintf("n1Record UUID mismatch for test case %d", idx))
			assert.Equal(t, n1Record.BookUUID, tc.expectedBookUUID, fmt.Sprintf("n1Record BookUUID mismatch for test case %d", idx))
			assert.Equal(t, n1Record.USN, tc.expectedUSN, fmt.Sprintf("n1Record USN mismatch for test case %d", idx))
			assert.Equal(t, n1Record.AddedOn, tc.expectedAddedOn, fmt.Sprintf("n1Record AddedOn mismatch for test case %d", idx))
			assert.Equal(t, n1Record.EditedOn, tc.expectedEditedOn, fmt.Sprintf("n1Record EditedOn mismatch for test case %d", idx))
			assert.Equal(t, n1Record.Body, tc.expectedBody, fmt.Sprintf("n1Record Body mismatch for test case %d", idx))
			assert.Equal(t, n1Record.Deleted, tc.expectedDeleted, fmt.Sprintf("n1Record Deleted mismatch for test case %d", idx))
			assert.Equal(t, n1Record.Dirty, tc.expectedDirty, fmt.Sprintf("n1Record Dirty mismatch for test case %d", idx))
		}()
	}
}

func TestCheckBookPristine(t *testing.T) {
	// set up
	db := database.InitTestDB(t, "../../tmp/.dnote", nil)
	defer database.CloseTestDB(t, db)

	database.MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, dirty) VALUES (?, ?, ?, ?)", "b1-uuid", "b1-label", 5, false)
	database.MustExec(t, "inserting b2", db, "INSERT INTO books (uuid, label, usn, dirty) VALUES (?, ?, ?, ?)", "b2-uuid", "b2-label", 6, false)
	database.MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, added_on, body, dirty) VALUES (?, ?, ?, ?, ?)", "n1-uuid", "b1-uuid", 1541108743, "n1 body", false)
	database.MustExec(t, "inserting n2", db, "INSERT INTO notes (uuid, book_uuid, added_on, body, dirty) VALUES (?, ?, ?, ?, ?)", "n2-uuid", "b1-uuid", 1541108743, "n2 body", false)
	database.MustExec(t, "inserting n3", db, "INSERT INTO notes (uuid, book_uuid, added_on, body, dirty) VALUES (?, ?, ?, ?, ?)", "n3-uuid", "b1-uuid", 1541108743, "n3 body", true)
	database.MustExec(t, "inserting n4", db, "INSERT INTO notes (uuid, book_uuid, added_on, body, dirty) VALUES (?, ?, ?, ?, ?)", "n4-uuid", "b2-uuid", 1541108743, "n4 body", false)
	database.MustExec(t, "inserting n5", db, "INSERT INTO notes (uuid, book_uuid, added_on, body, dirty) VALUES (?, ?, ?, ?, ?)", "n5-uuid", "b2-uuid", 1541108743, "n5 body", false)

	t.Run("b1", func(t *testing.T) {
		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}
		got, err := checkNotesPristine(tx, "b1-uuid")
		if err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		assert.Equal(t, got, false, "b1 should not be pristine")
	})

	t.Run("b2", func(t *testing.T) {
		// execute
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}
		got, err := checkNotesPristine(tx, "b2-uuid")
		if err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing").Error())
		}

		tx.Commit()

		// test
		assert.Equal(t, got, true, "b2 should be pristine")
	})
}

func TestCheckNoteInList(t *testing.T) {
	list := syncList{
		Notes: map[string]client.SyncFragNote{
			"n1-uuid": {
				UUID: "n1-uuid",
			},
			"n2-uuid": {
				UUID: "n2-uuid",
			},
		},
		Books: map[string]client.SyncFragBook{
			"b1-uuid": {
				UUID: "b1-uuid",
			},
			"b2-uuid": {
				UUID: "b2-uuid",
			},
		},
		ExpungedNotes: map[string]bool{
			"n3-uuid": true,
			"n4-uuid": true,
		},
		ExpungedBooks: map[string]bool{
			"b3-uuid": true,
			"b4-uuid": true,
		},
		MaxUSN:         1,
		MaxCurrentTime: 2,
	}

	testCases := []struct {
		uuid     string
		expected bool
	}{
		{
			uuid:     "n1-uuid",
			expected: true,
		},
		{
			uuid:     "n2-uuid",
			expected: true,
		},
		{
			uuid:     "n3-uuid",
			expected: true,
		},
		{
			uuid:     "n4-uuid",
			expected: true,
		},
		{
			uuid:     "nonexistent-note-uuid",
			expected: false,
		},
	}

	for idx, tc := range testCases {
		got := checkNoteInList(tc.uuid, &list)
		assert.Equal(t, got, tc.expected, fmt.Sprintf("result mismatch for test case %d", idx))
	}
}

func TestCheckBookInList(t *testing.T) {
	list := syncList{
		Notes: map[string]client.SyncFragNote{
			"n1-uuid": {
				UUID: "n1-uuid",
			},
			"n2-uuid": {
				UUID: "n2-uuid",
			},
		},
		Books: map[string]client.SyncFragBook{
			"b1-uuid": {
				UUID: "b1-uuid",
			},
			"b2-uuid": {
				UUID: "b2-uuid",
			},
		},
		ExpungedNotes: map[string]bool{
			"n3-uuid": true,
			"n4-uuid": true,
		},
		ExpungedBooks: map[string]bool{
			"b3-uuid": true,
			"b4-uuid": true,
		},
		MaxUSN:         1,
		MaxCurrentTime: 2,
	}

	testCases := []struct {
		uuid     string
		expected bool
	}{
		{
			uuid:     "b1-uuid",
			expected: true,
		},
		{
			uuid:     "b2-uuid",
			expected: true,
		},
		{
			uuid:     "b3-uuid",
			expected: true,
		},
		{
			uuid:     "b4-uuid",
			expected: true,
		},
		{
			uuid:     "nonexistent-book-uuid",
			expected: false,
		},
	}

	for idx, tc := range testCases {
		got := checkBookInList(tc.uuid, &list)
		assert.Equal(t, got, tc.expected, fmt.Sprintf("result mismatch for test case %d", idx))
	}
}

func TestCleanLocalNotes(t *testing.T) {
	// set up
	db := database.InitTestDB(t, "../../tmp/.dnote", nil)
	defer database.CloseTestDB(t, db)

	list := syncList{
		Notes: map[string]client.SyncFragNote{
			"n1-uuid": {
				UUID: "n1-uuid",
			},
			"n2-uuid": {
				UUID: "n2-uuid",
			},
		},
		Books: map[string]client.SyncFragBook{
			"b1-uuid": {
				UUID: "b1-uuid",
			},
			"b2-uuid": {
				UUID: "b2-uuid",
			},
		},
		ExpungedNotes: map[string]bool{
			"n3-uuid": true,
			"n4-uuid": true,
		},
		ExpungedBooks: map[string]bool{
			"b3-uuid": true,
			"b4-uuid": true,
		},
		MaxUSN:         1,
		MaxCurrentTime: 2,
	}

	b1UUID := "b1-uuid"
	database.MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b1UUID, "b1-label", 1, false, false)

	// exists in the list
	database.MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n1-uuid", b1UUID, 10, "n1 body", 1541108743, false, false)
	database.MustExec(t, "inserting n2", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n2-uuid", b1UUID, 0, "n2 body", 1541108743, false, true)
	// non-existent in the list but in valid state
	// (created in the cli and hasn't been uploaded)
	database.MustExec(t, "inserting n6", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n6-uuid", b1UUID, 0, "n6 body", 1541108743, false, true)
	// non-existent in the list and in an invalid state
	database.MustExec(t, "inserting n5", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n5-uuid", b1UUID, 7, "n5 body", 1541108743, true, true)
	database.MustExec(t, "inserting n9", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n9-uuid", b1UUID, 17, "n9 body", 1541108743, true, false)
	database.MustExec(t, "inserting n10", db, "INSERT INTO notes (uuid, book_uuid, usn, body, added_on, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", "n10-uuid", b1UUID, 0, "n10 body", 1541108743, false, false)

	// execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if err := cleanLocalNotes(tx, &list); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "executing").Error())
	}

	tx.Commit()

	// test
	var noteCount int
	database.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	assert.Equal(t, noteCount, 3, "note count mismatch")

	var n1, n2, n6 database.Note
	database.MustScan(t, "getting n1", db.QueryRow("SELECT dirty FROM notes WHERE uuid = ?", "n1-uuid"), &n1.Dirty)
	database.MustScan(t, "getting n2", db.QueryRow("SELECT dirty FROM notes WHERE uuid = ?", "n2-uuid"), &n2.Dirty)
	database.MustScan(t, "getting n6", db.QueryRow("SELECT dirty FROM notes WHERE uuid = ?", "n6-uuid"), &n6.Dirty)
}

func TestCleanLocalBooks(t *testing.T) {
	// set up
	db := database.InitTestDB(t, "../../tmp/.dnote", nil)
	defer database.CloseTestDB(t, db)

	list := syncList{
		Notes: map[string]client.SyncFragNote{
			"n1-uuid": {
				UUID: "n1-uuid",
			},
			"n2-uuid": {
				UUID: "n2-uuid",
			},
		},
		Books: map[string]client.SyncFragBook{
			"b1-uuid": {
				UUID: "b1-uuid",
			},
			"b2-uuid": {
				UUID: "b2-uuid",
			},
		},
		ExpungedNotes: map[string]bool{
			"n3-uuid": true,
			"n4-uuid": true,
		},
		ExpungedBooks: map[string]bool{
			"b3-uuid": true,
			"b4-uuid": true,
		},
		MaxUSN:         1,
		MaxCurrentTime: 2,
	}

	// existent in the server
	database.MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b1-uuid", "b1-label", 1, false, false)
	database.MustExec(t, "inserting b3", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b3-uuid", "b3-label", 0, false, true)
	// non-existent in the server but in valid state
	database.MustExec(t, "inserting b5", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b5-uuid", "b5-label", 0, true, true)
	// non-existent in the server and in an invalid state
	database.MustExec(t, "inserting b6", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b6-uuid", "b6-label", 10, true, true)
	database.MustExec(t, "inserting b7", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b7-uuid", "b7-label", 11, false, false)
	database.MustExec(t, "inserting b8", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", "b8-uuid", "b8-label", 0, false, false)

	// execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if err := cleanLocalBooks(tx, &list); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "executing").Error())
	}

	tx.Commit()

	// test
	var bookCount int
	database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	assert.Equal(t, bookCount, 3, "note count mismatch")

	var b1, b3, b5 database.Book
	database.MustScan(t, "getting b1", db.QueryRow("SELECT label FROM books WHERE uuid = ?", "b1-uuid"), &b1.Label)
	database.MustScan(t, "getting b3", db.QueryRow("SELECT label FROM books WHERE uuid = ?", "b3-uuid"), &b3.Label)
	database.MustScan(t, "getting b5", db.QueryRow("SELECT label FROM books WHERE uuid = ?", "b5-uuid"), &b5.Label)
}
