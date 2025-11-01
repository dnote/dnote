/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
)

//go:embed schema.sql
var defaultSchemaSQL string

// GetDefaultSchemaSQL returns the default schema SQL for tests
func GetDefaultSchemaSQL() string {
	return defaultSchemaSQL
}

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

// InitTestMemoryDB initializes an in-memory test database with the default schema.
func InitTestMemoryDB(t *testing.T) *DB {
	return InitTestMemoryDBRaw(t, "")
}

// InitTestFileDB initializes a file-based test database with the default schema.
func InitTestFileDB(t *testing.T) (*DB, string) {
	uuid := mustGenerateTestUUID(t)
	dbPath := filepath.Join(t.TempDir(), fmt.Sprintf("dnote-%s.db", uuid))
	db := InitTestFileDBRaw(t, dbPath)
	return db, dbPath
}

// InitTestFileDBRaw initializes a file-based test database at the specified path with the default schema.
func InitTestFileDBRaw(t *testing.T, dbPath string) *DB {
	db, err := Open(dbPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "opening database"))
	}

	if _, err := db.Exec(defaultSchemaSQL); err != nil {
		t.Fatal(errors.Wrap(err, "running schema sql"))
	}

	t.Cleanup(func() { db.Close() })
	return db
}

// InitTestMemoryDBRaw initializes an in-memory test database without marking migrations complete.
// If schemaPath is empty, uses the default schema. Used for migration testing.
func InitTestMemoryDBRaw(t *testing.T, schemaPath string) *DB {
	uuid := mustGenerateTestUUID(t)
	dbName := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid)

	db, err := Open(dbName)
	if err != nil {
		t.Fatal(errors.Wrap(err, "opening in-memory database"))
	}

	var schemaSQL string
	if schemaPath != "" {
		schemaSQL = string(utils.ReadFileAbs(schemaPath))
	} else {
		schemaSQL = defaultSchemaSQL
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatal(errors.Wrap(err, "running schema sql"))
	}

	t.Cleanup(func() { db.Close() })
	return db
}

// OpenTestDB opens the database connection to a test database
// without initializing any schema
func OpenTestDB(t *testing.T, dnoteDir string) *DB {
	dbPath := fmt.Sprintf("%s/%s/%s", dnoteDir, consts.DnoteDirName, consts.DnoteDBFileName)
	db, err := Open(dbPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "opening database connection to the test database"))
	}

	return db
}

// mustGenerateTestUUID generates a UUID for test databases and fails the test on error
func mustGenerateTestUUID(t *testing.T) string {
	uuid, err := utils.GenerateUUID()
	if err != nil {
		t.Fatal(errors.Wrap(err, "generating UUID for test database"))
	}
	return uuid
}
