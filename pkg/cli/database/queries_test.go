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

package database

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/pkg/errors"
)

func TestInsertSystem(t *testing.T) {
	testCases := []struct {
		key string
		val string
	}{
		{
			key: "foo",
			val: "1558089284",
		},
		{
			key: "baz",
			val: "quz",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("insert %s %s", tc.key, tc.val), func(t *testing.T) {
			// Setup
			db := InitTestDB(t, "../tmp/dnote-test.db", nil)
			defer CloseTestDB(t, db)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
			}

			if err := InsertSystem(tx, tc.key, tc.val); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, "executing for test case").Error())
			}

			tx.Commit()

			// test
			var key, val string
			MustScan(t, "getting the saved record",
				db.QueryRow("SELECT key, value FROM system WHERE key = ?", tc.key), &key, &val)

			assert.Equal(t, key, tc.key, "key mismatch for test case")
			assert.Equal(t, val, tc.val, "val mismatch for test case")
		})
	}
}

func TestUpsertSystem(t *testing.T) {
	testCases := []struct {
		key        string
		val        string
		countDelta int
	}{
		{
			key:        "foo",
			val:        "1558089284",
			countDelta: 1,
		},
		{
			key:        "baz",
			val:        "quz2",
			countDelta: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("insert %s %s", tc.key, tc.val), func(t *testing.T) {
			// Setup
			db := InitTestDB(t, "../tmp/dnote-test.db", nil)
			defer CloseTestDB(t, db)

			MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "baz", "quz")

			var initialSystemCount int
			MustScan(t, "counting records", db.QueryRow("SELECT count(*) FROM system"), &initialSystemCount)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
			}

			if err := UpsertSystem(tx, tc.key, tc.val); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, "executing for test case").Error())
			}

			tx.Commit()

			// test
			var key, val string
			MustScan(t, "getting the saved record",
				db.QueryRow("SELECT key, value FROM system WHERE key = ?", tc.key), &key, &val)
			var systemCount int
			MustScan(t, "counting records",
				db.QueryRow("SELECT count(*) FROM system"), &systemCount)

			assert.Equal(t, key, tc.key, "key mismatch")
			assert.Equal(t, val, tc.val, "val mismatch")
			assert.Equal(t, systemCount, initialSystemCount+tc.countDelta, "count mismatch")
		})
	}
}

func TestGetSystem(t *testing.T) {
	t.Run(fmt.Sprintf("get string value"), func(t *testing.T) {
		// Setup
		db := InitTestDB(t, "../tmp/dnote-test.db", nil)
		defer CloseTestDB(t, db)

		// execute
		MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "foo", "bar")

		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}
		var dest string
		if err := GetSystem(tx, "foo", &dest); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing for test case").Error())
		}
		tx.Commit()

		// test
		assert.Equal(t, dest, "bar", "dest mismatch")
	})

	t.Run(fmt.Sprintf("get int64 value"), func(t *testing.T) {
		// Setup
		db := InitTestDB(t, "../tmp/dnote-test.db", nil)
		defer CloseTestDB(t, db)

		// execute
		MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "foo", 1234)

		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}
		var dest int64
		if err := GetSystem(tx, "foo", &dest); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing for test case").Error())
		}
		tx.Commit()

		// test
		assert.Equal(t, dest, int64(1234), "dest mismatch")
	})
}

func TestUpdateSystem(t *testing.T) {
	testCases := []struct {
		key        string
		val        string
		countDelta int
	}{
		{
			key: "foo",
			val: "1558089284",
		},
		{
			key: "foo",
			val: "bar",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("update %s %s", tc.key, tc.val), func(t *testing.T) {
			// Setup
			db := InitTestDB(t, "../tmp/dnote-test.db", nil)
			defer CloseTestDB(t, db)

			MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "foo", "fuz")
			MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "baz", "quz")

			var initialSystemCount int
			MustScan(t, "counting records", db.QueryRow("SELECT count(*) FROM system"), &initialSystemCount)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
			}

			if err := UpdateSystem(tx, tc.key, tc.val); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, "executing for test case").Error())
			}

			tx.Commit()

			// test
			var key, val string
			MustScan(t, "getting the saved record",
				db.QueryRow("SELECT key, value FROM system WHERE key = ?", tc.key), &key, &val)
			var systemCount int
			MustScan(t, "counting records",
				db.QueryRow("SELECT count(*) FROM system"), &systemCount)

			assert.Equal(t, key, tc.key, "key mismatch")
			assert.Equal(t, val, tc.val, "val mismatch")
			assert.Equal(t, systemCount, initialSystemCount, "count mismatch")
		})
	}
}

