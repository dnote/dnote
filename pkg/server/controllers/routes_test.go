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
