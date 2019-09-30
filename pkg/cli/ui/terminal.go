/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/dnote/dnote/pkg/cli/log"
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
