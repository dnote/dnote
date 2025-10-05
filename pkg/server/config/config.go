/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package config

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/dnote/dnote/pkg/dirs"
	"github.com/dnote/dnote/pkg/server/assets"
	"github.com/pkg/errors"
)

const (
	// AppEnvProduction represents an app environment for production.
	AppEnvProduction string = "PRODUCTION"
	// DefaultDBDir is the default directory name for Dnote data
	DefaultDBDir = "dnote"
	// DefaultDBFilename is the default database filename
	DefaultDBFilename = "server.db"
)

var (
	// DefaultDBPath is the default path to the database file
	DefaultDBPath = filepath.Join(dirs.DataHome, DefaultDBDir, DefaultDBFilename)
)

var (
	// ErrDBMissingPath is an error for an incomplete configuration missing the database path
	ErrDBMissingPath = errors.New("DB Path is empty")
	// ErrWebURLInvalid is an error for an incomplete configuration with invalid web url
	ErrWebURLInvalid = errors.New("Invalid WebURL")
	// ErrPortInvalid is an error for an incomplete configuration with invalid port
	ErrPortInvalid = errors.New("Invalid Port")
)

func readBoolEnv(name string) bool {
	return os.Getenv(name) == "true"
}

// getOrEnv returns value if non-empty, otherwise env var, otherwise default
func getOrEnv(value, envKey, defaultVal string) string {
	if value != "" {
		return value
	}
	if env := os.Getenv(envKey); env != "" {
		return env
	}
	return defaultVal
}

// Config is an application configuration
type Config struct {
	AppEnv              string
	WebURL              string
	DisableRegistration bool
	Port                string
	DBPath              string
	AssetBaseURL        string
	HTTP500Page         []byte
	LogLevel            string
}

// Params are the configuration parameters for creating a new Config
type Params struct {
	AppEnv              string
	Port                string
	WebURL              string
	DBPath              string
	DisableRegistration bool
	LogLevel            string
}

// New constructs and returns a new validated config.
// Empty string params will fall back to environment variables and defaults.
func New(p Params) (Config, error) {
	c := Config{
		AppEnv:              getOrEnv(p.AppEnv, "APP_ENV", AppEnvProduction),
		Port:                getOrEnv(p.Port, "PORT", "3000"),
		WebURL:              getOrEnv(p.WebURL, "WebURL", ""),
		DBPath:              getOrEnv(p.DBPath, "DBPath", DefaultDBPath),
		DisableRegistration: p.DisableRegistration || readBoolEnv("DisableRegistration"),
		LogLevel:            getOrEnv(p.LogLevel, "LOG_LEVEL", "info"),
		AssetBaseURL:        "/static",
		HTTP500Page:         assets.MustGetHTTP500ErrorPage(),
	}

	if err := validate(c); err != nil {
		return Config{}, err
	}

	return c, nil
}

// IsProd checks if the app environment is configured to be production.
func (c Config) IsProd() bool {
	return c.AppEnv == AppEnvProduction
}

func validate(c Config) error {
	if _, err := url.ParseRequestURI(c.WebURL); err != nil {
		return errors.Wrapf(ErrWebURLInvalid, "'%s'", c.WebURL)
	}
	if c.Port == "" {
		return ErrPortInvalid
	}

	if c.DBPath == "" {
		return ErrDBMissingPath
	}

	return nil
}
