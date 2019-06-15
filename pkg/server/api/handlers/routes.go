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
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/api/clock"
	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/api/logger"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
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
		logger.Err(errors.Wrap(err, "getting credential").Error())
		return user, false, err
	}

	if sessionKey == "" {
		return user, false, nil
	}

	var session database.Session
	conn := db.Where("key = ?", sessionKey).First(&session)

	if conn.RecordNotFound() {
		return user, false, nil
	} else if err := conn.Error; err != nil {
		logger.Err(errors.Wrap(err, "finding session").Error())
		return user, false, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return user, false, nil
	}

	conn = db.Where("id = ?", session.UserID).First(&user)

	if conn.RecordNotFound() {
		return user, false, nil
	} else if err := conn.Error; err != nil {
		logger.Err(errors.Wrap(err, "finding user from token").Error())
		return user, false, err
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
		logger.Err(errors.Wrap(err, "finding token").Error())
		return user, token, false, err
	}

	if token.UsedAt != nil && time.Since(*token.UsedAt).Minutes() > 10 {
		return user, token, false, nil
	}

	if err := db.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		logger.Err(errors.Wrap(err, "finding user").Error())
		return user, token, false, err
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
			logger.Err(errors.Wrap(err, "authenticating with token").Error())
		}

		ctx := r.Context()

		if ok {
			ctx = context.WithValue(ctx, helpers.KeyToken, token)
		} else {
			// If token-based auth fails, fall back to session-based auth
			user, ok, err = authWithSession(r)
			if err != nil {
				// log the error and continue
				logger.Err(errors.Wrap(err, "authenticating with session").Error())
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

func logging(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func applyMiddleware(h http.Handler, rateLimit bool) http.Handler {
	ret := h
	ret = logging(ret)

	if rateLimit && os.Getenv("GO_ENV") == "PRODUCTION" {
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
		Route{"GET", "/health", app.checkHealth, false},
		Route{"GET", "/me", auth(app.getMe, nil), true},
		Route{"POST", "/verification-token", auth(app.createVerificationToken, nil), true},
		Route{"PATCH", "/verify-email", app.verifyEmail, true},
		Route{"GET", "/auth/{provider}", gothic.BeginAuthHandler, true},
		Route{"GET", "/auth/{provider}/callback", app.oauthCallbackHandler, true},
		Route{"PATCH", "/account/profile", auth(app.updateProfile, nil), true},
		Route{"PATCH", "/account/email", auth(app.updateEmail, nil), true},
		Route{"PATCH", "/account/password", auth(app.updatePassword, nil), true},
		Route{"GET", "/account/email-preference", tokenAuth(app.getEmailPreference, database.TokenTypeEmailPreference), true},
		Route{"PATCH", "/account/email-preference", tokenAuth(app.updateEmailPreference, database.TokenTypeEmailPreference), true},
		Route{"POST", "/subscriptions", auth(app.createSub, nil), true},
		Route{"PATCH", "/subscriptions", auth(app.updateSub, nil), true},
		Route{"POST", "/webhooks/stripe", app.stripeWebhook, true},
		Route{"GET", "/subscriptions", auth(app.getSub, nil), true},
		Route{"GET", "/stripe_source", auth(app.getStripeSource, nil), true},
		Route{"GET", "/notes", auth(app.getNotes, &proOnly), false},
		Route{"GET", "/demo/notes", app.getDemoNotes, true},
		Route{"GET", "/notes/{noteUUID}", auth(app.getNote, &proOnly), true},
		Route{"GET", "/demo/notes/{noteUUID}", app.getDemoNote, true},
		Route{"GET", "/calendar", auth(app.getCalendar, &proOnly), true},
		Route{"GET", "/demo/calendar", app.getDemoCalendar, true},
		Route{"GET", "/digests/{digestUUID}", auth(app.getDigest, &proOnly), true},
		Route{"GET", "/demo/digests/{digestUUID}", app.getDemoDigest, true},
		Route{"GET", "/digests", auth(app.getDigests, &proOnly), true},
		Route{"GET", "/demo/digests", app.getDemoDigests, true},
		//Route{"GET", "/books/{bookUUID}", cors(auth(app.getBook)), true},

		// routes for user migration to use encryption
		Route{"POST", "/legacy/signin", app.legacyPasswordLogin, true},
		Route{"POST", "/legacy/register", legacyAuth(app.legacyRegister), true},
		Route{"GET", "/legacy/me", legacyAuth(app.getMe), true},
		Route{"GET", "/legacy/notes", auth(app.legacyGetNotes, &proOnly), false},
		Route{"PATCH", "/legacy/migrate", auth(app.legacyMigrate, &proOnly), false},
		Route{"GET", "/auth/{provider}", gothic.BeginAuthHandler, true},
		Route{"GET", "/auth/{provider}/callback", app.oauthCallbackHandler, true},

		// v1
		Route{"POST", "/v1/sync", cors(app.Sync), true},
		Route{"GET", "/v1/sync/fragment", cors(auth(app.GetSyncFragment, &proOnly)), true},
		Route{"GET", "/v1/sync/state", cors(auth(app.GetSyncState, &proOnly)), true},

		Route{"OPTIONS", "/v1/books", cors(app.BooksOptions), false},
		Route{"GET", "/v1/demo/books", app.GetDemoBooks, true},
		Route{"GET", "/v1/books", cors(auth(app.GetBooks, &proOnly)), true},
		Route{"GET", "/v1/books/{bookUUID}", cors(auth(app.GetBook, &proOnly)), true},
		Route{"POST", "/v1/books", cors(app.CreateBook), false},
		Route{"PATCH", "/v1/books/{bookUUID}", cors(auth(app.UpdateBook, &proOnly)), false},
		Route{"DELETE", "/v1/books/{bookUUID}", cors(auth(app.DeleteBook, &proOnly)), false},

		Route{"OPTIONS", "/v1/notes", cors(app.NotesOptions), true},
		Route{"POST", "/v1/notes", cors(app.CreateNote), false},
		Route{"PATCH", "/v1/notes/{noteUUID}", auth(app.UpdateNote, &proOnly), false},
		Route{"DELETE", "/v1/notes/{noteUUID}", auth(app.DeleteNote, &proOnly), false},

		Route{"POST", "/v1/register", app.register, true},
		Route{"GET", "/v1/presignin", cors(app.presignin), true},
		Route{"POST", "/v1/signin", cors(app.signin), true},
		Route{"OPTIONS", "/v1/signout", cors(app.signoutOptions), true},
		Route{"POST", "/v1/signout", cors(app.signout), true},

		// v2
		Route{"OPTIONS", "/v2/notes", cors(app.NotesOptionsV2), true},
		Route{"POST", "/v2/notes", cors(auth(app.CreateNoteV2, &proOnly)), true},

		Route{"OPTIONS", "/v2/books", cors(app.BooksOptionsV2), true},
		Route{"POST", "/v2/books", cors(auth(app.CreateBookV2, &proOnly)), true},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(applyMiddleware(handler, route.RateLimit))
	}

	return router
}
