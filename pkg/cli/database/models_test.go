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
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/pkg/errors"
)

func TestNewNote(t *testing.T) {
	testCases := []struct {
		uuid     string
		bookUUID string
		body     string
		addedOn  int64
		editedOn int64
		usn      int
		public   bool
		deleted  bool
		dirty    bool
	}{
		{
			uuid:     "n1-uuid",
			bookUUID: "b1-uuid",
			body:     "n1-body",
			addedOn:  1542058875,
			editedOn: 0,
			usn:      0,
			public:   false,
			deleted:  false,
			dirty:    false,
		},
		{
			uuid:     "n2-uuid",
			bookUUID: "b2-uuid",
			body:     "n2-body",
			addedOn:  1542058875,
			editedOn: 1542058876,
			usn:      1008,
			public:   true,
			deleted:  true,
			dirty:    true,
		},
	}

	for idx, tc := range testCases {
		got := NewNote(tc.uuid, tc.bookUUID, tc.body, tc.addedOn, tc.editedOn, tc.usn, tc.public, tc.deleted, tc.dirty)

		assert.Equal(t, got.UUID, tc.uuid, fmt.Sprintf("UUID mismatch for test case %d", idx))
		assert.Equal(t, got.BookUUID, tc.bookUUID, fmt.Sprintf("BookUUID mismatch for test case %d", idx))
		assert.Equal(t, got.Body, tc.body, fmt.Sprintf("Body mismatch for test case %d", idx))
		assert.Equal(t, got.AddedOn, tc.addedOn, fmt.Sprintf("AddedOn mismatch for test case %d", idx))
		assert.Equal(t, got.EditedOn, tc.editedOn, fmt.Sprintf("EditedOn mismatch for test case %d", idx))
		assert.Equal(t, got.USN, tc.usn, fmt.Sprintf("USN mismatch for test case %d", idx))
		assert.Equal(t, got.Public, tc.public, fmt.Sprintf("Public mismatch for test case %d", idx))
		assert.Equal(t, got.Deleted, tc.deleted, fmt.Sprintf("Deleted mismatch for test case %d", idx))
		assert.Equal(t, got.Dirty, tc.dirty, fmt.Sprintf("Dirty mismatch for test case %d", idx))
	}
}

func TestNoteInsert(t *testing.T) {
	testCases := []struct {
		uuid     string
		bookUUID string
		body     string
		addedOn  int64
		editedOn int64
		usn      int
		public   bool
		deleted  bool
		dirty    bool
	}{
		{
			uuid:     "n1-uuid",
			bookUUID: "b1-uuid",
			body:     "n1-body",
			addedOn:  1542058875,
			editedOn: 0,
			usn:      0,
			public:   false,
			deleted:  false,
			dirty:    false,
		},
		{
			uuid:     "n2-uuid",
			bookUUID: "b2-uuid",
			body:     "n2-body",
			addedOn:  1542058875,
			editedOn: 1542058876,
			usn:      1008,
			public:   true,
			deleted:  true,
			dirty:    true,
		},
	}

	for idx, tc := range testCases {
		func() {
			// Setup
			db := InitTestDB(t, "../tmp/dnote-test.db", nil)
			defer CloseTestDB(t, db)

			n := Note{
				UUID:     tc.uuid,
				BookUUID: tc.bookUUID,
				Body:     tc.body,
				AddedOn:  tc.addedOn,
				EditedOn: tc.editedOn,
				USN:      tc.usn,
				Public:   tc.public,
				Deleted:  tc.deleted,
				Dirty:    tc.dirty,
			}

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
			}

			if err := n.Insert(tx); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
			}

			tx.Commit()

			// test
			var uuid, bookUUID, body string
			var addedOn, editedOn int64
			var usn int
			var public, deleted, dirty bool
			MustScan(t, "getting n1",
				db.QueryRow("SELECT uuid, book_uuid, body, added_on, edited_on, usn, public, deleted, dirty FROM notes WHERE uuid = ?", tc.uuid),
				&uuid, &bookUUID, &body, &addedOn, &editedOn, &usn, &public, &deleted, &dirty)

			assert.Equal(t, uuid, tc.uuid, fmt.Sprintf("uuid mismatch for test case %d", idx))
			assert.Equal(t, bookUUID, tc.bookUUID, fmt.Sprintf("bookUUID mismatch for test case %d", idx))
			assert.Equal(t, body, tc.body, fmt.Sprintf("body mismatch for test case %d", idx))
			assert.Equal(t, addedOn, tc.addedOn, fmt.Sprintf("addedOn mismatch for test case %d", idx))
			assert.Equal(t, editedOn, tc.editedOn, fmt.Sprintf("editedOn mismatch for test case %d", idx))
			assert.Equal(t, usn, tc.usn, fmt.Sprintf("usn mismatch for test case %d", idx))
			assert.Equal(t, public, tc.public, fmt.Sprintf("public mismatch for test case %d", idx))
			assert.Equal(t, deleted, tc.deleted, fmt.Sprintf("deleted mismatch for test case %d", idx))
			assert.Equal(t, dirty, tc.dirty, fmt.Sprintf("dirty mismatch for test case %d", idx))
		}()
	}
}

