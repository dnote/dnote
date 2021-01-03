// Package dirs provides base directory definitions for the system
package dirs

import (
	"os"
	"os/user"

	"github.com/pkg/errors"
)

var (
	// Home is the home directory of the user
	Home string
	// ConfigHome is the full path to the directory in which user-specific
	// configurations should be written.
	ConfigHome string
	// DataHome is the full path to the directory in which user-specific data
	// files should be written.
	DataHome string
	// CacheHome is the full path to the directory in which user-specific
	// non-essential cached data should be writte
	CacheHome string
)

func init() {
	Reload()
}

// Reload reloads the directory definitions
func Reload() {
	initDirs()
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(errors.Wrap(err, "getting home dir"))
	}

	return usr.HomeDir
}

func readPath(envName, defaultPath string) string {
	if dir := os.Getenv(envName); dir != "" {
		return dir
	}

	return defaultPath
}
