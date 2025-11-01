/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package middleware

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

func authWithToken(db *gorm.DB, r *http.Request, tokenType string) (database.User, database.Token, bool, error) {
	var user database.User
	var token database.Token

	query := r.URL.Query()
	tokenValue := query.Get("token")
	if tokenValue == "" {
		return user, token, false, nil
	}

	err := db.Where("value = ? AND type = ?", tokenValue, tokenType).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, token, false, nil
	} else if err != nil {
		return user, token, false, pkgErrors.Wrap(err, "finding token")
	}

	if token.UsedAt != nil && time.Since(*token.UsedAt).Minutes() > 10 {
		return user, token, false, nil
	}

	if err := db.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		return user, token, false, pkgErrors.Wrap(err, "finding user")
	}

	return user, token, true, nil
}

// AuthParams is the params for the authentication middleware
type AuthParams struct {
	RedirectGuestsToLogin bool
}

// Auth is an authentication middleware
func Auth(db *gorm.DB, next http.HandlerFunc, p *AuthParams) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok, err := AuthWithSession(db, r)
		if !ok {
			if p != nil && p.RedirectGuestsToLogin {

				q := url.Values{}
				q.Set("referrer", r.URL.Path)
				path := helpers.GetPath("/login", &q)

				http.Redirect(w, r, path, http.StatusFound)
				return
			}

			RespondUnauthorized(w)
			return
		}
		if err != nil {
			DoError(w, "authenticating with session", err, http.StatusInternalServerError)
			return
		}

		ctx := context.WithUser(r.Context(), &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TokenAuth is an authentication middleware with token
func TokenAuth(db *gorm.DB, next http.HandlerFunc, tokenType string, p *AuthParams) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, token, ok, err := authWithToken(db, r, tokenType)
		if err != nil {
			// log the error and continue
			log.ErrorWrap(err, "authenticating with token")
		}

		ctx := r.Context()

		if ok {
			ctx = context.WithToken(ctx, &token)
		} else {
			// If token-based auth fails, fall back to session-based auth
			user, ok, err = AuthWithSession(db, r)
			if err != nil {
				DoError(w, "authenticating with session", err, http.StatusInternalServerError)
				return
			}

			if !ok {
				RespondUnauthorized(w)
				return
			}
		}

		ctx = context.WithUser(ctx, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthWithSession performs user authentication with session
func AuthWithSession(db *gorm.DB, r *http.Request) (database.User, bool, error) {
	var user database.User

	sessionKey, err := GetCredential(r)
	if err != nil {
		return user, false, pkgErrors.Wrap(err, "getting credential")
	}
	if sessionKey == "" {
		return user, false, nil
	}

	var session database.Session
	err = db.Where("key = ?", sessionKey).First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, false, nil
	} else if err != nil {
		return user, false, pkgErrors.Wrap(err, "finding session")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return user, false, nil
	}

	err = db.Where("id = ?", session.UserID).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, false, nil
	} else if err != nil {
		return user, false, pkgErrors.Wrap(err, "finding user from token")
	}

	return user, true, nil
}

func GuestOnly(db *gorm.DB, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok, err := AuthWithSession(db, r)
		if err != nil {
			// log the error and continue
			log.ErrorWrap(err, "authenticating with session")
		}

		if ok {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
