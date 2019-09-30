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

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func init() {
	testutils.InitTestDB()

	templatePath := fmt.Sprintf("%s/mailer/templates/src", testutils.ServerPath)
	mailer.InitTemplates(&templatePath)
}

func TestClassicPresignin(t *testing.T) {
	db := database.DBConn
	defer testutils.ClearData()

	alice := database.Account{
		Email:              database.ToNullString("alice@example.com"),
		ClientKDFIteration: 100000,
	}
	bob := database.Account{
		Email:              database.ToNullString("bob@example.com"),
		ClientKDFIteration: 200000,
	}
	testutils.MustExec(t, db.Save(&alice), "saving alice")
	testutils.MustExec(t, db.Save(&bob), "saving bob")

	testCases := []struct {
		email             string
		expectedIteration int
	}{
		{
			email:             "alice@example.com",
			expectedIteration: 100000,
		},
		{
			email:             "bob@example.com",
			expectedIteration: 200000,
		},
		{
			email: "chuck@example.com",
			// If user does not exist, reply with a generic response
			expectedIteration: 100000,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("presignin %s", tc.email), func(t *testing.T) {

			// Setup
			server := httptest.NewServer(NewRouter(&App{
				Clock: clock.NewMock(),
			}))
			defer server.Close()

			endpoint := fmt.Sprintf("/classic/presignin?email=%s", tc.email)
			req := testutils.MakeReq(server, "GET", endpoint, "")

			// Execute
			res := testutils.HTTPDo(t, req)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusOK, "")

			var got PresigninResponse
			if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
				t.Fatal(errors.Wrap(err, "decoding payload"))
			}

			assert.Equal(t, got.Iteration, tc.expectedIteration, "Iteration mismatch")
		})
	}
}

func TestClassicPresignin_MissingParams(t *testing.T) {
	defer testutils.ClearData()

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	req := testutils.MakeReq(server, "GET", "/classic/presignin", "")

	// Execute
	res := testutils.HTTPDo(t, req)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusBadRequest, "Status mismatch")
}

func TestClassicSignin(t *testing.T) {
	db := database.DBConn
	defer testutils.ClearData()

	user := testutils.SetupUserData()
	alice := testutils.SetupClassicAccountData(user, "alice@example.com")
	testutils.MustExec(t, db.Save(&alice), "saving alice")

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	dat := fmt.Sprintf(`{"email": "%s", "auth_key": "%s"}`, "alice@example.com", "/XCYisXJ6/o+vf6NUEtmrdYzJYPz+T9oAUCtMpOjhzc=")
	req := testutils.MakeReq(server, "POST", "/classic/signin", dat)

	// Execute
	res := testutils.HTTPDo(t, req)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "Status mismatch")

	var sessionCount int
	var session database.Session
	testutils.MustExec(t, db.Model(&database.Session{}).Count(&sessionCount), "counting session")
	testutils.MustExec(t, db.First(&session), "getting session")

	var got SessionResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	assert.Equal(t, sessionCount, 1, "sessionCount mismatch")
	assert.Equal(t, got.Key, session.Key, "session Key mismatch")
	assert.Equal(t, got.ExpiresAt, session.ExpiresAt.Unix(), "session ExpiresAt mismatch")

	c := testutils.GetCookieByName(res.Cookies(), "id")
	assert.Equal(t, c.Value, session.Key, "session key mismatch")
	assert.Equal(t, c.Path, "/", "session path mismatch")
	assert.Equal(t, c.HttpOnly, true, "session HTTPOnly mismatch")
	assert.Equal(t, c.Expires.Unix(), session.ExpiresAt.Unix(), "session Expires mismatch")
}

func TestClassicSignin_Failure(t *testing.T) {
	db := database.DBConn
	defer testutils.ClearData()

	//password: correctbattery
	alice := database.Account{
		Email:              database.ToNullString("alice@example.com"),
		ClientKDFIteration: 10000,
		Salt:               "Vw57HhZTqeOo0hWGb+BLoQ==",
		// plain authKey: bKAcSKkGB4VrIaSckpZvHFIlqT6L+XMVY0CTsV2y5B8=
		AuthKeyHash: "jjSs8JCaYi6cRGFPYNQ7XAVwKSrNpF1I1bGye62+A5U=",
	}
	bob := database.Account{
		Email:              database.ToNullString("bob@example.com"),
		ClientKDFIteration: 10000,
		Salt:               "gShZ7X2AuYW1xZDkpavE3g==",
		// plain authKey: DN4d/teaq1I2bVYZ7QWaah4Fu7q2y2N4yJNZk76hFHw=
		AuthKeyHash: "fGOMHHAw9G7CH4Gv2EM1ZcZZklC1a55fS3QJ0qQVp4k=",
	}
	testutils.MustExec(t, db.Save(&alice), "saving alice")
	testutils.MustExec(t, db.Save(&bob), "saving bob")

	testCases := []struct {
		email   string
		authKey string
	}{
		// missing params
		{
			email:   "",
			authKey: "",
		},
		{
			email:   "",
			authKey: "GFSymYG+s64TyHSPD3TxxMLlBBurswhDWOZRmefSoGo=",
		},
		{
			email:   "alice@example.com",
			authKey: "",
		},
		// send incorrect authKey
		{
			email:   "alice@example.com",
			authKey: "GFSymYG+s64TyHSPD3TxxMLlBBurswhDWOZRmefSoGo=",
		},
		{
			email:   "alice@example.com",
			authKey: "D8b70qEl4CXlp2DqPQpjkLrxHfYZvrwHVA6W9wTDZ6E=",
		},
		// login with mixed credentials
		{
			email:   "alice@example.com",
			authKey: "DN4d/teaq1I2bVYZ7QWaah4Fu7q2y2N4yJNZk76hFHw=",
		},
		{
			email:   "bob@example.com",
			authKey: "bKAcSKkGB4VrIaSckpZvHFIlqT6L+XMVY0CTsV2y5B8=",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("signin %s %s", tc.email, tc.authKey), func(t *testing.T) {

			// Setup
			server := httptest.NewServer(NewRouter(&App{
				Clock: clock.NewMock(),
			}))
			defer server.Close()

			dat := fmt.Sprintf(`{"email": "%s", "auth_key": "%s"}`, tc.email, tc.authKey)
			req := testutils.MakeReq(server, "POST", "/classic/signin", dat)

			// Execute
			res := testutils.HTTPDo(t, req)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")
		})
	}
}
