package core

import (
	"encoding/json"
	"time"

	"github.com/dnote/actions"
	"github.com/dnote/cli/infra"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

func LogActionAddNote(ctx infra.DnoteCtx, noteUUID, bookName, content string, timestamp int64) error {
	b, err := json.Marshal(actions.AddNoteDataV2{
		NoteUUID: noteUUID,
		BookName: bookName,
		Content:  content,
		// TODO: support adding a public note
		Public: false,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := actions.Action{
		UUID:      uuid.NewV4().String(),
		Schema:    2,
		Type:      actions.ActionAddNote,
		Data:      b,
		Timestamp: timestamp,
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", actions.ActionAddNote)
	}

	return nil
}

func LogActionRemoveNote(ctx infra.DnoteCtx, noteUUID, bookName string) error {
	b, err := json.Marshal(actions.RemoveNoteDataV1{
		NoteUUID: noteUUID,
		BookName: bookName,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := actions.Action{
		UUID:      uuid.NewV4().String(),
		Schema:    1,
		Type:      actions.ActionRemoveNote,
		Data:      b,
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", actions.ActionRemoveNote)
	}

	return nil
}

func LogActionEditNote(ctx infra.DnoteCtx, noteUUID, bookName, content string, ts int64) error {
	b, err := json.Marshal(actions.EditNoteDataV1{
		NoteUUID: noteUUID,
		FromBook: bookName,
		Content:  content,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := actions.Action{
		UUID:      uuid.NewV4().String(),
		Schema:    2,
		Type:      actions.ActionEditNote,
		Data:      b,
		Timestamp: ts,
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", actions.ActionEditNote)
	}

	return nil
}

func LogActionAddBook(ctx infra.DnoteCtx, name string) error {
	b, err := json.Marshal(actions.AddBookDataV1{
		BookName: name,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := actions.Action{
		UUID:      uuid.NewV4().String(),
		Schema:    1,
		Type:      actions.ActionAddBook,
		Data:      b,
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", actions.ActionAddBook)
	}

	return nil
}

func LogActionRemoveBook(ctx infra.DnoteCtx, name string) error {
	b, err := json.Marshal(actions.RemoveBookDataV1{BookName: name})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := actions.Action{
		UUID:      uuid.NewV4().String(),
		Schema:    1,
		Type:      actions.ActionRemoveBook,
		Data:      b,
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", actions.ActionRemoveBook)
	}

	return nil
}
