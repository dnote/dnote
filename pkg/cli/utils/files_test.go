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

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test", "nested", "dir")

	// Create directory
	err := EnsureDir(testPath)
	assert.Equal(t, err, nil, "EnsureDir should succeed")

	// Verify it exists
	info, err := os.Stat(testPath)
	assert.Equal(t, err, nil, "directory should exist")
	assert.Equal(t, info.IsDir(), true, "should be a directory")

	// Call again on existing directory - should not error
	err = EnsureDir(testPath)
	assert.Equal(t, err, nil, "EnsureDir should succeed on existing directory")
}
