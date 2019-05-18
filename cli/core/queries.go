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
