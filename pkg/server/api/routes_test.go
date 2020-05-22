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

package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func TestNotSupportedVersions(t *testing.T) {
	testCases := []struct {
		path string
	}{
		// v1
		{
			path: "/v1",
		},
		{
			path: "/v1/foo",
		},
		{
			path: "/v1/bar/baz",
		},
		// v2
		{
			path: "/v2",
		},
		{
			path: "/v2/foo",
		},
		{
			path: "/v2/bar/baz",
		},
	}

	// setup
	server := MustNewServer(t, &app.App{
		DB:    &gorm.DB{},
		Clock: clock.NewMock(),
	})
	defer server.Close()

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			// execute
			req := testutils.MakeReq(server.URL, "GET", tc.path, "")
			res := testutils.HTTPDo(t, req)

			// test
			assert.Equal(t, res.StatusCode, http.StatusGone, "status code mismatch")
		})
	}
}

func TestNewRouter_AppValidate(t *testing.T) {
	c := config.Load()

	configWithoutWebURL := config.Load()
	configWithoutWebURL.WebURL = ""

	testCases := []struct {
		app         app.App
		expectedErr error
	}{
		{
			app: app.App{
				DB:             &gorm.DB{},
				Clock:          clock.NewMock(),
				EmailTemplates: mailer.Templates{},
				EmailBackend:   &testutils.MockEmailbackendImplementation{},
				Config:         c,
			},
			expectedErr: nil,
		},
		{
			app: app.App{
				DB:             nil,
				Clock:          clock.NewMock(),
				EmailTemplates: mailer.Templates{},
				EmailBackend:   &testutils.MockEmailbackendImplementation{},
				Config:         c,
			},
			expectedErr: app.ErrEmptyDB,
		},
		{
			app: app.App{
				DB:             &gorm.DB{},
				Clock:          nil,
				EmailTemplates: mailer.Templates{},
				EmailBackend:   &testutils.MockEmailbackendImplementation{},
				Config:         c,
			},
			expectedErr: app.ErrEmptyClock,
		},
		{
			app: app.App{
				DB:             &gorm.DB{},
				Clock:          clock.NewMock(),
				EmailTemplates: nil,
				EmailBackend:   &testutils.MockEmailbackendImplementation{},
				Config:         c,
			},
			expectedErr: app.ErrEmptyEmailTemplates,
		},
		{
			app: app.App{
				DB:             &gorm.DB{},
				Clock:          clock.NewMock(),
				EmailTemplates: mailer.Templates{},
				EmailBackend:   nil,
				Config:         c,
			},
			expectedErr: app.ErrEmptyEmailBackend,
		},
		{
			app: app.App{
				DB:             &gorm.DB{},
				Clock:          clock.NewMock(),
				EmailTemplates: mailer.Templates{},
				EmailBackend:   &testutils.MockEmailbackendImplementation{},
				Config:         configWithoutWebURL,
			},
			expectedErr: app.ErrEmptyWebURL,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			api := API{App: &tc.app}
			_, err := NewRouter(&api)

			assert.Equal(t, errors.Cause(err), tc.expectedErr, "error mismatch")
		})
	}
}
