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
	"time"

	"github.com/dnote/dnote/pkg/server/database"
)

// RepetitionRule is a presented digest rule
type RepetitionRule struct {
	UUID       string    `json:"uuid"`
	Title      string    `json:"title"`
	Enabled    bool      `json:"enabled"`
	Hour       int       `json:"hour" gorm:"index"`
	Minute     int       `json:"minute" gorm:"index"`
	Frequency  int       `json:"frequency"`
	BookDomain string    `json:"book_domain"`
	LastActive int       `json:"last_active"`
	Books      []Book    `json:"books"`
	NoteCount  int       `json:"note_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PresentRepetitionRule presents a digest rule
func PresentRepetitionRule(d database.RepetitionRule) RepetitionRule {
	ret := RepetitionRule{
		UUID:       d.UUID,
		Title:      d.Title,
		Enabled:    d.Enabled,
		Hour:       d.Hour,
		Minute:     d.Minute,
		Frequency:  d.Frequency,
		BookDomain: d.BookDomain,
		NoteCount:  d.NoteCount,
		LastActive: d.LastActive,
		Books:      PresentBooks(d.Books),
		CreatedAt:  FormatTS(d.CreatedAt),
		UpdatedAt:  FormatTS(d.UpdatedAt),
	}

	return ret
}

// PresentRepetitionRules presents a slice of digest rules
func PresentRepetitionRules(ds []database.RepetitionRule) []RepetitionRule {
	ret := []RepetitionRule{}

	for _, d := range ds {
		p := PresentRepetitionRule(d)
		ret = append(ret, p)
	}

	return ret
}
