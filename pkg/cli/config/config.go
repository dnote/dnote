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

package config

import (
	"fmt"
	"os"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config holds dnote configuration
type Config struct {
	Editor             string `yaml:"editor"`
	APIEndpoint        string `yaml:"apiEndpoint"`
	EnableUpgradeCheck bool   `yaml:"enableUpgradeCheck"`
}

func checkLegacyPath(ctx context.DnoteCtx) (string, bool) {
	legacyPath := fmt.Sprintf("%s/%s", ctx.Paths.LegacyDnote, consts.ConfigFilename)

	ok, err := utils.FileExists(legacyPath)
	if err != nil {
		log.Error(errors.Wrapf(err, "checking legacy dnote directory at %s", legacyPath).Error())
	}
	if ok {
		return legacyPath, true
	}

	return "", false
}

// GetPath returns the path to the dnote config file
func GetPath(ctx context.DnoteCtx) string {
	legacyPath, ok := checkLegacyPath(ctx)
	if ok {
		return legacyPath
	}

	return fmt.Sprintf("%s/%s/%s", ctx.Paths.Config, consts.DnoteDirName, consts.ConfigFilename)
}

// Read reads the config file
func Read(ctx context.DnoteCtx) (Config, error) {
	var ret Config

	configPath := GetPath(ctx)
	b, err := os.ReadFile(configPath)
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

	err = os.WriteFile(path, b, 0644)
	if err != nil {
		return errors.Wrap(err, "writing the config file")
	}

	return nil
}
