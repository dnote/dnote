/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package login

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/context"
)

func TestGetServerDisplayURL(t *testing.T) {
	testCases := []struct {
		apiEndpoint string
		expected    string
	}{
		{
			apiEndpoint: "https://api.getdnote.com",
			expected:    "https://www.getdnote.com",
		},
		{
			apiEndpoint: "https://dnote.mydomain.com/api",
			expected:    "https://dnote.mydomain.com",
		},
		{
			apiEndpoint: "https://mysubdomain.mydomain.com/dnote/api",
			expected:    "https://mysubdomain.mydomain.com",
		},
		{
			apiEndpoint: "https://dnote.mysubdomain.mydomain.com/api",
			expected:    "https://dnote.mysubdomain.mydomain.com",
		},
		{
			apiEndpoint: "some-string",
			expected:    "",
		},
		{
			apiEndpoint: "",
			expected:    "",
		},
		{
			apiEndpoint: "https://",
			expected:    "",
		},
		{
			apiEndpoint: "https://abc",
			expected:    "https://abc",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("for input %s", tc.apiEndpoint), func(t *testing.T) {
			got := getServerDisplayURL(context.DnoteCtx{APIEndpoint: tc.apiEndpoint})
			assert.Equal(t, got, tc.expected, "result mismatch")
		})
	}
}
