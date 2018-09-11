package core

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/dnote/actions"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

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

func LogActionRemoveNote(tx *sql.Tx, noteUUID, bookName string) error {
	b, err := json.Marshal(actions.RemoveNoteDataV1{
		NoteUUID: noteUUID,
		BookName: bookName,
	})
	if err != nil {
		return errors.Wrap(err, "marshalling data into JSON")
	}

	ts := time.Now().Unix()
	if err := LogAction(tx, 1, actions.ActionRemoveNote, string(b), ts); err != nil {
		return errors.Wrapf(err, "logging action")
	}

	return nil
}

func LogActionEditNote(tx *sql.Tx, noteUUID, bookName, content string, ts int64) error {
	b, err := json.Marshal(actions.EditNoteDataV2{
		NoteUUID: noteUUID,
		FromBook: bookName,
		Content:  &content,
	})
	if err != nil {
		return errors.Wrap(err, "marshalling data into JSON")
	}

	if err := LogAction(tx, 2, actions.ActionEditNote, string(b), ts); err != nil {
		return errors.Wrapf(err, "logging action")
	}

	return nil
}

func LogActionAddBook(tx *sql.Tx, name string) error {
	b, err := json.Marshal(actions.AddBookDataV1{
		BookName: name,
	})
	if err != nil {
		return errors.Wrap(err, "marshalling data into JSON")
	}

	_, err = tx.Exec("INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)",
		uuid.NewV4().String(), 2, actions.ActionAddBook, string(b), time.Now().Unix())
	if err != nil {
		return errors.Wrap(err, "inserting an action")
	}
	return nil
}

// LogActionRemoveBook logs an action for removing book
func LogActionRemoveBook(tx *sql.Tx, name string) error {
	b, err := json.Marshal(actions.RemoveBookDataV1{BookName: name})
	if err != nil {
		return errors.Wrap(err, "marshalling data into JSON")
	}

	ts := time.Now().Unix()
	if err := LogAction(tx, 1, actions.ActionRemoveBook, string(b), ts); err != nil {
		return errors.Wrapf(err, "logging action")
	}

	return nil
}
