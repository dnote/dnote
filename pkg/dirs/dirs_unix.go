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
