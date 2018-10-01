package core

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/dnote/actions"
	"github.com/pkg/errors"
)

// LogActionAddNote logs an action for adding a note
func LogActionAddNote(tx *sql.Tx, noteUUID, bookName, content string, timestamp int64) error {
	b, err := json.Marshal(actions.AddNoteDataV2{
		NoteUUID: noteUUID,
		BookName: bookName,
		Content:  content,
		// TODO: support adding a public note
		Public: false,
	})
	if err != nil {
		return errors.Wrap(err, "marshalling data into JSON")
	}

	if err := LogAction(tx, 2, actions.ActionAddNote, string(b), timestamp); err != nil {
		return errors.Wrapf(err, "logging action")
	}

	return nil
}

// LogActionRemoveNote logs an action for removing a book
func LogActionRemoveNote(tx *sql.Tx, noteUUID, bookName string) error {
	b, err := json.Marshal(actions.RemoveNoteDataV1{
		NoteUUID: noteUUID,
		BookName: bookName,
	})
	if err != nil {
		return errors.Wrap(err, "marshalling data into JSON")
	}

	ts := time.Now().UnixNano()
	if err := LogAction(tx, 1, actions.ActionRemoveNote, string(b), ts); err != nil {
		return errors.Wrapf(err, "logging action")
	}

	return nil
}

// LogActionEditNote logs an action for editing a note
func LogActionEditNote(tx *sql.Tx, noteUUID, bookName, content string, ts int64) error {
	b, err := json.Marshal(actions.EditNoteDataV3{
		NoteUUID: noteUUID,
		Content:  &content,
		BookName: nil,
		Public:   nil,
	})
	if err != nil {
		return errors.Wrap(err, "marshalling data into JSON")
	}

	if err := LogAction(tx, 3, actions.ActionEditNote, string(b), ts); err != nil {
		return errors.Wrapf(err, "logging action")
	}

	return nil
}

// LogActionAddBook logs an action for adding a book
func LogActionAddBook(tx *sql.Tx, name string) error {
	b, err := json.Marshal(actions.AddBookDataV1{
		BookName: name,
	})
	if err != nil {
		return errors.Wrap(err, "marshalling data into JSON")
	}

	ts := time.Now().UnixNano()
	if err := LogAction(tx, 1, actions.ActionAddBook, string(b), ts); err != nil {
		return errors.Wrapf(err, "logging action")
	}

	return nil
}

// LogActionRemoveBook logs an action for removing book
func LogActionRemoveBook(tx *sql.Tx, name string) error {
	b, err := json.Marshal(actions.RemoveBookDataV1{BookName: name})
	if err != nil {
		return errors.Wrap(err, "marshalling data into JSON")
	}

	ts := time.Now().UnixNano()
	if err := LogAction(tx, 1, actions.ActionRemoveBook, string(b), ts); err != nil {
		return errors.Wrapf(err, "logging action")
	}

	return nil
}
