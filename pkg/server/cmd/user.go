/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/dnote/dnote/pkg/prompt"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/pkg/errors"
)

// confirm prompts for user input to confirm a choice
func confirm(r io.Reader, question string, optimistic bool) (bool, error) {
	message := prompt.FormatQuestion(question, optimistic)
	fmt.Print(message + " ")

	confirmed, err := prompt.ReadYesNo(r, optimistic)
	if err != nil {
		return false, errors.Wrap(err, "reading stdin")
	}

	return confirmed, nil
}

func userCreateCmd(args []string) {
	fs := setupFlagSet("create", "dnote-server user create")

	email := fs.String("email", "", "User email address (required)")
	password := fs.String("password", "", "User password (required)")
	dbPath := fs.String("dbPath", "", "Path to SQLite database file (env: DBPath, default: $XDG_DATA_HOME/dnote/server.db)")

	fs.Parse(args)

	requireString(fs, *email, "email")
	requireString(fs, *password, "password")

	a, cleanup := setupAppWithDB(fs, *dbPath)
	defer cleanup()

	_, err := a.CreateUser(*email, *password, *password)
	if err != nil {
		log.ErrorWrap(err, "creating user")
		os.Exit(1)
	}

	fmt.Printf("User created successfully\n")
	fmt.Printf("Email: %s\n", *email)
}

func userRemoveCmd(args []string, stdin io.Reader) {
	fs := setupFlagSet("remove", "dnote-server user remove")

	email := fs.String("email", "", "User email address (required)")
	dbPath := fs.String("dbPath", "", "Path to SQLite database file (env: DBPath, default: $XDG_DATA_HOME/dnote/server.db)")

	fs.Parse(args)

	requireString(fs, *email, "email")

	a, cleanup := setupAppWithDB(fs, *dbPath)
	defer cleanup()

	// Check if user exists first
	_, err := a.GetUserByEmail(*email)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			fmt.Printf("Error: user with email %s not found\n", *email)
		} else {
			log.ErrorWrap(err, "finding user")
		}
		os.Exit(1)
	}

	// Show confirmation prompt
	ok, err := confirm(stdin, fmt.Sprintf("Remove user %s?", *email), false)
	if err != nil {
		log.ErrorWrap(err, "getting confirmation")
		os.Exit(1)
	}
	if !ok {
		fmt.Println("Aborted by user")
		os.Exit(0)
	}

	// Remove the user
	if err := a.RemoveUser(*email); err != nil {
		if errors.Is(err, app.ErrNotFound) {
			fmt.Printf("Error: user with email %s not found\n", *email)
		} else if errors.Is(err, app.ErrUserHasExistingResources) {
			fmt.Printf("Error: %s\n", err)
		} else {
			log.ErrorWrap(err, "removing user")
		}
		os.Exit(1)
	}

	fmt.Printf("User removed successfully\n")
	fmt.Printf("Email: %s\n", *email)
}

func userResetPasswordCmd(args []string) {
	fs := setupFlagSet("reset-password", "dnote-server user reset-password")

	email := fs.String("email", "", "User email address (required)")
	password := fs.String("password", "", "New password (required)")
	dbPath := fs.String("dbPath", "", "Path to SQLite database file (env: DBPath, default: $XDG_DATA_HOME/dnote/server.db)")

	fs.Parse(args)

	requireString(fs, *email, "email")
	requireString(fs, *password, "password")

	a, cleanup := setupAppWithDB(fs, *dbPath)
	defer cleanup()

	// Find the user
	user, err := a.GetUserByEmail(*email)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			fmt.Printf("Error: user with email %s not found\n", *email)
		} else {
			log.ErrorWrap(err, "finding user")
		}
		os.Exit(1)
	}

	// Update the password
	if err := app.UpdateUserPassword(a.DB, user, *password); err != nil {
		log.ErrorWrap(err, "updating password")
		os.Exit(1)
	}

	fmt.Printf("Password reset successfully\n")
	fmt.Printf("Email: %s\n", *email)
}

func userCmd(args []string) {
	if len(args) < 1 {
		fmt.Println(`Usage:
  dnote-server user [command]

Available commands:
  create: Create a new user
  remove: Remove a user
  reset-password: Reset a user's password`)
		os.Exit(1)
	}

	subcommand := args[0]
	subArgs := []string{}
	if len(args) > 1 {
		subArgs = args[1:]
	}

	switch subcommand {
	case "create":
		userCreateCmd(subArgs)
	case "remove":
		userRemoveCmd(subArgs, os.Stdin)
	case "reset-password":
		userResetPasswordCmd(subArgs)
	default:
		fmt.Printf("Unknown subcommand: %s\n\n", subcommand)
		fmt.Println(`Available commands:
  create: Create a new user
  remove: Remove a user (only if they have no notes or books)
  reset-password: Reset a user's password`)
		os.Exit(1)
	}
}
