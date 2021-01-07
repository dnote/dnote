/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/models"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdatePassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		user := models.SetUpUserData()
		models.SetUpAccountData(user, "alice@example.com", "oldpassword")

		// Execute
		dat := `{"old_password": "oldpassword", "new_password": "newpassword"}`
		req := models.MakeReq(server.URL, "PATCH", "/account/password", dat)
		res := models.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismsatch")

		var account models.Account
		models.MustExec(t, models.TestDB.Where("user_id = ?", user.ID).First(&account), "finding account")

		passwordErr := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte("newpassword"))
		assert.Equal(t, passwordErr, nil, "Password mismatch")
	})

	t.Run("old password mismatch", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, "alice@example.com", "oldpassword")

		// Execute
		dat := `{"old_password": "randompassword", "new_password": "newpassword"}`
		req := models.MakeReq(server.URL, "PATCH", "/account/password", dat)
		res := models.HTTPAuthDo(t, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "Status code mismsatch")

		var account models.Account
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&account), "finding account")
		assert.Equal(t, a.Password.String, account.Password.String, "password should not have been updated")
	})

	t.Run("password too short", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, "alice@example.com", "oldpassword")

		// Execute
		dat := `{"old_password": "oldpassword", "new_password": "a"}`
		req := models.MakeReq(server.URL, "PATCH", "/account/password", dat)
		res := models.HTTPAuthDo(t, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismsatch")

		var account models.Account
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&account), "finding account")
		assert.Equal(t, a.Password.String, account.Password.String, "password should not have been updated")
	})
}

func TestCreateVerificationToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		emailBackend := testutils.MockEmailbackendImplementation{}
		server := MustNewServer(t, &app.App{
			Clock:        clock.NewMock(),
			EmailBackend: &emailBackend,
		})
		defer server.Close()

		user := models.SetUpUserData()
		models.SetUpAccountData(user, "alice@example.com", "pass1234")

		// Execute
		req := models.MakeReq(server.URL, "POST", "/verification-token", "")
		res := models.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusCreated, "status code mismatch")

		var account models.Account
		var token models.Token
		var tokenCount int
		models.MustExec(t, models.TestDB.Where("user_id = ?", user.ID).First(&account), "finding account")
		models.MustExec(t, models.TestDB.Where("user_id = ? AND type = ?", user.ID, models.TokenTypeEmailVerification).First(&token), "finding token")
		models.MustExec(t, models.TestDB.Model(&models.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, false, "email_verified should not have been updated")
		assert.NotEqual(t, token.Value, "", "token Value mismatch")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.Equal(t, token.UsedAt, (*time.Time)(nil), "token UsedAt mismatch")
		assert.Equal(t, len(emailBackend.Emails), 1, "email queue count mismatch")
	})

	t.Run("already verified", func(t *testing.T) {

		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		user := models.SetUpUserData()
		a := models.SetUpAccountData(user, "alice@example.com", "pass1234")
		a.EmailVerified = true
		models.MustExec(t, models.TestDB.Save(&a), "preparing account")

		// Execute
		req := models.MakeReq(server.URL, "POST", "/verification-token", "")
		res := models.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusGone, "Status code mismatch")

		var account models.Account
		var tokenCount int
		models.MustExec(t, models.TestDB.Where("user_id = ?", user.ID).First(&account), "finding account")
		models.MustExec(t, models.TestDB.Model(&models.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, true, "email_verified should not have been updated")
		assert.Equal(t, tokenCount, 0, "token count mismatch")
	})
}

func TestVerifyEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		user := models.SetUpUserData()
		models.SetUpAccountData(user, "alice@example.com", "pass1234")
		tok := models.Token{
			UserID: user.ID,
			Type:   models.TokenTypeEmailVerification,
			Value:  "someTokenValue",
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")

		dat := `{"token": "someTokenValue"}`
		req := models.MakeReq(server.URL, "PATCH", "/verify-email", dat)

		// Execute
		res := models.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismatch")

		var account models.Account
		var token models.Token
		var tokenCount int
		models.MustExec(t, models.TestDB.Where("user_id = ?", user.ID).First(&account), "finding account")
		models.MustExec(t, models.TestDB.Where("user_id = ? AND type = ?", user.ID, models.TokenTypeEmailVerification).First(&token), "finding token")
		models.MustExec(t, models.TestDB.Model(&models.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, true, "email_verified mismatch")
		assert.NotEqual(t, token.Value, "", "token value should not have been updated")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.NotEqual(t, token.UsedAt, (*time.Time)(nil), "token should have been used")
	})

	t.Run("used token", func(t *testing.T) {

		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		user := models.SetUpUserData()
		models.SetUpAccountData(user, "alice@example.com", "pass1234")

		usedAt := time.Now().Add(time.Hour * -11).UTC()
		tok := models.Token{
			UserID: user.ID,
			Type:   models.TokenTypeEmailVerification,
			Value:  "someTokenValue",
			UsedAt: &usedAt,
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")

		dat := `{"token": "someTokenValue"}`
		req := models.MakeReq(server.URL, "PATCH", "/verify-email", dat)

		// Execute
		res := models.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "")

		var account models.Account
		var token models.Token
		var tokenCount int
		models.MustExec(t, models.TestDB.Where("user_id = ?", user.ID).First(&account), "finding account")
		models.MustExec(t, models.TestDB.Where("user_id = ? AND type = ?", user.ID, models.TokenTypeEmailVerification).First(&token), "finding token")
		models.MustExec(t, models.TestDB.Model(&models.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, false, "email_verified mismatch")
		assert.NotEqual(t, token.UsedAt, nil, "token used_at mismatch")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.NotEqual(t, token.UsedAt, (*time.Time)(nil), "token should have been used")
	})

	t.Run("expired token", func(t *testing.T) {

		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		user := models.SetUpUserData()
		models.SetUpAccountData(user, "alice@example.com", "pass1234")

		tok := models.Token{
			UserID: user.ID,
			Type:   models.TokenTypeEmailVerification,
			Value:  "someTokenValue",
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")
		models.MustExec(t, models.TestDB.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-31)), "Failed to prepare token created_at")

		dat := `{"token": "someTokenValue"}`
		req := models.MakeReq(server.URL, "PATCH", "/verify-email", dat)

		// Execute
		res := models.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusGone, "")

		var account models.Account
		var token models.Token
		var tokenCount int
		models.MustExec(t, models.TestDB.Where("user_id = ?", user.ID).First(&account), "finding account")
		models.MustExec(t, models.TestDB.Where("user_id = ? AND type = ?", user.ID, models.TokenTypeEmailVerification).First(&token), "finding token")
		models.MustExec(t, models.TestDB.Model(&models.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, false, "email_verified mismatch")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.Equal(t, token.UsedAt, (*time.Time)(nil), "token should have not been used")
	})

	t.Run("already verified", func(t *testing.T) {

		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		user := models.SetUpUserData()
		a := models.SetUpAccountData(user, "alice@example.com", "oldpass1234")
		a.EmailVerified = true
		models.MustExec(t, models.TestDB.Save(&a), "preparing account")

		tok := models.Token{
			UserID: user.ID,
			Type:   models.TokenTypeEmailVerification,
			Value:  "someTokenValue",
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")

		dat := `{"token": "someTokenValue"}`
		req := models.MakeReq(server.URL, "PATCH", "/verify-email", dat)

		// Execute
		res := models.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusConflict, "")

		var account models.Account
		var token models.Token
		var tokenCount int
		models.MustExec(t, models.TestDB.Where("user_id = ?", user.ID).First(&account), "finding account")
		models.MustExec(t, models.TestDB.Where("user_id = ? AND type = ?", user.ID, models.TokenTypeEmailVerification).First(&token), "finding token")
		models.MustExec(t, models.TestDB.Model(&models.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, true, "email_verified mismatch")
		assert.Equal(t, tokenCount, 1, "token count mismatch")
		assert.Equal(t, token.UsedAt, (*time.Time)(nil), "token should have not been used")
	})
}

func TestUpdateEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, "alice@example.com", "pass1234")
		a.EmailVerified = true
		models.MustExec(t, models.TestDB.Save(&a), "updating email_verified")

		// Execute
		dat := `{"email": "alice-new@example.com", "password": "pass1234"}`
		req := models.MakeReq(server.URL, "PATCH", "/account/profile", dat)
		res := models.HTTPAuthDo(t, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var user models.User
		var account models.Account
		models.MustExec(t, models.TestDB.Where("id = ?", u.ID).First(&user), "finding user")
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&account), "finding account")

		assert.Equal(t, account.Email.String, "alice-new@example.com", "email mismatch")
		assert.Equal(t, account.EmailVerified, false, "EmailVerified mismatch")
	})

	t.Run("password mismatch", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, "alice@example.com", "pass1234")
		a.EmailVerified = true
		models.MustExec(t, models.TestDB.Save(&a), "updating email_verified")

		// Execute
		dat := `{"email": "alice-new@example.com", "password": "wrongpassword"}`
		req := models.MakeReq(server.URL, "PATCH", "/account/profile", dat)
		res := models.HTTPAuthDo(t, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "Status code mismsatch")

		var user models.User
		var account models.Account
		models.MustExec(t, models.TestDB.Where("id = ?", u.ID).First(&user), "finding user")
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&account), "finding account")

		assert.Equal(t, account.Email.String, "alice@example.com", "email mismatch")
		assert.Equal(t, account.EmailVerified, true, "EmailVerified mismatch")
	})
}

