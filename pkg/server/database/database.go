/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// Use postgres
	_ "github.com/lib/pq"
)

var (
	// MigrationTableName is the name of the table that keeps track of migrations
	MigrationTableName = "migrations"
)

// Config holds the connection configuration
type Config struct {
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
	if os.Getenv("GO_ENV") == "PRODUCTION" {
		sslmode = "require"
	} else {
		sslmode = "disable"
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

var (
	// DBConn is the connection handle for the database
	DBConn *gorm.DB
)

const (
	// TokenTypeResetPassword is a type of a token for reseting password
	TokenTypeResetPassword = "reset_password"
	// TokenTypeEmailVerification is a type of a token for verifying email
	TokenTypeEmailVerification = "email_verification"
	// TokenTypeEmailPreference is a type of a token for updating email preference
	TokenTypeEmailPreference = "email_preference"
)

const (
	// BookDomainAll incidates that all books are eligible to be the source books
	BookDomainAll = "all"
	// BookDomainIncluding incidates that some specified books are eligible to be the source books
	BookDomainIncluding = "including"
	// BookDomainExluding incidates that all books except for some specified books are eligible to be the source books
	BookDomainExluding = "excluding"
)

// Open opens the connection with the database
func Open(c Config) {
	connStr, err := getPGConnectionString(c)
	if err != nil {
		panic(err)
	}

	DBConn, err = gorm.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

// Close closes database connection
func Close() {
	DBConn.Close()
}

// InitSchema migrates database schema to reflect the latest model definition
func InitSchema() {
	if err := DBConn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		panic(err)
	}

	if err := DBConn.AutoMigrate(
		Note{},
		Book{},
		User{},
		Account{},
		Notification{},
		Token{},
		EmailPreference{},
		Session{},
		Digest{},
		RepetitionRule{},
	).Error; err != nil {
		panic(err)
	}
}
