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

// Package context defines dnote context
package context

import (
	"net/http"

	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/clock"
)

// Paths contain directory definitions
type Paths struct {
	Home        string
	Config      string
	Data        string
	Cache       string
	LegacyDnote string
}

// DnoteCtx is a context holding the information of the current runtime
type DnoteCtx struct {
	Paths              Paths
	APIEndpoint        string
	Version            string
	DB                 *database.DB
	SessionKey         string
	SessionKeyExpiry   int64
	Editor             string
	Clock              clock.Clock
	EnableUpgradeCheck bool
	HTTPClient         *http.Client
}

// Redact replaces private information from the context with a set of
// placeholder values.
func Redact(ctx DnoteCtx) DnoteCtx {
	var sessionKey string
	if ctx.SessionKey != "" {
		sessionKey = "1"
	} else {
		sessionKey = "0"
	}
	ctx.SessionKey = sessionKey

	return ctx
}
