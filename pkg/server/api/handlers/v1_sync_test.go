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
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/pkg/errors"
)

func TestParseGetSyncFragmentQuery(t *testing.T) {
	testCases := []struct {
		input    string
		afterUSN int
		limit    int
		err      error
	}{
		{
			input:    `after_usn=50&limit=50`,
			afterUSN: 50,
			limit:    50,
			err:      nil,
		},
		{
			input:    `limit=50`,
			afterUSN: 0,
			limit:    50,
			err:      nil,
		},
		{
			input:    `after_usn=50`,
			afterUSN: 50,
			limit:    100,
			err:      nil,
		},
		{
			input:    `after_usn=50&limit=100`,
			afterUSN: 50,
			limit:    100,
			err:      nil,
		},
		{
			input:    "",
			afterUSN: 0,
			limit:    100,
			err:      nil,
		},
		{
			input:    "limit=101",
			afterUSN: 0,
			limit:    0,
			err: &queryParamError{
				key:     "limit",
				value:   "101",
				message: "maximum value is 100",
			},
		},
	}

	for idx, tc := range testCases {
		q, err := url.ParseQuery(tc.input)
		if err != nil {
			t.Fatal(errors.Wrap(err, "parsing test input"))
		}

		afterUSN, limit, err := parseGetSyncFragmentQuery(q)
		ok := reflect.DeepEqual(err, tc.err)
		assert.Equal(t, ok, true, fmt.Sprintf("err mismatch for test case %d. Expected: %+v. Got: %+v", idx, tc.err, err))

		assert.Equal(t, afterUSN, tc.afterUSN, fmt.Sprintf("afterUSN mismatch for test case %d", idx))
		assert.Equal(t, limit, tc.limit, fmt.Sprintf("limit mismatch for test case %d", idx))
	}
}
