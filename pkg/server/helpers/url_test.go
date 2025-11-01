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

package helpers

import (
	"net/url"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestGetPath(t *testing.T) {
	t.Run("without query", func(t *testing.T) {
		// execute
		got := GetPath("/some-path", nil)

		// test
		assert.Equal(t, got, "/some-path", "got mismatch")
	})

	t.Run("with query", func(t *testing.T) {
		// execute
		q := url.Values{}
		q.Set("foo", "bar")
		q.Set("baz", "/quz")
		got := GetPath("/some-path", &q)

		// test
		assert.Equal(t, got, "/some-path?baz=%2Fquz&foo=bar", "got mismatch")
	})
}
