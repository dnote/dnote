package routes

import (
	"net/http"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/config"
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

func registerRoutes(router *mux.Router, mw middleware, c config.Config, routes []Route) {
	for _, route := range routes {
		wrappedHandler := mw(route.Handler, c, route.RateLimit)

		router.
			Handle(route.Pattern, wrappedHandler).
			Methods(route.Method)
	}
}

// NewWebRoutes returns a new web routes
func NewWebRoutes(cfg config.Config, c *controllers.Controllers, cl clock.Clock) []Route {
	return []Route{
		{"GET", "/", http.HandlerFunc(c.Users.New), true},
	}
}

// NewAPIRoutes returns a new api routes
func NewAPIRoutes(cfg config.Config, c *controllers.Controllers, cl clock.Clock) []Route {
	return []Route{}
}

// RouteConfig is the configuration for routes
type RouteConfig struct {
	Controllers *controllers.Controllers
	WebRoutes   []Route
	APIRoutes   []Route
}

// New creates and returns a new router
func New(cfg config.Config, rc RouteConfig) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	webRouter := router.PathPrefix("/").Subrouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	registerRoutes(webRouter, WebMw, cfg, rc.WebRoutes)
	registerRoutes(apiRouter, APIMw, cfg, rc.APIRoutes)

	// static
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.StaticDir)))
	router.PathPrefix("/static/").Handler(staticHandler)

	// catch-all
	// router.PathPrefix("/").HandlerFunc(rc.Controllers.Static.NotFound)

	return LoggingMw(router)
}
