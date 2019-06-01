package sync

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/cli/testutils"
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
			local:  "foo\nbar",
			server: "foo\nbar",
			expected: `foo
bar
`,
		},
		{
			local:  "foo",
			server: "bar",
			expected: `<<<<<<< Local
foo
=======
bar
>>>>>>> Server
`,
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
			local:  "foo\n\nquz\nbaz",
			server: "foo\n\nbar\nbaz",
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
	}

	for idx, tc := range testCases {
		result := reportBodyConflict(tc.local, tc.server)

		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			testutils.AssertDeepEqual(t, result, tc.expected, "result mismatch")
		})
	}
}
