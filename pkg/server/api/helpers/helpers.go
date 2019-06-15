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

package helpers

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const (
	demoUserEmail = "demo@dnote.io"
)

// GetDemoUserID returns ID of the demo user
func GetDemoUserID() (int, error) {
	db := database.DBConn

	result := struct {
		UserID int
	}{}
	if err := db.Table("accounts").Select("user_id").Where("email = ?", demoUserEmail).Scan(&result).Error; err != nil {
		return result.UserID, errors.Wrap(err, "finding demo user")
	}

	return result.UserID, nil
}

// GenUUID generates a new uuid
func GenUUID() string {
	return uuid.NewV4().String()
}
