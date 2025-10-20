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
	"bytes"
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

var testServerBinary string

func init() {
	// Build server binary in temp directory
	tmpDir := os.TempDir()
	testServerBinary = fmt.Sprintf("%s/dnote-test-server", tmpDir)
	buildCmd := exec.Command("go", "build", "-tags", "fts5", "-o", testServerBinary, "../server")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("failed to build server: %v\n%s", err, out))
	}
}

func TestServerStart(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"
	port := "13456" // Use different port to avoid conflicts with main test server

	// Start server in background
	cmd := exec.Command(testServerBinary, "start", "--port", port)
	cmd.Env = append(os.Environ(),
		"DBPath="+tmpDB,
		"WebURL=http://localhost:"+port,
		"APP_ENV=PRODUCTION",
	)

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	// Ensure cleanup
	cleanup := func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
			cmd.Wait() // Wait for process to fully exit
		}
	}
	defer cleanup()

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
	cleanup()

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

func TestServerRootCommand(t *testing.T) {
	cmd := exec.Command(testServerBinary)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("server command failed: %v", err)
	}

	outputStr := string(output)
	assert.Equal(t, strings.Contains(outputStr, "Dnote server - a simple command line notebook"), true, "output should contain description")
	assert.Equal(t, strings.Contains(outputStr, "start: Start the server"), true, "output should contain start command")
	assert.Equal(t, strings.Contains(outputStr, "version: Print the version"), true, "output should contain version command")
}

func TestServerStartHelp(t *testing.T) {
	cmd := exec.Command(testServerBinary, "start", "--help")
	output, _ := cmd.CombinedOutput()

	outputStr := string(output)
	assert.Equal(t, strings.Contains(outputStr, "dnote-server start [flags]"), true, "output should contain usage")
	assert.Equal(t, strings.Contains(outputStr, "--appEnv"), true, "output should contain appEnv flag")
	assert.Equal(t, strings.Contains(outputStr, "--port"), true, "output should contain port flag")
	assert.Equal(t, strings.Contains(outputStr, "--webUrl"), true, "output should contain webUrl flag")
	assert.Equal(t, strings.Contains(outputStr, "--dbPath"), true, "output should contain dbPath flag")
	assert.Equal(t, strings.Contains(outputStr, "--disableRegistration"), true, "output should contain disableRegistration flag")
}

func TestServerStartInvalidConfig(t *testing.T) {
	cmd := exec.Command(testServerBinary, "start")
	// Set invalid WebURL to trigger validation failure
	cmd.Env = []string{"WebURL=not-a-valid-url"}

	output, err := cmd.CombinedOutput()

	// Should exit with non-zero status
	if err == nil {
		t.Fatal("expected command to fail with invalid config")
	}

	outputStr := string(output)
	assert.Equal(t, strings.Contains(outputStr, "Error:"), true, "output should contain error message")
	assert.Equal(t, strings.Contains(outputStr, "Invalid WebURL"), true, "output should mention invalid WebURL")
	assert.Equal(t, strings.Contains(outputStr, "dnote-server start [flags]"), true, "output should show usage")
	assert.Equal(t, strings.Contains(outputStr, "--webUrl"), true, "output should show flags")
}

func TestServerUnknownCommand(t *testing.T) {
	cmd := exec.Command(testServerBinary, "unknown")
	output, err := cmd.CombinedOutput()

	// Should exit with non-zero status
	if err == nil {
		t.Fatal("expected command to fail with unknown command")
	}

	outputStr := string(output)
	assert.Equal(t, strings.Contains(outputStr, "Unknown command"), true, "output should contain unknown command message")
	assert.Equal(t, strings.Contains(outputStr, "Dnote server - a simple command line notebook"), true, "output should show help")
}

func TestServerUserCreate(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"

	cmd := exec.Command(testServerBinary, "user", "create",
		"--dbPath", tmpDB,
		"--email", "test@example.com",
		"--password", "password123")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("user create failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	assert.Equal(t, strings.Contains(outputStr, "User created successfully"), true, "output should show success message")
	assert.Equal(t, strings.Contains(outputStr, "test@example.com"), true, "output should show email")

	// Verify user exists in database
	db, err := gorm.Open(sqlite.Open(tmpDB), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	var count int64
	db.Table("users").Count(&count)
	assert.Equal(t, count, int64(1), "should have created 1 user")
}

func TestServerUserCreateShortPassword(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"

	cmd := exec.Command(testServerBinary, "user", "create",
		"--dbPath", tmpDB,
		"--email", "test@example.com",
		"--password", "short")
	output, err := cmd.CombinedOutput()

	// Should fail with short password
	if err == nil {
		t.Fatal("expected command to fail with short password")
	}

	outputStr := string(output)
	assert.Equal(t, strings.Contains(outputStr, "password should be longer than 8 characters"), true, "output should show password error")
}

func TestServerUserResetPassword(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"

	// Create user first
	createCmd := exec.Command(testServerBinary, "user", "create",
		"--dbPath", tmpDB,
		"--email", "test@example.com",
		"--password", "oldpassword123")
	if output, err := createCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to create user: %v\nOutput: %s", err, output)
	}

	// Reset password
	resetCmd := exec.Command(testServerBinary, "user", "reset-password",
		"--dbPath", tmpDB,
		"--email", "test@example.com",
		"--password", "newpassword123")
	output, err := resetCmd.CombinedOutput()

	if err != nil {
		t.Fatalf("reset-password failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	assert.Equal(t, strings.Contains(outputStr, "Password reset successfully"), true, "output should show success message")
}

func TestServerUserRemove(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"

	// Create user first
	createCmd := exec.Command(testServerBinary, "user", "create",
		"--dbPath", tmpDB,
		"--email", "test@example.com",
		"--password", "password123")
	if output, err := createCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to create user: %v\nOutput: %s", err, output)
	}

	// Remove user with confirmation
	removeCmd := exec.Command(testServerBinary, "user", "remove",
		"--dbPath", tmpDB,
		"--email", "test@example.com")

	// Pipe "y" to stdin to confirm removal
	stdin, err := removeCmd.StdinPipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}

	// Capture output
	stdout, err := removeCmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	var stderr bytes.Buffer
	removeCmd.Stderr = &stderr

	// Start command
	if err := removeCmd.Start(); err != nil {
		t.Fatalf("failed to start remove command: %v", err)
	}

	// Wait for prompt and send "y" to confirm
	if err := assert.RespondToPrompt(stdout, stdin, "Remove user test@example.com?", "y\n", 10*time.Second); err != nil {
		t.Fatalf("failed to confirm removal: %v", err)
	}

	// Wait for command to finish
	if err := removeCmd.Wait(); err != nil {
		t.Fatalf("user remove failed: %v\nStderr: %s", err, stderr.String())
	}

	// Verify user was removed
	db, err := gorm.Open(sqlite.Open(tmpDB), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	var count int64
	db.Table("users").Count(&count)
	assert.Equal(t, count, int64(0), "should have 0 users after removal")
}

func TestServerUserCreateHelp(t *testing.T) {
	cmd := exec.Command(testServerBinary, "user", "create", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("help command failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)

	// Verify help shows double-dash flags for consistency with CLI
	assert.Equal(t, strings.Contains(outputStr, "--email"), true, "help should show --email (double dash)")
	assert.Equal(t, strings.Contains(outputStr, "--password"), true, "help should show --password (double dash)")
	assert.Equal(t, strings.Contains(outputStr, "--dbPath"), true, "help should show --dbPath (double dash)")
}
