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
	"errors"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	pkgErrors "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// validatePassword validates a password
func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

// TouchLastLoginAt updates the last login timestamp
func (a *App) TouchLastLoginAt(user database.User, tx *gorm.DB) error {
	t := a.Clock.Now()
	if err := tx.Model(&user).Update("last_login_at", &t).Error; err != nil {
		return pkgErrors.Wrap(err, "updating last_login_at")
	}

	return nil
}

// CreateUser creates a user
func (a *App) CreateUser(email, password string, passwordConfirmation string) (database.User, error) {
	if email == "" {
		return database.User{}, ErrEmailRequired
	}

	if err := validatePassword(password); err != nil {
		return database.User{}, err
	}

	if password != passwordConfirmation {
		return database.User{}, ErrPasswordConfirmationMismatch
	}

	tx := a.DB.Begin()

	var count int64
	if err := tx.Model(&database.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		tx.Rollback()
		return database.User{}, pkgErrors.Wrap(err, "counting user")
	}
	if count > 0 {
		tx.Rollback()
		return database.User{}, ErrDuplicateEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return database.User{}, pkgErrors.Wrap(err, "hashing password")
	}

	uuid, err := helpers.GenUUID()
	if err != nil {
		tx.Rollback()
		return database.User{}, pkgErrors.Wrap(err, "generating UUID")
	}

	user := database.User{
		UUID:     uuid,
		Email:    database.ToNullString(email),
		Password: database.ToNullString(string(hashedPassword)),
	}
	if err = tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return database.User{}, pkgErrors.Wrap(err, "saving user")
	}

	if err := a.TouchLastLoginAt(user, tx); err != nil {
		tx.Rollback()
		return database.User{}, pkgErrors.Wrap(err, "updating last login")
	}

	tx.Commit()

	return user, nil
}

// GetUserByEmail finds a user by email
func (a *App) GetUserByEmail(email string) (*database.User, error) {
	var user database.User
	err := a.DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (a *App) GetAllUsers() ([]database.User, error) {
	var users []database.User
	err := a.DB.Find(&users).Error
	if err != nil {
		return nil, pkgErrors.Wrap(err, "finding users")
	}

	return users, nil
}

// Authenticate authenticates a user
func (a *App) Authenticate(email, password string) (*database.User, error) {
	user, err := a.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
	if err != nil {
		return nil, ErrLoginInvalid
	}

	return user, nil
}

// UpdateUserPassword updates a user's password with validation
func UpdateUserPassword(db *gorm.DB, user *database.User, newPassword string) error {
	// Validate password
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return pkgErrors.Wrap(err, "hashing password")
	}

	// Update the password
	if err := db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return pkgErrors.Wrap(err, "updating password")
	}

	return nil
}

// RemoveUser removes a user from the system
// Returns an error if the user has any notes or books
func (a *App) RemoveUser(email string) error {
	// Find the user
	user, err := a.GetUserByEmail(email)
	if err != nil {
		return err
	}

	// Check if user has any notes
	var noteCount int64
	if err := a.DB.Model(&database.Note{}).Where("user_id = ? AND deleted = ?", user.ID, false).Count(&noteCount).Error; err != nil {
		return pkgErrors.Wrap(err, "counting notes")
	}
	if noteCount > 0 {
		return ErrUserHasExistingResources
	}

	// Check if user has any books
	var bookCount int64
	if err := a.DB.Model(&database.Book{}).Where("user_id = ? AND deleted = ?", user.ID, false).Count(&bookCount).Error; err != nil {
		return pkgErrors.Wrap(err, "counting books")
	}
	if bookCount > 0 {
		return ErrUserHasExistingResources
	}

	// Delete user
	if err := a.DB.Delete(&user).Error; err != nil {
		return pkgErrors.Wrap(err, "deleting user")
	}

	return nil
}

// SignIn signs in a user
func (a *App) SignIn(user *database.User) (*database.Session, error) {
	err := a.TouchLastLoginAt(*user, a.DB)
	if err != nil {
		log.ErrorWrap(err, "touching login timestamp")
	}

	session, err := a.CreateSession(user.ID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "creating session")
	}

	return &session, nil
}
