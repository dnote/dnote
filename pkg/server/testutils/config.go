/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package testutils

import (
	"os"
	"path/filepath"
)

// ProjectPath is the path of the proprietary test suite relative to the "GOPATH"
var ProjectPath string

// CLIPath is the path to the CLI project
var CLIPath string

// ServerPath is the path to the Dnote server project
var ServerPath string

func init() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		panic("GOPATH is not set up")
	}

	// Populate paths
	ProjectPath = filepath.Join(goPath, "src/gitlab.com/monomax/dnote-infra")
	CLIPath = filepath.Join(goPath, "src/github.com/dnote/dnote/pkg/cli")
	ServerPath = filepath.Join(goPath, "src/github.com/dnote/dnote/pkg/server")
}
