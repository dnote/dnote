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
	"os"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// API is a web API configuration
type API struct {
	App *app.App
}

// init sets up the application based on the configuration
func (a *API) init() error {
	if err := a.App.Validate(); err != nil {
		return errors.Wrap(err, "validating the app parameters")
	}

	return nil
}

func applyMiddleware(h http.HandlerFunc, rateLimit bool) http.Handler {
	ret := h
	ret = handlers.Logging(ret)

	if rateLimit && os.Getenv("GO_ENV") != "TEST" {
		ret = handlers.Limit(ret)
	}

	return ret
}

// NewRouter creates and returns a new router
func NewRouter(a *API) (*mux.Router, error) {
	if err := a.init(); err != nil {
		return nil, errors.Wrap(err, "initializing app")
	}

	proOnly := handlers.AuthParams{ProOnly: true}
	app := a.App

	var routes = []handlers.Route{
		// internal
		{Method: "GET", Pattern: "/health", HandlerFunc: a.checkHealth, RateLimit: false},
		{Method: "GET", Pattern: "/me", HandlerFunc: handlers.Auth(app, a.getMe, nil), RateLimit: true},
		{Method: "POST", Pattern: "/verification-token", HandlerFunc: handlers.Auth(app, a.createVerificationToken, nil), RateLimit: true},
		{Method: "PATCH", Pattern: "/verify-email", HandlerFunc: a.verifyEmail, RateLimit: true},
		{Method: "POST", Pattern: "/reset-token", HandlerFunc: a.createResetToken, RateLimit: true},
		{Method: "PATCH", Pattern: "/reset-password", HandlerFunc: a.resetPassword, RateLimit: true},
		{Method: "PATCH", Pattern: "/account/profile", HandlerFunc: handlers.Auth(app, a.updateProfile, nil), RateLimit: true},
		{Method: "PATCH", Pattern: "/account/password", HandlerFunc: handlers.Auth(app, a.updatePassword, nil), RateLimit: true},
		{Method: "GET", Pattern: "/account/email-preference", HandlerFunc: handlers.TokenAuth(app, a.getEmailPreference, database.TokenTypeEmailPreference, nil), RateLimit: true},
		{Method: "PATCH", Pattern: "/account/email-preference", HandlerFunc: handlers.TokenAuth(app, a.updateEmailPreference, database.TokenTypeEmailPreference, nil), RateLimit: true},
		{Method: "GET", Pattern: "/notes", HandlerFunc: handlers.Auth(app, a.getNotes, nil), RateLimit: false},
		{Method: "GET", Pattern: "/notes/{noteUUID}", HandlerFunc: a.getNote, RateLimit: true},
		{Method: "GET", Pattern: "/calendar", HandlerFunc: handlers.Auth(app, a.getCalendar, nil), RateLimit: true},

		// v3
		{Method: "GET", Pattern: "/v3/sync/fragment", HandlerFunc: handlers.Cors(handlers.Auth(app, a.GetSyncFragment, &proOnly)), RateLimit: false},
		{Method: "GET", Pattern: "/v3/sync/state", HandlerFunc: handlers.Cors(handlers.Auth(app, a.GetSyncState, &proOnly)), RateLimit: false},
		{Method: "OPTIONS", Pattern: "/v3/books", HandlerFunc: handlers.Cors(a.BooksOptions), RateLimit: true},
		{Method: "GET", Pattern: "/v3/books", HandlerFunc: handlers.Cors(handlers.Auth(app, a.GetBooks, &proOnly)), RateLimit: true},
		{Method: "GET", Pattern: "/v3/books/{bookUUID}", HandlerFunc: handlers.Cors(handlers.Auth(app, a.GetBook, &proOnly)), RateLimit: true},
		{Method: "POST", Pattern: "/v3/books", HandlerFunc: handlers.Cors(handlers.Auth(app, a.CreateBook, &proOnly)), RateLimit: false},
		{Method: "PATCH", Pattern: "/v3/books/{bookUUID}", HandlerFunc: handlers.Cors(handlers.Auth(app, a.UpdateBook, &proOnly)), RateLimit: false},
		{Method: "DELETE", Pattern: "/v3/books/{bookUUID}", HandlerFunc: handlers.Cors(handlers.Auth(app, a.DeleteBook, &proOnly)), RateLimit: false},
		{Method: "OPTIONS", Pattern: "/v3/notes", HandlerFunc: handlers.Cors(a.NotesOptions), RateLimit: true},
		{Method: "POST", Pattern: "/v3/notes", HandlerFunc: handlers.Cors(handlers.Auth(app, a.CreateNote, &proOnly)), RateLimit: false},
		{Method: "PATCH", Pattern: "/v3/notes/{noteUUID}", HandlerFunc: handlers.Auth(app, a.UpdateNote, &proOnly), RateLimit: false},
		{Method: "DELETE", Pattern: "/v3/notes/{noteUUID}", HandlerFunc: handlers.Auth(app, a.DeleteNote, &proOnly), RateLimit: false},
		{Method: "POST", Pattern: "/v3/signin", HandlerFunc: handlers.Cors(a.signin), RateLimit: true},
		{Method: "OPTIONS", Pattern: "/v3/signout", HandlerFunc: handlers.Cors(a.signoutOptions), RateLimit: true},
		{Method: "POST", Pattern: "/v3/signout", HandlerFunc: handlers.Cors(a.signout), RateLimit: true},
		{Method: "POST", Pattern: "/v3/register", HandlerFunc: a.register, RateLimit: true},
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
