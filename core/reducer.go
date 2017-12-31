package core

import (
	"github.com/dnote-io/cli/infra"
	"github.com/pkg/errors"
)

type addNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookUUID string `json:"book_uuid"`
	Content  string `json:"content"`
}

type editNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookUUID string `json:"book_uuid"`
	Content  string `json:"content"`
}

type removeNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookUUID string `json:"book_uuid"`
}

type addBookData struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type removeBookData struct {
	UUID string `json:"uuid"`
}

func parseAddNoteData(raw map[string]interface{}) (addNoteData, error) {
	ret := addNoteData{}

	noteUUID, ok := raw["note_uuid"].(string)
	if !ok {
		return ret, errors.Errorf("invalid note_uuid for action %s. Got %s", ActionAddNote, noteUUID)
	}
	bookUUID, ok := raw["book_uuid"].(string)
	if !ok {
		return ret, errors.Errorf("invalid book_uuid for action %s. Got %s", ActionAddNote, bookUUID)
	}
	content, ok := raw["content"].(string)
	if !ok {
		return ret, errors.Errorf("invalid content for action %s. Got %s", ActionAddNote, content)
	}

	ret.NoteUUID = noteUUID
	ret.BookUUID = bookUUID
	ret.Content = content

	return ret, nil
}

func parseAddBookData(raw map[string]interface{}) (addBookData, error) {
	ret := addBookData{}

	uuid, ok := raw["uuid"].(string)
	if !ok {
		return ret, errors.Errorf("invalid uuid for action %s. Got %s", ActionAddBook, uuid)
	}
	name, ok := raw["name"].(string)
	if !ok {
		return ret, errors.Errorf("invalid name for action %s. Got %s", ActionAddBook, name)
	}

	ret.UUID = uuid
	ret.Name = name

	return ret, nil
}

func parseRemoveNoteData(raw map[string]interface{}) (removeNoteData, error) {
	ret := removeNoteData{}

	noteUUID, ok := raw["note_uuid"].(string)
	if !ok {
		return ret, errors.Errorf("invalid note_uuid for action %s. Got %s", ActionRemoveNote, noteUUID)
	}
	bookUUID, ok := raw["book_uuid"].(string)
	if !ok {
		return ret, errors.Errorf("invalid book_uuid for action %s. Got %s", ActionRemoveNote, bookUUID)
	}

	ret.NoteUUID = noteUUID
	ret.BookUUID = bookUUID

	return ret, nil
}

func parseRemoveBookData(raw map[string]interface{}) (removeBookData, error) {
	ret := removeBookData{}

	uuid, ok := raw["uuid"].(string)
	if !ok {
		return ret, errors.Errorf("invalid uuid for action %s. Got %s", ActionRemoveBook, uuid)
	}

	ret.UUID = uuid

	return ret, nil
}

func parseEditNoteData(raw map[string]interface{}) (editNoteData, error) {
	ret := editNoteData{}

	noteUUID, ok := raw["note_uuid"].(string)
	if !ok {
		return ret, errors.Errorf("invalid note_uuid for action %s. Got %s", ActionEditNote, noteUUID)
	}
	bookUUID, ok := raw["book_uuid"].(string)
	if !ok {
		return ret, errors.Errorf("invalid book_uuid for action %s. Got %s", ActionEditNote, bookUUID)
	}
	content, ok := raw["content"].(string)
	if !ok {
		return ret, errors.Errorf("invalid content for action %s. Got %s", ActionEditNote, content)
	}

	ret.NoteUUID = noteUUID
	ret.BookUUID = bookUUID
	ret.Content = content

	return ret, nil
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
	data, err := parseAddNoteData(action.Data)
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
	book, err := GetBookByUUID(dnote, data.BookUUID)
	if err != nil {
		return errors.Wrap(err, "Failed to find the book by uuid")
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
	data, err := parseRemoveNoteData(action.Data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}
	book, err := GetBookByUUID(dnote, data.BookUUID)
	if err != nil {
		return errors.Wrap(err, "Failed to find the book by uuid")
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
	data, err := parseEditNoteData(action.Data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}
	book, err := GetBookByUUID(dnote, data.BookUUID)
	if err != nil {
		return errors.Wrap(err, "Failed to find the book by uuid")
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
	data, err := parseAddBookData(action.Data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}

	book := infra.Book{
		UUID:  data.UUID,
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
	data, err := parseRemoveBookData(action.Data)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the action data")
	}

	dnote, err := GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}

	for bookName, book := range dnote {
		if book.UUID == data.UUID {
			delete(dnote, bookName)
		}
	}

	err = WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	return nil
}
