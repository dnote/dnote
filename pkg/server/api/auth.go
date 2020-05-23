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

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/handlers"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/session"
	"github.com/dnote/dnote/pkg/server/token"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// GetMeResponse is the response for getMe endpoint
type GetMeResponse struct {
	User session.Session `json:"user"`
}

func (a *API) getMe(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handlers.DoError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var account database.Account
	if err := a.App.DB.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		handlers.DoError(w, "finding account", err, http.StatusInternalServerError)
		return
	}

	tx := a.App.DB.Begin()
	if err := a.App.TouchLastLoginAt(user, tx); err != nil {
		tx.Rollback()
		// In case of an error, gracefully continue to avoid disturbing the service
		log.ErrorWrap(err, "error touching last_login_at")
	}
	tx.Commit()

	response := GetMeResponse{
		User: session.New(user, account),
	}
	handlers.RespondJSON(w, http.StatusOK, response)
}

type createResetTokenPayload struct {
	Email string `json:"email"`
}

func (a *API) createResetToken(w http.ResponseWriter, r *http.Request) {
	var params createResetTokenPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var account database.Account
	conn := a.App.DB.Where("email = ?", params.Email).First(&account)
	if conn.RecordNotFound() {
		return
	}
	if err := conn.Error; err != nil {
		handlers.DoError(w, errors.Wrap(err, "finding account").Error(), nil, http.StatusInternalServerError)
		return
	}

	resetToken, err := token.Create(a.App.DB, account.UserID, database.TokenTypeResetPassword)
	if err != nil {
		handlers.DoError(w, errors.Wrap(err, "generating token").Error(), nil, http.StatusInternalServerError)
		return
	}

	if err := a.App.SendPasswordResetEmail(account.Email.String, resetToken.Value); err != nil {
		if errors.Cause(err) == mailer.ErrSMTPNotConfigured {
			handlers.RespondInvalidSMTPConfig(w)
		} else {
			handlers.DoError(w, errors.Wrap(err, "sending password reset email").Error(), nil, http.StatusInternalServerError)
		}

		return
	}
}

type resetPasswordPayload struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

func (a *API) resetPassword(w http.ResponseWriter, r *http.Request) {
	var params resetPasswordPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var token database.Token
	conn := a.App.DB.Where("value = ? AND type =? AND used_at IS NULL", params.Token, database.TokenTypeResetPassword).First(&token)
	if conn.RecordNotFound() {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}
	if err := conn.Error; err != nil {
		handlers.DoError(w, errors.Wrap(err, "finding token").Error(), nil, http.StatusInternalServerError)
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

	tx := a.App.DB.Begin()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		handlers.DoError(w, errors.Wrap(err, "hashing password").Error(), nil, http.StatusInternalServerError)
		return
	}

	var account database.Account
	if err := a.App.DB.Where("user_id = ?", token.UserID).First(&account).Error; err != nil {
		tx.Rollback()
		handlers.DoError(w, errors.Wrap(err, "finding user").Error(), nil, http.StatusInternalServerError)
		return
	}

	if err := tx.Model(&account).Update("password", string(hashedPassword)).Error; err != nil {
		tx.Rollback()
		handlers.DoError(w, errors.Wrap(err, "updating password").Error(), nil, http.StatusInternalServerError)
		return
	}
	if err := tx.Model(&token).Update("used_at", time.Now()).Error; err != nil {
		tx.Rollback()
		handlers.DoError(w, errors.Wrap(err, "updating password reset token").Error(), nil, http.StatusInternalServerError)
		return
	}

	if err := a.App.DeleteUserSessions(tx, account.UserID); err != nil {
		tx.Rollback()
		handlers.DoError(w, errors.Wrap(err, "deleting user sessions").Error(), nil, http.StatusInternalServerError)
		return
	}

	tx.Commit()

	var user database.User
	if err := a.App.DB.Where("id = ?", account.UserID).First(&user).Error; err != nil {
		handlers.DoError(w, errors.Wrap(err, "finding user").Error(), nil, http.StatusInternalServerError)
		return
	}

	a.respondWithSession(a.App.DB, w, user.ID, http.StatusOK)

	if err := a.App.SendPasswordResetAlertEmail(account.Email.String); err != nil {
		log.ErrorWrap(err, "sending password reset email")
	}
}
