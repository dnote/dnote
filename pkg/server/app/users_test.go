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

func TestGetAllUsers(t *testing.T) {
	t.Run("success with multiple users", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user1 := testutils.SetupUserData(db, "alice@example.com", "password123")
		user2 := testutils.SetupUserData(db, "bob@example.com", "password123")
		user3 := testutils.SetupUserData(db, "charlie@example.com", "password123")

		a := NewTest()
		a.DB = db

		users, err := a.GetAllUsers()

		assert.Equal(t, err, nil, "should not error")
		assert.Equal(t, len(users), 3, "should return 3 users")

		// Verify all users are returned
		emails := make(map[string]bool)
		for _, user := range users {
			emails[user.Email.String] = true
		}
		assert.Equal(t, emails["alice@example.com"], true, "alice should be in results")
		assert.Equal(t, emails["bob@example.com"], true, "bob should be in results")
		assert.Equal(t, emails["charlie@example.com"], true, "charlie should be in results")

		// Verify user details match
		for _, user := range users {
			if user.Email.String == "alice@example.com" {
				assert.Equal(t, user.ID, user1.ID, "alice ID mismatch")
			} else if user.Email.String == "bob@example.com" {
				assert.Equal(t, user.ID, user2.ID, "bob ID mismatch")
			} else if user.Email.String == "charlie@example.com" {
				assert.Equal(t, user.ID, user3.ID, "charlie ID mismatch")
			}
		}
	})

	t.Run("empty database", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		a := NewTest()
		a.DB = db

		users, err := a.GetAllUsers()

		assert.Equal(t, err, nil, "should not error")
		assert.Equal(t, len(users), 0, "should return 0 users")
	})

	t.Run("single user", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		user := testutils.SetupUserData(db, "alice@example.com", "password123")

		a := NewTest()
		a.DB = db

		users, err := a.GetAllUsers()

		assert.Equal(t, err, nil, "should not error")
		assert.Equal(t, len(users), 1, "should return 1 user")
		assert.Equal(t, users[0].Email.String, "alice@example.com", "email mismatch")
		assert.Equal(t, users[0].ID, user.ID, "user ID mismatch")
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
