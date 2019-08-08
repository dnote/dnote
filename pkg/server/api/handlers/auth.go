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
	"log"
	"net/http"

	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/api/operations"
	"github.com/dnote/dnote/pkg/server/database"
)

// Session represents user session
type Session struct {
	UUID          string `json:"uuid"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Pro           bool   `json:"pro"`
}

func makeSession(user database.User, account database.Account) Session {
	return Session{
		UUID:          user.UUID,
		Pro:           user.Cloud,
		Email:         account.Email.String,
		EmailVerified: account.EmailVerified,
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
