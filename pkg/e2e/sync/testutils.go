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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
	cliDatabase "github.com/dnote/dnote/pkg/cli/database"
	clitest "github.com/dnote/dnote/pkg/cli/testutils"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/controllers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/mailer"
	apitest "github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// testEnv holds the test environment for a single test
type testEnv struct {
	DB       *cliDatabase.DB
	CmdOpts  clitest.RunDnoteCmdOptions
	Server   *httptest.Server
	ServerDB *gorm.DB
	TmpDir   string
}

// setupTestEnv creates an isolated test environment with its own database and temp directory
func setupTestEnv(t *testing.T) testEnv {
	tmpDir := t.TempDir()

	// Create .dnote directory
	dnoteDir := filepath.Join(tmpDir, consts.DnoteDirName)
	if err := os.MkdirAll(dnoteDir, 0755); err != nil {
		t.Fatal(errors.Wrap(err, "creating dnote directory"))
	}

	// Create database at the expected path
	dbPath := filepath.Join(dnoteDir, consts.DnoteDBFileName)
	db := cliDatabase.InitTestFileDBRaw(t, dbPath)

	// Create server
	server, serverDB := setupNewServer(t)

	// Create config file with this server's endpoint
	apiEndpoint := fmt.Sprintf("%s/api", server.URL)
	updateConfigAPIEndpoint(t, tmpDir, apiEndpoint)

	// Create command options with XDG paths pointing to temp dir
	cmdOpts := clitest.RunDnoteCmdOptions{
		Env: []string{
			fmt.Sprintf("XDG_CONFIG_HOME=%s", tmpDir),
			fmt.Sprintf("XDG_DATA_HOME=%s", tmpDir),
			fmt.Sprintf("XDG_CACHE_HOME=%s", tmpDir),
		},
	}

	return testEnv{
		DB:       db,
		CmdOpts:  cmdOpts,
		Server:   server,
		ServerDB: serverDB,
		TmpDir:   tmpDir,
	}
}

// setupTestServer creates a test server with its own database
func setupTestServer(t *testing.T, serverTime time.Time) (*httptest.Server, *gorm.DB, error) {
	db := apitest.InitMemoryDB(t)

	mockClock := clock.NewMock()
	mockClock.SetNow(serverTime)

	a := app.NewTest()
	a.Clock = mockClock
	a.EmailTemplates = mailer.Templates{}
	a.EmailBackend = &apitest.MockEmailbackendImplementation{}
	a.DB = db

	server, err := controllers.NewServer(&a)
	if err != nil {
		return nil, nil, errors.Wrap(err, "initializing server")
	}

	return server, db, nil
}

// setupNewServer creates a new server and returns the server and database.
// This is useful when a test needs to switch to a new empty server.
func setupNewServer(t *testing.T) (*httptest.Server, *gorm.DB) {
	server, serverDB, err := setupTestServer(t, serverTime)
	if err != nil {
		t.Fatal(errors.Wrap(err, "setting up new test server"))
	}
	t.Cleanup(func() { server.Close() })

	return server, serverDB
}

// updateConfigAPIEndpoint updates the config file with the given API endpoint
func updateConfigAPIEndpoint(t *testing.T, tmpDir string, apiEndpoint string) {
	dnoteDir := filepath.Join(tmpDir, consts.DnoteDirName)
	configPath := filepath.Join(dnoteDir, consts.ConfigFilename)
	configContent := fmt.Sprintf("apiEndpoint: %s\n", apiEndpoint)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(errors.Wrap(err, "writing config file"))
	}
}

// switchToEmptyServer closes the current server and creates a new empty server,
// updating the config file to point to it.
func switchToEmptyServer(t *testing.T, env *testEnv) {
	// Close old server
	env.Server.Close()

	// Create new empty server
	env.Server, env.ServerDB = setupNewServer(t)

	// Update config file to point to new server
	apiEndpoint := fmt.Sprintf("%s/api", env.Server.URL)
	updateConfigAPIEndpoint(t, env.TmpDir, apiEndpoint)
}

// setupUser creates a test user in the server database
func setupUser(t *testing.T, env testEnv) database.User {
	user := apitest.SetupUserData(env.ServerDB, "alice@example.com", "pass1234")

	return user
}

// setupUserAndLogin creates a test user and logs them in on the CLI
func setupUserAndLogin(t *testing.T, env testEnv) database.User {
	user := setupUser(t, env)
	login(t, env.DB, env.ServerDB, user)

	return user
}

// login logs in the user in CLI
func login(t *testing.T, db *cliDatabase.DB, serverDB *gorm.DB, user database.User) {
	session := apitest.SetupSession(serverDB, user)

	cliDatabase.MustExec(t, "inserting session_key", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKey, session.Key)
	cliDatabase.MustExec(t, "inserting session_key_expiry", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKeyExpiry, session.ExpiresAt.Unix())
}

// apiCreateBook creates a book via the API and returns its UUID
func apiCreateBook(t *testing.T, env testEnv, user database.User, name, message string) string {
	res := doHTTPReq(t, env, "POST", "/v3/books", fmt.Sprintf(`{"name": "%s"}`, name), message, user)

	var resp controllers.CreateBookResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload for adding book"))
		return ""
	}

	return resp.Book.UUID
}

