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
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func assertSessionCookie(t *testing.T, c *http.Cookie, session database.Session) {
	assert.Equal(t, c.Value, session.Key, "session key mismatch")
	assert.Equal(t, c.Path, "/", "session path mismatch")
	assert.Equal(t, c.HttpOnly, true, "session HTTPOnly mismatch")
	assert.Equal(t, c.Expires.Unix(), session.ExpiresAt.Unix(), "session Expires mismatch")
}

func assertResponseSessionCookie(t *testing.T, res *http.Response) {
	var sessionCount int
	var session database.Session
	testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
	testutils.MustExec(t, testutils.DB.First(&session), "getting session")

	c := testutils.GetCookieByName(res.Cookies(), "id")
	assertSessionCookie(t, c, session)
}

func assertSessionResp(t *testing.T, res *http.Response) {
	// after register, should sign in user
	var got SessionResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	var sessionCount int
	var session database.Session
	testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
	testutils.MustExec(t, testutils.DB.First(&session), "getting session")

	assert.Equal(t, sessionCount, 1, "sessionCount mismatch")
	assert.Equal(t, got.Key, session.Key, "session Key mismatch")
	assert.Equal(t, got.ExpiresAt, session.ExpiresAt.Unix(), "session ExpiresAt mismatch")

	c := testutils.GetCookieByName(res.Cookies(), "id")
	assertSessionCookie(t, c, session)
}

func TestJoin(t *testing.T) {
	testCases := []struct {
		email       string
		password    string
		onPremise   bool
		expectedPro bool
	}{
		{
			email:       "alice@example.com",
			password:    "pass1234",
			onPremise:   false,
			expectedPro: false,
		},
		{
			email:       "bob@example.com",
			password:    "Y9EwmjH@Jq6y5a64MSACUoM4w7SAhzvY",
			onPremise:   false,
			expectedPro: false,
		},
		{
			email:       "chuck@example.com",
			password:    "e*H@kJi^vXbWEcD9T5^Am!Y@7#Po2@PC",
			onPremise:   false,
			expectedPro: false,
		},
		// on premise
		{
			email:       "dan@example.com",
			password:    "e*H@kJi^vXbWEcD9T5^Am!Y@7#Po2@PC",
			onPremise:   true,
			expectedPro: true,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("register %s %s", tc.email, tc.password), func(t *testing.T) {
			defer testutils.ClearData(testutils.DB)

			// Setup
			emailBackend := testutils.MockEmailbackendImplementation{}
			server := MustNewServer(t, &app.App{
				Clock:        clock.NewMock(),
				EmailBackend: &emailBackend,
				Config: config.Config{
					OnPremise:       tc.onPremise,
					PageTemplateDir: "../views",
				},
			})
			defer server.Close()

			dat := url.Values{}
			dat.Set("email", tc.email)
			dat.Set("password", tc.password)
			req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

			// Execute
			res := testutils.HTTPDo(t, req)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusCreated, "")

			var account database.Account
			testutils.MustExec(t, testutils.DB.Where("email = ?", tc.email).First(&account), "finding account")
			assert.Equal(t, account.Email.String, tc.email, "Email mismatch")
			assert.NotEqual(t, account.UserID, 0, "UserID mismatch")
			passwordErr := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte(tc.password))
			assert.Equal(t, passwordErr, nil, "Password mismatch")

			var user database.User
			testutils.MustExec(t, testutils.DB.Where("id = ?", account.UserID).First(&user), "finding user")
			assert.Equal(t, user.Cloud, tc.expectedPro, "Cloud mismatch")
			assert.Equal(t, user.MaxUSN, 0, "MaxUSN mismatch")

			// welcome email
			assert.Equalf(t, len(emailBackend.Emails), 1, "email queue count mismatch")
			assert.DeepEqual(t, emailBackend.Emails[0].To, []string{tc.email}, "email to mismatch")

			// after register, should sign in user
			assertResponseSessionCookie(t, res)
		})
	}
}

func TestJoinMissingParams(t *testing.T) {
	t.Run("missing email", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
		defer server.Close()

		dat := url.Values{}
		dat.Set("password", "SLMZFM5RmSjA5vfXnG5lPOnrpZSbtmV76cnAcrlr2yU")
		req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status mismatch")

		var accountCount, userCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Account{}).Count(&accountCount), "counting account")
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).Count(&userCount), "counting user")

		assert.Equal(t, accountCount, 0, "accountCount mismatch")
		assert.Equal(t, userCount, 0, "userCount mismatch")
	})

	t.Run("missing password", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
		defer server.Close()

		dat := url.Values{}
		dat.Set("email", "alice@example.com")
		req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status mismatch")

		var accountCount, userCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Account{}).Count(&accountCount), "counting account")
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).Count(&userCount), "counting user")

		assert.Equal(t, accountCount, 0, "accountCount mismatch")
		assert.Equal(t, userCount, 0, "userCount mismatch")
	})
}

