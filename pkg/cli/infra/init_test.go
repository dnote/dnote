/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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

package infra

import (
	"fmt"
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/config"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/pkg/errors"
)

func TestInitSystemKV(t *testing.T) {
	// Setup
	db := database.InitTestDB(t, "../tmp/dnote-test.db", nil)
	defer database.TeardownTestDB(t, db)

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
	db := database.InitTestDB(t, "../tmp/dnote-test.db", nil)
	defer database.TeardownTestDB(t, db)

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

func TestInit_APIEndpointChange(t *testing.T) {
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

	// First init.
	endpoint1 := "http://127.0.0.1:3001"
	ctx, err := Init("test-version", endpoint1)
	if err != nil {
		t.Fatal(errors.Wrap(err, "initializing"))
	}
	defer ctx.DB.Close()
	assert.Equal(t, ctx.APIEndpoint, endpoint1, "should use endpoint1 API endpoint")

	// Test that config was written with endpoint1.
	cf, err := config.Read(*ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "reading config"))
	}

	// Second init with different endpoint.
	endpoint2 := "http://127.0.0.1:3002"
	ctx2, err := Init("test-version", endpoint2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "initializing with override"))
	}
	defer ctx2.DB.Close()
	// Context must be using that endpoint.
	assert.Equal(t, ctx2.APIEndpoint, endpoint2, "should use endpoint2 API endpoint")

	// The config file shouldn't have been modified.
	cf2, err := config.Read(*ctx2)
	if err != nil {
		t.Fatal(errors.Wrap(err, "reading config after override"))
	}
	assert.Equal(t, cf2.APIEndpoint, cf.APIEndpoint, "config should still have original endpoint, not endpoint2")
}