// apiPatchBook updates a book via the API
func apiPatchBook(t *testing.T, env testEnv, user database.User, uuid, payload, message string) {
	doHTTPReq(t, env, "PATCH", fmt.Sprintf("/v3/books/%s", uuid), payload, message, user)
}

// apiDeleteBook deletes a book via the API
func apiDeleteBook(t *testing.T, env testEnv, user database.User, uuid, message string) {
	doHTTPReq(t, env, "DELETE", fmt.Sprintf("/v3/books/%s", uuid), "", message, user)
}

// apiCreateNote creates a note via the API and returns its UUID
func apiCreateNote(t *testing.T, env testEnv, user database.User, bookUUID, body, message string) string {
	res := doHTTPReq(t, env, "POST", "/v3/notes", fmt.Sprintf(`{"book_uuid": "%s", "content": "%s"}`, bookUUID, body), message, user)

	var resp controllers.CreateNoteResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload for adding note"))
		return ""
	}

	return resp.Result.UUID
}

// apiPatchNote updates a note via the API
func apiPatchNote(t *testing.T, env testEnv, user database.User, noteUUID, payload, message string) {
	doHTTPReq(t, env, "PATCH", fmt.Sprintf("/v3/notes/%s", noteUUID), payload, message, user)
}

// apiDeleteNote deletes a note via the API
func apiDeleteNote(t *testing.T, env testEnv, user database.User, noteUUID, message string) {
	doHTTPReq(t, env, "DELETE", fmt.Sprintf("/v3/notes/%s", noteUUID), "", message, user)
}

// doHTTPReq performs an authenticated HTTP request and checks for errors
func doHTTPReq(t *testing.T, env testEnv, method, path, payload, message string, user database.User) *http.Response {
	apiEndpoint := fmt.Sprintf("%s/api", env.Server.URL)
	endpoint := fmt.Sprintf("%s%s", apiEndpoint, path)

	req, err := http.NewRequest(method, endpoint, strings.NewReader(payload))
	if err != nil {
		panic(errors.Wrap(err, "constructing http request"))
	}

	res := apitest.HTTPAuthDo(t, env.ServerDB, req, user)
	if res.StatusCode >= 400 {
		bs, err := io.ReadAll(res.Body)
		if err != nil {
			panic(errors.Wrap(err, "parsing response body for error"))
		}

		t.Errorf("%s. HTTP status %d. Message: %s", message, res.StatusCode, string(bs))
	}

	return res
}

// setupFunc is a function that sets up test data and returns IDs for assertions
type setupFunc func(t *testing.T, env testEnv, user database.User) map[string]string

// assertFunc is a function that asserts the expected state after sync
type assertFunc func(t *testing.T, env testEnv, user database.User, ids map[string]string)

// testSyncCmd is a test helper that sets up a test environment, runs setup, syncs, and asserts
func testSyncCmd(t *testing.T, fullSync bool, setup setupFunc, assert assertFunc) {
	env := setupTestEnv(t)

	user := setupUserAndLogin(t, env)
	ids := setup(t, env, user)

	if fullSync {
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync", "-f")
	} else {
		clitest.RunDnoteCmd(t, env.CmdOpts, cliBinaryName, "sync")
	}

	assert(t, env, user, ids)
}

// systemState represents the expected state of the sync system
type systemState struct {
	clientNoteCount  int
	clientBookCount  int
	clientLastMaxUSN int
	clientLastSyncAt int64
	serverNoteCount  int64
	serverBookCount  int64
	serverUserMaxUSN int
}

// checkState compares the state of the client and the server with the given system state
func checkState(t *testing.T, clientDB *cliDatabase.DB, user database.User, serverDB *gorm.DB, expected systemState) {
	var clientBookCount, clientNoteCount int
	cliDatabase.MustScan(t, "counting client notes", clientDB.QueryRow("SELECT count(*) FROM notes"), &clientNoteCount)
	cliDatabase.MustScan(t, "counting client books", clientDB.QueryRow("SELECT count(*) FROM books"), &clientBookCount)
	assert.Equal(t, clientNoteCount, expected.clientNoteCount, "client note count mismatch")
	assert.Equal(t, clientBookCount, expected.clientBookCount, "client book count mismatch")

	var clientLastMaxUSN int
	var clientLastSyncAt int64
	cliDatabase.MustScan(t, "finding system last_max_usn", clientDB.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastMaxUSN), &clientLastMaxUSN)
	cliDatabase.MustScan(t, "finding system last_sync_at", clientDB.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastSyncAt), &clientLastSyncAt)
	assert.Equal(t, clientLastMaxUSN, expected.clientLastMaxUSN, "client last_max_usn mismatch")
	assert.Equal(t, clientLastSyncAt, expected.clientLastSyncAt, "client last_sync_at mismatch")

	var serverBookCount, serverNoteCount int64
	apitest.MustExec(t, serverDB.Model(&database.Note{}).Count(&serverNoteCount), "counting server notes")
	apitest.MustExec(t, serverDB.Model(&database.Book{}).Count(&serverBookCount), "counting api notes")
	assert.Equal(t, serverNoteCount, expected.serverNoteCount, "server note count mismatch")
	assert.Equal(t, serverBookCount, expected.serverBookCount, "server book count mismatch")
	var serverUser database.User
	apitest.MustExec(t, serverDB.Where("id = ?", user.ID).First(&serverUser), "finding user")
	assert.Equal(t, serverUser.MaxUSN, expected.serverUserMaxUSN, "user max_usn mismatch")
}
