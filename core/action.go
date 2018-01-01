package core

import (
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
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

func LogActionAddNote(ctx infra.DnoteCtx, noteUUID, bookUUID, content string) error {
	action := Action{
		Type: ActionAddNote,
		Data: map[string]interface{}{
			"note_uuid": noteUUID,
			"book_uuid": bookUUID,
			"content":   content,
		},
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionAddNote)
	}

	return nil
}

func LogActionRemoveNote(ctx infra.DnoteCtx, noteUUID, bookUUID string) error {
	action := Action{
		Type: ActionRemoveNote,
		Data: map[string]interface{}{
			"note_uuid": noteUUID,
			"book_uuid": bookUUID,
		},
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionRemoveNote)
	}

	return nil
}

func LogActionEditNote(ctx infra.DnoteCtx, noteUUID, bookUUID, content string) error {
	action := Action{
		Type: ActionEditNote,
		Data: map[string]interface{}{
			"book_uuid": bookUUID,
			"note_uuid": noteUUID,
			"content":   content,
		},
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionEditNote)
	}

	return nil
}

func LogActionAddBook(ctx infra.DnoteCtx, uuid, name string) error {
	action := Action{
		Type: ActionAddBook,
		Data: map[string]interface{}{
			"uuid": uuid,
			"name": name,
		},
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionAddBook)
	}

	return nil
}

func LogActionRemoveBook(ctx infra.DnoteCtx, uuid string) error {
	action := Action{
		Type: ActionRemoveBook,
		Data: map[string]interface{}{
			"uuid": uuid,
		},
		Timestamp: time.Now().Unix(),
	}

	if err := LogAction(ctx, action); err != nil {
		return errors.Wrapf(err, "Failed to log action type %s", ActionRemoveBook)
	}

	return nil
}