func TestNoteUpdate(t *testing.T) {
	testCases := []struct {
		uuid        string
		bookUUID    string
		body        string
		addedOn     int64
		editedOn    int64
		usn         int
		public      bool
		deleted     bool
		dirty       bool
		newBookUUID string
		newBody     string
		newEditedOn int64
		newUSN      int
		newPublic   bool
		newDeleted  bool
		newDirty    bool
	}{
		{
			uuid:        "n1-uuid",
			bookUUID:    "b1-uuid",
			body:        "n1-body",
			addedOn:     1542058875,
			editedOn:    0,
			usn:         0,
			public:      false,
			deleted:     false,
			dirty:       false,
			newBookUUID: "b1-uuid",
			newBody:     "n1-body edited",
			newEditedOn: 1542058879,
			newUSN:      0,
			newPublic:   false,
			newDeleted:  false,
			newDirty:    false,
		},
		{
			uuid:        "n1-uuid",
			bookUUID:    "b1-uuid",
			body:        "n1-body",
			addedOn:     1542058875,
			editedOn:    0,
			usn:         0,
			public:      false,
			deleted:     false,
			dirty:       true,
			newBookUUID: "b2-uuid",
			newBody:     "n1-body",
			newEditedOn: 1542058879,
			newUSN:      0,
			newPublic:   true,
			newDeleted:  false,
			newDirty:    false,
		},
		{
			uuid:        "n1-uuid",
			bookUUID:    "b1-uuid",
			body:        "n1-body",
			addedOn:     1542058875,
			editedOn:    0,
			usn:         10,
			public:      false,
			deleted:     false,
			dirty:       false,
			newBookUUID: "",
			newBody:     "",
			newEditedOn: 1542058879,
			newUSN:      151,
			newPublic:   false,
			newDeleted:  true,
			newDirty:    false,
		},
		{
			uuid:        "n1-uuid",
			bookUUID:    "b1-uuid",
			body:        "n1-body",
			addedOn:     1542058875,
			editedOn:    0,
			usn:         0,
			public:      false,
			deleted:     false,
			dirty:       false,
			newBookUUID: "",
			newBody:     "",
			newEditedOn: 1542058879,
			newUSN:      15,
			newPublic:   false,
			newDeleted:  true,
			newDirty:    false,
		},
	}

	for idx, tc := range testCases {
		func() {
			// Setup
			db := InitTestDB(t, "../tmp/dnote-test.db", nil)
			defer CloseTestDB(t, db)

			n1 := Note{
				UUID:     tc.uuid,
				BookUUID: tc.bookUUID,
				Body:     tc.body,
				AddedOn:  tc.addedOn,
				EditedOn: tc.editedOn,
				USN:      tc.usn,
				Public:   tc.public,
				Deleted:  tc.deleted,
				Dirty:    tc.dirty,
			}
			n2 := Note{
				UUID:     "n2-uuid",
				BookUUID: "b10-uuid",
				Body:     "n2 body",
				AddedOn:  1542058875,
				EditedOn: 0,
				USN:      39,
				Public:   false,
				Deleted:  false,
				Dirty:    false,
			}

			MustExec(t, fmt.Sprintf("inserting n1 for test case %d", idx), db, "INSERT INTO notes (uuid, book_uuid, usn, added_on, edited_on, body, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", n1.UUID, n1.BookUUID, n1.USN, n1.AddedOn, n1.EditedOn, n1.Body, n1.Public, n1.Deleted, n1.Dirty)
			MustExec(t, fmt.Sprintf("inserting n2 for test case %d", idx), db, "INSERT INTO notes (uuid, book_uuid, usn, added_on, edited_on, body, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", n2.UUID, n2.BookUUID, n2.USN, n2.AddedOn, n2.EditedOn, n2.Body, n2.Public, n2.Deleted, n2.Dirty)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
			}

			n1.BookUUID = tc.newBookUUID
			n1.Body = tc.newBody
			n1.EditedOn = tc.newEditedOn
			n1.USN = tc.newUSN
			n1.Public = tc.newPublic
			n1.Deleted = tc.newDeleted
			n1.Dirty = tc.newDirty

			if err := n1.Update(tx); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
			}

			tx.Commit()

			// test
			var n1Record, n2Record Note
			MustScan(t, "getting n1",
				db.QueryRow("SELECT uuid, book_uuid, body, added_on, edited_on, usn, public, deleted, dirty FROM notes WHERE uuid = ?", tc.uuid),
				&n1Record.UUID, &n1Record.BookUUID, &n1Record.Body, &n1Record.AddedOn, &n1Record.EditedOn, &n1Record.USN, &n1Record.Public, &n1Record.Deleted, &n1Record.Dirty)
			MustScan(t, "getting n2",
				db.QueryRow("SELECT uuid, book_uuid, body, added_on, edited_on, usn, public, deleted, dirty FROM notes WHERE uuid = ?", n2.UUID),
				&n2Record.UUID, &n2Record.BookUUID, &n2Record.Body, &n2Record.AddedOn, &n2Record.EditedOn, &n2Record.USN, &n2Record.Public, &n2Record.Deleted, &n2Record.Dirty)

			assert.Equal(t, n1Record.UUID, n1.UUID, fmt.Sprintf("n1 uuid mismatch for test case %d", idx))
			assert.Equal(t, n1Record.BookUUID, tc.newBookUUID, fmt.Sprintf("n1 bookUUID mismatch for test case %d", idx))
			assert.Equal(t, n1Record.Body, tc.newBody, fmt.Sprintf("n1 body mismatch for test case %d", idx))
			assert.Equal(t, n1Record.AddedOn, n1.AddedOn, fmt.Sprintf("n1 addedOn mismatch for test case %d", idx))
			assert.Equal(t, n1Record.EditedOn, tc.newEditedOn, fmt.Sprintf("n1 editedOn mismatch for test case %d", idx))
			assert.Equal(t, n1Record.USN, tc.newUSN, fmt.Sprintf("n1 usn mismatch for test case %d", idx))
			assert.Equal(t, n1Record.Public, tc.newPublic, fmt.Sprintf("n1 public mismatch for test case %d", idx))
			assert.Equal(t, n1Record.Deleted, tc.newDeleted, fmt.Sprintf("n1 deleted mismatch for test case %d", idx))
			assert.Equal(t, n1Record.Dirty, tc.newDirty, fmt.Sprintf("n1 dirty mismatch for test case %d", idx))

			assert.Equal(t, n2Record.UUID, n2.UUID, fmt.Sprintf("n2 uuid mismatch for test case %d", idx))
			assert.Equal(t, n2Record.BookUUID, n2.BookUUID, fmt.Sprintf("n2 bookUUID mismatch for test case %d", idx))
			assert.Equal(t, n2Record.Body, n2.Body, fmt.Sprintf("n2 body mismatch for test case %d", idx))
			assert.Equal(t, n2Record.AddedOn, n2.AddedOn, fmt.Sprintf("n2 addedOn mismatch for test case %d", idx))
			assert.Equal(t, n2Record.EditedOn, n2.EditedOn, fmt.Sprintf("n2 editedOn mismatch for test case %d", idx))
			assert.Equal(t, n2Record.USN, n2.USN, fmt.Sprintf("n2 usn mismatch for test case %d", idx))
			assert.Equal(t, n2Record.Public, n2.Public, fmt.Sprintf("n2 public mismatch for test case %d", idx))
			assert.Equal(t, n2Record.Deleted, n2.Deleted, fmt.Sprintf("n2 deleted mismatch for test case %d", idx))
			assert.Equal(t, n2Record.Dirty, n2.Dirty, fmt.Sprintf("n2 dirty mismatch for test case %d", idx))
		}()
	}
}

