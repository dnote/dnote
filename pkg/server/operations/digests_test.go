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
	// "fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestCreateDigest(t *testing.T) {
	t.Run("no previous digest", func(t *testing.T) {
		defer testutils.ClearData()

		db := testutils.DB

		user := testutils.SetupUserData()
		rule := database.RepetitionRule{UserID: user.ID}
		testutils.MustExec(t, testutils.DB.Save(&rule), "preparing rule")

		result, err := CreateDigest(db, rule, nil)
		if err != nil {
			t.Fatal(errors.Wrap(err, "performing"))
		}

		assert.Equal(t, result.Version, 1, "Version mismatch")
	})

	t.Run("with previous digest", func(t *testing.T) {
		defer testutils.ClearData()

		db := testutils.DB

		user := testutils.SetupUserData()
		rule := database.RepetitionRule{UserID: user.ID}
		testutils.MustExec(t, testutils.DB.Save(&rule), "preparing rule")

		d := database.Digest{UserID: user.ID, RuleID: rule.ID, Version: 8}
		testutils.MustExec(t, testutils.DB.Save(&d), "preparing digest")

		result, err := CreateDigest(db, rule, nil)
		if err != nil {
			t.Fatal(errors.Wrap(err, "performing"))
		}

		assert.Equal(t, result.Version, 9, "Version mismatch")
	})
}
