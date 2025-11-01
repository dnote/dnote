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
