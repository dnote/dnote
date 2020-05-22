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
	"net/http"
	"strings"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func paginate(conn *gorm.DB, page int) *gorm.DB {
	limit := 30

	// Paginate
	if page > 0 {
		offset := limit * (page - 1)
		conn = conn.Offset(offset)
	}

	conn = conn.Limit(limit)

	return conn
}

func getBookIDs(books []database.Book) []int {
	ret := []int{}

	for _, book := range books {
		ret = append(ret, book.ID)
	}

	return ret
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password should be longer than 8 characters")
	}

	return nil
}

func getClientType(r *http.Request) string {
	origin := r.Header.Get("Origin")

	if strings.HasPrefix(origin, "moz-extension://") {
		return "firefox-extension"
	}

	if strings.HasPrefix(origin, "chrome-extension://") {
		return "chrome-extension"
	}

	userAgent := r.Header.Get("User-Agent")
	if strings.HasPrefix(userAgent, "Go-http-client") {
		return "cli"
	}

	return "web"
}

// notSupported is the handler for the route that is no longer supported
func (a *API) notSupported(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API version is not supported. Please upgrade your client.", http.StatusGone)
	return
}
