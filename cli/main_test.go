/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/dnote/dnote/cli/core"
	"github.com/dnote/dnote/cli/infra"
	"github.com/dnote/dnote/cli/testutils"
	"github.com/dnote/dnote/cli/utils"
	"github.com/pkg/errors"
)

var binaryName = "test-dnote"

func TestMain(m *testing.M) {
	if err := exec.Command("go", "build", "--tags", "fts5", "-o", binaryName).Run(); err != nil {
		log.Print(errors.Wrap(err, "building a binary").Error())
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestInit(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv(t, "./tmp", "./testutils/fixtures/schema.sql", true)
	defer testutils.TeardownEnv(ctx)

	// Execute
	testutils.RunDnoteCmd(t, ctx, binaryName)

	// Test
	if !utils.FileExists(fmt.Sprintf("%s", ctx.DnoteDir)) {
		t.Errorf("dnote directory was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, core.ConfigFilename)) {
		t.Errorf("config file was not initialized")
	}

	db := ctx.DB

	var notesTableCount, booksTableCount, systemTableCount int
	testutils.MustScan(t, "counting notes",
		db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = ? AND name = ?", "table", "notes"), &notesTableCount)
	testutils.MustScan(t, "counting books",
		db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = ? AND name = ?", "table", "books"), &booksTableCount)
	testutils.MustScan(t, "counting system",
		db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = ? AND name = ?", "table", "system"), &systemTableCount)

	testutils.AssertEqual(t, notesTableCount, 1, "notes table count mismatch")
	testutils.AssertEqual(t, booksTableCount, 1, "books table count mismatch")
	testutils.AssertEqual(t, systemTableCount, 1, "system table count mismatch")

	// test that all default system configurations are generated
	var lastUpgrade, lastMaxUSN, lastSyncAt string
	testutils.MustScan(t, "scanning last upgrade",
		db.QueryRow("SELECT value FROM system WHERE key = ?", infra.SystemLastUpgrade), &lastUpgrade)
	testutils.MustScan(t, "scanning last max usn",
		db.QueryRow("SELECT value FROM system WHERE key = ?", infra.SystemLastMaxUSN), &lastMaxUSN)
	testutils.MustScan(t, "scanning last sync at",
		db.QueryRow("SELECT value FROM system WHERE key = ?", infra.SystemLastSyncAt), &lastSyncAt)

	testutils.AssertNotEqual(t, lastUpgrade, "", "last upgrade should not be empty")
	testutils.AssertNotEqual(t, lastMaxUSN, "", "last max usn should not be empty")
	testutils.AssertNotEqual(t, lastSyncAt, "", "last sync at should not be empty")
}

func TestAddNote_NewBook_BodyFlag(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv(t, "./tmp", "./testutils/fixtures/schema.sql", true)
	defer testutils.TeardownEnv(ctx)

	// Execute
	testutils.RunDnoteCmd(t, ctx, binaryName, "add", "js", "-c", "foo")

	// Test
	db := ctx.DB

	var noteCount, bookCount int
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)

	testutils.AssertEqualf(t, bookCount, 1, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 1, "note count mismatch")

	var book core.Book
	testutils.MustScan(t, "getting book", db.QueryRow("SELECT uuid, dirty FROM books where label = ?", "js"), &book.UUID, &book.Dirty)
	var note core.Note
	testutils.MustScan(t, "getting note",
		db.QueryRow("SELECT uuid, body, added_on, dirty FROM notes where book_uuid = ?", book.UUID), &note.UUID, &note.Body, &note.AddedOn, &note.Dirty)

	testutils.AssertEqual(t, book.Dirty, true, "Book dirty mismatch")

	testutils.AssertNotEqual(t, note.UUID, "", "Note should have UUID")
	testutils.AssertEqual(t, note.Body, "foo", "Note body mismatch")
	testutils.AssertEqual(t, note.Dirty, true, "Note dirty mismatch")
	testutils.AssertNotEqual(t, note.AddedOn, int64(0), "Note added_on mismatch")
}

func TestAddNote_ExistingBook_BodyFlag(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv(t, "./tmp", "./testutils/fixtures/schema.sql", true)
	defer testutils.TeardownEnv(ctx)

	testutils.Setup3(t, ctx)

	// Execute
	testutils.RunDnoteCmd(t, ctx, binaryName, "add", "js", "-c", "foo")

	// Test
	db := ctx.DB

	var noteCount, bookCount int
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)

	testutils.AssertEqualf(t, bookCount, 1, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 2, "note count mismatch")

	var n1, n2 core.Note
	testutils.MustScan(t, "getting n1",
		db.QueryRow("SELECT uuid, body, added_on, dirty FROM notes WHERE book_uuid = ? AND uuid = ?", "js-book-uuid", "43827b9a-c2b0-4c06-a290-97991c896653"), &n1.UUID, &n1.Body, &n1.AddedOn, &n1.Dirty)
	testutils.MustScan(t, "getting n2",
		db.QueryRow("SELECT uuid, body, added_on, dirty FROM notes WHERE book_uuid = ? AND body = ?", "js-book-uuid", "foo"), &n2.UUID, &n2.Body, &n2.AddedOn, &n2.Dirty)

	var book core.Book
	testutils.MustScan(t, "getting book", db.QueryRow("SELECT dirty FROM books where label = ?", "js"), &book.Dirty)

	testutils.AssertEqual(t, book.Dirty, false, "Book dirty mismatch")

	testutils.AssertNotEqual(t, n1.UUID, "", "n1 should have UUID")
	testutils.AssertEqual(t, n1.Body, "Booleans have toString()", "n1 body mismatch")
	testutils.AssertEqual(t, n1.AddedOn, int64(1515199943), "n1 added_on mismatch")
	testutils.AssertEqual(t, n1.Dirty, false, "n1 dirty mismatch")

	testutils.AssertNotEqual(t, n2.UUID, "", "n2 should have UUID")
	testutils.AssertEqual(t, n2.Body, "foo", "n2 body mismatch")
	testutils.AssertEqual(t, n2.Dirty, true, "n2 dirty mismatch")
}

func TestEditNote_BodyFlag(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv(t, "./tmp", "./testutils/fixtures/schema.sql", true)
	defer testutils.TeardownEnv(ctx)

	testutils.Setup4(t, ctx)

	// Execute
	testutils.RunDnoteCmd(t, ctx, binaryName, "edit", "js", "2", "-c", "foo bar")

	// Test
	db := ctx.DB

	var noteCount, bookCount int
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)

	testutils.AssertEqualf(t, bookCount, 1, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 2, "note count mismatch")

	var n1, n2 core.Note
	testutils.MustScan(t, "getting n1",
		db.QueryRow("SELECT uuid, body, added_on, dirty FROM notes where book_uuid = ? AND uuid = ?", "js-book-uuid", "43827b9a-c2b0-4c06-a290-97991c896653"), &n1.UUID, &n1.Body, &n1.AddedOn, &n1.Dirty)
	testutils.MustScan(t, "getting n2",
		db.QueryRow("SELECT uuid, body, added_on, dirty FROM notes where book_uuid = ? AND uuid = ?", "js-book-uuid", "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f"), &n2.UUID, &n2.Body, &n2.AddedOn, &n2.Dirty)

	testutils.AssertEqual(t, n1.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "n1 should have UUID")
	testutils.AssertEqual(t, n1.Body, "Booleans have toString()", "n1 body mismatch")
	testutils.AssertEqual(t, n1.Dirty, false, "n1 dirty mismatch")

	testutils.AssertEqual(t, n2.UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "Note should have UUID")
	testutils.AssertEqual(t, n2.Body, "foo bar", "Note body mismatch")
	testutils.AssertEqual(t, n2.Dirty, true, "n2 dirty mismatch")
	testutils.AssertNotEqual(t, n2.EditedOn, 0, "Note edited_on mismatch")
}

func TestRemoveNote(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv(t, "./tmp", "./testutils/fixtures/schema.sql", true)
	defer testutils.TeardownEnv(ctx)

	testutils.Setup2(t, ctx)

	// Execute
	testutils.WaitDnoteCmd(t, ctx, testutils.UserConfirm, binaryName, "remove", "js", "1")

	// Test
	db := ctx.DB

	var noteCount, bookCount, jsNoteCount, linuxNoteCount int
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	testutils.MustScan(t, "counting js notes", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "js-book-uuid"), &jsNoteCount)
	testutils.MustScan(t, "counting linux notes", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "linux-book-uuid"), &linuxNoteCount)

	testutils.AssertEqualf(t, bookCount, 2, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 3, "note count mismatch")
	testutils.AssertEqual(t, jsNoteCount, 2, "js book should have 2 notes")
	testutils.AssertEqual(t, linuxNoteCount, 1, "linux book book should have 1 note")

	var b1, b2 core.Book
	var n1, n2, n3 core.Note
	testutils.MustScan(t, "getting b1",
		db.QueryRow("SELECT label, deleted, usn FROM books WHERE uuid = ?", "js-book-uuid"),
		&b1.Label, &b1.Deleted, &b1.USN)
	testutils.MustScan(t, "getting b2",
		db.QueryRow("SELECT label, deleted, usn FROM books WHERE uuid = ?", "linux-book-uuid"),
		&b2.Label, &b2.Deleted, &b2.USN)
	testutils.MustScan(t, "getting n1",
		db.QueryRow("SELECT uuid, body, added_on, deleted, dirty, usn FROM notes WHERE book_uuid = ? AND uuid = ?", "js-book-uuid", "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f"),
		&n1.UUID, &n1.Body, &n1.AddedOn, &n1.Deleted, &n1.Dirty, &n1.USN)
	testutils.MustScan(t, "getting n2",
		db.QueryRow("SELECT uuid, body, added_on, deleted, dirty, usn FROM notes WHERE book_uuid = ? AND uuid = ?", "js-book-uuid", "43827b9a-c2b0-4c06-a290-97991c896653"),
		&n2.UUID, &n2.Body, &n2.AddedOn, &n2.Deleted, &n2.Dirty, &n2.USN)
	testutils.MustScan(t, "getting n3",
		db.QueryRow("SELECT uuid, body, added_on, deleted, dirty, usn FROM notes WHERE book_uuid = ? AND uuid = ?", "linux-book-uuid", "3e065d55-6d47-42f2-a6bf-f5844130b2d2"),
		&n3.UUID, &n3.Body, &n3.AddedOn, &n3.Deleted, &n3.Dirty, &n3.USN)

	testutils.AssertEqual(t, b1.Label, "js", "b1 label mismatch")
	testutils.AssertEqual(t, b1.Deleted, false, "b1 deleted mismatch")
	testutils.AssertEqual(t, b1.Dirty, false, "b1 Dirty mismatch")
	testutils.AssertEqual(t, b1.USN, 111, "b1 usn mismatch")

	testutils.AssertEqual(t, b2.Label, "linux", "b2 label mismatch")
	testutils.AssertEqual(t, b2.Deleted, false, "b2 deleted mismatch")
	testutils.AssertEqual(t, b2.Dirty, false, "b2 Dirty mismatch")
	testutils.AssertEqual(t, b2.USN, 122, "b2 usn mismatch")

	testutils.AssertEqual(t, n1.UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "n1 should have UUID")
	testutils.AssertEqual(t, n1.Body, "", "n1 body mismatch")
	testutils.AssertEqual(t, n1.Deleted, true, "n1 deleted mismatch")
	testutils.AssertEqual(t, n1.Dirty, true, "n1 Dirty mismatch")
	testutils.AssertEqual(t, n1.USN, 11, "n1 usn mismatch")

	testutils.AssertEqual(t, n2.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "n2 should have UUID")
	testutils.AssertEqual(t, n2.Body, "n2 body", "n2 body mismatch")
	testutils.AssertEqual(t, n2.Deleted, false, "n2 deleted mismatch")
	testutils.AssertEqual(t, n2.Dirty, false, "n2 Dirty mismatch")
	testutils.AssertEqual(t, n2.USN, 12, "n2 usn mismatch")

	testutils.AssertEqual(t, n3.UUID, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", "n3 should have UUID")
	testutils.AssertEqual(t, n3.Body, "n3 body", "n3 body mismatch")
	testutils.AssertEqual(t, n3.Deleted, false, "n3 deleted mismatch")
	testutils.AssertEqual(t, n3.Dirty, false, "n3 Dirty mismatch")
	testutils.AssertEqual(t, n3.USN, 13, "n3 usn mismatch")
}

func TestRemoveBook(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv(t, "./tmp", "./testutils/fixtures/schema.sql", true)
	defer testutils.TeardownEnv(ctx)

	testutils.Setup2(t, ctx)

	// Execute
	testutils.WaitDnoteCmd(t, ctx, testutils.UserConfirm, binaryName, "remove", "-b", "js")

	// Test
	db := ctx.DB

	var noteCount, bookCount, jsNoteCount, linuxNoteCount int
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	testutils.MustScan(t, "counting js notes", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "js-book-uuid"), &jsNoteCount)
	testutils.MustScan(t, "counting linux notes", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "linux-book-uuid"), &linuxNoteCount)

	testutils.AssertEqualf(t, bookCount, 2, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 3, "note count mismatch")
	testutils.AssertEqual(t, jsNoteCount, 2, "js book should have 2 notes")
	testutils.AssertEqual(t, linuxNoteCount, 1, "linux book book should have 1 note")

	var b1, b2 core.Book
	var n1, n2, n3 core.Note
	testutils.MustScan(t, "getting b1",
		db.QueryRow("SELECT label, dirty, deleted, usn FROM books WHERE uuid = ?", "js-book-uuid"),
		&b1.Label, &b1.Dirty, &b1.Deleted, &b1.USN)
	testutils.MustScan(t, "getting b2",
		db.QueryRow("SELECT label, dirty, deleted, usn FROM books WHERE uuid = ?", "linux-book-uuid"),
		&b2.Label, &b2.Dirty, &b2.Deleted, &b2.USN)
	testutils.MustScan(t, "getting n1",
		db.QueryRow("SELECT uuid, body, added_on, dirty, deleted, usn FROM notes WHERE book_uuid = ? AND uuid = ?", "js-book-uuid", "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f"),
		&n1.UUID, &n1.Body, &n1.AddedOn, &n1.Deleted, &n1.Dirty, &n1.USN)
	testutils.MustScan(t, "getting n2",
		db.QueryRow("SELECT uuid, body, added_on, dirty, deleted, usn FROM notes WHERE book_uuid = ? AND uuid = ?", "js-book-uuid", "43827b9a-c2b0-4c06-a290-97991c896653"),
		&n2.UUID, &n2.Body, &n2.AddedOn, &n2.Deleted, &n2.Dirty, &n2.USN)
	testutils.MustScan(t, "getting n3",
		db.QueryRow("SELECT uuid, body, added_on, dirty, deleted, usn FROM notes WHERE book_uuid = ? AND uuid = ?", "linux-book-uuid", "3e065d55-6d47-42f2-a6bf-f5844130b2d2"),
		&n3.UUID, &n3.Body, &n3.AddedOn, &n3.Deleted, &n3.Dirty, &n3.USN)

	testutils.AssertNotEqual(t, b1.Label, "js", "b1 label mismatch")
	testutils.AssertEqual(t, b1.Dirty, true, "b1 Dirty mismatch")
	testutils.AssertEqual(t, b1.Deleted, true, "b1 deleted mismatch")
	testutils.AssertEqual(t, b1.USN, 111, "b1 usn mismatch")

	testutils.AssertEqual(t, b2.Label, "linux", "b2 label mismatch")
	testutils.AssertEqual(t, b2.Dirty, false, "b2 Dirty mismatch")
	testutils.AssertEqual(t, b2.Deleted, false, "b2 deleted mismatch")
	testutils.AssertEqual(t, b2.USN, 122, "b2 usn mismatch")

	testutils.AssertEqual(t, n1.UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "n1 should have UUID")
	testutils.AssertEqual(t, n1.Body, "", "n1 body mismatch")
	testutils.AssertEqual(t, n1.Dirty, true, "n1 Dirty mismatch")
	testutils.AssertEqual(t, n1.Deleted, true, "n1 deleted mismatch")
	testutils.AssertEqual(t, n1.USN, 11, "n1 usn mismatch")

	testutils.AssertEqual(t, n2.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "n2 should have UUID")
	testutils.AssertEqual(t, n2.Body, "", "n2 body mismatch")
	testutils.AssertEqual(t, n2.Dirty, true, "n2 Dirty mismatch")
	testutils.AssertEqual(t, n2.Deleted, true, "n2 deleted mismatch")
	testutils.AssertEqual(t, n2.USN, 12, "n2 usn mismatch")

	testutils.AssertEqual(t, n3.UUID, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", "n3 should have UUID")
	testutils.AssertEqual(t, n3.Body, "n3 body", "n3 body mismatch")
	testutils.AssertEqual(t, n3.Dirty, false, "n3 Dirty mismatch")
	testutils.AssertEqual(t, n3.Deleted, false, "n3 deleted mismatch")
	testutils.AssertEqual(t, n3.USN, 13, "n3 usn mismatch")
}
