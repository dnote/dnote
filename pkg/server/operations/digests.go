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

package operations

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// CreateDigest creates a new digest
func CreateDigest(db *gorm.DB, rule database.RepetitionRule, notes []database.Note) (database.Digest, error) {
	var maxVersion int
	if err := db.Raw("SELECT COALESCE(max(version), 0) FROM digests WHERE rule_id = ?", rule.ID).Row().Scan(&maxVersion); err != nil {
		return database.Digest{}, errors.Wrap(err, "finding max version")
	}

	digest := database.Digest{
		RuleID:  rule.ID,
		UserID:  rule.UserID,
		Version: maxVersion + 1,
		Notes:   notes,
	}
	if err := db.Save(&digest).Error; err != nil {
		return database.Digest{}, errors.Wrap(err, "saving digest")
	}

	return digest, nil
}
