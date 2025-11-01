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

//go:build linux || darwin || freebsd


package dirs

import (
	"path/filepath"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestDirs(t *testing.T) {
	home := Home
	assert.NotEqual(t, home, "", "home is empty")

	configHome := filepath.Join(home, ".config")
	dataHome := filepath.Join(home, ".local", "share")
	cacheHome := filepath.Join(home, ".cache")

	testCases := []struct {
		got      string
		expected string
	}{
		{
			got:      ConfigHome,
			expected: configHome,
		},
		{
			got:      DataHome,
			expected: dataHome,
		},
		{
			got:      CacheHome,
			expected: cacheHome,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.got, tc.expected, "result mismatch")
	}
}

func TestCustomDirs(t *testing.T) {
	testCases := []envTestCase{
		{
			envKey:   "XDG_CONFIG_HOME",
			envVal:   "~/custom/config",
			got:      &ConfigHome,
			expected: "~/custom/config",
		},
		{
			envKey:   "XDG_DATA_HOME",
			envVal:   "~/custom/data",
			got:      &DataHome,
			expected: "~/custom/data",
		},
		{
			envKey:   "XDG_CACHE_HOME",
			envVal:   "~/custom/cache",
			got:      &CacheHome,
			expected: "~/custom/cache",
		},
	}

	testCustomDirs(t, testCases)
}
