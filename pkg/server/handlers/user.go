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
	"net/http"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type updateProfilePayload struct {
	Email string `json:"email"`
}

// updateProfile updates user
func (a *API) updateProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var params updateProfilePayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, errors.Wrap(err, "invalid params").Error(), http.StatusBadRequest)
		return
	}

	// Validate
	if len(params.Email) > 60 {
		http.Error(w, "Email is too long", http.StatusBadRequest)
		return
	}

	var account database.Account
	err = a.App.DB.Where("user_id = ?", user.ID).First(&account).Error
	if err != nil {
		HandleError(w, "finding account", err, http.StatusInternalServerError)
		return
	}

	tx := a.App.DB.Begin()
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		HandleError(w, "saving user", err, http.StatusInternalServerError)
		return
	}

	// check if email was changed
	if params.Email != account.Email.String {
		account.EmailVerified = false
	}
	account.Email.String = params.Email

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		HandleError(w, "saving account", err, http.StatusInternalServerError)
		return
	}

	tx.Commit()

	respondWithSession(a.App.DB, w, user.ID, http.StatusOK)
}

type updateEmailPayload struct {
	NewEmail        string `json:"new_email"`
	NewCipherKeyEnc string `json:"new_cipher_key_enc"`
	OldAuthKey      string `json:"old_auth_key"`
	NewAuthKey      string `json:"new_auth_key"`
}

func respondWithCalendar(db *gorm.DB, w http.ResponseWriter, userID int) {
	rows, err := db.Table("notes").Select("COUNT(id), date(to_timestamp(added_on/1000000000)) AS added_date").
		Where("user_id = ?", userID).
		Group("added_date").
		Order("added_date DESC").Rows()

	if err != nil {
		HandleError(w, "Failed to count lessons", err, http.StatusInternalServerError)
		return
	}

	payload := map[string]int{}

	for rows.Next() {
		var count int
		var d time.Time

		if err := rows.Scan(&count, &d); err != nil {
			HandleError(w, "counting notes", err, http.StatusInternalServerError)
		}
		payload[d.Format("2006-1-2")] = count
	}

	respondJSON(w, http.StatusOK, payload)
}

func (a *API) getCalendar(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	respondWithCalendar(a.App.DB, w, user.ID)
}

func (a *API) createVerificationToken(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var account database.Account
	err := a.App.DB.Where("user_id = ?", user.ID).First(&account).Error
	if err != nil {
		HandleError(w, "finding account", err, http.StatusInternalServerError)
		return
	}

	if account.EmailVerified {
		http.Error(w, "Email already verified", http.StatusGone)
		return
	}
	if account.Email.String == "" {
		http.Error(w, "Email not set", http.StatusUnprocessableEntity)
		return
	}

	tokenValue, err := generateVerificationCode()
	if err != nil {
		HandleError(w, "generating verification code", err, http.StatusInternalServerError)
		return
	}

	token := database.Token{
		UserID: account.UserID,
		Value:  tokenValue,
		Type:   database.TokenTypeEmailVerification,
	}

	if err := a.App.DB.Save(&token).Error; err != nil {
		HandleError(w, "saving token", err, http.StatusInternalServerError)
		return
	}

	if err := a.App.SendVerificationEmail(account.Email.String, tokenValue); err != nil {
		if errors.Cause(err) == mailer.ErrSMTPNotConfigured {
			respondInvalidSMTPConfig(w)
		} else {
			HandleError(w, errors.Wrap(err, "sending verification email").Error(), nil, http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
}

type verifyEmailPayload struct {
	Token string `json:"token"`
}

func (a *API) verifyEmail(w http.ResponseWriter, r *http.Request) {
	var params verifyEmailPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		HandleError(w, "decoding payload", err, http.StatusInternalServerError)
		return
	}

	var token database.Token
	if err := a.App.DB.
		Where("value = ? AND type = ?", params.Token, database.TokenTypeEmailVerification).
		First(&token).Error; err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	if token.UsedAt != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	// Expire after ttl
	if time.Since(token.CreatedAt).Minutes() > 30 {
		http.Error(w, "This link has been expired. Please request a new link.", http.StatusGone)
		return
	}

	var account database.Account
	if err := a.App.DB.Where("user_id = ?", token.UserID).First(&account).Error; err != nil {
		HandleError(w, "finding account", err, http.StatusInternalServerError)
		return
	}
	if account.EmailVerified {
		http.Error(w, "Already verified", http.StatusConflict)
		return
	}

	tx := a.App.DB.Begin()
	account.EmailVerified = true
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		HandleError(w, "updating email_verified", err, http.StatusInternalServerError)
		return
	}
	if err := tx.Model(&token).Update("used_at", time.Now()).Error; err != nil {
		tx.Rollback()
		HandleError(w, "updating reset token", err, http.StatusInternalServerError)
		return
	}
	tx.Commit()

	var user database.User
	if err := a.App.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		HandleError(w, "finding user", err, http.StatusInternalServerError)
		return
	}

	session := makeSession(user, account)
	respondJSON(w, http.StatusOK, session)
}

type updateEmailPreferencePayload struct {
	DigestWeekly bool `json:"digest_weekly"`
}

func (a *API) updateEmailPreference(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var params updateEmailPreferencePayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		HandleError(w, "decoding payload", err, http.StatusInternalServerError)
		return
	}

	var frequency database.EmailPreference
	if err := a.App.DB.Where(database.EmailPreference{UserID: user.ID}).FirstOrCreate(&frequency).Error; err != nil {
		HandleError(w, "finding frequency", err, http.StatusInternalServerError)
		return
	}

	tx := a.App.DB.Begin()

	frequency.DigestWeekly = params.DigestWeekly
	if err := tx.Save(&frequency).Error; err != nil {
		tx.Rollback()
		HandleError(w, "saving frequency", err, http.StatusInternalServerError)
		return
	}

	token, ok := r.Context().Value(helpers.KeyToken).(database.Token)
	if ok {
		// Use token if the user was authenticated by token
		if err := tx.Model(&token).Update("used_at", time.Now()).Error; err != nil {
			tx.Rollback()
			HandleError(w, "updating reset token", err, http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	respondJSON(w, http.StatusOK, frequency)
}

func (a *API) getEmailPreference(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var pref database.EmailPreference
	if err := a.App.DB.Where(database.EmailPreference{UserID: user.ID}).First(&pref).Error; err != nil {
		HandleError(w, "finding pref", err, http.StatusInternalServerError)
		return
	}

	presented := presenters.PresentEmailPreference(pref)
	respondJSON(w, http.StatusOK, presented)
}

type updatePasswordPayload struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (a *API) updatePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var params updatePasswordPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if params.OldPassword == "" || params.NewPassword == "" {
		http.Error(w, "invalid params", http.StatusBadRequest)
		return
	}

	var account database.Account
	if err := a.App.DB.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		HandleError(w, "getting user", nil, http.StatusInternalServerError)
		return
	}

	password := []byte(params.OldPassword)
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password.String), password); err != nil {
		log.WithFields(log.Fields{
			"user_id": user.ID,
		}).Warn("invalid password update attempt")
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	if err := validatePassword(params.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, errors.Wrap(err, "hashing password").Error(), http.StatusInternalServerError)
		return
	}

	if err := a.App.DB.Model(&account).Update("password", string(hashedNewPassword)).Error; err != nil {
		http.Error(w, errors.Wrap(err, "updating password").Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
