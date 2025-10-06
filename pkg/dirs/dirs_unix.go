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

//go:build linux || darwin || freebsd


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
