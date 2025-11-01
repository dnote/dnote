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

