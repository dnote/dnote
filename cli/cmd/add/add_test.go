/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package add

import (
	"testing"

	"github.com/dnote/dnote/cli/testutils"
)

func TestValidateBookName(t *testing.T) {
	testCases := []struct {
		input    string
		expected error
	}{
		{
			input:    "javascript",
			expected: nil,
		},
		{
			input:    "node.js",
			expected: nil,
		},
		{
			input:    "foo bar",
			expected: nil,
		},
		{
			input:    "123",
			expected: ErrNumericBookName,
		},
		{
			input:    "+123",
			expected: nil,
		},
		{
			input:    "-123",
			expected: nil,
		},
		{
			input:    "+javascript",
			expected: nil,
		},
		{
			input:    "0",
			expected: ErrNumericBookName,
		},
		{
			input:    "0333",
			expected: ErrNumericBookName,
		},

		// reserved book names
		{
			input:    "trash",
			expected: ErrBookNameReserved,
		},
		{
			input:    "conflicts",
			expected: ErrBookNameReserved,
		},
	}

	for _, tc := range testCases {
		actual := validateBookName(tc.input)

		testutils.AssertEqual(t, actual, tc.expected, "result does not match")
	}
}
