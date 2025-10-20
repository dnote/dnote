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

package app

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func TestValidatePassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  nil,
		},
		{
			name:     "valid password exactly 8 chars",
			password: "12345678",
			wantErr:  nil,
		},
		{
			name:     "password too short",
			password: "1234567",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  ErrPasswordTooShort,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePassword(tc.password)
			assert.Equal(t, err, tc.wantErr, "error mismatch")
		})
	}
}

func TestCreateUser_ProValue(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	a := NewTest()
	a.DB = db
	if _, err := a.CreateUser("alice@example.com", "pass1234", "pass1234"); err != nil {
		t.Fatal(errors.Wrap(err, "executing"))
	}

	var userCount int64
	var userRecord database.User
	testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting user")
	testutils.MustExec(t, db.First(&userRecord), "finding user")

	assert.Equal(t, userCount, int64(1), "book count mismatch")

}

func TestGetUserByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "password123")

		a := NewTest()
		a.DB = db

		foundUser, err := a.GetUserByEmail("alice@example.com")

		assert.Equal(t, err, nil, "should not error")
		assert.Equal(t, foundUser.Email.String, "alice@example.com", "email mismatch")
		assert.Equal(t, foundUser.ID, user.ID, "user ID mismatch")
	})

	t.Run("not found", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		a := NewTest()
		a.DB = db

		user, err := a.GetUserByEmail("nonexistent@example.com")

		assert.Equal(t, err, ErrNotFound, "should return ErrNotFound")
		assert.Equal(t, user, (*database.User)(nil), "user should be nil")
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		a := NewTest()
		a.DB = db
		if _, err := a.CreateUser("alice@example.com", "pass1234", "pass1234"); err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		var userCount int64
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting user")
		assert.Equal(t, userCount, int64(1), "user count mismatch")

		var userRecord database.User
		testutils.MustExec(t, db.First(&userRecord), "finding user")

		assert.Equal(t, userRecord.Email.String, "alice@example.com", "user email mismatch")

		passwordErr := bcrypt.CompareHashAndPassword([]byte(userRecord.Password.String), []byte("pass1234"))
		assert.Equal(t, passwordErr, nil, "Password mismatch")
	})

	t.Run("duplicate email", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		testutils.SetupUserData(db, "alice@example.com", "somepassword")

		a := NewTest()
		a.DB = db
		_, err := a.CreateUser("alice@example.com", "newpassword", "newpassword")

		assert.Equal(t, err, ErrDuplicateEmail, "error mismatch")

		var userCount int64
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting user")

		assert.Equal(t, userCount, int64(1), "user count mismatch")
	})
}

func TestUpdateUserPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "oldpassword123")

		err := UpdateUserPassword(db, &user, "newpassword123")

		assert.Equal(t, err, nil, "should not error")

		// Verify password was updated in database
		var updatedUser database.User
		testutils.MustExec(t, db.Where("id = ?", user.ID).First(&updatedUser), "finding updated user")

		// Verify new password works
		passwordErr := bcrypt.CompareHashAndPassword([]byte(updatedUser.Password.String), []byte("newpassword123"))
		assert.Equal(t, passwordErr, nil, "New password should match")

		// Verify old password no longer works
		oldPasswordErr := bcrypt.CompareHashAndPassword([]byte(updatedUser.Password.String), []byte("oldpassword123"))
		assert.NotEqual(t, oldPasswordErr, nil, "Old password should not match")
	})

	t.Run("password too short", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "oldpassword123")

		err := UpdateUserPassword(db, &user, "short")

		assert.Equal(t, err, ErrPasswordTooShort, "should return ErrPasswordTooShort")

		// Verify password was NOT updated in database
		var unchangedUser database.User
		testutils.MustExec(t, db.Where("id = ?", user.ID).First(&unchangedUser), "finding unchanged user")

		// Verify old password still works
		passwordErr := bcrypt.CompareHashAndPassword([]byte(unchangedUser.Password.String), []byte("oldpassword123"))
		assert.Equal(t, passwordErr, nil, "Old password should still match")
	})

	t.Run("empty password", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "oldpassword123")

		err := UpdateUserPassword(db, &user, "")

		assert.Equal(t, err, ErrPasswordTooShort, "should return ErrPasswordTooShort")

		// Verify password was NOT updated in database
		var unchangedUser database.User
		testutils.MustExec(t, db.Where("id = ?", user.ID).First(&unchangedUser), "finding unchanged user")

		// Verify old password still works
		passwordErr := bcrypt.CompareHashAndPassword([]byte(unchangedUser.Password.String), []byte("oldpassword123"))
		assert.Equal(t, passwordErr, nil, "Old password should still match")
	})

	t.Run("transaction rollback", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "oldpassword123")

		// Start a transaction and rollback to verify UpdateUserPassword respects transactions
		tx := db.Begin()
		err := UpdateUserPassword(tx, &user, "newpassword123")
		assert.Equal(t, err, nil, "should not error")
		tx.Rollback()

		// Verify password was NOT updated after rollback
		var unchangedUser database.User
		testutils.MustExec(t, db.Where("id = ?", user.ID).First(&unchangedUser), "finding unchanged user")

		// Verify old password still works
		passwordErr := bcrypt.CompareHashAndPassword([]byte(unchangedUser.Password.String), []byte("oldpassword123"))
		assert.Equal(t, passwordErr, nil, "Old password should still match after rollback")
	})

	t.Run("transaction commit", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "oldpassword123")

		// Start a transaction and commit to verify UpdateUserPassword respects transactions
		tx := db.Begin()
		err := UpdateUserPassword(tx, &user, "newpassword123")
		assert.Equal(t, err, nil, "should not error")
		tx.Commit()

		// Verify password was updated after commit
		var updatedUser database.User
		testutils.MustExec(t, db.Where("id = ?", user.ID).First(&updatedUser), "finding updated user")

		// Verify new password works
		passwordErr := bcrypt.CompareHashAndPassword([]byte(updatedUser.Password.String), []byte("newpassword123"))
		assert.Equal(t, passwordErr, nil, "New password should match after commit")
	})
}

func TestRemoveUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		testutils.SetupUserData(db, "alice@example.com", "password123")

		a := NewTest()
		a.DB = db

		err := a.RemoveUser("alice@example.com")

		assert.Equal(t, err, nil, "should not error")

		// Verify user was deleted
		var userCount int64
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting users")
		assert.Equal(t, userCount, int64(0), "user should be deleted")
	})

	t.Run("user not found", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		a := NewTest()
		a.DB = db

		err := a.RemoveUser("nonexistent@example.com")

		assert.Equal(t, err, ErrNotFound, "should return ErrNotFound")
	})

	t.Run("user has notes", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "password123")

		book := database.Book{UserID: user.ID, Label: "testbook", Deleted: false}
		testutils.MustExec(t, db.Save(&book), "creating book")

		note := database.Note{UserID: user.ID, BookUUID: book.UUID, Body: "test note", Deleted: false}
		testutils.MustExec(t, db.Save(&note), "creating note")

		a := NewTest()
		a.DB = db

		err := a.RemoveUser("alice@example.com")

		assert.Equal(t, err, ErrUserHasExistingResources, "should return ErrUserHasExistingResources")

		// Verify user was NOT deleted
		var userCount int64
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting users")
		assert.Equal(t, userCount, int64(1), "user should not be deleted")

	})

	t.Run("user has books", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "password123")

		book := database.Book{UserID: user.ID, Label: "testbook", Deleted: false}
		testutils.MustExec(t, db.Save(&book), "creating book")

		a := NewTest()
		a.DB = db

		err := a.RemoveUser("alice@example.com")

		assert.Equal(t, err, ErrUserHasExistingResources, "should return ErrUserHasExistingResources")

		// Verify user was NOT deleted
		var userCount int64
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting users")
		assert.Equal(t, userCount, int64(1), "user should not be deleted")

	})

	t.Run("user has deleted notes and books", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "password123")

		book := database.Book{UserID: user.ID, Label: "testbook", Deleted: false}
		testutils.MustExec(t, db.Save(&book), "creating book")

		note := database.Note{UserID: user.ID, BookUUID: book.UUID, Body: "test note", Deleted: false}
		testutils.MustExec(t, db.Save(&note), "creating note")

		// Soft delete the note and book
		testutils.MustExec(t, db.Model(&note).Update("deleted", true), "soft deleting note")
		testutils.MustExec(t, db.Model(&book).Update("deleted", true), "soft deleting book")

		a := NewTest()
		a.DB = db

		err := a.RemoveUser("alice@example.com")

		assert.Equal(t, err, nil, "should not error when user only has deleted notes and books")

		// Verify user was deleted
		var userCount int64
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting users")
		assert.Equal(t, userCount, int64(0), "user should be deleted")

	})
}
