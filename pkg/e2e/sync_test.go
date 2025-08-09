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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	cliDatabase "github.com/dnote/dnote/pkg/cli/database"
	clitest "github.com/dnote/dnote/pkg/cli/testutils"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/controllers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/mailer"
	apitest "github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

var cliBinaryName string
var server *httptest.Server
var serverTime = time.Date(2017, time.March, 14, 21, 15, 0, 0, time.UTC)

var tmpDirPath string
var dnoteCmdOpts clitest.RunDnoteCmdOptions
var paths context.Paths

var testDir = "./tmp/.dnote"

func init() {
	tmpDirPath = fmt.Sprintf("%s/tmp", testDir)
	cliBinaryName = fmt.Sprintf("%s/test/cli/test-cli", testDir)
	dnoteCmdOpts = clitest.RunDnoteCmdOptions{
		Env: []string{
			fmt.Sprintf("XDG_CONFIG_HOME=%s", tmpDirPath),
			fmt.Sprintf("XDG_DATA_HOME=%s", tmpDirPath),
			fmt.Sprintf("XDG_CACHE_HOME=%s", tmpDirPath),
		},
	}

	paths = context.Paths{
		Data:   tmpDirPath,
		Cache:  tmpDirPath,
		Config: tmpDirPath,
	}
}

func clearTmp(t *testing.T) {
	if err := os.RemoveAll(tmpDirPath); err != nil {
		t.Fatal("cleaning tmp dir")
	}
}

func TestMain(m *testing.M) {
	// Set up server database
	apitest.InitTestDB()

	mockClock := clock.NewMock()
	mockClock.SetNow(serverTime)

	var err error
	server, err = controllers.NewServer(&app.App{
		Clock:          mockClock,
		EmailTemplates: mailer.Templates{},
		EmailBackend:   &apitest.MockEmailbackendImplementation{},
		DB:             apitest.DB,
		Config: config.Config{
			WebURL: os.Getenv("WebURL"),
		},
	})
	if err != nil {
		panic(errors.Wrap(err, "initializing router"))
	}

	defer server.Close()

	// Build binaries
	apiEndpoint := fmt.Sprintf("%s/api", server.URL)
	ldflags := fmt.Sprintf("-X main.apiEndpoint=%s", apiEndpoint)

	cmd := exec.Command("go", "build", "--tags", "fts5", "-o", cliBinaryName, "-ldflags", ldflags, "github.com/dnote/dnote/pkg/cli")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Print(errors.Wrap(err, "building a CLI binary").Error())
		log.Print(stderr.String())
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// helpers
func setupUser(t *testing.T, ctx *context.DnoteCtx) database.User {
	user := apitest.SetupUserData()
	apitest.SetupAccountData(user, "alice@example.com", "pass1234")

	// log in the user in CLI
	session := apitest.SetupSession(t, user)
	cliDatabase.MustExec(t, "inserting session_key", ctx.DB, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKey, session.Key)
	cliDatabase.MustExec(t, "inserting session_key_expiry", ctx.DB, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKeyExpiry, session.ExpiresAt.Unix())

	return user
}

func apiCreateBook(t *testing.T, user database.User, name, message string) string {
	res := doHTTPReq(t, "POST", "/v3/books", fmt.Sprintf(`{"name": "%s"}`, name), message, user)

	var resp controllers.CreateBookResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload for adding book"))
		return ""
	}

	return resp.Book.UUID
}

func apiPatchBook(t *testing.T, user database.User, uuid, payload, message string) {
	doHTTPReq(t, "PATCH", fmt.Sprintf("/v3/books/%s", uuid), payload, message, user)
}

func apiDeleteBook(t *testing.T, user database.User, uuid, message string) {
	doHTTPReq(t, "DELETE", fmt.Sprintf("/v3/books/%s", uuid), "", message, user)
}

func apiCreateNote(t *testing.T, user database.User, bookUUID, body, message string) string {
	res := doHTTPReq(t, "POST", "/v3/notes", fmt.Sprintf(`{"book_uuid": "%s", "content": "%s"}`, bookUUID, body), message, user)

	var resp controllers.CreateNoteResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload for adding note"))
		return ""
	}

	return resp.Result.UUID
}

func apiPatchNote(t *testing.T, user database.User, noteUUID, payload, message string) {
	doHTTPReq(t, "PATCH", fmt.Sprintf("/v3/notes/%s", noteUUID), payload, message, user)
}

func apiDeleteNote(t *testing.T, user database.User, noteUUID, message string) {
	doHTTPReq(t, "DELETE", fmt.Sprintf("/v3/notes/%s", noteUUID), "", message, user)
}

func doHTTPReq(t *testing.T, method, path, payload, message string, user database.User) *http.Response {
	apiEndpoint := fmt.Sprintf("%s/api", server.URL)
	endpoint := fmt.Sprintf("%s%s", apiEndpoint, path)

	req, err := http.NewRequest(method, endpoint, strings.NewReader(payload))
	if err != nil {
		panic(errors.Wrap(err, "constructing http request"))
	}

	res := apitest.HTTPAuthDo(t, req, user)
	if res.StatusCode >= 400 {
		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(errors.Wrap(err, "parsing response body for error"))
		}

		t.Errorf("%s. HTTP status %d. Message: %s", message, res.StatusCode, string(bs))
	}

	return res
}

type setupFunc func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string
type assertFunc func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string)

func testSyncCmd(t *testing.T, fullSync bool, setup setupFunc, assert assertFunc) {
	// clean up
	apitest.ClearData(apitest.DB)
	defer apitest.ClearData(apitest.DB)

	clearTmp(t)

	ctx := context.InitTestCtx(t, paths, nil)
	defer context.TeardownTestCtx(t, ctx)

	user := setupUser(t, &ctx)
	ids := setup(t, ctx, user)

	if fullSync {
		clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync", "-f")
	} else {
		clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
	}

	assert(t, ctx, user, ids)
}

type systemState struct {
	clientNoteCount  int
	clientBookCount  int
	clientLastMaxUSN int
	clientLastSyncAt int64
	serverNoteCount  int
	serverBookCount  int
	serverUserMaxUSN int
}

// checkState compares the state of the client and the server with the given system state
func checkState(t *testing.T, ctx context.DnoteCtx, user database.User, expected systemState) {
	serverDB := apitest.DB
	clientDB := ctx.DB

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

	var serverBookCount, serverNoteCount int
	apitest.MustExec(t, serverDB.Model(&database.Note{}).Count(&serverNoteCount), "counting server notes")
	apitest.MustExec(t, serverDB.Model(&database.Book{}).Count(&serverBookCount), "counting api notes")
	assert.Equal(t, serverNoteCount, expected.serverNoteCount, "server note count mismatch")
	assert.Equal(t, serverBookCount, expected.serverBookCount, "server book count mismatch")
	var serverUser database.User
	apitest.MustExec(t, serverDB.Where("id = ?", user.ID).First(&serverUser), "finding user")
	assert.Equal(t, serverUser.MaxUSN, expected.serverUserMaxUSN, "user max_usn mismatch")
}

// tests
func TestSync_Empty(t *testing.T) {
	setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
		return map[string]string{}
	}

	assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
		// Test
		checkState(t, ctx, user, systemState{
			clientNoteCount:  0,
			clientBookCount:  0,
			clientLastMaxUSN: 0,
			clientLastSyncAt: serverTime.Unix(),
			serverNoteCount:  0,
			serverBookCount:  0,
			serverUserMaxUSN: 0,
		})
	}

	testSyncCmd(t, false, setup, assert)
	testSyncCmd(t, true, setup, assert)
}

