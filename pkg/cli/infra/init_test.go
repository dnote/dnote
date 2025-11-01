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

package infra

import (
	"fmt"
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/config"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/dirs"
	"github.com/pkg/errors"
)

func TestInitSystemKV(t *testing.T) {
	// Setup
	db := database.InitTestMemoryDB(t)

	var originalCount int
	database.MustScan(t, "counting system configs", db.QueryRow("SELECT count(*) FROM system"), &originalCount)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	if err := initSystemKV(tx, "testKey", "testVal"); err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "executing"))
	}

	tx.Commit()

	// Test
	var count int
	database.MustScan(t, "counting system configs", db.QueryRow("SELECT count(*) FROM system"), &count)
	assert.Equal(t, count, originalCount+1, "system count mismatch")

	var val string
	database.MustScan(t, "getting system value",
		db.QueryRow("SELECT value FROM system WHERE key = ?", "testKey"), &val)
	assert.Equal(t, val, "testVal", "system value mismatch")
}

func TestInitSystemKV_existing(t *testing.T) {
	// Setup
	db := database.InitTestMemoryDB(t)

	database.MustExec(t, "inserting a system config", db, "INSERT INTO system (key, value) VALUES (?, ?)", "testKey", "testVal")

	var originalCount int
	database.MustScan(t, "counting system configs", db.QueryRow("SELECT count(*) FROM system"), &originalCount)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	if err := initSystemKV(tx, "testKey", "newTestVal"); err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "executing"))
	}

	tx.Commit()

	// Test
	var count int
	database.MustScan(t, "counting system configs", db.QueryRow("SELECT count(*) FROM system"), &count)
	assert.Equal(t, count, originalCount, "system count mismatch")

	var val string
	database.MustScan(t, "getting system value",
		db.QueryRow("SELECT value FROM system WHERE key = ?", "testKey"), &val)
	assert.Equal(t, val, "testVal", "system value should not have been updated")
}

func TestInit_APIEndpoint(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "dnote-init-test-*")
	if err != nil {
		t.Fatal(errors.Wrap(err, "creating temp dir"))
	}
	defer os.RemoveAll(tmpDir)

	// Set up environment to use our temp directory
	t.Setenv("XDG_CONFIG_HOME", fmt.Sprintf("%s/config", tmpDir))
	t.Setenv("XDG_DATA_HOME", fmt.Sprintf("%s/data", tmpDir))
	t.Setenv("XDG_CACHE_HOME", fmt.Sprintf("%s/cache", tmpDir))

	// Force dirs package to reload with new environment
	dirs.Reload()

	// Initialize - should create config with default apiEndpoint
	ctx, err := Init("test-version", "", "")
	if err != nil {
		t.Fatal(errors.Wrap(err, "initializing"))
	}
	defer ctx.DB.Close()

	// Read the config that was created
	cf, err := config.Read(*ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "reading config"))
	}

	// Context should use the apiEndpoint from config
	assert.Equal(t, ctx.APIEndpoint, DefaultAPIEndpoint, "context should use apiEndpoint from config")
	assert.Equal(t, cf.APIEndpoint, DefaultAPIEndpoint, "context should use apiEndpoint from config")
}
