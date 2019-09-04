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
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	demoUserEmail = "demo@dnote.io"
)

func generateRandomToken(bits int) (string, error) {
	b := make([]byte, bits)

	_, err := crand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "generating random bytes")
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func generateResetToken() (string, error) {
	ret, err := generateRandomToken(16)
	if err != nil {
		return "", errors.Wrap(err, "generating random token")
	}

	return ret, nil
}

func generateVerificationCode() (string, error) {
	ret, err := generateRandomToken(16)
	if err != nil {
		return "", errors.Wrap(err, "generating random token")
	}

	return ret, nil
}

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

func getClientType(origin string) string {
	if strings.HasPrefix(origin, "moz-extension://") {
		return "firefox-extension"
	}

	if strings.HasPrefix(origin, "chrome-extension://") {
		return "chrome-extension"
	}

	return "web"
}

// handleError logs the error and responds with the given status code with a generic status text
func handleError(w http.ResponseWriter, msg string, err error, statusCode int) {
	var message string
	if err == nil {
		message = msg
	} else {
		message = errors.Wrap(err, msg).Error()
	}

	log.WithFields(log.Fields{
		"statusCode": statusCode,
	}).Error(message)

	statusText := http.StatusText(statusCode)
	http.Error(w, statusText, statusCode)
}

// respondJSON encodes the given payload into a JSON format and writes it to the given response writer
func respondJSON(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		handleError(w, "encoding response", err, http.StatusInternalServerError)
	}
}

// notSupported is the handler for the route that is no longer supported
func (a *App) notSupported(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API version is not supported. Please upgrade your client.", http.StatusGone)
	return
}
