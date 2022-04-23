/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"os"
)

const (
	// AppEnvProduction represents an app environment for production.
	AppEnvProduction string = "PRODUCTION"
)

var (
	// ErrDBMissingHost is an error for an incomplete configuration missing the host
	ErrDBMissingHost = errors.New("DB Host is empty")
	// ErrDBMissingPort is an error for an incomplete configuration missing the port
	ErrDBMissingPort = errors.New("DB Port is empty")
	// ErrDBMissingName is an error for an incomplete configuration missing the name
	ErrDBMissingName = errors.New("DB Name is empty")
	// ErrDBMissingUser is an error for an incomplete configuration missing the user
	ErrDBMissingUser = errors.New("DB User is empty")
	// ErrWebURLInvalid is an error for an incomplete configuration with invalid web url
	ErrWebURLInvalid = errors.New("Invalid WebURL")
	// ErrPortInvalid is an error for an incomplete configuration with invalid port
	ErrPortInvalid = errors.New("Invalid Port")
)

// PostgresConfig holds the postgres connection configuration.
type PostgresConfig struct {
	SSLMode  string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func readBoolEnv(name string) bool {
	if os.Getenv(name) == "true" {
		return true
	}

	return false
}

// checkSSLMode checks if SSL is required for the database connection
func checkSSLMode() bool {
	// TODO: deprecate DB_NOSSL in favor of DBSkipSSL
	if os.Getenv("DB_NOSSL") != "" {
		return true
	}

	if os.Getenv("DBSkipSSL") == "true" {
		return true
	}

	return os.Getenv("GO_ENV") != "PRODUCTION"
}

func loadDBConfig() PostgresConfig {
	var sslmode string
	if checkSSLMode() {
		sslmode = "disable"
	} else {
		sslmode = "require"
	}

	return PostgresConfig{
		SSLMode:  sslmode,
		Host:     os.Getenv("DBHost"),
		Port:     os.Getenv("DBPort"),
		Name:     os.Getenv("DBName"),
		User:     os.Getenv("DBUser"),
		Password: os.Getenv("DBPassword"),
	}
}

// Config is an application configuration
type Config struct {
	AppEnv              string
	WebURL              string
	OnPremise           bool
	DisableRegistration bool
	Port                string
	DB                  PostgresConfig
	PageTemplateDir     string
	StaticDir           string
	AssetBaseURL        string
}

func getAppEnv() string {
	// DEPRECATED
	goEnv := os.Getenv("GO_ENV")
	if goEnv != "" {
		return goEnv
	}

	return os.Getenv("APP_ENV")
}

// Load constructs and returns a new config based on the environment variables.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	c := Config{
		AppEnv:              getAppEnv(),
		WebURL:              os.Getenv("WebURL"),
		Port:                port,
		OnPremise:           readBoolEnv("OnPremise"),
		DisableRegistration: readBoolEnv("DisableRegistration"),
		DB:                  loadDBConfig(),
		AssetBaseURL:        "",
	}

	if err := validate(c); err != nil {
		panic(err)
	}

	return c
}

// SetOnPremise sets the OnPremise value
func (c *Config) SetOnPremise(val bool) {
	c.OnPremise = val
}

// SetPageTemplateDir sets page template dir for the config
func (c *Config) SetPageTemplateDir(d string) {
	c.PageTemplateDir = d
}

// SetStaticDir sets static dir for the confi
func (c *Config) SetStaticDir(d string) {
	c.StaticDir = d
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

	if c.DB.Host == "" {
		return ErrDBMissingHost
	}
	if c.DB.Port == "" {
		return ErrDBMissingPort
	}
	if c.DB.Name == "" {
		return ErrDBMissingName
	}
	if c.DB.User == "" {
		return ErrDBMissingUser
	}

	return nil
}

// GetConnectionStr returns a postgres connection string.
func (c PostgresConfig) GetConnectionStr() string {
	return fmt.Sprintf(
		"sslmode=%s host=%s port=%s dbname=%s user=%s password=%s",
		c.SSLMode, c.Host, c.Port, c.Name, c.User, c.Password)
}
