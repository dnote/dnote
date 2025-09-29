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
	"log"
	"net/http"

	"github.com/dnote/dnote/pkg/server/database/migrations"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

// Migrate runs the migrations
func Migrate(db *gorm.DB) error {
	migrations := &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FileSystem(http.FS(migrations.Files)),
	}

	migrate.SetTable(MigrationTableName)

	n, err := migrate.Exec(db.DB(), "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "running migrations")
	}

	log.Printf("Performed %d migrations", n)

	return nil
}
