package login

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/context"
)

func TestGetServerDisplayURL(t *testing.T) {
	testCases := []struct {
		apiEndpoint string
		expected    string
	}{
		{
			apiEndpoint: "https://api.getdnote.com",
			expected:    "https://www.getdnote.com",
		},
		{
			apiEndpoint: "https://dnote.mydomain.com/api",
			expected:    "https://dnote.mydomain.com",
		},
		{
			apiEndpoint: "https://mysubdomain.mydomain.com/dnote/api",
			expected:    "https://mysubdomain.mydomain.com",
		},
		{
			apiEndpoint: "https://dnote.mysubdomain.mydomain.com/api",
			expected:    "https://dnote.mysubdomain.mydomain.com",
		},
		{
			apiEndpoint: "some-string",
			expected:    "",
		},
		{
			apiEndpoint: "",
			expected:    "",
		},
		{
			apiEndpoint: "https://",
			expected:    "",
		},
		{
			apiEndpoint: "https://abc",
			expected:    "https://abc",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("for input %s", tc.apiEndpoint), func(t *testing.T) {
			got := getServerDisplayURL(context.DnoteCtx{APIEndpoint: tc.apiEndpoint})
			assert.Equal(t, got, tc.expected, "result mismatch")
		})
	}
}
