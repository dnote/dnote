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
		{"GET", "/", mw.Auth(a.DB, c.Users.Settings, redirectGuest), true},
		{"GET", "/about", mw.Auth(a.DB, c.Users.About, redirectGuest), true},
		{"GET", "/login", mw.GuestOnly(a.DB, c.Users.NewLogin), true},
		{"POST", "/login", mw.GuestOnly(a.DB, c.Users.Login), true},
		{"POST", "/logout", c.Users.Logout, true},

		{"GET", "/password-reset", c.Users.PasswordResetView.ServeHTTP, true},
		{"PATCH", "/password-reset", c.Users.PasswordReset, true},
		{"GET", "/password-reset/{token}", c.Users.PasswordResetConfirm, true},
		{"POST", "/reset-token", c.Users.CreateResetToken, true},
		{"PATCH", "/account/profile", mw.Auth(a.DB, c.Users.ProfileUpdate, nil), true},
		{"PATCH", "/account/password", mw.Auth(a.DB, c.Users.PasswordUpdate, nil), true},

		{"GET", "/health", c.Health.Index, true},
	}

	if !a.DisableRegistration {
		ret = append(ret, Route{"GET", "/join", c.Users.New, true})
		ret = append(ret, Route{"POST", "/join", c.Users.Create, true})
	}

	return ret
}

// NewAPIRoutes returns a new api routes
func NewAPIRoutes(a *app.App, c *Controllers) []Route {
	return []Route{
		// v3
		{"GET", "/v3/sync/fragment", mw.Auth(a.DB, c.Sync.GetSyncFragment, nil), false},
		{"GET", "/v3/sync/state", mw.Auth(a.DB, c.Sync.GetSyncState, nil), false},
		{"POST", "/v3/signin", c.Users.V3Login, true},
		{"POST", "/v3/signout", c.Users.V3Logout, true},
		{"OPTIONS", "/v3/signout", c.Users.logoutOptions, true},
		{"GET", "/v3/notes", mw.Auth(a.DB, c.Notes.V3Index, nil), true},
		{"GET", "/v3/notes/{noteUUID}", mw.Auth(a.DB, c.Notes.V3Show, nil), true},
		{"POST", "/v3/notes", mw.Auth(a.DB, c.Notes.V3Create, nil), true},
		{"DELETE", "/v3/notes/{noteUUID}", mw.Auth(a.DB, c.Notes.V3Delete, nil), true},
		{"PATCH", "/v3/notes/{noteUUID}", mw.Auth(a.DB, c.Notes.V3Update, nil), true},
		{"OPTIONS", "/v3/notes", c.Notes.IndexOptions, true},
		{"GET", "/v3/books", mw.Auth(a.DB, c.Books.V3Index, nil), true},
		{"GET", "/v3/books/{bookUUID}", mw.Auth(a.DB, c.Books.V3Show, nil), true},
		{"POST", "/v3/books", mw.Auth(a.DB, c.Books.V3Create, nil), true},
		{"PATCH", "/v3/books/{bookUUID}", mw.Auth(a.DB, c.Books.V3Update, nil), true},
		{"DELETE", "/v3/books/{bookUUID}", mw.Auth(a.DB, c.Books.V3Delete, nil), true},
		{"OPTIONS", "/v3/books", c.Books.IndexOptions, true},
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
