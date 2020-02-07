package config

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"os"
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
	WebURL              string
	OnPremise           bool
	DisableRegistration bool
	Port                string
	DB                  PostgresConfig
}

// Load constructs and returns a new config based on the environment variables.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	c := Config{
		WebURL:              os.Getenv("WebURL"),
		Port:                port,
		OnPremise:           readBoolEnv("OnPremise"),
		DisableRegistration: readBoolEnv("DisableRegistration"),
		DB:                  loadDBConfig(),
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
