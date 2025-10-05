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
		assert.Equal(t, userCount, int64(1), "book count mismatch")

		var accountCount int64
		var accountRecord database.Account
		testutils.MustExec(t, db.Model(&database.Account{}).Count(&accountCount), "counting account")
		testutils.MustExec(t, db.First(&accountRecord), "finding account")

		assert.Equal(t, accountCount, int64(1), "account count mismatch")
		assert.Equal(t, accountRecord.Email.String, "alice@example.com", "account email mismatch")

		passwordErr := bcrypt.CompareHashAndPassword([]byte(accountRecord.Password.String), []byte("pass1234"))
		assert.Equal(t, passwordErr, nil, "Password mismatch")
	})

	t.Run("duplicate email", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		aliceUser := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, aliceUser, "alice@example.com", "somepassword")

		a := NewTest()
		a.DB = db
		_, err := a.CreateUser("alice@example.com", "newpassword", "newpassword")

		assert.Equal(t, err, ErrDuplicateEmail, "error mismatch")

		var userCount, accountCount int64
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting user")
		testutils.MustExec(t, db.Model(&database.Account{}).Count(&accountCount), "counting account")

		assert.Equal(t, userCount, int64(1), "user count mismatch")
		assert.Equal(t, accountCount, int64(1), "account count mismatch")
	})
}
