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
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
)

// ErrInvalidAuthHeader is an error for invalid format of Authorization HTTP header
var ErrInvalidAuthHeader = errors.New("Invalid authorization header")

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
		return authHeader{}, ErrInvalidAuthHeader
	}

	parsed := authHeader{
		scheme:     parts[0],
		credential: parts[1],
	}

	return parsed, nil
}

func respondUnauthorized(w http.ResponseWriter) {
	unsetSessionCookie(w)
	w.Header().Add("WWW-Authenticate", `Bearer realm="Dnote Pro", charset="UTF-8"`)
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func legacyAuth(next http.HandlerFunc) http.HandlerFunc {
	db := database.DBConn

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("api_key")
		if err != nil {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		apiKey := c.Value
		var user database.User
		if db.Where("api_key = ?", apiKey).First(&user).RecordNotFound() {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), helpers.KeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func authWithSession(r *http.Request) (database.User, bool, error) {
	db := database.DBConn
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

func authWithToken(r *http.Request, tokenType string) (database.User, database.Token, bool, error) {
	db := database.DBConn
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

type authMiddlewareParams struct {
	ProOnly bool
}

func auth(next http.HandlerFunc, p *authMiddlewareParams) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok, err := authWithSession(r)
		if !ok || err != nil {
			respondUnauthorized(w)
			return
		}

		if p != nil && p.ProOnly {
			if !user.Cloud {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}

		ctx := context.WithValue(r.Context(), helpers.KeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func tokenAuth(next http.HandlerFunc, tokenType string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, token, ok, err := authWithToken(r, tokenType)
		if err != nil {
			// log the error and continue
			log.ErrorWrap(err, "authenticating with token")
		}

		ctx := r.Context()

		if ok {
			ctx = context.WithValue(ctx, helpers.KeyToken, token)
		} else {
			// If token-based auth fails, fall back to session-based auth
			user, ok, err = authWithSession(r)
			if err != nil {
				// log the error and continue
				log.ErrorWrap(err, "authenticating with session")
			}

			if !ok {
				respondUnauthorized(w)
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
			"remoteAddr": r.RemoteAddr,
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

// App is an application configuration
type App struct {
	Clock            clock.Clock
	StripeAPIBackend stripe.Backend
}

// init sets up the application based on the configuration
func (a *App) init() {
	stripe.Key = os.Getenv("StripeSecretKey")

	if a.StripeAPIBackend != nil {
		stripe.SetBackend(stripe.APIBackend, a.StripeAPIBackend)
	}
}

// NewRouter creates and returns a new router
func NewRouter(app *App) *mux.Router {
	app.init()

	proOnly := authMiddlewareParams{ProOnly: true}

	var routes = []Route{
		// internal
		{"GET", "/health", app.checkHealth, false},
		{"GET", "/me", auth(app.getMe, nil), true},
		{"POST", "/verification-token", auth(app.createVerificationToken, nil), true},
		{"PATCH", "/verify-email", app.verifyEmail, true},
		{"POST", "/reset-token", app.createResetToken, true},
		{"PATCH", "/reset-password", app.resetPassword, true},
		{"PATCH", "/account/profile", auth(app.updateProfile, nil), true},
		{"PATCH", "/account/password", auth(app.updatePassword, nil), true},
		{"GET", "/account/email-preference", tokenAuth(app.getEmailPreference, database.TokenTypeEmailPreference), true},
		{"PATCH", "/account/email-preference", tokenAuth(app.updateEmailPreference, database.TokenTypeEmailPreference), true},
		{"POST", "/subscriptions", auth(app.createSub, nil), true},
		{"PATCH", "/subscriptions", auth(app.updateSub, nil), true},
		{"POST", "/webhooks/stripe", app.stripeWebhook, true},
		{"GET", "/subscriptions", auth(app.getSub, nil), true},
		{"GET", "/stripe_source", auth(app.getStripeSource, nil), true},
		{"PATCH", "/stripe_source", auth(app.updateStripeSource, nil), true},
		{"GET", "/notes", auth(app.getNotes, &proOnly), false},
		{"GET", "/notes/{noteUUID}", auth(app.getNote, &proOnly), true},
		{"GET", "/calendar", auth(app.getCalendar, &proOnly), true},

		// migration of classic users
		{"GET", "/classic/presignin", cors(app.classicPresignin), true},
		{"POST", "/classic/signin", cors(app.classicSignin), true},
		{"PATCH", "/classic/migrate", auth(app.classicMigrate, &proOnly), true},
		{"GET", "/classic/notes", auth(app.classicGetNotes, nil), true},
		{"PATCH", "/classic/set-password", auth(app.classicSetPassword, nil), true},

		// v3
		{"GET", "/v3/sync/fragment", cors(auth(app.GetSyncFragment, &proOnly)), true},
		{"GET", "/v3/sync/state", cors(auth(app.GetSyncState, &proOnly)), true},
		{"OPTIONS", "/v3/books", cors(app.BooksOptions), true},
		{"GET", "/v3/books", cors(auth(app.GetBooks, &proOnly)), true},
		{"GET", "/v3/books/{bookUUID}", cors(auth(app.GetBook, &proOnly)), true},
		{"POST", "/v3/books", cors(auth(app.CreateBook, &proOnly)), true},
		{"PATCH", "/v3/books/{bookUUID}", cors(auth(app.UpdateBook, &proOnly)), false},
		{"DELETE", "/v3/books/{bookUUID}", cors(auth(app.DeleteBook, &proOnly)), false},
		{"GET", "/v3/demo/books", app.GetDemoBooks, true},
		{"OPTIONS", "/v3/notes", cors(app.NotesOptions), true},
		{"POST", "/v3/notes", cors(auth(app.CreateNote, &proOnly)), true},
		{"PATCH", "/v3/notes/{noteUUID}", auth(app.UpdateNote, &proOnly), false},
		{"DELETE", "/v3/notes/{noteUUID}", auth(app.DeleteNote, &proOnly), false},
		{"POST", "/v3/signin", cors(app.signin), true},
		{"OPTIONS", "/v3/signout", cors(app.signoutOptions), true},
		{"POST", "/v3/signout", cors(app.signout), true},
		{"POST", "/v3/register", app.register, true},
	}

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/v1").Handler(applyMiddleware(app.notSupported, true))
	router.PathPrefix("/v2").Handler(applyMiddleware(app.notSupported, true))

	for _, route := range routes {
		handler := route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(applyMiddleware(handler, route.RateLimit))
	}

	return router
}
