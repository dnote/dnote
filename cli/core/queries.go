/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package core

import (
	"database/sql"

	"github.com/dnote/dnote/cli/infra"
	"github.com/pkg/errors"
)

// NoteInfo is a basic information about a note
type NoteInfo struct {
	RowID     int
	BookLabel string
	UUID      string
	Content   string
	AddedOn   int64
	EditedOn  int64
}

// GetNoteInfo returns a NoteInfo for the note with the given noteRowID
func GetNoteInfo(ctx infra.DnoteCtx, noteRowID string) (NoteInfo, error) {
	var ret NoteInfo

	db := ctx.DB
	err := db.QueryRow(`SELECT books.label, notes.uuid, notes.body, notes.added_on, notes.edited_on, notes.rowid
			FROM notes
			INNER JOIN books ON books.uuid = notes.book_uuid
			WHERE notes.rowid = ? AND notes.deleted = false`, noteRowID).
		Scan(&ret.BookLabel, &ret.UUID, &ret.Content, &ret.AddedOn, &ret.EditedOn, &ret.RowID)
	if err == sql.ErrNoRows {
		return ret, errors.Errorf("note %s not found", noteRowID)
	} else if err != nil {
		return ret, errors.Wrap(err, "querying the note")
	}

	return ret, nil

}
