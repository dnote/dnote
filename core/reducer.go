package core

import (
	"encoding/json"
	"sort"

	"github.com/dnote/actions"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
)

// ReduceAll reduces all actions
func ReduceAll(ctx infra.DnoteCtx, ats []actions.Action) error {
	for _, action := range ats {
		if err := Reduce(ctx, action); err != nil {
			return errors.Wrap(err, "Failed to reduce action")
		}
	}

	return nil
}

// Reduce transitions the local dnote state by consuming the action returned
// from the server
func Reduce(ctx infra.DnoteCtx, action actions.Action) error {
	var err error

	switch action.Type {
	case actions.ActionAddNote:
		err = handleAddNote(ctx, action)
	case actions.ActionRemoveNote:
		err = handleRemoveNote(ctx, action)
	case actions.ActionEditNote:
		err = handleEditNote(ctx, action)
	case actions.ActionAddBook:
		err = handleAddBook(ctx, action)
	case actions.ActionRemoveBook:
		err = handleRemoveBook(ctx, action)
	default:
		return errors.Errorf("Unsupported action %s", action.Type)
	}

	if err != nil {
		return errors.Wrapf(err, "Failed to process the action %s", action.Type)
	}

	return nil
}

func handleAddNote(ctx infra.DnoteCtx, action actions.Action) error {
	var data actions.AddNoteDataV1
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	log.Debug("reducing add_note. action: %+v. data: %+v\n", action, data)

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

func handleRemoveNote(ctx infra.DnoteCtx, action actions.Action) error {
	var data actions.RemoveNoteDataV1
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	log.Debug("reducing remove_note. action: %+v. data: %+v\n", action, data)

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

func handleEditNoteV1(ctx infra.DnoteCtx, action actions.Action) error {
	var data actions.EditNoteDataV1
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	log.Debug("reducing edit_note v1. action: %+v. data: %+v\n", action, data)

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}
	fromBook, ok := dnote[data.FromBook]
	if !ok {
		return errors.Errorf("Origin book with a name %s is not found", data.FromBook)
	}

	if data.ToBook == "" {
		for idx, note := range fromBook.Notes {
			if note.UUID == data.NoteUUID {
				note.Content = data.Content
				note.EditedOn = action.Timestamp
				dnote[fromBook.Name].Notes[idx] = note
			}
		}
	} else {
		// Change the book

		toBook, ok := dnote[data.ToBook]
		if !ok {
			return errors.Errorf("Destination book with a name %s is not found", data.FromBook)
		}

		var index int
		var note infra.Note

		// Find the note
		for idx := range fromBook.Notes {
			note = fromBook.Notes[idx]

			if note.UUID == data.NoteUUID {
				index = idx
			}
		}

		note.Content = data.Content
		note.EditedOn = action.Timestamp

		dnote[fromBook.Name] = GetUpdatedBook(dnote[fromBook.Name], append(fromBook.Notes[:index], fromBook.Notes[index+1:]...))
		dnote[toBook.Name] = GetUpdatedBook(dnote[toBook.Name], append(toBook.Notes, note))
	}

	err = WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	return nil
}

func handleEditNoteV2(ctx infra.DnoteCtx, action actions.Action) error {
	var data actions.EditNoteDataV2
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	log.Debug("reducing edit_note v2. action: %+v. data: %+v\n", action, data)

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}

	fromBook, ok := dnote[data.FromBook]
	if !ok {
		return errors.Errorf("Origin book with a name %s is not found", data.FromBook)
	}

	if data.ToBook == nil {
		for idx, note := range fromBook.Notes {
			if note.UUID == data.NoteUUID {
				if data.Content != nil {
					note.Content = *data.Content
				}
				if data.Public != nil {
					note.Public = *data.Public
				}

				note.EditedOn = action.Timestamp
				dnote[fromBook.Name].Notes[idx] = note
			}
		}
	} else {
		// Change the book
		toBook := *data.ToBook

		dstBook, ok := dnote[toBook]
		if !ok {
			return errors.Errorf("Destination book with a name %s is not found", toBook)
		}

		var index int
		var note infra.Note

		// Find the note
		for idx := range fromBook.Notes {
			note = fromBook.Notes[idx]

			if note.UUID == data.NoteUUID {
				index = idx
			}
		}

		if data.Content != nil {
			note.Content = *data.Content
		}
		if data.Public != nil {
			note.Public = *data.Public
		}
		note.EditedOn = action.Timestamp

		dnote[fromBook.Name] = GetUpdatedBook(dnote[fromBook.Name], append(fromBook.Notes[:index], fromBook.Notes[index+1:]...))
		dnote[toBook] = GetUpdatedBook(dnote[toBook], append(dstBook.Notes, note))
	}

	err = WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	return nil
}

func handleEditNote(ctx infra.DnoteCtx, action actions.Action) error {
	if action.Schema == 1 {
		return handleEditNoteV1(ctx, action)
	} else if action.Schema == 2 {
		return handleEditNoteV2(ctx, action)
	}

	return errors.Errorf("Unsupported schema version for editing note: %d", action.Schema)
}

func handleAddBook(ctx infra.DnoteCtx, action actions.Action) error {
	var data actions.AddBookDataV1
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	log.Debug("reducing add_book. action: %+v. data: %+v\n", action, data)

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

func handleRemoveBook(ctx infra.DnoteCtx, action actions.Action) error {
	var data actions.RemoveBookDataV1
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	log.Debug("reducing remove_book. action: %+v. data: %+v\n", action, data)

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}

	for bookName := range dnote {
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
