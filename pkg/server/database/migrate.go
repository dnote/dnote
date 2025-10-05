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
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/dnote/dnote/pkg/server/database/migrations"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type migrationFile struct {
	filename string
	version  int
}

// validateMigrationFilename checks if filename follows format: NNN-description.sql
func validateMigrationFilename(name string) error {
	// Check .sql extension
	if !strings.HasSuffix(name, ".sql") {
		return errors.Errorf("invalid migration filename: must end with .sql")
	}

	name = strings.TrimSuffix(name, ".sql")
	parts := strings.SplitN(name, "-", 2)
	if len(parts) != 2 {
		return errors.Errorf("invalid migration filename: must be NNN-description.sql")
	}

	version, description := parts[0], parts[1]

	// Validate version is 3 digits
	if len(version) != 3 {
		return errors.Errorf("invalid migration filename: version must be 3 digits, got %s", version)
	}
	for _, c := range version {
		if c < '0' || c > '9' {
			return errors.Errorf("invalid migration filename: version must be numeric, got %s", version)
		}
	}

	// Validate description is not empty
	if description == "" {
		return errors.Errorf("invalid migration filename: description is required")
	}

	return nil
}

// Migrate runs the migrations using the embedded migration files
func Migrate(db *gorm.DB) error {
	return migrate(db, migrations.Files)
}

// getMigrationFiles reads, validates, and sorts migration files
func getMigrationFiles(fsys fs.FS) ([]migrationFile, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, errors.Wrap(err, "reading migration directory")
	}

	var migrations []migrationFile
	seen := make(map[int]string)
	for _, e := range entries {
		name := e.Name()

		if err := validateMigrationFilename(name); err != nil {
			return nil, err
		}

		// Parse version
		var v int
		fmt.Sscanf(name, "%d", &v)

		// Check for duplicate version numbers
		if existing, found := seen[v]; found {
			return nil, errors.Errorf("duplicate migration version %d: %s and %s", v, existing, name)
		}
		seen[v] = name

		migrations = append(migrations, migrationFile{
			filename: name,
			version:  v,
		})
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	return migrations, nil
}

// migrate runs migrations from the provided filesystem
func migrate(db *gorm.DB, fsys fs.FS) error {
	if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS schema_migrations (
					version INTEGER PRIMARY KEY,
					applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
	`).Error; err != nil {
		return errors.Wrap(err, "initializing migration table")
	}

	// Get current version
	var version int
	if err := db.Raw("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version).Error; err != nil {
		return errors.Wrap(err, "reading current version")
	}

	// Read and validate migration files
	migrations, err := getMigrationFiles(fsys)
	if err != nil {
		return err
	}

	var filenames []string
	for _, m := range migrations {
		filenames = append(filenames, m.filename)
	}

	log.WithFields(log.Fields{
		"version": version,
	}).Info("Database schema version.")

	log.WithFields(log.Fields{
		"files": filenames,
	}).Debug("Database migration files.")

	// Apply pending migrations
	for _, m := range migrations {
		if m.version <= version {
			continue
		}

		log.WithFields(log.Fields{
			"file": m.filename,
		}).Info("Applying migration.")

		sql, err := fs.ReadFile(fsys, m.filename)
		if err != nil {
			return errors.Wrapf(err, "reading migration file %s", m.filename)
		}

		if len(strings.TrimSpace(string(sql))) == 0 {
			return errors.Errorf("migration file %s is empty", m.filename)
		}

		if err := db.Exec(string(sql)).Error; err != nil {
			return fmt.Errorf("migration %s failed: %w", m.filename, err)
		}

		if err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", m.version).Error; err != nil {
			return errors.Wrapf(err, "recording migration %s", m.filename)
		}

		log.WithFields(log.Fields{
			"file": m.filename,
		}).Info("Migrate success.")
	}

	return nil
}
