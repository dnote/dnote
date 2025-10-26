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

package sync

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
	cliDatabase "github.com/dnote/dnote/pkg/cli/database"
	clitest "github.com/dnote/dnote/pkg/cli/testutils"
	"github.com/google/uuid"
)

// TestSync_EmptyFragmentPreservesLastMaxUSN verifies that last_max_usn is not reset to 0
// when sync receives an empty response from the server.
//
// Scenario: Client has orphaned note (references non-existent book). During sync:
// 1. Downloads data successfully (last_max_usn=3)
// 2. Upload fails (orphaned note -> 500 error, triggers retry stepSync)
// 3. Retry stepSync gets 0 fragments (already at latest USN)
// 4. last_max_usn should stay at 3, not reset to 0
func TestSync_EmptyFragmentPreservesLastMaxUSN(t *testing.T) {
	env := setupTestEnv(t)
	user := setupUserAndLogin(t, env)

	// Create data on server (max_usn=3)
	bookUUID := apiCreateBook(t, env, user, "javascript", "creating book via API")
	apiCreateNote(t, env, user, bookUUID, "note1 content", "creating note1 via API")
	apiCreateNote(t, env, user, bookUUID, "note2 content", "creating note2 via API")

	// Create orphaned note locally (will fail to upload)
	orphanedNote := cliDatabase.Note{
		UUID:     uuid.New().String(),
		BookUUID: uuid.New().String(), // non-existent book
		Body:     "orphaned note content",
		AddedOn:  1234567890,
		EditedOn: 0,
		USN:      0,
		Deleted:  false,
		Dirty:    true,
	}
	if err := orphanedNote.Insert(env.DB); err != nil {
		t.Fatal(err)
	}

	// Run sync (downloads data, upload fails, retry gets 0 fragments)
	clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync")

	// Verify last_max_usn is preserved at 3, not reset to 0
	var lastMaxUSN int
	cliDatabase.MustScan(t, "finding system last_max_usn",
		env.DB.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastMaxUSN),
		&lastMaxUSN)

	assert.Equal(t, lastMaxUSN, 3, "last_max_usn should be 3 after syncing")
}
