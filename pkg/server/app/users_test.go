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

package app

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		onPremise   bool
		expectedPro bool
	}{
		{
			onPremise:   true,
			expectedPro: true,
		},
		{
			onPremise:   false,
			expectedPro: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("self hosting %t", tc.onPremise), func(t *testing.T) {
			c := config.Load()
			c.SetOnPremise(tc.onPremise)

			defer testutils.ClearData(testutils.DB)

			a := NewTest(&App{
				Config: c,
			})
			if _, err := a.CreateUser("alice@example.com", "pass1234"); err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			var userCount int
			var userRecord database.User
			testutils.MustExec(t, testutils.DB.Model(&database.User{}).Count(&userCount), "counting user")
			testutils.MustExec(t, testutils.DB.First(&userRecord), "finding user")

			assert.Equal(t, userCount, 1, "book count mismatch")
			assert.Equal(t, userRecord.Cloud, tc.expectedPro, "user pro mismatch")
		})
	}
}
