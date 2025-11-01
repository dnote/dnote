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
