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

package dirs

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

type envTestCase struct {
	envKey   string
	envVal   string
	got      *string
	expected string
}

func testCustomDirs(t *testing.T, testCases []envTestCase) {
	for _, tc := range testCases {
		t.Setenv(tc.envKey, tc.envVal)

		Reload()

		assert.Equal(t, *tc.got, tc.expected, "result mismatch")
	}
}
