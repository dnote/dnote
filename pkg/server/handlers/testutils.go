/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package handlers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

// MustNewServer is a test utility function to initialize a new server
// with the given app paratmers
func MustNewServer(t *testing.T, a *app.App) *httptest.Server {
	a.WebURL = os.Getenv("WebURL")
	a.DB = testutils.DB

	// If email backend was not provided, use the default mock backend
	if a.EmailBackend == nil {
		a.EmailBackend = &testutils.MockEmailbackendImplementation{}
	}

	emailTmplDir := os.Getenv("DNOTE_TEST_EMAIL_TEMPLATE_DIR")
	a.EmailTemplates = mailer.NewTemplates(&emailTmplDir)

	api := API{App: a}

	r, err := api.NewRouter()
	if err != nil {
		t.Fatal(errors.Wrap(err, "initializing server"))
	}

	server := httptest.NewServer(r)

	return server
}
