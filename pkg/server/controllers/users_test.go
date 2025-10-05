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

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func assertResponseSessionCookie(t *testing.T, db *gorm.DB, res *http.Response) {
	var sessionCount int64
	var session database.Session
	testutils.MustExec(t, db.Model(&database.Session{}).Count(&sessionCount), "counting session")
	testutils.MustExec(t, db.First(&session), "getting session")

	c := testutils.GetCookieByName(res.Cookies(), "id")
	assert.Equal(t, c.Value, session.Key, "session key mismatch")
	assert.Equal(t, c.Path, "/", "session path mismatch")
	assert.Equal(t, c.HttpOnly, true, "session HTTPOnly mismatch")
	assert.Equal(t, c.Expires.Unix(), session.ExpiresAt.Unix(), "session Expires mismatch")
}

func TestJoin(t *testing.T) {
	testCases := []struct {
		email                string
		password             string
		passwordConfirmation string
	}{
		{
			email:                "alice@example.com",
			password:             "pass1234",
			passwordConfirmation: "pass1234",
		},
		{
			email:                "bob@example.com",
			password:             "Y9EwmjH@Jq6y5a64MSACUoM4w7SAhzvY",
			passwordConfirmation: "Y9EwmjH@Jq6y5a64MSACUoM4w7SAhzvY",
		},
		{
			email:                "chuck@example.com",
			password:             "e*H@kJi^vXbWEcD9T5^Am!Y@7#Po2@PC",
			passwordConfirmation: "e*H@kJi^vXbWEcD9T5^Am!Y@7#Po2@PC",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("register %s %s", tc.email, tc.password), func(t *testing.T) {
			db := testutils.InitMemoryDB(t)

			// Setup
			emailBackend := testutils.MockEmailbackendImplementation{}
			a := app.NewTest()
			a.Clock = clock.NewMock()
			a.EmailBackend = &emailBackend
			a.DB = db
			server := MustNewServer(t, &a)
			defer server.Close()

			dat := url.Values{}
			dat.Set("email", tc.email)
			dat.Set("password", tc.password)
			dat.Set("password_confirmation", tc.passwordConfirmation)
			req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

			// Execute
			res := testutils.HTTPDo(t, req)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusFound, "")

			var account database.Account
			testutils.MustExec(t, db.Where("email = ?", tc.email).First(&account), "finding account")
			assert.Equal(t, account.Email.String, tc.email, "Email mismatch")
			assert.NotEqual(t, account.UserID, 0, "UserID mismatch")
			passwordErr := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte(tc.password))
			assert.Equal(t, passwordErr, nil, "Password mismatch")

			var user database.User
			testutils.MustExec(t, db.Where("id = ?", account.UserID).First(&user), "finding user")
			assert.Equal(t, user.MaxUSN, 0, "MaxUSN mismatch")

			// welcome email
			assert.Equalf(t, len(emailBackend.Emails), 1, "email queue count mismatch")
			assert.DeepEqual(t, emailBackend.Emails[0].To, []string{tc.email}, "email to mismatch")

			// after register, should sign in user
			assertResponseSessionCookie(t, db, res)
		})
	}
}

func TestJoinError(t *testing.T) {
	t.Run("missing email", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		dat := url.Values{}
		dat.Set("password", "SLMZFM5RmSjA5vfXnG5lPOnrpZSbtmV76cnAcrlr2yU")
		req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status mismatch")

		var accountCount, userCount int64
		testutils.MustExec(t, db.Model(&database.Account{}).Count(&accountCount), "counting account")
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting user")

		assert.Equal(t, accountCount, int64(0), "accountCount mismatch")
		assert.Equal(t, userCount, int64(0), "userCount mismatch")
	})

	t.Run("missing password", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		dat := url.Values{}
		dat.Set("email", "alice@example.com")
		req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status mismatch")

		var accountCount, userCount int64
		testutils.MustExec(t, db.Model(&database.Account{}).Count(&accountCount), "counting account")
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting user")

		assert.Equal(t, accountCount, int64(0), "accountCount mismatch")
		assert.Equal(t, userCount, int64(0), "userCount mismatch")
	})

	t.Run("password confirmation mismatch", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		dat := url.Values{}
		dat.Set("email", "alice@example.com")
		dat.Set("password", "pass1234")
		dat.Set("password_confirmation", "1234pass")
		req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status mismatch")

		var accountCount, userCount int64
		testutils.MustExec(t, db.Model(&database.Account{}).Count(&accountCount), "counting account")
		testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting user")

		assert.Equal(t, accountCount, int64(0), "accountCount mismatch")
		assert.Equal(t, userCount, int64(0), "userCount mismatch")
	})
}

