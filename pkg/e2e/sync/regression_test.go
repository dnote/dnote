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

// TestSync_LastMaxUSNResetBug reproduces the bug where last_max_usn gets reset to 0
// when a retry stepSync gets no fragments from the server.
//
// Bug scenario:
// 1. Client syncs and downloads data from server (last_max_usn should be updated to server's max USN)
// 2. Client has a dirty note that fails to upload (e.g., 500 error due to orphaned note)
// 3. Sync triggers retry stepSync because of the upload error
// 4. Retry stepSync queries after_usn=<server_max_usn> and gets 0 fragments
// 5. processFragments returns maxUSN=0 (initialized to 0, no fragments to process)
// 6. saveSyncState(tx, serverTime, 0) overwrites last_max_usn to 0
// 7. Transaction commits with last_max_usn=0 instead of server's max USN
func TestSync_LastMaxUSNResetBug(t *testing.T) {
	env := setupTestEnv(t)
	user := setupUserAndLogin(t, env)

	// Step 1: Create data on server so client has something to download
	bookUUID := apiCreateBook(t, env, user, "javascript", "creating book via API")
	apiCreateNote(t, env, user, bookUUID, "note1 content", "creating note1 via API")
	apiCreateNote(t, env, user, bookUUID, "note2 content", "creating note2 via API")

	// At this point, server has:
	// - 1 book (USN=1)
	// - 2 notes (USN=2,3)
	// - server max_usn=3

	// Step 2: Create a dirty note on client that references a non-existent book
	// This will cause a 500 error when trying to upload because the server
	// will try to find the book and fail
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

	// Step 3: Run sync
	// Expected flow:
	// 1. First stepSync downloads 1 book + 2 notes from server
	//    - In transaction: last_max_usn should be set to 3
	// 2. sendChanges tries to upload the orphaned note
	//    - Server returns 500 error (book not found)
	//    - Sets isBehind=true
	// 3. Retry stepSync queries after_usn=3 (from transaction state)
	//    - Server returns 0 fragments (no new data after USN 3)
	//    - processFragments returns maxUSN=0 (no fragments, so initialized value)
	//    - saveSyncState(tx, serverTime, 0) overwrites last_max_usn to 0
	// 4. Transaction commits
	//
	// BUG: last_max_usn ends up at 0 instead of 3
	clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync")

	// Step 4: Verify the bug - last_max_usn should be 0 (bug) instead of 3 (correct)
	var lastMaxUSN int
	cliDatabase.MustScan(t, "finding system last_max_usn",
		env.DB.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastMaxUSN),
		&lastMaxUSN)

	// This assertion FAILS with the current buggy code (lastMaxUSN will be 0)
	// After fixing the bug, this should pass (lastMaxUSN should be 3)
	//
	// BUG: The retry stepSync's empty response causes saveSyncState(tx, serverTime, 0)
	// which overwrites the last_max_usn that was set during the first stepSync.
	if lastMaxUSN == 0 {
		t.Logf("BUG REPRODUCED: last_max_usn is 0 (should be 3)")
		t.Logf("This happens because:")
		t.Logf("  1. First stepSync sets last_max_usn=3 in transaction")
		t.Logf("  2. Upload fails with 500 error, triggers retry stepSync")
		t.Logf("  3. Retry stepSync queries after_usn=3, gets 0 fragments")
		t.Logf("  4. processFragments returns maxUSN=0 (no fragments)")
		t.Logf("  5. saveSyncState(tx, serverTime, 0) overwrites last_max_usn to 0")
		t.Logf("  6. Transaction commits with last_max_usn=0")
		t.Logf("")
		t.Logf("Expected: last_max_usn should be 3 (server's max USN)")
		t.Logf("Actual:   last_max_usn is 0")
		t.Logf("")
		t.Logf("See /home/device10/development/dnote-memory-bank/artifacts-orphaned-notes-2025-10-25.md")
		t.Logf("for detailed investigation and proposed fixes.")
	}

	// When the bug is fixed, uncomment this assertion:
	assert.Equal(t, lastMaxUSN, 3, "last_max_usn should be 3 after syncing")

	// For now, just verify we can reproduce the bug
	if lastMaxUSN != 0 {
		t.Errorf("Failed to reproduce bug: last_max_usn is %d, expected 0 (indicating bug was NOT reproduced)", lastMaxUSN)
	}
}
