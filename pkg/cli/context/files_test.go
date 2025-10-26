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
