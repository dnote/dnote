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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/api/presenters"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	testutils.InitTestDB()

}

func TestUpdatePassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		user := testutils.SetupUserData()
		testutils.SetupAccountData(user, "alice@example.com", "oldpassword")

		// Execute
		dat := `{"old_password": "oldpassword", "new_password": "newpassword"}`
		req := testutils.MakeReq(server, "PATCH", "/account/password", dat)
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismsatch")

		var account database.Account
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")

		passwordErr := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte("newpassword"))
		assert.Equal(t, passwordErr, nil, "Password mismatch")
	})

	t.Run("old password mismatch", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		a := testutils.SetupAccountData(u, "alice@example.com", "oldpassword")

		// Execute
		dat := `{"old_password": "randompassword", "new_password": "newpassword"}`
		req := testutils.MakeReq(server, "PATCH", "/account/password", dat)
		res := testutils.HTTPAuthDo(t, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "Status code mismsatch")

		var account database.Account
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&account), "finding account")
		assert.Equal(t, a.Password.String, account.Password.String, "password should not have been updated")
	})

	t.Run("password too short", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		a := testutils.SetupAccountData(u, "alice@example.com", "oldpassword")

		// Execute
		dat := `{"old_password": "oldpassword", "new_password": "a"}`
		req := testutils.MakeReq(server, "PATCH", "/account/password", dat)
		res := testutils.HTTPAuthDo(t, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismsatch")

		var account database.Account
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&account), "finding account")
		assert.Equal(t, a.Password.String, account.Password.String, "password should not have been updated")
	})
}

func TestCreateVerificationToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup

		// TODO: send emails in the background using job queue to avoid coupling the
		// handler itself to the mailer
		templatePath := fmt.Sprintf("%s/mailer/templates/src", testutils.ServerPath)
		mailer.InitTemplates(&templatePath)

		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		user := testutils.SetupUserData()
		testutils.SetupAccountData(user, "alice@example.com", "pass1234")

		// Execute
		req := testutils.MakeReq(server, "POST", "/verification-token", "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusCreated, "status code mismatch")

		var account database.Account
		var token database.Token
		var tokenCount int
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, false, "email_verified should not have been updated")
		assert.NotEqual(t, token.Value, "", "token Value mismatch")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.Equal(t, token.UsedAt, (*time.Time)(nil), "token UsedAt mismatch")
	})

	t.Run("already verified", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		user := testutils.SetupUserData()
		a := testutils.SetupAccountData(user, "alice@example.com", "pass1234")
		a.EmailVerified = true
		testutils.MustExec(t, db.Save(&a), "preparing account")

		// Execute
		req := testutils.MakeReq(server, "POST", "/verification-token", "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusGone, "Status code mismatch")

		var account database.Account
		var tokenCount int
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, true, "email_verified should not have been updated")
		assert.Equal(t, tokenCount, 0, "token count mismatch")
	})
}

func TestVerifyEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		user := testutils.SetupUserData()
		testutils.SetupAccountData(user, "alice@example.com", "pass1234")
		tok := database.Token{
			UserID: user.ID,
			Type:   database.TokenTypeEmailVerification,
			Value:  "someTokenValue",
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		dat := `{"token": "someTokenValue"}`
		req := testutils.MakeReq(server, "PATCH", "/verify-email", dat)

		// Execute
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismatch")

		var account database.Account
		var token database.Token
		var tokenCount int
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, true, "email_verified mismatch")
		assert.NotEqual(t, token.Value, "", "token value should not have been updated")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.NotEqual(t, token.UsedAt, (*time.Time)(nil), "token should have been used")
	})

	t.Run("used token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		user := testutils.SetupUserData()
		testutils.SetupAccountData(user, "alice@example.com", "pass1234")

		usedAt := time.Now().Add(time.Hour * -11).UTC()
		tok := database.Token{
			UserID: user.ID,
			Type:   database.TokenTypeEmailVerification,
			Value:  "someTokenValue",
			UsedAt: &usedAt,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		dat := `{"token": "someTokenValue"}`
		req := testutils.MakeReq(server, "PATCH", "/verify-email", dat)

		// Execute
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "")

		var account database.Account
		var token database.Token
		var tokenCount int
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, false, "email_verified mismatch")
		assert.NotEqual(t, token.UsedAt, nil, "token used_at mismatch")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.NotEqual(t, token.UsedAt, (*time.Time)(nil), "token should have been used")
	})

	t.Run("expired token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		user := testutils.SetupUserData()
		testutils.SetupAccountData(user, "alice@example.com", "pass1234")

		tok := database.Token{
			UserID: user.ID,
			Type:   database.TokenTypeEmailVerification,
			Value:  "someTokenValue",
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")
		testutils.MustExec(t, db.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-31)), "Failed to prepare token created_at")

		dat := `{"token": "someTokenValue"}`
		req := testutils.MakeReq(server, "PATCH", "/verify-email", dat)

		// Execute
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusGone, "")

		var account database.Account
		var token database.Token
		var tokenCount int
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, false, "email_verified mismatch")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.Equal(t, token.UsedAt, (*time.Time)(nil), "token should have not been used")
	})

	t.Run("already verified", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		user := testutils.SetupUserData()
		a := testutils.SetupAccountData(user, "alice@example.com", "oldpass1234")
		a.EmailVerified = true
		testutils.MustExec(t, db.Save(&a), "preparing account")

		tok := database.Token{
			UserID: user.ID,
			Type:   database.TokenTypeEmailVerification,
			Value:  "someTokenValue",
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		dat := `{"token": "someTokenValue"}`
		req := testutils.MakeReq(server, "PATCH", "/verify-email", dat)

		// Execute
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusConflict, "")

		var account database.Account
		var token database.Token
		var tokenCount int
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, true, "email_verified mismatch")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.Equal(t, token.UsedAt, (*time.Time)(nil), "token should have not been used")
	})
}

func TestUpdateEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		a := testutils.SetupAccountData(u, "alice@example.com", "pass1234")
		a.EmailVerified = true
		testutils.MustExec(t, db.Save(&a), "updating email_verified")

		// Execute
		dat := `{"email": "alice-new@example.com"}`
		req := testutils.MakeReq(server, "PATCH", "/account/profile", dat)
		res := testutils.HTTPAuthDo(t, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var user database.User
		var account database.Account
		testutils.MustExec(t, db.Where("id = ?", u.ID).First(&user), "finding user")
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&account), "finding account")

		assert.Equal(t, account.Email.String, "alice-new@example.com", "email mismatch")
		assert.Equal(t, account.EmailVerified, false, "EmailVerified mismatch")
	})
}

func TestUpdateEmailPreference(t *testing.T) {
	t.Run("with login", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupEmailPreferenceData(u, false)

		// Execute
		dat := `{"digest_weekly": true}`
		req := testutils.MakeReq(server, "PATCH", "/account/email-preference", dat)
		res := testutils.HTTPAuthDo(t, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var preference database.EmailPreference
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&preference), "finding account")
		assert.Equal(t, preference.DigestWeekly, true, "preference mismatch")
	})

	t.Run("with an unused token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupEmailPreferenceData(u, false)
		tok := database.Token{
			UserID: u.ID,
			Type:   database.TokenTypeEmailPreference,
			Value:  "someTokenValue",
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		// Execute
		dat := `{"digest_weekly": true}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someTokenValue")
		req := testutils.MakeReq(server, "PATCH", url, dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var preference database.EmailPreference
		var preferenceCount int
		var token database.Token
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		testutils.MustExec(t, db.Model(database.EmailPreference{}).Count(&preferenceCount), "counting preference")
		testutils.MustExec(t, db.Where("id = ?", tok.ID).First(&token), "failed to find token")

		assert.Equal(t, preferenceCount, 1, "preference count mismatch")
		assert.Equal(t, preference.DigestWeekly, true, "email mismatch")
		assert.NotEqual(t, token.UsedAt, (*time.Time)(nil), "token should have been used")
	})

	t.Run("with nonexistent token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupEmailPreferenceData(u, true)
		tok := database.Token{
			UserID: u.ID,
			Type:   database.TokenTypeEmailPreference,
			Value:  "someTokenValue",
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		dat := `{"digest_weekly": false}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someNonexistentToken")
		req := testutils.MakeReq(server, "PATCH", url, dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var preference database.EmailPreference
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.DigestWeekly, true, "email mismatch")
	})

	t.Run("with expired token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupEmailPreferenceData(u, true)

		usedAt := time.Now().Add(-11 * time.Minute)
		tok := database.Token{
			UserID: u.ID,
			Type:   database.TokenTypeEmailPreference,
			Value:  "someTokenValue",
			UsedAt: &usedAt,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		// Execute
		dat := `{"digest_weekly": false}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someTokenValue")
		req := testutils.MakeReq(server, "PATCH", url, dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var preference database.EmailPreference
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.DigestWeekly, true, "email mismatch")
	})

	t.Run("with a used but unexpired token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupEmailPreferenceData(u, true)
		usedAt := time.Now().Add(-9 * time.Minute)
		tok := database.Token{
			UserID: u.ID,
			Type:   database.TokenTypeEmailPreference,
			Value:  "someTokenValue",
			UsedAt: &usedAt,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		dat := `{"digest_weekly": false}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someTokenValue")
		req := testutils.MakeReq(server, "PATCH", url, dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var preference database.EmailPreference
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.DigestWeekly, false, "DigestWeekly mismatch")
	})

	t.Run("no user and no token", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupEmailPreferenceData(u, true)

		// Execute
		dat := `{"digest_weekly": false}`
		req := testutils.MakeReq(server, "PATCH", "/account/email-preference", dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var preference database.EmailPreference
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.DigestWeekly, true, "email mismatch")
	})

	t.Run("create a record if not exists", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		u := testutils.SetupUserData()
		tok := database.Token{
			UserID: u.ID,
			Type:   database.TokenTypeEmailPreference,
			Value:  "someTokenValue",
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		// Execute
		dat := `{"digest_weekly": false}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someTokenValue")
		req := testutils.MakeReq(server, "PATCH", url, dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var preferenceCount int
		testutils.MustExec(t, db.Model(database.EmailPreference{}).Count(&preferenceCount), "counting preference")
		assert.Equal(t, preferenceCount, 1, "preference count mismatch")

		var preference database.EmailPreference
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.DigestWeekly, false, "email mismatch")
	})
}

func TestGetEmailPreference(t *testing.T) {
	defer testutils.ClearData()

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	u := testutils.SetupUserData()
	pref := testutils.SetupEmailPreferenceData(u, true)

	// Execute
	req := testutils.MakeReq(server, "GET", "/account/email-preference", "")
	res := testutils.HTTPAuthDo(t, req, u)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var got presenters.EmailPreference
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	expected := presenters.EmailPreference{
		DigestWeekly: pref.DigestWeekly,
		CreatedAt:    presenters.FormatTS(pref.CreatedAt),
		UpdatedAt:    presenters.FormatTS(pref.UpdatedAt),
	}
	assert.DeepEqual(t, got, expected, "payload mismatch")
}