func TestSync_oneway(t *testing.T) {
	t.Run("cli to api only", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) {
			apiDB := apitest.DB
			apitest.MustExec(t, apiDB.Model(&user).Update("max_usn", 0), "updating user max_usn")

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js1")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css1")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js2")
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User) {
			apiDB := apitest.DB
			cliDB := ctx.DB

			// test client
			checkState(t, ctx, user, systemState{
				clientNoteCount:  3,
				clientBookCount:  2,
				clientLastMaxUSN: 5,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  3,
				serverBookCount:  2,
				serverUserMaxUSN: 5,
			})

			var cliBookJS, cliBookCSS cliDatabase.Book
			var cliNote1JS, cliNote2JS, cliNote1CSS cliDatabase.Note
			cliDatabase.MustScan(t, "finding cli book js", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cli book css", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote2JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js2"), &cliNote2JS.UUID, &cliNote2JS.Body, &cliNote2JS.USN)
			cliDatabase.MustScan(t, "finding cliNote1CSS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css1"), &cliNote1CSS.UUID, &cliNote1CSS.Body, &cliNote1CSS.USN)

			// assert on usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			assert.NotEqual(t, cliNote2JS.USN, 0, "cliNote2JS USN mismatch")
			assert.NotEqual(t, cliNote1CSS.USN, 0, "cliNote1CSS USN mismatch")
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliBookCSS.USN, 0, "cliBookCSS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote2JS.Body, "js2", "cliNote2JS Body mismatch")
			assert.Equal(t, cliNote1CSS.Body, "css1", "cliNote1CSS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote2JS.Deleted, false, "cliNote2JS Deleted mismatch")
			assert.Equal(t, cliNote1CSS.Deleted, false, "cliNote1CSS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")

			// test server
			var apiBookJS, apiBookCSS database.Book
			var apiNote1JS, apiNote2JS, apiNote1CSS database.Note
			apitest.MustExec(t, apiDB.Model(&database.Note{}).Where("uuid = ?", cliNote1JS.UUID).First(&apiNote1JS), "getting js1 note")
			apitest.MustExec(t, apiDB.Model(&database.Note{}).Where("uuid = ?", cliNote2JS.UUID).First(&apiNote2JS), "getting js2 note")
			apitest.MustExec(t, apiDB.Model(&database.Note{}).Where("uuid = ?", cliNote1CSS.UUID).First(&apiNote1CSS), "getting css1 note")
			apitest.MustExec(t, apiDB.Model(&database.Book{}).Where("uuid = ?", cliBookJS.UUID).First(&apiBookJS), "getting js book")
			apitest.MustExec(t, apiDB.Model(&database.Book{}).Where("uuid = ?", cliBookCSS.UUID).First(&apiBookCSS), "getting css book")

			// assert usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS usn mismatch")
			assert.NotEqual(t, apiNote2JS.USN, 0, "apiNote2JS usn mismatch")
			assert.NotEqual(t, apiNote1CSS.USN, 0, "apiNote1CSS usn mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS usn mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS usn mismatch")
			// client must have generated uuids
			assert.NotEqual(t, apiNote1JS.UUID, "", "apiNote1JS UUID mismatch")
			assert.NotEqual(t, apiNote2JS.UUID, "", "apiNote2JS UUID mismatch")
			assert.NotEqual(t, apiNote1CSS.UUID, "", "apiNote1CSS UUID mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiNote2JS.Deleted, false, "apiNote2JS Deleted mismatch")
			assert.Equal(t, apiNote1CSS.Deleted, false, "apiNote1CSS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")
			// assert on body and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote2JS.Body, "js2", "apiNote2JS Body mismatch")
			assert.Equal(t, apiNote1CSS.Body, "css1", "apiNote1CSS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
		}

		t.Run("stepSync", func(t *testing.T) {
			clearTmp(t)
			defer apitest.ClearData(apitest.DB)

			ctx := context.InitTestCtx(t, paths, nil)
			defer context.TeardownTestCtx(t, ctx)
			user := setupUser(t, &ctx)
			setup(t, ctx, user)

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			assert(t, ctx, user)
		})

		t.Run("fullSync", func(t *testing.T) {
			clearTmp(t)
			defer apitest.ClearData(apitest.DB)

			ctx := context.InitTestCtx(t, paths, nil)
			defer context.TeardownTestCtx(t, ctx)
			user := setupUser(t, &ctx)
			setup(t, ctx, user)

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync", "-f")

			assert(t, ctx, user)
		})
	})

	t.Run("cli to api with edit and delete", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) {
			apiDB := apitest.DB
			apitest.MustExec(t, apiDB.Model(&user).Update("max_usn", 0), "updating user max_usn")

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js1")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css1")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js2")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js3")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css2")

			var nid, nid2 string
			cliDB := ctx.DB
			cliDatabase.MustScan(t, "getting id of note to edit", cliDB.QueryRow("SELECT rowid FROM notes WHERE body = ?", "js3"), &nid)
			cliDatabase.MustScan(t, "getting id of note to delete", cliDB.QueryRow("SELECT rowid FROM notes WHERE body = ?", "css2"), &nid2)

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", nid, "-c", "js3-edited")
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "css", nid2)

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css3")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css4")
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  6,
				clientBookCount:  2,
				clientLastMaxUSN: 8,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  6,
				serverBookCount:  2,
				serverUserMaxUSN: 8,
			})

			// test cli
			var cliN1, cliN2, cliN3, cliN4, cliN5, cliN6 cliDatabase.Note
			var cliB1, cliB2 cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliN1", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliN1.UUID, &cliN1.Body, &cliN1.USN)
			cliDatabase.MustScan(t, "finding cliN2", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js2"), &cliN2.UUID, &cliN2.Body, &cliN2.USN)
			cliDatabase.MustScan(t, "finding cliN3", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js3-edited"), &cliN3.UUID, &cliN3.Body, &cliN3.USN)
			cliDatabase.MustScan(t, "finding cliN4", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css1"), &cliN4.UUID, &cliN4.Body, &cliN4.USN)
			cliDatabase.MustScan(t, "finding cliN5", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css3"), &cliN5.UUID, &cliN5.Body, &cliN5.USN)
			cliDatabase.MustScan(t, "finding cliN6", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css4"), &cliN6.UUID, &cliN6.Body, &cliN6.USN)
			cliDatabase.MustScan(t, "finding cliB1", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliB1.UUID, &cliB1.Label, &cliB1.USN)
			cliDatabase.MustScan(t, "finding cliB2", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliB2.UUID, &cliB2.Label, &cliB2.USN)

			// assert on usn
			assert.NotEqual(t, cliN1.USN, 0, "cliN1 USN mismatch")
			assert.NotEqual(t, cliN2.USN, 0, "cliN2 USN mismatch")
			assert.NotEqual(t, cliN3.USN, 0, "cliN3 USN mismatch")
			assert.NotEqual(t, cliN4.USN, 0, "cliN4 USN mismatch")
			assert.NotEqual(t, cliN5.USN, 0, "cliN5 USN mismatch")
			assert.NotEqual(t, cliN6.USN, 0, "cliN6 USN mismatch")
			assert.NotEqual(t, cliB1.USN, 0, "cliB1 USN mismatch")
			assert.NotEqual(t, cliB2.USN, 0, "cliB2 USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliN1.Body, "js1", "cliN1 Body mismatch")
			assert.Equal(t, cliN2.Body, "js2", "cliN2 Body mismatch")
			assert.Equal(t, cliN3.Body, "js3-edited", "cliN3 Body mismatch")
			assert.Equal(t, cliN4.Body, "css1", "cliN4 Body mismatch")
			assert.Equal(t, cliN5.Body, "css3", "cliN5 Body mismatch")
			assert.Equal(t, cliN6.Body, "css4", "cliN6 Body mismatch")
			assert.Equal(t, cliB1.Label, "js", "cliB1 Label mismatch")
			assert.Equal(t, cliB2.Label, "css", "cliB2 Label mismatch")
			// assert on deleted
			assert.Equal(t, cliN1.Deleted, false, "cliN1 Deleted mismatch")
			assert.Equal(t, cliN2.Deleted, false, "cliN2 Deleted mismatch")
			assert.Equal(t, cliN3.Deleted, false, "cliN3 Deleted mismatch")
			assert.Equal(t, cliN4.Deleted, false, "cliN4 Deleted mismatch")
			assert.Equal(t, cliN5.Deleted, false, "cliN5 Deleted mismatch")
			assert.Equal(t, cliN6.Deleted, false, "cliN6 Deleted mismatch")
			assert.Equal(t, cliB1.Deleted, false, "cliB1 Deleted mismatch")
			assert.Equal(t, cliB2.Deleted, false, "cliB2 Deleted mismatch")

			// test api
			var apiN1, apiN2, apiN3, apiN4, apiN5, apiN6 database.Note
			var apiB1, apiB2 database.Book
			apitest.MustExec(t, apiDB.Where("uuid = ?", cliN1.UUID).First(&apiN1), "finding apiN1")
			apitest.MustExec(t, apiDB.Where("uuid = ?", cliN2.UUID).First(&apiN2), "finding apiN2")
			apitest.MustExec(t, apiDB.Where("uuid = ?", cliN3.UUID).First(&apiN3), "finding apiN3")
			apitest.MustExec(t, apiDB.Where("uuid = ?", cliN4.UUID).First(&apiN4), "finding apiN4")
			apitest.MustExec(t, apiDB.Where("uuid = ?", cliN5.UUID).First(&apiN5), "finding apiN5")
			apitest.MustExec(t, apiDB.Where("uuid = ?", cliN6.UUID).First(&apiN6), "finding apiN6")
			apitest.MustExec(t, apiDB.Where("uuid = ?", cliB1.UUID).First(&apiB1), "finding apiB1")
			apitest.MustExec(t, apiDB.Where("uuid = ?", cliB2.UUID).First(&apiB2), "finding apiB2")

			// assert on usn
			assert.NotEqual(t, apiN1.USN, 0, "apiN1 usn mismatch")
			assert.NotEqual(t, apiN2.USN, 0, "apiN2 usn mismatch")
			assert.NotEqual(t, apiN3.USN, 0, "apiN3 usn mismatch")
			assert.NotEqual(t, apiN4.USN, 0, "apiN4 usn mismatch")
			assert.NotEqual(t, apiN5.USN, 0, "apiN5 usn mismatch")
			assert.NotEqual(t, apiN6.USN, 0, "apiN6 usn mismatch")
			assert.NotEqual(t, apiB1.USN, 0, "apiB1 usn mismatch")
			assert.NotEqual(t, apiB2.USN, 0, "apiB2 usn mismatch")
			// client must have generated uuids
			assert.NotEqual(t, apiN1.UUID, "", "apiN1 UUID mismatch")
			assert.NotEqual(t, apiN2.UUID, "", "apiN2 UUID mismatch")
			assert.NotEqual(t, apiN3.UUID, "", "apiN3 UUID mismatch")
			assert.NotEqual(t, apiN4.UUID, "", "apiN4 UUID mismatch")
			assert.NotEqual(t, apiN5.UUID, "", "apiN5 UUID mismatch")
			assert.NotEqual(t, apiN6.UUID, "", "apiN6 UUID mismatch")
			assert.NotEqual(t, apiB1.UUID, "", "apiB1 UUID mismatch")
			assert.NotEqual(t, apiB2.UUID, "", "apiB2 UUID mismatch")
			// assert on deleted
			assert.Equal(t, apiN1.Deleted, false, "apiN1 Deleted mismatch")
			assert.Equal(t, apiN2.Deleted, false, "apiN2 Deleted mismatch")
			assert.Equal(t, apiN3.Deleted, false, "apiN3 Deleted mismatch")
			assert.Equal(t, apiN4.Deleted, false, "apiN4 Deleted mismatch")
			assert.Equal(t, apiN5.Deleted, false, "apiN5 Deleted mismatch")
			assert.Equal(t, apiN6.Deleted, false, "apiN6 Deleted mismatch")
			assert.Equal(t, apiB1.Deleted, false, "apiB1 Deleted mismatch")
			assert.Equal(t, apiB2.Deleted, false, "apiB2 Deleted mismatch")
			// assert on body and labels
			assert.Equal(t, apiN1.Body, "js1", "apiN1 Body mismatch")
			assert.Equal(t, apiN2.Body, "js2", "apiN2 Body mismatch")
			assert.Equal(t, apiN3.Body, "js3-edited", "apiN3 Body mismatch")
			assert.Equal(t, apiN4.Body, "css1", "apiN4 Body mismatch")
			assert.Equal(t, apiN5.Body, "css3", "apiN5 Body mismatch")
			assert.Equal(t, apiN6.Body, "css4", "apiN6 Body mismatch")
			assert.Equal(t, apiB1.Label, "js", "apiB1 Label mismatch")
			assert.Equal(t, apiB2.Label, "css", "apiB2 Label mismatch")
		}

		t.Run("stepSync", func(t *testing.T) {
			clearTmp(t)
			defer apitest.ClearData(apitest.DB)

			ctx := context.InitTestCtx(t, paths, nil)
			defer context.TeardownTestCtx(t, ctx)
			user := setupUser(t, &ctx)
			setup(t, ctx, user)

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			assert(t, ctx, user)
		})

		t.Run("fullSync", func(t *testing.T) {
			clearTmp(t)
			defer apitest.ClearData(apitest.DB)

			ctx := context.InitTestCtx(t, paths, nil)
			defer context.TeardownTestCtx(t, ctx)
			user := setupUser(t, &ctx)
			setup(t, ctx, user)

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync", "-f")

			assert(t, ctx, user)
		})
	})

	t.Run("api to cli", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			apiDB := apitest.DB

			apitest.MustExec(t, apiDB.Model(&user).Update("max_usn", 0), "updating user max_usn")

			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")
			cssBookUUID := apiCreateBook(t, user, "css", "adding css book")
			cssNote1UUID := apiCreateNote(t, user, cssBookUUID, "css1", "adding css note 1")
			jsNote2UUID := apiCreateNote(t, user, jsBookUUID, "js2", "adding js note 2")
			cssNote2UUID := apiCreateNote(t, user, cssBookUUID, "css2", "adding css note 2")
			linuxBookUUID := apiCreateBook(t, user, "linux", "adding linux book")
			linuxNote1UUID := apiCreateNote(t, user, linuxBookUUID, "linux1", "adding linux note 1")
			apiPatchNote(t, user, jsNote2UUID, fmt.Sprintf(`{"book_uuid": "%s"}`, linuxBookUUID), "moving js note 2 to linux")
			apiDeleteNote(t, user, jsNote1UUID, "deleting js note 1")
			cssNote3UUID := apiCreateNote(t, user, cssBookUUID, "css3", "adding css note 3")
			bashBookUUID := apiCreateBook(t, user, "bash", "adding bash book")
			bashNote1UUID := apiCreateNote(t, user, bashBookUUID, "bash1", "adding bash note 1")

			// delete the linux book and its two notes
			apiDeleteBook(t, user, linuxBookUUID, "deleting linux book")

			apiPatchNote(t, user, cssNote2UUID, fmt.Sprintf(`{"content": "%s"}`, "css2-edited"), "editing css 2 body")
			bashNote2UUID := apiCreateNote(t, user, bashBookUUID, "bash2", "adding bash note 2")
			linuxBook2UUID := apiCreateBook(t, user, "linux", "adding new linux book")
			linux2Note1UUID := apiCreateNote(t, user, linuxBookUUID, "linux-new-1", "adding linux note 1")
			apiDeleteBook(t, user, jsBookUUID, "deleting js book")

			return map[string]string{
				"jsBookUUID":      jsBookUUID,
				"jsNote1UUID":     jsNote1UUID,
				"jsNote2UUID":     jsNote2UUID,
				"cssBookUUID":     cssBookUUID,
				"cssNote1UUID":    cssNote1UUID,
				"cssNote2UUID":    cssNote2UUID,
				"cssNote3UUID":    cssNote3UUID,
				"linuxBookUUID":   linuxBookUUID,
				"linuxNote1UUID":  linuxNote1UUID,
				"bashBookUUID":    bashBookUUID,
				"bashNote1UUID":   bashNote1UUID,
				"bashNote2UUID":   bashNote2UUID,
				"linuxBook2UUID":  linuxBook2UUID,
				"linux2Note1UUID": linux2Note1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  6,
				clientBookCount:  3,
				clientLastMaxUSN: 21,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  9,
				serverBookCount:  5,
				serverUserMaxUSN: 21,
			})

			// test server
			var apiNote1JS, apiNote2JS, apiNote1CSS, apiNote2CSS, apiNote3CSS, apiNote1Bash, apiNote2Bash, apiNote1Linux, apiNote2Linux, apiNote1LinuxDup database.Note
			var apiBookJS, apiBookCSS, apiBookBash, apiBookLinux, apiBookLinuxDup database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding api js note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote2UUID"]).First(&apiNote2JS), "finding api js note 2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssNote1UUID"]).First(&apiNote1CSS), "finding api css note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssNote2UUID"]).First(&apiNote2CSS), "finding api css note 2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssNote3UUID"]).First(&apiNote3CSS), "finding api css note 3")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["linuxNote1UUID"]).First(&apiNote1Linux), "finding api linux note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote2UUID"]).First(&apiNote2Linux), "finding api linux note 2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["bashNote1UUID"]).First(&apiNote1Bash), "finding api bash note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["bashNote2UUID"]).First(&apiNote2Bash), "finding api bash note 2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["linux2Note1UUID"]).First(&apiNote1LinuxDup), "finding api linux 2 note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding api js book")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding api css book")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["bashBookUUID"]).First(&apiBookBash), "finding api bash book")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["linuxBookUUID"]).First(&apiBookLinux), "finding api linux book")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["linuxBook2UUID"]).First(&apiBookLinuxDup), "finding api linux book 2")

			// assert on server Label
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiNote2JS.USN, 0, "apiNote2JS USN mismatch")
			assert.NotEqual(t, apiNote1CSS.USN, 0, "apiNote1CSS USN mismatch")
			assert.NotEqual(t, apiNote2CSS.USN, 0, "apiNote2CSS USN mismatch")
			assert.NotEqual(t, apiNote3CSS.USN, 0, "apiNote3CSS USN mismatch")
			assert.NotEqual(t, apiNote1Linux.USN, 0, "apiNote1Linux USN mismatch")
			assert.NotEqual(t, apiNote2Linux.USN, 0, "apiNote2Linux USN mismatch")
			assert.NotEqual(t, apiNote1Bash.USN, 0, "apiNote1Bash USN mismatch")
			assert.NotEqual(t, apiNote2Bash.USN, 0, "apiNote2Bash USN mismatch")
			assert.NotEqual(t, apiNote1LinuxDup.USN, 0, "apiNote1LinuxDup USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apibookJS USN mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apibookCSS USN mismatch")
			assert.NotEqual(t, apiBookBash.USN, 0, "apibookBash USN mismatch")
			assert.NotEqual(t, apiBookLinux.USN, 0, "apibookLinux USN mismatch")
			assert.NotEqual(t, apiBookLinuxDup.USN, 0, "apiBookLinuxDup USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote2JS.Body, "", "apiNote2JS Body mismatch")
			assert.Equal(t, apiNote1CSS.Body, "css1", "apiNote1CSS Body mismatch")
			assert.Equal(t, apiNote2CSS.Body, "css2-edited", "apiNote2CSS Body mismatch")
			assert.Equal(t, apiNote3CSS.Body, "css3", "apiNote3CSS Body mismatch")
			assert.Equal(t, apiNote1Linux.Body, "", "apiNote1Linux Body mismatch")
			assert.Equal(t, apiNote2Linux.Body, "", "apiNote2Linux Body mismatch")
			assert.Equal(t, apiNote1Bash.Body, "bash1", "apiNote1Bash Body mismatch")
			assert.Equal(t, apiNote2Bash.Body, "bash2", "apiNote2Bash Body mismatch")
			assert.Equal(t, apiNote1LinuxDup.Body, "linux-new-1", "apiNote1LinuxDup Body mismatch")
			assert.Equal(t, apiBookJS.Label, "", "apibookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apibookCSS Label mismatch")
			assert.Equal(t, apiBookBash.Label, "bash", "apibookBash Label mismatch")
			assert.Equal(t, apiBookLinux.Label, "", "apibookLinux Label mismatch")
			assert.Equal(t, apiBookLinuxDup.Label, "linux", "apiBookLinuxDup Label mismatch")
			// assert on uuids
			assert.NotEqual(t, apiNote1JS.UUID, "", "apiNote1JS UUID mismatch")
			assert.NotEqual(t, apiNote2JS.UUID, "", "apiNote2JS UUID mismatch")
			assert.NotEqual(t, apiNote1CSS.UUID, "", "apiNote1CSS UUID mismatch")
			assert.NotEqual(t, apiNote2CSS.UUID, "", "apiNote2CSS UUID mismatch")
			assert.NotEqual(t, apiNote3CSS.UUID, "", "apiNote3CSS UUID mismatch")
			assert.NotEqual(t, apiNote1Linux.UUID, "", "apiNote1Linux UUID mismatch")
			assert.NotEqual(t, apiNote2Linux.UUID, "", "apiNote2Linux UUID mismatch")
			assert.NotEqual(t, apiNote1Bash.UUID, "", "apiNote1Bash UUID mismatch")
			assert.NotEqual(t, apiNote2Bash.UUID, "", "apiNote2Bash UUID mismatch")
			assert.NotEqual(t, apiNote2Bash.UUID, "", "apiNote2Bash UUID mismatch")
			assert.NotEqual(t, apiBookJS.UUID, "", "apibookJS UUID mismatch")
			assert.NotEqual(t, apiBookCSS.UUID, "", "apibookCSS UUID mismatch")
			assert.NotEqual(t, apiBookBash.UUID, "", "apibookBash UUID mismatch")
			assert.NotEqual(t, apiBookLinux.UUID, "", "apibookLinux UUID mismatch")
			assert.NotEqual(t, apiBookLinuxDup.UUID, "", "apiBookLinuxDup UUID mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiNote2JS.Deleted, true, "apiNote2JS Deleted mismatch")
			assert.Equal(t, apiNote1CSS.Deleted, false, "apiNote1CSS Deleted mismatch")
			assert.Equal(t, apiNote2CSS.Deleted, false, "apiNote2CSS Deleted mismatch")
			assert.Equal(t, apiNote3CSS.Deleted, false, "apiNote3CSS Deleted mismatch")
			assert.Equal(t, apiNote1Linux.Deleted, true, "apiNote1Linux Deleted mismatch")
			assert.Equal(t, apiNote2Linux.Deleted, true, "apiNote2Linux Deleted mismatch")
			assert.Equal(t, apiNote1Bash.Deleted, false, "apiNote1Bash Deleted mismatch")
			assert.Equal(t, apiNote2Bash.Deleted, false, "apiNote2Bash Deleted mismatch")
			assert.Equal(t, apiNote2Bash.Deleted, false, "apiNote2Bash Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, true, "apibookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apibookCSS Deleted mismatch")
			assert.Equal(t, apiBookBash.Deleted, false, "apibookBash Deleted mismatch")
			assert.Equal(t, apiBookLinux.Deleted, true, "apibookLinux Deleted mismatch")
			assert.Equal(t, apiBookLinuxDup.Deleted, false, "apiBookLinuxDup Deleted mismatch")

			// test client
			var cliBookCSS, cliBookBash, cliBookLinux cliDatabase.Book
			var cliNote1CSS, cliNote2CSS, cliNote3CSS, cliNote1Bash, cliNote2Bash, cliNote1Linux cliDatabase.Note
			cliDatabase.MustScan(t, "finding cli book css", cliDB.QueryRow("SELECT label FROM books WHERE uuid = ?", ids["cssBookUUID"]), &cliBookCSS.Label)
			cliDatabase.MustScan(t, "finding cli book bash", cliDB.QueryRow("SELECT label FROM books WHERE uuid = ?", ids["bashBookUUID"]), &cliBookBash.Label)
			cliDatabase.MustScan(t, "finding cli book linux2", cliDB.QueryRow("SELECT label FROM books WHERE uuid = ?", ids["linuxBook2UUID"]), &cliBookLinux.Label)
			cliDatabase.MustScan(t, "finding cliNote1CSS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", apiNote1CSS.UUID), &cliNote1CSS.Body, &cliNote1CSS.USN)
			cliDatabase.MustScan(t, "finding cliNote2CSS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", apiNote2CSS.UUID), &cliNote2CSS.Body, &cliNote2CSS.USN)
			cliDatabase.MustScan(t, "finding cliNote3CSS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", apiNote3CSS.UUID), &cliNote3CSS.Body, &cliNote3CSS.USN)
			cliDatabase.MustScan(t, "finding cliNote1Bash", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", apiNote1Bash.UUID), &cliNote1Bash.Body, &cliNote1Bash.USN)
			cliDatabase.MustScan(t, "finding cliNote2Bash", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", apiNote2Bash.UUID), &cliNote2Bash.Body, &cliNote2Bash.USN)
			cliDatabase.MustScan(t, "finding cliNote2Bash", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", apiNote1LinuxDup.UUID), &cliNote1Linux.Body, &cliNote1Linux.USN)

			// assert on usn
			assert.NotEqual(t, cliNote1CSS.USN, 0, "cliNote1CSS usn mismatch")
			assert.NotEqual(t, cliNote2CSS.USN, 0, "cliNote2CSS usn mismatch")
			assert.NotEqual(t, cliNote3CSS.USN, 0, "cliNote3CSS usn mismatch")
			assert.NotEqual(t, cliNote1Bash.USN, 0, "cliNote1Bash usn mismatch")
			assert.NotEqual(t, cliNote2Bash.USN, 0, "cliNote2Bash usn mismatch")
			assert.NotEqual(t, cliNote1Linux.USN, 0, "cliNote1Linux usn mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1CSS.Body, "css1", "cliNote1CSS Body mismatch")
			assert.Equal(t, cliNote2CSS.Body, "css2-edited", "cliNote2CSS Body mismatch")
			assert.Equal(t, cliNote3CSS.Body, "css3", "cliNote3CSS Body mismatch")
			assert.Equal(t, cliNote1Bash.Body, "bash1", "cliNote1Bash Body mismatch")
			assert.Equal(t, cliNote2Bash.Body, "bash2", "cliNote2Bash Body mismatch")
			assert.Equal(t, cliNote1Linux.Body, "linux-new-1", "cliNote1Linux Body mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			assert.Equal(t, cliBookBash.Label, "bash", "cliBookBash Label mismatch")
			assert.Equal(t, cliBookLinux.Label, "linux", "cliBookLinux Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1CSS.Deleted, false, "cliNote1CSS Deleted mismatch")
			assert.Equal(t, cliNote2CSS.Deleted, false, "cliNote2CSS Deleted mismatch")
			assert.Equal(t, cliNote3CSS.Deleted, false, "cliNote3CSS Deleted mismatch")
			assert.Equal(t, cliNote1Bash.Deleted, false, "cliNote1Bash Deleted mismatch")
			assert.Equal(t, cliNote2Bash.Deleted, false, "cliNote2Bash Deleted mismatch")
			assert.Equal(t, cliNote1Linux.Deleted, false, "cliNote1Linux Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
			assert.Equal(t, cliBookBash.Deleted, false, "cliBookBash Deleted mismatch")
			assert.Equal(t, cliBookLinux.Deleted, false, "cliBookLinux Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})
}

func TestSync_twoway(t *testing.T) {
	t.Run("once", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			cliDB := ctx.DB

			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")
			cssBookUUID := apiCreateBook(t, user, "css", "adding css book")
			cssNote1UUID := apiCreateNote(t, user, cssBookUUID, "css1", "adding css note 1")
			jsNote2UUID := apiCreateNote(t, user, jsBookUUID, "js2", "adding js note 2")
			cssNote2UUID := apiCreateNote(t, user, cssBookUUID, "css2", "adding css note 2")
			linuxBookUUID := apiCreateBook(t, user, "linux", "adding linux book")
			linuxNote1UUID := apiCreateNote(t, user, linuxBookUUID, "linux1", "adding linux note 1")
			apiPatchNote(t, user, jsNote2UUID, fmt.Sprintf(`{"book_uuid": "%s"}`, linuxBookUUID), "moving js note 2 to linux")
			apiDeleteNote(t, user, jsNote1UUID, "deleting js note 1")
			cssNote3UUID := apiCreateNote(t, user, cssBookUUID, "css3", "adding css note 3")
			bashBookUUID := apiCreateBook(t, user, "bash", "adding bash book")
			bashNote1UUID := apiCreateNote(t, user, bashBookUUID, "bash1", "adding bash note 1")
			apiDeleteBook(t, user, linuxBookUUID, "deleting linux book")
			apiPatchNote(t, user, cssNote2UUID, fmt.Sprintf(`{"content": "%s"}`, "css2-edited"), "editing css 2 body")
			bashNote2UUID := apiCreateNote(t, user, bashBookUUID, "bash2", "adding bash note 2")
			linuxBook2UUID := apiCreateBook(t, user, "linux", "adding new linux book")
			linux2Note1UUID := apiCreateNote(t, user, linuxBookUUID, "linux-new-1", "adding linux note 1")
			apiDeleteBook(t, user, jsBookUUID, "deleting js book")

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js3")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "algorithms", "-c", "algorithms1")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js4")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "algorithms", "-c", "algorithms2")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "math", "-c", "math1")

			var nid string
			cliDatabase.MustScan(t, "getting id of note to remove", cliDB.QueryRow("SELECT rowid FROM notes WHERE body = ?", "js3"), &nid)

			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "algorithms")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css4")
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js", nid)

			return map[string]string{
				"jsBookUUID":      jsBookUUID,
				"jsNote1UUID":     jsNote1UUID,
				"jsNote2UUID":     jsNote2UUID,
				"cssBookUUID":     cssBookUUID,
				"cssNote1UUID":    cssNote1UUID,
				"cssNote2UUID":    cssNote2UUID,
				"cssNote3UUID":    cssNote3UUID,
				"linuxBookUUID":   linuxBookUUID,
				"linuxNote1UUID":  linuxNote1UUID,
				"bashBookUUID":    bashBookUUID,
				"bashNote1UUID":   bashNote1UUID,
				"bashNote2UUID":   bashNote2UUID,
				"linuxBook2UUID":  linuxBook2UUID,
				"linux2Note1UUID": linux2Note1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  9,
				clientBookCount:  6,
				clientLastMaxUSN: 27,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  12,
				serverBookCount:  8,
				serverUserMaxUSN: 27,
			})

			// test client
			var cliNote1CSS, cliNote2CSS, cliNote3CSS, cliNote1CSS2, cliNote1Bash, cliNote2Bash, cliNote1Linux, cliNote1Math, cliNote1JS cliDatabase.Note
			var cliBookCSS, cliBookCSS2, cliBookBash, cliBookLinux, cliBookMath, cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1CSS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css1"), &cliNote1CSS.UUID, &cliNote1CSS.Body, &cliNote1CSS.USN)
			cliDatabase.MustScan(t, "finding cliNote2CSS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css2-edited"), &cliNote2CSS.UUID, &cliNote2CSS.Body, &cliNote2CSS.USN)
			cliDatabase.MustScan(t, "finding cliNote3CSS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css3"), &cliNote3CSS.UUID, &cliNote3CSS.Body, &cliNote3CSS.USN)
			cliDatabase.MustScan(t, "finding cliNote1CSS2", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css4"), &cliNote1CSS2.UUID, &cliNote1CSS2.Body, &cliNote1CSS2.USN)
			cliDatabase.MustScan(t, "finding cliNote1Bash", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "bash1"), &cliNote1Bash.UUID, &cliNote1Bash.Body, &cliNote1Bash.USN)
			cliDatabase.MustScan(t, "finding cliNote2Bash", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "bash2"), &cliNote2Bash.UUID, &cliNote2Bash.Body, &cliNote2Bash.USN)
			cliDatabase.MustScan(t, "finding cliNote1Linux", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "linux-new-1"), &cliNote1Linux.UUID, &cliNote1Linux.Body, &cliNote1Linux.USN)
			cliDatabase.MustScan(t, "finding cliNote1Math", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "math1"), &cliNote1Math.UUID, &cliNote1Math.Body, &cliNote1Math.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js4"), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS2", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css_2"), &cliBookCSS2.UUID, &cliBookCSS2.Label, &cliBookCSS2.USN)
			cliDatabase.MustScan(t, "finding cliBookBash", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "bash"), &cliBookBash.UUID, &cliBookBash.Label, &cliBookBash.USN)
			cliDatabase.MustScan(t, "finding cliBookLinux", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "linux"), &cliBookLinux.UUID, &cliBookLinux.Label, &cliBookLinux.USN)
			cliDatabase.MustScan(t, "finding cliBookMath", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "math"), &cliBookMath.UUID, &cliBookMath.Label, &cliBookMath.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliNote1CSS.USN, 0, "cliNote1CSS USN mismatch")
			assert.NotEqual(t, cliNote2CSS.USN, 0, "cliNote2CSS USN mismatch")
			assert.NotEqual(t, cliNote3CSS.USN, 0, "cliNote3CSS USN mismatch")
			assert.NotEqual(t, cliNote1CSS2.USN, 0, "cliNote1CSS2 USN mismatch")
			assert.NotEqual(t, cliNote1Bash.USN, 0, "cliNote1Bash USN mismatch")
			assert.NotEqual(t, cliNote2Bash.USN, 0, "cliNote2Bash USN mismatch")
			assert.NotEqual(t, cliNote1Linux.USN, 0, "cliNote1Linux USN mismatch")
			assert.NotEqual(t, cliNote1Math.USN, 0, "cliNote1Math USN mismatch")
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			assert.NotEqual(t, cliBookCSS.USN, 0, "cliBookCSS USN mismatch")
			assert.NotEqual(t, cliBookCSS2.USN, 0, "cliBookCSS2 USN mismatch")
			assert.NotEqual(t, cliBookBash.USN, 0, "cliBookBash USN mismatch")
			assert.NotEqual(t, cliBookMath.USN, 0, "cliBookMath USN mismatch")
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1CSS.Body, "css1", "cliNote1CSS Body mismatch")
			assert.Equal(t, cliNote2CSS.Body, "css2-edited", "cliNote2CSS Body mismatch")
			assert.Equal(t, cliNote3CSS.Body, "css3", "cliNote3CSS Body mismatch")
			assert.Equal(t, cliNote1CSS2.Body, "css4", "cliNote1CSS2 Body mismatch")
			assert.Equal(t, cliNote1Bash.Body, "bash1", "cliNote1Bash Body mismatch")
			assert.Equal(t, cliNote2Bash.Body, "bash2", "cliNote2Bash Body mismatch")
			assert.Equal(t, cliNote1Linux.Body, "linux-new-1", "cliNote1Linux Body mismatch")
			assert.Equal(t, cliNote1Math.Body, "math1", "cliNote1Math Body mismatch")
			assert.Equal(t, cliNote1JS.Body, "js4", "cliNote1JS Body mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			assert.Equal(t, cliBookCSS2.Label, "css_2", "cliBookCSS2 Label mismatch")
			assert.Equal(t, cliBookBash.Label, "bash", "cliBookBash Label mismatch")
			assert.Equal(t, cliBookMath.Label, "math", "cliBookMath Label mismatch")
			assert.Equal(t, cliBookLinux.Label, "linux", "cliBookLinux Label mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1CSS.Deleted, false, "cliNote1CSS Deleted mismatch")
			assert.Equal(t, cliNote2CSS.Deleted, false, "cliNote2CSS Deleted mismatch")
			assert.Equal(t, cliNote3CSS.Deleted, false, "cliNote3CSS Deleted mismatch")
			assert.Equal(t, cliNote1CSS2.Deleted, false, "cliNote1CSS2 Deleted mismatch")
			assert.Equal(t, cliNote1Bash.Deleted, false, "cliNote1Bash Deleted mismatch")
			assert.Equal(t, cliNote2Bash.Deleted, false, "cliNote2Bash Deleted mismatch")
			assert.Equal(t, cliNote1Linux.Deleted, false, "cliNote1Linux Deleted mismatch")
			assert.Equal(t, cliNote1Math.Deleted, false, "cliNote1Math Deleted mismatch")
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
			assert.Equal(t, cliBookCSS2.Deleted, false, "cliBookCSS2 Deleted mismatch")
			assert.Equal(t, cliBookBash.Deleted, false, "cliBookBash Deleted mismatch")
			assert.Equal(t, cliBookMath.Deleted, false, "cliBookMath Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")

			// test server
			var apiNote1JS, apiNote1CSS, apiNote2CSS, apiNote3CSS, apiNote1Linux, apiNote2Linux, apiNote1Bash, apiNote2Bash, apiNote1LinuxDup, apiNote1CSS2, apiNote1Math, apiNote1JS2 database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding api js note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssNote1UUID"]).First(&apiNote1CSS), "finding api css note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssNote2UUID"]).First(&apiNote2CSS), "finding api css note 2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssNote3UUID"]).First(&apiNote3CSS), "finding api css note 3")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["linuxNote1UUID"]).First(&apiNote1Linux), "finding api linux note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote2UUID"]).First(&apiNote2Linux), "finding api linux note 2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["bashNote1UUID"]).First(&apiNote1Bash), "finding api bash note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["bashNote2UUID"]).First(&apiNote2Bash), "finding api bash note 2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["linux2Note1UUID"]).First(&apiNote1LinuxDup), "finding api linux 2 note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1CSS2.UUID).First(&apiNote1CSS2), "finding apiNote1CSS2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1Math.UUID).First(&apiNote1Math), "finding apiNote1Math")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1JS.UUID).First(&apiNote1JS2), "finding apiNote1JS2")
			var apiBookJS, apiBookCSS, apiBookLinux, apiBookBash, apiBookLinuxDup, apiBookCSS2, apiBookMath, apiBookJS2 database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding api js book")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding api css book")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["bashBookUUID"]).First(&apiBookBash), "finding api bash book")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["linuxBookUUID"]).First(&apiBookLinux), "finding api linux book")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["linuxBook2UUID"]).First(&apiBookLinuxDup), "finding api linux book 2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookCSS2.UUID).First(&apiBookCSS2), "finding apiBookCSS2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookMath.UUID).First(&apiBookMath), "finding apiBookMath")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookJS.UUID).First(&apiBookJS2), "finding apiBookJS2")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS usn mismatch")
			assert.NotEqual(t, apiNote1CSS.USN, 0, "apiNote1CSS usn mismatch")
			assert.NotEqual(t, apiNote2CSS.USN, 0, "apiNote2CSS usn mismatch")
			assert.NotEqual(t, apiNote3CSS.USN, 0, "apiNote3CSS usn mismatch")
			assert.NotEqual(t, apiNote1Linux.USN, 0, "apiNote1Linux usn mismatch")
			assert.NotEqual(t, apiNote2Linux.USN, 0, "apiNote2Linux usn mismatch")
			assert.NotEqual(t, apiNote1Bash.USN, 0, "apiNote1Bash usn mismatch")
			assert.NotEqual(t, apiNote2Bash.USN, 0, "apiNote2Bash usn mismatch")
			assert.NotEqual(t, apiNote1LinuxDup.USN, 0, "apiNote1LinuxDup usn mismatch")
			assert.NotEqual(t, apiNote1CSS2.USN, 0, "apiNoteCSS2 usn mismatch")
			assert.NotEqual(t, apiNote1Math.USN, 0, "apiNote1Math usn mismatch")
			assert.NotEqual(t, apiNote1JS2.USN, 0, "apiNote1JS2 usn mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS usn mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS usn mismatch")
			assert.NotEqual(t, apiBookLinux.USN, 0, "apiBookLinux usn mismatch")
			assert.NotEqual(t, apiBookBash.USN, 0, "apiBookBash usn mismatch")
			assert.NotEqual(t, apiBookLinuxDup.USN, 0, "apiBookLinuxDup usn mismatch")
			assert.NotEqual(t, apiBookCSS2.USN, 0, "apiBookCSS2 usn mismatch")
			assert.NotEqual(t, apiBookMath.USN, 0, "apiBookMath usn mismatch")
			assert.NotEqual(t, apiBookJS2.USN, 0, "apiBookJS2 usn mismatch")
			// assert on note bodys
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1CSS.Body, "css1", "apiNote1CSS Body mismatch")
			assert.Equal(t, apiNote2CSS.Body, "css2-edited", "apiNote2CSS Body mismatch")
			assert.Equal(t, apiNote3CSS.Body, "css3", "apiNote3CSS Body mismatch")
			assert.Equal(t, apiNote1Linux.Body, "", "apiNote1Linux Body mismatch")
			assert.Equal(t, apiNote2Linux.Body, "", "apiNote2Linux Body mismatch")
			assert.Equal(t, apiNote1Bash.Body, "bash1", "apiNote1Bash Body mismatch")
			assert.Equal(t, apiNote2Bash.Body, "bash2", "apiNote2Bash Body mismatch")
			assert.Equal(t, apiNote1LinuxDup.Body, "linux-new-1", "apiNote1LinuxDup Body mismatch")
			assert.Equal(t, apiNote1CSS2.Body, "css4", "apiNote1CSS2 Body mismatch")
			assert.Equal(t, apiNote1Math.Body, "math1", "apiNote1Math Body mismatch")
			assert.Equal(t, apiNote1JS2.Body, "js4", "apiNote1JS2 Body mismatch")
			// client must have generated uuids
			assert.NotEqual(t, apiNote1CSS2.UUID, "", "apiNote1CSS2 uuid mismatch")
			assert.NotEqual(t, apiNote1Math.UUID, "", "apiNote1Math uuid mismatch")
			assert.NotEqual(t, apiNote1JS2.UUID, "", "apiNote1JS2 uuid mismatch")
			// assert on labels
			assert.Equal(t, apiBookJS.Label, "", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			assert.Equal(t, apiBookLinux.Label, "", "apiBookLinux Label mismatch")
			assert.Equal(t, apiBookBash.Label, "bash", "apiBookBash Label mismatch")
			assert.Equal(t, apiBookLinuxDup.Label, "linux", "apiBookLinuxDup Label mismatch")
			assert.Equal(t, apiBookCSS2.Label, "css_2", "apiBookCSS2 Label mismatch")
			assert.Equal(t, apiBookMath.Label, "math", "apiBookMath Label mismatch")
			assert.Equal(t, apiBookJS2.Label, "js", "apiBookJS2 Label mismatch")
			// assert on note deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiNote1CSS.Deleted, false, "apiNote1CSS Deleted mismatch")
			assert.Equal(t, apiNote2CSS.Deleted, false, "apiNote2CSS Deleted mismatch")
			assert.Equal(t, apiNote3CSS.Deleted, false, "apiNote3CSS Deleted mismatch")
			assert.Equal(t, apiNote1Linux.Deleted, true, "apiNote1Linux Deleted mismatch")
			assert.Equal(t, apiNote2Linux.Deleted, true, "apiNote2Linux Deleted mismatch")
			assert.Equal(t, apiNote1Bash.Deleted, false, "apiNote1Bash Deleted mismatch")
			assert.Equal(t, apiNote2Bash.Deleted, false, "apiNote2Bash Deleted mismatch")
			assert.Equal(t, apiNote1LinuxDup.Deleted, false, "apiNote1LinuxDup Deleted mismatch")
			assert.Equal(t, apiNote1CSS2.Deleted, false, "apiNote1CSS2 Deleted mismatch")
			assert.Equal(t, apiNote1Math.Deleted, false, "apiNote1Math Deleted mismatch")
			assert.Equal(t, apiNote1JS2.Deleted, false, "apiNote1JS2 Deleted mismatch")
			// assert on book deleted
			assert.Equal(t, apiBookJS.Deleted, true, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")
			assert.Equal(t, apiBookLinux.Deleted, true, "apiBookLinux Deleted mismatch")
			assert.Equal(t, apiBookBash.Deleted, false, "apiBookBash Deleted mismatch")
			assert.Equal(t, apiBookLinuxDup.Deleted, false, "apiBookLinuxDup Deleted mismatch")
			assert.Equal(t, apiBookCSS2.Deleted, false, "apiBookCSS2 Deleted mismatch")
			assert.Equal(t, apiBookMath.Deleted, false, "apiBookMath Deleted mismatch")
			assert.Equal(t, apiBookJS2.Deleted, false, "apiBookJS2 Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("twice", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {

			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")
			cssBookUUID := apiCreateBook(t, user, "css", "adding css book")
			cssNote1UUID := apiCreateNote(t, user, cssBookUUID, "css1", "adding css note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js2")
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "math", "-c", "math1")

			var nid string
			cliDB := ctx.DB
			cliDatabase.MustScan(t, "getting id of note to edit", cliDB.QueryRow("SELECT rowid FROM notes WHERE body = ?", "math1"), &nid)

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "math", nid, "-c", "math1-edited")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			scssBookUUID := apiCreateBook(t, user, "scss", "adding a scss book")
			apiPatchNote(t, user, cssNote1UUID, fmt.Sprintf(`{"book_uuid": "%s"}`, scssBookUUID), "moving css note 1 to scss")

			var n1UUID string
			cliDatabase.MustScan(t, "getting math1-edited note UUID", cliDB.QueryRow("SELECT uuid FROM notes WHERE body = ?", "math1-edited"), &n1UUID)
			apiPatchNote(t, user, n1UUID, fmt.Sprintf(`{"content": "%s", "public": true}`, "math1-edited"), "editing math1 note")

			cssNote2UUID := apiCreateNote(t, user, cssBookUUID, "css2", "adding css note 2")
			apiDeleteBook(t, user, cssBookUUID, "deleting css book")

			bashBookUUID := apiCreateBook(t, user, "bash", "adding a bash book")
			algorithmsBookUUID := apiCreateBook(t, user, "algorithms", "adding a algorithms book")

			// 4. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js3")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "algorithms", "-c", "algorithms1")

			return map[string]string{
				"jsBookUUID":         jsBookUUID,
				"jsNote1UUID":        jsNote1UUID,
				"cssBookUUID":        cssBookUUID,
				"scssBookUUID":       scssBookUUID,
				"cssNote1UUID":       cssNote1UUID,
				"cssNote2UUID":       cssNote2UUID,
				"bashBookUUID":       bashBookUUID,
				"algorithmsBookUUID": algorithmsBookUUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			apiDB := apitest.DB
			cliDB := ctx.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  5,
				clientBookCount:  6,
				clientLastMaxUSN: 17,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  6,
				serverBookCount:  7,
				serverUserMaxUSN: 17,
			})

			// test client
			var cliNote1JS, cliNote2JS, cliNote1SCSS, cliNote1Math, cliNote1Alg2 cliDatabase.Note
			var cliBookJS, cliBookSCSS, cliBookMath, cliBookBash, cliBookAlg, cliBookAlg2 cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote2JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js3"), &cliNote2JS.UUID, &cliNote2JS.Body, &cliNote2JS.USN)
			cliDatabase.MustScan(t, "finding cliNote1SCSS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css1"), &cliNote1SCSS.UUID, &cliNote1SCSS.Body, &cliNote1SCSS.USN)
			cliDatabase.MustScan(t, "finding cliNote1Math", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "math1-edited"), &cliNote1Math.UUID, &cliNote1Math.Body, &cliNote1Math.USN)
			cliDatabase.MustScan(t, "finding cliNote1Alg2", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "algorithms1"), &cliNote1Alg2.UUID, &cliNote1Alg2.Body, &cliNote1Alg2.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliBookSCSS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "scss"), &cliBookSCSS.UUID, &cliBookSCSS.Label, &cliBookSCSS.USN)
			cliDatabase.MustScan(t, "finding cliBookMath", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "math"), &cliBookMath.UUID, &cliBookMath.Label, &cliBookMath.USN)
			cliDatabase.MustScan(t, "finding cliBookBash", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "bash"), &cliBookBash.UUID, &cliBookBash.Label, &cliBookBash.USN)
			cliDatabase.MustScan(t, "finding cliBookAlg", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "algorithms"), &cliBookAlg.UUID, &cliBookAlg.Label, &cliBookAlg.USN)
			cliDatabase.MustScan(t, "finding cliBookAlg2", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "algorithms_2"), &cliBookAlg2.UUID, &cliBookAlg2.Label, &cliBookAlg2.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			assert.NotEqual(t, cliNote2JS.USN, 0, "cliNote2JS USN mismatch")
			assert.NotEqual(t, cliNote1SCSS.USN, 0, "cliNote1SCSS USN mismatch")
			assert.NotEqual(t, cliNote1Math.USN, 0, "cliNote1Math USN mismatch")
			assert.NotEqual(t, cliNote1Alg2.USN, 0, "cliNote1Alg2 USN mismatch")
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliBookSCSS.USN, 0, "cliBookSCSS USN mismatch")
			assert.NotEqual(t, cliBookMath.USN, 0, "cliBookMath USN mismatch")
			assert.NotEqual(t, cliBookBash.USN, 0, "cliBookBash USN mismatch")
			assert.NotEqual(t, cliBookAlg.USN, 0, "cliBookAlg USN mismatch")
			assert.NotEqual(t, cliBookAlg2.USN, 0, "cliBookAlg2 USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote2JS.Body, "js3", "cliNote2JS Body mismatch")
			assert.Equal(t, cliNote1SCSS.Body, "css1", "cliNote1SCSS Body mismatch")
			assert.Equal(t, cliNote1Math.Body, "math1-edited", "cliNote1Math Body mismatch")
			assert.Equal(t, cliNote1Alg2.Body, "algorithms1", "cliNote1Alg2 Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookSCSS.Label, "scss", "cliBookSCSS Label mismatch")
			assert.Equal(t, cliBookMath.Label, "math", "cliBookMath Label mismatch")
			assert.Equal(t, cliBookBash.Label, "bash", "cliBookBash Label mismatch")
			assert.Equal(t, cliBookAlg.Label, "algorithms", "cliBookAlg Label mismatch")
			assert.Equal(t, cliBookAlg2.Label, "algorithms_2", "cliBookAlg2 Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote2JS.Deleted, false, "cliNote2JS Deleted mismatch")
			assert.Equal(t, cliNote1SCSS.Deleted, false, "cliNote1SCSS Deleted mismatch")
			assert.Equal(t, cliNote1Math.Deleted, false, "cliNote1Math Deleted mismatch")
			assert.Equal(t, cliNote1Alg2.Deleted, false, "cliNote1Alg2 Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookSCSS.Deleted, false, "cliBookSCSS Deleted mismatch")
			assert.Equal(t, cliBookMath.Deleted, false, "cliBookMath Deleted mismatch")
			assert.Equal(t, cliBookBash.Deleted, false, "cliBookBash Deleted mismatch")
			assert.Equal(t, cliBookAlg.Deleted, false, "cliBookAlg Deleted mismatch")
			assert.Equal(t, cliBookAlg2.Deleted, false, "cliBookAlg2 Deleted mismatch")

			// test server
			var apiNote1JS, apiNote2JS, apiNote1SCSS, apiNote2CSS, apiNote1Math, apiNote1Alg database.Note
			var apiBookJS, apiBookCSS, apiBookSCSS, apiBookMath, apiBookBash, apiBookAlg, apiBookAlg2 database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssNote2UUID"]).First(&apiNote2CSS), "finding apiNote2CSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote2JS.UUID).First(&apiNote2JS), "finding apiNote2JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssNote1UUID"]).First(&apiNote1SCSS), "finding apiNote1SCSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1Math.UUID).First(&apiNote1Math), "finding apiNote1Math")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1Alg2.UUID).First(&apiNote1Alg), "finding apiNote1Alg")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding apiBookCSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["bashBookUUID"]).First(&apiBookBash), "finding apiBookBash")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["scssBookUUID"]).First(&apiBookSCSS), "finding apiBookSCSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["algorithmsBookUUID"]).First(&apiBookAlg), "finding apiBookAlg")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookAlg2.UUID).First(&apiBookAlg2), "finding apiBookAlg2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookMath.UUID).First(&apiBookMath), "finding apiBookMath")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS usn mismatch")
			assert.NotEqual(t, apiNote2JS.USN, 0, "apiNote1JS usn mismatch")
			assert.NotEqual(t, apiNote1SCSS.USN, 0, "apiNote1JS usn mismatch")
			assert.NotEqual(t, apiNote2CSS.USN, 0, "apiNote1JS usn mismatch")
			assert.NotEqual(t, apiNote1Math.USN, 0, "apiNote1JS usn mismatch")
			assert.NotEqual(t, apiNote1Alg.USN, 0, "apiNote1JS usn mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBook1Alg usn mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS usn mismatch")
			assert.NotEqual(t, apiBookSCSS.USN, 0, "apibookSCSS usn mismatch")
			assert.NotEqual(t, apiBookMath.USN, 0, "apiBookMath usn mismatch")
			assert.NotEqual(t, apiBookBash.USN, 0, "apiBookBash usn mismatch")
			assert.NotEqual(t, apiBookAlg.USN, 0, "apiBookAlg usn mismatch")
			assert.NotEqual(t, apiBookAlg2.USN, 0, "apiBookAlg2 usn mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote2JS.Body, "js3", "apiNote2JS Body mismatch")
			assert.Equal(t, apiNote1SCSS.Body, "css1", "apiNote1SCSS Body mismatch")
			assert.Equal(t, apiNote2CSS.Body, "", "apiNote2CSS Body mismatch")
			assert.Equal(t, apiNote1Math.Body, "math1-edited", "apiNote1Math Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "", "apiBookCSS Label mismatch")
			assert.Equal(t, apiBookSCSS.Label, "scss", "apiBookSCSS Label mismatch")
			assert.Equal(t, apiBookMath.Label, "math", "apiBookMath Label mismatch")
			assert.Equal(t, apiBookBash.Label, "bash", "apiBookBash Label mismatch")
			assert.Equal(t, apiBookAlg.Label, "algorithms", "apiBookAlg Label mismatch")
			assert.Equal(t, apiBookAlg2.Label, "algorithms_2", "apiBookAlg2 Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiNote2JS.Deleted, false, "apiNote2JS Deleted mismatch")
			assert.Equal(t, apiNote1SCSS.Deleted, false, "apiNote1SCSS Deleted mismatch")
			assert.Equal(t, apiNote2CSS.Deleted, true, "apiNote2CSS Deleted mismatch")
			assert.Equal(t, apiNote1Math.Deleted, false, "apiNote1Math Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, true, "apiBookCSS Deleted mismatch")
			assert.Equal(t, apiBookSCSS.Deleted, false, "apiBookSCSS Deleted mismatch")
			assert.Equal(t, apiBookMath.Deleted, false, "apiBookMath Deleted mismatch")
			assert.Equal(t, apiBookBash.Deleted, false, "apiBookBash Deleted mismatch")
			assert.Equal(t, apiBookAlg.Deleted, false, "apiBookAlg Deleted mismatch")
			assert.Equal(t, apiBookAlg2.Deleted, false, "apiBookAlg2 Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("three times", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css1")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			goBookUUID := apiCreateBook(t, user, "go", "adding a go book")
			goNote1UUID := apiCreateNote(t, user, goBookUUID, "go1", "adding go note 1")

			// 4. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "html", "-c", "html1")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
				"goBookUUID":  goBookUUID,
				"goNote1UUID": goNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  4,
				clientBookCount:  4,
				clientLastMaxUSN: 8,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  4,
				serverBookCount:  4,
				serverUserMaxUSN: 8,
			})

			// test client
			var cliNote1JS, cliNote1CSS, cliNote1Go, cliNote1HTML cliDatabase.Note
			var cliBookJS, cliBookCSS, cliBookGo, cliBookHTML cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote1CSS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css1"), &cliNote1CSS.UUID, &cliNote1CSS.Body, &cliNote1CSS.USN)
			cliDatabase.MustScan(t, "finding cliNote1Go", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "go1"), &cliNote1Go.UUID, &cliNote1Go.Body, &cliNote1Go.USN)
			cliDatabase.MustScan(t, "finding cliNote1HTML", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "html1"), &cliNote1HTML.UUID, &cliNote1HTML.Body, &cliNote1HTML.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
			cliDatabase.MustScan(t, "finding cliBookGo", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "go"), &cliBookGo.UUID, &cliBookGo.Label, &cliBookGo.USN)
			cliDatabase.MustScan(t, "finding cliBookHTML", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "html"), &cliBookHTML.UUID, &cliBookHTML.Label, &cliBookHTML.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			assert.NotEqual(t, cliNote1CSS.USN, 0, "cliNote1CSS USN mismatch")
			assert.NotEqual(t, cliNote1Go.USN, 0, "cliNote1Go USN mismatch")
			assert.NotEqual(t, cliNote1HTML.USN, 0, "cliNote1HTML USN mismatch")
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliBookCSS.USN, 0, "cliBookCSS USN mismatch")
			assert.NotEqual(t, cliBookGo.USN, 0, "cliBookGo USN mismatch")
			assert.NotEqual(t, cliBookHTML.USN, 0, "cliBookHTML USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1CSS.Body, "css1", "cliNote1CSS Body mismatch")
			assert.Equal(t, cliNote1Go.Body, "go1", "cliNote1Go Body mismatch")
			assert.Equal(t, cliNote1HTML.Body, "html1", "cliNote1HTML Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			assert.Equal(t, cliBookGo.Label, "go", "cliBookGo Label mismatch")
			assert.Equal(t, cliBookHTML.Label, "html", "cliBookHTML Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote1CSS.Deleted, false, "cliNote1CSS Deleted mismatch")
			assert.Equal(t, cliNote1Go.Deleted, false, "cliNote1Go Deleted mismatch")
			assert.Equal(t, cliNote1HTML.Deleted, false, "cliNote1HTML Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
			assert.Equal(t, cliBookGo.Deleted, false, "cliBookGo Deleted mismatch")
			assert.Equal(t, cliBookHTML.Deleted, false, "cliBookHTML Deleted mismatch")

			// test server
			var apiNote1JS, apiNote1CSS, apiNote1Go, apiNote1HTML database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["goNote1UUID"]).First(&apiNote1Go), "finding apiNote1Go")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1CSS.UUID).First(&apiNote1CSS), "finding apiNote1CSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1HTML.UUID).First(&apiNote1HTML), "finding apiNote1HTML")
			var apiBookJS, apiBookCSS, apiBookGo, apiBookHTML database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["goBookUUID"]).First(&apiBookGo), "finding apiBookGo")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookCSS.UUID).First(&apiBookCSS), "finding apiBookCSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookHTML.UUID).First(&apiBookHTML), "finding apiBookHTML")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiNote1CSS.USN, 0, "apiNote1CSS USN mismatch")
			assert.NotEqual(t, apiNote1Go.USN, 0, "apiNote1Go USN mismatch")
			assert.NotEqual(t, apiNote1HTML.USN, 0, "apiNote1HTM USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			assert.NotEqual(t, apiBookGo.USN, 0, "apiBookGo USN mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS USN mismatch")
			assert.NotEqual(t, apiBookHTML.USN, 0, "apiBookHTML USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1CSS.Body, "css1", "apiNote1CSS Body mismatch")
			assert.Equal(t, apiNote1Go.Body, "go1", "apiNote1Go Body mismatch")
			assert.Equal(t, apiNote1HTML.Body, "html1", "apiNote1HTM Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookGo.Label, "go", "apiBookGo Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			assert.Equal(t, apiBookHTML.Label, "html", "apiBookHTML Label mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})
}

func TestSync(t *testing.T) {
	t.Run("client adds a book and a note", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js1")

			return map[string]string{}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 2,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 2,
			})

			// test client
			// assert on bodys and labels
			var cliNote1JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1JS.UUID).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookJS.UUID).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client deletes a book", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  0,
				clientLastMaxUSN: 5,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 5,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, true, "apiBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client deletes a note", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			cliDB := ctx.DB

			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			var nid string
			cliDatabase.MustScan(t, "getting id of note to remove", cliDB.QueryRow("SELECT rowid FROM notes WHERE uuid = ?", jsNote1UUID), &nid)
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js", nid)

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  1,
				clientLastMaxUSN: 3,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 3,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE label = ?", "js"), &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client edits a note", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			cliDB := ctx.DB

			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			var nid string
			cliDatabase.MustScan(t, "getting id of note to edit", cliDB.QueryRow("SELECT rowid FROM notes WHERE uuid = ?", jsNote1UUID), &nid)
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", nid, "-c", "js1-edited")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 3,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 3,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1-edited", "apiNote1JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliNote1JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1-edited", "cliNote1JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client edits a book by renaming it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", "-n", "js-edited")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 3,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 3,
			})

			// test client
			var cliBookJS cliDatabase.Book
			var cliNote1JS cliDatabase.Note
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, book_uuid, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.BookUUID, &cliNote1JS.USN)

			// assert on usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, cliBookJS.UUID, "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js-edited", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			// assert on dirty
			assert.Equal(t, cliNote1JS.Dirty, false, "cliNote1JS Dirty mismatch")
			assert.Equal(t, cliBookJS.Dirty, false, "cliBookJS Dirty mismatch")

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js-edited", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server adds a book", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")

			return map[string]string{
				"jsBookUUID": jsBookUUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  1,
				clientLastMaxUSN: 1,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  0,
				serverBookCount:  1,
				serverUserMaxUSN: 1,
			})

			// test server
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			// assert on bodys and labels
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE label = ?", "js"), &cliBookJS.Label, &cliBookJS.USN)
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server edits a book by renaming it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiPatchBook(t, user, jsBookUUID, fmt.Sprintf(`{"name": "%s"}`, "js-new-label"), "editing js book")

			return map[string]string{
				"jsBookUUID": jsBookUUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  1,
				clientLastMaxUSN: 2,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  0,
				serverBookCount:  1,
				serverUserMaxUSN: 2,
			})

			// test server
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiBookJS.Label, "js-new-label", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			// assert on bodys and labels
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)
			assert.Equal(t, cliBookJS.Label, "js-new-label", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server deletes a book", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiDeleteBook(t, user, jsBookUUID, "deleting js book")

			return map[string]string{
				"jsBookUUID": jsBookUUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  0,
				clientLastMaxUSN: 2,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  0,
				serverBookCount:  1,
				serverUserMaxUSN: 2,
			})

			// test server
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiBookJS.Label, "", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiBookJS.Deleted, true, "apiBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server adds a note", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 2,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 2,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliNote1JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server edits a note body", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiPatchNote(t, user, jsNote1UUID, fmt.Sprintf(`{"content": "%s"}`, "js1-edited"), "editing js note 1")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 3,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 3,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1-edited", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliNote1JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1-edited", "cliNote1JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server moves a note to another book", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")
			cssBookUUID := apiCreateBook(t, user, "css", "adding css book")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiPatchNote(t, user, jsNote1UUID, fmt.Sprintf(`{"book_uuid": "%s"}`, cssBookUUID), "moving js note 1 to css book")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
				"cssBookUUID": cssBookUUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  2,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  2,
				serverUserMaxUSN: 4,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS, apiBookCSS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding apiBookCSS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["cssBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")

			// test client
			var cliNote1JS cliDatabase.Note
			var cliBookJS, cliBookCSS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["cssBookUUID"]), &cliBookCSS.Label, &cliBookCSS.USN)
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server deletes a note", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiDeleteNote(t, user, jsNote1UUID, "deleting js note 1")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  1,
				clientLastMaxUSN: 3,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 3,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)
			// assert on bodys and labels
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client and server deletes the same book", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiDeleteBook(t, user, jsBookUUID, "deleting js book")

			// 4. on cli
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  0,
				clientLastMaxUSN: 6,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 6,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, true, "apiBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client and server deletes the same note", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			cliDB := ctx.DB

			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiDeleteNote(t, user, jsNote1UUID, "deleting js note 1")

			// 4. on cli
			var nid string
			cliDatabase.MustScan(t, "getting id of note to remove", cliDB.QueryRow("SELECT rowid FROM notes WHERE body = ?", "js1"), &nid)
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js", nid)

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			apiDB := apitest.DB
			cliDB := ctx.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  1,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 4,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server and client adds a note with same body", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 4. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js1")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  2,
				clientBookCount:  1,
				clientLastMaxUSN: 3,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  2,
				serverBookCount:  1,
				serverUserMaxUSN: 3,
			})

			// test client
			var cliNote1JS, cliNote2JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote2JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ? and uuid != ?", "js1", ids["jsNote1UUID"]), &cliNote2JS.UUID, &cliNote2JS.Body, &cliNote2JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote2JS.Body, "js1", "cliNote2JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote2JS.Deleted, false, "cliNote2JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")

			// test server
			var apiNote1JS, apiNote2JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote2JS.UUID).First(&apiNote2JS), "finding apiNote2JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiNote2JS.USN, 0, "apiNote2JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote2JS.Body, "js1", "apiNote2JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server and client adds a book with same label", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js1")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  2,
				clientBookCount:  2,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  2,
				serverBookCount:  2,
				serverUserMaxUSN: 4,
			})

			// test client
			var cliNote1JS, cliNote1JS2 cliDatabase.Note
			var cliBookJS, cliBookJS2 cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS2",
				cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ? AND uuid !=?", "js1", ids["jsNote1UUID"]), &cliNote1JS2.UUID, &cliNote1JS2.Body, &cliNote1JS2.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS2", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js_2"), &cliBookJS2.UUID, &cliBookJS2.Label, &cliBookJS2.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS2.Body, "js1", "cliNote1JS2 Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookJS2.Label, "js_2", "cliBookJS2 Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote1JS2.Deleted, false, "cliNote1JS2 Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookJS2.Deleted, false, "cliBookJS2 Deleted mismatch")

			// test server
			var apiNote1JS, apiNote1JS2 database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1JS2.UUID).First(&apiNote1JS2), "finding apiNote1JS2")
			var apiBookJS, apiBookJS2 database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookJS2.UUID).First(&apiBookJS2), "finding apiBookJS2")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiNote1JS2.USN, 0, "apiNote1JS2 USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			assert.NotEqual(t, apiBookJS2.USN, 0, "apiBookJS2 USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS2.Body, "js1", "apiNote1JS2 Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookJS2.Label, "js_2", "apiBookJS2 Label mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server and client adds two sets of books with same labels", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")
			cssBookUUID := apiCreateBook(t, user, "css", "adding css book")
			cssNote1UUID := apiCreateNote(t, user, cssBookUUID, "css1", "adding css note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js1")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css1")

			return map[string]string{
				"jsBookUUID":   jsBookUUID,
				"jsNote1UUID":  jsNote1UUID,
				"cssBookUUID":  cssBookUUID,
				"cssNote1UUID": cssNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  4,
				clientBookCount:  4,
				clientLastMaxUSN: 8,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  4,
				serverBookCount:  4,
				serverUserMaxUSN: 8,
			})

			// test client
			var cliNote1JS, cliNote1JS2, cliNote1CSS, cliNote1CSS2 cliDatabase.Note
			var cliBookJS, cliBookJS2, cliBookCSS, cliBookCSS2 cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS2",
				cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ? AND uuid != ?", "js1", ids["jsNote1UUID"]), &cliNote1JS2.UUID, &cliNote1JS2.Body, &cliNote1JS2.USN)
			cliDatabase.MustScan(t, "finding cliNote1CSS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE uuid = ?", ids["cssNote1UUID"]), &cliNote1CSS.UUID, &cliNote1CSS.Body, &cliNote1CSS.USN)
			cliDatabase.MustScan(t, "finding cliNote1CSS2",
				cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ? AND uuid != ?", "css1", ids["cssNote1UUID"]), &cliNote1CSS2.UUID, &cliNote1CSS2.Body, &cliNote1CSS2.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS2", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js_2"), &cliBookJS2.UUID, &cliBookJS2.Label, &cliBookJS2.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS2", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css_2"), &cliBookCSS2.UUID, &cliBookCSS2.Label, &cliBookCSS2.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS2.Body, "js1", "cliNote1JS2 Body mismatch")
			assert.Equal(t, cliNote1CSS.Body, "css1", "cliNote1CSS Body mismatch")
			assert.Equal(t, cliNote1CSS2.Body, "css1", "cliNote1CSS2 Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookJS2.Label, "js_2", "cliBookJS2 Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			assert.Equal(t, cliBookCSS2.Label, "css_2", "cliBookCSS2 Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote1JS2.Deleted, false, "cliNote1JS2 Deleted mismatch")
			assert.Equal(t, cliNote1CSS.Deleted, false, "cliNote1CSS Deleted mismatch")
			assert.Equal(t, cliNote1CSS2.Deleted, false, "cliNote1CSS2 Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookJS2.Deleted, false, "cliBookJS2 Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
			assert.Equal(t, cliBookCSS2.Deleted, false, "cliBookCSS2 Deleted mismatch")

			// test server
			var apiNote1JS, apiNote1JS2, apiNote1CSS, apiNote1CSS2 database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1JS2.UUID).First(&apiNote1JS2), "finding apiNote1JS2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssNote1UUID"]).First(&apiNote1CSS), "finding apiNote1CSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1CSS2.UUID).First(&apiNote1CSS2), "finding apiNote1CSS2")
			var apiBookJS, apiBookJS2, apiBookCSS, apiBookCSS2 database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookJS2.UUID).First(&apiBookJS2), "finding apiBookJS2")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding apiBookCSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookCSS2.UUID).First(&apiBookCSS2), "finding apiBookCSS2")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiNote1JS2.USN, 0, "apiNote1JS2 USN mismatch")
			assert.NotEqual(t, apiNote1CSS.USN, 0, "apiNote1CSS USN mismatch")
			assert.NotEqual(t, apiNote1CSS2.USN, 0, "apiNote1CSS2 USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			assert.NotEqual(t, apiBookJS2.USN, 0, "apiBookJS2 USN mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS USN mismatch")
			assert.NotEqual(t, apiBookCSS2.USN, 0, "apiBookCSS2 USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS2.Body, "js1", "apiNote1JS2 Body mismatch")
			assert.Equal(t, apiNote1CSS.Body, "css1", "apiNote1CSS2 Body mismatch")
			assert.Equal(t, apiNote1CSS2.Body, "css1", "apiNote1CSS2 Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookJS2.Label, "js_2", "apiBookJS2 Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			assert.Equal(t, apiBookCSS2.Label, "css_2", "apiBookCSS2 Label mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server and client adds notes to the same book", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 4. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js2")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  2,
				clientBookCount:  1,
				clientLastMaxUSN: 3,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  2,
				serverBookCount:  1,
				serverUserMaxUSN: 3,
			})

			// test client
			var cliNote1JS, cliNote2JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote2JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js2"), &cliNote2JS.UUID, &cliNote2JS.Body, &cliNote2JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote2JS.Body, "js2", "cliNote2JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote2JS.Deleted, false, "cliNote2JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")

			// test server
			var apiNote1JS, apiNote2JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote2JS.UUID).First(&apiNote2JS), "finding apiNote2JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiNote2JS.USN, 0, "apiNote2JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote2JS.Body, "js2", "apiNote2JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server and client adds a book with the same label and notes in it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "js", "-c", "js2")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  2,
				clientBookCount:  2,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  2,
				serverBookCount:  2,
				serverUserMaxUSN: 4,
			})

			// test client
			var cliNote1JS, cliNote2JS cliDatabase.Note
			var cliBookJS, cliBookJS2 cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote2JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js2"), &cliNote2JS.UUID, &cliNote2JS.Body, &cliNote2JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS2", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js_2"), &cliBookJS2.UUID, &cliBookJS2.Label, &cliBookJS2.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote2JS.Body, "js2", "cliNote2JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookJS2.Label, "js_2", "cliBookJS2 Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote2JS.Deleted, false, "cliNote2JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookJS2.Deleted, false, "cliBookJS2 Deleted mismatch")

			// test server
			var apiNote1JS, apiNote2JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote2JS.UUID).First(&apiNote2JS), "finding apiNote2JS")
			var apiBookJS, apiBookJS2 database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookJS2.UUID).First(&apiBookJS2), "finding apiBookJS2")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiNote2JS.USN, 0, "apiNote2JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			assert.NotEqual(t, apiBookJS2.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote2JS.Body, "js2", "apiNote2JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookJS2.Label, "js_2", "apiBookJS2 USN mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client and server edits bodys of the same note", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			cliDB := ctx.DB

			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			var nid string
			cliDatabase.MustScan(t, "getting id of note to edit", cliDB.QueryRow("SELECT rowid FROM notes WHERE uuid = ?", jsNote1UUID), &nid)
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", nid, "-c", "js1-edited-from-client")

			// 3. on server
			apiPatchNote(t, user, jsNote1UUID, fmt.Sprintf(`{"content": "%s"}`, "js1-edited-from-server"), "editing js note 1")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			resolvedBody := "<<<<<<< Local\njs1-edited-from-client\n=======\njs1-edited-from-server\n>>>>>>> Server\n"

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 4,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, resolvedBody, "apiNote1JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliNote1JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, resolvedBody, "cliNote1JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("clients deletes a note and server edits its body", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			cliDB := ctx.DB

			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			var nid string
			cliDatabase.MustScan(t, "getting id of note to remove", cliDB.QueryRow("SELECT rowid FROM notes WHERE uuid = ?", jsNote1UUID), &nid)
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js", nid)

			// 3. on server
			apiPatchNote(t, user, jsNote1UUID, fmt.Sprintf(`{"content": "%s"}`, "js1-edited"), "editing js note 1")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 3,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 3,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1-edited", "apiNote1JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliNote1JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1-edited", "cliNote1JS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("clients deletes a note and server moves it to another book", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			cliDB := ctx.DB

			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")
			cssBookUUID := apiCreateBook(t, user, "css", "adding css book")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			var nid string
			cliDatabase.MustScan(t, "getting id of note to remove", cliDB.QueryRow("SELECT rowid FROM notes WHERE uuid = ?", jsNote1UUID), &nid)
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js", nid)

			// 3. on server
			apiPatchNote(t, user, jsNote1UUID, fmt.Sprintf(`{"book_uuid": "%s"}`, cssBookUUID), "moving js note 1 to css book")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
				"cssBookUUID": cssBookUUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  2,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  2,
				serverUserMaxUSN: 4,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS, apiBookCSS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding apiBookCSS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["cssBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")

			// test client
			var cliNote1JS cliDatabase.Note
			var cliBookJS, cliBookCSS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["cssBookUUID"]), &cliBookCSS.Label, &cliBookCSS.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, ids["cssNote1UUID"], "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server deletes a note and client edits it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			cliDB := ctx.DB

			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiDeleteNote(t, user, jsNote1UUID, "deleting js note 1")

			// 4. on cli
			var nid string
			cliDatabase.MustScan(t, "getting id of note to edit", cliDB.QueryRow("SELECT rowid FROM notes WHERE body = ?", "js1"), &nid)
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", nid, "-c", "js1-edited")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 4,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1-edited", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliNote1JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT body, book_uuid, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.Body, &cliNote1JS.BookUUID, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1-edited", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, ids["jsBookUUID"], "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server deletes a book and client edits it by renaming it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiDeleteNote(t, user, jsNote1UUID, "deleting js note 1")

			// 4. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", "-n", "js-edited")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  1,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 4,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js-edited", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliBookJS.Label, "js-edited", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("server deletes a book and client edits a note in it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			cliDB := ctx.DB

			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiDeleteBook(t, user, jsBookUUID, "deleting js book")

			// 4. on cli
			var nid string
			cliDatabase.MustScan(t, "getting id of note to edit", cliDB.QueryRow("SELECT rowid FROM notes WHERE body = ?", "js1"), &nid)
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", nid, "-c", "js1-edited")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 6,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 6,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1-edited", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliNote1JS cliDatabase.Note
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT body, book_uuid, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.Body, &cliNote1JS.BookUUID, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1-edited", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, ids["jsBookUUID"], "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client deletes a book and server edits it by renaming it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js")

			// 3. on server
			apiPatchBook(t, user, jsBookUUID, fmt.Sprintf(`{"name": "%s"}`, "js-edited"), "editing js book")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  1,
				clientLastMaxUSN: 5,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 5,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js-edited", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")

			// test client
			var cliBookJS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.Label, &cliBookJS.USN)

			// test usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliBookJS.Label, "js-edited", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client deletes a book and server edits a note in it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiPatchNote(t, user, jsNote1UUID, fmt.Sprintf(`{"content": "%s"}`, "js1-edited"), "editing js1 note")

			// 4. on cli
			clitest.WaitDnoteCmd(t, dnoteCmdOpts, clitest.UserConfirm, cliBinaryName, "remove", "js")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  0,
				clientBookCount:  0,
				clientLastMaxUSN: 6,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 6,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["jsBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, true, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, true, "apiBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client and server edit a book by renaming it to a same name", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", "-n", "js-edited")

			// 3. on server
			apiPatchBook(t, user, jsBookUUID, fmt.Sprintf(`{"name": "%s"}`, "js-edited"), "editing js book")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 4,
			})

			// test client
			var cliBookJS cliDatabase.Book
			var cliNote1JS cliDatabase.Note
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, book_uuid, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.BookUUID, &cliNote1JS.USN)

			// assert on usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, cliBookJS.UUID, "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js-edited", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			// assert on dirty
			assert.Equal(t, cliNote1JS.Dirty, false, "cliNote1JS Dirty mismatch")
			assert.Equal(t, cliBookJS.Dirty, false, "cliBookJS Dirty mismatch")

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js-edited", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client and server edit a book by renaming it to different names", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", "-n", "js-edited-client")

			// 3. on server
			apiPatchBook(t, user, jsBookUUID, fmt.Sprintf(`{"name": "%s"}`, "js-edited-server"), "editing js book")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			// In this case, server's change wins and overwrites that of client's

			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  1,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  1,
				serverUserMaxUSN: 4,
			})

			// test client
			var cliBookJS cliDatabase.Book
			var cliNote1JS cliDatabase.Note
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE uuid = ?", ids["jsBookUUID"]), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, book_uuid, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.BookUUID, &cliNote1JS.USN)

			// assert on usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, cliBookJS.UUID, "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js-edited-server", "cliBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			// assert on dirty
			assert.Equal(t, cliNote1JS.Dirty, false, "cliNote1JS Dirty mismatch")
			assert.Equal(t, cliBookJS.Dirty, false, "cliBookJS Dirty mismatch")

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js-edited-server", "apiBookJS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client moves a note", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			cssBookUUID := apiCreateBook(t, user, "css", "adding a css book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "1", "-b", "css")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"cssBookUUID": cssBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  2,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  2,
				serverUserMaxUSN: 4,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS, apiBookCSS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding apiBookCSS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["cssBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")

			// test client
			var cliBookJS, cliBookCSS cliDatabase.Book
			var cliNote1JS cliDatabase.Note
			cliDatabase.MustScan(t, "finding cli book js", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cli book css", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, book_uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.BookUUID, &cliNote1JS.Body, &cliNote1JS.USN)

			// assert on usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliBookCSS.USN, 0, "cliBookCSS USN mismatch")
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, cliBookCSS.UUID, "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
			// assert on dirty
			assert.Equal(t, cliNote1JS.Dirty, false, "cliNote1JS Dirty mismatch")
			assert.Equal(t, cliBookJS.Dirty, false, "cliBookJS Dirty mismatch")
			assert.Equal(t, cliBookCSS.Dirty, false, "cliBookCSS Dirty mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client and server each moves a note to a same book", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			cssBookUUID := apiCreateBook(t, user, "css", "adding a css book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiPatchNote(t, user, jsNote1UUID, fmt.Sprintf(`{"book_uuid": "%s"}`, cssBookUUID), "moving js note 1 to css book")

			// 3. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "1", "-b", "css")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"cssBookUUID": cssBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  2,
				clientLastMaxUSN: 5,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  2,
				serverUserMaxUSN: 5,
			})

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS, apiBookCSS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding apiBookCSS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, ids["cssBookUUID"], "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")

			// test client
			var cliBookJS, cliBookCSS cliDatabase.Book
			var cliNote1JS cliDatabase.Note
			cliDatabase.MustScan(t, "finding cli book js", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cli book css", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, book_uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.BookUUID, &cliNote1JS.Body, &cliNote1JS.USN)

			// assert on usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliBookCSS.USN, 0, "cliBookCSS USN mismatch")
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, cliBookCSS.UUID, "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
			// assert on dirty
			assert.Equal(t, cliNote1JS.Dirty, false, "cliNote1JS Dirty mismatch")
			assert.Equal(t, cliBookJS.Dirty, false, "cliBookJS Dirty mismatch")
			assert.Equal(t, cliBookCSS.Dirty, false, "cliBookCSS Dirty mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client and server each moves a note to different books", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			cssBookUUID := apiCreateBook(t, user, "css", "adding a css book")
			linuxBookUUID := apiCreateBook(t, user, "linux", "adding a linux book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			apiPatchNote(t, user, jsNote1UUID, fmt.Sprintf(`{"book_uuid": "%s"}`, cssBookUUID), "moving js note 1 to css book")

			// 3. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "1", "-b", "linux")

			return map[string]string{
				"jsBookUUID":    jsBookUUID,
				"cssBookUUID":   cssBookUUID,
				"jsNote1UUID":   jsNote1UUID,
				"linuxBookUUID": linuxBookUUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			expectedNote1JSBody := `<<<<<<< Local
Moved to the book linux
=======
Moved to the book css
>>>>>>> Server

js1`

			checkState(t, ctx, user, systemState{
				clientNoteCount:  1,
				clientBookCount:  4,
				clientLastMaxUSN: 7,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  1,
				serverBookCount:  4,
				serverUserMaxUSN: 7,
			})

			// test client
			var cliBookJS, cliBookCSS, cliBookLinux, cliBookConflicts cliDatabase.Book
			var cliNote1JS cliDatabase.Note
			cliDatabase.MustScan(t, "finding cli book js", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cli book css", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
			cliDatabase.MustScan(t, "finding cliBookLinux", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "linux"), &cliBookLinux.UUID, &cliBookLinux.Label, &cliBookLinux.USN)
			cliDatabase.MustScan(t, "finding cliBookConflicts", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "conflicts"), &cliBookConflicts.UUID, &cliBookConflicts.Label, &cliBookConflicts.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, book_uuid, body, usn FROM notes WHERE uuid = ?", ids["jsNote1UUID"]), &cliNote1JS.UUID, &cliNote1JS.BookUUID, &cliNote1JS.Body, &cliNote1JS.USN)

			// assert on usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliBookCSS.USN, 0, "cliBookCSS USN mismatch")
			assert.NotEqual(t, cliBookLinux.USN, 0, "cliBookLinux USN mismatch")
			assert.NotEqual(t, cliBookConflicts.USN, 0, "cliBookConflicts USN mismatch")
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, expectedNote1JSBody, "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, cliBookConflicts.UUID, "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			assert.Equal(t, cliBookLinux.Label, "linux", "cliBookLinux Label mismatch")
			assert.Equal(t, cliBookConflicts.Label, "conflicts", "cliBookConflicts Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
			assert.Equal(t, cliBookLinux.Deleted, false, "cliBookLinux Deleted mismatch")
			assert.Equal(t, cliBookConflicts.Deleted, false, "cliBookConflicts Deleted mismatch")
			// assert on dirty
			assert.Equal(t, cliNote1JS.Dirty, false, "cliNote1JS Dirty mismatch")
			assert.Equal(t, cliBookJS.Dirty, false, "cliBookJS Dirty mismatch")
			assert.Equal(t, cliBookCSS.Dirty, false, "cliBookCSS Dirty mismatch")
			assert.Equal(t, cliBookLinux.Dirty, false, "cliBookLinux Dirty mismatch")
			assert.Equal(t, cliBookConflicts.Dirty, false, "cliBookConflicts Dirty mismatch")

			// test server
			var apiNote1JS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			var apiBookJS, apiBookCSS, apiBookLinux, apiBookConflicts database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding apiBookCSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["linuxBookUUID"]).First(&apiBookLinux), "finding apiBookLinux")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookConflicts.UUID).First(&apiBookConflicts), "finding apiBookConflicts")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS USN mismatch")
			assert.NotEqual(t, apiBookConflicts.USN, 0, "apiBookConflicts USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, expectedNote1JSBody, "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, apiBookConflicts.UUID, "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			assert.Equal(t, apiBookLinux.Label, "linux", "apiBookLinux Label mismatch")
			assert.Equal(t, apiBookConflicts.Label, "conflicts", "apiBookConflicts Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")
			assert.Equal(t, apiBookLinux.Deleted, false, "apiBookLinux Deleted mismatch")
			assert.Equal(t, apiBookConflicts.Deleted, false, "apiBookConflicts Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client adds a new book and moves a note into it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css1")

			cliDB := ctx.DB
			var nid string
			cliDatabase.MustScan(t, "getting id of note to edit", cliDB.QueryRow("SELECT rowid FROM notes WHERE body = ?", "js1"), &nid)
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", "js", nid, "-b", "css")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  2,
				clientBookCount:  2,
				clientLastMaxUSN: 5,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  2,
				serverBookCount:  2,
				serverUserMaxUSN: 5,
			})

			// test client
			var cliBookJS, cliBookCSS cliDatabase.Book
			var cliNote1JS, cliNote1CSS cliDatabase.Note
			cliDatabase.MustScan(t, "finding cli book js", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cli book css", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, book_uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.BookUUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote1CSS", cliDB.QueryRow("SELECT uuid, book_uuid, body, usn FROM notes WHERE body = ?", "css1"), &cliNote1CSS.UUID, &cliNote1CSS.BookUUID, &cliNote1CSS.Body, &cliNote1CSS.USN)

			// assert on usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliBookCSS.USN, 0, "cliBookCSS USN mismatch")
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			assert.NotEqual(t, cliNote1CSS.USN, 0, "cliNote1CSS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, cliBookCSS.UUID, "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliNote1CSS.Body, "css1", "cliNote1CSS Body mismatch")
			assert.Equal(t, cliNote1CSS.BookUUID, cliBookCSS.UUID, "cliNote1CSS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote1CSS.Deleted, false, "cliNote1CSS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
			// assert on dirty
			assert.Equal(t, cliNote1JS.Dirty, false, "cliNote1JS Dirty mismatch")
			assert.Equal(t, cliNote1CSS.Dirty, false, "cliNote1CSS Dirty mismatch")
			assert.Equal(t, cliBookJS.Dirty, false, "cliBookJS Dirty mismatch")
			assert.Equal(t, cliBookCSS.Dirty, false, "cliBookCSS Dirty mismatch")

			// test server
			var apiNote1JS, apiNote1CSS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1CSS.UUID).First(&apiNote1CSS), "finding apiNote1CSS")
			var apiBookJS, apiBookCSS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookCSS.UUID).First(&apiBookCSS), "finding apiBookCSS")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiNote1CSS.USN, 0, "apiNote1CSS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, apiBookCSS.UUID, "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiNote1CSS.Body, "css1", "apiNote1CSS Body mismatch")
			assert.Equal(t, apiNote1CSS.BookUUID, apiBookCSS.UUID, "apiNote1CSS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiNote1CSS.Deleted, false, "apiNote1CSS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})

	t.Run("client adds a duplicate book and moves a note into it", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			// 1. on server
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			// 2. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")

			// 3. on server
			cssBookUUID := apiCreateBook(t, user, "css", "adding a css book")

			// 3. on cli
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css1")

			var nid string
			cliDatabase.MustScan(t, "getting id of note to edit", ctx.DB.QueryRow("SELECT rowid FROM notes WHERE body = ?", "js1"), &nid)
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "edit", nid, "-b", "css")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"cssBookUUID": cssBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  2,
				clientBookCount:  3,
				clientLastMaxUSN: 6,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  2,
				serverBookCount:  3,
				serverUserMaxUSN: 6,
			})

			// test client
			var cliBookJS, cliBookCSS, cliBookCSS2 cliDatabase.Book
			var cliNote1JS, cliNote1CSS cliDatabase.Note
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS2", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css_2"), &cliBookCSS2.UUID, &cliBookCSS2.Label, &cliBookCSS2.USN)
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, book_uuid, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.BookUUID, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote1CSS", cliDB.QueryRow("SELECT uuid, body, book_uuid, usn FROM notes WHERE body = ?", "css1"), &cliNote1CSS.UUID, &cliNote1CSS.Body, &cliNote1CSS.BookUUID, &cliNote1CSS.USN)

			// assert on usn
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliBookCSS.USN, 0, "cliBookCSS USN mismatch")
			assert.NotEqual(t, cliBookCSS2.USN, 0, "cliBookCSS2 USN mismatch")
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			assert.NotEqual(t, cliNote1CSS.USN, 0, "cliNote1CSS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1JS.BookUUID, cliBookCSS2.UUID, "cliNote1JS BookUUID mismatch")
			assert.Equal(t, cliNote1CSS.Body, "css1", "cliNote1CSS Body mismatch")
			assert.Equal(t, cliNote1CSS.BookUUID, cliBookCSS2.UUID, "cliNote1CSS BookUUID mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			assert.Equal(t, cliBookCSS2.Label, "css_2", "cliBookCSS2 Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote1CSS.Deleted, false, "cliNote1CSS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")
			assert.Equal(t, cliBookCSS2.Deleted, false, "cliBookCSS2 Deleted mismatch")
			// assert on dirty
			assert.Equal(t, cliNote1JS.Dirty, false, "cliNote1JS Dirty mismatch")
			assert.Equal(t, cliNote1CSS.Dirty, false, "cliNote1CSS Dirty mismatch")
			assert.Equal(t, cliBookJS.Dirty, false, "cliBookJS Dirty mismatch")
			assert.Equal(t, cliBookCSS.Dirty, false, "cliBookCSS Dirty mismatch")
			assert.Equal(t, cliBookCSS2.Dirty, false, "cliBookCSS2 Dirty mismatch")

			// test server
			var apiNote1JS, apiNote1CSS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding apiNote1JS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1CSS.UUID).First(&apiNote1CSS), "finding apiNote1CSS")
			var apiBookJS, apiBookCSS, apiBookCSS2 database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding apiBookJS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["cssBookUUID"]).First(&apiBookCSS), "finding apiBookCSS")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookCSS2.UUID).First(&apiBookCSS2), "finding apiBookCSS2")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS USN mismatch")
			assert.NotEqual(t, apiNote1CSS.USN, 0, "apiNote1CSS USN mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS USN mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS USN mismatch")
			assert.NotEqual(t, apiBookCSS2.USN, 0, "apiBookCSS2 USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1JS.BookUUID, apiBookCSS2.UUID, "apiNote1JS BookUUID mismatch")
			assert.Equal(t, apiNote1CSS.Body, "css1", "apiNote1CSS Body mismatch")
			assert.Equal(t, apiNote1CSS.BookUUID, apiBookCSS2.UUID, "apiNote1CSS BookUUID mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			assert.Equal(t, apiBookCSS2.Label, "css_2", "apiBookCSS2 Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiNote1CSS.Deleted, false, "apiNote1CSS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")
			assert.Equal(t, apiBookCSS2.Deleted, false, "apiBookCSS2 Deleted mismatch")
		}

		testSyncCmd(t, false, setup, assert)
		testSyncCmd(t, true, setup, assert)
	})
}

