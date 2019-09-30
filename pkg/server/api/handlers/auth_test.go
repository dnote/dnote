/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"golang.org/x/crypto/bcrypt"
)

func TestGetMe(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	u := testutils.SetupUserData()
	testutils.SetupAccountData(u, "alice@example.com", "somepassword")

	dat := `{"email": "alice@example.com"}`
	req := testutils.MakeReq(server, "POST", "/reset-token", dat)

	// Execute
	res := testutils.HTTPAuthDo(t, req, u)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismtach")

	var user database.User
	testutils.MustExec(t, db.Where("id = ?", u.ID).First(&user), "finding user")
	assert.Equal(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")
}

func TestCreateResetToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "somepassword")

		dat := `{"email": "alice@example.com"}`
		req := testutils.MakeReq(server, "POST", "/reset-token", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismtach")

		var tokenCount int
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting tokens")

		var resetToken database.Token
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", u.ID, database.TokenTypeResetPassword).First(&resetToken), "finding reset token")

		assert.Equal(t, tokenCount, 1, "reset_token count mismatch")
		assert.NotEqual(t, resetToken.Value, nil, "reset_token value mismatch")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "reset_token UsedAt mismatch")
	})

	t.Run("nonexistent email", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "somepassword")

		dat := `{"email": "bob@example.com"}`
		req := testutils.MakeReq(server, "POST", "/reset-token", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismtach")

		var tokenCount int
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting tokens")
		assert.Equal(t, tokenCount, 0, "reset_token count mismatch")
	})
}

func TestResetPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		a := testutils.SetupAccountData(u, "alice@example.com", "oldpassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")
		otherTok := database.Token{
			UserID: u.ID,
			Value:  "somerandomvalue",
			Type:   database.TokenTypeEmailVerification,
		}
		testutils.MustExec(t, db.Save(&otherTok), "preparing another token")

		dat := `{"token": "MivFxYiSMMA4An9dP24DNQ==", "password": "newpassword"}`
		req := testutils.MakeReq(server, "PATCH", "/reset-password", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismatch")

		var resetToken, verificationToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "finding reset token")
		testutils.MustExec(t, db.Where("value = ?", "somerandomvalue").First(&verificationToken), "finding reset token")
		testutils.MustExec(t, db.Where("id = ?", a.ID).First(&account), "finding account")

		assert.NotEqual(t, resetToken.UsedAt, nil, "reset_token UsedAt mismatch")
		passwordErr := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte("newpassword"))
		assert.Equal(t, passwordErr, nil, "Password mismatch")
		assert.Equal(t, verificationToken.UsedAt, (*time.Time)(nil), "verificationToken UsedAt mismatch")
	})

	t.Run("nonexistent token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		a := testutils.SetupAccountData(u, "alice@example.com", "somepassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		dat := `{"token": "-ApMnyvpg59uOU5b-Kf5uQ==", "password": "oldpassword"}`
		req := testutils.MakeReq(server, "PATCH", "/reset-password", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "finding reset token")
		testutils.MustExec(t, db.Where("id = ?", a.ID).First(&account), "finding account")

		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})

	t.Run("expired token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		a := testutils.SetupAccountData(u, "alice@example.com", "somepassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")
		testutils.MustExec(t, db.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := `{"token": "MivFxYiSMMA4An9dP24DNQ==", "password": "oldpassword"}`
		req := testutils.MakeReq(server, "PATCH", "/reset-password", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusGone, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		testutils.MustExec(t, db.Where("id = ?", a.ID).First(&account), "failed to find account")
		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})

	t.Run("used token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		a := testutils.SetupAccountData(u, "alice@example.com", "somepassword")

		usedAt := time.Now().Add(time.Hour * -11).UTC()
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
			UsedAt: &usedAt,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")
		testutils.MustExec(t, db.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := `{"token": "MivFxYiSMMA4An9dP24DNQ==", "password": "oldpassword"}`
		req := testutils.MakeReq(server, "PATCH", "/reset-password", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		testutils.MustExec(t, db.Where("id = ?", a.ID).First(&account), "failed to find account")
		assert.Equal(t, a.Password, account.Password, "password should not have been updated")

		if resetToken.UsedAt.Year() != usedAt.Year() ||
			resetToken.UsedAt.Month() != usedAt.Month() ||
			resetToken.UsedAt.Day() != usedAt.Day() ||
			resetToken.UsedAt.Hour() != usedAt.Hour() ||
			resetToken.UsedAt.Minute() != usedAt.Minute() ||
			resetToken.UsedAt.Second() != usedAt.Second() {
			t.Errorf("used_at should be %+v but got: %+v", usedAt, resetToken.UsedAt)
		}
	})

	t.Run("using wrong type token: email_verification", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		a := testutils.SetupAccountData(u, "alice@example.com", "somepassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeEmailVerification,
		}
		testutils.MustExec(t, db.Save(&tok), "Failed to prepare reset_token")
		testutils.MustExec(t, db.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := `{"token": "MivFxYiSMMA4An9dP24DNQ==", "password": "oldpassword"}`
		req := testutils.MakeReq(server, "PATCH", "/reset-password", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		testutils.MustExec(t, db.Where("id = ?", a.ID).First(&account), "failed to find account")

		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})
}
