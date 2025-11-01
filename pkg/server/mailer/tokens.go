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

package mailer

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func generateRandomToken(bits int) (string, error) {
	b := make([]byte, bits)

	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "generating random bytes")
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// GetToken returns an token of the given kind for the user
// by first looking up any unused record and creating one if none exists.
func GetToken(db *gorm.DB, userID int, kind string) (database.Token, error) {
	var tok database.Token
	err := db.
		Where("user_id = ? AND type =? AND used_at IS NULL", userID, kind).
		First(&tok).Error

	tokenVal, genErr := generateRandomToken(16)
	if genErr != nil {
		return tok, errors.Wrap(genErr, "generating token value")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		tok = database.Token{
			UserID: userID,
			Type:   kind,
			Value:  tokenVal,
		}
		if err := db.Save(&tok).Error; err != nil {
			return tok, errors.Wrap(err, "saving token")
		}

		return tok, nil
	} else if err != nil {
		return tok, errors.Wrap(err, "finding token")
	}

	return tok, nil
}
