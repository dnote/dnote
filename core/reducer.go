package core

import (
	"encoding/json"

	"github.com/dnote-io/cli/infra"
	"github.com/pkg/errors"
)

type AddNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookName string `json:"book_name"`
	Content  string `json:"content"`
}

type EditNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookName string `json:"book_name"`
	Content  string `json:"content"`
}

type RemoveNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookName string `json:"book_name"`
}

type AddBookData struct {
	Name string `json:"name"`
}

type RemoveBookData struct {
	Name string `json:"name"`
}

// ReduceAll reduces all actions
func ReduceAll(ctx infra.DnoteCtx, actions []Action) error {
	for _, action := range actions {
		if err := Reduce(ctx, action); err != nil {
			return errors.Wrap(err, "Failed to reduce action")
		}
	}

	// After having consumed all actions, if bookmark is less than last_action,
	// bring it forward to be in sync with the server
	ts, err := ReadTimestamp(ctx)
	if ts.Bookmark < ts.LastAction {
		ts.Bookmark = ts.LastAction

		err = WriteTimestamp(ctx, ts)
		if err != nil {
			return errors.Wrap(err, "Failed to update last_sync")
		}
	}

	return nil
}

// Reduce transitions the local dnote state by consuming the action returned
// from the server
func Reduce(ctx infra.DnoteCtx, action Action) error {
	var err error

	switch action.Type {
	case ActionAddNote:
		err = handleAddNote(ctx, action)
	case ActionRemoveNote:
		err = handleRemoveNote(ctx, action)
	case ActionEditNote:
		err = handleEditNote(ctx, action)
	case ActionAddBook:
		err = handleAddBook(ctx, action)
	case ActionRemoveBook:
		err = handleRemoveBook(ctx, action)
	default:
		return errors.Errorf("Unsupported action %s", action.Type)
	}

	if err != nil {
		return errors.Wrap(err, "Failed to process the action")
	}

	// Update timestamp
	ts, err := ReadTimestamp(ctx)
	if ts.Bookmark < action.Timestamp {
		ts.Bookmark = action.Timestamp

		err = WriteTimestamp(ctx, ts)
		if err != nil {
			return errors.Wrap(err, "Failed to update last_sync")
		}
	}

	return nil
}

func handleAddNote(ctx infra.DnoteCtx, action Action) error {
	var data AddNoteData
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	note := infra.Note{
		UUID:    data.NoteUUID,
		Content: data.Content,
	}

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}
	book, ok := dnote[data.BookName]
	if !ok {
		return errors.Errorf("Book with a name %s is not found", data.BookName)
	}

	// Check duplicate
	for _, note := range book.Notes {
		if note.UUID == data.NoteUUID {
			return errors.New("Duplicate note exists")
		}
	}

	notes := append(dnote[book.Name].Notes, note)
	dnote[book.Name] = GetUpdatedBook(dnote[book.Name], notes)

	err = WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	return nil
}

func handleRemoveNote(ctx infra.DnoteCtx, action Action) error {
	var data RemoveNoteData
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}
	book, ok := dnote[data.BookName]
	if !ok {
		return errors.Errorf("Book with a name %s is not found", data.BookName)
	}

	notes := FilterNotes(book.Notes, func(note infra.Note) bool {
		return note.UUID != data.NoteUUID
	})
	dnote[book.Name] = GetUpdatedBook(dnote[book.Name], notes)

	err = WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	return nil
}

func handleEditNote(ctx infra.DnoteCtx, action Action) error {
	var data EditNoteData
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}
	book, ok := dnote[data.BookName]
	if !ok {
		return errors.Errorf("Book with a name %s is not found", data.BookName)
	}

	for idx, note := range book.Notes {
		if note.UUID == data.NoteUUID {
			note.Content = data.Content
			dnote[book.Name].Notes[idx] = note
		}
	}

	err = WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	return nil
}

func handleAddBook(ctx infra.DnoteCtx, action Action) error {
	var data AddBookData
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}

	book := infra.Book{
		Name:  data.Name,
		Notes: []infra.Note{},
	}
	dnote[data.Name] = book

	err = WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	return nil
}

func handleRemoveBook(ctx infra.DnoteCtx, action Action) error {
	var data RemoveBookData
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}

	for bookName, _ := range dnote {
		if bookName == data.Name {
			delete(dnote, bookName)
		}
	}

	err = WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	return nil
}