func TestNoteUpdateUUID(t *testing.T) {
	testCases := []struct {
		newUUID string
	}{
		{
			newUUID: "n1-new-uuid",
		},
		{
			newUUID: "n2-new-uuid",
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("testCase%d", idx), func(t *testing.T) {
			// Setup
			db := InitTestDB(t, "../tmp/dnote-test.db", nil)
			defer CloseTestDB(t, db)

			n1 := Note{
				UUID:     "n1-uuid",
				BookUUID: "b1-uuid",
				AddedOn:  1542058874,
				Body:     "n1-body",
				USN:      1,
				Deleted:  true,
				Dirty:    false,
			}
			n2 := Note{
				UUID:     "n2-uuid",
				BookUUID: "b1-uuid",
				AddedOn:  1542058874,
				Body:     "n2-body",
				USN:      1,
				Deleted:  true,
				Dirty:    false,
			}

			MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", n1.UUID, n1.BookUUID, n1.Body, n1.AddedOn, n1.USN, n1.Deleted, n1.Dirty)
			MustExec(t, "inserting n2", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?)", n2.UUID, n2.BookUUID, n2.Body, n2.AddedOn, n2.USN, n2.Deleted, n2.Dirty)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
			}
			if err := n1.UpdateUUID(tx, tc.newUUID); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, "executing").Error())
			}

			tx.Commit()

			// test
			var n1Record, n2Record Note
			MustScan(t, "getting n1",
				db.QueryRow("SELECT uuid, body, usn, deleted, dirty FROM notes WHERE body = ?", "n1-body"),
				&n1Record.UUID, &n1Record.Body, &n1Record.USN, &n1Record.Deleted, &n1Record.Dirty)
			MustScan(t, "getting n2",
				db.QueryRow("SELECT uuid, body, usn, deleted, dirty FROM notes WHERE body = ?", "n2-body"),
				&n2Record.UUID, &n2Record.Body, &n2Record.USN, &n2Record.Deleted, &n2Record.Dirty)

			assert.Equal(t, n1.UUID, tc.newUUID, "n1 original reference uuid mismatch")
			assert.Equal(t, n1Record.UUID, tc.newUUID, "n1 uuid mismatch")
			assert.Equal(t, n2Record.UUID, n2.UUID, "n2 uuid mismatch")
		})
	}
}

