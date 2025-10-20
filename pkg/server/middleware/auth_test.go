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

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestGuestOnly(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(GuestOnly(db, handler))
	defer server.Close()

	t.Run("guest", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusOK, "status code mismatch")
	})

	t.Run("logged in", func(t *testing.T) {
		user := testutils.SetupUserData(db, "user@test.com", "password123")
		req := testutils.MakeReq(server.URL, "GET", "/", "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		assert.Equal(t, res.StatusCode, http.StatusFound, "status code mismatch")
		assert.Equal(t, res.Header.Get("Location"), "/", "location mismatch")
	})

	t.Run("error getting credential", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")
		req.Header.Set("Authorization", "InvalidFormat")
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusOK, "status code mismatch")
	})
}

func TestAuth(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db, "alice@test.com", "pass1234")

	session := database.Session{
		Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	testutils.MustExec(t, db.Save(&session), "preparing session")
	expiredSession := database.Session{
		Key:       "Vvgm3eBXfXGEFWERI7faiRJ3DAzJw+7DdT9J1LEyNfI=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(-time.Hour * 24),
	}
	testutils.MustExec(t, db.Save(&expiredSession), "preparing expired session")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	t.Run("valid session with header", func(t *testing.T) {
		server := httptest.NewServer(Auth(db, handler, nil))
		defer server.Close()

		req := testutils.MakeReq(server.URL, "GET", "/", "")
		req.Header.Set("Authorization", "Bearer "+session.Key)
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusOK, "status code mismatch")
	})

	t.Run("expired session with header", func(t *testing.T) {
		server := httptest.NewServer(Auth(db, handler, nil))
		defer server.Close()

		req := testutils.MakeReq(server.URL, "GET", "/", "")
		req.Header.Set("Authorization", "Bearer "+expiredSession.Key)
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})

	t.Run("invalid session with header", func(t *testing.T) {
		server := httptest.NewServer(Auth(db, handler, nil))
		defer server.Close()

		req := testutils.MakeReq(server.URL, "GET", "/", "")
		req.Header.Set("Authorization", "Bearer someInvalidSessionKey=")
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})

	t.Run("valid session with cookie", func(t *testing.T) {
		server := httptest.NewServer(Auth(db, handler, nil))
		defer server.Close()

		req := testutils.MakeReq(server.URL, "GET", "/", "")
		req.AddCookie(&http.Cookie{
			Name:     "id",
			Value:    session.Key,
			HttpOnly: true,
		})
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusOK, "status code mismatch")
	})

	t.Run("expired session with cookie", func(t *testing.T) {
		server := httptest.NewServer(Auth(db, handler, nil))
		defer server.Close()

		req := testutils.MakeReq(server.URL, "GET", "/", "")
		req.AddCookie(&http.Cookie{
			Name:     "id",
			Value:    expiredSession.Key,
			HttpOnly: true,
		})
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})

	t.Run("no auth", func(t *testing.T) {
		server := httptest.NewServer(Auth(db, handler, nil))
		defer server.Close()

		req := testutils.MakeReq(server.URL, "GET", "/", "")
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})

	t.Run("redirect guests to login", func(t *testing.T) {
		server := httptest.NewServer(Auth(db, handler, &AuthParams{RedirectGuestsToLogin: true}))
		defer server.Close()

		req := testutils.MakeReq(server.URL, "GET", "/settings", "")
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusFound, "status code mismatch")
		assert.Equal(t, res.Header.Get("Location"), "/login?referrer=%2Fsettings", "location mismatch")
	})
}

func TestTokenAuth(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db, "user@test.com", "password123")
	tok := database.Token{
		UserID: user.ID,
		Type:   database.TokenTypeResetPassword,
		Value:  "xpwFnc0MdllFUePDq9DLeQ==",
	}
	testutils.MustExec(t, db.Save(&tok), "preparing token")
	session := database.Session{
		Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	testutils.MustExec(t, db.Save(&session), "preparing session")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(TokenAuth(db, handler, database.TokenTypeResetPassword, nil))
	defer server.Close()

	t.Run("with token", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/?token=xpwFnc0MdllFUePDq9DLeQ==", "")
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusOK, "status code mismatch")
	})

	t.Run("with invalid token", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/?token=someRandomToken==", "")
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})

	t.Run("with session header", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")
		req.Header.Set("Authorization", "Bearer "+session.Key)
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusOK, "status code mismatch")
	})

	t.Run("with invalid session", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")
		req.Header.Set("Authorization", "Bearer someInvalidSessionKey=")
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})

	t.Run("without anything", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")
		res := testutils.HTTPDo(t, req)

		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})
}

func TestWithAccount(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	t.Run("authenticated user", func(t *testing.T) {
		user := testutils.SetupUserData(db, "alice@test.com", "pass1234")

		server := httptest.NewServer(Auth(db, handler, nil))
		defer server.Close()

		req := testutils.MakeReq(server.URL, "GET", "/", "")
		res := testutils.HTTPAuthDo(t, db, req, user)

		assert.Equal(t, res.StatusCode, http.StatusOK, "status code mismatch")
	})
}
