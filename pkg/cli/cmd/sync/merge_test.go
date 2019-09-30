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

package sync

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestReportConflict(t *testing.T) {
	testCases := []struct {
		local    string
		server   string
		expected string
	}{
		{
			local:    "\n",
			server:   "\n",
			expected: "\n",
		},
		{
			local:    "",
			server:   "",
			expected: "",
		},
		{
			local:    "foo",
			server:   "foo",
			expected: "foo",
		},
		{
			local:    "foo\nbar",
			server:   "foo\nbar",
			expected: "foo\nbar",
		},
		{
			local:  "foo-local",
			server: "foo-server",
			expected: `<<<<<<< Local
foo-local
=======
foo-server
>>>>>>> Server
`,
		},
		{
			local:  "foo\n",
			server: "bar\n",
			expected: `<<<<<<< Local
foo
=======
bar
>>>>>>> Server
`,
		},
		{
			local:  "foo\n",
			server: "\n",
			expected: `<<<<<<< Local
foo
=======

>>>>>>> Server
`,
		},

		{
			local:  "\n",
			server: "foo\n",
			expected: `<<<<<<< Local

=======
foo
>>>>>>> Server
`,
		},
		{
			local:  "foo\n\nquz\nbaz\n",
			server: "foo\n\nbar\nbaz\n",
			expected: `foo

<<<<<<< Local
quz
=======
bar
>>>>>>> Server
baz
`,
		},
		{
			local:  "foo\n\nquz\nbaz\n\nqux quz\nfuz\n",
			server: "foo\n\nbar\nbaz\n\nqux quz\nfuuz\n",
			expected: `foo

<<<<<<< Local
quz
=======
bar
>>>>>>> Server
baz

qux quz
<<<<<<< Local
fuz
=======
fuuz
>>>>>>> Server
`,
		},
		{
			local:  "foo\nquz\nbaz\nbar\n",
			server: "foo\nquzz\nbazz\nbar\n",
			expected: `foo
<<<<<<< Local
quz
baz
=======
quzz
bazz
>>>>>>> Server
bar
`,
		},
	}

	for idx, tc := range testCases {
		result := reportBodyConflict(tc.local, tc.server)

		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			assert.DeepEqual(t, result, tc.expected, "result mismatch")
		})
	}
}