func TestNoteExpunge(t *testing.T) {
	// Setup
	db := InitTestDB(t, "../tmp/dnote-test.db", nil)
	defer CloseTestDB(t, db)

	n1 := Note{
		UUID:     "n1-uuid",
		BookUUID: "b9-uuid",
		Body:     "n1 body",
		AddedOn:  1542058874,
		EditedOn: 0,
		USN:      22,
		Public:   false,
		Deleted:  false,
		Dirty:    false,
	}
	n2 := Note{
		UUID:     "n2-uuid",
		BookUUID: "b10-uuid",
		Body:     "n2 body",
		AddedOn:  1542058875,
		EditedOn: 0,
		USN:      39,
		Public:   false,
		Deleted:  false,
		Dirty:    false,
	}

	MustExec(t, "inserting n1", db, "INSERT INTO notes (uuid, book_uuid, usn, added_on, edited_on, body, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", n1.UUID, n1.BookUUID, n1.USN, n1.AddedOn, n1.EditedOn, n1.Body, n1.Public, n1.Deleted, n1.Dirty)
	MustExec(t, "inserting n2", db, "INSERT INTO notes (uuid, book_uuid, usn, added_on, edited_on, body, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", n2.UUID, n2.BookUUID, n2.USN, n2.AddedOn, n2.EditedOn, n2.Body, n2.Public, n2.Deleted, n2.Dirty)

	// execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if err := n1.Expunge(tx); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "executing").Error())
	}

	tx.Commit()

	// test
	var noteCount int
	MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)

	assert.Equalf(t, noteCount, 1, "note count mismatch")

	var n2Record Note
	MustScan(t, "getting n2",
		db.QueryRow("SELECT uuid, book_uuid, body, added_on, edited_on, usn, public, deleted, dirty FROM notes WHERE uuid = ?", n2.UUID),
		&n2Record.UUID, &n2Record.BookUUID, &n2Record.Body, &n2Record.AddedOn, &n2Record.EditedOn, &n2Record.USN, &n2Record.Public, &n2Record.Deleted, &n2Record.Dirty)

	assert.Equal(t, n2Record.UUID, n2.UUID, "n2 uuid mismatch")
	assert.Equal(t, n2Record.BookUUID, n2.BookUUID, "n2 bookUUID mismatch")
	assert.Equal(t, n2Record.Body, n2.Body, "n2 body mismatch")
	assert.Equal(t, n2Record.AddedOn, n2.AddedOn, "n2 addedOn mismatch")
	assert.Equal(t, n2Record.EditedOn, n2.EditedOn, "n2 editedOn mismatch")
	assert.Equal(t, n2Record.USN, n2.USN, "n2 usn mismatch")
	assert.Equal(t, n2Record.Public, n2.Public, "n2 public mismatch")
	assert.Equal(t, n2Record.Deleted, n2.Deleted, "n2 deleted mismatch")
	assert.Equal(t, n2Record.Dirty, n2.Dirty, "n2 dirty mismatch")
}

