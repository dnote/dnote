/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
