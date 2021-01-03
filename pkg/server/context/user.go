package context

import (
	"context"

	"github.com/dnote/dnote/pkg/server/database"
)

const (
	userKey privateKey = "user"
)

type privateKey string

// WithUser creates a new context with the given user
func WithUser(ctx context.Context, user *database.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User retrieves a user from the given context. It returns a pointer to
// a user. If the context does not contain a user, it returns nil.
func User(ctx context.Context) *database.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*database.User); ok {
			return user
		}
	}

	return nil
}
