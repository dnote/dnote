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

package session

import (
	"github.com/dnote/dnote/pkg/server/database"
)

// Session represents user session
type Session struct {
	UUID          string `json:"uuid"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Pro           bool   `json:"pro"`
}

// New returns a new session for the given user
func New(user database.User, account database.Account) Session {
	return Session{
		UUID:          user.UUID,
		Pro:           user.Cloud,
		Email:         account.Email.String,
		EmailVerified: account.EmailVerified,
	}
}
