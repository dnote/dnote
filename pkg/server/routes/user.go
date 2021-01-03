package routes

import (
	"github.com/dnote/dnote/pkg/server/app"
	"net/http"
)

func userMw(inner http.Handler, app *app.App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
