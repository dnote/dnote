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
