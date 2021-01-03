// +build linux darwin

package dirs

import (
	"path/filepath"
)

// The environment variable names for the XDG base directory specification
var (
	envConfigHome = "XDG_CONFIG_HOME"
	envDataHome   = "XDG_DATA_HOME"
	envCacheHome  = "XDG_CACHE_HOME"
)

func initDirs() {
	Home = getHomeDir()
	ConfigHome = readPath(envConfigHome, getConfigHome(Home))
	DataHome = readPath(envDataHome, getDataHome(Home))
	CacheHome = readPath(envCacheHome, getCacheHome(Home))
}

func getConfigHome(homeDir string) string {
	return filepath.Join(homeDir, ".config")
}

func getDataHome(homeDir string) string {
	return filepath.Join(homeDir, ".local/share")
}

func getCacheHome(homeDir string) string {
	return filepath.Join(homeDir, ".cache")
}
