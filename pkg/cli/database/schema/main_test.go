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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
)

func TestRun(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "schema.sql")

	// Run the function
	if err := run(tmpDir, outputPath); err != nil {
		t.Fatalf("run() failed: %v", err)
	}

	// Verify schema.sql was created
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("reading schema.sql: %v", err)
	}

	schema := string(content)

	// Verify it has the header
	assert.Equal(t, strings.HasPrefix(schema, "-- This is the final state"), true, "schema.sql should have header comment")

	// Verify schema contains expected tables
	expectedTables := []string{
		"CREATE TABLE books",
		"CREATE TABLE system",
		"CREATE TABLE \"notes\"",
		"CREATE VIRTUAL TABLE note_fts",
	}

	for _, expected := range expectedTables {
		assert.Equal(t, strings.Contains(schema, expected), true, fmt.Sprintf("schema should contain %s", expected))
	}

	// Verify schema contains triggers
	expectedTriggers := []string{
		"CREATE TRIGGER notes_after_insert",
		"CREATE TRIGGER notes_after_delete",
		"CREATE TRIGGER notes_after_update",
	}

	for _, expected := range expectedTriggers {
		assert.Equal(t, strings.Contains(schema, expected), true, fmt.Sprintf("schema should contain %s", expected))
	}

	// Verify schema does not contain sqlite internal tables
	assert.Equal(t, strings.Contains(schema, "sqlite_sequence"), false, "schema should not contain sqlite_sequence")

	// Verify system key-value pairs for schema versions are present
	expectedSchemaKey := fmt.Sprintf("INSERT INTO system (key, value) VALUES ('%s',", consts.SystemSchema)
	assert.Equal(t, strings.Contains(schema, expectedSchemaKey), true, "schema should contain schema version INSERT statement")

	expectedRemoteSchemaKey := fmt.Sprintf("INSERT INTO system (key, value) VALUES ('%s',", consts.SystemRemoteSchema)
	assert.Equal(t, strings.Contains(schema, expectedRemoteSchemaKey), true, "schema should contain remote_schema version INSERT statement")
}
