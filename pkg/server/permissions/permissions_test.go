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

package permissions

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestViewNote(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db)
	anotherUser := testutils.SetupUserData(db)

	b1 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")

	note := database.Note{
		UUID:     testutils.MustUUID(t),
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "note content",
		Deleted:  false,
	}
	testutils.MustExec(t, db.Save(&note), "preparing note")

	t.Run("owner accessing note", func(t *testing.T) {
		result := ViewNote(&user, note)
		assert.Equal(t, result, true, "result mismatch")
	})

	t.Run("non-owner accessing note", func(t *testing.T) {
		result := ViewNote(&anotherUser, note)
		assert.Equal(t, result, false, "result mismatch")
	})

	t.Run("guest accessing note", func(t *testing.T) {
		result := ViewNote(nil, note)
		assert.Equal(t, result, false, "result mismatch")
	})
}
