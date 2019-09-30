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
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Session represents user session
type Session struct {
	UUID          string `json:"uuid"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Pro           bool   `json:"pro"`
	Classic       bool   `json:"classic"`
}

func makeSession(user database.User, account database.Account) Session {
	classic := account.AuthKeyHash != ""

	return Session{
		UUID:          user.UUID,
		Pro:           user.Cloud,
		Email:         account.Email.String,
		EmailVerified: account.EmailVerified,
		Classic:       classic,
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

type createResetTokenPayload struct {
	Email string `json:"email"`
}

func (a *App) createResetToken(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	var params createResetTokenPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var account database.Account
	conn := db.Where("email = ?", params.Email).First(&account)
	if conn.RecordNotFound() {
		return
	}
	if err := conn.Error; err != nil {
		handleError(w, errors.Wrap(err, "finding account").Error(), nil, http.StatusInternalServerError)
		return
	}

	if account.AuthKeyHash != "" {
		http.Error(w, "Please migrate your account from Dnote classic before resetting password", http.StatusBadRequest)
		return
	}

	resetToken, err := generateResetToken()
	if err != nil {
		handleError(w, errors.Wrap(err, "generating token").Error(), nil, http.StatusInternalServerError)
		return
	}

	token := database.Token{
		UserID: account.UserID,
		Value:  resetToken,
		Type:   database.TokenTypeResetPassword,
	}

	if err := db.Save(&token).Error; err != nil {
		handleError(w, errors.Wrap(err, "saving token").Error(), nil, http.StatusInternalServerError)
		return
	}

	subject := "Reset your password"
	data := struct {
		Subject string
		Token   string
	}{
		subject,
		resetToken,
	}
	email := mailer.NewEmail("noreply@getdnote.com", []string{params.Email}, subject)
	if err := email.ParseTemplate(mailer.EmailTypeResetPassword, data); err != nil {
		handleError(w, errors.Wrap(err, "parsing template").Error(), nil, http.StatusInternalServerError)
		return
	}

	if err := email.Send(); err != nil {
		handleError(w, errors.Wrap(err, "sending email").Error(), nil, http.StatusInternalServerError)
		return
	}
}

type resetPasswordPayload struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

func (a *App) resetPassword(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	var params resetPasswordPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var token database.Token
	conn := db.Where("value = ? AND type =? AND used_at IS NULL", params.Token, database.TokenTypeResetPassword).First(&token)
	if conn.RecordNotFound() {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}
	if err := conn.Error; err != nil {
		handleError(w, errors.Wrap(err, "finding token").Error(), nil, http.StatusInternalServerError)
		return
	}

	if token.UsedAt != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	// Expire after 10 minutes
	if time.Since(token.CreatedAt).Minutes() > 10 {
		http.Error(w, "This link has been expired. Please request a new password reset link.", http.StatusGone)
		return
	}

	tx := db.Begin()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		handleError(w, errors.Wrap(err, "hashing password").Error(), nil, http.StatusInternalServerError)
		return
	}

	var account database.Account
	if err := db.Where("user_id = ?", token.UserID).First(&account).Error; err != nil {
		tx.Rollback()
		handleError(w, errors.Wrap(err, "finding user").Error(), nil, http.StatusInternalServerError)
		return
	}

	if err := tx.Model(&account).Update("password", string(hashedPassword)).Error; err != nil {
		tx.Rollback()
		handleError(w, errors.Wrap(err, "updating password").Error(), nil, http.StatusInternalServerError)
		return
	}
	if err := tx.Model(&token).Update("used_at", time.Now()).Error; err != nil {
		tx.Rollback()
		handleError(w, errors.Wrap(err, "updating password reset token").Error(), nil, http.StatusInternalServerError)
		return
	}

	tx.Commit()

	var user database.User
	if err := db.Where("id = ?", account.UserID).First(&user).Error; err != nil {
		handleError(w, errors.Wrap(err, "finding user").Error(), nil, http.StatusInternalServerError)
		return
	}

	respondWithSession(w, user.ID, http.StatusOK)
}
