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

package context

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/clock"
)

// InitTestCtx initializes a test context
func InitTestCtx(t *testing.T, dnoteDir string, dbOpts *database.TestDBOptions) DnoteCtx {
	dbPath := fmt.Sprintf("%s/%s", dnoteDir, consts.DnoteDBFileName)

	db := database.InitTestDB(t, dbPath, dbOpts)

	return DnoteCtx{
		DB:       db,
		DnoteDir: dnoteDir,
		// Use a mock clock to test times
		Clock: clock.NewMock(),
	}
}

// TeardownTestCtx cleans up the test context
func TeardownTestCtx(t *testing.T, ctx DnoteCtx) {
	database.CloseTestDB(t, ctx.DB)
}