func TestJoinDuplicateEmail(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	// Setup
	a := app.NewTest()
	a.Clock = clock.NewMock()
	a.DB = db
	server := MustNewServer(t, &a)
	defer server.Close()

	u := testutils.SetupUserData(db)
	testutils.SetupAccountData(db, u, "alice@example.com", "somepassword")

	dat := url.Values{}
	dat.Set("email", "alice@example.com")
	dat.Set("password", "foobarbaz")
	dat.Set("password_confirmation", "foobarbaz")
	req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

	// Execute
	res := testutils.HTTPDo(t, req)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusBadRequest, "status code mismatch")

	var accountCount, userCount, verificationTokenCount int64
	testutils.MustExec(t, db.Model(&database.Account{}).Count(&accountCount), "counting account")
	testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting user")
	testutils.MustExec(t, db.Model(&database.Token{}).Count(&verificationTokenCount), "counting verification token")

	var user database.User
	testutils.MustExec(t, db.Where("id = ?", u.ID).First(&user), "finding user")

	assert.Equal(t, accountCount, int64(1), "account count mismatch")
	assert.Equal(t, userCount, int64(1), "user count mismatch")
	assert.Equal(t, verificationTokenCount, int64(0), "verification_token should not have been created")
	assert.Equal(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")
}

func TestJoinDisabled(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	// Setup
	a := app.NewTest()
	a.Clock = clock.NewMock()
	a.DB = db
	a.DisableRegistration = true
	server := MustNewServer(t, &a)
	defer server.Close()

	dat := url.Values{}
	dat.Set("email", "alice@example.com")
	dat.Set("password", "foobarbaz")
	req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

	// Execute
	res := testutils.HTTPDo(t, req)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusNotFound, "status code mismatch")

	var accountCount, userCount int64
	testutils.MustExec(t, db.Model(&database.Account{}).Count(&accountCount), "counting account")
	testutils.MustExec(t, db.Model(&database.User{}).Count(&userCount), "counting user")

	assert.Equal(t, accountCount, int64(0), "account count mismatch")
	assert.Equal(t, userCount, int64(0), "user count mismatch")
}

