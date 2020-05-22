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

package permissions

import (
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestMain(m *testing.M) {
	testutils.InitTestDB()

	code := m.Run()
	testutils.ClearData(testutils.DB)

	os.Exit(code)
}

func TestViewNote(t *testing.T) {
	user := testutils.SetupUserData()
	anotherUser := testutils.SetupUserData()

	defer testutils.ClearData(testutils.DB)

	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")

	privateNote := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "privateNote content",
		Deleted:  false,
		Public:   false,
	}
	testutils.MustExec(t, testutils.DB.Save(&privateNote), "preparing privateNote")

	publicNote := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "privateNote content",
		Deleted:  false,
		Public:   true,
	}
	testutils.MustExec(t, testutils.DB.Save(&publicNote), "preparing privateNote")

	t.Run("owner accessing private note", func(t *testing.T) {
		result := ViewNote(&user, privateNote)
		assert.Equal(t, result, true, "result mismatch")
	})

	t.Run("owner accessing public note", func(t *testing.T) {
		result := ViewNote(&user, publicNote)
		assert.Equal(t, result, true, "result mismatch")
	})

	t.Run("non-owner accessing private note", func(t *testing.T) {
		result := ViewNote(&anotherUser, privateNote)
		assert.Equal(t, result, false, "result mismatch")
	})

	t.Run("non-owner accessing public note", func(t *testing.T) {
		result := ViewNote(&anotherUser, publicNote)
		assert.Equal(t, result, true, "result mismatch")
	})

	t.Run("guest accessing private note", func(t *testing.T) {
		result := ViewNote(nil, privateNote)
		assert.Equal(t, result, false, "result mismatch")
	})

	t.Run("guest accessing public note", func(t *testing.T) {
		result := ViewNote(nil, publicNote)
		assert.Equal(t, result, true, "result mismatch")
	})
}
