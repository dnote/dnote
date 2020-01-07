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

package web

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestInit(t *testing.T) {
	mockIndexHTML := []byte("<html></html>")
	mockRobotsTxt := []byte("Allow: *")
	mockServiceWorkerJs := []byte("function() {}")
	mockStaticFileSystem := http.Dir(".")

	testCases := []struct {
		ctx         Context
		expectedErr error
	}{
		{
			ctx: Context{
				DB:               testutils.DB,
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: nil,
		},
		{
			ctx: Context{
				DB:               nil,
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: ErrEmptyDatabase,
		},
		{
			ctx: Context{
				DB:               testutils.DB,
				IndexHTML:        nil,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: ErrEmptyIndexHTML,
		},
		{
			ctx: Context{
				DB:               testutils.DB,
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        nil,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: ErrEmptyRobotsTxt,
		},
		{
			ctx: Context{
				DB:               testutils.DB,
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  nil,
				StaticFileSystem: mockStaticFileSystem,
			},
			expectedErr: ErrEmptyServiceWorkerJS,
		},
		{
			ctx: Context{
				DB:               testutils.DB,
				IndexHTML:        mockIndexHTML,
				RobotsTxt:        mockRobotsTxt,
				ServiceWorkerJs:  mockServiceWorkerJs,
				StaticFileSystem: nil,
			},
			expectedErr: ErrEmptyStaticFileSystem,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			_, err := Init(tc.ctx)

			assert.Equal(t, errors.Cause(err), tc.expectedErr, "error mismatch")
		})
	}
}
