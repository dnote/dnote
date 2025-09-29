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
	"fmt"
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/pkg/errors"
)

// InitTestCtx initializes a test context
func InitTestCtx(t *testing.T, paths Paths, dbOpts *database.TestDBOptions) DnoteCtx {
	dbPath := fmt.Sprintf("%s/%s/%s", paths.Data, consts.DnoteDirName, consts.DnoteDBFileName)

	db := database.InitTestDB(t, dbPath, dbOpts)

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
