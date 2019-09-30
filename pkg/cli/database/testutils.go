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
	"os"
	"path/filepath"
	"testing"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
)

var defaultSchemaSQL = `CREATE TABLE books
		(
			uuid text PRIMARY KEY,
			label text NOT NULL
		, dirty bool DEFAULT false, usn int DEFAULT 0 NOT NULL, deleted bool DEFAULT false);
CREATE TABLE system
		(
			key string NOT NULL,
			value text NOT NULL
		);
CREATE UNIQUE INDEX idx_books_label ON books(label);
CREATE UNIQUE INDEX idx_books_uuid ON books(uuid);
CREATE TABLE IF NOT EXISTS "notes"
		(
			uuid text NOT NULL,
			book_uuid text NOT NULL,
			body text NOT NULL,
			added_on integer NOT NULL,
			edited_on integer DEFAULT 0,
			public bool DEFAULT false,
			dirty bool DEFAULT false,
			usn int DEFAULT 0 NOT NULL,
			deleted bool DEFAULT false
		);
CREATE VIRTUAL TABLE note_fts USING fts5(content=notes, body, tokenize="porter unicode61 categories 'L* N* Co Ps Pe'")
/* note_fts(body) */;
CREATE TABLE IF NOT EXISTS 'note_fts_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'note_fts_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'note_fts_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'note_fts_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE TRIGGER notes_after_insert AFTER INSERT ON notes BEGIN
				INSERT INTO note_fts(rowid, body) VALUES (new.rowid, new.body);
			END;
CREATE TRIGGER notes_after_delete AFTER DELETE ON notes BEGIN
				INSERT INTO note_fts(note_fts, rowid, body) VALUES ('delete', old.rowid, old.body);
			END;
CREATE TRIGGER notes_after_update AFTER UPDATE ON notes BEGIN
				INSERT INTO note_fts(note_fts, rowid, body) VALUES ('delete', old.rowid, old.body);
				INSERT INTO note_fts(rowid, body) VALUES (new.rowid, new.body);
			END;
CREATE TABLE actions
		(
			uuid text PRIMARY KEY,
			schema integer NOT NULL,
			type text NOT NULL,
			data text NOT NULL,
			timestamp integer NOT NULL
		);
CREATE UNIQUE INDEX idx_notes_uuid ON notes(uuid);
CREATE INDEX idx_notes_book_uuid ON notes(book_uuid);`

// MustScan scans the given row and fails a test in case of any errors
func MustScan(t *testing.T, message string, row *sql.Row, args ...interface{}) {
	err := row.Scan(args...)
	if err != nil {
		t.Fatal(errors.Wrap(errors.Wrap(err, "scanning a row"), message))
	}
}

// MustExec executes the given SQL query and fails a test if an error occurs
func MustExec(t *testing.T, message string, db *DB, query string, args ...interface{}) sql.Result {
	result, err := db.Exec(query, args...)
	if err != nil {
		t.Fatal(errors.Wrap(errors.Wrap(err, "executing sql"), message))
	}

	return result
}

// TestDBOptions contains options for test database
type TestDBOptions struct {
	SchemaSQLPath string
	SkipMigration bool
}

// InitTestDB opens a test database connection
func InitTestDB(t *testing.T, dbPath string, options *TestDBOptions) *DB {
	db, err := Open(dbPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "opening database connection"))
	}

	dir, _ := filepath.Split(dbPath)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		t.Fatal(errors.Wrap(err, "creating the directory for test database file"))
	}

	var schemaSQL string
	if options != nil && options.SchemaSQLPath != "" {
		b := utils.ReadFileAbs(options.SchemaSQLPath)
		schemaSQL = string(b)
	} else {
		schemaSQL = defaultSchemaSQL
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatal(errors.Wrap(err, "running schema sql"))
	}

	if options == nil || !options.SkipMigration {
		MarkMigrationComplete(t, db)
	}

	return db
}

// CloseTestDB closes the test database
func CloseTestDB(t *testing.T, db *DB) {
	if err := db.Close(); err != nil {
		t.Fatal(errors.Wrap(err, "closing database"))
	}

	if err := os.RemoveAll(db.Filepath); err != nil {
		t.Fatal(errors.Wrap(err, "removing database file"))
	}
}

// OpenTestDB opens the database connection to the test database
func OpenTestDB(t *testing.T, dnoteDir string) *DB {
	dbPath := fmt.Sprintf("%s/%s", dnoteDir, consts.DnoteDBFileName)
	db, err := Open(dbPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "opening database connection to the test database"))
	}

	return db
}

// MarkMigrationComplete marks all migrations as complete in the database
func MarkMigrationComplete(t *testing.T, db *DB) {
	if _, err := db.Exec("INSERT INTO system (key, value) VALUES (? , ?);", consts.SystemSchema, 12); err != nil {
		t.Fatal(errors.Wrap(err, "inserting schema"))
	}
	if _, err := db.Exec("INSERT INTO system (key, value) VALUES (? , ?);", consts.SystemRemoteSchema, 1); err != nil {
		t.Fatal(errors.Wrap(err, "inserting remote schema"))
	}
}
