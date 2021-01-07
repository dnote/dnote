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
	"github.com/dnote/dnote/pkg/server/session"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func TestGetMe(t *testing.T) {
	models.InitTestDB()
	defer models.ClearTestData(models.TestDB)

	// Setup
	server := MustNewServer(t, &app.App{
		Clock: clock.NewMock(),
	})
	defer server.Close()

	u1 := models.SetUpUserData()
	a1 := models.SetUpAccountData(u1, "alice@example.com", "somepassword")

	u2 := models.SetUpUserData()
	models.MustExec(t, models.TestDB.Model(&u2).Update("cloud", false), "preparing u2 cloud")
	a2 := models.SetUpAccountData(u2, "bob@example.com", "somepassword")

	testCases := []struct {
		user        models.User
		account     models.Account
		expectedPro bool
	}{
		{
			user:        u1,
			account:     a1,
			expectedPro: true,
		},
		{
			user:        u2,
			account:     a2,
			expectedPro: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("user pro %t", tc.expectedPro), func(t *testing.T) {
			// Execute
			req := testutils.MakeReq(server.URL, "GET", "/me", "")
			res := models.HTTPAuthDo(t, req, tc.user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, "")

			var payload GetMeResponse
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			expectedPayload := GetMeResponse{
				User: session.Session{
					UUID:          tc.user.UUID,
					Pro:           tc.expectedPro,
					Email:         tc.account.Email.String,
					EmailVerified: tc.account.EmailVerified,
				},
			}
			assert.DeepEqual(t, payload, expectedPayload, "payload mismatch")

			var user models.User
			models.MustExec(t, models.TestDB.Where("id = ?", tc.user.ID).First(&user), "finding user")
			assert.NotEqual(t, user.LastLoginAt, nil, "LastLoginAt mismatch")
		})
	}
}

func TestCreateResetToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		models.SetUpAccountData(u, "alice@example.com", "somepassword")

		dat := `{"email": "alice@example.com"}`
		req := testutils.MakeReq(server.URL, "POST", "/reset-token", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismtach")

		var tokenCount int
		models.MustExec(t, models.TestDB.Model(&models.Token{}).Count(&tokenCount), "counting tokens")

		var resetToken models.Token
		models.MustExec(t, models.TestDB.Where("user_id = ? AND type = ?", u.ID, models.TokenTypeResetPassword).First(&resetToken), "finding reset token")

		assert.Equal(t, tokenCount, 1, "reset_token count mismatch")
		assert.NotEqual(t, resetToken.Value, nil, "reset_token value mismatch")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "reset_token UsedAt mismatch")
	})

	t.Run("nonexistent email", func(t *testing.T) {

		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		models.SetUpAccountData(u, "alice@example.com", "somepassword")

		dat := `{"email": "bob@example.com"}`
		req := testutils.MakeReq(server.URL, "POST", "/reset-token", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismtach")

		var tokenCount int
		models.MustExec(t, models.TestDB.Model(&models.Token{}).Count(&tokenCount), "counting tokens")
		assert.Equal(t, tokenCount, 0, "reset_token count mismatch")
	})
}

func TestResetPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, "alice@example.com", "oldpassword")
		tok := models.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   models.TokenTypeResetPassword,
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")
		otherTok := models.Token{
			UserID: u.ID,
			Value:  "somerandomvalue",
			Type:   models.TokenTypeEmailVerification,
		}
		models.MustExec(t, models.TestDB.Save(&otherTok), "preparing another token")

		dat := `{"token": "MivFxYiSMMA4An9dP24DNQ==", "password": "newpassword"}`
		req := testutils.MakeReq(server.URL, "PATCH", "/reset-password", dat)

		s1 := models.Session{
			Key:       "some-session-key-1",
			UserID:    u.ID,
			ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
		}
		models.MustExec(t, models.TestDB.Save(&s1), "preparing user session 1")

		s2 := &models.Session{
			Key:       "some-session-key-2",
			UserID:    u.ID,
			ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
		}
		models.MustExec(t, models.TestDB.Save(&s2), "preparing user session 2")

		anotherUser := models.SetUpUserData()
		models.MustExec(t, models.TestDB.Save(&models.Session{
			Key:       "some-session-key-3",
			UserID:    anotherUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
		}), "preparing anotherUser session 1")

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismatch")

		var resetToken, verificationToken models.Token
		var account models.Account
		models.MustExec(t, models.TestDB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "finding reset token")
		models.MustExec(t, models.TestDB.Where("value = ?", "somerandomvalue").First(&verificationToken), "finding reset token")
		models.MustExec(t, models.TestDB.Where("id = ?", a.ID).First(&account), "finding account")

		assert.NotEqual(t, resetToken.UsedAt, nil, "reset_token UsedAt mismatch")
		passwordErr := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte("newpassword"))
		assert.Equal(t, passwordErr, nil, "Password mismatch")
		assert.Equal(t, verificationToken.UsedAt, (*time.Time)(nil), "verificationToken UsedAt mismatch")

		var s1Count, s2Count int
		models.MustExec(t, models.TestDB.Model(&models.Session{}).Where("id = ?", s1.ID).Count(&s1Count), "counting s1")
		models.MustExec(t, models.TestDB.Model(&models.Session{}).Where("id = ?", s2.ID).Count(&s2Count), "counting s2")

		assert.Equal(t, s1Count, 0, "s1 should have been deleted")
		assert.Equal(t, s2Count, 0, "s2 should have been deleted")

		var userSessionCount, anotherUserSessionCount int
		models.MustExec(t, models.TestDB.Model(&models.Session{}).Where("user_id = ?", u.ID).Count(&userSessionCount), "counting user session")
		models.MustExec(t, models.TestDB.Model(&models.Session{}).Where("user_id = ?", anotherUser.ID).Count(&anotherUserSessionCount), "counting anotherUser session")

		assert.Equal(t, userSessionCount, 1, "should have created a new user session")
		assert.Equal(t, anotherUserSessionCount, 1, "anotherUser session count mismatch")
	})

	t.Run("nonexistent token", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, "alice@example.com", "somepassword")
		tok := models.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   models.TokenTypeResetPassword,
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")

		dat := `{"token": "-ApMnyvpg59uOU5b-Kf5uQ==", "password": "oldpassword"}`
		req := testutils.MakeReq(server.URL, "PATCH", "/reset-password", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken models.Token
		var account models.Account
		models.MustExec(t, models.TestDB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "finding reset token")
		models.MustExec(t, models.TestDB.Where("id = ?", a.ID).First(&account), "finding account")

		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})

	t.Run("expired token", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, "alice@example.com", "somepassword")
		tok := models.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   models.TokenTypeResetPassword,
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")
		models.MustExec(t, models.TestDB.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := `{"token": "MivFxYiSMMA4An9dP24DNQ==", "password": "oldpassword"}`
		req := testutils.MakeReq(server.URL, "PATCH", "/reset-password", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusGone, "Status code mismatch")

		var resetToken models.Token
		var account models.Account
		models.MustExec(t, models.TestDB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		models.MustExec(t, models.TestDB.Where("id = ?", a.ID).First(&account), "failed to find account")
		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})

	t.Run("used token", func(t *testing.T) {
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, "alice@example.com", "somepassword")

		usedAt := time.Now().Add(time.Hour * -11).UTC()
		tok := models.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   models.TokenTypeResetPassword,
			UsedAt: &usedAt,
		}
		models.MustExec(t, models.TestDB.Save(&tok), "preparing token")
		models.MustExec(t, models.TestDB.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := `{"token": "MivFxYiSMMA4An9dP24DNQ==", "password": "oldpassword"}`
		req := testutils.MakeReq(server.URL, "PATCH", "/reset-password", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken models.Token
		var account models.Account
		models.MustExec(t, models.TestDB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		models.MustExec(t, models.TestDB.Where("id = ?", a.ID).First(&account), "failed to find account")
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
		defer models.ClearTestData(models.TestDB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, "alice@example.com", "somepassword")
		tok := models.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   models.TokenTypeEmailVerification,
		}
		models.MustExec(t, models.TestDB.Save(&tok), "Failed to prepare reset_token")
		models.MustExec(t, models.TestDB.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := `{"token": "MivFxYiSMMA4An9dP24DNQ==", "password": "oldpassword"}`
		req := testutils.MakeReq(server.URL, "PATCH", "/reset-password", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken models.Token
		var account models.Account
		models.MustExec(t, models.TestDB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		models.MustExec(t, models.TestDB.Where("id = ?", a.ID).First(&account), "failed to find account")

		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})
}
