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

package database

import (
	"database/sql"

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
func GetNoteInfo(db *DB, noteRowID string) (NoteInfo, error) {
	var ret NoteInfo

	err := db.QueryRow(`SELECT books.label, notes.uuid, notes.body, notes.added_on, notes.edited_on, notes.rowid
			FROM notes
			INNER JOIN books ON books.uuid = notes.book_uuid
			WHERE notes.rowid = ? AND notes.deleted = false`, noteRowID).
		Scan(&ret.BookLabel, &ret.UUID, &ret.Content, &ret.AddedOn, &ret.EditedOn, &ret.RowID)
	if err == sql.ErrNoRows {
		return ret, errors.Errorf("note %s not found", noteRowID)
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
