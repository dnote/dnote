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

package app

import (
	"time"

	"github.com/dnote/dnote/pkg/server/crypt"
	"github.com/dnote/dnote/pkg/server/database"
	"gorm.io/gorm"
	"github.com/pkg/errors"
)

// CreateSession returns a new session for the user of the given id
func (a *App) CreateSession(userID int) (database.Session, error) {
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

	if err := a.DB.Save(&session).Error; err != nil {
		return database.Session{}, errors.Wrap(err, "saving session")
	}

	return session, nil
}

// DeleteUserSessions deletes all existing sessions for the given user. It effectively
// invalidates all existing sessions.
func (a *App) DeleteUserSessions(db *gorm.DB, userID int) error {
	if err := db.Where("user_id = ?", userID).Delete(&database.Session{}).Error; err != nil {
		return errors.Wrap(err, "deleting sessions")
	}

	return nil
}

// DeleteSession deletes the session that match the given info
func (a *App) DeleteSession(sessionKey string) error {
	if err := a.DB.Where("key = ?", sessionKey).Delete(&database.Session{}).Error; err != nil {
		return errors.Wrap(err, "deleting the session")
	}

	return nil
}
