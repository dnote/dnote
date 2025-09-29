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

package app

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
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
func (a *App) CreateUser(email, password string, passwordConfirmation string) (database.User, error) {
	if email == "" {
		return database.User{}, ErrEmailRequired
	}

	if len(password) < 8 {
		return database.User{}, ErrPasswordTooShort
	}

	if password != passwordConfirmation {
		return database.User{}, ErrPasswordConfirmationMismatch
	}

	tx := a.DB.Begin()

	var count int
	if err := tx.Model(database.Account{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return database.User{}, errors.Wrap(err, "counting user")
	}
	if count > 0 {
		return database.User{}, ErrDuplicateEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return database.User{}, errors.Wrap(err, "hashing password")
	}

	// Grant all privileges if self-hosting
	var pro bool
	if a.Config.OnPremises {
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

// Authenticate authenticates a user
func (a *App) Authenticate(email, password string) (*database.User, error) {
	var account database.Account
	conn := a.DB.Where("email = ?", email).First(&account)
	if conn.RecordNotFound() {
		return nil, ErrNotFound
	} else if conn.Error != nil {
		return nil, conn.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(account.Password.String), []byte(password))
	if err != nil {
		return nil, ErrLoginInvalid
	}

	var user database.User
	err = a.DB.Where("id = ?", account.UserID).First(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, "finding user")
	}

	return &user, nil
}

// SignIn signs in a user
func (a *App) SignIn(user *database.User) (*database.Session, error) {
	err := a.TouchLastLoginAt(*user, a.DB)
	if err != nil {
		log.ErrorWrap(err, "touching login timestamp")
	}

	session, err := a.CreateSession(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "creating session")
	}

	return &session, nil
}
