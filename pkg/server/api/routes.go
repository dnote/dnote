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
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/handlers"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/gorilla/mux"
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
func NewRouter(a *API) (*mux.Router, error) {
	if err := a.init(); err != nil {
		return nil, errors.Wrap(err, "initializing app")
	}

	proOnly := handlers.AuthParams{ProOnly: true}
	app := a.App

	var routes = []Route{
		// internal
		{"GET", "/health", a.checkHealth, false},
		{"GET", "/me", handlers.Auth(app, a.getMe, nil), true},
		{"POST", "/verification-token", handlers.Auth(app, a.createVerificationToken, nil), true},
		{"PATCH", "/verify-email", a.verifyEmail, true},
		{"POST", "/reset-token", a.createResetToken, true},
		{"PATCH", "/reset-password", a.resetPassword, true},
		{"PATCH", "/account/profile", handlers.Auth(app, a.updateProfile, nil), true},
		{"PATCH", "/account/password", handlers.Auth(app, a.updatePassword, nil), true},
		{"GET", "/account/email-preference", handlers.TokenAuth(app, a.getEmailPreference, database.TokenTypeEmailPreference, nil), true},
		{"PATCH", "/account/email-preference", handlers.TokenAuth(app, a.updateEmailPreference, database.TokenTypeEmailPreference, nil), true},
		{"POST", "/subscriptions", handlers.Auth(app, a.createSub, nil), true},
		{"PATCH", "/subscriptions", handlers.Auth(app, a.updateSub, nil), true},
		{"POST", "/webhooks/stripe", a.stripeWebhook, true},
		{"GET", "/subscriptions", handlers.Auth(app, a.getSub, nil), true},
		{"GET", "/stripe_source", handlers.Auth(app, a.getStripeSource, nil), true},
		{"PATCH", "/stripe_source", handlers.Auth(app, a.updateStripeSource, nil), true},
		{"GET", "/notes", handlers.Auth(app, a.getNotes, nil), false},
		{"GET", "/notes/{noteUUID}", a.getNote, true},
		{"GET", "/calendar", handlers.Auth(app, a.getCalendar, nil), true},

		// v3
		{"GET", "/v3/sync/fragment", handlers.Cors(handlers.Auth(app, a.GetSyncFragment, &proOnly)), false},
		{"GET", "/v3/sync/state", handlers.Cors(handlers.Auth(app, a.GetSyncState, &proOnly)), false},
		{"OPTIONS", "/v3/books", handlers.Cors(a.BooksOptions), true},
		{"GET", "/v3/books", handlers.Cors(handlers.Auth(app, a.GetBooks, &proOnly)), true},
		{"GET", "/v3/books/{bookUUID}", handlers.Cors(handlers.Auth(app, a.GetBook, &proOnly)), true},
		{"POST", "/v3/books", handlers.Cors(handlers.Auth(app, a.CreateBook, &proOnly)), false},
		{"PATCH", "/v3/books/{bookUUID}", handlers.Cors(handlers.Auth(app, a.UpdateBook, &proOnly)), false},
		{"DELETE", "/v3/books/{bookUUID}", handlers.Cors(handlers.Auth(app, a.DeleteBook, &proOnly)), false},
		{"OPTIONS", "/v3/notes", handlers.Cors(a.NotesOptions), true},
		{"POST", "/v3/notes", handlers.Cors(handlers.Auth(app, a.CreateNote, &proOnly)), false},
		{"PATCH", "/v3/notes/{noteUUID}", handlers.Auth(app, a.UpdateNote, &proOnly), false},
		{"DELETE", "/v3/notes/{noteUUID}", handlers.Auth(app, a.DeleteNote, &proOnly), false},
		{"POST", "/v3/signin", handlers.Cors(a.signin), true},
		{"OPTIONS", "/v3/signout", handlers.Cors(a.signoutOptions), true},
		{"POST", "/v3/signout", handlers.Cors(a.signout), true},
		{"POST", "/v3/register", a.register, true},
	}

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/v1").Handler(applyMiddleware(handlers.NotSupported, true))
	router.PathPrefix("/v2").Handler(applyMiddleware(handlers.NotSupported, true))

	for _, route := range routes {
		handler := route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(applyMiddleware(handler, route.RateLimit))
	}

	return router, nil
}
