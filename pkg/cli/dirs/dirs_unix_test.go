// +build linux darwin

package dirs

import (
	"path/filepath"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestDirs(t *testing.T) {
	home := Home
	assert.NotEqual(t, home, "", "home is empty")

	configHome := filepath.Join(home, ".config")
	dataHome := filepath.Join(home, ".local", "share")
	cacheHome := filepath.Join(home, ".cache")

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

func TestCustomDirs(t *testing.T) {
	testCases := []envTestCase{
		{
			envKey:   "XDG_CONFIG_HOME",
			envVal:   "~/custom/config",
			got:      &ConfigHome,
			expected: "~/custom/config",
		},
		{
			envKey:   "XDG_DATA_HOME",
			envVal:   "~/custom/data",
			got:      &DataHome,
			expected: "~/custom/data",
		},
		{
			envKey:   "XDG_CACHE_HOME",
			envVal:   "~/custom/cache",
			got:      &CacheHome,
			expected: "~/custom/cache",
		},
	}

	testCustomDirs(t, testCases)
}
