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

package token

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// generateRandom generates random bits of given length
func generateRandom(bits int) (string, error) {
	b := make([]byte, bits)

	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "reading random bytes")
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// Create generates a new token in the database
func Create(db *gorm.DB, userID int, kind string) (database.Token, error) {
	val, err := generateRandom(16)
	if err != nil {
		return database.Token{}, errors.Wrap(err, "generating random bytes")
	}

	token := database.Token{
		UserID: userID,
		Value:  val,
		Type:   kind,
	}
	if err := db.Save(&token).Error; err != nil {
		return database.Token{}, errors.Wrap(err, "creating a token for unsubscribing")
	}

	return token, nil
}
