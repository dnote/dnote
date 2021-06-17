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

func assertResponseSessionCookie(t *testing.T, res *http.Response) {
	var sessionCount int
	var session database.Session
	testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
	testutils.MustExec(t, testutils.DB.First(&session), "getting session")

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
		onPremise            bool
		expectedPro          bool
	}{
		{
			email:                "alice@example.com",
			password:             "pass1234",
			passwordConfirmation: "pass1234",
			onPremise:            false,
			expectedPro:          false,
		},
		{
			email:                "bob@example.com",
			password:             "Y9EwmjH@Jq6y5a64MSACUoM4w7SAhzvY",
			passwordConfirmation: "Y9EwmjH@Jq6y5a64MSACUoM4w7SAhzvY",
			onPremise:            false,
			expectedPro:          false,
		},
		{
			email:                "chuck@example.com",
			password:             "e*H@kJi^vXbWEcD9T5^Am!Y@7#Po2@PC",
			passwordConfirmation: "e*H@kJi^vXbWEcD9T5^Am!Y@7#Po2@PC",
			onPremise:            false,
			expectedPro:          false,
		},
		// on premise
		{
			email:                "dan@example.com",
			password:             "e*H@kJi^vXbWEcD9T5^Am!Y@7#Po2@PC",
			passwordConfirmation: "e*H@kJi^vXbWEcD9T5^Am!Y@7#Po2@PC",
			onPremise:            true,
			expectedPro:          true,
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
			dat.Set("password_confirmation", tc.passwordConfirmation)
			req := testutils.MakeFormReq(server.URL, "POST", "/join", dat)

			// Execute
			res := testutils.HTTPDo(t, req)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusFound, "")

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

func TestJoniError(t *testing.T) {
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

	t.Run("password confirmation mismatch", func(t *testing.T) {
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
		dat.Set("password", "pass1234")
		dat.Set("password_confirmation", "1234pass")
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
	dat.Set("password_confirmation", "foobarbaz")
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

func TestLogin(t *testing.T) {
	testutils.RunForWebAndAPI(t, "success", func(t *testing.T, target testutils.EndpointType) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "pass1234")
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
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.NotEqual(t, user.LastLoginAt, nil, "LastLoginAt mismatch")

		if target == testutils.EndpointWeb {
			assertResponseSessionCookie(t, res)
		} else {
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

			assertResponseSessionCookie(t, res)
		}
	})

	testutils.RunForWebAndAPI(t, "wrong password", func(t *testing.T, target testutils.EndpointType) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "pass1234")
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
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.Equal(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})

	testutils.RunForWebAndAPI(t, "wrong email", func(t *testing.T, target testutils.EndpointType) {
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
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.DeepEqual(t, user.LastLoginAt, (*time.Time)(nil), "LastLoginAt mismatch")

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})

	testutils.RunForWebAndAPI(t, "nonexistent email", func(t *testing.T, target testutils.EndpointType) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{
			Clock: clock.NewMock(),
			Config: config.Config{
				PageTemplateDir: "../views",
			},
		})
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

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})
}

