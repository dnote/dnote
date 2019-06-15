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

package mailer

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/pkg/errors"
)

func generateRandomToken(bits int) (string, error) {
	b := make([]byte, bits)

	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "generating random bytes")
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// GetEmailPreferenceToken returns an unused email frequency token for the user
// by first looking up any existing record and creating one if none exists.
func GetEmailPreferenceToken(user database.User) (database.Token, error) {
	db := database.DBConn

	var tok database.Token
	conn := db.
		Where("user_id = ? AND type =? AND used_at IS NULL", user.ID, database.TokenTypeEmailPreference).
		First(&tok)

	tokenVal, err := generateRandomToken(16)
	if err != nil {
		return tok, errors.Wrap(err, "generating token value")
	}

	if conn.RecordNotFound() {
		tok = database.Token{
			UserID: user.ID,
			Type:   database.TokenTypeEmailPreference,
			Value:  tokenVal,
		}
		if err := db.Save(&tok).Error; err != nil {
			return tok, errors.Wrap(err, "saving token")
		}

		return tok, nil
	} else if err := conn.Error; err != nil {
		return tok, errors.Wrap(err, "finding token")
	}

	return tok, nil
}
