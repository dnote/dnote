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
	"time"

	"github.com/dnote/dnote/pkg/server/api/crypt"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// CreateSession returns a new session for the user of the given id
func CreateSession(db *gorm.DB, userID int) (database.Session, error) {
	key, err := crypt.GetRandomStr(32)
	if err != nil {
		return database.Session{}, errors.Wrap(err, "generating key")
	}

	session := database.Session{
		UserID:     userID,
		Key:        key,
		LastUsedAt: time.Now(),
		ExpiresAt:  time.Now().Add(24 * 100 * time.Hour),
	}

	if err := db.Save(&session).Error; err != nil {
		return database.Session{}, errors.Wrap(err, "saving session")
	}

	return session, nil
}

// DeleteUserSessions deletes all existing sessions for the given user. It effectively
// invalidates all existing sessions.
func DeleteUserSessions(db *gorm.DB, userID int) error {
	if err := db.Where("user_id = ?", userID).Delete(&database.Session{}).Error; err != nil {
		return errors.Wrap(err, "deleting sessions")
	}

	return nil
}

// DeleteSession deletes the session that match the given info
func DeleteSession(db *gorm.DB, sessionKey string) error {
	if err := db.Debug().Where("key = ?", sessionKey).Delete(&database.Session{}).Error; err != nil {
		return errors.Wrap(err, "deleting the session")
	}

	return nil
}
