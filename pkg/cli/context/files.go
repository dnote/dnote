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
