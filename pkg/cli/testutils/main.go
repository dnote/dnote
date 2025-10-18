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

// Package testutils provides utilities used in tests
package testutils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

// Login simulates a logged in user by inserting credentials in the local database
func Login(t *testing.T, ctx *context.DnoteCtx) {
	db := ctx.DB

	database.MustExec(t, "inserting sessionKey", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKey, "someSessionKey")
	database.MustExec(t, "inserting sessionKeyExpiry", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKeyExpiry, time.Now().Add(24*time.Hour).Unix())

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

// waitForPrompt waits for an expected prompt to appear in stdout with a timeout.
// Returns an error if the prompt is not found within the timeout period.
// Handles prompts with or without newlines by reading character by character.
func waitForPrompt(stdout io.Reader, expectedPrompt string, timeout time.Duration) error {
	type result struct {
		found bool
		err   error
	}
	resultCh := make(chan result, 1)

	go func() {
		reader := bufio.NewReader(stdout)
		var buffer strings.Builder
		found := false

		for {
			b, err := reader.ReadByte()
			if err != nil {
				resultCh <- result{found: found, err: err}
				return
			}

			buffer.WriteByte(b)
			if strings.Contains(buffer.String(), expectedPrompt) {
				found = true
				break
			}
		}

		resultCh <- result{found: found, err: nil}
	}()

	select {
	case res := <-resultCh:
		if res.err != nil && res.err != io.EOF {
			return errors.Wrap(res.err, "reading stdout")
		}
		if !res.found {
			return errors.Errorf("expected prompt '%s' not found in stdout", expectedPrompt)
		}
		return nil
	case <-time.After(timeout):
		return errors.Errorf("timeout waiting for prompt '%s'", expectedPrompt)
	}
}

// MustWaitForPrompt waits for an expected prompt with a default timeout.
// Fails the test if the prompt is not found or an error occurs.
func MustWaitForPrompt(t *testing.T, stdout io.Reader, expectedPrompt string) {
	if err := waitForPrompt(stdout, expectedPrompt, promptTimeout); err != nil {
		t.Fatal(err)
	}
}

// userRespondToPrompt is a helper that waits for a prompt and sends a response.
func userRespondToPrompt(stdout io.Reader, stdin io.WriteCloser, expectedPrompt, response, action string) error {
	if err := waitForPrompt(stdout, expectedPrompt, promptTimeout); err != nil {
		return err
	}

	if _, err := io.WriteString(stdin, response); err != nil {
		return errors.Wrapf(err, "indicating %s in stdin", action)
	}

	return nil
}

// userConfirmOutput simulates confirmation from the user by writing to stdin.
// It waits for the expected prompt with a timeout to prevent deadlocks.
func userConfirmOutput(stdout io.Reader, stdin io.WriteCloser, expectedPrompt string) error {
	return userRespondToPrompt(stdout, stdin, expectedPrompt, "y\n", "confirmation")
}

// userCancelOutput simulates cancellation from the user by writing to stdin.
// It waits for the expected prompt with a timeout to prevent deadlocks.
func userCancelOutput(stdout io.Reader, stdin io.WriteCloser, expectedPrompt string) error {
	return userRespondToPrompt(stdout, stdin, expectedPrompt, "n\n", "cancellation")
}

// ConfirmRemoveNote waits for prompt for removing a note and confirms.
func ConfirmRemoveNote(stdout io.Reader, stdin io.WriteCloser) error {
	return userConfirmOutput(stdout, stdin, PromptRemoveNote)
}

// ConfirmRemoveBook waits for prompt for deleting a book confirms.
func ConfirmRemoveBook(stdout io.Reader, stdin io.WriteCloser) error {
	return userConfirmOutput(stdout, stdin, PromptDeleteBook)
}

// UserConfirmEmptyServerSync waits for an empty server prompt and confirms.
func UserConfirmEmptyServerSync(stdout io.Reader, stdin io.WriteCloser) error {
	return userConfirmOutput(stdout, stdin, PromptEmptyServer)
}

// UserCancelEmptyServerSync waits for an empty server prompt and confirms.
func UserCancelEmptyServerSync(stdout io.Reader, stdin io.WriteCloser) error {
	return userCancelOutput(stdout, stdin, PromptEmptyServer)
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
