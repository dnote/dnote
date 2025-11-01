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
