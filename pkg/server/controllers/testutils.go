/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/pkg/errors"
)

// MustNewServer is a test utility function to initialize a new server
// with the given app paratmers
func MustNewServer(t *testing.T, appParams *app.App) *httptest.Server {
	a := app.NewTest(appParams)

	ctl := New(&a)
	rc := RouteConfig{
		WebRoutes:   NewWebRoutes(&a, ctl),
		APIRoutes:   NewAPIRoutes(&a, ctl),
		Controllers: ctl,
	}
	r, err := NewRouter(&a, rc)
	if err != nil {
		t.Fatal(errors.Wrap(err, "initializing router"))
	}

	server := httptest.NewServer(r)

	return server
}
