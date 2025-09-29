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

package database

import (
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// Use postgres
	_ "github.com/lib/pq"
)

var (
	// MigrationTableName is the name of the table that keeps track of migrations
	MigrationTableName = "migrations"
)

// InitSchema migrates database schema to reflect the latest model definition
func InitSchema(db *gorm.DB) {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(
		Note{},
		Book{},
		User{},
		Account{},
		Notification{},
		Token{},
		EmailPreference{},
		Session{},
	).Error; err != nil {
		panic(err)
	}
}

// Open initializes the database connection
func Open(c config.Config) *gorm.DB {
	db, err := gorm.Open("postgres", c.DB.GetConnectionStr())
	if err != nil {
		panic(errors.Wrap(err, "opening database conection"))
	}

	return db
}
