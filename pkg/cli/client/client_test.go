/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/testutils"
	"github.com/pkg/errors"
)

// startCommonTestServer starts a test HTTP server that simulates a common set of senarios
func startCommonTestServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// internal server error
		if r.URL.String() == "/bad-api/v3/signout" && r.Method == "POST" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// catch-all
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><div id="app-root"></div></body></html>`))
	}))

	return ts
}

func TestSignIn(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/api/v3/signin" && r.Method == "POST" {
			var payload SigninPayload

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				t.Fatalf(errors.Wrap(err, "decoding payload in the test server").Error())
				return
			}

			if payload.Email == "alice@example.com" && payload.Passowrd == "pass1234" {
				resp := testutils.MustMarshalJSON(t, SigninResponse{
					Key:       "somekey",
					ExpiresAt: int64(1596439890),
				})

				w.Header().Set("Content-Type", "application/json")
				w.Write(resp)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}

			return
		}
	}))
	defer ts.Close()

	commonTs := startCommonTestServer()
	defer commonTs.Close()

	correctEndpoint := fmt.Sprintf("%s/api", ts.URL)

	t.Run("success", func(t *testing.T) {
		result, err := Signin(context.DnoteCtx{APIEndpoint: correctEndpoint}, "alice@example.com", "pass1234")
		if err != nil {
			t.Errorf("got signin request error: %+v", err.Error())
		}

		assert.Equal(t, result.Key, "somekey", "Key mismatch")
		assert.Equal(t, result.ExpiresAt, int64(1596439890), "ExpiresAt mismatch")
	})

	t.Run("failure", func(t *testing.T) {
		result, err := Signin(context.DnoteCtx{APIEndpoint: correctEndpoint}, "alice@example.com", "incorrectpassword")

		assert.Equal(t, err, ErrInvalidLogin, "err mismatch")
		assert.Equal(t, result.Key, "", "Key mismatch")
		assert.Equal(t, result.ExpiresAt, int64(0), "ExpiresAt mismatch")
	})

	t.Run("server error", func(t *testing.T) {
		endpoint := fmt.Sprintf("%s/bad-api", ts.URL)
		result, err := Signin(context.DnoteCtx{APIEndpoint: endpoint}, "alice@example.com", "pass1234")
		if err == nil {
			t.Error("error should have been returned")
		}

		assert.Equal(t, result.Key, "", "Key mismatch")
		assert.Equal(t, result.ExpiresAt, int64(0), "ExpiresAt mismatch")
	})

	t.Run("accidentally pointing to a catch-all handler", func(t *testing.T) {
		endpoint := fmt.Sprintf("%s", ts.URL)
		result, err := Signin(context.DnoteCtx{APIEndpoint: endpoint}, "alice@example.com", "pass1234")

		assert.Equal(t, errors.Cause(err), ErrContentTypeMismatch, "error cause mismatch")
		assert.Equal(t, result.Key, "", "Key mismatch")
		assert.Equal(t, result.ExpiresAt, int64(0), "ExpiresAt mismatch")
	})
}

func TestSignOut(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/api/v3/signout" && r.Method == "POST" {
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer ts.Close()

	commonTs := startCommonTestServer()
	defer commonTs.Close()

	correctEndpoint := fmt.Sprintf("%s/api", ts.URL)

	t.Run("success", func(t *testing.T) {
		err := Signout(context.DnoteCtx{SessionKey: "somekey", APIEndpoint: correctEndpoint}, "alice@example.com")
		if err != nil {
			t.Errorf("got signin request error: %+v", err.Error())
		}
	})

	t.Run("server error", func(t *testing.T) {
		endpoint := fmt.Sprintf("%s/bad-api", commonTs.URL)
		err := Signout(context.DnoteCtx{SessionKey: "somekey", APIEndpoint: endpoint}, "alice@example.com")
		if err == nil {
			t.Error("error should have been returned")
		}
	})

	t.Run("accidentally pointing to a catch-all handler", func(t *testing.T) {
		endpoint := fmt.Sprintf("%s", commonTs.URL)
		err := Signout(context.DnoteCtx{SessionKey: "somekey", APIEndpoint: endpoint}, "alice@example.com")

		assert.Equal(t, errors.Cause(err), ErrContentTypeMismatch, "error cause mismatch")
	})
}