func TestFullSync(t *testing.T) {
	t.Run("consecutively with stepSync", func(t *testing.T) {
		setup := func(t *testing.T, ctx context.DnoteCtx, user database.User) map[string]string {
			jsBookUUID := apiCreateBook(t, user, "js", "adding a js book")
			jsNote1UUID := apiCreateNote(t, user, jsBookUUID, "js1", "adding js note 1")

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "add", "css", "-c", "css1")

			return map[string]string{
				"jsBookUUID":  jsBookUUID,
				"jsNote1UUID": jsNote1UUID,
			}
		}

		assert := func(t *testing.T, ctx context.DnoteCtx, user database.User, ids map[string]string) {
			cliDB := ctx.DB
			apiDB := apitest.DB

			checkState(t, ctx, user, systemState{
				clientNoteCount:  2,
				clientBookCount:  2,
				clientLastMaxUSN: 4,
				clientLastSyncAt: serverTime.Unix(),
				serverNoteCount:  2,
				serverBookCount:  2,
				serverUserMaxUSN: 4,
			})

			// test client
			var cliNote1JS, cliNote1CSS cliDatabase.Note
			var cliBookJS, cliBookCSS cliDatabase.Book
			cliDatabase.MustScan(t, "finding cliNote1JS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "js1"), &cliNote1JS.UUID, &cliNote1JS.Body, &cliNote1JS.USN)
			cliDatabase.MustScan(t, "finding cliNote1CSS", cliDB.QueryRow("SELECT uuid, body, usn FROM notes WHERE body = ?", "css1"), &cliNote1CSS.UUID, &cliNote1CSS.Body, &cliNote1CSS.USN)
			cliDatabase.MustScan(t, "finding cliBookJS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "js"), &cliBookJS.UUID, &cliBookJS.Label, &cliBookJS.USN)
			cliDatabase.MustScan(t, "finding cliBookCSS", cliDB.QueryRow("SELECT uuid, label, usn FROM books WHERE label = ?", "css"), &cliBookCSS.UUID, &cliBookCSS.Label, &cliBookCSS.USN)

			// test usn
			assert.NotEqual(t, cliNote1JS.USN, 0, "cliNote1JS USN mismatch")
			assert.NotEqual(t, cliNote1CSS.USN, 0, "cliNote1CSS USN mismatch")
			assert.NotEqual(t, cliBookJS.USN, 0, "cliBookJS USN mismatch")
			assert.NotEqual(t, cliBookCSS.USN, 0, "cliBookCSS USN mismatch")
			// assert on bodys and labels
			assert.Equal(t, cliNote1JS.Body, "js1", "cliNote1JS Body mismatch")
			assert.Equal(t, cliNote1CSS.Body, "css1", "cliNote1CSS Body mismatch")
			assert.Equal(t, cliBookJS.Label, "js", "cliBookJS Label mismatch")
			assert.Equal(t, cliBookCSS.Label, "css", "cliBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, cliNote1JS.Deleted, false, "cliNote1JS Deleted mismatch")
			assert.Equal(t, cliNote1CSS.Deleted, false, "cliNote1CSS Deleted mismatch")
			assert.Equal(t, cliBookJS.Deleted, false, "cliBookJS Deleted mismatch")
			assert.Equal(t, cliBookCSS.Deleted, false, "cliBookCSS Deleted mismatch")

			// test server
			var apiNote1JS, apiNote1CSS database.Note
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsNote1UUID"]).First(&apiNote1JS), "finding api js note 1")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliNote1CSS.UUID).First(&apiNote1CSS), "finding api css note 1")
			var apiBookJS, apiBookCSS database.Book
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, ids["jsBookUUID"]).First(&apiBookJS), "finding api js book")
			apitest.MustExec(t, apiDB.Where("user_id = ? AND uuid = ?", user.ID, cliBookCSS.UUID).First(&apiBookCSS), "finding api css book")

			// assert on usn
			assert.NotEqual(t, apiNote1JS.USN, 0, "apiNote1JS usn mismatch")
			assert.NotEqual(t, apiNote1CSS.USN, 0, "apiNote1CSS usn mismatch")
			assert.NotEqual(t, apiBookJS.USN, 0, "apiBookJS usn mismatch")
			assert.NotEqual(t, apiBookCSS.USN, 0, "apiBookCSS usn mismatch")
			// assert on bodys and labels
			assert.Equal(t, apiNote1JS.Body, "js1", "apiNote1JS Body mismatch")
			assert.Equal(t, apiNote1CSS.Body, "css1", "apiNote1CSS Body mismatch")
			assert.Equal(t, apiBookJS.Label, "js", "apiBookJS Label mismatch")
			assert.Equal(t, apiBookCSS.Label, "css", "apiBookCSS Label mismatch")
			// assert on deleted
			assert.Equal(t, apiNote1JS.Deleted, false, "apiNote1JS Deleted mismatch")
			assert.Equal(t, apiNote1CSS.Deleted, false, "apiNote1CSS Deleted mismatch")
			assert.Equal(t, apiBookJS.Deleted, false, "apiBookJS Deleted mismatch")
			assert.Equal(t, apiBookCSS.Deleted, false, "apiBookCSS Deleted mismatch")
		}

		t.Run("stepSync then fullSync", func(t *testing.T) {
			// clean up
			os.RemoveAll(tmpDirPath)
			apitest.ClearData(apitest.DB)
			defer apitest.ClearData(apitest.DB)

			ctx := context.InitTestCtx(t, paths, nil)
			defer context.TeardownTestCtx(t, ctx)
			user := setupUser(t, &ctx)
			ids := setup(t, ctx, user)

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			assert(t, ctx, user, ids)
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync", "-f")
			assert(t, ctx, user, ids)
		})

		t.Run("fullSync then stepSync", func(t *testing.T) {
			// clean up
			os.RemoveAll(tmpDirPath)
			apitest.ClearData(apitest.DB)
			defer apitest.ClearData(apitest.DB)

			ctx := context.InitTestCtx(t, paths, nil)
			defer context.TeardownTestCtx(t, ctx)

			user := setupUser(t, &ctx)
			ids := setup(t, ctx, user)

			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync", "-f")
			assert(t, ctx, user, ids)
			clitest.RunDnoteCmd(t, dnoteCmdOpts, cliBinaryName, "sync")
			assert(t, ctx, user, ids)
		})
	})
}
