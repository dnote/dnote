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

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
)

// Route represents a single route
type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	RateLimit   bool
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

// getCredential extracts a session key from the request from the request header. Concretely,
// it first looks at the 'Cookie' and then the 'Authorization' header. If no credential is found,
// it returns an empty string.
func getCredential(r *http.Request) (string, error) {
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

// AuthWithSession performs user authentication with session
func AuthWithSession(db *gorm.DB, r *http.Request, p *AuthMiddlewareParams) (database.User, bool, error) {
	var user database.User

	sessionKey, err := getCredential(r)
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

func authWithToken(db *gorm.DB, r *http.Request, tokenType string, p *AuthMiddlewareParams) (database.User, database.Token, bool, error) {
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

// AuthMiddlewareParams is the params for the authentication middleware
type AuthMiddlewareParams struct {
	ProOnly bool
}

func (a *API) auth(next http.HandlerFunc, p *AuthMiddlewareParams) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok, err := AuthWithSession(a.App.DB, r, p)
		if !ok {
			respondUnauthorized(w)
			return
		}
		if err != nil {
			HandleError(w, "authenticating with session", err, http.StatusInternalServerError)
			return
		}

		if p != nil && p.ProOnly {
			if !user.Cloud {
				respondForbidden(w)
				return
			}
		}

		ctx := context.WithValue(r.Context(), helpers.KeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) tokenAuth(next http.HandlerFunc, tokenType string, p *AuthMiddlewareParams) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, token, ok, err := authWithToken(a.App.DB, r, tokenType, p)
		if err != nil {
			// log the error and continue
			log.ErrorWrap(err, "authenticating with token")
		}

		ctx := r.Context()

		if ok {
			ctx = context.WithValue(ctx, helpers.KeyToken, token)
		} else {
			// If token-based auth fails, fall back to session-based auth
			user, ok, err = AuthWithSession(a.App.DB, r, p)
			if err != nil {
				HandleError(w, "authenticating with session", err, http.StatusInternalServerError)
				return
			}

			if !ok {
				respondUnauthorized(w)
				return
			}
		}

		if p != nil && p.ProOnly {
			if !user.Cloud {
				respondForbidden(w)
				return
			}
		}

		ctx = context.WithValue(ctx, helpers.KeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow browser extensions
		if strings.HasPrefix(origin, "moz-extension://") || strings.HasPrefix(origin, "chrome-extension://") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		next.ServeHTTP(w, r)
	})
}

// logResponseWriter wraps http.ResponseWriter to expose HTTP status code for logging.
// The optional interfaces of http.ResponseWriter are lost because of the wrapping, and
// such interfaces should be implemented if needed. (i.e. http.Pusher, http.Flusher, etc.)
type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func logging(inner http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := logResponseWriter{w, http.StatusOK}
		inner.ServeHTTP(&lw, r)

		log.WithFields(log.Fields{
			"origin":     r.Header.Get("Origin"),
			"remoteAddr": lookupIP(r),
			"uri":        r.RequestURI,
			"statusCode": lw.statusCode,
			"method":     r.Method,
			"duration":   fmt.Sprintf("%dms", time.Since(start)/1000000),
			"userAgent":  r.Header.Get("User-Agent"),
		}).Info("incoming request")
	}
}

func applyMiddleware(h http.HandlerFunc, rateLimit bool) http.Handler {
	ret := h
	ret = logging(ret)

	if rateLimit && os.Getenv("GO_ENV") != "TEST" {
		ret = limit(ret)
	}

	return ret
}

// API is a web API configuration
type API struct {
	App *app.App
}

// init sets up the application based on the configuration
func (a *API) init() error {
	if err := a.App.Validate(); err != nil {
		return errors.Wrap(err, "validating the app parameters")
	}

	stripe.Key = os.Getenv("StripeSecretKey")

	if a.App.StripeAPIBackend != nil {
		stripe.SetBackend(stripe.APIBackend, a.App.StripeAPIBackend)
	}

	return nil
}

