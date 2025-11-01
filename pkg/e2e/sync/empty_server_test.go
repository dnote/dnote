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

package sync

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
	cliDatabase "github.com/dnote/dnote/pkg/cli/database"
	clitest "github.com/dnote/dnote/pkg/cli/testutils"
	"github.com/dnote/dnote/pkg/server/database"
	apitest "github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestSync_EmptyServer(t *testing.T) {
	t.Run("sync to empty server after syncing to non-empty server", func(t *testing.T) {
		// Test server data loss/wipe scenario (disaster recovery):
		// Verify empty server detection works when the server loses all its data

		env := setupTestEnv(t)

		user := setupUserAndLogin(t, env)

		// Step 1: Create local data and sync to server
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "js", "-c", "js1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "css", "-c", "css1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync")

		// Verify sync succeeded
		checkState(t, env.DB, user, env.ServerDB, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  2,
			serverBookCount:  2,
			serverUserMaxUSN: 4,
		})

		// Step 2: Switch to a completely new empty server
		switchToEmptyServer(t, &env)

		// Recreate user and session on new server
		user = setupUserAndLogin(t, env)

		// Step 3: Sync again - should detect empty server and prompt user
		// User confirms with "y"
		clitest.MustWaitDnoteCmd(t, env.CmdOpts, clitest.UserConfirmEmptyServerSync, cliBinaryName, "sync")

		// Step 4: Verify data was uploaded to the empty server
		checkState(t, env.DB, user, env.ServerDB, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  2,
			serverBookCount:  2,
			serverUserMaxUSN: 4,
		})

		// Verify the content is correct on both client and server
		var cliNote1JS, cliNote1CSS cliDatabase.Note
		var cliBookJS, cliBookCSS cliDatabase.Book
		cliDatabase.MustScan(t, "finding cliNote1JS", env.DB.QueryRow("SELECT uuid, body FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.Body)
		cliDatabase.MustScan(t, "finding cliNote1CSS", env.DB.QueryRow("SELECT uuid, body FROM notes WHERE body = ?", "css1"), &cliNote1CSS.UUID, &cliNote1CSS.Body)
		cliDatabase.MustScan(t, "finding cliBookJS", env.DB.QueryRow("SELECT uuid, label FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label)
		cliDatabase.MustScan(t, "finding cliBookCSS", env.DB.QueryRow("SELECT uuid, label FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label)

		assert.Equal(t, cliNote1JS.Body, "js1", "js note body mismatch")
		assert.Equal(t, cliNote1CSS.Body, "css1", "css note body mismatch")
		assert.Equal(t, cliBookJS.Label, "js", "js book label mismatch")
		assert.Equal(t, cliBookCSS.Label, "css", "css book label mismatch")

		// Verify on server side
		var serverNoteJS, serverNoteCSS database.Note
		var serverBookJS, serverBookCSS database.Book
		apitest.MustExec(t, env.ServerDB.Where("body = ?", "js1").First(&serverNoteJS), "finding server note js1")
		apitest.MustExec(t, env.ServerDB.Where("body = ?", "css1").First(&serverNoteCSS), "finding server note css1")
		apitest.MustExec(t, env.ServerDB.Where("label = ?", "js").First(&serverBookJS), "finding server book js")
		apitest.MustExec(t, env.ServerDB.Where("label = ?", "css").First(&serverBookCSS), "finding server book css")

		assert.Equal(t, serverNoteJS.Body, "js1", "server js note body mismatch")
		assert.Equal(t, serverNoteCSS.Body, "css1", "server css note body mismatch")
		assert.Equal(t, serverBookJS.Label, "js", "server js book label mismatch")
		assert.Equal(t, serverBookCSS.Label, "css", "server css book label mismatch")
	})

	t.Run("user cancels empty server prompt", func(t *testing.T) {
		env := setupTestEnv(t)

		user := setupUserAndLogin(t, env)

		// Step 1: Create local data and sync to server
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "js", "-c", "js1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "css", "-c", "css1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync")

		// Verify initial sync succeeded
		checkState(t, env.DB, user, env.ServerDB, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  2,
			serverBookCount:  2,
			serverUserMaxUSN: 4,
		})

		// Step 2: Switch to empty server
		switchToEmptyServer(t, &env)
		user = setupUserAndLogin(t, env)

		// Step 3: Sync again but user cancels with "n"
		output, err := clitest.WaitDnoteCmd(t, env.CmdOpts, clitest.UserCancelEmptyServerSync, cliBinaryName, "sync")
		if err == nil {
			t.Fatal("Expected sync to fail when user cancels, but it succeeded")
		}

		// Verify the prompt appeared
		if !strings.Contains(output, clitest.PromptEmptyServer) {
			t.Fatalf("Expected empty server warning in output, got: %s", output)
		}

		// Step 4: Verify local state unchanged (transaction rolled back)
		checkState(t, env.DB, user, env.ServerDB, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  0,
			serverBookCount:  0,
			serverUserMaxUSN: 0,
		})

		// Verify items still have original USN and dirty=false
		var book cliDatabase.Book
		var note cliDatabase.Note
		cliDatabase.MustScan(t, "checking book state", env.DB.QueryRow("SELECT usn, dirty FROM books WHERE label = ?", "js"), &book.USN, &book.Dirty)
		cliDatabase.MustScan(t, "checking note state", env.DB.QueryRow("SELECT usn, dirty FROM notes WHERE body = ?", "js1"), &note.USN, &note.Dirty)

		assert.NotEqual(t, book.USN, 0, "book USN should not be reset")
		assert.NotEqual(t, note.USN, 0, "note USN should not be reset")
		assert.Equal(t, book.Dirty, false, "book should not be marked dirty")
		assert.Equal(t, note.Dirty, false, "note should not be marked dirty")
	})

	t.Run("all local data is marked deleted - should not upload", func(t *testing.T) {
		// Test edge case: Server MaxUSN=0, local MaxUSN>0, but all items are deleted=true
		// Should NOT prompt because there's nothing to upload

		env := setupTestEnv(t)

		user := setupUserAndLogin(t, env)

		// Step 1: Create local data and sync to server
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "js", "-c", "js1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "css", "-c", "css1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync")

		// Verify initial sync succeeded
		checkState(t, env.DB, user, env.ServerDB, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  2,
			serverBookCount:  2,
			serverUserMaxUSN: 4,
		})

		// Step 2: Delete all local notes and books (mark as deleted)
		cliDatabase.MustExec(t, "marking all books deleted", env.DB, "UPDATE books SET deleted = 1")
		cliDatabase.MustExec(t, "marking all notes deleted", env.DB, "UPDATE notes SET deleted = 1")

		// Step 3: Switch to empty server
		switchToEmptyServer(t, &env)
		user = setupUserAndLogin(t, env)

		// Step 4: Sync - should NOT prompt because bookCount=0 and noteCount=0 (counting only deleted=0)
		// This should complete without user interaction
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync")

		// Verify no data was uploaded (server still empty, but client still has deleted items)
		// Check server is empty
		var serverNoteCount, serverBookCount int64
		apitest.MustExec(t, env.ServerDB.Model(&database.Note{}).Count(&serverNoteCount), "counting server notes")
		apitest.MustExec(t, env.ServerDB.Model(&database.Book{}).Count(&serverBookCount), "counting server books")
		assert.Equal(t, serverNoteCount, int64(0), "server should have no notes")
		assert.Equal(t, serverBookCount, int64(0), "server should have no books")

		// Check client still has the deleted items locally
		var clientNoteCount, clientBookCount int
		cliDatabase.MustScan(t, "counting client notes", env.DB.QueryRow("SELECT count(*) FROM notes WHERE deleted = 1"), &clientNoteCount)
		cliDatabase.MustScan(t, "counting client books", env.DB.QueryRow("SELECT count(*) FROM books WHERE deleted = 1"), &clientBookCount)
		assert.Equal(t, clientNoteCount, 2, "client should still have 2 deleted notes")
		assert.Equal(t, clientBookCount, 2, "client should still have 2 deleted books")

		// Verify lastMaxUSN was reset to 0
		var lastMaxUSN int
		cliDatabase.MustScan(t, "getting lastMaxUSN", env.DB.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastMaxUSN), &lastMaxUSN)
		assert.Equal(t, lastMaxUSN, 0, "lastMaxUSN should be reset to 0")
	})

	t.Run("race condition - other client uploads first", func(t *testing.T) {
		// This test exercises a race condition that can occur during sync:
		// While Client A is waiting for user input, Client B uploads data to the server.
		//
		// The empty server scenario is the natural place to test this because
		// an empty server detection triggers a prompt, at which point the test
		// can make client B upload data. We trigger the race condition deterministically.
		//
		// Test flow:
		// - Client A detects empty server and prompts user
		// - While waiting for confirmation, Client B uploads the same data via API
		// - Client A continues and handles the 409 conflict gracefully by:
		//   1. Detecting the 409 error when trying to CREATE books that already exist
		//   2. Running stepSync to pull the server's books (js, css)
		//   3. mergeBook renames local conflicts (js→js_2, css→css_2)
		//   4. Retrying sendChanges to upload the renamed books
		// - Result: Both clients' data is preserved (4 books total)

		env := setupTestEnv(t)

		user := setupUserAndLogin(t, env)

		// Step 1: Create local data and sync to establish lastMaxUSN > 0
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "js", "-c", "js1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "css", "-c", "css1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync")

		// Verify initial sync succeeded
		checkState(t, env.DB, user, env.ServerDB, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  2,
			serverBookCount:  2,
			serverUserMaxUSN: 4,
		})

		// Step 2: Switch to new empty server to simulate switching to empty server
		switchToEmptyServer(t, &env)

		// Create user on new server and login
		user = setupUserAndLogin(t, env)

		// Step 3: Trigger sync which will detect empty server and prompt user
		// Inside the callback (before confirming), we simulate Client B uploading via API.
		// We wait for the empty server prompt to ensure Client B uploads AFTER
		// GetSyncState but BEFORE the sync decision, creating the race condition deterministically
		raceCallback := func(stdout io.Reader, stdin io.WriteCloser) error {
			// First, wait for the prompt to ensure Client A has obtained the sync state from the server.
			clitest.MustWaitForPrompt(t, stdout, clitest.PromptEmptyServer)

			// Now Client B uploads the same data via API (after Client A got the sync state from the server
			// but before its sync decision)
			// This creates the race condition: Client A thinks server is empty, but Client B uploads data
			jsBookUUID := apiCreateBook(t, env, user, "js", "client B creating js book")
			cssBookUUID := apiCreateBook(t, env, user, "css", "client B creating css book")
			apiCreateNote(t, env, user, jsBookUUID, "js1", "client B creating js note")
			apiCreateNote(t, env, user, cssBookUUID, "css1", "client B creating css note")

			// Now user confirms
			if _, err := io.WriteString(stdin, "y\n"); err != nil {
				return errors.Wrap(err, "confirming sync")
			}

			return nil
		}

		// Step 4: Client A runs sync with race condition
		// The 409 conflict is automatically handled:
		// - When 409 is detected, isBehind flag is set
		// - stepSync pulls Client B's data
		// - mergeBook renames Client A's books to js_2, css_2
		// - Renamed books are uploaded
		// - Both clients' data is preserved.
		clitest.MustWaitDnoteCmd(t, env.CmdOpts, raceCallback, cliBinaryName, "sync")

		// Verify final state - both clients' data preserved
		checkState(t, env.DB, user, env.ServerDB, systemState{
			clientNoteCount:  4, // Both clients' notes
			clientBookCount:  4, // js, css, js_2, css_2
			clientLastMaxUSN: 8, // 4 from Client B + 4 from Client A's renamed books/notes
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  4,
			serverBookCount:  4,
			serverUserMaxUSN: 8,
		})

		// Verify server has both clients' books
		var svrBookJS, svrBookCSS, svrBookJS2, svrBookCSS2 database.Book
		apitest.MustExec(t, env.ServerDB.Where("label = ?", "js").First(&svrBookJS), "finding server book 'js'")
		apitest.MustExec(t, env.ServerDB.Where("label = ?", "css").First(&svrBookCSS), "finding server book 'css'")
		apitest.MustExec(t, env.ServerDB.Where("label = ?", "js_2").First(&svrBookJS2), "finding server book 'js_2'")
		apitest.MustExec(t, env.ServerDB.Where("label = ?", "css_2").First(&svrBookCSS2), "finding server book 'css_2'")

		assert.Equal(t, svrBookJS.Label, "js", "server should have book 'js' (Client B)")
		assert.Equal(t, svrBookCSS.Label, "css", "server should have book 'css' (Client B)")
		assert.Equal(t, svrBookJS2.Label, "js_2", "server should have book 'js_2' (Client A renamed)")
		assert.Equal(t, svrBookCSS2.Label, "css_2", "server should have book 'css_2' (Client A renamed)")

		// Verify client has all books
		var cliBookJS, cliBookCSS, cliBookJS2, cliBookCSS2 cliDatabase.Book
		cliDatabase.MustScan(t, "finding client book 'js'", env.DB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
		cliDatabase.MustScan(t, "finding client book 'css'", env.DB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
		cliDatabase.MustScan(t, "finding client book 'js_2'", env.DB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js_2"), &cliBookJS2.UUID, &cliBookJS2.Label, &cliBookJS2.USN)
		cliDatabase.MustScan(t, "finding client book 'css_2'", env.DB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css_2"), &cliBookCSS2.UUID, &cliBookCSS2.Label, &cliBookCSS2.USN)

		// Verify client UUIDs match server
		assert.Equal(t, cliBookJS.UUID, svrBookJS.UUID, "client 'js' UUID should match server")
		assert.Equal(t, cliBookCSS.UUID, svrBookCSS.UUID, "client 'css' UUID should match server")
		assert.Equal(t, cliBookJS2.UUID, svrBookJS2.UUID, "client 'js_2' UUID should match server")
		assert.Equal(t, cliBookCSS2.UUID, svrBookCSS2.UUID, "client 'css_2' UUID should match server")

		// Verify all items have non-zero USN (synced successfully)
		assert.NotEqual(t, cliBookJS.USN, 0, "client 'js' should have non-zero USN")
		assert.NotEqual(t, cliBookCSS.USN, 0, "client 'css' should have non-zero USN")
		assert.NotEqual(t, cliBookJS2.USN, 0, "client 'js_2' should have non-zero USN")
		assert.NotEqual(t, cliBookCSS2.USN, 0, "client 'css_2' should have non-zero USN")
	})

	t.Run("sync to server A, then B, then back to A, then back to B", func(t *testing.T) {
		// Test switching between two actual servers to verify:
		// 1. Empty server detection works when switching to empty server
		// 2. No false detection when switching back to non-empty servers
		// 3. Both servers maintain independent state across multiple switches

		env := setupTestEnv(t)

		// Create Server A with its own database
		serverA, serverDBA, err := setupTestServer(t, serverTime)
		if err != nil {
			t.Fatal(errors.Wrap(err, "setting up server A"))
		}
		defer serverA.Close()

		// Create Server B with its own database
		serverB, serverDBB, err := setupTestServer(t, serverTime)
		if err != nil {
			t.Fatal(errors.Wrap(err, "setting up server B"))
		}
		defer serverB.Close()

		// Step 1: Set up user on Server A and sync
		apiEndpointA := fmt.Sprintf("%s/api", serverA.URL)

		userA := apitest.SetupUserData(serverDBA, "alice@example.com", "pass1234")
		sessionA := apitest.SetupSession(serverDBA, userA)
		cliDatabase.MustExec(t, "inserting session_key", env.DB, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKey, sessionA.Key)
		cliDatabase.MustExec(t, "inserting session_key_expiry", env.DB, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKeyExpiry, sessionA.ExpiresAt.Unix())

		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "js", "-c", "js1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "add", "css", "-c", "css1")
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync", "--apiEndpoint", apiEndpointA)

		// Verify sync to Server A succeeded
		checkState(t, env.DB, userA, serverDBA, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  2,
			serverBookCount:  2,
			serverUserMaxUSN: 4,
		})

		// Step 2: Switch to Server B (empty) and sync
		apiEndpointB := fmt.Sprintf("%s/api", serverB.URL)

		// Set up user on Server B
		userB := apitest.SetupUserData(serverDBB, "alice@example.com", "pass1234")
		sessionB := apitest.SetupSession(serverDBB, userB)
		cliDatabase.MustExec(t, "updating session_key for B", env.DB, "UPDATE system SET value = ? WHERE key = ?", sessionB.Key, consts.SystemSessionKey)
		cliDatabase.MustExec(t, "updating session_key_expiry for B", env.DB, "UPDATE system SET value = ? WHERE key = ?", sessionB.ExpiresAt.Unix(), consts.SystemSessionKeyExpiry)

		// Should detect empty server and prompt
		clitest.MustWaitDnoteCmd(t, env.CmdOpts, clitest.UserConfirmEmptyServerSync, cliBinaryName, "sync", "--apiEndpoint", apiEndpointB)

		// Verify Server B now has data
		checkState(t, env.DB, userB, serverDBB, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  2,
			serverBookCount:  2,
			serverUserMaxUSN: 4,
		})

		// Step 3: Switch back to Server A and sync
		cliDatabase.MustExec(t, "updating session_key back to A", env.DB, "UPDATE system SET value = ? WHERE key = ?", sessionA.Key, consts.SystemSessionKey)
		cliDatabase.MustExec(t, "updating session_key_expiry back to A", env.DB, "UPDATE system SET value = ? WHERE key = ?", sessionA.ExpiresAt.Unix(), consts.SystemSessionKeyExpiry)

		// Should NOT trigger empty server detection (Server A has MaxUSN > 0)
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync", "--apiEndpoint", apiEndpointA)

		// Verify Server A still has its data
		checkState(t, env.DB, userA, serverDBA, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  2,
			serverBookCount:  2,
			serverUserMaxUSN: 4,
		})

		// Step 4: Switch back to Server B and sync again
		cliDatabase.MustExec(t, "updating session_key back to B", env.DB, "UPDATE system SET value = ? WHERE key = ?", sessionB.Key, consts.SystemSessionKey)
		cliDatabase.MustExec(t, "updating session_key_expiry back to B", env.DB, "UPDATE system SET value = ? WHERE key = ?", sessionB.ExpiresAt.Unix(), consts.SystemSessionKeyExpiry)

		// Should NOT trigger empty server detection (Server B now has MaxUSN > 0 from Step 2)
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync", "--apiEndpoint", apiEndpointB)

		// Verify both servers maintain independent state
		checkState(t, env.DB, userB, serverDBB, systemState{
			clientNoteCount:  2,
			clientBookCount:  2,
			clientLastMaxUSN: 4,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  2,
			serverBookCount:  2,
			serverUserMaxUSN: 4,
		})
	})

	t.Run("two clients with identical copied database sync to empty server", func(t *testing.T) {
		// Suppose we have two clients and server becomes empty (migration).
		// After the first client sync to empty server, the second client should trigger full sync.
		// Without the full sync, client2 will do step sync asking for changes after its stale USN,
		// get nothing from server, and potentially orphan notes during full sync.

		// Step 1: Create client1 with data and sync to ORIGINAL server
		env1 := setupTestEnv(t)
		user := setupUserAndLogin(t, env1)

		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "add", "js", "-c", "js1")
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "add", "css", "-c", "css1")
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "sync")
		// Add more data to create a higher USN
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "add", "go", "-c", "go1")
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "add", "rust", "-c", "rust1")
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "sync")
		// Verify initial sync succeeded (now with 4 notes, 4 books, USN=8)
		checkState(t, env1.DB, user, env1.ServerDB, systemState{
			clientNoteCount:  4,
			clientBookCount:  4,
			clientLastMaxUSN: 8,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  4,
			serverBookCount:  4,
			serverUserMaxUSN: 8,
		})

		// Step 2: Create client2 by copying client1's database (simulating same DB on two devices)
		env2 := setupTestEnv(t)
		// Copy the database file from client1 to client2
		dbPath1 := env1.DB.Filepath
		dbPath2 := env2.DB.Filepath
		// Close both DBs before copying
		env1.DB.Close()
		env2.DB.Close()

		// Copy the database file
		input, err := os.ReadFile(dbPath1)
		if err != nil {
			t.Fatal(errors.Wrap(err, "reading client1 database"))
		}
		if err := os.WriteFile(dbPath2, input, 0644); err != nil {
			t.Fatal(errors.Wrap(err, "writing client2 database"))
		}

		// Reopen databases
		env1.DB, err = cliDatabase.Open(dbPath1)
		if err != nil {
			t.Fatal(errors.Wrap(err, "reopening client1 database"))
		}
		env2.DB, err = cliDatabase.Open(dbPath2)
		if err != nil {
			t.Fatal(errors.Wrap(err, "reopening client2 database"))
		}

		// Verify client2 has identical data and USN=8 (stale) - same as client1
		// Note: at this point there's no server to compare against, we just check counts
		var client2MaxUSN, client2NoteCount, client2BookCount int
		cliDatabase.MustScan(t, "getting client2 maxUSN",
			env2.DB.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastMaxUSN),
			&client2MaxUSN)
		cliDatabase.MustScan(t, "counting client2 notes",
			env2.DB.QueryRow("SELECT count(*) FROM notes WHERE deleted = 0"),
			&client2NoteCount)
		cliDatabase.MustScan(t, "counting client2 books",
			env2.DB.QueryRow("SELECT count(*) FROM books WHERE deleted = 0"),
			&client2BookCount)

		assert.Equal(t, client2MaxUSN, 8, "client2 should have same maxUSN=8 as client1")
		assert.Equal(t, client2NoteCount, 4, "client2 should have 4 notes")
		assert.Equal(t, client2BookCount, 4, "client2 should have 4 books")

		// Step 3: Switch client1 to new empty server
		switchToEmptyServer(t, &env1)

		// Point client2 to the same new server
		env2.Server = env1.Server
		env2.ServerDB = env1.ServerDB

		// Update client2's API endpoint config to point to env1's server
		apiEndpoint := fmt.Sprintf("%s/api", env1.Server.URL)
		updateConfigAPIEndpoint(t, env2.TmpDir, apiEndpoint)

		// Create same user on new server
		user = setupUserAndLogin(t, env1)

		// Setup session for client2 (same user, same server)
		login(t, env2.DB, env2.ServerDB, user)

		// Step 4: Client1 syncs ONLY FIRST 2 BOOKS to empty server (simulates partial upload)
		// This creates the stale USN scenario: client2 has maxUSN=8, but server will only have maxUSN=4
		clitest.MustWaitDnoteCmd(t, env1.CmdOpts, clitest.UserConfirmEmptyServerSync,
			cliBinaryName, "sync")

		// Delete the last 2 books from client1 to prevent them being on server
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "remove", "go", "-y")
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "remove", "rust", "-y")

		// Sync deletions to server
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "sync")

		// Verify server has 2 active books/notes (go/rust deleted)
		var serverNoteCount, serverBookCount int64
		apitest.MustExec(t, env1.ServerDB.Model(&database.Note{}).Where("deleted = ?", false).Count(&serverNoteCount), "counting active server notes")
		apitest.MustExec(t, env1.ServerDB.Model(&database.Book{}).Where("deleted = ?", false).Count(&serverBookCount), "counting active server books")
		assert.Equal(t, int(serverNoteCount), 2, "server should have 2 active notes (go/rust deleted)")
		assert.Equal(t, int(serverBookCount), 2, "server should have 2 active books (go/rust deleted)")

		// Step 5: Client2 syncs
		// CRITICAL: Client2 has lastMaxUSN=8 (from copied DB), but server's max_usn is now ~8 but only has 2 books
		// Client2 will ask for changes after USN=8, get nothing, then try to upload its 4 books
		// This should trigger the orphaned notes scenario or require full sync
		clitest.RunDnoteCmd(t, env2.CmdOpts, cliBinaryName, "sync")

		// Step 6: Verify client2 has all data and NO orphaned notes
		var orphanedCount int
		cliDatabase.MustScan(t, "checking for orphaned notes",
			env2.DB.QueryRow(`
				SELECT COUNT(*) FROM notes
				WHERE deleted = 0
				AND book_uuid NOT IN (SELECT uuid FROM books WHERE deleted = 0)
			`), &orphanedCount)

		assert.Equal(t, orphanedCount, 0, "client2 should have no orphaned notes")

		// Verify client2 converged with server state
		// Note: checkState counts ALL records (including deleted ones)
		// During full sync, cleanLocalBooks/cleanLocalNotes DELETE local records not on server
		// So client2 ends up with only the 2 active books/notes
		// Server has 4 total (2 active + 2 deleted)
		var client2LastMaxUSN, client2LastSyncAt int
		var serverUserMaxUSN int
		cliDatabase.MustScan(t, "getting client2 lastMaxUSN",
			env2.DB.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastMaxUSN),
			&client2LastMaxUSN)
		var lastSyncAtStr string
		cliDatabase.MustScan(t, "getting client2 lastSyncAt",
			env2.DB.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastSyncAt),
			&lastSyncAtStr)
		lastSyncAtInt, _ := strconv.ParseInt(lastSyncAtStr, 10, 64)
		client2LastSyncAt = int(lastSyncAtInt)

		apitest.MustExec(t, env2.ServerDB.Table("users").Select("max_usn").Where("id = ?", user.ID).Scan(&serverUserMaxUSN), "getting server user max_usn")

		checkState(t, env2.DB, user, env2.ServerDB, systemState{
			clientNoteCount:  2, // Only active notes (deleted ones removed by cleanLocalNotes)
			clientBookCount:  2, // Only active books (deleted ones removed by cleanLocalBooks)
			clientLastMaxUSN: client2LastMaxUSN,
			clientLastSyncAt: int64(client2LastSyncAt),
			serverNoteCount:  4, // 2 active + 2 deleted
			serverBookCount:  4, // 2 active + 2 deleted
			serverUserMaxUSN: serverUserMaxUSN,
		})

		// Verify both clients have the expected books (css, js only - go/rust deleted)
		var client1BookCSS, client1BookJS, client2BookCSS, client2BookJS cliDatabase.Book
		cliDatabase.MustScan(t, "finding client1 book 'css'",
			env1.DB.QueryRow("SELECT uuid, label FROM books WHERE label = ? AND deleted = 0", "css"),
			&client1BookCSS.UUID, &client1BookCSS.Label)
		cliDatabase.MustScan(t, "finding client1 book 'js'",
			env1.DB.QueryRow("SELECT uuid, label FROM books WHERE label = ? AND deleted = 0", "js"),
			&client1BookJS.UUID, &client1BookJS.Label)
		cliDatabase.MustScan(t, "finding client2 book 'css'",
			env2.DB.QueryRow("SELECT uuid, label FROM books WHERE label = ? AND deleted = 0", "css"),
			&client2BookCSS.UUID, &client2BookCSS.Label)
		cliDatabase.MustScan(t, "finding client2 book 'js'",
			env2.DB.QueryRow("SELECT uuid, label FROM books WHERE label = ? AND deleted = 0", "js"),
			&client2BookJS.UUID, &client2BookJS.Label)

		assert.Equal(t, client1BookCSS.Label, "css", "client1 should have css book")
		assert.Equal(t, client1BookJS.Label, "js", "client1 should have js book")
		assert.Equal(t, client2BookCSS.Label, "css", "client2 should have css book")
		assert.Equal(t, client2BookJS.Label, "js", "client2 should have js book")

		// Verify go and rust books are deleted/absent on both clients
		var client2BookGo, client2BookRust cliDatabase.Book
		errGo := env2.DB.QueryRow("SELECT uuid, label FROM books WHERE label = ? AND deleted = 0", "go").Scan(&client2BookGo.UUID, &client2BookGo.Label)
		assert.Equal(t, errGo, sql.ErrNoRows, "client2 should not have non-deleted 'go' book")
		errRust := env2.DB.QueryRow("SELECT uuid, label FROM books WHERE label = ? AND deleted = 0", "rust").Scan(&client2BookRust.UUID, &client2BookRust.Label)
		assert.Equal(t, errRust, sql.ErrNoRows, "client2 should not have non-deleted 'rust' book")
	})

	t.Run("client with local data syncs after another client uploads to empty server - should not orphan notes", func(t *testing.T) {
		// This test reproduces the scenario where:
		// 1. Client1 has local data and syncs to original server
		// 2. Client2 has DIFFERENT local data and syncs to SAME original server
		// 3. Both clients switch to NEW empty server
		// 4. Client1 uploads to the new empty server (sets FullSyncBefore)
		// 5. Client2 syncs - should trigger full sync AND upload its local data
		//    WITHOUT orphaning notes due to cleanLocalBooks deleting them first

		// Step 1: Create client1 with local data on original server
		env1 := setupTestEnv(t)
		user := setupUserAndLogin(t, env1)
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "add", "client1-book", "-c", "client1-note")
		clitest.RunDnoteCmd(t, env1.CmdOpts, cliBinaryName, "sync")

		// Step 2: Create client2 with DIFFERENT local data on SAME original server
		env2 := setupTestEnv(t)
		// Point env2 to env1's server (the original server)
		env2.Server = env1.Server
		env2.ServerDB = env1.ServerDB
		apiEndpoint := fmt.Sprintf("%s/api", env1.Server.URL)
		updateConfigAPIEndpoint(t, env2.TmpDir, apiEndpoint)

		// Login client2 to the same server
		login(t, env2.DB, env2.ServerDB, user)
		clitest.RunDnoteCmd(t, env2.CmdOpts, cliBinaryName, "add", "client2-book", "-c", "client2-note")
		clitest.RunDnoteCmd(t, env2.CmdOpts, cliBinaryName, "sync")

		// Step 3: Both clients switch to NEW empty server
		switchToEmptyServer(t, &env1)
		env2.Server = env1.Server
		env2.ServerDB = env1.ServerDB
		apiEndpoint = fmt.Sprintf("%s/api", env1.Server.URL)
		updateConfigAPIEndpoint(t, env2.TmpDir, apiEndpoint)

		// Create same user on new server
		user = setupUserAndLogin(t, env1)
		login(t, env2.DB, env2.ServerDB, user)

		// Step 4: Client1 uploads to empty server
		clitest.MustWaitDnoteCmd(t, env1.CmdOpts, clitest.UserConfirmEmptyServerSync, cliBinaryName, "sync")

		// Verify server has client1's data and FullSyncBefore is set
		var serverUser database.User
		apitest.MustExec(t, env1.ServerDB.Where("id = ?", user.ID).First(&serverUser), "getting server user state")
		assert.Equal(t, serverUser.MaxUSN > 0, true, "server should have data after client1 upload")
		assert.Equal(t, serverUser.FullSyncBefore > 0, true, "server should have FullSyncBefore set")

		// Step 5: Client2 syncs - should trigger full sync due to FullSyncBefore
		// CRITICAL: Client2 has local data (client2-book, client2-note) that should be uploaded
		// Without the fix, cleanLocalBooks will delete client2-book before upload, orphaning client2-note
		clitest.RunDnoteCmd(t, env2.CmdOpts, cliBinaryName, "sync")

		// Step 6: Verify NO orphaned notes on client2
		var orphanedCount int
		cliDatabase.MustScan(t, "checking for orphaned notes on client2",
			env2.DB.QueryRow(`
				SELECT COUNT(*) FROM notes
				WHERE deleted = 0
				AND book_uuid NOT IN (SELECT uuid FROM books WHERE deleted = 0)
			`), &orphanedCount)

		assert.Equal(t, orphanedCount, 0, "client2 should have no orphaned notes after sync")

		// Step 7: Verify client2's data was uploaded to server
		var client2BookOnServer database.Book
		err := env2.ServerDB.Where("label = ? AND deleted = ?", "client2-book", false).First(&client2BookOnServer).Error
		assert.Equal(t, err, nil, "client2-book should exist on server")

		var client2NoteOnServer database.Note
		err = env2.ServerDB.Where("body = ? AND deleted = ?", "client2-note", false).First(&client2NoteOnServer).Error
		assert.Equal(t, err, nil, "client2-note should exist on server")

		// Step 8: Verify server has data from BOTH clients
		// Note: Both clients had synced to original server, so they each have 2 books + 2 notes locally.
		// When switching to new empty server, client1 uploads 2 books + 2 notes (USN 1-4).
		// Client2 then does full sync, downloads client1's uploads, marks its local data as dirty,
		// and uploads its version of the same 2 books + 2 notes with potentially different UUIDs.
		// The exact count depends on UUID conflict resolution, but we verify both original books exist.
		var serverBookCount, serverNoteCount int64
		apitest.MustExec(t, env2.ServerDB.Model(&database.Book{}).Where("deleted = ?", false).Count(&serverBookCount), "counting active server books")
		apitest.MustExec(t, env2.ServerDB.Model(&database.Note{}).Where("deleted = ?", false).Count(&serverNoteCount), "counting active server notes")

		// The main assertion: both original client books should exist
		var client1BookExists, client2BookExists bool
		err = env2.ServerDB.Model(&database.Book{}).Where("label = ? AND deleted = ?", "client1-book", false).First(&database.Book{}).Error
		client1BookExists = (err == nil)
		err = env2.ServerDB.Model(&database.Book{}).Where("label = ? AND deleted = ?", "client2-book", false).First(&database.Book{}).Error
		client2BookExists = (err == nil)

		assert.Equal(t, client1BookExists, true, "server should have client1-book")
		assert.Equal(t, client2BookExists, true, "server should have client2-book")
		assert.Equal(t, serverBookCount >= 2, true, "server should have at least 2 books")
		assert.Equal(t, serverNoteCount >= 2, true, "server should have at least 2 notes")
	})
}
