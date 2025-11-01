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

// Package testutils provides utilities used in tests
package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
)

// Prompts for user input
const (
	PromptRemoveNote  = "remove this note?"
	PromptDeleteBook  = "delete book"
	PromptEmptyServer = "The server is empty but you have local data"
)

// Timeout for waiting for prompts in tests
const promptTimeout = 10 * time.Second

// LoginDB sets up login credentials in the database for tests
func LoginDB(t *testing.T, db *database.DB) {
	database.MustExec(t, "inserting sessionKey", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKey, "someSessionKey")
	database.MustExec(t, "inserting sessionKeyExpiry", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKeyExpiry, time.Now().Add(24*time.Hour).Unix())
}

// Login simulates a logged in user by inserting credentials in the local database
func Login(t *testing.T, ctx *context.DnoteCtx) {
	LoginDB(t, ctx.DB)

	ctx.SessionKey = "someSessionKey"
	ctx.SessionKeyExpiry = time.Now().Add(24 * time.Hour).Unix()
}

// RemoveDir cleans up the test env represented by the given context
func RemoveDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(errors.Wrap(err, "removing the directory"))
	}
}

// CopyFixture writes the content of the given fixture to the filename inside the dnote dir
func CopyFixture(t *testing.T, ctx context.DnoteCtx, fixturePath string, filename string) {
	fp, err := filepath.Abs(fixturePath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "getting the absolute path for fixture"))
	}

	dp, err := filepath.Abs(filepath.Join(ctx.Paths.LegacyDnote, filename))
	if err != nil {
		t.Fatal(errors.Wrap(err, "getting the absolute path dnote dir"))
	}

	err = utils.CopyFile(fp, dp)
	if err != nil {
		t.Fatal(errors.Wrap(err, "copying the file"))
	}
}

// WriteFile writes a file with the given content and  filename inside the dnote dir
func WriteFile(ctx context.DnoteCtx, content []byte, filename string) {
	dp, err := filepath.Abs(filepath.Join(ctx.Paths.LegacyDnote, filename))
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(dp, content, 0644); err != nil {
		panic(err)
	}
}

// ReadFile reads the content of the file with the given name in dnote dir
func ReadFile(ctx context.DnoteCtx, filename string) []byte {
	path := filepath.Join(ctx.Paths.LegacyDnote, filename)

	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return b
}

// ReadJSON reads JSON fixture to the struct at the destination address
func ReadJSON(path string, destination interface{}) {
	var dat []byte
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(errors.Wrap(err, "Failed to load fixture payload"))
	}
	if err := json.Unmarshal(dat, destination); err != nil {
		panic(errors.Wrap(err, "Failed to get event"))
	}
}

// NewDnoteCmd returns a new Dnote command and a pointer to stderr
func NewDnoteCmd(opts RunDnoteCmdOptions, binaryName string, arg ...string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer, error) {
	var stderr, stdout bytes.Buffer

	binaryPath, err := filepath.Abs(binaryName)
	if err != nil {
		return &exec.Cmd{}, &stderr, &stdout, errors.Wrap(err, "getting the absolute path to the test binary")
	}

	cmd := exec.Command(binaryPath, arg...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	cmd.Env = opts.Env

	return cmd, &stderr, &stdout, nil
}

// RunDnoteCmdOptions is an option for RunDnoteCmd
type RunDnoteCmdOptions struct {
	Env []string
}

// RunDnoteCmd runs a dnote command
func RunDnoteCmd(t *testing.T, opts RunDnoteCmdOptions, binaryName string, arg ...string) {
	t.Logf("running: %s %s", binaryName, strings.Join(arg, " "))

	cmd, stderr, stdout, err := NewDnoteCmd(opts, binaryName, arg...)
	if err != nil {
		t.Logf("\n%s", stdout)
		t.Fatal(errors.Wrap(err, "getting command").Error())
	}

	cmd.Env = append(cmd.Env, "DNOTE_DEBUG=1")

	if err := cmd.Run(); err != nil {
		t.Logf("\n%s", stdout)
		t.Fatal(errors.Wrapf(err, "running command %s", stderr.String()))
	}

	// Print stdout if and only if test fails later
	t.Logf("\n%s", stdout)
}

// WaitDnoteCmd runs a dnote command and passes stdout to the callback.
func WaitDnoteCmd(t *testing.T, opts RunDnoteCmdOptions, runFunc func(io.Reader, io.WriteCloser) error, binaryName string, arg ...string) (string, error) {
	t.Logf("running: %s %s", binaryName, strings.Join(arg, " "))

	binaryPath, err := filepath.Abs(binaryName)
	if err != nil {
		return "", errors.Wrap(err, "getting absolute path to test binary")
	}

	cmd := exec.Command(binaryPath, arg...)
	cmd.Env = opts.Env

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, "getting stdout pipe")
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", errors.Wrap(err, "getting stdin")
	}
	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		return "", errors.Wrap(err, "starting command")
	}

	var output bytes.Buffer
	tee := io.TeeReader(stdout, &output)

	err = runFunc(tee, stdin)
	if err != nil {
		t.Logf("\n%s", output.String())
		return output.String(), errors.Wrap(err, "running callback")
	}

	io.Copy(&output, stdout)

	if err := cmd.Wait(); err != nil {
		t.Logf("\n%s", output.String())
		return output.String(), errors.Wrapf(err, "command failed: %s", stderr.String())
	}

	t.Logf("\n%s", output.String())
	return output.String(), nil
}

