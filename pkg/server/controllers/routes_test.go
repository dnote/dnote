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

package controllers

import (
	"net/http"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestNotSupportedVersions(t *testing.T) {
	testCases := []struct {
		path string
	}{
		// v1
		{
			path: "/api/v1",
		},
		{
			path: "/api/v1/foo",
		},
		{
			path: "/api/v1/bar/baz",
		},
		// v2
		{
			path: "/api/v2",
		},
		{
			path: "/api/v2/foo",
		},
		{
			path: "/api/v2/bar/baz",
		},
	}

	// setup
	db := testutils.InitMemoryDB(t)
	a := app.NewTest()
	a.Clock = clock.NewMock()
	a.DB = db
	server := MustNewServer(t, &a)
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
