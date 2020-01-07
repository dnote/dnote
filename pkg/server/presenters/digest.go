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

// Digest is a presented digest
type Digest struct {
	UUID           string         `json:"uuid"`
	Version        int            `json:"version"`
	RepetitionRule RepetitionRule `json:"repetition_rule"`
	Notes          []DigestNote   `json:"notes"`
	IsRead         bool           `json:"is_read"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// DigestNote is a presented note inside a digest
type DigestNote struct {
	Note
	IsReviewed bool `json:"is_reviewed"`
}

func presentDigestNote(note database.Note) DigestNote {
	ret := DigestNote{
		Note:       PresentNote(note),
		IsReviewed: note.NoteReview.UUID != "",
	}

	return ret
}

func presentDigestNotes(notes []database.Note) []DigestNote {
	ret := []DigestNote{}

	for _, note := range notes {
		n := presentDigestNote(note)
		ret = append(ret, n)
	}

	return ret
}

// PresentDigest presents a digest
func PresentDigest(digest database.Digest) Digest {
	ret := Digest{
		UUID:           digest.UUID,
		Notes:          presentDigestNotes(digest.Notes),
		Version:        digest.Version,
		RepetitionRule: PresentRepetitionRule(digest.Rule),
		IsRead:         len(digest.Receipts) > 0,
		CreatedAt:      digest.CreatedAt,
		UpdatedAt:      digest.UpdatedAt,
	}

	return ret
}

// PresentDigests presetns digests
func PresentDigests(digests []database.Digest) []Digest {
	ret := []Digest{}

	for _, digest := range digests {
		p := PresentDigest(digest)
		ret = append(ret, p)
	}

	return ret
}
