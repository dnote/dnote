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
=======
quzz
>>>>>>> Server
<<<<<<< Local
baz
=======
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
