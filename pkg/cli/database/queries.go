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

	"github.com/dnote/dnote/pkg/clock"
	"github.com/pkg/errors"
)

// GetSystem scans the given system configuration record onto the destination
func GetSystem(db *DB, key string, dest interface{}) error {
	if err := db.QueryRow("SELECT value FROM system WHERE key = ?", key).Scan(dest); err != nil {
		return errors.Wrap(err, "finding system configuration record")
	}

	return nil
}

// InsertSystem inserets a system configuration
func InsertSystem(db *DB, key, val string) error {
	if _, err := db.Exec("INSERT INTO system (key, value) VALUES (? , ?);", key, val); err != nil {
		return errors.Wrap(err, "saving system config")
	}

	return nil
}

// UpsertSystem inserts or updates a system configuration
func UpsertSystem(db *DB, key, val string) error {
	var count int
	if err := db.QueryRow("SELECT count(*) FROM system WHERE key = ?", key).Scan(&count); err != nil {
		return errors.Wrap(err, "counting system record")
	}

	if count == 0 {
		if _, err := db.Exec("INSERT INTO system (key, value) VALUES (? , ?);", key, val); err != nil {
			return errors.Wrap(err, "saving system config")
		}
	} else {
		if _, err := db.Exec("UPDATE system SET value = ? WHERE key = ?", val, key); err != nil {
			return errors.Wrap(err, "updating system config")
		}
	}

	return nil
}

// UpdateSystem updates a system configuration
func UpdateSystem(db *DB, key, val interface{}) error {
	if _, err := db.Exec("UPDATE system SET value = ? WHERE key = ?", val, key); err != nil {
		return errors.Wrap(err, "updating system config")
	}

	return nil
}

// DeleteSystem delets the given system record
func DeleteSystem(db *DB, key string) error {
	if _, err := db.Exec("DELETE FROM system WHERE key = ?", key); err != nil {
		return errors.Wrap(err, "deleting system config")
	}

	return nil
}

// NoteInfo is a basic information about a note
type NoteInfo struct {
	RowID     int
	BookLabel string
	UUID      string
	Content   string
	AddedOn   int64
	EditedOn  int64
}

// GetNoteInfo returns a NoteInfo for the note with the given noteRowID
func GetNoteInfo(db *DB, noteRowID int) (NoteInfo, error) {
	var ret NoteInfo

	err := db.QueryRow(`SELECT books.label, notes.uuid, notes.body, notes.added_on, notes.edited_on, notes.rowid
			FROM notes
			INNER JOIN books ON books.uuid = notes.book_uuid
			WHERE notes.rowid = ? AND notes.deleted = false`, noteRowID).
		Scan(&ret.BookLabel, &ret.UUID, &ret.Content, &ret.AddedOn, &ret.EditedOn, &ret.RowID)
	if err == sql.ErrNoRows {
		return ret, errors.Errorf("note %d not found", noteRowID)
	} else if err != nil {
		return ret, errors.Wrap(err, "querying the note")
	}

	return ret, nil
}

// BookInfo is a basic information about a book
type BookInfo struct {
	RowID int
	UUID  string
	Name  string
}

// GetBookInfo returns a BookInfo for the book with the given uuid
func GetBookInfo(db *DB, uuid string) (BookInfo, error) {
	var ret BookInfo

	err := db.QueryRow(`SELECT books.rowid, books.uuid, books.label
			FROM books
			WHERE books.uuid = ? AND books.deleted = false`, uuid).
		Scan(&ret.RowID, &ret.UUID, &ret.Name)
	if err == sql.ErrNoRows {
		return ret, errors.Errorf("book %s not found", uuid)
	} else if err != nil {
		return ret, errors.Wrap(err, "querying the note")
	}

	return ret, nil
}

// GetBookUUID returns a uuid of a book given a label
func GetBookUUID(db *DB, label string) (string, error) {
	var ret string
	err := db.QueryRow("SELECT uuid FROM books WHERE label = ?", label).Scan(&ret)
	if err == sql.ErrNoRows {
		return ret, errors.Errorf("book '%s' not found", label)
	} else if err != nil {
		return ret, errors.Wrap(err, "querying the book")
	}

	return ret, nil
}

// UpdateBookName updates a book name
func UpdateBookName(db *DB, uuid string, name string) error {
	_, err := db.Exec(`UPDATE books
		SET label = ?, dirty = ?
		WHERE uuid = ?`, name, true, uuid)
	if err != nil {
		return errors.Wrap(err, "updating the book")
	}

	return nil
}

// GetActiveNote gets the note which has the given rowid and is not deleted
func GetActiveNote(db *DB, rowid int) (Note, error) {
	var ret Note

	err := db.QueryRow(`SELECT
		rowid,
		uuid,
		book_uuid,
		body,
		added_on,
		edited_on,
		usn,
		public,
		deleted,
		dirty
	FROM notes WHERE rowid = ? AND deleted = false;`, rowid).Scan(
		&ret.RowID,
		&ret.UUID,
		&ret.BookUUID,
		&ret.Body,
		&ret.AddedOn,
		&ret.EditedOn,
		&ret.USN,
		&ret.Public,
		&ret.Deleted,
		&ret.Dirty,
	)

	if err == sql.ErrNoRows {
		return ret, err
	} else if err != nil {
		return ret, errors.Wrap(err, "finding the note")
	}

	return ret, nil
}

// UpdateNoteContent updates the note content and marks the note as dirty
func UpdateNoteContent(db *DB, c clock.Clock, rowID int, content string) error {
	ts := c.Now().UnixNano()

	_, err := db.Exec(`UPDATE notes
			SET body = ?, edited_on = ?, dirty = ?
			WHERE rowid = ?`, content, ts, true, rowID)
	if err != nil {
		return errors.Wrap(err, "updating the note")
	}

	return nil
}

// UpdateNoteBook moves the note to a different book and marks the note as dirty
func UpdateNoteBook(db *DB, c clock.Clock, rowID int, bookUUID string) error {
	ts := c.Now().UnixNano()

	_, err := db.Exec(`UPDATE notes
			SET book_uuid = ?, edited_on = ?, dirty = ?
			WHERE rowid = ?`, bookUUID, ts, true, rowID)
	if err != nil {
		return errors.Wrap(err, "updating the note")
	}

	return nil
}
