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

	"github.com/dnote/dnote/pkg/server/api/operations"
	"github.com/dnote/dnote/pkg/server/database"
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

func unsetSessionCookie(w http.ResponseWriter) {
	expire := time.Now().Add(time.Hour * -24 * 30)
	cookie := http.Cookie{
		Name:     "id",
		Value:    "",
		Expires:  expire,
		Path:     "/",
		HttpOnly: true,
	}

	w.Header().Set("Cache-Control", "no-cache")
	http.SetCookie(w, &cookie)
}

func touchLastLoginAt(user database.User) error {
	db := database.DBConn

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

func (a *App) signin(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	var params signinPayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		handleError(w, "decoding payload", err, http.StatusInternalServerError)
		return
	}
	if params.Email == "" || params.Password == "" {
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	}

	var account database.Account
	conn := db.Where("email = ?", params.Email).First(&account)
	if conn.RecordNotFound() {
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	} else if conn.Error != nil {
		handleError(w, "getting user", err, http.StatusInternalServerError)
		return
	}

	password := []byte(params.Password)
	err = bcrypt.CompareHashAndPassword([]byte(account.Password.String), password)
	if err != nil {
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	}

	var user database.User
	err = db.Where("id = ?", account.UserID).First(&user).Error
	if err != nil {
		handleError(w, "finding user", err, http.StatusInternalServerError)
		return
	}

	err = operations.TouchLastLoginAt(user, db)
	if err != nil {
		http.Error(w, errors.Wrap(err, "touching login timestamp").Error(), http.StatusInternalServerError)
		return
	}

	respondWithSession(w, account.UserID, http.StatusOK)
}

func (a *App) signoutOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Version")
}

func (a *App) signout(w http.ResponseWriter, r *http.Request) {
	key, err := getCredential(r)
	if err != nil {
		handleError(w, "getting credential", nil, http.StatusInternalServerError)
		return
	}

	if key == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = operations.DeleteSession(database.DBConn, key)
	if err != nil {
		handleError(w, "deleting session", nil, http.StatusInternalServerError)
		return
	}

	unsetSessionCookie(w)
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

func parseRegisterPaylaod(r *http.Request) (registerPayload, bool) {
	var ret registerPayload
	if err := json.NewDecoder(r.Body).Decode(&ret); err != nil {
		return ret, false
	}
	if err := validateRegisterPayload(ret); err != nil {
		return ret, false
	}

	return ret, true
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	params, ok := parseRegisterPaylaod(r)
	if !ok {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var count int
	if err := db.Model(database.Account{}).Where("email = ?", params.Email).Count(&count).Error; err != nil {
		handleError(w, "checking duplicate user", err, http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Duplicate email", http.StatusBadRequest)
		return
	}

	user, err := operations.CreateUser(params.Email, params.Password)
	if err != nil {
		handleError(w, "creating user", err, http.StatusInternalServerError)
		return
	}

	respondWithSession(w, user.ID, http.StatusCreated)
}

// respondWithSession makes a HTTP response with the session from the user with the given userID.
// It sets the HTTP-Only cookie for browser clients and also sends a JSON response for non-browser clients.
func respondWithSession(w http.ResponseWriter, userID int, statusCode int) {
	db := database.DBConn

	session, err := operations.CreateSession(db, userID)
	if err != nil {
		handleError(w, "creating session", nil, http.StatusBadRequest)
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
		handleError(w, "encoding response", err, http.StatusInternalServerError)
		return
	}
}