func TestLogin(t *testing.T) {
	testutils.RunForWebAndAPI(t, "success", func(t *testing.T, target testutils.EndpointType) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)

		u := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, u, "alice@example.com", "pass1234")
		defer server.Close()

		// Execute
		var req *http.Request
		if target == testutils.EndpointWeb {
			dat := url.Values{}
			dat.Set("email", "alice@example.com")
			dat.Set("password", "pass1234")
			req = testutils.MakeFormReq(server.URL, "POST", "/login", dat)
		} else {
			dat := `{"email": "alice@example.com", "password": "pass1234"}`
			req = testutils.MakeReq(server.URL, "POST", "/api/v3/signin", dat)
		}

		res := testutils.HTTPDo(t, req)

		// Test
		if target == testutils.EndpointWeb {
			assert.StatusCodeEquals(t, res, http.StatusFound, "")
		} else {
			assert.StatusCodeEquals(t, res, http.StatusOK, "")
		}

		var user database.User
		testutils.MustExec(t, db.Model(&database.User{}).First(&user), "finding user")
		assert.NotEqual(t, user.LastLoginAt, nil, "LastLoginAt mismatch")

		if target == testutils.EndpointWeb {
			assertResponseSessionCookie(t, db, res)
		} else {
			// after register, should sign in user
			var got SessionResponse
			if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			var sessionCount int64
			var session database.Session
			testutils.MustExec(t, db.Model(&database.Session{}).Count(&sessionCount), "counting session")
			testutils.MustExec(t, db.First(&session), "getting session")

			assert.Equal(t, sessionCount, int64(1), "sessionCount mismatch")
			assert.Equal(t, got.Key, session.Key, "session Key mismatch")
			assert.Equal(t, got.ExpiresAt, session.ExpiresAt.Unix(), "session ExpiresAt mismatch")

			assertResponseSessionCookie(t, db, res)
		}
	})

	testutils.RunForWebAndAPI(t, "wrong password", func(t *testing.T, target testutils.EndpointType) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)

		u := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, u, "alice@example.com", "pass1234")
		defer server.Close()

		var req *http.Request
		if target == testutils.EndpointWeb {
			dat := url.Values{}
			dat.Set("email", "alice@example.com")
			dat.Set("password", "wrongpassword1234")
			req = testutils.MakeFormReq(server.URL, "POST", "/login", dat)
		} else {
			dat := `{"email": "alice@example.com", "password": "wrongpassword1234"}`
			req = testutils.MakeReq(server.URL, "POST", "/api/v3/signin", dat)
		}

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var user database.User
		testutils.MustExec(t, db.Model(&database.User{}).First(&user), "finding user")
		assert.Equal(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")

		var sessionCount int64
		testutils.MustExec(t, db.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, int64(0), "sessionCount mismatch")
	})

	testutils.RunForWebAndAPI(t, "wrong email", func(t *testing.T, target testutils.EndpointType) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, u, "alice@example.com", "pass1234")

		var req *http.Request
		if target == testutils.EndpointWeb {
			dat := url.Values{}
			dat.Set("email", "bob@example.com")
			dat.Set("password", "foobarbaz")
			req = testutils.MakeFormReq(server.URL, "POST", "/login", dat)
		} else {
			dat := `{"email": "bob@example.com", "password": "foobarbaz"}`
			req = testutils.MakeReq(server.URL, "POST", "/api/v3/signin", dat)
		}

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var user database.User
		testutils.MustExec(t, db.Model(&database.User{}).First(&user), "finding user")
		assert.DeepEqual(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")

		var sessionCount int64
		testutils.MustExec(t, db.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, int64(0), "sessionCount mismatch")
	})

	testutils.RunForWebAndAPI(t, "nonexistent email", func(t *testing.T, target testutils.EndpointType) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		var req *http.Request
		if target == testutils.EndpointWeb {
			dat := url.Values{}
			dat.Set("email", "nonexistent@example.com")
			dat.Set("password", "pass1234")
			req = testutils.MakeFormReq(server.URL, "POST", "/login", dat)
		} else {
			dat := `{"email": "nonexistent@example.com", "password": "pass1234"}`
			req = testutils.MakeReq(server.URL, "POST", "/api/v3/signin", dat)
		}

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var sessionCount int64
		testutils.MustExec(t, db.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, int64(0), "sessionCount mismatch")
	})
}

