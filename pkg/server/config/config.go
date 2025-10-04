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
)

var (
	// ErrDBMissingPath is an error for an incomplete configuration missing the database path
	ErrDBMissingPath = errors.New("DB Path is empty")
	// ErrWebURLInvalid is an error for an incomplete configuration with invalid web url
	ErrWebURLInvalid = errors.New("Invalid WebURL")
	// ErrPortInvalid is an error for an incomplete configuration with invalid port
	ErrPortInvalid = errors.New("Invalid Port")
)

// DBConfig holds the database connection configuration.
type DBConfig struct {
	Path string
}

func readBoolEnv(name string) bool {
	if os.Getenv(name) == "true" {
		return true
	}

	return false
}

func LoadDBConfig() DBConfig {
	path := os.Getenv("DBPath")
	if path == "" {
		path = filepath.Join(dirs.DataHome, "dnote", "server.db")
	}

	return DBConfig{
		Path: path,
	}
}

// Config is an application configuration
type Config struct {
	AppEnv              string
	WebURL              string
	DisableRegistration bool
	Port                string
	DB                  DBConfig
	AssetBaseURL        string
	HTTP500Page         []byte
}

func getAppEnv() string {
	// DEPRECATED
	goEnv := os.Getenv("GO_ENV")
	if goEnv != "" {
		return goEnv
	}

	return os.Getenv("APP_ENV")
}

func checkDeprecatedEnvVars() {
}

// Load constructs and returns a new config based on the environment variables.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	checkDeprecatedEnvVars()

	c := Config{
		AppEnv:              getAppEnv(),
		WebURL:              os.Getenv("WebURL"),
		Port:                port,
		DisableRegistration: readBoolEnv("DisableRegistration"),
		DB:                  LoadDBConfig(),
		AssetBaseURL:        "",
		HTTP500Page:         assets.MustGetHTTP500ErrorPage(),
	}

	if err := validate(c); err != nil {
		panic(err)
	}

	return c
}

// SetAssetBaseURL sets static dir for the confi
func (c *Config) SetAssetBaseURL(d string) {
	c.AssetBaseURL = d
}

// IsProd checks if the app environment is configured to be production.
func (c Config) IsProd() bool {
	return c.AppEnv == AppEnvProduction
}

func validate(c Config) error {
	if _, err := url.ParseRequestURI(c.WebURL); err != nil {
		return errors.Wrapf(ErrWebURLInvalid, "provided: '%s'", c.WebURL)
	}
	if c.Port == "" {
		return ErrPortInvalid
	}

	if c.DB.Path == "" {
		return ErrDBMissingPath
	}

	return nil
}
