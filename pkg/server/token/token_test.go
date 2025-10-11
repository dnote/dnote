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

package token

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		kind string
	}{
		{
			kind: database.TokenTypeResetPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("token type %s", tc.kind), func(t *testing.T) {
			db := testutils.InitMemoryDB(t)

			// Set up
			u := testutils.SetupUserData(db)

			// Execute
			tok, err := Create(db, u.ID, tc.kind)
			if err != nil {
				t.Fatal(errors.Wrap(err, "performing"))
			}

			// Test
			var count int64
			testutils.MustExec(t, db.Model(&database.Token{}).Count(&count), "counting token")
			assert.Equalf(t, count, int64(1), "error mismatch")

			var tokenRecord database.Token
			testutils.MustExec(t, db.First(&tokenRecord), "finding token")
			assert.Equalf(t, tokenRecord.UserID, tok.UserID, "UserID mismatch")
			assert.Equalf(t, tokenRecord.Value, tok.Value, "Value mismatch")
			assert.Equalf(t, tokenRecord.Type, tok.Type, "Type mismatch")
		})
	}
}
