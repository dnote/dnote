package controllers

import (
	"net/http"
	"os"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/handlers"
	"github.com/gorilla/mux"
)

// Route represents a single route
type Route struct {
	Method    string
	Pattern   string
	Handler   http.Handler
	RateLimit bool
}

// RouteConfig is the configuration for routes
type RouteConfig struct {
	Controllers *Controllers
	WebRoutes   []Route
	APIRoutes   []Route
}

// NewWebRoutes returns a new web routes
func NewWebRoutes(app *app.App, c *Controllers) []Route {
	ret := []Route{
		{"GET", "/", handlers.Auth(app, http.HandlerFunc(c.Notes.Index), &handlers.AuthParams{RedirectGuestsToLogin: true}), true},
		{"GET", "/login", c.Users.LoginView, true},
		{"POST", "/login", http.HandlerFunc(c.Users.Login), true},
		{"POST", "/logout", http.HandlerFunc(c.Users.Logout), true},
	}

	if !app.Config.DisableRegistration {
		ret = append(ret, Route{"GET", "/join", http.HandlerFunc(c.Users.New), true})
		ret = append(ret, Route{"POST", "/join", http.HandlerFunc(c.Users.Create), true})
	}

	return ret
}

// NewAPIRoutes returns a new api routes
func NewAPIRoutes(c *Controllers) []Route {
	return []Route{
		{"POST", "/v1/login", handlers.Cors(c.Users.V3Login), true},
		{"POST", "/v1/logout", handlers.Cors(c.Users.V3Logout), true},
	}
}

func applyMiddleware(h http.HandlerFunc, rateLimit bool) http.Handler {
	ret := h
	ret = handlers.Logging(ret)

	if rateLimit && os.Getenv("GO_ENV") != "TEST" {
		ret = handlers.Limit(ret)
	}

	return ret
}

func registerRoutes(router *mux.Router, mw handlers.Middleware, app *app.App, routes []Route) {
	for _, route := range routes {
		wrappedHandler := mw(route.Handler, app, route.RateLimit)

		router.
			Handle(route.Pattern, wrappedHandler).
			Methods(route.Method)
	}
}

// NewRouter creates and returns a new router
func NewRouter(app *app.App, rc RouteConfig) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	webRouter := router.PathPrefix("/").Subrouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	registerRoutes(webRouter, handlers.WebMw, app, rc.WebRoutes)
	registerRoutes(apiRouter, handlers.APIMw, app, rc.APIRoutes)

	// static
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(app.Config.StaticDir)))
	router.PathPrefix("/static/").Handler(staticHandler)

	// catch-all
	router.PathPrefix("/").HandlerFunc(rc.Controllers.Static.NotFound)

	return handlers.Logging(router)
}
