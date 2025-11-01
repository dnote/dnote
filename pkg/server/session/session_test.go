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

package session

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
)

func TestNew(t *testing.T) {
	u1 := database.User{
		UUID:  "0f5f0054-d23f-4be1-b5fb-57673109e9cb",
		Email: database.ToNullString("alice@example.com"),
	}

	u2 := database.User{
		UUID:  "718a1041-bbe6-496e-bbe4-ea7e572c295e",
		Email: database.ToNullString("bob@example.com"),
	}

	testCases := []struct {
		user database.User
	}{
		{
			user: u1,
		},
		{
			user: u2,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("user %d", idx), func(t *testing.T) {
			// Execute
			got := New(tc.user)
			expected := Session{
				UUID:  tc.user.UUID,
				Email: tc.user.Email.String,
			}

			assert.DeepEqual(t, got, expected, "result mismatch")
		})
	}
}
