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
	"log"

	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

// Migrate runs the migrations
func Migrate() error {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "../database/migrations/"),
	}

	migrate.SetTable(MigrationTableName)

	db := DBConn.DB()
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "running migrations")
	}

	log.Printf("Performed %d migrations", n)

	return nil
}
