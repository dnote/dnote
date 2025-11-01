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

package app

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestIncremenetUserUSN(t *testing.T) {
	testCases := []struct {
		maxUSN         int
		expectedMaxUSN int
	}{
		{
			maxUSN:         1,
			expectedMaxUSN: 2,
		},
		{
			maxUSN:         1988,
			expectedMaxUSN: 1989,
		},
	}

	// set up
	for idx, tc := range testCases {
		func() {
			db := testutils.InitMemoryDB(t)

			user := testutils.SetupUserData(db, "user@test.com", "password123")
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.maxUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			// execute
			tx := db.Begin()
			nextUSN, err := incrementUserUSN(tx, user.ID)
			if err != nil {
				t.Fatal(errors.Wrap(err, "incrementing the user usn"))
			}
			tx.Commit()

			// test
			var userRecord database.User
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), fmt.Sprintf("finding user for test case %d", idx))

			assert.Equal(t, userRecord.MaxUSN, tc.expectedMaxUSN, fmt.Sprintf("user max_usn mismatch for case %d", idx))
			assert.Equal(t, nextUSN, tc.expectedMaxUSN, fmt.Sprintf("next_usn mismatch for case %d", idx))
		}()
	}
}