func TestNewBook(t *testing.T) {
	testCases := []struct {
		uuid    string
		label   string
		usn     int
		deleted bool
		dirty   bool
	}{
		{
			uuid:    "b1-uuid",
			label:   "b1-label",
			usn:     0,
			deleted: false,
			dirty:   false,
		},
		{
			uuid:    "b2-uuid",
			label:   "b2-label",
			usn:     1008,
			deleted: false,
			dirty:   true,
		},
	}

	for idx, tc := range testCases {
		got := NewBook(tc.uuid, tc.label, tc.usn, tc.deleted, tc.dirty)

		assert.Equal(t, got.UUID, tc.uuid, fmt.Sprintf("UUID mismatch for test case %d", idx))
		assert.Equal(t, got.Label, tc.label, fmt.Sprintf("Label mismatch for test case %d", idx))
		assert.Equal(t, got.USN, tc.usn, fmt.Sprintf("USN mismatch for test case %d", idx))
		assert.Equal(t, got.Deleted, tc.deleted, fmt.Sprintf("Deleted mismatch for test case %d", idx))
		assert.Equal(t, got.Dirty, tc.dirty, fmt.Sprintf("Dirty mismatch for test case %d", idx))
	}
}

func TestBookInsert(t *testing.T) {
	testCases := []struct {
		uuid    string
		label   string
		usn     int
		deleted bool
		dirty   bool
	}{
		{
			uuid:    "b1-uuid",
			label:   "b1-label",
			usn:     10808,
			deleted: false,
			dirty:   false,
		},
		{
			uuid:    "b1-uuid",
			label:   "b1-label",
			usn:     10808,
			deleted: false,
			dirty:   true,
		},
	}

	for idx, tc := range testCases {
		func() {
			// Setup
			db := InitTestDB(t, "../tmp/dnote-test.db", nil)
			defer CloseTestDB(t, db)

			b := Book{
				UUID:    tc.uuid,
				Label:   tc.label,
				USN:     tc.usn,
				Dirty:   tc.dirty,
				Deleted: tc.deleted,
			}

			// execute

			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
			}

			if err := b.Insert(tx); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
			}

			tx.Commit()

			// test
			var uuid, label string
			var usn int
			var deleted, dirty bool
			MustScan(t, "getting b1",
				db.QueryRow("SELECT uuid, label, usn, deleted, dirty FROM books WHERE uuid = ?", tc.uuid),
				&uuid, &label, &usn, &deleted, &dirty)

			assert.Equal(t, uuid, tc.uuid, fmt.Sprintf("uuid mismatch for test case %d", idx))
			assert.Equal(t, label, tc.label, fmt.Sprintf("label mismatch for test case %d", idx))
			assert.Equal(t, usn, tc.usn, fmt.Sprintf("usn mismatch for test case %d", idx))
			assert.Equal(t, deleted, tc.deleted, fmt.Sprintf("deleted mismatch for test case %d", idx))
			assert.Equal(t, dirty, tc.dirty, fmt.Sprintf("dirty mismatch for test case %d", idx))
		}()
	}
}

