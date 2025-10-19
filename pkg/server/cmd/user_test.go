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

package cmd

import (
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

	var account database.Account
	testutils.MustExec(t, db.Where("email = ?", "test@example.com").First(&account), "finding account")
	assert.Equal(t, account.Email.String, "test@example.com", "email mismatch")
}

func TestUserRemoveCmd(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"

	// Create a user first
	db := testutils.InitDB(tmpDB)
	user := testutils.SetupUserData(db)
	testutils.SetupAccountData(db, user, "test@example.com", "password123")
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
	user := testutils.SetupUserData(db)
	account := testutils.SetupAccountData(db, user, "test@example.com", "oldpassword123")
	oldPasswordHash := account.Password.String
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

	var updatedAccount database.Account
	testutils.MustExec(t, db2.Where("email = ?", "test@example.com").First(&updatedAccount), "finding account")

	// Verify password hash changed
	assert.Equal(t, updatedAccount.Password.String != oldPasswordHash, true, "password hash should be different")
	assert.Equal(t, len(updatedAccount.Password.String) > 0, true, "password should be set")

	// Verify new password works
	err := bcrypt.CompareHashAndPassword([]byte(updatedAccount.Password.String), []byte("newpassword123"))
	assert.Equal(t, err, nil, "new password should match")

	// Verify old password doesn't work
	err = bcrypt.CompareHashAndPassword([]byte(updatedAccount.Password.String), []byte("oldpassword123"))
	assert.Equal(t, err != nil, true, "old password should not match")
}
