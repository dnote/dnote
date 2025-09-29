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

package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/log"
	"github.com/pkg/errors"
)

// Route represents a single route
type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	RateLimit   bool
}

// RespondForbidden responds with forbidden
func RespondForbidden(w http.ResponseWriter) {
	http.Error(w, "forbidden", http.StatusForbidden)
}

// RespondUnauthorized responds with unauthorized
func RespondUnauthorized(w http.ResponseWriter) {
	UnsetSessionCookie(w)
	w.Header().Add("WWW-Authenticate", `Bearer realm="Dnote", charset="UTF-8"`)
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

// RespondNotFound responds with not found
func RespondNotFound(w http.ResponseWriter) {
	http.Error(w, "not found", http.StatusNotFound)
}

// RespondInvalidSMTPConfig responds with invalid SMTP config error
func RespondInvalidSMTPConfig(w http.ResponseWriter) {
	http.Error(w, "SMTP is not configured", http.StatusInternalServerError)
}

// UnsetSessionCookie unsets the session cookie
func UnsetSessionCookie(w http.ResponseWriter) {
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

// DoError logs the error and responds with the given status code with a generic status text
func DoError(w http.ResponseWriter, msg string, err error, statusCode int) {
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

// NotSupported is the handler for the route that is no longer supported
func NotSupported(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API version is not supported. Please upgrade your client.", http.StatusGone)
	return
}

// getSessionKeyFromCookie reads and returns a session key from the cookie sent by the
// request. If no session key is found, it returns an empty string
func getSessionKeyFromCookie(r *http.Request) (string, error) {
	c, err := r.Cookie("id")

	if err == http.ErrNoCookie {
		return "", nil
	} else if err != nil {
		return "", errors.Wrap(err, "reading cookie")
	}

	return c.Value, nil
}

type authHeader struct {
	scheme     string
	credential string
}

func parseAuthHeader(h string) (authHeader, error) {
	parts := strings.Split(h, " ")

	if len(parts) != 2 {
		return authHeader{}, errors.New("Invalid authorization header")
	}

	parsed := authHeader{
		scheme:     parts[0],
		credential: parts[1],
	}

	return parsed, nil
}

// getSessionKeyFromAuth reads and returns a session key from the Authorization header
func getSessionKeyFromAuth(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", nil
	}

	payload, err := parseAuthHeader(h)
	if err != nil {
		return "", errors.Wrap(err, "parsing the authorization header")
	}
	if payload.scheme != "Bearer" {
		return "", errors.New("unsupported scheme")
	}

	return payload.credential, nil
}

// GetCredential extracts a session key from the request from the request header. Concretely,
// it first looks at the 'Cookie' and then the 'Authorization' header. If no credential is found,
// it returns an empty string.
func GetCredential(r *http.Request) (string, error) {
	ret, err := getSessionKeyFromCookie(r)
	if err != nil {
		return "", errors.Wrap(err, "getting session key from cookie")
	}
	if ret != "" {
		return ret, nil
	}

	ret, err = getSessionKeyFromAuth(r)
	if err != nil {
		return "", errors.Wrap(err, "getting session key from Authorization header")
	}

	return ret, nil
}
