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

package controllers

import (
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/assets"
	mw "github.com/dnote/dnote/pkg/server/middleware"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Route represents a single route
type Route struct {
	Method    string
	Pattern   string
	Handler   http.HandlerFunc
	RateLimit bool
}

// RouteConfig is the configuration for routes
type RouteConfig struct {
	Controllers *Controllers
	WebRoutes   []Route
	APIRoutes   []Route
}

// NewWebRoutes returns a new web routes
func NewWebRoutes(a *app.App, c *Controllers) []Route {
	redirectGuest := &mw.AuthParams{RedirectGuestsToLogin: true}

	ret := []Route{
		{"GET", "/", mw.Auth(a, c.Users.Settings, redirectGuest), true},
		{"GET", "/about", mw.Auth(a, c.Users.About, redirectGuest), true},
		{"GET", "/login", mw.GuestOnly(a, c.Users.NewLogin), true},
		{"POST", "/login", mw.GuestOnly(a, c.Users.Login), true},
		{"POST", "/logout", c.Users.Logout, true},

		{"GET", "/password-reset", c.Users.PasswordResetView.ServeHTTP, true},
		{"PATCH", "/password-reset", c.Users.PasswordReset, true},
		{"GET", "/password-reset/{token}", c.Users.PasswordResetConfirm, true},
		{"POST", "/reset-token", c.Users.CreateResetToken, true},
		{"POST", "/verification-token", mw.Auth(a, c.Users.CreateEmailVerificationToken, redirectGuest), true},
		{"GET", "/verify-email/{token}", mw.Auth(a, c.Users.VerifyEmail, redirectGuest), true},
		{"PATCH", "/account/profile", mw.Auth(a, c.Users.ProfileUpdate, nil), true},
		{"PATCH", "/account/password", mw.Auth(a, c.Users.PasswordUpdate, nil), true},

		{"GET", "/health", c.Health.Index, true},
	}

	if !a.Config.DisableRegistration {
		ret = append(ret, Route{"GET", "/join", c.Users.New, true})
		ret = append(ret, Route{"POST", "/join", c.Users.Create, true})
	}

	return ret
}

// NewAPIRoutes returns a new api routes
func NewAPIRoutes(a *app.App, c *Controllers) []Route {

	proOnly := mw.AuthParams{ProOnly: true}

	return []Route{
		// v3
		{"GET", "/v3/sync/fragment", mw.Cors(mw.Auth(a, c.Sync.GetSyncFragment, &proOnly)), false},
		{"GET", "/v3/sync/state", mw.Cors(mw.Auth(a, c.Sync.GetSyncState, &proOnly)), false},
		{"POST", "/v3/signin", mw.Cors(c.Users.V3Login), true},
		{"POST", "/v3/signout", mw.Cors(c.Users.V3Logout), true},
		{"OPTIONS", "/v3/signout", mw.Cors(c.Users.logoutOptions), true},
		{"GET", "/v3/notes", mw.Cors(mw.Auth(a, c.Notes.V3Index, nil)), true},
		{"GET", "/v3/notes/{noteUUID}", c.Notes.V3Show, true},
		{"POST", "/v3/notes", mw.Cors(mw.Auth(a, c.Notes.V3Create, nil)), true},
		{"DELETE", "/v3/notes/{noteUUID}", mw.Cors(mw.Auth(a, c.Notes.V3Delete, nil)), true},
		{"PATCH", "/v3/notes/{noteUUID}", mw.Cors(mw.Auth(a, c.Notes.V3Update, nil)), true},
		{"OPTIONS", "/v3/notes", mw.Cors(c.Notes.IndexOptions), true},
		{"GET", "/v3/books", mw.Cors(mw.Auth(a, c.Books.V3Index, nil)), true},
		{"GET", "/v3/books/{bookUUID}", mw.Cors(mw.Auth(a, c.Books.V3Show, nil)), true},
		{"POST", "/v3/books", mw.Cors(mw.Auth(a, c.Books.V3Create, nil)), true},
		{"PATCH", "/v3/books/{bookUUID}", mw.Cors(mw.Auth(a, c.Books.V3Update, nil)), true},
		{"DELETE", "/v3/books/{bookUUID}", mw.Cors(mw.Auth(a, c.Books.V3Delete, nil)), true},
		{"OPTIONS", "/v3/books", mw.Cors(c.Books.IndexOptions), true},
	}
}

func registerRoutes(router *mux.Router, wrapper mw.Middleware, app *app.App, routes []Route) {
	for _, route := range routes {
		wrappedHandler := wrapper(route.Handler, app, route.RateLimit)

		router.
			Handle(route.Pattern, wrappedHandler).
			Methods(route.Method)
	}
}

// NewRouter creates and returns a new router
func NewRouter(app *app.App, rc RouteConfig) (http.Handler, error) {
	if err := app.Validate(); err != nil {
		return nil, errors.Wrap(err, "validating the app parameters")
	}

	router := mux.NewRouter().StrictSlash(true)

	webRouter := router.PathPrefix("/").Subrouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	registerRoutes(webRouter, mw.WebMw, app, rc.WebRoutes)
	registerRoutes(apiRouter, mw.APIMw, app, rc.APIRoutes)

	router.PathPrefix("/api/v1").Handler(mw.ApplyLimit(mw.NotSupported, true))
	router.PathPrefix("/api/v2").Handler(mw.ApplyLimit(mw.NotSupported, true))

	// static
	staticFs, err := assets.GetStaticFS()
	if err != nil {
		return nil, errors.Wrap(err, "getting the filesystem for static files")
	}

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.FS(staticFs)))
	router.PathPrefix("/static/").Handler(staticHandler)

	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nAllow: /"))
	})

	// catch-all
	router.PathPrefix("/").HandlerFunc(rc.Controllers.Static.NotFound)

	return mw.Global(router), nil
}
