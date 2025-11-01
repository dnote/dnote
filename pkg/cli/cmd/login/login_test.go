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
