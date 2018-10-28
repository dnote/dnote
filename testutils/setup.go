package testutils

import (
	"github.com/dnote/cli/infra"
	"testing"
)

// Setup1 sets up a dnote env #1
// dnote4.json
func Setup1(t *testing.T, ctx infra.DnoteCtx) {
	db := ctx.DB

	b1UUID := "js-book-uuid"
	b2UUID := "linux-book-uuid"

	MustExec(t, "setting up book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "js")
	MustExec(t, "setting up book 2", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "linux")

	MustExec(t, "setting up note 1", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", "43827b9a-c2b0-4c06-a290-97991c896653", b1UUID, "Booleans have toString()", 1515199943)
}

// Setup2 sets up a dnote env #2
// dnote3.json
func Setup2(t *testing.T, ctx infra.DnoteCtx) {
	db := ctx.DB

	b1UUID := "js-book-uuid"
	b2UUID := "linux-book-uuid"

	MustExec(t, "setting up book 1", db, "INSERT INTO books (uuid, label, usn) VALUES (?, ?, ?)", b1UUID, "js", 111)
	MustExec(t, "setting up book 2", db, "INSERT INTO books (uuid, label, usn) VALUES (?, ?, ?)", b2UUID, "linux", 122)

	MustExec(t, "setting up note 1", db, "INSERT INTO notes (id, uuid, book_uuid, content, added_on, usn) VALUES (?, ?, ?, ?, ?, ?)", 1, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", b1UUID, "n1 content", 1515199951, 11)
	MustExec(t, "setting up note 2", db, "INSERT INTO notes (id, uuid, book_uuid, content, added_on, usn) VALUES (?, ?, ?, ?, ?, ?)", 2, "43827b9a-c2b0-4c06-a290-97991c896653", b1UUID, "n2 content", 1515199943, 12)
	MustExec(t, "setting up note 3", db, "INSERT INTO notes (id, uuid, book_uuid, content, added_on, usn) VALUES (?, ?, ?, ?, ?, ?)", 3, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", b2UUID, "n3 content", 1515199961, 13)
}

// Setup3 sets up a dnote env #1
// dnote1.json
func Setup3(t *testing.T, ctx infra.DnoteCtx) {
	db := ctx.DB

	b1UUID := "js-book-uuid"

	MustExec(t, "setting up book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "js")

	MustExec(t, "setting up note 1", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", "43827b9a-c2b0-4c06-a290-97991c896653", b1UUID, "Booleans have toString()", 1515199943)
}

// Setup4 sets up a dnote env #1
// dnote2.json
func Setup4(t *testing.T, ctx infra.DnoteCtx) {
	db := ctx.DB

	b1UUID := "js-book-uuid"

	MustExec(t, "setting up book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "js")

	MustExec(t, "setting up note 1", db, "INSERT INTO notes (id, uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?, ?)", 1, "43827b9a-c2b0-4c06-a290-97991c896653", b1UUID, "Booleans have toString()", 1515199943)
	MustExec(t, "setting up note 2", db, "INSERT INTO notes (id, uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?, ?)", 2, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", b1UUID, "Date object implements mathematical comparisons", 1515199951)
}
