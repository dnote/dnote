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

package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/pkg/errors"
)

// MustNewServer is a test utility function to initialize a new server
// with the given app
func MustNewServer(t *testing.T, a *app.App) *httptest.Server {
	server, err := NewServer(a)
	if err != nil {
		t.Fatal(errors.Wrap(err, "initializing router"))
	}

	return server
}

func NewServer(a *app.App) (*httptest.Server, error) {
	ctl := New(a)
	rc := RouteConfig{
		WebRoutes:   NewWebRoutes(a, ctl),
		APIRoutes:   NewAPIRoutes(a, ctl),
		Controllers: ctl,
	}
	r, err := NewRouter(a, rc)
	if err != nil {
		return nil, errors.Wrap(err, "initializing router")
	}

	server := httptest.NewServer(r)

	return server, nil
}
