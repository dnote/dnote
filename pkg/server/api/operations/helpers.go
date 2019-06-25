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
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// incrementUserUSN increment the given user's max_usn by 1
// and returns the new, incremented max_usn
func incrementUserUSN(tx *gorm.DB, userID int) (int, error) {
	if err := tx.Table("users").Where("id = ?", userID).Update("max_usn", gorm.Expr("max_usn + 1")).Error; err != nil {
		return 0, errors.Wrap(err, "incrementing user max_usn")
	}

	var user database.User
	if err := tx.Select("max_usn").Where("id = ?", userID).First(&user).Error; err != nil {
		return 0, errors.Wrap(err, "getting the updated user max_usn")
	}

	return user.MaxUSN, nil
}
