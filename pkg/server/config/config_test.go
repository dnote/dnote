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
				DB: PostgresConfig{
					Host: "mockHost",
					Port: "5432",
					Name: "mockDB",
					User: "mockUser",
				},
				WebURL: "http://mock.url",
				Port:   "3000",
			},
			expectedErr: nil,
		},
		{
			config: Config{
				DB: PostgresConfig{
					Port: "5432",
					Name: "mockDB",
					User: "mockUser",
				},
				WebURL: "http://mock.url",
				Port:   "3000",
			},
			expectedErr: ErrDBMissingHost,
		},
		{
			config: Config{
				DB: PostgresConfig{
					Host: "mockHost",
					Name: "mockDB",
					User: "mockUser",
				},
				WebURL: "http://mock.url",
				Port:   "3000",
			},
			expectedErr: ErrDBMissingPort,
		},
		{
			config: Config{
				DB: PostgresConfig{
					Host: "mockHost",
					Port: "5432",
					User: "mockUser",
				},
				WebURL: "http://mock.url",
				Port:   "3000",
			},
			expectedErr: ErrDBMissingName,
		},
		{
			config: Config{
				DB: PostgresConfig{
					Host: "mockHost",
					Port: "5432",
					Name: "mockDB",
				},
				WebURL: "http://mock.url",
				Port:   "3000",
			},
			expectedErr: ErrDBMissingUser,
		},
		{
			config: Config{
				DB: PostgresConfig{
					Host: "mockHost",
					Port: "5432",
					Name: "mockDB",
					User: "mockUser",
				},
			},
			expectedErr: ErrWebURLInvalid,
		},
		{
			config: Config{
				DB: PostgresConfig{
					Host: "mockHost",
					Port: "5432",
					Name: "mockDB",
					User: "mockUser",
				},
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
