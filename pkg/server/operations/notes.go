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
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/permissions"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// GetNote retrieves a note for the given user
func GetNote(db *gorm.DB, uuid string, user *database.User) (database.Note, bool, error) {
	zeroNote := database.Note{}
	if !helpers.ValidateUUID(uuid) {
		return zeroNote, false, nil
	}

	var note database.Note
	err := database.PreloadNote(db.Where("notes.uuid = ? AND deleted = ?", uuid, false)).Find(&note).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return zeroNote, false, nil
	} else if err != nil {
		return zeroNote, false, errors.Wrap(err, "finding note")
	}

	if ok := permissions.ViewNote(user, note); !ok {
		return zeroNote, false, nil
	}

	return note, true, nil
}
