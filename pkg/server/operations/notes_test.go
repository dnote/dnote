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

package operations

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestGetNote(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db, "user@test.com", "password123")
	anotherUser := testutils.SetupUserData(db, "another@test.com", "password123")

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

	var noteRecord database.Note
	testutils.MustExec(t, db.Where("uuid = ?", note.UUID).Preload("Book").Preload("User").First(&noteRecord), "finding note")

	testCases := []struct {
		name         string
		user         database.User
		note         database.Note
		expectedOK   bool
		expectedNote database.Note
	}{
		{
			name:         "owner accessing note",
			user:         user,
			note:         note,
			expectedOK:   true,
			expectedNote: noteRecord,
		},
		{
			name:         "non-owner accessing note",
			user:         anotherUser,
			note:         note,
			expectedOK:   false,
			expectedNote: database.Note{},
		},
		{
			name:         "guest accessing note",
			user:         database.User{},
			note:         note,
			expectedOK:   false,
			expectedNote: database.Note{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			note, ok, err := GetNote(db, tc.note.UUID, &tc.user)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			assert.Equal(t, ok, tc.expectedOK, "ok mismatch")
			assert.DeepEqual(t, note, tc.expectedNote, "note mismatch")
		})
	}
}

func TestGetNote_nonexistent(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db, "user@test.com", "password123")

	b1 := database.Book{
		UUID:   testutils.MustUUID(t),
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")

	n1 := database.Note{
		UUID:     "4fd19336-671e-4ff3-8f22-662b80e22edc",
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "n1 content",
		Deleted:  false,
	}
	testutils.MustExec(t, db.Save(&n1), "preparing n1")

	nonexistentUUID := "4fd19336-671e-4ff3-8f22-662b80e22edd"
	note, ok, err := GetNote(db, nonexistentUUID, &user)
	if err != nil {
		t.Fatal(errors.Wrap(err, "executing"))
	}

	assert.Equal(t, ok, false, "ok mismatch")
	assert.DeepEqual(t, note, database.Note{}, "note mismatch")
}
