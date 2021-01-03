// +build windows

package dirs

import (
	"path/filepath"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestDirs(t *testing.T) {
	home := Home
	assert.NotEqual(t, home, "", "home is empty")

	configHome := filepath.Join(home, ".dnote")
	dataHome := filepath.Join(home, ".dnote")
	cacheHome := filepath.Join(home, ".dnote")

	testCases := []struct {
		got      string
		expected string
	}{
		{
			got:      ConfigHome,
			expected: configHome,
		},
		{
			got:      DataHome,
			expected: dataHome,
		},
		{
			got:      CacheHome,
			expected: cacheHome,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.got, tc.expected, "result mismatch")
	}
}
