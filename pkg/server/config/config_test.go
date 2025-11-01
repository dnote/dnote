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

package config

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/pkg/errors"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		config      Config
		expectedErr error
	}{
		{
			config: Config{
				DBPath: "test.db",
				WebURL: "http://mock.url",
				Port:   "3000",
			},
			expectedErr: nil,
		},
		{
			config: Config{
				DBPath: "",
				WebURL: "http://mock.url",
				Port:   "3000",
			},
			expectedErr: ErrDBMissingPath,
		},
		{
			config: Config{
				DBPath: "test.db",
			},
			expectedErr: ErrWebURLInvalid,
		},
		{
			config: Config{
				DBPath: "test.db",
				WebURL: "http://mock.url",
			},
			expectedErr: ErrPortInvalid,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			err := validate(tc.config)

			assert.Equal(t, errors.Cause(err), tc.expectedErr, "error mismatch")
		})
	}
}
