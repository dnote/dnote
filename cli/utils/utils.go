/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/dnote/dnote/cli/infra"
	"github.com/dnote/dnote/cli/log"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh/terminal"
)

// GenerateUUID returns a uid
func GenerateUUID() string {
	return uuid.NewV4().String()
}

func getInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrap(err, "reading stdin")
	}

	return strings.Trim(input, "\r\n"), nil
}

// PromptInput prompts the user input and saves the result to the destination
func PromptInput(message string, dest *string) error {
	log.Askf(message, false)

	input, err := getInput()
	if err != nil {
		return errors.Wrap(err, "getting user input")
	}

	*dest = input

	return nil
}

// PromptPassword prompts the user input a password and saves the result to the destination.
// The input is masked, meaning it is not echoed on the terminal.
func PromptPassword(message string, dest *string) error {
	log.Askf(message, true)

	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return errors.Wrap(err, "getting user input")
	}

	fmt.Println("")

	*dest = string(password)

	return nil
}

// AskConfirmation prompts for user input to confirm a choice
func AskConfirmation(question string, optimistic bool) (bool, error) {
	var choices string
	if optimistic {
		choices = "(Y/n)"
	} else {
		choices = "(y/N)"
	}

	message := fmt.Sprintf("%s %s", question, choices)

	var input string
	if err := PromptInput(message, &input); err != nil {
		return false, errors.Wrap(err, "Failed to get user input")
	}

	confirmed := input == "y"

	if optimistic {
		confirmed = confirmed || input == ""
	}

	return confirmed, nil
}

// FileExists checks if the file exists at the given path
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// CopyDir copies a directory from src to dest, recursively copying nested
// directories
func CopyDir(src, dest string) error {
	srcPath := filepath.Clean(src)
	destPath := filepath.Clean(dest)

	fi, err := os.Stat(srcPath)
	if err != nil {
		return errors.Wrap(err, "getting the file info for the input")
	}
	if !fi.IsDir() {
		return errors.New("source is not a directory")
	}

	_, err = os.Stat(dest)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "looking up the destination")
	}

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return errors.Wrap(err, "creating destination")
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.Wrap(err, "reading the directory listing for the input")
	}

	for _, entry := range entries {
		srcEntryPath := filepath.Join(srcPath, entry.Name())
		destEntryPath := filepath.Join(destPath, entry.Name())

		if entry.IsDir() {
			if err = CopyDir(srcEntryPath, destEntryPath); err != nil {
				return errors.Wrapf(err, "copying %s", entry.Name())
			}
		} else {
			if err = CopyFile(srcEntryPath, destEntryPath); err != nil {
				return errors.Wrapf(err, "copying %s", entry.Name())
			}
		}
	}

	return nil
}

func getReq(ctx infra.DnoteCtx, path, method, body string) (*http.Request, error) {
	endpoint := fmt.Sprintf("%s%s", ctx.APIEndpoint, path)
	req, err := http.NewRequest(method, endpoint, strings.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "constructing http request")
	}

	req.Header.Set("CLI-Version", ctx.Version)

	return req, nil
}

// DoAuthorizedReq does a http request to the given path in the api endpoint as a user,
// with the appropriate headers. The given path should include the preceding slash.
func DoAuthorizedReq(ctx infra.DnoteCtx, hc http.Client, method, path, body string) (*http.Response, error) {
	if ctx.SessionKey == "" {
		return nil, errors.New("no session key found")
	}

	req, err := getReq(ctx, path, method, body)
	if err != nil {
		return nil, errors.Wrap(err, "getting request")
	}

	credential := fmt.Sprintf("Bearer %s", ctx.SessionKey)
	req.Header.Set("Authorization", credential)

	res, err := hc.Do(req)
	if err != nil {
		return res, errors.Wrap(err, "making http request")
	}

	return res, nil
}

// DoReq does a http request to the given path in the api endpoint
func DoReq(ctx infra.DnoteCtx, method, path, body string) (*http.Response, error) {
	req, err := getReq(ctx, path, method, body)
	if err != nil {
		return nil, errors.Wrap(err, "getting request")
	}

	hc := http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return res, errors.Wrap(err, "making http request")
	}

	return res, nil
}

// CopyFile copies a file from the src to dest
func CopyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "opening the input file")
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return errors.Wrap(err, "creating the output file")
	}

	if _, err = io.Copy(out, in); err != nil {
		return errors.Wrap(err, "copying the file content")
	}

	if err = out.Sync(); err != nil {
		return errors.Wrap(err, "flushing the output file to disk")
	}

	fi, err := os.Stat(src)
	if err != nil {
		return errors.Wrap(err, "getting the file info for the input file")
	}

	if err = os.Chmod(dest, fi.Mode()); err != nil {
		return errors.Wrap(err, "copying permission to the output file")
	}

	// Close the output file
	if err = out.Close(); err != nil {
		return errors.Wrap(err, "closing the output file")
	}

	return nil
}
