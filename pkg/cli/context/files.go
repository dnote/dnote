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

package context

import (
	"path/filepath"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
)

// InitDnoteDirs creates the dnote directories if they don't already exist.
func InitDnoteDirs(paths Paths) error {
	if paths.Config != "" {
		configDir := filepath.Join(paths.Config, consts.DnoteDirName)
		if err := utils.EnsureDir(configDir); err != nil {
			return errors.Wrap(err, "initializing config dir")
		}
	}
	if paths.Data != "" {
		dataDir := filepath.Join(paths.Data, consts.DnoteDirName)
		if err := utils.EnsureDir(dataDir); err != nil {
			return errors.Wrap(err, "initializing data dir")
		}
	}
	if paths.Cache != "" {
		cacheDir := filepath.Join(paths.Cache, consts.DnoteDirName)
		if err := utils.EnsureDir(cacheDir); err != nil {
			return errors.Wrap(err, "initializing cache dir")
		}
	}

	return nil
}
