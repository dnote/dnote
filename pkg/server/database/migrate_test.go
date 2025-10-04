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
	"io/fs"
	"testing"
	"testing/fstest"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// unsortedFS wraps fstest.MapFS to return entries in reverse order
type unsortedFS struct {
	fstest.MapFS
}

func (u unsortedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := u.MapFS.ReadDir(name)
	if err != nil {
		return nil, err
	}
	// Reverse the entries to ensure they're not in sorted order
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
	return entries, nil
}

// errorFS returns an error on ReadDir
type errorFS struct{}

func (e errorFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (e errorFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return nil, fs.ErrPermission
}

func TestMigrate_createsSchemaTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	migrationsFs := fstest.MapFS{}
	migrate(db, migrationsFs)

	// Verify schema_migrations table exists
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM schema_migrations").Scan(&count).Error; err != nil {
		t.Fatalf("schema_migrations table not found: %v", err)
	}
}

func TestMigrate_idempotency(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Set up table before migration
	if err := db.Exec("CREATE TABLE counter (value INTEGER)").Error; err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Create migration that inserts a row
	migrationsFs := fstest.MapFS{
		"001-insert-data.sql": &fstest.MapFile{
			Data: []byte("INSERT INTO counter (value) VALUES (100);"),
		},
	}

	// Run migration first time
	if err := migrate(db, migrationsFs); err != nil {
		t.Fatalf("first migration failed: %v", err)
	}
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM counter").Scan(&count).Error; err != nil {
		t.Fatalf("failed to count rows: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}

	// Run migration second time - it should not run the SQL again
	if err := migrate(db, migrationsFs); err != nil {
		t.Fatalf("second migration failed: %v", err)
	}
	if err := db.Raw("SELECT COUNT(*) FROM counter").Scan(&count).Error; err != nil {
		t.Fatalf("failed to count rows: %v", err)
	}
	if count != 1 {
		t.Errorf("migration ran twice: expected 1 row, got %d", count)
	}
}

func TestMigrate_ordering(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Create table before migrations
	if err := db.Exec("CREATE TABLE log (value INTEGER)").Error; err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Create migrations with unsorted filesystem
	migrationsFs := unsortedFS{
		MapFS: fstest.MapFS{
			"010-tenth.sql": &fstest.MapFile{
				Data: []byte("INSERT INTO log (value) VALUES (3);"),
			},
			"001-first.sql": &fstest.MapFile{
				Data: []byte("INSERT INTO log (value) VALUES (1);"),
			},
			"002-second.sql": &fstest.MapFile{
				Data: []byte("INSERT INTO log (value) VALUES (2);"),
			},
		},
	}

	// Run migrations
	if err := migrate(db, migrationsFs); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	// Verify migrations ran in correct order (1, 2, 3)
	var values []int
	if err := db.Raw("SELECT value FROM log ORDER BY rowid").Scan(&values).Error; err != nil {
		t.Fatalf("failed to query log: %v", err)
	}

	expected := []int{1, 2, 3}
	if len(values) != len(expected) {
		t.Fatalf("expected %d rows, got %d", len(expected), len(values))
	}

	for i, v := range values {
		if v != expected[i] {
			t.Errorf("row %d: expected value %d, got %d", i, expected[i], v)
		}
	}
}

func TestMigrate_duplicateVersion(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Create migrations with duplicate version numbers
	migrationsFs := fstest.MapFS{
		"001-first.sql": &fstest.MapFile{
			Data: []byte("SELECT 1;"),
		},
		"001-second.sql": &fstest.MapFile{
			Data: []byte("SELECT 2;"),
		},
	}

	// Should return error for duplicate version
	err = migrate(db, migrationsFs)
	if err == nil {
		t.Fatal("expected error for duplicate version numbers, got nil")
	}
}

func TestMigrate_initTableError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Close the database connection to cause exec to fail
	sqlDB, _ := db.DB()
	sqlDB.Close()

	migrationsFs := fstest.MapFS{
		"001-init.sql": &fstest.MapFile{
			Data: []byte("SELECT 1;"),
		},
	}

	// Should return error for table initialization failure
	err = migrate(db, migrationsFs)
	if err == nil {
		t.Fatal("expected error for table initialization failure, got nil")
	}
}

func TestMigrate_readDirError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Use filesystem that fails on ReadDir
	err = migrate(db, errorFS{})
	if err == nil {
		t.Fatal("expected error for ReadDir failure, got nil")
	}
}

func TestMigrate_sqlError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Create migration with invalid SQL
	migrationsFs := fstest.MapFS{
		"001-bad-sql.sql": &fstest.MapFile{
			Data: []byte("INVALID SQL SYNTAX HERE;"),
		},
	}

	// Should return error for SQL execution failure
	err = migrate(db, migrationsFs)
	if err == nil {
		t.Fatal("expected error for invalid SQL, got nil")
	}
}

func TestMigrate_emptyFile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{"completely empty", "", true},
		{"only whitespace", "   \n\t  ", true},
		{"only comments", "-- just a comment", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrationsFs := fstest.MapFS{
				"001-empty.sql": &fstest.MapFile{
					Data: []byte(tt.data),
				},
			}

			err = migrate(db, migrationsFs)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMigrate_invalidFilename(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"valid format", "001-init.sql", false},
		{"no leading zeros", "1-init.sql", true},
		{"two digits", "01-init.sql", true},
		{"no dash", "001init.sql", true},
		{"no description", "001-.sql", true},
		{"no extension", "001-init.", true},
		{"wrong extension", "001-init.txt", true},
		{"non-numeric version number", "0a1-init.sql", true},
		{"underscore separator", "001_init.sql", true},
		{"multiple dashes in description", "001-add-feature-v2.sql", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrationsFs := fstest.MapFS{
				tt.filename: &fstest.MapFile{
					Data: []byte("SELECT 1;"),
				},
			}

			err := migrate(db, migrationsFs)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
