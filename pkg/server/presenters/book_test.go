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

package presenters

import (
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
)

func TestPresentBook(t *testing.T) {
	createdAt := time.Date(2025, 1, 15, 10, 30, 45, 123456789, time.UTC)
	updatedAt := time.Date(2025, 2, 20, 14, 45, 30, 987654321, time.UTC)

	testCases := []struct {
		name     string
		input    database.Book
		expected Book
	}{
		{
			name: "basic book",
			input: database.Book{
				Model: database.Model{
					ID:        1,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				UUID:   "a1b2c3d4-e5f6-4789-a012-3456789abcde",
				UserID: 42,
				Label:  "JavaScript",
				USN:    100,
			},
			expected: Book{
				UUID:      "a1b2c3d4-e5f6-4789-a012-3456789abcde",
				USN:       100,
				CreatedAt: FormatTS(createdAt),
				UpdatedAt: FormatTS(updatedAt),
				Label:     "JavaScript",
			},
		},
		{
			name: "book with special characters in label",
			input: database.Book{
				Model: database.Model{
					ID:        2,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				UUID:   "f1e2d3c4-b5a6-4987-b654-321fedcba098",
				UserID: 99,
				Label:  "C++",
				USN:    200,
			},
			expected: Book{
				UUID:      "f1e2d3c4-b5a6-4987-b654-321fedcba098",
				USN:       200,
				CreatedAt: FormatTS(createdAt),
				UpdatedAt: FormatTS(updatedAt),
				Label:     "C++",
			},
		},
		{
			name: "book with empty label",
			input: database.Book{
				Model: database.Model{
					ID:        3,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				UUID:   "12345678-90ab-4cde-8901-234567890abc",
				UserID: 1,
				Label:  "",
				USN:    0,
			},
			expected: Book{
				UUID:      "12345678-90ab-4cde-8901-234567890abc",
				USN:       0,
				CreatedAt: FormatTS(createdAt),
				UpdatedAt: FormatTS(updatedAt),
				Label:     "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := PresentBook(tc.input)

			assert.Equal(t, got.UUID, tc.expected.UUID, "UUID mismatch")
			assert.Equal(t, got.USN, tc.expected.USN, "USN mismatch")
			assert.Equal(t, got.Label, tc.expected.Label, "Label mismatch")
			assert.Equal(t, got.CreatedAt, tc.expected.CreatedAt, "CreatedAt mismatch")
			assert.Equal(t, got.UpdatedAt, tc.expected.UpdatedAt, "UpdatedAt mismatch")
		})
	}
}

func TestPresentBooks(t *testing.T) {
	createdAt1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt1 := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	createdAt2 := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	updatedAt2 := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name     string
		input    []database.Book
		expected []Book
	}{
		{
			name:     "empty slice",
			input:    []database.Book{},
			expected: []Book{},
		},
		{
			name: "single book",
			input: []database.Book{
				{
					Model: database.Model{
						ID:        1,
						CreatedAt: createdAt1,
						UpdatedAt: updatedAt1,
					},
					UUID:   "9a8b7c6d-5e4f-4321-9876-543210fedcba",
					UserID: 1,
					Label:  "Go",
					USN:    10,
				},
			},
			expected: []Book{
				{
					UUID:      "9a8b7c6d-5e4f-4321-9876-543210fedcba",
					USN:       10,
					CreatedAt: FormatTS(createdAt1),
					UpdatedAt: FormatTS(updatedAt1),
					Label:     "Go",
				},
			},
		},
		{
			name: "multiple books",
			input: []database.Book{
				{
					Model: database.Model{
						ID:        1,
						CreatedAt: createdAt1,
						UpdatedAt: updatedAt1,
					},
					UUID:   "9a8b7c6d-5e4f-4321-9876-543210fedcba",
					UserID: 1,
					Label:  "Go",
					USN:    10,
				},
				{
					Model: database.Model{
						ID:        2,
						CreatedAt: createdAt2,
						UpdatedAt: updatedAt2,
					},
					UUID:   "abcdef01-2345-4678-9abc-def012345678",
					UserID: 1,
					Label:  "Python",
					USN:    20,
				},
			},
			expected: []Book{
				{
					UUID:      "9a8b7c6d-5e4f-4321-9876-543210fedcba",
					USN:       10,
					CreatedAt: FormatTS(createdAt1),
					UpdatedAt: FormatTS(updatedAt1),
					Label:     "Go",
				},
				{
					UUID:      "abcdef01-2345-4678-9abc-def012345678",
					USN:       20,
					CreatedAt: FormatTS(createdAt2),
					UpdatedAt: FormatTS(updatedAt2),
					Label:     "Python",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := PresentBooks(tc.input)

			assert.Equal(t, len(got), len(tc.expected), "Length mismatch")

			for i := range got {
				assert.Equal(t, got[i].UUID, tc.expected[i].UUID, "UUID mismatch")
				assert.Equal(t, got[i].USN, tc.expected[i].USN, "USN mismatch")
				assert.Equal(t, got[i].Label, tc.expected[i].Label, "Label mismatch")
				assert.Equal(t, got[i].CreatedAt, tc.expected[i].CreatedAt, "CreatedAt mismatch")
				assert.Equal(t, got[i].UpdatedAt, tc.expected[i].UpdatedAt, "UpdatedAt mismatch")
			}
		})
	}
}
