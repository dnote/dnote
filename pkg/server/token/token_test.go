/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
			u := testutils.SetupUserData(db, "user@test.com", "password123")

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
