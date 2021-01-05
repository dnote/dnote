package context

import (
	"context"

	"github.com/dnote/dnote/pkg/server/database"
)

const (
	userKey  privateKey = "user"
	tokenKey privateKey = "token"
)

type privateKey string

// WithUser creates a new context with the given user
func WithUser(ctx context.Context, user *database.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// WithToken creates a new context with the given user
func WithToken(ctx context.Context, tok *database.Token) context.Context {
	return context.WithValue(ctx, tokenKey, tok)
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

// Token retrieves a token from the given context.
func Token(ctx context.Context) *database.Token {
	if temp := ctx.Value(tokenKey); temp != nil {
		if tok, ok := temp.(*database.Token); ok {
			return tok
		}
	}

	return nil
}
