package core

import (
	"encoding/json"
	"time"

	"github.com/dnote-io/cli/infra"
	"github.com/pkg/errors"
)

var (
	ActionAddNote    = "add_note"
	ActionRemoveNote = "remove_note"
	ActionEditNote   = "edit_note"
	ActionAddBook    = "add_book"
	ActionRemoveBook = "remove_book"
)

type Action struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

func LogActionAddNote(ctx infra.DnoteCtx, noteUUID, bookName, content string) error {
	b, err := json.Marshal(AddNoteData{
		NoteUUID: noteUUID,
		BookName: bookName,
		Content:  content,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := Action{
		Type:      ActionAddNote,
		Data:      b,
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionAddNote)
	}

	return nil
}

func LogActionRemoveNote(ctx infra.DnoteCtx, noteUUID, bookName string) error {
	b, err := json.Marshal(RemoveNoteData{
		NoteUUID: noteUUID,
		BookName: bookName,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := Action{
		Type:      ActionRemoveNote,
		Data:      b,
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionRemoveNote)
	}

	return nil
}

func LogActionEditNote(ctx infra.DnoteCtx, noteUUID, bookName, content string) error {
	b, err := json.Marshal(EditNoteData{
		NoteUUID: noteUUID,
		BookName: bookName,
		Content:  content,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := Action{
		Type:      ActionEditNote,
		Data:      b,
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionEditNote)
	}

	return nil
}

func LogActionAddBook(ctx infra.DnoteCtx, name string) error {
	b, err := json.Marshal(AddBookData{
		Name: name,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := Action{
		Type:      ActionAddBook,
		Data:      b,
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionAddBook)
	}

	return nil
}

func LogActionRemoveBook(ctx infra.DnoteCtx, name string) error {
	b, err := json.Marshal(RemoveBookData{Name: name})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data into JSON")
	}

	action := Action{
		Type:      ActionRemoveBook,
		Data:      b,
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionRemoveBook)
	}

	return nil
}
