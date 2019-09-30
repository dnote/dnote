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

package testutils

import (
	"testing"

	"github.com/dnote/dnote/pkg/cli/database"
)

// Setup1 sets up a dnote env #1
func Setup1(t *testing.T, db *database.DB) {
	b1UUID := "js-book-uuid"
	b2UUID := "linux-book-uuid"

	database.MustExec(t, "setting up book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "js")
	database.MustExec(t, "setting up book 2", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "linux")

	database.MustExec(t, "setting up note 1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on) VALUES (?, ?, ?, ?)", "43827b9a-c2b0-4c06-a290-97991c896653", b1UUID, "Booleans have toString()", 1515199943)
}

// Setup2 sets up a dnote env #2
func Setup2(t *testing.T, db *database.DB) {
	b1UUID := "js-book-uuid"
	b2UUID := "linux-book-uuid"

	database.MustExec(t, "setting up book 1", db, "INSERT INTO books (uuid, label, usn) VALUES (?, ?, ?)", b1UUID, "js", 111)
	database.MustExec(t, "setting up book 2", db, "INSERT INTO books (uuid, label, usn) VALUES (?, ?, ?)", b2UUID, "linux", 122)

	database.MustExec(t, "setting up note 1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, usn) VALUES (?, ?, ?, ?, ?)", "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", b1UUID, "n1 body", 1515199951, 11)
	database.MustExec(t, "setting up note 2", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, usn) VALUES (?, ?, ?, ?, ?)", "43827b9a-c2b0-4c06-a290-97991c896653", b1UUID, "n2 body", 1515199943, 12)
	database.MustExec(t, "setting up note 3", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, usn) VALUES (?, ?, ?, ?, ?)", "3e065d55-6d47-42f2-a6bf-f5844130b2d2", b2UUID, "n3 body", 1515199961, 13)
}

// Setup3 sets up a dnote env #3
func Setup3(t *testing.T, db *database.DB) {
	b1UUID := "js-book-uuid"

	database.MustExec(t, "setting up book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "js")

	database.MustExec(t, "setting up note 1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on) VALUES (?, ?, ?, ?)", "43827b9a-c2b0-4c06-a290-97991c896653", b1UUID, "Booleans have toString()", 1515199943)
}

// Setup4 sets up a dnote env #4
func Setup4(t *testing.T, db *database.DB) {
	b1UUID := "js-book-uuid"

	database.MustExec(t, "setting up book 1", db, "INSERT INTO books (uuid, label, usn) VALUES (?, ?, ?)", b1UUID, "js", 111)

	database.MustExec(t, "setting up note 1", db, "INSERT INTO notes (rowid, uuid, book_uuid, body, added_on, usn) VALUES (?, ?, ?, ?, ?, ?)", 1, "43827b9a-c2b0-4c06-a290-97991c896653", b1UUID, "Booleans have toString()", 1515199943, 11)
	database.MustExec(t, "setting up note 2", db, "INSERT INTO notes (rowid, uuid, book_uuid, body, added_on, usn) VALUES (?, ?, ?, ?, ?, ?)", 2, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", b1UUID, "Date object implements mathematical comparisons", 1515199951, 12)
}

// Setup5 sets up a dnote env #2
func Setup5(t *testing.T, db *database.DB) {
	b1UUID := "js-book-uuid"
	b2UUID := "linux-book-uuid"

	database.MustExec(t, "setting up book 1", db, "INSERT INTO books (uuid, label, usn) VALUES (?, ?, ?)", b1UUID, "js", 111)
	database.MustExec(t, "setting up book 2", db, "INSERT INTO books (uuid, label, usn) VALUES (?, ?, ?)", b2UUID, "linux", 122)

	database.MustExec(t, "setting up note 1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, usn) VALUES (?, ?, ?, ?, ?)", "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", b1UUID, "n1 body", 1515199951, 11)
	database.MustExec(t, "setting up note 2", db, "INSERT INTO notes (uuid, book_uuid, body, added_on, usn) VALUES (?, ?, ?, ?, ?)", "43827b9a-c2b0-4c06-a290-97991c896653", b1UUID, "n2 body", 1515199943, 12)
}
