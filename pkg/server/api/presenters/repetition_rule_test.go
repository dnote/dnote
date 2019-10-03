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

package presenters

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
)

func TestPresentRepetitionRule(t *testing.T) {
	b1 := database.Book{UUID: "1cf8794f-4d61-4a9d-a9da-18f8db9e53cc", Label: "foo"}
	b2 := database.Book{UUID: "ede00f3b-eab1-469c-ae12-c60cebeeef17", Label: "bar"}
	d1 := database.RepetitionRule{
		UUID:       "c725afb5-8bf1-4581-a0e7-0f683c15f3d0",
		Title:      "test title",
		Enabled:    true,
		Hour:       1,
		Minute:     2,
		BookDomain: database.BookDomainAll,
		Books:      []database.Book{b1, b2},
	}

	testCases := []struct {
		input    database.RepetitionRule
		expected RepetitionRule
	}{
		{
			input: d1,
			expected: RepetitionRule{
				UUID:       d1.UUID,
				Title:      d1.Title,
				Enabled:    d1.Enabled,
				Hour:       d1.Hour,
				Minute:     d1.Minute,
				BookDomain: d1.BookDomain,
				Books: []Book{
					{
						UUID:      b1.UUID,
						USN:       b1.USN,
						CreatedAt: b1.CreatedAt,
						UpdatedAt: b1.UpdatedAt,
						Label:     b1.Label,
					},
					{
						UUID:      b2.UUID,
						USN:       b2.USN,
						CreatedAt: b2.CreatedAt,
						UpdatedAt: b2.UpdatedAt,
						Label:     b2.Label,
					},
				},
				CreatedAt: d1.CreatedAt,
				UpdatedAt: d1.UpdatedAt,
			},
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			result := PresentRepetitionRule(tc.input)

			assert.DeepEqual(t, result, tc.expected, "result mismatch")
		})
	}
}
