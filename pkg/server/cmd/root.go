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
	"os"
)

func rootCmd() {
	fmt.Printf(`Dnote server - a simple command line notebook

Usage:
  dnote-server [command] [flags]

Available commands:
  start: Start the server (use 'dnote-server start --help' for flags)
  user: Manage users (use 'dnote-server user' for subcommands)
  version: Print the version
`)
}

// Execute is the main entry point for the CLI
func Execute() {
	if len(os.Args) < 2 {
		rootCmd()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "start":
		startCmd(os.Args[2:])
	case "user":
		userCmd(os.Args[2:])
	case "version":
		versionCmd()
	default:
		fmt.Printf("Unknown command %s\n", cmd)
		rootCmd()
		os.Exit(1)
	}
}
