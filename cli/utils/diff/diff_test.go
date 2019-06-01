package diff

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/cli/testutils"
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
			testutils.AssertDeepEqual(t, result, tc.expected, "result mismatch")
		})
	}
}
