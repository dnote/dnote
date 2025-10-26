/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
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

package sync

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/pkg/errors"
)

var cliBinaryName string
var serverTime = time.Date(2017, time.March, 14, 21, 15, 0, 0, time.UTC)

var testDir = "./tmp/"

func init() {
	cliBinaryName = fmt.Sprintf("%s/test-cli", testDir)
}

func TestMain(m *testing.M) {
	// Build CLI binary without hardcoded API endpoint
	// Each test will create its own server and config file
	cmd := exec.Command("go", "build", "--tags", "fts5", "-o", cliBinaryName, "github.com/dnote/dnote/pkg/cli")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Print(errors.Wrap(err, "building a CLI binary").Error())
		log.Print(stderr.String())
		os.Exit(1)
	}

	os.Exit(m.Run())
}