func TestBookUpdate(t *testing.T) {
	testCases := []struct {
		uuid       string
		label      string
		usn        int
		deleted    bool
		dirty      bool
		newLabel   string
		newUSN     int
		newDeleted bool
		newDirty   bool
	}{
		{
			uuid:       "b1-uuid",
			label:      "b1-label",
			usn:        0,
			deleted:    false,
			dirty:      false,
			newLabel:   "b1-label-edited",
			newUSN:     0,
			newDeleted: false,
			newDirty:   true,
		},
		{
			uuid:       "b1-uuid",
			label:      "b1-label",
			usn:        0,
			deleted:    false,
			dirty:      false,
			newLabel:   "",
			newUSN:     10,
			newDeleted: true,
			newDirty:   false,
		},
	}

	for idx, tc := range testCases {
		func() {
			// Setup
			db := InitTestDB(t, "../tmp/dnote-test.db", nil)
			defer CloseTestDB(t, db)

			b1 := Book{
				UUID:    "b1-uuid",
				Label:   "b1-label",
				USN:     1,
				Deleted: true,
				Dirty:   false,
			}
			b2 := Book{
				UUID:    "b2-uuid",
				Label:   "b2-label",
				USN:     1,
				Deleted: true,
				Dirty:   false,
			}

			MustExec(t, fmt.Sprintf("inserting b1 for test case %d", idx), db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b1.UUID, b1.Label, b1.USN, b1.Deleted, b1.Dirty)
			MustExec(t, fmt.Sprintf("inserting b2 for test case %d", idx), db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b2.UUID, b2.Label, b2.USN, b2.Deleted, b2.Dirty)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("beginning a transaction for test case %d", idx)).Error())
			}

			b1.Label = tc.newLabel
			b1.USN = tc.newUSN
			b1.Deleted = tc.newDeleted
			b1.Dirty = tc.newDirty

			if err := b1.Update(tx); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, fmt.Sprintf("executing for test case %d", idx)).Error())
			}

			tx.Commit()

			// test
			var b1Record, b2Record Book
			MustScan(t, "getting b1",
				db.QueryRow("SELECT uuid, label, usn, deleted, dirty FROM books WHERE uuid = ?", tc.uuid),
				&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Deleted, &b1Record.Dirty)
			MustScan(t, "getting b2",
				db.QueryRow("SELECT uuid, label, usn, deleted, dirty FROM books WHERE uuid = ?", b2.UUID),
				&b2Record.UUID, &b2Record.Label, &b2Record.USN, &b2Record.Deleted, &b2Record.Dirty)

			assert.Equal(t, b1Record.UUID, b1.UUID, fmt.Sprintf("b1 uuid mismatch for test case %d", idx))
			assert.Equal(t, b1Record.Label, tc.newLabel, fmt.Sprintf("b1 label mismatch for test case %d", idx))
			assert.Equal(t, b1Record.USN, tc.newUSN, fmt.Sprintf("b1 usn mismatch for test case %d", idx))
			assert.Equal(t, b1Record.Deleted, tc.newDeleted, fmt.Sprintf("b1 deleted mismatch for test case %d", idx))
			assert.Equal(t, b1Record.Dirty, tc.newDirty, fmt.Sprintf("b1 dirty mismatch for test case %d", idx))

			assert.Equal(t, b2Record.UUID, b2.UUID, fmt.Sprintf("b2 uuid mismatch for test case %d", idx))
			assert.Equal(t, b2Record.Label, b2.Label, fmt.Sprintf("b2 label mismatch for test case %d", idx))
			assert.Equal(t, b2Record.USN, b2.USN, fmt.Sprintf("b2 usn mismatch for test case %d", idx))
			assert.Equal(t, b2Record.Deleted, b2.Deleted, fmt.Sprintf("b2 deleted mismatch for test case %d", idx))
			assert.Equal(t, b2Record.Dirty, b2.Dirty, fmt.Sprintf("b2 dirty mismatch for test case %d", idx))
		}()
	}
}

