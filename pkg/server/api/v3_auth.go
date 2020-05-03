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
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// ErrLoginFailure is an error for failed login
var ErrLoginFailure = errors.New("Wrong email and password combination")

// SessionResponse is a response containing a session information
type SessionResponse struct {
	Key       string `json:"key"`
	ExpiresAt int64  `json:"expires_at"`
}

func setSessionCookie(w http.ResponseWriter, key string, expires time.Time) {
	cookie := http.Cookie{
		Name:     "id",
		Value:    key,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func touchLastLoginAt(db *gorm.DB, user database.User) error {
	t := time.Now()
	if err := db.Model(&user).Update(database.User{LastLoginAt: &t}).Error; err != nil {
		return errors.Wrap(err, "updating last_login_at")
	}

	return nil
}

type signinPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) signin(w http.ResponseWriter, r *http.Request) {
	var params signinPayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		handlers.DoError(w, "decoding payload", err, http.StatusInternalServerError)
		return
	}
	if params.Email == "" || params.Password == "" {
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	}

	var account database.Account
	conn := a.App.DB.Where("email = ?", params.Email).First(&account)
	if conn.RecordNotFound() {
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	} else if conn.Error != nil {
		handlers.DoError(w, "getting user", err, http.StatusInternalServerError)
		return
	}

	password := []byte(params.Password)
	err = bcrypt.CompareHashAndPassword([]byte(account.Password.String), password)
	if err != nil {
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	}

	var user database.User
	err = a.App.DB.Where("id = ?", account.UserID).First(&user).Error
	if err != nil {
		handlers.DoError(w, "finding user", err, http.StatusInternalServerError)
		return
	}

	err = a.App.TouchLastLoginAt(user, a.App.DB)
	if err != nil {
		http.Error(w, errors.Wrap(err, "touching login timestamp").Error(), http.StatusInternalServerError)
		return
	}

	a.respondWithSession(a.App.DB, w, account.UserID, http.StatusOK)
}

func (a *API) signoutOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Version")
}

func (a *API) signout(w http.ResponseWriter, r *http.Request) {
	key, err := handlers.GetCredential(r)
	if err != nil {
		handlers.DoError(w, "getting credential", nil, http.StatusInternalServerError)
		return
	}

	if key == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = a.App.DeleteSession(key)
	if err != nil {
		handlers.DoError(w, "deleting session", nil, http.StatusInternalServerError)
		return
	}

	handlers.UnsetSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

type registerPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func validateRegisterPayload(p registerPayload) error {
	if p.Email == "" {
		return errors.New("email is required")
	}
	if len(p.Password) < 8 {
		return errors.New("Password should be longer than 8 characters")
	}

	return nil
}

func parseRegisterPaylaod(r *http.Request) (registerPayload, error) {
	var ret registerPayload
	if err := json.NewDecoder(r.Body).Decode(&ret); err != nil {
		return ret, errors.Wrap(err, "decoding json")
	}

	return ret, nil
}

func (a *API) register(w http.ResponseWriter, r *http.Request) {
	if a.App.Config.DisableRegistration {
		handlers.RespondForbidden(w)
		return
	}

	params, err := parseRegisterPaylaod(r)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if err := validateRegisterPayload(params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var count int
	if err := a.App.DB.Model(database.Account{}).Where("email = ?", params.Email).Count(&count).Error; err != nil {
		handlers.DoError(w, "checking duplicate user", err, http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Duplicate email", http.StatusBadRequest)
		return
	}

	user, err := a.App.CreateUser(params.Email, params.Password)
	if err != nil {
		handlers.DoError(w, "creating user", err, http.StatusInternalServerError)
		return
	}

	a.respondWithSession(a.App.DB, w, user.ID, http.StatusCreated)

	if err := a.App.SendWelcomeEmail(params.Email); err != nil {
		log.ErrorWrap(err, "sending welcome email")
	}
}

// respondWithSession makes a HTTP response with the session from the user with the given userID.
// It sets the HTTP-Only cookie for browser clients and also sends a JSON response for non-browser clients.
func (a *API) respondWithSession(db *gorm.DB, w http.ResponseWriter, userID int, statusCode int) {
	session, err := a.App.CreateSession(userID)
	if err != nil {
		handlers.DoError(w, "creating session", nil, http.StatusBadRequest)
		return
	}

	setSessionCookie(w, session.Key, session.ExpiresAt)

	response := SessionResponse{
		Key:       session.Key,
		ExpiresAt: session.ExpiresAt.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handlers.DoError(w, "encoding response", err, http.StatusInternalServerError)
		return
	}
}
