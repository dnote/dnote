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

package permissions

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestCheckPlanAllowance_Pro(t *testing.T) {
	t.Run("has less than 5 books", func(t *testing.T) {
		defer testutils.ClearData()

		user := testutils.SetupUserData()
		testutils.MustExec(t, testutils.DB.Model(&user).Update("cloud", true), "preparing user")

		testutils.PrepareBooks(t, user, 3)

		ok, err := CheckPlanAllowance(testutils.DB, user)
		if err != nil {
			t.Fatal(errors.Wrap(err, "performing"))
		}

		assert.Equal(t, ok, true, "result mismatch")
	})

	t.Run("has more than 5 books", func(t *testing.T) {
		defer testutils.ClearData()

		user := testutils.SetupUserData()
		testutils.MustExec(t, testutils.DB.Model(&user).Update("cloud", true), "preparing user")

		testutils.PrepareBooks(t, user, 6)

		ok, err := CheckPlanAllowance(testutils.DB, user)
		if err != nil {
			t.Fatal(errors.Wrap(err, "performing"))
		}

		assert.Equal(t, ok, true, "result mismatch")
	})
}

func TestCheckPlanAllowance_Core(t *testing.T) {
	t.Run("has less than 5 books", func(t *testing.T) {
		defer testutils.ClearData()

		user := testutils.SetupUserData()
		testutils.MustExec(t, testutils.DB.Model(&user).Update("cloud", false), "preparing user")

		testutils.PrepareBooks(t, user, 3)

		ok, err := CheckPlanAllowance(testutils.DB, user)
		if err != nil {
			t.Fatal(errors.Wrap(err, "performing"))
		}

		assert.Equal(t, ok, true, "result mismatch")
	})

	t.Run("has 5 books", func(t *testing.T) {
		defer testutils.ClearData()

		user := testutils.SetupUserData()
		testutils.MustExec(t, testutils.DB.Model(&user).Update("cloud", false), "preparing user")

		testutils.PrepareBooks(t, user, 5)

		ok, err := CheckPlanAllowance(testutils.DB, user)
		if err != nil {
			t.Fatal(errors.Wrap(err, "performing"))
		}

		assert.Equal(t, ok, false, "result mismatch")
	})

	t.Run("has more than 5 books", func(t *testing.T) {
		defer testutils.ClearData()

		user := testutils.SetupUserData()
		testutils.MustExec(t, testutils.DB.Model(&user).Update("cloud", false), "preparing user")

		testutils.PrepareBooks(t, user, 6)

		ok, err := CheckPlanAllowance(testutils.DB, user)
		if err != nil {
			t.Fatal(errors.Wrap(err, "performing"))
		}

		assert.Equal(t, ok, false, "result mismatch")
	})
}