// NewRouter creates and returns a new router
func (a *API) NewRouter() (*mux.Router, error) {
	if err := a.init(); err != nil {
		return nil, errors.Wrap(err, "initializing app")
	}

	proOnly := AuthMiddlewareParams{ProOnly: true}

	var routes = []Route{
		// internal
		{"GET", "/health", a.checkHealth, false},
		{"GET", "/me", a.auth(a.getMe, nil), true},
		{"POST", "/verification-token", a.auth(a.createVerificationToken, nil), true},
		{"PATCH", "/verify-email", a.verifyEmail, true},
		{"POST", "/reset-token", a.createResetToken, true},
		{"PATCH", "/reset-password", a.resetPassword, true},
		{"PATCH", "/account/profile", a.auth(a.updateProfile, nil), true},
		{"PATCH", "/account/password", a.auth(a.updatePassword, nil), true},
		{"GET", "/account/email-preference", a.tokenAuth(a.getEmailPreference, database.TokenTypeEmailPreference, nil), true},
		{"PATCH", "/account/email-preference", a.tokenAuth(a.updateEmailPreference, database.TokenTypeEmailPreference, nil), true},
		{"POST", "/subscriptions", a.auth(a.createSub, nil), true},
		{"PATCH", "/subscriptions", a.auth(a.updateSub, nil), true},
		{"POST", "/webhooks/stripe", a.stripeWebhook, true},
		{"GET", "/subscriptions", a.auth(a.getSub, nil), true},
		{"GET", "/stripe_source", a.auth(a.getStripeSource, nil), true},
		{"PATCH", "/stripe_source", a.auth(a.updateStripeSource, nil), true},
		{"GET", "/notes", a.auth(a.getNotes, nil), false},
		{"GET", "/notes/{noteUUID}", a.getNote, true},
		{"GET", "/calendar", a.auth(a.getCalendar, nil), true},
		{"GET", "/repetition_rules", a.auth(a.getRepetitionRules, nil), true},
		{"GET", "/repetition_rules/{repetitionRuleUUID}", a.tokenAuth(a.getRepetitionRule, database.TokenTypeRepetition, &proOnly), true},
		{"POST", "/repetition_rules", a.auth(a.createRepetitionRule, &proOnly), true},
		{"PATCH", "/repetition_rules/{repetitionRuleUUID}", a.tokenAuth(a.updateRepetitionRule, database.TokenTypeRepetition, &proOnly), true},
		{"DELETE", "/repetition_rules/{repetitionRuleUUID}", a.auth(a.deleteRepetitionRule, &proOnly), true},
		{"GET", "/digests/{digestUUID}", a.auth(a.getDigest, nil), true},
		{"GET", "/digests", a.auth(a.getDigests, nil), true},
		{"POST", "/note_review", a.auth(a.createNoteReview, nil), true},
		{"DELETE", "/note_review", a.auth(a.deleteNoteReview, nil), true},

		// migration of classic users
		{"GET", "/classic/presignin", cors(a.classicPresignin), true},
		{"POST", "/classic/signin", cors(a.classicSignin), true},
		{"PATCH", "/classic/migrate", a.auth(a.classicMigrate, &proOnly), true},
		{"GET", "/classic/notes", a.auth(a.classicGetNotes, nil), true},
		{"PATCH", "/classic/set-password", a.auth(a.classicSetPassword, nil), true},

		// v3
		{"GET", "/v3/sync/fragment", cors(a.auth(a.GetSyncFragment, nil)), false},
		{"GET", "/v3/sync/state", cors(a.auth(a.GetSyncState, nil)), false},
		{"OPTIONS", "/v3/books", cors(a.BooksOptions), true},
		{"GET", "/v3/books", cors(a.auth(a.GetBooks, nil)), true},
		{"GET", "/v3/books/{bookUUID}", cors(a.auth(a.GetBook, nil)), true},
		{"POST", "/v3/books", cors(a.auth(a.CreateBook, nil)), false},
		{"PATCH", "/v3/books/{bookUUID}", cors(a.auth(a.UpdateBook, nil)), false},
		{"DELETE", "/v3/books/{bookUUID}", cors(a.auth(a.DeleteBook, nil)), false},
		{"OPTIONS", "/v3/notes", cors(a.NotesOptions), true},
		{"POST", "/v3/notes", cors(a.auth(a.CreateNote, nil)), false},
		{"PATCH", "/v3/notes/{noteUUID}", a.auth(a.UpdateNote, nil), false},
		{"DELETE", "/v3/notes/{noteUUID}", a.auth(a.DeleteNote, nil), false},
		{"POST", "/v3/signin", cors(a.signin), true},
		{"OPTIONS", "/v3/signout", cors(a.signoutOptions), true},
		{"POST", "/v3/signout", cors(a.signout), true},
		{"POST", "/v3/register", a.register, true},
	}

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/v1").Handler(applyMiddleware(a.notSupported, true))
	router.PathPrefix("/v2").Handler(applyMiddleware(a.notSupported, true))

	for _, route := range routes {
		handler := route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(applyMiddleware(handler, route.RateLimit))
	}

	return router, nil
}
