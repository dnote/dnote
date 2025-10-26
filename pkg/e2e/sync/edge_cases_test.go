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
	"io"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
	cliDatabase "github.com/dnote/dnote/pkg/cli/database"
	clitest "github.com/dnote/dnote/pkg/cli/testutils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

// TestSync_ConcurrentInitialSync reproduces the issue where two clients with identical
// local data syncing simultaneously to an empty server results in 500 errors.
//
// This demonstrates the race condition:
// - Client1 starts sync to empty server, gets empty server state
// - Client2 syncs.
// - Client1 tries to create same books → 409 "duplicate"
// - Client1 tries to create notes with wrong UUIDs → 500 "record not found"
// - stepSync recovers by renaming local books with _2 suffix
func TestSync_ConcurrentInitialSync(t *testing.T) {
	env := setupTestEnv(t)

	user := setupUserAndLogin(t, env)

	// Step 1: Create local data and sync
	clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "javascript", "-c", "js note from client1")
	clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync")
	checkState(t, env.DB, user, env.ServerDB, systemState{
		clientNoteCount:  1,
		clientBookCount:  1,
		clientLastMaxUSN: 2,
		clientLastSyncAt: serverTime.Unix(),
		serverNoteCount:  1,
		serverBookCount:  1,
		serverUserMaxUSN: 2,
	})

	// Step 2: Switch to new empty server to simulate concurrent initial sync scenario
	switchToEmptyServer(t, &env)
	user = setupUserAndLogin(t, env)

	// Set up client2 with separate database
	client2DB, client2DBPath := cliDatabase.InitTestFileDB(t)
	login(t, client2DB, env.ServerDB, user)
	client2DB.Close() // Close so CLI can access the database

	// Step 3: Client1 syncs to empty server, but during sync Client2 uploads same data
	// This simulates the race condition deterministically
	raceCallback := func(stdout io.Reader, stdin io.WriteCloser) error {
		// Wait for empty server prompt to ensure Client1 has called GetSyncState
		clitest.MustWaitForPrompt(t, stdout, clitest.PromptEmptyServer)

		// Now Client2 creates the same book and note via CLI (creating the race condition)
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "--dbPath", client2DBPath, "add", "javascript", "-c", "js note from client2")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "--dbPath", client2DBPath, "sync")

		// User confirms sync
		if _, err := io.WriteString(stdin, "y\n"); err != nil {
			return errors.Wrap(err, "confirming sync")
		}

		return nil
	}

	// Client1 continues sync - will hit 409 conflict, then 500 error, then recover
	clitest.MustWaitDnoteCmd(t, env.CmdOpts, raceCallback, cliBinaryName, "sync")

	// After sync:
	// - Server has 2 books: "javascript" (from client2) and "javascript_2" (from client1 renamed)
	// - Server has 2 notes
	// - Both clients should converge to the same state
	expectedState := systemState{
		clientNoteCount:  2, // both notes
		clientBookCount:  2, // javascript and javascript_2
		clientLastMaxUSN: 4, // 2 books + 2 notes
		clientLastSyncAt: serverTime.Unix(),
		serverNoteCount:  2,
		serverBookCount:  2,
		serverUserMaxUSN: 4,
	}
	checkState(t, env.DB, user, env.ServerDB, expectedState)

	// Client2 syncs again to download client1's data
	clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "--dbPath", client2DBPath, "sync")
	client2DB = clitest.MustOpenDatabase(t, client2DBPath)
	defer client2DB.Close()

	// Client2 should have converged to the same state as client1
	checkState(t, client2DB, user, env.ServerDB, expectedState)

	// Verify no orphaned notes on server
	var orphanedCount int
	if err := env.ServerDB.Raw(`
		SELECT COUNT(*) FROM notes
		WHERE book_uuid NOT IN (SELECT uuid FROM books)
	`).Scan(&orphanedCount).Error; err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, orphanedCount, 0, "server should have no orphaned notes after sync")
}
