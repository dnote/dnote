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

package diff

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestDo(t *testing.T) {
	testCases := []struct {
		s1       string
		s2       string
		expected []diffmatchpatch.Diff
	}{
		{
			s1:       "",
			s2:       "",
			expected: []diffmatchpatch.Diff{},
		},
		{
			s1: "",
			s2: "foo",
			expected: []diffmatchpatch.Diff{
				{
					Type: diffmatchpatch.DiffInsert,
					Text: "foo",
				},
			},
		},
		{
			s1: "foo",
			s2: "",
			expected: []diffmatchpatch.Diff{
				{
					Type: DiffDelete,
					Text: "foo",
				},
			},
		},
		{
			s1: "foo",
			s2: "bar",
			expected: []diffmatchpatch.Diff{
				{
					Type: DiffDelete,
					Text: "foo",
				},
				{
					Type: DiffInsert,
					Text: "bar",
				},
			},
		},
		{
			s1: "foo\nbar\nbaz",
			s2: "foo\nbar\nquz",
			expected: []diffmatchpatch.Diff{
				{
					Type: DiffEqual,
					Text: "foo\nbar\n",
				},
				{
					Type: DiffDelete,
					Text: "baz",
				},
				{
					Type: DiffInsert,
					Text: "quz",
				},
			},
		},
		{
			s1: "fuz\nbar\nbaz\nquz",
			s2: "foo\nbar\nbaz\nqux",
			expected: []diffmatchpatch.Diff{
				{
					Type: DiffDelete,
					Text: "fuz\n",
				},
				{
					Type: DiffInsert,
					Text: "foo\n",
				},
				{
					Type: DiffEqual,
					Text: "bar\nbaz\n",
				},
				{
					Type: DiffDelete,
					Text: "quz",
				},
				{
					Type: DiffInsert,
					Text: "qux",
				},
			},
		},
		{
			s1: "foo bar\nhello dnote\nbaz quz",
			s2: "foo bar\nhello foo\nbaz quz",
			expected: []diffmatchpatch.Diff{
				{
					Type: DiffEqual,
					Text: "foo bar\n",
				},
				{
					Type: DiffDelete,
					Text: "hello dnote\n",
				},
				{
					Type: DiffInsert,
					Text: "hello foo\n",
				},
				{
					Type: DiffEqual,
					Text: "baz quz",
				},
			},
		},
	}

	for idx, tc := range testCases {
		result := Do(tc.s1, tc.s2)

		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			assert.DeepEqual(t, result, tc.expected, "result mismatch")
		})
	}
}