func TestLogout(t *testing.T) {
	setupLogoutTest := func(t *testing.T, db *gorm.DB) (*httptest.Server, *database.Session, *database.Session) {
		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)

		aliceUser := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, aliceUser, "alice@example.com", "pass1234")
		anotherUser := testutils.SetupUserData(db)

		session1ExpiresAt := time.Now().Add(time.Hour * 24)
		session1 := database.Session{
			Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
			UserID:    aliceUser.ID,
			ExpiresAt: session1ExpiresAt,
		}
		testutils.MustExec(t, db.Save(&session1), "preparing session1")
		session2 := database.Session{
			Key:       "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=",
			UserID:    anotherUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}
		testutils.MustExec(t, db.Save(&session2), "preparing session2")

		return server, &session1, &session2
	}

	testutils.RunForWebAndAPI(t, "authenticated", func(t *testing.T, target testutils.EndpointType) {
		db := testutils.InitMemoryDB(t)

		server, session1, _ := setupLogoutTest(t, db)
		defer server.Close()

		// Execute
		var req *http.Request
		if target == testutils.EndpointWeb {
			dat := url.Values{}
			req = testutils.MakeFormReq(server.URL, "POST", "/logout", dat)
			req.AddCookie(&http.Cookie{Name: "id", Value: "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=", Expires: session1.ExpiresAt, Path: "/", HttpOnly: true})
		} else {
			req = testutils.MakeReq(server.URL, "POST", "/api/v3/signout", "")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session1.Key))
		}

		res := testutils.HTTPDo(t, req)

		// Test
		if target == testutils.EndpointWeb {
			assert.StatusCodeEquals(t, res, http.StatusFound, "Status mismatch")
		} else {
			assert.StatusCodeEquals(t, res, http.StatusNoContent, "Status mismatch")
		}

		var sessionCount int64
		var s2 database.Session
		testutils.MustExec(t, db.Model(&database.Session{}).Count(&sessionCount), "counting session")
		testutils.MustExec(t, db.Where("key = ?", "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=").First(&s2), "getting s2")

		assert.Equal(t, sessionCount, int64(1), "sessionCount mismatch")

		if target == testutils.EndpointWeb {
			c := testutils.GetCookieByName(res.Cookies(), "id")
			assert.Equal(t, c.Value, "", "session key mismatch")
			assert.Equal(t, c.Path, "/", "session path mismatch")
			assert.Equal(t, c.HttpOnly, true, "session HTTPOnly mismatch")
			if c.Expires.After(time.Now()) {
				t.Error("session cookie is not expired")
			}
		}
	})

	testutils.RunForWebAndAPI(t, "unauthenticated", func(t *testing.T, target testutils.EndpointType) {
		db := testutils.InitMemoryDB(t)

		server, _, _ := setupLogoutTest(t, db)
		defer server.Close()

		// Execute
		var req *http.Request
		if target == testutils.EndpointWeb {
			dat := url.Values{}
			req = testutils.MakeFormReq(server.URL, "POST", "/logout", dat)
		} else {
			req = testutils.MakeReq(server.URL, "POST", "/api/v3/signout", "")
		}

		res := testutils.HTTPDo(t, req)

		// Test
		if target == testutils.EndpointWeb {
			assert.StatusCodeEquals(t, res, http.StatusFound, "Status mismatch")
		} else {
			assert.StatusCodeEquals(t, res, http.StatusNoContent, "Status mismatch")
		}

		var sessionCount int64
		var postSession1, postSession2 database.Session
		testutils.MustExec(t, db.Model(&database.Session{}).Count(&sessionCount), "counting session")
		testutils.MustExec(t, db.Where("key = ?", "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=").First(&postSession1), "getting postSession1")
		testutils.MustExec(t, db.Where("key = ?", "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=").First(&postSession2), "getting postSession2")

		// two existing sessions should remain
		assert.Equal(t, sessionCount, int64(2), "sessionCount mismatch")

		c := testutils.GetCookieByName(res.Cookies(), "id")
		assert.Equal(t, c, (*http.Cookie)(nil), "id cookie should have not been set")
	})
}

func TestResetPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "oldpassword")
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

		s1 := database.Session{
			Key:       "some-session-key-1",
			UserID:    u.ID,
			ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
		}
		testutils.MustExec(t, db.Save(&s1), "preparing user session 1")

		s2 := &database.Session{
			Key:       "some-session-key-2",
			UserID:    u.ID,
			ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
		}
		testutils.MustExec(t, db.Save(&s2), "preparing user session 2")

		anotherUser := testutils.SetupUserData(db)
		testutils.MustExec(t, db.Save(&database.Session{
			Key:       "some-session-key-3",
			UserID:    anotherUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
		}), "preparing anotherUser session 1")

		// Execute
		dat := url.Values{}
		dat.Set("token", "MivFxYiSMMA4An9dP24DNQ==")
		dat.Set("password", "newpassword")
		dat.Set("password_confirmation", "newpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusFound, "Status code mismatch")

		var resetToken, verificationToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "finding reset token")
		testutils.MustExec(t, db.Where("value = ?", "somerandomvalue").First(&verificationToken), "finding reset token")
		testutils.MustExec(t, db.Where("id = ?", acc.ID).First(&account), "finding account")

		assert.NotEqual(t, resetToken.UsedAt, nil, "reset_token UsedAt mismatch")
		passwordErr := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte("newpassword"))
		assert.Equal(t, passwordErr, nil, "Password mismatch")
		assert.Equal(t, verificationToken.UsedAt, (*time.Time)(nil), "verificationToken UsedAt mismatch")

		var s1Count, s2Count int64
		testutils.MustExec(t, db.Model(&database.Session{}).Where("id = ?", s1.ID).Count(&s1Count), "counting s1")
		testutils.MustExec(t, db.Model(&database.Session{}).Where("id = ?", s2.ID).Count(&s2Count), "counting s2")

		assert.Equal(t, s1Count, int64(0), "s1 should have been deleted")
		assert.Equal(t, s2Count, int64(0), "s2 should have been deleted")

		var userSessionCount, anotherUserSessionCount int64
		testutils.MustExec(t, db.Model(&database.Session{}).Where("user_id = ?", u.ID).Count(&userSessionCount), "counting user session")
		testutils.MustExec(t, db.Model(&database.Session{}).Where("user_id = ?", anotherUser.ID).Count(&anotherUserSessionCount), "counting anotherUser session")

		assert.Equal(t, userSessionCount, int64(0), "should have deleted a user session")
		assert.Equal(t, anotherUserSessionCount, int64(1), "anotherUser session count mismatch")
	})

	t.Run("nonexistent token", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "somepassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		dat := url.Values{}
		dat.Set("token", "-ApMnyvpg59uOU5b-Kf5uQ==")
		dat.Set("password", "oldpassword")
		dat.Set("password_confirmation", "oldpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "finding reset token")
		testutils.MustExec(t, db.Where("id = ?", acc.ID).First(&account), "finding account")

		assert.Equal(t, acc.Password, account.Password, "password should not have been updated")
		assert.Equal(t, acc.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})

	t.Run("expired token", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "somepassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")
		testutils.MustExec(t, db.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := url.Values{}
		dat.Set("token", "MivFxYiSMMA4An9dP24DNQ==")
		dat.Set("password", "oldpassword")
		dat.Set("password_confirmation", "oldpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusGone, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		testutils.MustExec(t, db.Where("id = ?", acc.ID).First(&account), "failed to find account")
		assert.Equal(t, acc.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})

	t.Run("used token", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "somepassword")

		usedAt := time.Now().Add(time.Hour * -11).UTC()
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
			UsedAt: &usedAt,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")
		testutils.MustExec(t, db.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := url.Values{}
		dat.Set("token", "MivFxYiSMMA4An9dP24DNQ==")
		dat.Set("password", "oldpassword")
		dat.Set("password_confirmation", "oldpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		testutils.MustExec(t, db.Where("id = ?", acc.ID).First(&account), "failed to find account")
		assert.Equal(t, acc.Password, account.Password, "password should not have been updated")

		resetTokenUsedAtUTC := resetToken.UsedAt.UTC()
		if resetTokenUsedAtUTC.Year() != usedAt.Year() ||
			resetTokenUsedAtUTC.Month() != usedAt.Month() ||
			resetTokenUsedAtUTC.Day() != usedAt.Day() ||
			resetTokenUsedAtUTC.Hour() != usedAt.Hour() ||
			resetTokenUsedAtUTC.Minute() != usedAt.Minute() ||
			resetTokenUsedAtUTC.Second() != usedAt.Second() {
			t.Errorf("used_at should be %+v but got: %+v", usedAt, resetToken.UsedAt)
		}
	})

	t.Run("using wrong type token: email_verification", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "somepassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeEmailVerification,
		}
		testutils.MustExec(t, db.Save(&tok), "Failed to prepare reset_token")
		testutils.MustExec(t, db.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := url.Values{}
		dat.Set("token", "MivFxYiSMMA4An9dP24DNQ==")
		dat.Set("password", "oldpassword")
		dat.Set("password_confirmation", "oldpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, db.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		testutils.MustExec(t, db.Where("id = ?", acc.ID).First(&account), "failed to find account")

		assert.Equal(t, acc.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})
}

func TestCreateResetToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, u, "alice@example.com", "somepassword")

		// Execute
		dat := url.Values{}
		dat.Set("email", "alice@example.com")
		req := testutils.MakeFormReq(server.URL, "POST", "/reset-token", dat)

		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusFound, "Status code mismtach")

		var tokenCount int64
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting tokens")

		var resetToken database.Token
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", u.ID, database.TokenTypeResetPassword).First(&resetToken), "finding reset token")

		assert.Equal(t, tokenCount, int64(1), "reset_token count mismatch")
		assert.NotEqual(t, resetToken.Value, nil, "reset_token value mismatch")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "reset_token UsedAt mismatch")
	})

	t.Run("nonexistent email", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, u, "alice@example.com", "somepassword")

		// Execute
		dat := url.Values{}
		dat.Set("email", "bob@example.com")
		req := testutils.MakeFormReq(server.URL, "POST", "/reset-token", dat)

		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismtach")

		var tokenCount int64
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting tokens")
		assert.Equal(t, tokenCount, int64(0), "reset_token count mismatch")
	})
}

func TestUpdatePassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		user := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, user, "alice@example.com", "oldpassword")

		// Execute
		dat := url.Values{}
		dat.Set("old_password", "oldpassword")
		dat.Set("new_password", "newpassword")
		dat.Set("new_password_confirmation", "newpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/account/password", dat)

		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusFound, "Status code mismsatch")

		var account database.Account
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")

		passwordErr := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte("newpassword"))
		assert.Equal(t, passwordErr, nil, "Password mismatch")
	})

	t.Run("old password mismatch", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)
		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "oldpassword")

		// Execute
		dat := url.Values{}
		dat.Set("old_password", "randompassword")
		dat.Set("new_password", "newpassword")
		dat.Set("new_password_confirmation", "newpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/account/password", dat)

		res := testutils.HTTPAuthDo(t, db, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "Status code mismsatch")

		var account database.Account
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&account), "finding account")
		assert.Equal(t, acc.Password.String, account.Password.String, "password should not have been updated")
	})

	t.Run("password too short", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "oldpassword")

		// Execute
		dat := url.Values{}
		dat.Set("old_password", "oldpassword")
		dat.Set("new_password", "a")
		dat.Set("new_password_confirmation", "a")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/account/password", dat)

		res := testutils.HTTPAuthDo(t, db, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismsatch")

		var account database.Account
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&account), "finding account")
		assert.Equal(t, acc.Password.String, account.Password.String, "password should not have been updated")
	})

	t.Run("password confirmation mismatch", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "oldpassword")

		// Execute
		dat := url.Values{}
		dat.Set("old_password", "oldpassword")
		dat.Set("new_password", "newpassword1")
		dat.Set("new_password_confirmation", "newpassword2")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/account/password", dat)

		res := testutils.HTTPAuthDo(t, db, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismsatch")

		var account database.Account
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&account), "finding account")
		assert.Equal(t, acc.Password.String, account.Password.String, "password should not have been updated")
	})
}

func TestUpdateEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "pass1234")
		acc.EmailVerified = true
		testutils.MustExec(t, db.Save(&acc), "updating email_verified")

		// Execute
		dat := url.Values{}
		dat.Set("email", "alice-new@example.com")
		dat.Set("password", "pass1234")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/account/profile", dat)

		res := testutils.HTTPAuthDo(t, db, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusFound, "Status code mismatch")

		var user database.User
		var account database.Account
		testutils.MustExec(t, db.Where("id = ?", u.ID).First(&user), "finding user")
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&account), "finding account")

		assert.Equal(t, account.Email.String, "alice-new@example.com", "email mismatch")
		assert.Equal(t, account.EmailVerified, false, "EmailVerified mismatch")
	})

	t.Run("password mismatch", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		u := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, u, "alice@example.com", "pass1234")
		acc.EmailVerified = true
		testutils.MustExec(t, db.Save(&acc), "updating email_verified")

		// Execute
		dat := url.Values{}
		dat.Set("email", "alice-new@example.com")
		dat.Set("password", "wrongpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/account/profile", dat)

		res := testutils.HTTPAuthDo(t, db, req, u)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "Status code mismsatch")

		var user database.User
		var account database.Account
		testutils.MustExec(t, db.Where("id = ?", u.ID).First(&user), "finding user")
		testutils.MustExec(t, db.Where("user_id = ?", u.ID).First(&account), "finding account")

		assert.Equal(t, account.Email.String, "alice@example.com", "email mismatch")
		assert.Equal(t, account.EmailVerified, true, "EmailVerified mismatch")
	})
}

func TestVerifyEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		user := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, user, "alice@example.com", "pass1234")
		tok := database.Token{
			UserID: user.ID,
			Type:   database.TokenTypeEmailVerification,
			Value:  "someTokenValue",
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		// Execute
		req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/verify-email/%s", "someTokenValue"), "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusFound, "Status code mismatch")

		var account database.Account
		var token database.Token
		var tokenCount int64
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, true, "email_verified mismatch")
		assert.NotEqual(t, token.Value, "", "token value should not have been updated")
		assert.Equal(t, tokenCount, int64(1), "token count mismatch")
		assert.NotEqual(t, token.UsedAt, (*time.Time)(nil), "token should have been used")
	})

	t.Run("used token", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		user := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, user, "alice@example.com", "pass1234")

		usedAt := time.Now().Add(time.Hour * -11).UTC()
		tok := database.Token{
			UserID: user.ID,
			Type:   database.TokenTypeEmailVerification,
			Value:  "someTokenValue",
			UsedAt: &usedAt,
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		// Execute
		req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/verify-email/%s", "someTokenValue"), "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "")

		var account database.Account
		var token database.Token
		var tokenCount int64
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, false, "email_verified mismatch")
		assert.NotEqual(t, token.UsedAt, nil, "token used_at mismatch")
		assert.Equal(t, tokenCount, int64(1), "token count mismatch")
		assert.NotEqual(t, token.UsedAt, (*time.Time)(nil), "token should have been used")
	})

	t.Run("expired token", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		user := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, user, "alice@example.com", "pass1234")

		tok := database.Token{
			UserID: user.ID,
			Type:   database.TokenTypeEmailVerification,
			Value:  "someTokenValue",
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")
		testutils.MustExec(t, db.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-31)), "Failed to prepare token created_at")

		// Execute
		req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/verify-email/%s", "someTokenValue"), "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusGone, "")

		var account database.Account
		var token database.Token
		var tokenCount int64
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, false, "email_verified mismatch")
		assert.Equal(t, tokenCount, int64(1), "token count mismatch")
		assert.Equal(t, token.UsedAt, (*time.Time)(nil), "token should have not been used")
	})

	t.Run("already verified", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		user := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, user, "alice@example.com", "oldpass1234")
		acc.EmailVerified = true
		testutils.MustExec(t, db.Save(&acc), "preparing account")

		tok := database.Token{
			UserID: user.ID,
			Type:   database.TokenTypeEmailVerification,
			Value:  "someTokenValue",
		}
		testutils.MustExec(t, db.Save(&tok), "preparing token")

		// Execute
		req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/verify-email/%s", "someTokenValue"), "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusConflict, "")

		var account database.Account
		var token database.Token
		var tokenCount int64
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, true, "email_verified mismatch")
		assert.Equal(t, tokenCount, int64(1), "token count mismatch")
		assert.Equal(t, token.UsedAt, (*time.Time)(nil), "token should have not been used")
	})
}

func TestCreateVerificationToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)

		// Setup
		emailBackend := testutils.MockEmailbackendImplementation{}
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		a.EmailBackend = &emailBackend
		server := MustNewServer(t, &a)
		defer server.Close()

		user := testutils.SetupUserData(db)
		testutils.SetupAccountData(db, user, "alice@example.com", "pass1234")

		// Execute
		req := testutils.MakeReq(server.URL, "POST", "/verification-token", "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusFound, "status code mismatch")

		var account database.Account
		var token database.Token
		var tokenCount int64
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Where("user_id = ? AND type = ?", user.ID, database.TokenTypeEmailVerification).First(&token), "finding token")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, false, "email_verified should not have been updated")
		assert.NotEqual(t, token.Value, "", "token Value mismatch")
		assert.Equal(t, tokenCount, int64(1), "token count mismatch")
		assert.Equal(t, token.UsedAt, (*time.Time)(nil), "token UsedAt mismatch")
		assert.Equal(t, len(emailBackend.Emails), 1, "email queue count mismatch")
	})

	t.Run("already verified", func(t *testing.T) {
		db := testutils.InitMemoryDB(t)
		// Setup
		a := app.NewTest()
		a.Clock = clock.NewMock()
		a.DB = db
		server := MustNewServer(t, &a)
		defer server.Close()

		user := testutils.SetupUserData(db)
		acc := testutils.SetupAccountData(db, user, "alice@example.com", "pass1234")
		acc.EmailVerified = true
		testutils.MustExec(t, db.Save(&acc), "preparing account")

		// Execute
		req := testutils.MakeReq(server.URL, "POST", "/verification-token", "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusConflict, "Status code mismatch")

		var account database.Account
		var tokenCount int64
		testutils.MustExec(t, db.Where("user_id = ?", user.ID).First(&account), "finding account")
		testutils.MustExec(t, db.Model(&database.Token{}).Count(&tokenCount), "counting token")

		assert.Equal(t, account.EmailVerified, true, "email_verified should not have been updated")
		assert.Equal(t, tokenCount, int64(0), "token count mismatch")
	})
}
