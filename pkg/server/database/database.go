/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package database

import (
	"os"
	"path/filepath"
	"time"

	"github.com/dnote/dnote/pkg/server/log"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// MigrationTableName is the name of the table that keeps track of migrations
	MigrationTableName = "migrations"
)

// getDBLogLevel converts application log level to GORM log level
func getDBLogLevel(level string) logger.LogLevel {
	switch level {
	case log.LevelDebug:
		return logger.Info
	case log.LevelInfo:
		return logger.Info
	case log.LevelWarn:
		return logger.Warn
	case log.LevelError:
		return logger.Error
	default:
		return logger.Error
	}
}

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

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(getDBLogLevel(log.GetLevel())),
	})
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
				log.ErrorWrap(err, "WAL checkpoint failed")
			}
		}
	}()
}

// StartPeriodicVacuum runs full VACUUM on a schedule to reclaim space and defragment.
// VACUUM acquires an exclusive lock and blocks all database operations briefly.
func StartPeriodicVacuum(db *gorm.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := db.Exec("VACUUM").Error; err != nil {
				log.ErrorWrap(err, "VACUUM failed")
			}
		}
	}()
}
