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

package permissions

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestViewNote(t *testing.T) {
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