func TestBookUpdateUUID(t *testing.T) {
	testCases := []struct {
		newUUID string
	}{
		{
			newUUID: "b1-new-uuid",
		},
		{
			newUUID: "b2-new-uuid",
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("testCase%d", idx), func(t *testing.T) {

			// Setup
			db := InitTestDB(t, "../tmp/dnote-test.db", nil)
			defer CloseTestDB(t, db)

			b1 := Book{
				UUID:    "b1-uuid",
				Label:   "b1-label",
				USN:     1,
				Deleted: true,
				Dirty:   false,
			}
			b2 := Book{
				UUID:    "b2-uuid",
				Label:   "b2-label",
				USN:     1,
				Deleted: true,
				Dirty:   false,
			}

			MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b1.UUID, b1.Label, b1.USN, b1.Deleted, b1.Dirty)
			MustExec(t, "inserting b2", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b2.UUID, b2.Label, b2.USN, b2.Deleted, b2.Dirty)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
			}
			if err := b1.UpdateUUID(tx, tc.newUUID); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, "executing").Error())
			}

			tx.Commit()

			// test
			var b1Record, b2Record Book
			MustScan(t, "getting b1",
				db.QueryRow("SELECT uuid, label, usn, deleted, dirty FROM books WHERE label = ?", "b1-label"),
				&b1Record.UUID, &b1Record.Label, &b1Record.USN, &b1Record.Deleted, &b1Record.Dirty)
			MustScan(t, "getting b2",
				db.QueryRow("SELECT uuid, label, usn, deleted, dirty FROM books WHERE label = ?", "b2-label"),
				&b2Record.UUID, &b2Record.Label, &b2Record.USN, &b2Record.Deleted, &b2Record.Dirty)

			assert.Equal(t, b1.UUID, tc.newUUID, "b1 original reference uuid mismatch")
			assert.Equal(t, b1Record.UUID, tc.newUUID, "b1 uuid mismatch")
			assert.Equal(t, b2Record.UUID, b2.UUID, "b2 uuid mismatch")
		})
	}
}

