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

package app

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/pkg/errors"
)

func (a *App) getExistingDigestReceipt(userID, digestID int) (*database.DigestReceipt, error) {
	var ret database.DigestReceipt
	conn := a.DB.Where("user_id = ? AND digest_id = ?", userID, digestID).First(&ret)

	if conn.RecordNotFound() {
		return nil, nil
	}
	if err := conn.Error; err != nil {
		return nil, errors.Wrap(err, "querying existing digest receipt")
	}

	return &ret, nil
}

// GetUserDigestByUUID retrives a digest by the uuid for the given user
func (a *App) GetUserDigestByUUID(userID int, uuid string) (*database.Digest, error) {
	var ret database.Digest
	conn := a.DB.Where("user_id = ? AND uuid = ?", userID, uuid).First(&ret)

	if conn.RecordNotFound() {
		return nil, nil
	}
	if err := conn.Error; err != nil {
		return nil, errors.Wrap(err, "finding digest")
	}

	return &ret, nil
}

// MarkDigestRead creates a new digest receipt. If one already exists for
// the given digest and the user, it is a noop.
func (a *App) MarkDigestRead(digest database.Digest, user database.User) (database.DigestReceipt, error) {
	db := a.DB

	existing, err := a.getExistingDigestReceipt(user.ID, digest.ID)
	if err != nil {
		return database.DigestReceipt{}, errors.Wrap(err, "checking existing digest receipt")
	}
	if existing != nil {
		return *existing, nil
	}

	dat := database.DigestReceipt{
		UserID:   user.ID,
		DigestID: digest.ID,
	}
	if err := db.Create(&dat).Error; err != nil {
		return database.DigestReceipt{}, errors.Wrap(err, "creating digest receipt")
	}

	return dat, nil
}
