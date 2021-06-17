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

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"testing"
// 	"time"
//
// 	"github.com/dnote/dnote/pkg/assert"
// 	"github.com/dnote/dnote/pkg/clock"
// 	"github.com/dnote/dnote/pkg/server/app"
// 	"github.com/dnote/dnote/pkg/server/database"
// 	"github.com/dnote/dnote/pkg/server/session"
// 	"github.com/dnote/dnote/pkg/server/testutils"
// 	"github.com/pkg/errors"
// 	"golang.org/x/crypto/bcrypt"
// )
//
// func TestGetMe(t *testing.T) {
// 	testutils.InitTestDB()
// 	defer testutils.ClearData(testutils.DB)
//
// 	// Setup
// 	server := MustNewServer(t, &app.App{
// 		Clock: clock.NewMock(),
// 	})
// 	defer server.Close()
//
// 	u1 := testutils.SetupUserData()
// 	a1 := testutils.SetupAccountData(u1, "alice@example.com", "somepassword")
//
// 	u2 := testutils.SetupUserData()
// 	testutils.MustExec(t, testutils.DB.Model(&u2).Update("cloud", false), "preparing u2 cloud")
// 	a2 := testutils.SetupAccountData(u2, "bob@example.com", "somepassword")
//
// 	testCases := []struct {
// 		user        database.User
// 		account     database.Account
// 		expectedPro bool
// 	}{
// 		{
// 			user:        u1,
// 			account:     a1,
// 			expectedPro: true,
// 		},
// 		{
// 			user:        u2,
// 			account:     a2,
// 			expectedPro: false,
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(fmt.Sprintf("user pro %t", tc.expectedPro), func(t *testing.T) {
// 			// Execute
// 			req := testutils.MakeReq(server.URL, "GET", "/me", "")
// 			res := testutils.HTTPAuthDo(t, req, tc.user)
//
// 			// Test
// 			assert.StatusCodeEquals(t, res, http.StatusOK, "")
//
// 			var payload GetMeResponse
// 			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
// 				t.Fatal(errors.Wrap(err, "decoding payload"))
// 			}
//
// 			expectedPayload := GetMeResponse{
// 				User: session.Session{
// 					UUID:          tc.user.UUID,
// 					Pro:           tc.expectedPro,
// 					Email:         tc.account.Email.String,
// 					EmailVerified: tc.account.EmailVerified,
// 				},
// 			}
// 			assert.DeepEqual(t, payload, expectedPayload, "payload mismatch")
//
// 			var user database.User
// 			testutils.MustExec(t, testutils.DB.Where("id = ?", tc.user.ID).First(&user), "finding user")
// 			assert.NotEqual(t, user.LastLoginAt, nil, "LastLoginAt mismatch")
// 		})
// 	}
// }
//
// func TestCreateResetToken(t *testing.T) {
// 	t.Run("success", func(t *testing.T) {
// 		defer testutils.ClearData(testutils.DB)
//
// 		// Setup
// 		server := MustNewServer(t, &app.App{
//
// 			Clock: clock.NewMock(),
// 		})
// 		defer server.Close()
//
// 		u := testutils.SetupUserData()
// 		testutils.SetupAccountData(u, "alice@example.com", "somepassword")
//
// 		dat := `{"email": "alice@example.com"}`
// 		req := testutils.MakeReq(server.URL, "POST", "/reset-token", dat)
//
// 		// Execute
// 		res := testutils.HTTPDo(t, req)
//
// 		// Test
// 		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismtach")
//
// 		var tokenCount int
// 		testutils.MustExec(t, testutils.DB.Model(&database.Token{}).Count(&tokenCount), "counting tokens")
//
// 		var resetToken database.Token
// 		testutils.MustExec(t, testutils.DB.Where("user_id = ? AND type = ?", u.ID, database.TokenTypeResetPassword).First(&resetToken), "finding reset token")
//
// 		assert.Equal(t, tokenCount, 1, "reset_token count mismatch")
// 		assert.NotEqual(t, resetToken.Value, nil, "reset_token value mismatch")
// 		assert.Equal(t, resetToken.UsedAt, (*time.Time)(nil), "reset_token UsedAt mismatch")
// 	})
//
// 	t.Run("nonexistent email", func(t *testing.T) {
//
// 		defer testutils.ClearData(testutils.DB)
//
// 		// Setup
// 		server := MustNewServer(t, &app.App{
//
// 			Clock: clock.NewMock(),
// 		})
// 		defer server.Close()
//
// 		u := testutils.SetupUserData()
// 		testutils.SetupAccountData(u, "alice@example.com", "somepassword")
//
// 		dat := `{"email": "bob@example.com"}`
// 		req := testutils.MakeReq(server.URL, "POST", "/reset-token", dat)
//
// 		// Execute
// 		res := testutils.HTTPDo(t, req)
//
// 		// Test
// 		assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismtach")
//
// 		var tokenCount int
// 		testutils.MustExec(t, testutils.DB.Model(&database.Token{}).Count(&tokenCount), "counting tokens")
// 		assert.Equal(t, tokenCount, 0, "reset_token count mismatch")
// 	})
// }
//
