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
	"flag"
	"fmt"
	"os"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Helper functions for user commands

// setupUserFlags creates a FlagSet with standard usage format
func setupUserFlags(name, usageCmd string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf(`Usage:
  %s [flags]

Flags:
`, usageCmd)
		fs.PrintDefaults()
	}
	return fs
}

// requireString validates that a required string flag is not empty
func requireString(fs *flag.FlagSet, value, fieldName string) {
	if value == "" {
		fmt.Printf("Error: %s is required\n", fieldName)
		fs.Usage()
		os.Exit(1)
	}
}

// setupAppWithDB creates config, initializes app, and returns cleanup function
func setupAppWithDB(fs *flag.FlagSet, dbPath string) (*app.App, func()) {
	cfg, err := config.New(config.Params{
		DBPath: dbPath,
	})
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		fs.Usage()
		os.Exit(1)
	}

	a := initApp(cfg)
	cleanup := func() {
		sqlDB, err := a.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	return &a, cleanup
}


func userCreateCmd(args []string) {
	fs := setupUserFlags("create", "dnote-server user create")

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

func userRemoveCmd(args []string) {
	fs := setupUserFlags("remove", "dnote-server user remove")

	email := fs.String("email", "", "User email address (required)")
	dbPath := fs.String("dbPath", "", "Path to SQLite database file (env: DBPath, default: $XDG_DATA_HOME/dnote/server.db)")

	fs.Parse(args)

	requireString(fs, *email, "email")

	a, cleanup := setupAppWithDB(fs, *dbPath)
	defer cleanup()

	// Find the account and user
	account, err := a.GetAccountByEmail(*email)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			fmt.Printf("Error: user with email %s not found\n", *email)
		} else {
			log.ErrorWrap(err, "finding account")
		}
		os.Exit(1)
	}

	// Check if user has any notes
	var noteCount int64
	if err := a.DB.Model(&database.Note{}).Where("user_id = ? AND deleted = ?", account.UserID, false).Count(&noteCount).Error; err != nil {
		log.ErrorWrap(err, "counting notes")
		os.Exit(1)
	}
	if noteCount > 0 {
		fmt.Printf("Error: cannot remove user with existing notes (found %d notes)\n", noteCount)
		os.Exit(1)
	}

	// Check if user has any books
	var bookCount int64
	if err := a.DB.Model(&database.Book{}).Where("user_id = ? AND deleted = ?", account.UserID, false).Count(&bookCount).Error; err != nil {
		log.ErrorWrap(err, "counting books")
		os.Exit(1)
	}
	if bookCount > 0 {
		fmt.Printf("Error: cannot remove user with existing books (found %d books)\n", bookCount)
		os.Exit(1)
	}

	// Delete account and user
	tx := a.DB.Begin()
	if err := tx.Delete(&account).Error; err != nil {
		tx.Rollback()
		log.ErrorWrap(err, "deleting account")
		os.Exit(1)
	}

	var user database.User
	if err := tx.Where("id = ?", account.UserID).First(&user).Error; err != nil {
		tx.Rollback()
		log.ErrorWrap(err, "finding user")
		os.Exit(1)
	}

	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		log.ErrorWrap(err, "deleting user")
		os.Exit(1)
	}

	tx.Commit()

	fmt.Printf("User removed successfully\n")
	fmt.Printf("Email: %s\n", *email)
}

func userResetPasswordCmd(args []string) {
	fs := setupUserFlags("reset-password", "dnote-server user reset-password")

	email := fs.String("email", "", "User email address (required)")
	password := fs.String("password", "", "New password (required)")
	dbPath := fs.String("dbPath", "", "Path to SQLite database file (env: DBPath, default: $XDG_DATA_HOME/dnote/server.db)")

	fs.Parse(args)

	requireString(fs, *email, "email")
	requireString(fs, *password, "password")

	a, cleanup := setupAppWithDB(fs, *dbPath)
	defer cleanup()

	// Validate password
	if err := app.ValidatePassword(*password); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// Find the account
	account, err := a.GetAccountByEmail(*email)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			fmt.Printf("Error: user with email %s not found\n", *email)
		} else {
			log.ErrorWrap(err, "finding account")
		}
		os.Exit(1)
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.ErrorWrap(err, "hashing password")
		os.Exit(1)
	}

	// Update the password
	if err := a.DB.Model(&account).Update("password", string(hashedPassword)).Error; err != nil {
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
  remove: Remove a user (only if they have no notes or books)
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
		userRemoveCmd(subArgs)
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