func MustWaitDnoteCmd(t *testing.T, opts RunDnoteCmdOptions, runFunc func(io.Reader, io.WriteCloser) error, binaryName string, arg ...string) string {
	output, err := WaitDnoteCmd(t, opts, runFunc, binaryName, arg...)
	if err != nil {
		t.Fatal(err)
	}

	return output
}

// MustWaitForPrompt waits for an expected prompt with a default timeout.
// Fails the test if the prompt is not found or an error occurs.
func MustWaitForPrompt(t *testing.T, stdout io.Reader, expectedPrompt string) {
	if err := assert.WaitForPrompt(stdout, expectedPrompt, promptTimeout); err != nil {
		t.Fatal(err)
	}
}

// ConfirmRemoveNote waits for prompt for removing a note and confirms.
func ConfirmRemoveNote(stdout io.Reader, stdin io.WriteCloser) error {
	return assert.RespondToPrompt(stdout, stdin, PromptRemoveNote, "y\n", promptTimeout)
}

// ConfirmRemoveBook waits for prompt for deleting a book confirms.
func ConfirmRemoveBook(stdout io.Reader, stdin io.WriteCloser) error {
	return assert.RespondToPrompt(stdout, stdin, PromptDeleteBook, "y\n", promptTimeout)
}

// UserConfirmEmptyServerSync waits for an empty server prompt and confirms.
func UserConfirmEmptyServerSync(stdout io.Reader, stdin io.WriteCloser) error {
	return assert.RespondToPrompt(stdout, stdin, PromptEmptyServer, "y\n", promptTimeout)
}

// UserCancelEmptyServerSync waits for an empty server prompt and cancels.
func UserCancelEmptyServerSync(stdout io.Reader, stdin io.WriteCloser) error {
	return assert.RespondToPrompt(stdout, stdin, PromptEmptyServer, "n\n", promptTimeout)
}

// UserContent simulates content from the user by writing to stdin.
// This is used for piped input where no prompt is shown.
func UserContent(stdout io.Reader, stdin io.WriteCloser) error {
	longText := `Lorem ipsum dolor sit amet, consectetur adipiscing elit,
	sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`

	if _, err := io.WriteString(stdin, longText); err != nil {
		return errors.Wrap(err, "creating note from stdin")
	}

	// stdin needs to close so stdin reader knows to stop reading
	// otherwise test case would wait until test timeout
	stdin.Close()

	return nil
}

// MustMarshalJSON marshalls the given interface into JSON.
// If there is any error, it fails the test.
func MustMarshalJSON(t *testing.T, v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("%s: marshalling data: %s", t.Name(), err.Error())
	}

	return b
}

// MustUnmarshalJSON marshalls the given interface into JSON.
// If there is any error, it fails the test.
func MustUnmarshalJSON(t *testing.T, data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		t.Fatalf("%s: unmarshalling data: %s", t.Name(), err.Error())
	}
}

// MustGenerateUUID generates the uuid. If error occurs, it fails the test.
func MustGenerateUUID(t *testing.T) string {
	ret, err := utils.GenerateUUID()
	if err != nil {
		t.Fatal(errors.Wrap(err, "generating uuid").Error())
	}

	return ret
}

func MustOpenDatabase(t *testing.T, dbPath string) *database.DB {
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "opening database"))
	}

	return db
}