func TestJoinDuplicateEmail(t *testing.T) {
	defer testutils.ClearData(testutils.DB)

	// Setup
	server := MustNewServer(t, &app.App{
		Clock: clock.NewMock(),
		Config: config.Config{
			PageTemplateDir: "../views",
		},
	})
	defer server.Close()

	u := testutils.SetupUserData()
	testutils.SetupAccountData(u, "alice@example.com", "somepassword")

	dat := url.Values{}
	dat.Set("email", "alice@example.com")
	dat.Set("password", "foobarbaz")
	req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

	// Execute
	res := testutils.HTTPDo(t, req)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusBadRequest, "status code mismatch")

	var accountCount, userCount, verificationTokenCount int
	testutils.MustExec(t, testutils.DB.Model(&database.Account{}).Count(&accountCount), "counting account")
	testutils.MustExec(t, testutils.DB.Model(&database.User{}).Count(&userCount), "counting user")
	testutils.MustExec(t, testutils.DB.Model(&database.Token{}).Count(&verificationTokenCount), "counting verification token")

	var user database.User
	testutils.MustExec(t, testutils.DB.Where("id = ?", u.ID).First(&user), "finding user")

	assert.Equal(t, accountCount, 1, "account count mismatch")
	assert.Equal(t, userCount, 1, "user count mismatch")
	assert.Equal(t, verificationTokenCount, 0, "verification_token should not have been created")
	assert.Equal(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")
}

func TestJoinDisabled(t *testing.T) {
	defer testutils.ClearData(testutils.DB)

	// Setup
	server := MustNewServer(t, &app.App{
		Clock: clock.NewMock(),
		Config: config.Config{
			PageTemplateDir:     "../views",
			DisableRegistration: true,
		},
	})
	defer server.Close()

	dat := url.Values{}
	dat.Set("email", "alice@example.com")
	dat.Set("password", "foobarbaz")
	req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

	// Execute
	res := testutils.HTTPDo(t, req)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusNotFound, "status code mismatch")

	var accountCount, userCount int
	testutils.MustExec(t, testutils.DB.Model(&database.Account{}).Count(&accountCount), "counting account")
	testutils.MustExec(t, testutils.DB.Model(&database.User{}).Count(&userCount), "counting user")

	assert.Equal(t, accountCount, 0, "account count mismatch")
	assert.Equal(t, userCount, 0, "user count mismatch")
}

func setupLoginTest(t *testing.T) *httptest.Server {
	// Setup
	server := MustNewServer(t, &app.App{
		Clock: clock.NewMock(),
		Config: config.Config{
			PageTemplateDir: "../views",
		},
	})

	u := testutils.SetupUserData()
	testutils.SetupAccountData(u, "alice@example.com", "pass1234")

	return server
}

func TestV3Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		server := setupLoginTest(t)
		defer server.Close()

		// Execute
		dat := `{"email": "alice@example.com", "password": "pass1234"}`
		req := testutils.MakeReq(server.URL, "POST", "/api/v3/signin", dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var user database.User
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.NotEqual(t, user.LastLoginAt, nil, "LastLoginAt mismatch")

		assertSessionResp(t, res)
	})

	t.Run("wrong password", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := setupLoginTest(t)
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "pass1234")

		dat := `{"email": "alice@example.com", "password": "wrongpassword1234"}`
		req := testutils.MakeReq(server.URL, "POST", "/api/v3/signin", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var user database.User
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.Equal(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})

	t.Run("wrong email", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "pass1234")

		dat := `{"email": "bob@example.com", "password": "foobarbaz"}`
		req := testutils.MakeReq(server.URL, "POST", "/api/v3/signin", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var user database.User
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.DeepEqual(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})

	t.Run("nonexistent email", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
		defer server.Close()

		dat := `{"email": "nonexistent@example.com", "password": "pass1234"}`
		req := testutils.MakeReq(server.URL, "POST", "/api/v3/signin", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})
}

func TestWebLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		server := setupLoginTest(t)
		defer server.Close()

		// Execute
		dat := url.Values{}
		dat.Set("email", "alice@example.com")
		dat.Set("password", "pass1234")
		req := testutils.MakeFormReq(server.URL, "POST", "/login", dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var user database.User
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.NotEqual(t, user.LastLoginAt, nil, "LastLoginAt mismatch")

		assertResponseSessionCookie(t, res)
	})

	t.Run("wrong password", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := setupLoginTest(t)
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "pass1234")

		dat := url.Values{}
		dat.Set("email", "alice@example.com")
		dat.Set("password", "wrongpassword1234")
		req := testutils.MakeFormReq(server.URL, "POST", "/login", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var user database.User
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.Equal(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})

	t.Run("wrong email", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "pass1234")

		dat := url.Values{}
		dat.Set("email", "bob@example.com")
		dat.Set("password", "foobarbaz")
		req := testutils.MakeFormReq(server.URL, "POST", "/login", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var user database.User
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.DeepEqual(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})

	t.Run("nonexistent email", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
		defer server.Close()

		dat := url.Values{}
		dat.Set("email", "nonexistent@example.com")
		dat.Set("password", "pass1234")
		req := testutils.MakeFormReq(server.URL, "POST", "/login", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})
}

