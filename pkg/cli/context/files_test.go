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

package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
)

func assertDirsExist(t *testing.T, paths Paths) {
	configDir := filepath.Join(paths.Config, consts.DnoteDirName)
	info, err := os.Stat(configDir)
	assert.Equal(t, err, nil, "config dir should exist")
	assert.Equal(t, info.IsDir(), true, "config should be a directory")

	dataDir := filepath.Join(paths.Data, consts.DnoteDirName)
	info, err = os.Stat(dataDir)
	assert.Equal(t, err, nil, "data dir should exist")
	assert.Equal(t, info.IsDir(), true, "data should be a directory")

	cacheDir := filepath.Join(paths.Cache, consts.DnoteDirName)
	info, err = os.Stat(cacheDir)
	assert.Equal(t, err, nil, "cache dir should exist")
	assert.Equal(t, info.IsDir(), true, "cache should be a directory")
}

func TestInitDnoteDirs(t *testing.T) {
	tmpDir := t.TempDir()

	paths := Paths{
		Config: filepath.Join(tmpDir, "config"),
		Data:   filepath.Join(tmpDir, "data"),
		Cache:  filepath.Join(tmpDir, "cache"),
	}

	// Initialize directories
	err := InitDnoteDirs(paths)
	assert.Equal(t, err, nil, "InitDnoteDirs should succeed")
	assertDirsExist(t, paths)

	// Call again - should be idempotent
	err = InitDnoteDirs(paths)
	assert.Equal(t, err, nil, "InitDnoteDirs should succeed when dirs already exist")
	assertDirsExist(t, paths)
}
