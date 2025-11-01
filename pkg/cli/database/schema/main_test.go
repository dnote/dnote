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
