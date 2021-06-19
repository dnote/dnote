package views

import (
	"fmt"
	"testing"
	"time"

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

	ctx := viewCtx{}

	for _, tc := range testCases {
		got := ctx.toDateTime(tc.year, tc.month)

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

	ctx := viewCtx{}

	for _, tc := range testCases {
		got := ctx.getFullMonthName(tc.input)

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

	ctx := viewCtx{}

	for _, tc := range testCases {
		got := ctx.excerpt(tc.str, tc.maxLength)
		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}

func TestTimeAgo(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		input    time.Time
		expected string
	}{
		{
			input:    now.Add(-2 * time.Hour),
			expected: "2 hours ago",
		},
		{
			input:    now.Add(-2*time.Hour - 59*time.Minute),
			expected: "2 hours ago",
		},
		{
			input:    now.Add(-23 * time.Hour),
			expected: "23 hours ago",
		},
		{
			input:    now.Add(-23*time.Hour - 59*time.Minute),
			expected: "23 hours ago",
		},
		{
			input:    now.Add(-24 * time.Hour),
			expected: "1 day ago",
		},
		{
			input:    now.Add(-47 * time.Hour),
			expected: "1 day ago",
		},
		{
			input:    now.Add(-48 * time.Hour),
			expected: "2 days ago",
		},

		{
			input:    now.Add(-24 * time.Hour * 7),
			expected: "1 week ago",
		},
		{
			input:    now.Add(-24 * time.Hour * 7 * 2),
			expected: "2 weeks ago",
		},

		{
			input:    now.Add(-24 * time.Hour * 7 * 4),
			expected: "1 month ago",
		},
		{
			input:    now.Add(-24 * time.Hour * 7 * 7),
			expected: "1 month ago",
		},
		{
			input:    now.Add(-24 * time.Hour * 7 * 8),
			expected: "2 months ago",
		},

		{
			input:    now.Add(-24 * time.Hour * 7 * 52),
			expected: "1 year ago",
		},
		{
			input:    now.Add(-24 * time.Hour * 7 * 55),
			expected: "1 year ago",
		},
		{
			input:    now.Add(-24 * time.Hour * 7 * 52 * 2),
			expected: "2 years ago",
		},
	}

	ctx := newViewCtx(Config{})

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input.String()), func(t *testing.T) {
			got := ctx.timeAgo(tc.input)
			assert.Equal(t, got, tc.expected, "result mismatch")
		})
	}
}
