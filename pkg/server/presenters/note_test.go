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

func TestPresentNote(t *testing.T) {
	createdAt := time.Date(2025, 1, 15, 10, 30, 45, 123456789, time.UTC)
	updatedAt := time.Date(2025, 2, 20, 14, 45, 30, 987654321, time.UTC)

	input := database.Note{
		Model: database.Model{
			ID:        1,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		UUID:     "a1b2c3d4-e5f6-4789-a012-3456789abcde",
		UserID:   42,
		BookUUID: "f1e2d3c4-b5a6-4987-b654-321fedcba098",
		Body:     "Test note content",
		AddedOn:  1234567890,
		USN:      100,
		Book: database.Book{
			UUID:  "f1e2d3c4-b5a6-4987-b654-321fedcba098",
			Label: "JavaScript",
		},
		User: database.User{
			UUID: "9a8b7c6d-5e4f-4321-9876-543210fedcba",
		},
	}

	got := PresentNote(input)

	assert.Equal(t, got.UUID, "a1b2c3d4-e5f6-4789-a012-3456789abcde", "UUID mismatch")
	assert.Equal(t, got.Body, "Test note content", "Body mismatch")
	assert.Equal(t, got.AddedOn, int64(1234567890), "AddedOn mismatch")
	assert.Equal(t, got.USN, 100, "USN mismatch")
	assert.Equal(t, got.CreatedAt, FormatTS(createdAt), "CreatedAt mismatch")
	assert.Equal(t, got.UpdatedAt, FormatTS(updatedAt), "UpdatedAt mismatch")
	assert.Equal(t, got.Book.UUID, "f1e2d3c4-b5a6-4987-b654-321fedcba098", "Book UUID mismatch")
	assert.Equal(t, got.Book.Label, "JavaScript", "Book Label mismatch")
	assert.Equal(t, got.User.UUID, "9a8b7c6d-5e4f-4321-9876-543210fedcba", "User UUID mismatch")
}

func TestPresentNotes(t *testing.T) {
	createdAt1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt1 := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	createdAt2 := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	updatedAt2 := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)

	input := []database.Note{
		{
			Model: database.Model{
				ID:        1,
				CreatedAt: createdAt1,
				UpdatedAt: updatedAt1,
			},
			UUID:     "a1b2c3d4-e5f6-4789-a012-3456789abcde",
			UserID:   1,
			BookUUID: "f1e2d3c4-b5a6-4987-b654-321fedcba098",
			Body:     "First note",
			AddedOn:  1000000000,
			USN:      10,
			Book: database.Book{
				UUID:  "f1e2d3c4-b5a6-4987-b654-321fedcba098",
				Label: "Go",
			},
			User: database.User{
				UUID: "9a8b7c6d-5e4f-4321-9876-543210fedcba",
			},
		},
		{
			Model: database.Model{
				ID:        2,
				CreatedAt: createdAt2,
				UpdatedAt: updatedAt2,
			},
			UUID:     "12345678-90ab-4cde-8901-234567890abc",
			UserID:   1,
			BookUUID: "abcdef01-2345-4678-9abc-def012345678",
			Body:     "Second note",
			AddedOn:  2000000000,
			USN:      20,
			Book: database.Book{
				UUID:  "abcdef01-2345-4678-9abc-def012345678",
				Label: "Python",
			},
			User: database.User{
				UUID: "9a8b7c6d-5e4f-4321-9876-543210fedcba",
			},
		},
	}

	got := PresentNotes(input)

	assert.Equal(t, len(got), 2, "Length mismatch")
	assert.Equal(t, got[0].UUID, "a1b2c3d4-e5f6-4789-a012-3456789abcde", "Note 0 UUID mismatch")
	assert.Equal(t, got[0].Body, "First note", "Note 0 Body mismatch")
	assert.Equal(t, got[1].UUID, "12345678-90ab-4cde-8901-234567890abc", "Note 1 UUID mismatch")
	assert.Equal(t, got[1].Body, "Second note", "Note 1 Body mismatch")
}
