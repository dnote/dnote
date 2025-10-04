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
