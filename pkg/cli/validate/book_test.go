/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package validate

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
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
			input:    "",
			expected: ErrBookNameEmpty,
		},
		{
			input:    "foo bar",
			expected: ErrBookNameHasSpace,
		},
		{
			input:    "123",
			expected: ErrBookNameNumeric,
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
			expected: ErrBookNameNumeric,
		},
		{
			input:    "0333",
			expected: ErrBookNameNumeric,
		},
		{
			input:    " javascript",
			expected: ErrBookNameHasSpace,
		},
		{
			input:    "java script",
			expected: ErrBookNameHasSpace,
		},
		{
			input:    "javascript (1)",
			expected: ErrBookNameHasSpace,
		},
		{
			input:    "javascript ",
			expected: ErrBookNameHasSpace,
		},
		{
			input:    "javascript (1) (2) (3)",
			expected: ErrBookNameHasSpace,
		},

		// multiline
		{
			input:    "\n",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "\n\n",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "foo\n",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "foo\nbar\n",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "foo\nbar\nbaz",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "\r\n",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "\r\n\r\n",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "foo\r\n",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "foo\r\nbar\r\n",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "foo\r\nbar\r\nbaz",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "\n\r\n",
			expected: ErrBookNameMultiline,
		},
		{
			input:    "foo\nbar\r\n",
			expected: ErrBookNameMultiline,
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
		actual := BookName(tc.input)

		assert.Equal(t, actual, tc.expected, fmt.Sprintf("result does not match for the input '%s'", tc.input))
	}
}
