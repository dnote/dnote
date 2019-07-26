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

	"github.com/dnote/dnote/pkg/server/api/crypt"
	"github.com/dnote/dnote/pkg/server/api/operations"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/pkg/errors"
)

// ErrLoginFailure is an error for failed login
var ErrLoginFailure = errors.New("Wrong email and password combination")

// SessionResponse is a response containing a session information
type SessionResponse struct {
	Key          string `json:"key"`
	ExpiresAt    int64  `json:"expires_at"`
	CipherKeyEnc string `json:"cipher_key_enc"`
}

type signinPayload struct {
	Email   string `json:"email"`
	AuthKey string `json:"auth_key"`
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

func (a *App) signin(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	var params signinPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handleError(w, "decoding payload", err, http.StatusInternalServerError)
		return
	}

	if params.Email == "" || params.AuthKey == "" {
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	}

	var account database.Account
	conn := db.Where("email = ?", params.Email).First(&account)
	if conn.RecordNotFound() {
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	} else if err := conn.Error; err != nil {
		handleError(w, "getting user", err, http.StatusInternalServerError)
		return
	}

	authKeyHash := crypt.HashAuthKey(params.AuthKey, account.Salt, account.ServerKDFIteration)
	if account.AuthKeyHash != authKeyHash {
		log.WithFields(log.Fields{
			"account_id": account.ID,
		}).Error("Sign in password mismatch")
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	}

	respondWithSession(w, account.UserID, account.CipherKeyEnc)
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
	Email        string `json:"email"`
	AuthKey      string `json:"auth_key"`
	Iteration    int    `json:"iteration"`
	CipherKeyEnc string `json:"cipher_key_enc"`
}

func validateRegisterPayload(p registerPayload) error {
	if p.Email == "" {
		return errors.New("email is required")
	}
	if p.AuthKey == "" {
		return errors.New("auth_key is required")
	}
	if p.Iteration == 0 {
		return errors.New("iteration is required")
	}
	if p.CipherKeyEnc == "" {
		return errors.New("cipher_key_enc is required")
	}

	return nil
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	var params registerPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handleError(w, "decoding payload", err, http.StatusInternalServerError)
		return
	}
	if err := validateRegisterPayload(params); err != nil {
		handleError(w, "validating payload", err, http.StatusBadRequest)
		return
	}

	var count int
	if err := db.Model(database.Account{}).Where("email = ?", params.Email).Count(&count).Error; err != nil {
		handleError(w, "checking duplicate", err, http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Duplicate email", http.StatusBadRequest)
		return
	}

	tx := db.Begin()

	user, err := operations.CreateUser(tx, params.Email, params.AuthKey, params.CipherKeyEnc, params.Iteration)
	if err != nil {
		tx.Rollback()

		handleError(w, "creating user", nil, http.StatusBadRequest)
		return
	}

	var account database.Account
	if err := tx.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		tx.Rollback()
		handleError(w, "finding account", nil, http.StatusBadRequest)
		return
	}

	tx.Commit()

	respondWithSession(w, user.ID, account.CipherKeyEnc)
}

// respondWithSession makes a HTTP response with the session from the user with the given userID.
// It sets the HTTP-Only cookie for browser clients and also sends a JSON response for non-browser clients.
func respondWithSession(w http.ResponseWriter, userID int, cipherKeyEnc string) {
	db := database.DBConn

	session, err := operations.CreateSession(db, userID)
	if err != nil {
		handleError(w, "creating session", nil, http.StatusBadRequest)
		return
	}

	setSessionCookie(w, session.Key, session.ExpiresAt)

	response := SessionResponse{
		Key:          session.Key,
		ExpiresAt:    session.ExpiresAt.Unix(),
		CipherKeyEnc: cipherKeyEnc,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handleError(w, "encoding response", err, http.StatusInternalServerError)
		return
	}
}

// PresigninResponse is a response for presignin
type PresigninResponse struct {
	Iteration int `json:"iteration"`
}

func (a *App) presignin(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	q := r.URL.Query()
	email := q.Get("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	var account database.Account
	conn := db.Where("email = ?", email).First(&account)
	if !conn.RecordNotFound() && conn.Error != nil {
		handleError(w, "getting user", conn.Error, http.StatusInternalServerError)
		return
	}

	var response PresigninResponse
	if conn.RecordNotFound() {
		response = PresigninResponse{
			Iteration: 100000,
		}
	} else {
		response = PresigninResponse{
			Iteration: account.ClientKDFIteration,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handleError(w, "encoding response", nil, http.StatusInternalServerError)
		return
	}
}
