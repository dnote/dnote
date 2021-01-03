package routes

import (
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/controllers"
	"github.com/gorilla/mux"
)

// Route represents a single route
type Route struct {
	Method    string
	Pattern   string
	Handler   http.Handler
	RateLimit bool
}

func registerRoutes(router *mux.Router, mw middleware, app *app.App, routes []Route) {
	for _, route := range routes {
		wrappedHandler := mw(route.Handler, app, route.RateLimit)

		router.
			Handle(route.Pattern, wrappedHandler).
			Methods(route.Method)
	}
}

// NewWebRoutes returns a new web routes
func NewWebRoutes(app *app.App, c *controllers.Controllers) []Route {
	return []Route{
		{"GET", "/", Auth(app, http.HandlerFunc(c.Users.New), nil), true},
		{"GET", "/new", http.HandlerFunc(c.Users.New), true},
	}
}

// NewAPIRoutes returns a new api routes
func NewAPIRoutes(c *controllers.Controllers) []Route {
	return []Route{}
}

// Config is the configuration for routes
type Config struct {
	Controllers *controllers.Controllers
	WebRoutes   []Route
	APIRoutes   []Route
}

// New creates and returns a new router
func New(app *app.App, rc Config) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	webRouter := router.PathPrefix("/").Subrouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	registerRoutes(webRouter, WebMw, app, rc.WebRoutes)
	registerRoutes(apiRouter, APIMw, app, rc.APIRoutes)

	// static
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(app.Config.StaticDir)))
	router.PathPrefix("/static/").Handler(staticHandler)

	// catch-all
	router.PathPrefix("/").HandlerFunc(rc.Controllers.Static.NotFound)

	return LoggingMw(router)
}
