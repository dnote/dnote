/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestServerStart(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"
	port := "3456" // Use non-standard port to avoid conflicts

	// Start server in background
	cmd := exec.Command("go", "run", "-tags", "fts5", "../server", "-port", port, "start")
	cmd.Env = append(os.Environ(),
		"DBPath="+tmpDB,
		"WebURL=http://localhost:"+port,
		"APP_ENV=PRODUCTION",
	)
	// Capture output for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	// Wait for server to start and migrations to run
	time.Sleep(3 * time.Second)

	// Verify server responds to health check
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", port))
	if err != nil {
		t.Fatalf("failed to reach server health endpoint: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, 200, "health endpoint should return 200")

	// Kill server before checking database to avoid locks
	if cmd.Process != nil {
		cmd.Process.Kill()
		cmd.Wait() // Clean up zombie process
	}

	// Verify database file was created
	if _, err := os.Stat(tmpDB); os.IsNotExist(err) {
		t.Fatalf("database file was not created at %s", tmpDB)
	}

	// Verify migrations ran by checking database
	db, err := gorm.Open(sqlite.Open(tmpDB), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Verify migrations ran
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM schema_migrations").Scan(&count).Error; err != nil {
		t.Fatalf("schema_migrations table not found: %v", err)
	}
	if count == 0 {
		t.Fatal("no migrations were run")
	}

	// Verify FTS table exists and is functional
	if err := db.Exec("SELECT * FROM notes_fts LIMIT 1").Error; err != nil {
		t.Fatalf("notes_fts table not found or not functional: %v", err)
	}
}

func TestServerVersion(t *testing.T) {
	cmd := exec.Command("go", "run", "-tags", "fts5", "../server", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "dnote-server-") {
		t.Errorf("expected version output to contain 'dnote-server-', got: %s", outputStr)
	}
}
