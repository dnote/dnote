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
	"os"
	"path/filepath"
	"testing"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/pkg/errors"
)

// createTestDirectories creates test directories for the given paths
func createTestDirectories(t *testing.T, paths Paths) {
	if paths.Config != "" {
		configDir := filepath.Join(paths.Config, consts.DnoteDirName)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatal(errors.Wrap(err, "creating test config directory"))
		}
	}
	if paths.Data != "" {
		dataDir := filepath.Join(paths.Data, consts.DnoteDirName)
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			t.Fatal(errors.Wrap(err, "creating test data directory"))
		}
	}
	if paths.Cache != "" {
		cacheDir := filepath.Join(paths.Cache, consts.DnoteDirName)
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			t.Fatal(errors.Wrap(err, "creating test cache directory"))
		}
	}
}

// InitTestCtx initializes a test context with an in-memory database
func InitTestCtx(t *testing.T, paths Paths) DnoteCtx {
	db := database.InitTestMemoryDB(t)
	createTestDirectories(t, paths)

	return DnoteCtx{
		DB:    db,
		Paths: paths,
		Clock: clock.NewMock(), // Use a mock clock to test times
	}
}

// InitTestCtxWithDB initializes a test context with the provided database.
// Used when you need full control over database initialization (e.g. migration tests).
func InitTestCtxWithDB(t *testing.T, paths Paths, db *database.DB) DnoteCtx {
	createTestDirectories(t, paths)

	return DnoteCtx{
		DB:    db,
		Paths: paths,
		Clock: clock.NewMock(), // Use a mock clock to test times
	}
}

// InitTestCtxWithFileDB initializes a test context with a file-based database
// at the expected XDG path. This is used for e2e tests that spawn CLI processes
// which need to access the database file.
func InitTestCtxWithFileDB(t *testing.T, paths Paths) DnoteCtx {
	createTestDirectories(t, paths)

	dbPath := filepath.Join(paths.Data, consts.DnoteDirName, consts.DnoteDBFileName)
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "opening database"))
	}

	if _, err := db.Exec(database.GetDefaultSchemaSQL()); err != nil {
		t.Fatal(errors.Wrap(err, "running schema sql"))
	}

	database.MarkMigrationComplete(t, db)

	return DnoteCtx{
		DB:    db,
		Paths: paths,
		Clock: clock.NewMock(), // Use a mock clock to test times
	}
}

// TeardownTestCtx cleans up the test context
func TeardownTestCtx(t *testing.T, ctx DnoteCtx) {
	database.TeardownTestDB(t, ctx.DB)

	if err := os.RemoveAll(ctx.Paths.Data); err != nil {
		t.Fatal(errors.Wrap(err, "removing test data directory"))
	}
	if err := os.RemoveAll(ctx.Paths.Config); err != nil {
		t.Fatal(errors.Wrap(err, "removing test config directory"))
	}
	if err := os.RemoveAll(ctx.Paths.Cache); err != nil {
		t.Fatal(errors.Wrap(err, "removing test cache directory"))
	}
}
