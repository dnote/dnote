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

package dirs

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

type envTestCase struct {
	envKey   string
	envVal   string
	got      *string
	expected string
}

func testCustomDirs(t *testing.T, testCases []envTestCase) {
	for _, tc := range testCases {
		t.Setenv(tc.envKey, tc.envVal)

		Reload()

		assert.Equal(t, *tc.got, tc.expected, "result mismatch")
	}
}