func TestLogout(t *testing.T) {
	setupLogoutTest := func(t *testing.T) (*httptest.Server, *database.Session, *database.Session) {
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

	testutils.RunForWebAndAPI(t, "authenticated", func(t *testing.T, target testutils.EndpointType) {
		defer testutils.ClearData(testutils.DB)

		server, session1, _ := setupLogoutTest(t)
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

		var sessionCount int
		var s2 database.Session
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		testutils.MustExec(t, testutils.DB.Where("key = ?", "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=").First(&s2), "getting s2")

		assert.Equal(t, sessionCount, 1, "sessionCount mismatch")

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
		defer testutils.ClearData(testutils.DB)

		server, _, _ := setupLogoutTest(t)
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

		var sessionCount int
		var postSession1, postSession2 database.Session
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		testutils.MustExec(t, testutils.DB.Where("key = ?", "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=").First(&postSession1), "getting postSession1")
		testutils.MustExec(t, testutils.DB.Where("key = ?", "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=").First(&postSession2), "getting postSession2")

		// two existing sessions should remain
		assert.Equal(t, sessionCount, 2, "sessionCount mismatch")

		c := testutils.GetCookieByName(res.Cookies(), "id")
		assert.Equal(t, c, (*http.Cookie)(nil), "id cookie should have not been set")
	})
}

func TestResetPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
		a := testutils.SetupAccountData(u, "alice@example.com", "oldpassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
		}
		testutils.MustExec(t, testutils.DB.Save(&tok), "preparing token")
		otherTok := database.Token{
			UserID: u.ID,
			Value:  "somerandomvalue",
			Type:   database.TokenTypeEmailVerification,
		}
		testutils.MustExec(t, testutils.DB.Save(&otherTok), "preparing another token")

		s1 := database.Session{
			Key:       "some-session-key-1",
			UserID:    u.ID,
			ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
		}
		testutils.MustExec(t, testutils.DB.Save(&s1), "preparing user session 1")

		s2 := &database.Session{
			Key:       "some-session-key-2",
			UserID:    u.ID,
			ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
		}
		testutils.MustExec(t, testutils.DB.Save(&s2), "preparing user session 2")

		anotherUser := testutils.SetupUserData()
		testutils.MustExec(t, testutils.DB.Save(&database.Session{
			Key:       "some-session-key-3",
			UserID:    anotherUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
		}), "preparing anotherUser session 1")

		// Execute
		dat := url.Values{}
		dat.Set("token", "MivFxYiSMMA4An9dP24DNQ==")
		dat.Set("password", "newpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismatch")

		var resetToken, verificationToken database.Token
		var account database.Account
		testutils.MustExec(t, testutils.DB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "finding reset token")
		testutils.MustExec(t, testutils.DB.Where("value = ?", "somerandomvalue").First(&verificationToken), "finding reset token")
		testutils.MustExec(t, testutils.DB.Where("id = ?", a.ID).First(&account), "finding account")

		assert.NotEqual(t, resetToken.UsedAt, nil, "reset_token UsedAt mismatch")
		passwordErr := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte("newpassword"))
		assert.Equal(t, passwordErr, nil, "Password mismatch")
		assert.Equal(t, verificationToken.UsedAt, (*time.Time)(nil), "verificationToken UsedAt mismatch")

		var s1Count, s2Count int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Where("id = ?", s1.ID).Count(&s1Count), "counting s1")
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Where("id = ?", s2.ID).Count(&s2Count), "counting s2")

		assert.Equal(t, s1Count, 0, "s1 should have been deleted")
		assert.Equal(t, s2Count, 0, "s2 should have been deleted")

		var userSessionCount, anotherUserSessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Where("user_id = ?", u.ID).Count(&userSessionCount), "counting user session")
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Where("user_id = ?", anotherUser.ID).Count(&anotherUserSessionCount), "counting anotherUser session")

		assert.Equal(t, userSessionCount, 1, "should have created a new user session")
		assert.Equal(t, anotherUserSessionCount, 1, "anotherUser session count mismatch")
	})

	t.Run("nonexistent token", func(t *testing.T) {
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
		a := testutils.SetupAccountData(u, "alice@example.com", "somepassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
		}
		testutils.MustExec(t, testutils.DB.Save(&tok), "preparing token")

		dat := url.Values{}
		dat.Set("token", "-ApMnyvpg59uOU5b-Kf5uQ==")
		dat.Set("password", "oldpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, testutils.DB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "finding reset token")
		testutils.MustExec(t, testutils.DB.Where("id = ?", a.ID).First(&account), "finding account")

		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})

	t.Run("expired token", func(t *testing.T) {
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
		a := testutils.SetupAccountData(u, "alice@example.com", "somepassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
		}
		testutils.MustExec(t, testutils.DB.Save(&tok), "preparing token")
		testutils.MustExec(t, testutils.DB.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := url.Values{}
		dat.Set("token", "MivFxYiSMMA4An9dP24DNQ==")
		dat.Set("password", "oldpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusGone, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, testutils.DB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		testutils.MustExec(t, testutils.DB.Where("id = ?", a.ID).First(&account), "failed to find account")
		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})

	t.Run("used token", func(t *testing.T) {
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
		a := testutils.SetupAccountData(u, "alice@example.com", "somepassword")

		usedAt := time.Now().Add(time.Hour * -11).UTC()
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeResetPassword,
			UsedAt: &usedAt,
		}
		testutils.MustExec(t, testutils.DB.Save(&tok), "preparing token")
		testutils.MustExec(t, testutils.DB.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := url.Values{}
		dat.Set("token", "MivFxYiSMMA4An9dP24DNQ==")
		dat.Set("password", "oldpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, testutils.DB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		testutils.MustExec(t, testutils.DB.Where("id = ?", a.ID).First(&account), "failed to find account")
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
		a := testutils.SetupAccountData(u, "alice@example.com", "somepassword")
		tok := database.Token{
			UserID: u.ID,
			Value:  "MivFxYiSMMA4An9dP24DNQ==",
			Type:   database.TokenTypeEmailVerification,
		}
		testutils.MustExec(t, testutils.DB.Save(&tok), "Failed to prepare reset_token")
		testutils.MustExec(t, testutils.DB.Model(&tok).Update("created_at", time.Now().Add(time.Minute*-11)), "Failed to prepare reset_token created_at")

		dat := url.Values{}
		dat.Set("token", "MivFxYiSMMA4An9dP24DNQ==")
		dat.Set("password", "oldpassword")
		req := testutils.MakeFormReq(server.URL, "PATCH", "/password-reset", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status code mismatch")

		var resetToken database.Token
		var account database.Account
		testutils.MustExec(t, testutils.DB.Where("value = ?", "MivFxYiSMMA4An9dP24DNQ==").First(&resetToken), "failed to find reset_token")
		testutils.MustExec(t, testutils.DB.Where("id = ?", a.ID).First(&account), "failed to find account")

		assert.Equal(t, a.Password, account.Password, "password should not have been updated")
		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "used_at should be nil")
	})
}
