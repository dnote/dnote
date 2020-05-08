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

package app

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/token"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// TouchLastLoginAt updates the last login timestamp
func (a *App) TouchLastLoginAt(user database.User, tx *gorm.DB) error {
	t := a.Clock.Now()
	if err := tx.Model(&user).Update(database.User{LastLoginAt: &t}).Error; err != nil {
		return errors.Wrap(err, "updating last_login_at")
	}

	return nil
}

func createEmailPreference(user database.User, tx *gorm.DB) error {
	p := database.EmailPreference{
		UserID: user.ID,
	}
	if err := tx.Save(&p).Error; err != nil {
		return errors.Wrap(err, "inserting email preference")
	}

	return nil
}

// CreateUser creates a user
func (a *App) CreateUser(email, password string) (database.User, error) {
	tx := a.DB.Begin()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return database.User{}, errors.Wrap(err, "hashing password")
	}

	// Grant all privileges if self-hosting
	var pro bool
	if a.Config.OnPremise {
		pro = true
	} else {
		pro = false
	}

	user := database.User{
		Cloud: pro,
	}
	if err = tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return database.User{}, errors.Wrap(err, "saving user")
	}
	account := database.Account{
		Email:    database.ToNullString(email),
		Password: database.ToNullString(string(hashedPassword)),
		UserID:   user.ID,
	}
	if err = tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return database.User{}, errors.Wrap(err, "saving account")
	}

	if _, err := token.Create(tx, user.ID, database.TokenTypeEmailPreference); err != nil {
		tx.Rollback()
		return database.User{}, errors.Wrap(err, "creating email verificaiton token")
	}
	if err := createEmailPreference(user, tx); err != nil {
		tx.Rollback()
		return database.User{}, errors.Wrap(err, "creating email preference")
	}
	if err := a.TouchLastLoginAt(user, tx); err != nil {
		tx.Rollback()
		return database.User{}, errors.Wrap(err, "updating last login")
	}

	tx.Commit()

	return user, nil
}
