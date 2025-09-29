/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
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
