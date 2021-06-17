package views

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestToDateTime(t *testing.T) {
	testCases := []struct {
		year     int
		month    int
		expected string
	}{
		{
			year:     2010,
			month:    10,
			expected: "2010-10",
		},
		{
			year:     2010,
			month:    8,
			expected: "2010-08",
		},
	}

	for _, tc := range testCases {
		got := toDateTime(tc.year, tc.month)

		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}

func TestGetFullMonthName(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{
			input:    1,
			expected: "January",
		},
		{
			input:    12,
			expected: "December",
		},
	}

	for _, tc := range testCases {
		got := getFullMonthName(tc.input)

		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}

func TestExcerpt(t *testing.T) {
	testCases := []struct {
		str       string
		maxLength int
		expected  string
	}{
		{
			str:       "hello world",
			maxLength: 5,
			expected:  "hello...",
		},
		{
			str:       "hello world",
			maxLength: 1,
			expected:  "...",
		},
		{
			str:       "hello world",
			maxLength: 7,
			expected:  "hello...",
		},
		{
			str:       "foo bar baz",
			maxLength: 9,
			expected:  "foo bar...",
		},
	}

	for _, tc := range testCases {
		got := exerpt(tc.str, tc.maxLength)
		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}
