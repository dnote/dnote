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
