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
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// MigrationTableName is the name of the table that keeps track of migrations
	MigrationTableName = "migrations"
)

// InitSchema migrates database schema to reflect the latest model definition
func InitSchema(db *gorm.DB) {
	if err := db.AutoMigrate(
		&User{},
		&Book{},
		&Note{},
		&Token{},
		&Session{},
	); err != nil {
		panic(err)
	}
}

// Open initializes the database connection
func Open(dbPath string) *gorm.DB {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(errors.Wrapf(err, "creating database directory at %s", dir))
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(errors.Wrap(err, "opening database conection"))
	}

	// Get underlying *sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		panic(errors.Wrap(err, "getting underlying database connection"))
	}

	// Configure connection pool for SQLite with WAL mode
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(0) // Doesn't expire.

	// Apply performance PRAGMAs
	pragmas := []string{
		"PRAGMA journal_mode=WAL",   // Enable WAL mode for better concurrency
		"PRAGMA synchronous=NORMAL", // Balance between safety and speed
		"PRAGMA cache_size=-64000",  // 64MB cache (negative = KB)
		"PRAGMA busy_timeout=5000",  // Wait up to 5s for locks
		"PRAGMA foreign_keys=ON",    // Enforce foreign key constraints
		"PRAGMA temp_store=MEMORY",  // Store temp tables in memory
	}

	for _, pragma := range pragmas {
		if err := db.Exec(pragma).Error; err != nil {
			panic(errors.Wrapf(err, "executing pragma: %s", pragma))
		}
	}

	return db
}

// StartWALCheckpointing starts a background goroutine that periodically
// checkpoints the WAL file to prevent it from growing unbounded
func StartWALCheckpointing(db *gorm.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			// TRUNCATE mode removes the WAL file after checkpointing
			if err := db.Exec("PRAGMA wal_checkpoint(TRUNCATE)").Error; err != nil {
				// Log error but don't panic - this is a background maintenance task
				// TODO: Use proper logging once available
				_ = err
			}
		}
	}()
}

// StartPeriodicVacuum runs full VACUUM on a schedule to reclaim space and defragment.
// WARNING: VACUUM acquires an exclusive lock and blocks all database operations briefly.
func StartPeriodicVacuum(db *gorm.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := db.Exec("VACUUM").Error; err != nil {
				// Log error but don't panic - this is a background maintenance task
				// TODO: Use proper logging once available
				_ = err
			}
		}
	}()
}
