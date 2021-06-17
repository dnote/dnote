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
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

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
	assert.Equal(t, c.Value, session.Key, "session key mismatch")
	assert.Equal(t, c.Path, "/", "session path mismatch")
	assert.Equal(t, c.HttpOnly, true, "session HTTPOnly mismatch")
	assert.Equal(t, c.Expires.Unix(), session.ExpiresAt.Unix(), "session Expires mismatch")
}

func TestSignIn(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "pass1234")

		dat := `{"email": "alice@example.com", "password": "pass1234"}`
		req := testutils.MakeReq(server.URL, "POST", "/v3/signin", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var user database.User
		testutils.MustExec(t, testutils.DB.Model(&database.User{}).First(&user), "finding user")
		assert.NotEqual(t, user.LastLoginAt, nil, "LastLoginAt mismatch")

		// after register, should sign in user
		assertSessionResp(t, res)
	})

	t.Run("wrong password", func(t *testing.T) {
		defer testutils.ClearData(testutils.DB)

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "pass1234")

		dat := `{"email": "alice@example.com", "password": "wrongpassword1234"}`
		req := testutils.MakeReq(server.URL, "POST", "/v3/signin", dat)

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
		})
		defer server.Close()

		u := testutils.SetupUserData()
		testutils.SetupAccountData(u, "alice@example.com", "pass1234")

		dat := `{"email": "bob@example.com", "password": "pass1234"}`
		req := testutils.MakeReq(server.URL, "POST", "/v3/signin", dat)

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
		})
		defer server.Close()

		dat := `{"email": "nonexistent@example.com", "password": "pass1234"}`
		req := testutils.MakeReq(server.URL, "POST", "/v3/signin", dat)

		// Execute
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")

		var sessionCount int
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		assert.Equal(t, sessionCount, 0, "sessionCount mismatch")
	})
}

func TestSignout(t *testing.T) {
	t.Run("authenticated", func(t *testing.T) {

		defer testutils.ClearData(testutils.DB)

		aliceUser := testutils.SetupUserData()
		testutils.SetupAccountData(aliceUser, "alice@example.com", "pass1234")
		anotherUser := testutils.SetupUserData()

		session1 := database.Session{
			Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
			UserID:    aliceUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}
		testutils.MustExec(t, testutils.DB.Save(&session1), "preparing session1")
		session2 := database.Session{
			Key:       "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=",
			UserID:    anotherUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}
		testutils.MustExec(t, testutils.DB.Save(&session2), "preparing session2")

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		// Execute
		req := testutils.MakeReq(server.URL, "POST", "/v3/signout", "")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU="))
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNoContent, "Status mismatch")

		var sessionCount int
		var s2 database.Session
		testutils.MustExec(t, testutils.DB.Model(&database.Session{}).Count(&sessionCount), "counting session")
		testutils.MustExec(t, testutils.DB.Where("key = ?", "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=").First(&s2), "getting s2")

		assert.Equal(t, sessionCount, 1, "sessionCount mismatch")

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

		aliceUser := testutils.SetupUserData()
		testutils.SetupAccountData(aliceUser, "alice@example.com", "pass1234")
		anotherUser := testutils.SetupUserData()

		session1 := database.Session{
			Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
			UserID:    aliceUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}
		testutils.MustExec(t, testutils.DB.Save(&session1), "preparing session1")
		session2 := database.Session{
			Key:       "MDCpbvCRg7W2sH6S870wqLqZDZTObYeVd0PzOekfo/A=",
			UserID:    anotherUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}
		testutils.MustExec(t, testutils.DB.Save(&session2), "preparing session2")

		// Setup
		server := MustNewServer(t, &app.App{

			Clock: clock.NewMock(),
		})
		defer server.Close()

		// Execute
		req := testutils.MakeReq(server.URL, "POST", "/v3/signout", "")
		res := testutils.HTTPDo(t, req)

		// Test
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
	})
}
