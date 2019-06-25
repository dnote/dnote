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

func generateResetToken() (string, error) {
	ret, err := crypt.GetRandomStr(16)
	if err != nil {
		return "", errors.Wrap(err, "generating random token")
	}

	return ret, nil
}

func generateVerificationCode() (string, error) {
	ret, err := crypt.GetRandomStr(16)
	if err != nil {
		return "", errors.Wrap(err, "generating random token")
	}

	return ret, nil
}

// TouchLastLoginAt updates the last login timestamp
func TouchLastLoginAt(user database.User, tx *gorm.DB) error {
	t := time.Now()
	if err := tx.Model(&user).Update(database.User{LastLoginAt: &t}).Error; err != nil {
		return errors.Wrap(err, "updating last_login_at")
	}

	return nil
}

func createEmailVerificaitonToken(user database.User, tx *gorm.DB) error {
	verificationCode, err := generateVerificationCode()
	if err != nil {
		return errors.Wrap(err, "generating verification code")
	}
	token := database.Token{
		UserID: user.ID,
		Type:   database.TokenTypeEmailVerification,
		Value:  verificationCode,
	}
	if err := tx.Save(&token).Error; err != nil {
		return errors.Wrap(err, "saving verification token")
	}

	return nil
}

func createEmailPreference(user database.User, tx *gorm.DB) error {
	EmailPreference := database.EmailPreference{
		UserID:       user.ID,
		DigestWeekly: true,
	}
	if err := tx.Save(&EmailPreference).Error; err != nil {
		return errors.Wrap(err, "inserting email preference")
	}

	return nil
}

// CreateUser creates a user
func CreateUser(tx *gorm.DB, email, authKey, cipherKeyEnc string, iteration int) (database.User, error) {
	salt, err := crypt.GetRandomStr(16)
	if err != nil {
		return database.User{}, errors.Wrap(err, "generating salt")
	}

	user := database.User{}
	if err = tx.Save(&user).Error; err != nil {
		return database.User{}, errors.Wrap(err, "saving user")
	}
	account := database.Account{
		// TODO: email should not be nullable.
		Email:              database.ToNullString(email),
		UserID:             user.ID,
		AuthKeyHash:        crypt.HashAuthKey(authKey, salt, crypt.ServerKDFIteration),
		Salt:               salt,
		ClientKDFIteration: iteration,
		ServerKDFIteration: crypt.ServerKDFIteration,
		CipherKeyEnc:       cipherKeyEnc,
	}
	if err = tx.Save(&account).Error; err != nil {
		return database.User{}, errors.Wrap(err, "saving account")
	}

	if err := createEmailVerificaitonToken(user, tx); err != nil {
		return database.User{}, errors.Wrap(err, "creating email verificaiton token")
	}
	if err := createEmailPreference(user, tx); err != nil {
		return database.User{}, errors.Wrap(err, "creating email preference")
	}
	if err := TouchLastLoginAt(user, tx); err != nil {
		return database.User{}, errors.Wrap(err, "updating last login")
	}
	return user, nil
}

// LegacyRegisterUser migrates the given user to the encrypted user
func LegacyRegisterUser(tx *gorm.DB, userID int, email, authKey string, cipherKeyEnc string, iteration int) error {
	salt, err := crypt.GetRandomStr(16)
	if err != nil {
		return errors.Wrap(err, "generating salt")
	}

	var account database.Account
	if err := tx.Where("user_id = ?", userID).First(&account).Error; err != nil {
		return errors.Wrap(err, "finding account")
	}

	account.Email = database.ToNullString(email)
	account.AuthKeyHash = crypt.HashAuthKey(authKey, salt, crypt.ServerKDFIteration)
	account.Salt = salt
	account.ClientKDFIteration = iteration
	account.ServerKDFIteration = crypt.ServerKDFIteration
	account.CipherKeyEnc = cipherKeyEnc
	account.Password = database.ToNullString("")

	if err = tx.Save(&account).Error; err != nil {
		return errors.Wrap(err, "saving account")
	}

	return nil
}
