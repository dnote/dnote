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
