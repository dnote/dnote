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

package operations

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func init() {
	testutils.InitTestDB()
}

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
			defer testutils.ClearData()
			db := database.DBConn

			user := testutils.SetupUserData()
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