func setupLogoutTest(t *testing.T) (*httptest.Server, *database.Session, *database.Session) {
	// Setup
	server := MustNewServer(t, &app.App{
		Clock: clock.NewMock(),
		Config: config.Config{
			PageTemplateDir: "../views",
		},
	})

	aliceUser := testutils.SetupUserData()
	testutils.SetupAccountData(aliceUser, "alice@example.com", "pass1234")
	anotherUser := testutils.SetupUserData()

	session1ExpiresAt := time.Now().Add(time.Hour * 24)
	session1 := database.Session{
		Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
		UserID:    aliceUser.ID,
		ExpiresAt: session1ExpiresAt,
	}
	testutils.MustExec(t, testutils.DB.Save(&session1), "preparing session1")
	session2 := database.Session{
		Key:       "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=",
		UserID:    anotherUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	testutils.MustExec(t, testutils.DB.Save(&session2), "preparing session2")

	return server, &session1, &session2
}

func assertLogoutAuthenticated(t *testing.T, res *http.Response) {
	assert.StatusCodeEquals(t, res, http.StatusNoContent, "Status mismatch")

	var sessionCount int
	var s2 database.Session
	testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
	testutils.MustExec(t, testutils.DB.Where("key = ?", "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=").First(&s2), "getting s2")

	assert.Equal(t, sessionCount, 1, "sessionCount mismatch")
}

func assertLogoutUnauthenticated(t *testing.T, res *http.Response) {
	assert.StatusCodeEquals(t, res, http.StatusNoContent, "Status mismatch")

	var sessionCount int
	var postSession1, postSession2 database.Session
	testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
	testutils.MustExec(t, testutils.DB.Where("key = ?", "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=").First(&postSession1), "getting postSession1")
	testutils.MustExec(t, testutils.DB.Where("key = ?", "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=").First(&postSession2), "getting postSession2")

	// two existing sessions should remain
	assert.Equal(t, sessionCount, 2, "sessionCount mismatch")

	c := testutils.GetCookieByName(res.Cookies(), "id")
	assert.Equal(t, c, (*http.Cookie)(nil), "id cookie should have not been set")
}

func TestV3Logout(t *testing.T) {
	t.Run("authenticated", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		server, session1, _ := setupLogoutTest(t)
		defer server.Close()

		// Execute
		req := testutils.MakeReq(server.URL, "POST", "/api/v3/signout", "")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session1.Key))
		res := testutils.HTTPDo(t, req)

		// Test
		assertLogoutAuthenticated(t, res)
	})

	t.Run("unauthenticated", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		server, _, _ := setupLogoutTest(t)
		defer server.Close()

		// Execute
		req := testutils.MakeReq(server.URL, "POST", "/api/v3/signout", "")
		res := testutils.HTTPDo(t, req)

		// Test
		assertLogoutUnauthenticated(t, res)
	})
}

func TestWebLogout(t *testing.T) {
	t.Run("authenticated", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		server, session1, _ := setupLogoutTest(t)
		defer server.Close()

		// Execute
		dat := url.Values{}
		req := testutils.MakeFormReq(server.URL, "POST", "/logout", dat)
		req.AddCookie(&http.Cookie{Name: "id", Value: "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=", Expires: session1.ExpiresAt, Path: "/", HttpOnly: true})

		res := testutils.HTTPDo(t, req)

		// Test
		assertLogoutAuthenticated(t, res)

		c := testutils.GetCookieByName(res.Cookies(), "id")
		assert.Equal(t, c.Value, "", "session key mismatch")
		assert.Equal(t, c.Path, "/", "session path mismatch")
		assert.Equal(t, c.HttpOnly, true, "session HTTPOnly mismatch")
		if c.Expires.After(time.Now()) {
			t.Error("session cookie is not expired")
		}
	})

	t.Run("unauthenticated", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		server, _, _ := setupLogoutTest(t)
		defer server.Close()

		// Execute
		dat := url.Values{}
		req := testutils.MakeFormReq(server.URL, "POST", "/logout", dat)
		res := testutils.HTTPDo(t, req)

		// Test
		assertLogoutUnauthenticated(t, res)
	})
}
