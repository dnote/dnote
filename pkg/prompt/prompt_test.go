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

package prompt

import (
	"strings"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestFormatQuestion(t *testing.T) {
	testCases := []struct {
		question   string
		optimistic bool
		expected   string
	}{
		{
			question:   "Are you sure?",
			optimistic: false,
			expected:   "Are you sure? (y/N)",
		},
		{
			question:   "Continue?",
			optimistic: true,
			expected:   "Continue? (Y/n)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.question, func(t *testing.T) {
			result := FormatQuestion(tc.question, tc.optimistic)
			assert.Equal(t, result, tc.expected, "formatted question mismatch")
		})
	}
}

func TestReadYesNo(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		optimistic bool
		expected   bool
	}{
		{
			name:       "pessimistic with y",
			input:      "y\n",
			optimistic: false,
			expected:   true,
		},
		{
			name:       "pessimistic with Y (uppercase)",
			input:      "Y\n",
			optimistic: false,
			expected:   true,
		},
		{
			name:       "pessimistic with n",
			input:      "n\n",
			optimistic: false,
			expected:   false,
		},
		{
			name:       "pessimistic with empty",
			input:      "\n",
			optimistic: false,
			expected:   false,
		},
		{
			name:       "pessimistic with whitespace",
			input:      "  \n",
			optimistic: false,
			expected:   false,
		},
		{
			name:       "optimistic with y",
			input:      "y\n",
			optimistic: true,
			expected:   true,
		},
		{
			name:       "optimistic with n",
			input:      "n\n",
			optimistic: true,
			expected:   false,
		},
		{
			name:       "optimistic with empty",
			input:      "\n",
			optimistic: true,
			expected:   true,
		},
		{
			name:       "optimistic with whitespace",
			input:      "  \n",
			optimistic: true,
			expected:   true,
		},
		{
			name:       "invalid input defaults to no",
			input:      "maybe\n",
			optimistic: false,
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a reader with test input
			reader := strings.NewReader(tc.input)

			// Test ReadYesNo
			result, err := ReadYesNo(reader, tc.optimistic)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Equal(t, result, tc.expected, "ReadYesNo result mismatch")
		})
	}
}

func TestReadYesNo_Error(t *testing.T) {
	// Test error case with EOF (empty reader)
	reader := strings.NewReader("")

	_, err := ReadYesNo(reader, false)
	if err == nil {
		t.Fatal("expected error when reading from empty reader")
	}
}

