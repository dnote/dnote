/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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
