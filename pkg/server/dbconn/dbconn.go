package dbconn

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Config holds the connection configuration
type Config struct {
	SkipSSL  bool
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// ErrConfigMissingHost is an error for an incomplete configuration missing the host
var ErrConfigMissingHost = errors.New("Host is empty")

// ErrConfigMissingPort is an error for an incomplete configuration missing the port
var ErrConfigMissingPort = errors.New("Port is empty")

// ErrConfigMissingName is an error for an incomplete configuration missing the name
var ErrConfigMissingName = errors.New("Name is empty")

// ErrConfigMissingUser is an error for an incomplete configuration missing the user
var ErrConfigMissingUser = errors.New("User is empty")

func validateConfig(c Config) error {
	if c.Host == "" {
		return ErrConfigMissingHost
	}
	if c.Port == "" {
		return ErrConfigMissingPort
	}
	if c.Name == "" {
		return ErrConfigMissingName
	}
	if c.User == "" {
		return ErrConfigMissingUser
	}

	return nil
}

func getPGConnectionString(c Config) (string, error) {
	if err := validateConfig(c); err != nil {
		return "", errors.Wrap(err, "invalid database config")
	}

	var sslmode string
	if c.SkipSSL {
		sslmode = "disable"
	} else {
		sslmode = "require"
	}

	return fmt.Sprintf(
		"sslmode=%s host=%s port=%s dbname=%s user=%s password=%s",
		sslmode,
		c.Host,
		c.Port,
		c.Name,
		c.User,
		c.Password,
	), nil
}

// Open opens the connection with the database
func Open(c Config) *gorm.DB {
	connStr, err := getPGConnectionString(c)
	if err != nil {
		panic(errors.Wrap(err, "getting connection string"))
	}

	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		panic(errors.Wrap(err, "opening database connection"))
	}

	return conn
}
