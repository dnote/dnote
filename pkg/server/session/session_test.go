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

package session

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
)

func TestNew(t *testing.T) {
	u1 := database.User{UUID: "0f5f0054-d23f-4be1-b5fb-57673109e9cb"}
	a1 := database.Account{Email: database.ToNullString("alice@example.com"), EmailVerified: false}

	u2 := database.User{UUID: "718a1041-bbe6-496e-bbe4-ea7e572c295e"}
	a2 := database.Account{Email: database.ToNullString("bob@example.com"), EmailVerified: false}

	testCases := []struct {
		user    database.User
		account database.Account
	}{
		{
			user:    u1,
			account: a1,
		},
		{
			user:    u2,
			account: a2,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("user %d", idx), func(t *testing.T) {
			// Execute
			got := New(tc.user, tc.account)
			expected := Session{
				UUID:          tc.user.UUID,
				Email:         tc.account.Email.String,
				EmailVerified: tc.account.EmailVerified,
			}

			assert.DeepEqual(t, got, expected, "result mismatch")
		})
	}
}
