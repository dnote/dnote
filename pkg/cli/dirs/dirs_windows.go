// +build windows

package dirs

import (
	"path/filepath"
)

func initDirs() {
	Home = getHomeDir()
	ConfigHome = filepath.Join(Home, ".dnote")
	DataHome = filepath.Join(Home, ".dnote")
	CacheHome = filepath.Join(Home, ".dnote")
}
