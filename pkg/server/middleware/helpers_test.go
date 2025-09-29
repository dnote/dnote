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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestGetSessionKeyFromCookie(t *testing.T) {
	testCases := []struct {
		cookie   *http.Cookie
		expected string
	}{
		{
			cookie: &http.Cookie{
				Name:     "id",
				Value:    "foo",
				HttpOnly: true,
			},
			expected: "foo",
		},
		{
			cookie:   nil,
			expected: "",
		},
		{
			cookie: &http.Cookie{
				Name:     "foo",
				Value:    "bar",
				HttpOnly: true,
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		// set up
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(errors.Wrap(err, "constructing request"))
		}

		if tc.cookie != nil {
			r.AddCookie(tc.cookie)
		}

		// execute
		got, err := getSessionKeyFromCookie(r)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}

func TestGetSessionKeyFromAuth(t *testing.T) {
	testCases := []struct {
		authHeaderStr string
		expected      string
	}{
		{
			authHeaderStr: "Bearer foo",
			expected:      "foo",
		},
	}

	for _, tc := range testCases {
		// set up
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(errors.Wrap(err, "constructing request"))
		}

		r.Header.Set("Authorization", tc.authHeaderStr)

		// execute
		got, err := getSessionKeyFromAuth(r)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}

func mustMakeRequest(t *testing.T) *http.Request {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(errors.Wrap(err, "constructing request"))
	}

	return r
}

func TestGetCredential(t *testing.T) {
	r1 := mustMakeRequest(t)
	r2 := mustMakeRequest(t)
	r2.Header.Set("Authorization", "Bearer foo")
	r3 := mustMakeRequest(t)
	r3.Header.Set("Authorization", "Bearer bar")

	r4 := mustMakeRequest(t)
	c4 := http.Cookie{
		Name:     "id",
		Value:    "foo",
		HttpOnly: true,
	}
	r4.AddCookie(&c4)

	r5 := mustMakeRequest(t)
	c5 := http.Cookie{
		Name:     "id",
		Value:    "foo",
		HttpOnly: true,
	}
	r5.AddCookie(&c5)
	r5.Header.Set("Authorization", "Bearer foo")

	testCases := []struct {
		request  *http.Request
		expected string
	}{
		{
			request:  r1,
			expected: "",
		},
		{
			request:  r2,
			expected: "foo",
		},
		{
			request:  r3,
			expected: "bar",
		},
		{
			request:  r4,
			expected: "foo",
		},
		{
			request:  r5,
			expected: "foo",
		},
	}

	for _, tc := range testCases {
		// execute
		got, err := GetCredential(tc.request)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}

func TestAuthMiddleware(t *testing.T) {
	defer testutils.ClearData(testutils.DB)

	user := testutils.SetupUserData()
	testutils.SetupAccountData(user, "alice@test.com", "pass1234")

	session := database.Session{
		Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	testutils.MustExec(t, testutils.DB.Save(&session), "preparing session")
	session2 := database.Session{
		Key:       "Vvgm3eBXfXGEFWERI7faiRJ3DAzJw+7DdT9J1LEyNfI=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(-time.Hour * 24),
	}
	testutils.MustExec(t, testutils.DB.Save(&session2), "preparing session")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	a := &app.App{DB: testutils.DB}
	server := httptest.NewServer(Auth(a, handler, nil))
	defer server.Close()

	t.Run("with header", func(t *testing.T) {
		testCases := []struct {
			header         string
			expectedStatus int
		}{
			{
				header:         fmt.Sprintf("Bearer %s", session.Key),
				expectedStatus: http.StatusOK,
			},
			{
				header:         fmt.Sprintf("Bearer %s", session2.Key),
				expectedStatus: http.StatusUnauthorized,
			},
			{
				header:         fmt.Sprintf("Bearer someInvalidSessionKey="),
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.header, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", "/", "")
				req.Header.Set("Authorization", tc.header)

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})

	t.Run("with cookie", func(t *testing.T) {
		testCases := []struct {
			cookie         *http.Cookie
			expectedStatus int
		}{
			{
				cookie: &http.Cookie{
					Name:     "id",
					Value:    session.Key,
					HttpOnly: true,
				},
				expectedStatus: http.StatusOK,
			},
			{
				cookie: &http.Cookie{
					Name:     "id",
					Value:    session2.Key,
					HttpOnly: true,
				},
				expectedStatus: http.StatusUnauthorized,
			},
			{
				cookie: &http.Cookie{
					Name:     "id",
					Value:    "someInvalidSessionKey=",
					HttpOnly: true,
				},
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.cookie.Value, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", "/", "")
				req.AddCookie(tc.cookie)

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})

	t.Run("without anything", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")

		// execute
		res := testutils.HTTPDo(t, req)

		// test
		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})
}

func TestAuthMiddleware_ProOnly(t *testing.T) {
	defer testutils.ClearData(testutils.DB)

	user := testutils.SetupUserData()
	testutils.MustExec(t, testutils.DB.Model(&user).Update("cloud", false), "preparing session")
	session := database.Session{
		Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	testutils.MustExec(t, testutils.DB.Save(&session), "preparing session")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	a := &app.App{DB: testutils.DB}
	server := httptest.NewServer(Auth(a, handler, &AuthParams{
		ProOnly: true,
	}))

	defer server.Close()

	t.Run("with header", func(t *testing.T) {
		testCases := []struct {
			header         string
			expectedStatus int
		}{
			{
				header:         fmt.Sprintf("Bearer %s", session.Key),
				expectedStatus: http.StatusForbidden,
			},
			{
				header:         fmt.Sprintf("Bearer someInvalidSessionKey="),
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.header, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", "/", "")
				req.Header.Set("Authorization", tc.header)

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})

	t.Run("with cookie", func(t *testing.T) {
		testCases := []struct {
			cookie         *http.Cookie
			expectedStatus int
		}{
			{
				cookie: &http.Cookie{
					Name:     "id",
					Value:    session.Key,
					HttpOnly: true,
				},
				expectedStatus: http.StatusForbidden,
			},
			{
				cookie: &http.Cookie{
					Name:     "id",
					Value:    "someInvalidSessionKey=",
					HttpOnly: true,
				},
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.cookie.Value, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", "/", "")
				req.AddCookie(tc.cookie)

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})
}

func TestAuthMiddleware_RedirectGuestsToLogin(t *testing.T) {
	defer testutils.ClearData(testutils.DB)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	a := &app.App{DB: testutils.DB}
	server := httptest.NewServer(Auth(a, handler, &AuthParams{
		RedirectGuestsToLogin: true,
	}))

	defer server.Close()

	t.Run("guest", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")

		// execute
		res := testutils.HTTPDo(t, req)

		// test
		assert.Equal(t, res.StatusCode, http.StatusFound, "status code mismatch")
		assert.Equal(t, res.Header.Get("Location"), "/login?referrer=%2F", "location header mismatch")
	})

	t.Run("logged in user", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")

		user := testutils.SetupUserData()
		testutils.SetupAccountData(user, "alice@test.com", "pass1234")

		testutils.MustExec(t, testutils.DB.Model(&user).Update("cloud", false), "preparing session")
		session := database.Session{
			Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		}
		testutils.MustExec(t, testutils.DB.Save(&session), "preparing session")

		// execute
		res := testutils.HTTPAuthDo(t, req, user)
		req.Header.Set("Authorization", session.Key)

		// test
		assert.Equal(t, res.StatusCode, http.StatusOK, "status code mismatch")
		assert.Equal(t, res.Header.Get("Location"), "", "location header mismatch")
	})

}

func TestTokenAuthMiddleWare(t *testing.T) {
	defer testutils.ClearData(testutils.DB)

	user := testutils.SetupUserData()
	tok := database.Token{
		UserID: user.ID,
		Type:   database.TokenTypeEmailPreference,
		Value:  "xpwFnc0MdllFUePDq9DLeQ==",
	}
	testutils.MustExec(t, testutils.DB.Save(&tok), "preparing token")
	session := database.Session{
		Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	testutils.MustExec(t, testutils.DB.Save(&session), "preparing session")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	a := &app.App{DB: testutils.DB}
	server := httptest.NewServer(TokenAuth(a, handler, database.TokenTypeEmailPreference, nil))
	defer server.Close()

	t.Run("with token", func(t *testing.T) {
		testCases := []struct {
			token          string
			expectedStatus int
		}{
			{
				token:          "xpwFnc0MdllFUePDq9DLeQ==",
				expectedStatus: http.StatusOK,
			},
			{
				token:          "someRandomToken==",
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.token, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/?token=%s", tc.token), "")

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})

	t.Run("with session header", func(t *testing.T) {
		testCases := []struct {
			header         string
			expectedStatus int
		}{
			{
				header:         fmt.Sprintf("Bearer %s", session.Key),
				expectedStatus: http.StatusOK,
			},
			{
				header:         fmt.Sprintf("Bearer someInvalidSessionKey="),
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.header, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", "/", "")
				req.Header.Set("Authorization", tc.header)

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})

	t.Run("with session cookie", func(t *testing.T) {
		testCases := []struct {
			cookie         *http.Cookie
			expectedStatus int
		}{
			{
				cookie: &http.Cookie{
					Name:     "id",
					Value:    session.Key,
					HttpOnly: true,
				},
				expectedStatus: http.StatusOK,
			},
			{
				cookie: &http.Cookie{
					Name:     "id",
					Value:    "someInvalidSessionKey=",
					HttpOnly: true,
				},
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.cookie.Value, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", "/", "")
				req.AddCookie(tc.cookie)

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})

	t.Run("without anything", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")

		// execute
		res := testutils.HTTPDo(t, req)

		// test
		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})
}

func TestTokenAuthMiddleWare_ProOnly(t *testing.T) {
	defer testutils.ClearData(testutils.DB)

	user := testutils.SetupUserData()
	testutils.MustExec(t, testutils.DB.Model(&user).Update("cloud", false), "preparing session")
	tok := database.Token{
		UserID: user.ID,
		Type:   database.TokenTypeEmailPreference,
		Value:  "xpwFnc0MdllFUePDq9DLeQ==",
	}
	testutils.MustExec(t, testutils.DB.Save(&tok), "preparing token")
	session := database.Session{
		Key:       "A9xgggqzTHETy++GDi1NpDNe0iyqosPm9bitdeNGkJU=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	testutils.MustExec(t, testutils.DB.Save(&session), "preparing session")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	a := &app.App{DB: testutils.DB}
	server := httptest.NewServer(TokenAuth(a, handler, database.TokenTypeEmailPreference, &AuthParams{
		ProOnly: true,
	}))

	defer server.Close()

	t.Run("with token", func(t *testing.T) {
		testCases := []struct {
			token          string
			expectedStatus int
		}{
			{
				token:          "xpwFnc0MdllFUePDq9DLeQ==",
				expectedStatus: http.StatusForbidden,
			},
			{
				token:          "someRandomToken==",
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.token, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/?token=%s", tc.token), "")

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})

	t.Run("with session header", func(t *testing.T) {
		testCases := []struct {
			header         string
			expectedStatus int
		}{
			{
				header:         fmt.Sprintf("Bearer %s", session.Key),
				expectedStatus: http.StatusForbidden,
			},
			{
				header:         fmt.Sprintf("Bearer someInvalidSessionKey="),
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.header, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", "/", "")
				req.Header.Set("Authorization", tc.header)

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})

	t.Run("with session cookie", func(t *testing.T) {
		testCases := []struct {
			cookie         *http.Cookie
			expectedStatus int
		}{
			{
				cookie: &http.Cookie{
					Name:     "id",
					Value:    session.Key,
					HttpOnly: true,
				},
				expectedStatus: http.StatusForbidden,
			},
			{
				cookie: &http.Cookie{
					Name:     "id",
					Value:    "someInvalidSessionKey=",
					HttpOnly: true,
				},
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.cookie.Value, func(t *testing.T) {
				req := testutils.MakeReq(server.URL, "GET", "/", "")
				req.AddCookie(tc.cookie)

				// execute
				res := testutils.HTTPDo(t, req)

				// test
				assert.Equal(t, res.StatusCode, tc.expectedStatus, "status code mismatch")
			})
		}
	})

	t.Run("without anything", func(t *testing.T) {
		req := testutils.MakeReq(server.URL, "GET", "/", "")

		// execute
		res := testutils.HTTPDo(t, req)

		// test
		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "status code mismatch")
	})
}
