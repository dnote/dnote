/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package config

import (
	"fmt"
	"io/ioutil"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config holds dnote configuration
type Config struct {
	Editor      string `yaml:"editor"`
	APIEndpoint string `yaml:"apiEndpoint"`
}

// GetPath returns the path to the dnote config file
func GetPath(ctx context.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, consts.ConfigFilename)
}

// Read reads the config file
func Read(ctx context.DnoteCtx) (Config, error) {
	var ret Config

	configPath := GetPath(ctx)
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return ret, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, errors.Wrap(err, "unmarshalling config")
	}

	return ret, nil
}

// Write writes the config to the config file
func Write(ctx context.DnoteCtx, cf Config) error {
	path := GetPath(ctx)

	b, err := yaml.Marshal(cf)
	if err != nil {
		return errors.Wrap(err, "marshalling config into YAML")
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return errors.Wrap(err, "writing the config file")
	}

	return nil
}
