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

package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"golang.org/x/crypto/bcrypt"
)

func TestUserCreateCmd(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"

	// Call the function directly
	userCreateCmd([]string{"--dbPath", tmpDB, "--email", "test@example.com", "--password", "password123"})

	// Verify user was created in database
	db := testutils.InitDB(tmpDB)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	var count int64
	testutils.MustExec(t, db.Model(&database.User{}).Count(&count), "counting users")
	assert.Equal(t, count, int64(1), "should have 1 user")

	var user database.User
	testutils.MustExec(t, db.Where("email = ?", "test@example.com").First(&user), "finding user")
	assert.Equal(t, user.Email.String, "test@example.com", "email mismatch")
}

func TestUserRemoveCmd(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"

	// Create a user first
	db := testutils.InitDB(tmpDB)
	testutils.SetupUserData(db, "test@example.com", "password123")
	sqlDB, _ := db.DB()
	sqlDB.Close()

	// Remove the user with mock stdin that responds "y"
	mockStdin := strings.NewReader("y\n")
	userRemoveCmd([]string{"--dbPath", tmpDB, "--email", "test@example.com"}, mockStdin)

	// Verify user was removed
	db2 := testutils.InitDB(tmpDB)
	defer func() {
		sqlDB2, _ := db2.DB()
		sqlDB2.Close()
	}()

	var count int64
	testutils.MustExec(t, db2.Model(&database.User{}).Count(&count), "counting users")
	assert.Equal(t, count, int64(0), "should have 0 users")
}

func TestUserResetPasswordCmd(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"

	// Create a user first
	db := testutils.InitDB(tmpDB)
	user := testutils.SetupUserData(db, "test@example.com", "oldpassword123")
	oldPasswordHash := user.Password.String
	sqlDB, _ := db.DB()
	sqlDB.Close()

	// Reset password
	userResetPasswordCmd([]string{"--dbPath", tmpDB, "--email", "test@example.com", "--password", "newpassword123"})

	// Verify password was changed
	db2 := testutils.InitDB(tmpDB)
	defer func() {
		sqlDB2, _ := db2.DB()
		sqlDB2.Close()
	}()

	var updatedUser database.User
	testutils.MustExec(t, db2.Where("email = ?", "test@example.com").First(&updatedUser), "finding user")

	// Verify password hash changed
	assert.Equal(t, updatedUser.Password.String != oldPasswordHash, true, "password hash should be different")
	assert.Equal(t, len(updatedUser.Password.String) > 0, true, "password should be set")

	// Verify new password works
	err := bcrypt.CompareHashAndPassword([]byte(updatedUser.Password.String), []byte("newpassword123"))
	assert.Equal(t, err, nil, "new password should match")

	// Verify old password doesn't work
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password.String), []byte("oldpassword123"))
	assert.Equal(t, err != nil, true, "old password should not match")
}

func TestUserListCmd(t *testing.T) {
	t.Run("multiple users", func(t *testing.T) {
		tmpDB := t.TempDir() + "/test.db"

		// Create multiple users
		db := testutils.InitDB(tmpDB)
		user1 := testutils.SetupUserData(db, "alice@example.com", "password123")
		user2 := testutils.SetupUserData(db, "bob@example.com", "password123")
		user3 := testutils.SetupUserData(db, "charlie@example.com", "password123")
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Capture output
		var buf bytes.Buffer
		userListCmd([]string{"--dbPath", tmpDB}, &buf)

		// Verify output matches expected format
		output := strings.TrimSpace(buf.String())
		lines := strings.Split(output, "\n")

		expectedLine1 := fmt.Sprintf("%s,alice@example.com,%s", user1.UUID, user1.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"))
		expectedLine2 := fmt.Sprintf("%s,bob@example.com,%s", user2.UUID, user2.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"))
		expectedLine3 := fmt.Sprintf("%s,charlie@example.com,%s", user3.UUID, user3.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"))

		assert.Equal(t, lines[0], expectedLine1, "line 1 should match")
		assert.Equal(t, lines[1], expectedLine2, "line 2 should match")
		assert.Equal(t, lines[2], expectedLine3, "line 3 should match")
	})

	t.Run("empty database", func(t *testing.T) {
		tmpDB := t.TempDir() + "/test.db"

		// Initialize empty database
		db := testutils.InitDB(tmpDB)
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Capture output
		var buf bytes.Buffer
		userListCmd([]string{"--dbPath", tmpDB}, &buf)

		// Verify no output
		output := buf.String()
		assert.Equal(t, output, "", "should have no output for empty database")
	})
}
