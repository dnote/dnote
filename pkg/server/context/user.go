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
