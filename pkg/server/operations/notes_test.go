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

package operations

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestGetNote(t *testing.T) {
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

	var privateNoteRecord, publicNoteRecord database.Note
	testutils.MustExec(t, testutils.DB.Where("uuid = ?", privateNote.UUID).Preload("Book").Preload("User").First(&privateNoteRecord), "finding privateNote")
	testutils.MustExec(t, testutils.DB.Where("uuid = ?", publicNote.UUID).Preload("Book").Preload("User").First(&publicNoteRecord), "finding publicNote")

	testCases := []struct {
		name         string
		user         database.User
		note         database.Note
		expectedOK   bool
		expectedNote database.Note
	}{
		{
			name:         "owner accessing private note",
			user:         user,
			note:         privateNote,
			expectedOK:   true,
			expectedNote: privateNoteRecord,
		},
		{
			name:         "non-owner accessing private note",
			user:         anotherUser,
			note:         privateNote,
			expectedOK:   false,
			expectedNote: database.Note{},
		},
		{
			name:         "non-owner accessing public note",
			user:         anotherUser,
			note:         publicNote,
			expectedOK:   true,
			expectedNote: publicNoteRecord,
		},
		{
			name:         "guest accessing private note",
			user:         database.User{},
			note:         privateNote,
			expectedOK:   false,
			expectedNote: database.Note{},
		},
		{
			name:         "guest accessing public note",
			user:         database.User{},
			note:         publicNote,
			expectedOK:   true,
			expectedNote: publicNoteRecord,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			note, ok, err := GetNote(testutils.DB, tc.note.UUID, &tc.user)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			assert.Equal(t, ok, tc.expectedOK, "ok mismatch")
			assert.DeepEqual(t, note, tc.expectedNote, "note mismatch")
		})
	}
}

func TestGetNote_nonexistent(t *testing.T) {
	user := testutils.SetupUserData()

	defer testutils.ClearData(testutils.DB)

	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")

	n1UUID := "4fd19336-671e-4ff3-8f22-662b80e22edc"
	n1 := database.Note{
		UUID:     n1UUID,
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n1 content",
		Deleted:  false,
		Public:   false,
	}
	testutils.MustExec(t, testutils.DB.Save(&n1), "preparing n1")

	nonexistentUUID := "4fd19336-671e-4ff3-8f22-662b80e22edd"
	note, ok, err := GetNote(testutils.DB, nonexistentUUID, &user)
	if err != nil {
		t.Fatal(errors.Wrap(err, "executing"))
	}

	assert.Equal(t, ok, false, "ok mismatch")
	assert.DeepEqual(t, note, database.Note{}, "note mismatch")
}
