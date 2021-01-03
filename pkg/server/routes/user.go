package routes

import (
	"net/http"
)

func userMw(inner http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// ctx = context.WithUser(ctx, user)
		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
