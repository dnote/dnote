package routes

import (
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/log"
	"net/http"
)

func userMw(inner http.Handler, app *app.App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _, err := AuthWithSession(app, r, nil)
		if err != nil {
			log.ErrorWrap(err, "authenticating with session")
			inner.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, &user)
		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
