package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func authWithToken(db *gorm.DB, r *http.Request, tokenType string, p *AuthParams) (database.User, database.Token, bool, error) {
	var user database.User
	var token database.Token

	query := r.URL.Query()
	tokenValue := query.Get("token")
	if tokenValue == "" {
		return user, token, false, nil
	}

	conn := db.Where("value = ? AND type = ?", tokenValue, tokenType).First(&token)
	if conn.RecordNotFound() {
		return user, token, false, nil
	} else if err := conn.Error; err != nil {
		return user, token, false, errors.Wrap(err, "finding token")
	}

	if token.UsedAt != nil && time.Since(*token.UsedAt).Minutes() > 10 {
		return user, token, false, nil
	}

	if err := db.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		return user, token, false, errors.Wrap(err, "finding user")
	}

	return user, token, true, nil
}

// Cors allows browser extensions to load resources
func Cors(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow browser extensions
		if strings.HasPrefix(origin, "moz-extension://") || strings.HasPrefix(origin, "chrome-extension://") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		next.ServeHTTP(w, r)
	})
}

// AuthParams is the params for the authentication middleware
type AuthParams struct {
	ProOnly               bool
	RedirectGuestsToLogin bool
}

// Auth is an authentication middleware
func Auth(a *app.App, next http.HandlerFunc, p *AuthParams) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok, err := AuthWithSession(a.DB, r, p)
		if !ok {
			if p != nil && p.RedirectGuestsToLogin {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			RespondUnauthorized(w)
			return
		}
		if err != nil {
			DoError(w, "authenticating with session", err, http.StatusInternalServerError)
			return
		}

		if p != nil && p.ProOnly {
			if !user.Cloud {
				RespondForbidden(w)
				return
			}
		}

		ctx := context.WithValue(r.Context(), helpers.KeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TokenAuth is an authentication middleware with token
func TokenAuth(a *app.App, next http.HandlerFunc, tokenType string, p *AuthParams) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, token, ok, err := authWithToken(a.DB, r, tokenType, p)
		if err != nil {
			// log the error and continue
			log.ErrorWrap(err, "authenticating with token")
		}

		ctx := r.Context()

		if ok {
			ctx = context.WithValue(ctx, helpers.KeyToken, token)
		} else {
			// If token-based auth fails, fall back to session-based auth
			user, ok, err = AuthWithSession(a.DB, r, p)
			if err != nil {
				DoError(w, "authenticating with session", err, http.StatusInternalServerError)
				return
			}

			if !ok {
				RespondUnauthorized(w)
				return
			}
		}

		if p != nil && p.ProOnly {
			if !user.Cloud {
				RespondForbidden(w)
				return
			}
		}

		ctx = context.WithValue(ctx, helpers.KeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthWithSession performs user authentication with session
func AuthWithSession(db *gorm.DB, r *http.Request, p *AuthParams) (database.User, bool, error) {
	var user database.User

	sessionKey, err := GetCredential(r)
	if err != nil {
		return user, false, errors.Wrap(err, "getting credential")
	}
	if sessionKey == "" {
		return user, false, nil
	}

	var session database.Session
	conn := db.Where("key = ?", sessionKey).First(&session)

	if conn.RecordNotFound() {
		return user, false, nil
	} else if err := conn.Error; err != nil {
		return user, false, errors.Wrap(err, "finding session")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return user, false, nil
	}

	conn = db.Where("id = ?", session.UserID).First(&user)

	if conn.RecordNotFound() {
		return user, false, nil
	} else if err := conn.Error; err != nil {
		return user, false, errors.Wrap(err, "finding user from token")
	}

	return user, true, nil
}