func TestUpdateEmailPreference(t *testing.T) {
	t.Run("with login", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		models.SetupEmailPreferenceData(u, false)

		// Execute
		dat := `{"inactive_reminder": true}`
		req := models.MakeReq(server.URL, "PATCH", "/account/email-preference", dat)
		res := models.HTTPAuthDo(t, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var preference models.EmailPreference
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&preference), "finding account")
		assert.Equal(t, preference.InactiveReminder, true, "preference mismatch")
	})

	t.Run("with an unused token", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		models.SetupEmailPreferenceData(u, false)
		tok := models.Token{
			UserID: u.ID,
			Type:   models.TokenTypeEmailPreference,
			Value:  "someTokenValue",
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")

		// Execute
		dat := `{"inactive_reminder": true}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someTokenValue")
		req := models.MakeReq(server.URL, "PATCH", url, dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var preference models.EmailPreference
		var preferenceCount int
		var token models.Token
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		models.MustExec(t, models.TestDB.Model(models.EmailPreference{}).Count(&preferenceCount), "counting preference")
		models.MustExec(t, models.TestDB.Where("id = ?", tok.ID).First(&token), "failed to find token")

		assert.Equal(t, preferenceCount, 1, "preference count mismatch")
		assert.Equal(t, preference.InactiveReminder, true, "email mismatch")
		assert.NotEqual(t, token.UsedAt, (*time.Time)(nil), "token should have been used")
	})

	t.Run("with nonexistent token", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		models.SetupEmailPreferenceData(u, true)
		tok := models.Token{
			UserID: u.ID,
			Type:   models.TokenTypeEmailPreference,
			Value:  "someTokenValue",
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")

		dat := `{"inactive_reminder": false}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someNonexistentToken")
		req := models.MakeReq(server.URL, "PATCH", url, dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var preference models.EmailPreference
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.InactiveReminder, true, "email mismatch")
	})

	t.Run("with expired token", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		models.SetupEmailPreferenceData(u, true)

		usedAt := time.Now().Add(-11 * time.Minute)
		tok := models.Token{
			UserID: u.ID,
			Type:   models.TokenTypeEmailPreference,
			Value:  "someTokenValue",
			UsedAt: &usedAt,
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")

		// Execute
		dat := `{"inactive_reminder": false}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someTokenValue")
		req := models.MakeReq(server.URL, "PATCH", url, dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var preference models.EmailPreference
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.InactiveReminder, true, "email mismatch")
	})

	t.Run("with a used but unexpired token", func(t *testing.T) {

		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		models.SetupEmailPreferenceData(u, true)
		usedAt := time.Now().Add(-9 * time.Minute)
		tok := models.Token{
			UserID: u.ID,
			Type:   models.TokenTypeEmailPreference,
			Value:  "someTokenValue",
			UsedAt: &usedAt,
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")

		dat := `{"inactive_reminder": false}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someTokenValue")
		req := models.MakeReq(server.URL, "PATCH", url, dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var preference models.EmailPreference
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.InactiveReminder, false, "InactiveReminder mismatch")
	})

	t.Run("no user and no token", func(t *testing.T) {

		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		models.SetupEmailPreferenceData(u, true)

		// Execute
		dat := `{"inactive_reminder": false}`
		req := models.MakeReq(server.URL, "PATCH", "/account/email-preference", dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var preference models.EmailPreference
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.InactiveReminder, true, "email mismatch")
	})

	t.Run("create a record if not exists", func(t *testing.T) {

		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		tok := models.Token{
			UserID: u.ID,
			Type:   models.TokenTypeEmailPreference,
			Value:  "someTokenValue",
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")

		// Execute
		dat := `{"inactive_reminder": false}`
		url := fmt.Sprintf("/account/email-preference?token=%s", "someTokenValue")
		req := models.MakeReq(server.URL, "PATCH", url, dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var preferenceCount int
		models.MustExec(t, models.TestDB.Model(models.EmailPreference{}).Count(&preferenceCount), "counting preference")
		assert.Equal(t, preferenceCount, 1, "preference count mismatch")

		var preference models.EmailPreference
		models.MustExec(t, models.TestDB.Where("user_id = ?", u.ID).First(&preference), "finding preference")
		assert.Equal(t, preference.InactiveReminder, false, "email mismatch")
	})
}

func TestGetEmailPreference(t *testing.T) {
	defer models.ClearTestData(models.TestDB)
	// Setup
	server := MustNewServer(t, &app.App{

		Clock: clock.NewMock(),
	})
	defer server.Close()

	u := models.SetUpUserData()
	pref := models.SetupEmailPreferenceData(u, true)

	// Execute
	req := models.MakeReq(server.URL, "GET", "/account/email-preference", "")
	res := models.HTTPAuthDo(t, req, u)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var got presenters.EmailPreference
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	expected := presenters.EmailPreference{
		InactiveReminder: pref.InactiveReminder,
		ProductUpdate:    pref.ProductUpdate,
		CreatedAt:        presenters.FormatTS(pref.CreatedAt),
		UpdatedAt:        presenters.FormatTS(pref.UpdatedAt),
	}
	assert.DeepEqual(t, got, expected, "payload mismatch")
}
