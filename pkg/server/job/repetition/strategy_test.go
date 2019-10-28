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

package repetition

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func init() {
	testutils.InitTestDB()
}

func TestApplyBookDomain(t *testing.T) {
	defer testutils.ClearData()

	db := database.DBConn

	user := testutils.SetupUserData()
	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")
	b2 := database.Book{
		UserID: user.ID,
		Label:  "css",
	}
	testutils.MustExec(t, db.Save(&b2), "preparing b2")
	b3 := database.Book{
		UserID: user.ID,
		Label:  "golang",
	}
	testutils.MustExec(t, db.Save(&b3), "preparing b3")

	n1 := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
	}
	testutils.MustExec(t, db.Save(&n1), "preparing n1")
	n2 := database.Note{
		UserID:   user.ID,
		BookUUID: b2.UUID,
	}
	testutils.MustExec(t, db.Save(&n2), "preparing n2")
	n3 := database.Note{
		UserID:   user.ID,
		BookUUID: b3.UUID,
	}
	testutils.MustExec(t, db.Save(&n3), "preparing n3")

	var n1Record, n2Record, n3Record database.Note
	testutils.MustExec(t, db.Where("uuid = ?", n1.UUID).First(&n1Record), "finding n1")
	testutils.MustExec(t, db.Where("uuid = ?", n2.UUID).First(&n2Record), "finding n2")
	testutils.MustExec(t, db.Where("uuid = ?", n3.UUID).First(&n3Record), "finding n3")

	t.Run("book domain all", func(t *testing.T) {
		rule := database.RepetitionRule{
			UserID:     user.ID,
			BookDomain: database.BookDomainAll,
		}

		conn, err := applyBookDomain(db, rule)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing").Error())
		}

		var result []database.Note
		testutils.MustExec(t, conn.Order("id ASC").Find(&result), "finding notes")

		expected := []database.Note{n1Record, n2Record, n3Record}
		assert.DeepEqual(t, result, expected, "result mismatch")
	})

	t.Run("book domain exclude", func(t *testing.T) {
		rule := database.RepetitionRule{
			UserID:     user.ID,
			BookDomain: database.BookDomainExluding,
			Books:      []database.Book{b1},
		}
		testutils.MustExec(t, db.Save(&rule), "preparing rule")

		conn, err := applyBookDomain(db.Debug(), rule)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing").Error())
		}

		var result []database.Note
		testutils.MustExec(t, conn.Order("id ASC").Find(&result), "finding notes")

		expected := []database.Note{n2Record, n3Record}
		assert.DeepEqual(t, result, expected, "result mismatch")
	})
}
