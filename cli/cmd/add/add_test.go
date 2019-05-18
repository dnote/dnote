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
