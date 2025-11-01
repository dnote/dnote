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

package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/prompt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

func readInput() (string, error) {
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

	input, err := readInput()
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

// Confirm prompts for user input to confirm a choice
func Confirm(question string, optimistic bool) (bool, error) {
	message := prompt.FormatQuestion(question, optimistic)

	// Use log.Askf for colored prompt in CLI
	log.Askf(message, false)

	confirmed, err := prompt.ReadYesNo(os.Stdin, optimistic)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get user input")
	}

	return confirmed, nil
}

// Grab text from stdin content
func ReadStdInput() (string, error) {
	var lines []string

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	err := s.Err()
	if err != nil {
		return "", errors.Wrap(err, "reading pipe")
	}

	return strings.Join(lines, "\n"), nil
}
