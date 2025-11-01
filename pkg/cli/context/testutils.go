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

package context

import (
	"path/filepath"
	"testing"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/pkg/errors"
)

// getDefaultTestPaths creates default test paths with all paths pointing to a temp directory
func getDefaultTestPaths(t *testing.T) Paths {
	tmpDir := t.TempDir()
	return Paths{
		Home:   tmpDir,
		Cache:  tmpDir,
		Config: tmpDir,
		Data:   tmpDir,
	}
}


// InitTestCtx initializes a test context with an in-memory database
// and a temporary directory for all paths
func InitTestCtx(t *testing.T) DnoteCtx {
	paths := getDefaultTestPaths(t)
	db := database.InitTestMemoryDB(t)

	if err := InitDnoteDirs(paths); err != nil {
		t.Fatal(errors.Wrap(err, "creating test directories"))
	}

	return DnoteCtx{
		DB:    db,
		Paths: paths,
		Clock: clock.NewMock(), // Use a mock clock to test times
	}
}

// InitTestCtxWithDB initializes a test context with the provided database
// and a temporary directory for all paths.
// Used when you need full control over database initialization (e.g. migration tests).
func InitTestCtxWithDB(t *testing.T, db *database.DB) DnoteCtx {
	paths := getDefaultTestPaths(t)

	if err := InitDnoteDirs(paths); err != nil {
		t.Fatal(errors.Wrap(err, "creating test directories"))
	}

	return DnoteCtx{
		DB:    db,
		Paths: paths,
		Clock: clock.NewMock(), // Use a mock clock to test times
	}
}

// InitTestCtxWithFileDB initializes a test context with a file-based database
// at the expected path.
func InitTestCtxWithFileDB(t *testing.T) DnoteCtx {
	paths := getDefaultTestPaths(t)

	if err := InitDnoteDirs(paths); err != nil {
		t.Fatal(errors.Wrap(err, "creating test directories"))
	}

	dbPath := filepath.Join(paths.Data, consts.DnoteDirName, consts.DnoteDBFileName)
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "opening database"))
	}

	if _, err := db.Exec(database.GetDefaultSchemaSQL()); err != nil {
		t.Fatal(errors.Wrap(err, "running schema sql"))
	}

	t.Cleanup(func() { db.Close() })

	return DnoteCtx{
		DB:    db,
		Paths: paths,
		Clock: clock.NewMock(), // Use a mock clock to test times
	}
}
