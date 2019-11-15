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

package infra

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/pkg/errors"
)

func TestInitSystemKV(t *testing.T) {
	// Setup
	db := database.InitTestDB(t, "../tmp/dnote-test.db", nil)
	defer database.CloseTestDB(t, db)

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
	defer database.CloseTestDB(t, db)

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
