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

	// Use postgres
	_ "github.com/lib/pq"
)

var (
	// MigrationTableName is the name of the table that keeps track of migrations
	MigrationTableName = "migrations"
)

func getPGConnectionString() string {
	ret := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s",
		os.Getenv("DBHost"),
		os.Getenv("DBPort"),
		os.Getenv("DBName"),
		os.Getenv("DBUser"),
		os.Getenv("DBPassword"),
	)

	if os.Getenv("GO_ENV") != "PRODUCTION" {
		ret = fmt.Sprintf("%s sslmode=disable", ret)
	}

	return ret
}

var (
	// DBConn is the connection handle for the database
	DBConn *gorm.DB
)

const (
	// TokenTypeEmailVerification is a type of a token for verifying email
	TokenTypeEmailVerification = "email_verification"
	// TokenTypeEmailPreference is a type of a token for updating email preference
	TokenTypeEmailPreference = "email_preference"
)

// InitDB opens the connection with the database
func InitDB() {
	var err error

	connStr := getPGConnectionString()

	DBConn, err = gorm.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

// CloseDB closes database connection
func CloseDB() {
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
	).Error; err != nil {
		panic(err)
	}
}
