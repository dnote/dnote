package dirs

import (
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

type envTestCase struct {
	envKey   string
	envVal   string
	got      *string
	expected string
}

func testCustomDirs(t *testing.T, testCases []envTestCase) {
	for _, tc := range testCases {
		os.Setenv(tc.envKey, tc.envVal)

		Reload()

		assert.Equal(t, *tc.got, tc.expected, "result mismatch")
	}
}