func TestGetActiveNote(t *testing.T) {
	t.Run("not deleted", func(t *testing.T) {
		// set up
		db := InitTestDB(t, "../tmp/dnote-test.db", nil)
		defer CloseTestDB(t, db)

		n1UUID := "n1-uuid"
		MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, edited_on, usn, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", n1UUID, "b1-uuid", "n1 content", 1542058875, 1542058876, 1, true, false, true)

		var n1RowID int
		MustScan(t, "getting rowid", db.QueryRow("SELECT rowid FROM notes WHERE uuid = ?", n1UUID), &n1RowID)

		// execute
		got, err := GetActiveNote(db, n1RowID)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		// test
		assert.Equal(t, got.RowID, n1RowID, "RowID mismatch")
		assert.Equal(t, got.UUID, n1UUID, "UUID mismatch")
		assert.Equal(t, got.BookUUID, "b1-uuid", "BookUUID mismatch")
		assert.Equal(t, got.Body, "n1 content", "Body mismatch")
		assert.Equal(t, got.AddedOn, int64(1542058875), "AddedOn mismatch")
		assert.Equal(t, got.EditedOn, int64(1542058876), "EditedOn mismatch")
		assert.Equal(t, got.USN, 1, "USN mismatch")
		assert.Equal(t, got.Public, true, "Public mismatch")
		assert.Equal(t, got.Deleted, false, "Deleted mismatch")
		assert.Equal(t, got.Dirty, true, "Dirty mismatch")
	})

	t.Run("deleted", func(t *testing.T) {
		// set up
		db := InitTestDB(t, "../tmp/dnote-test.db", nil)
		defer CloseTestDB(t, db)

		n1UUID := "n1-uuid"
		MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, edited_on, usn, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", n1UUID, "b1-uuid", "n1 content", 1542058875, 1542058876, 1, true, true, true)

		var n1RowID int
		MustScan(t, "getting rowid", db.QueryRow("SELECT rowid FROM notes WHERE uuid = ?", n1UUID), &n1RowID)

		// execute
		_, err := GetActiveNote(db, n1RowID)

		// test
		if err == nil {
			t.Error("Should have returned an error")
		}
		if err != nil && err != sql.ErrNoRows {
			t.Error(errors.Wrap(err, "executing"))
		}
	})
}

func TestUpdateNoteContent(t *testing.T) {
	// set up
	db := InitTestDB(t, "../tmp/dnote-test.db", nil)
	defer CloseTestDB(t, db)

	uuid := "n1-uuid"
	MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, edited_on, usn, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", uuid, "b1-uuid", "n1 content", 1542058875, 0, 1, false, false, false)

	var rowid int
	MustScan(t, "getting rowid", db.QueryRow("SELECT rowid FROM notes WHERE uuid = ?", uuid), &rowid)

	// execute
	c := clock.NewMock()
	now := time.Date(2017, time.March, 14, 21, 15, 0, 0, time.UTC)
	c.SetNow(now)

	err := UpdateNoteContent(db, c, rowid, "n1 content updated")
	if err != nil {
		t.Fatal(errors.Wrap(err, "executing"))
	}

	var content string
	var editedOn int
	var dirty bool

	MustScan(t, "getting the note record", db.QueryRow("SELECT body, edited_on, dirty FROM notes WHERE rowid = ?", rowid), &content, &editedOn, &dirty)

	assert.Equal(t, content, "n1 content updated", "content mismatch")
	assert.Equal(t, int64(editedOn), now.UnixNano(), "editedOn mismatch")
	assert.Equal(t, dirty, true, "dirty mismatch")
}

func TestUpdateNoteBook(t *testing.T) {
	// set up
	db := InitTestDB(t, "../tmp/dnote-test.db", nil)
	defer CloseTestDB(t, db)

	b1UUID := "b1-uuid"
	b2UUID := "b2-uuid"
	MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b1UUID, "b1-label", 8, false, false)
	MustExec(t, "inserting b2", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b2UUID, "b2-label", 9, false, false)

	uuid := "n1-uuid"
	MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, edited_on, usn, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", uuid, b1UUID, "n1 content", 1542058875, 0, 1, false, false, false)

	var rowid int
	MustScan(t, "getting rowid", db.QueryRow("SELECT rowid FROM notes WHERE uuid = ?", uuid), &rowid)

	// execute
	c := clock.NewMock()
	now := time.Date(2017, time.March, 14, 21, 15, 0, 0, time.UTC)
	c.SetNow(now)

	err := UpdateNoteBook(db, c, rowid, b2UUID)
	if err != nil {
		t.Fatal(errors.Wrap(err, "executing"))
	}

	var bookUUID string
	var editedOn int
	var dirty bool

	MustScan(t, "getting the note record", db.QueryRow("SELECT book_uuid, edited_on, dirty FROM notes WHERE rowid = ?", rowid), &bookUUID, &editedOn, &dirty)

	assert.Equal(t, bookUUID, b2UUID, "content mismatch")
	assert.Equal(t, int64(editedOn), now.UnixNano(), "editedOn mismatch")
	assert.Equal(t, dirty, true, "dirty mismatch")
}

func TestUpdateBookName(t *testing.T) {
	// set up
	db := InitTestDB(t, "../tmp/dnote-test.db", nil)
	defer CloseTestDB(t, db)

	b1UUID := "b1-uuid"
	MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b1UUID, "b1-label", 8, false, false)

	// execute
	err := UpdateBookName(db, b1UUID, "b1-label-edited")
	if err != nil {
		t.Fatal(errors.Wrap(err, "executing"))
	}

	// test
	var b1 Book
	MustScan(t, "getting the note record", db.QueryRow("SELECT uuid, label, dirty, usn, deleted FROM books WHERE uuid = ?", b1UUID), &b1.UUID, &b1.Label, &b1.Dirty, &b1.USN, &b1.Deleted)
	assert.Equal(t, b1.UUID, b1UUID, "UUID mismatch")
	assert.Equal(t, b1.Label, "b1-label-edited", "Label mismatch")
	assert.Equal(t, b1.Dirty, true, "Dirty mismatch")
	assert.Equal(t, b1.USN, 8, "USN mismatch")
	assert.Equal(t, b1.Deleted, false, "Deleted mismatch")
}
