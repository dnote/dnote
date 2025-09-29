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

//go:build linux || darwin


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
