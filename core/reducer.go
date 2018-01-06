package core

import (
	"encoding/json"
	"sort"

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
	BookName string `json:"book_name"`
}

type RemoveBookData struct {
	BookName string `json:"book_name"`
}

// ReduceAll reduces all actions
func ReduceAll(ctx infra.DnoteCtx, actions []Action) error {
	for _, action := range actions {
		if err := Reduce(ctx, action); err != nil {
			return errors.Wrap(err, "Failed to reduce action")
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
		return errors.Wrapf(err, "Failed to process the action %s", action.Type)
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
		AddedOn: action.Timestamp,
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

	sort.SliceStable(notes, func(i, j int) bool {
		return notes[i].AddedOn < notes[j].AddedOn
	})

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
			note.EditedOn = action.Timestamp
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

	_, exists := dnote[data.BookName]
	if exists {
		// If book already exists, another machine added a book with the same name.
		// noop
		return nil
	}

	book := infra.Book{
		Name:  data.BookName,
		Notes: []infra.Note{},
	}
	dnote[data.BookName] = book

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
		if bookName == data.BookName {
			delete(dnote, bookName)
		}
	}

	err = WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	return nil
}
