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

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/api/operations"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Session represents user session
type Session struct {
	ID              int    `json:"id"`
	GithubName      string `json:"github_name"`
	GithubAccountID string `json:"github_account_id"`
	APIKey          string `json:"api_key"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	EmailVerified   bool   `json:"email_verified"`
	Provider        string `json:"provider"`
	Cloud           bool   `json:"cloud"`
	Legacy          bool   `json:"legacy"`
	Encrypted       bool   `json:"encrypted"`
	CipherKeyEnc    string `json:"cipher_key_enc"`
}

func makeSession(user database.User, account database.Account) Session {
	legacy := account.AuthKeyHash == ""

	return Session{
		// TODO: remove ID and use UUID
		ID:              user.ID,
		GithubName:      account.Nickname,
		GithubAccountID: account.AccountID,
		APIKey:          user.APIKey,
		Cloud:           user.Cloud,
		Email:           account.Email.String,
		EmailVerified:   account.EmailVerified,
		Name:            user.Name,
		Provider:        account.Provider,
		Legacy:          legacy,
		Encrypted:       user.Encrypted,
		CipherKeyEnc:    account.CipherKeyEnc,
	}
}

func (a *App) getMe(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	db := database.DBConn

	var account database.Account
	if err := db.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		handleError(w, "finding account", err, http.StatusInternalServerError)
		return
	}

	session := makeSession(user, account)

	response := struct {
		User Session `json:"user"`
	}{
		User: session,
	}

	tx := db.Begin()
	if err := operations.TouchLastLoginAt(user, tx); err != nil {
		tx.Rollback()
		// In case of an error, gracefully continue to avoid disturbing the service
		log.Println("error touching last_login_at", err.Error())
	}
	tx.Commit()

	respondJSON(w, response)
}

// OauthCallbackHandler handler
func (a *App) oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	githubUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		handleError(w, "completing user uath", err, http.StatusInternalServerError)
		return
	}

	db := database.DBConn
	tx := db.Begin()

	currentUser, err := findUserFromOauth(githubUser, tx)
	if err != nil {
		tx.Rollback()
		handleError(w, "Failed to upsert user", err, http.StatusInternalServerError)
		return
	}
	err = operations.TouchLastLoginAt(currentUser, tx)
	if err != nil {
		tx.Rollback()
		handleError(w, "touching login timestamp", err, http.StatusInternalServerError)
		return
	}

	tx.Commit()

	setAuthCookie(w, currentUser)
	http.Redirect(w, r, "/app/legacy/register", 301)
}

// helpers
// setAuthCookie sets 'api_key' cookie in the HTTP response for a given user
func setAuthCookie(w http.ResponseWriter, currentUser database.User) {
	expire := time.Now().Add(time.Hour * 24 * 90)
	cookie := http.Cookie{
		Name:     "api_key",
		Value:    currentUser.APIKey,
		Expires:  expire,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func findUserFromOauth(oauthUser goth.User, tx *gorm.DB) (database.User, error) {
	var user database.User
	var account database.Account

	conn := tx.Where("account_id = ?", oauthUser.UserID).First(&account)
	if err := conn.Error; err != nil {
		return user, errors.Wrap(err, "finding account")
	}

	conn = tx.Where("id = ?", account.UserID).First(&user)
	if err := conn.Error; err != nil {
		return user, errors.Wrap(err, "finding user")
	}

	return user, nil
}

type legacyPasswordLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) legacyPasswordLogin(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	var params legacyPasswordLoginPayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		handleError(w, "decoding payload", err, http.StatusInternalServerError)
		return
	}

	var account database.Account
	conn := db.Where("email = ?", params.Email).First(&account)
	if conn.RecordNotFound() {
		http.Error(w, "Wrong email and password combination", http.StatusUnauthorized)
		return
	} else if conn.Error != nil {
		handleError(w, "getting user", err, http.StatusInternalServerError)
		return
	}

	password := []byte(params.Password)
	err = bcrypt.CompareHashAndPassword([]byte(account.Password.String), password)
	if err != nil {
		http.Error(w, "Wrong email and password combination", http.StatusUnauthorized)
		return
	}

	var user database.User
	err = db.Where("id = ?", account.UserID).First(&user).Error
	if err != nil {
		handleError(w, "finding user", err, http.StatusInternalServerError)
		return
	}

	tx := db.Begin()

	err = operations.TouchLastLoginAt(user, tx)
	if err != nil {
		tx.Rollback()
		handleError(w, "touching login timestamp", err, http.StatusInternalServerError)
		return
	}

	tx.Commit()

	session := makeSession(user, account)
	response := struct {
		User Session `json:"user"`
	}{
		User: session,
	}

	setAuthCookie(w, user)
	respondJSON(w, response)
}

type legacyRegisterPayload struct {
	Email        string `json:"email"`
	AuthKey      string `json:"auth_key"`
	CipherKeyEnc string `json:"cipher_key_enc"`
	Iteration    int    `json:"iteration"`
}

func validateLegacyRegisterPayload(p legacyRegisterPayload) error {
	if p.Email == "" {
		return errors.New("email is required")
	}
	if p.AuthKey == "" {
		return errors.New("auth_key is required")
	}
	if p.CipherKeyEnc == "" {
		return errors.New("cipher_key_enc is required")
	}
	if p.Iteration == 0 {
		return errors.New("iteration is required")
	}

	return nil
}

func (a *App) legacyRegister(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	db := database.DBConn

	var params legacyRegisterPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handleError(w, "decoding payload", err, http.StatusInternalServerError)
		return
	}
	if err := validateLegacyRegisterPayload(params); err != nil {
		handleError(w, "validating payload", err, http.StatusBadRequest)
		return
	}

	tx := db.Begin()

	err := operations.LegacyRegisterUser(tx, user.ID, params.Email, params.AuthKey, params.CipherKeyEnc, params.Iteration)
	if err != nil {
		tx.Rollback()
		handleError(w, "creating user", err, http.StatusBadRequest)
		return
	}

	tx.Commit()

	var account database.Account
	if err := db.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		handleError(w, "finding account", err, http.StatusInternalServerError)
		return
	}

	respondWithSession(w, user.ID, account.CipherKeyEnc)
}

func (a *App) legacyMigrate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	db := database.DBConn

	if err := db.Model(&user).Update("encrypted = ?", true).Error; err != nil {
		handleError(w, "updating user", err, http.StatusInternalServerError)
		return
	}
}