func TestBookExpunge(t *testing.T) {
	// Setup
	db := InitTestDB(t, "../tmp/dnote-test.db", nil)
	defer CloseTestDB(t, db)

	b1 := Book{
		UUID:    "b1-uuid",
		Label:   "b1-label",
		USN:     1,
		Deleted: true,
		Dirty:   false,
	}
	b2 := Book{
		UUID:    "b2-uuid",
		Label:   "b2-label",
		USN:     1,
		Deleted: true,
		Dirty:   false,
	}

	MustExec(t, "inserting b1", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b1.UUID, b1.Label, b1.USN, b1.Deleted, b1.Dirty)
	MustExec(t, "inserting b2", db, "INSERT INTO books (uuid, label, usn, deleted, dirty) VALUES (?, ?, ?, ?, ?)", b2.UUID, b2.Label, b2.USN, b2.Deleted, b2.Dirty)

	// execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if err := b1.Expunge(tx); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "executing").Error())
	}

	tx.Commit()

	// test
	var bookCount int
	MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)

	assert.Equalf(t, bookCount, 1, "book count mismatch")

	var b2Record Book
	MustScan(t, "getting b2",
		db.QueryRow("SELECT uuid, label, usn, deleted, dirty FROM books WHERE uuid = ?", "b2-uuid"),
		&b2Record.UUID, &b2Record.Label, &b2Record.USN, &b2Record.Deleted, &b2Record.Dirty)

	assert.Equal(t, b2Record.UUID, b2.UUID, "b2 uuid mismatch")
	assert.Equal(t, b2Record.Label, b2.Label, "b2 label mismatch")
	assert.Equal(t, b2Record.USN, b2.USN, "b2 usn mismatch")
	assert.Equal(t, b2Record.Deleted, b2.Deleted, "b2 deleted mismatch")
	assert.Equal(t, b2Record.Dirty, b2.Dirty, "b2 dirty mismatch")
}

// TestNoteFTS tests that note full text search indices stay in sync with the notes after insert, update and delete
func TestNoteFTS(t *testing.T) {
	// set up
	db := InitTestDB(t, "../tmp/dnote-test.db", nil)
	defer CloseTestDB(t, db)

	// execute - insert
	n := Note{
		UUID:     "n1-uuid",
		BookUUID: "b1-uuid",
		Body:     "foo bar",
		AddedOn:  1542058875,
		EditedOn: 0,
		USN:      0,
		Public:   false,
		Deleted:  false,
		Dirty:    false,
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if err := n.Insert(tx); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "inserting").Error())
	}

	tx.Commit()

	// test
	var noteCount, noteFtsCount, noteSearchCount int
	MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	MustScan(t, "counting note_fts", db.QueryRow("SELECT count(*) FROM note_fts"), &noteFtsCount)
	MustScan(t, "counting search results", db.QueryRow("SELECT count(*) FROM note_fts WHERE note_fts MATCH ?", "foo"), &noteSearchCount)

	assert.Equal(t, noteCount, 1, "noteCount mismatch")
	assert.Equal(t, noteFtsCount, 1, "noteFtsCount mismatch")
	assert.Equal(t, noteSearchCount, 1, "noteSearchCount mismatch")

	// execute - update
	tx, err = db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	n.Body = "baz quz"
	if err := n.Update(tx); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "updating").Error())
	}

	tx.Commit()

	// test
	MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	MustScan(t, "counting note_fts", db.QueryRow("SELECT count(*) FROM note_fts"), &noteFtsCount)
	assert.Equal(t, noteCount, 1, "noteCount mismatch")
	assert.Equal(t, noteFtsCount, 1, "noteFtsCount mismatch")

	MustScan(t, "counting search results", db.QueryRow("SELECT count(*) FROM note_fts WHERE note_fts MATCH ?", "foo"), &noteSearchCount)
	assert.Equal(t, noteSearchCount, 0, "noteSearchCount for foo mismatch")
	MustScan(t, "counting search results", db.QueryRow("SELECT count(*) FROM note_fts WHERE note_fts MATCH ?", "baz"), &noteSearchCount)
	assert.Equal(t, noteSearchCount, 1, "noteSearchCount for baz mismatch")

	// execute - delete
	tx, err = db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if err := n.Expunge(tx); err != nil {
		tx.Rollback()
		t.Fatalf(errors.Wrap(err, "expunging").Error())
	}

	tx.Commit()

	// test
	MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	MustScan(t, "counting note_fts", db.QueryRow("SELECT count(*) FROM note_fts"), &noteFtsCount)

	assert.Equal(t, noteCount, 0, "noteCount mismatch")
	assert.Equal(t, noteFtsCount, 0, "noteFtsCount mismatch")
}
