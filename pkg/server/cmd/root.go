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
