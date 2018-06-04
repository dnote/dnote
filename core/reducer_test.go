package core

import (
	"encoding/json"
	"testing"

	"github.com/dnote-io/cli/testutils"
	"github.com/pkg/errors"
)

func TestReduceAddNote(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("../tmp")

	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)
	testutils.WriteFile(ctx, "../testutils/fixtures/dnote4.json", "dnote")

	// Execute
	b, err := json.Marshal(&AddNoteData{
		Content:  "new content",
		BookName: "js",
		NoteUUID: "06896551-8a06-4996-89cc-0d866308b0f6",
	})
	action := Action{
		Type:      ActionAddNote,
		Data:      b,
		Timestamp: 1517629805,
	}
	if err = Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	book := dnote["js"]
	otherBook := dnote["linux"]
	existingNote := book.Notes[0]
	newNote := book.Notes[1]

	testutils.AssertEqual(t, len(book.Notes), 2, "notes length mismatch")
	testutils.AssertEqual(t, len(otherBook.Notes), 0, "other book notes length mismatch")
	testutils.AssertEqual(t, existingNote.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "existing note uuid mismatch")
	testutils.AssertEqual(t, existingNote.Content, "Booleans have toString()", "existing note content mismatch")
	testutils.AssertEqual(t, newNote.UUID, "06896551-8a06-4996-89cc-0d866308b0f6", "new note uuid mismatch")
	testutils.AssertEqual(t, newNote.Content, "new content", "new note content mismatch")
	testutils.AssertEqual(t, newNote.AddedOn, int64(1517629805), "new note added_on mismatch")
}

func TestReduceAddNote_SortByAddedOn(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("../tmp")

	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)
	testutils.WriteFile(ctx, "../testutils/fixtures/dnote3.json", "dnote")

	// Execute
	b, err := json.Marshal(&AddNoteData{
		Content:  "new content",
		BookName: "js",
		NoteUUID: "06896551-8a06-4996-89cc-0d866308b0f6",
	})
	action := Action{
		Type:      ActionAddNote,
		Data:      b,
		Timestamp: 1515199944,
	}
	if err = Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	book := dnote["js"]
	otherBook := dnote["linux"]
	note1 := book.Notes[0]
	note2 := book.Notes[1]
	note3 := book.Notes[2]

	testutils.AssertEqual(t, len(book.Notes), 3, "notes length mismatch")
	testutils.AssertEqual(t, len(otherBook.Notes), 1, "other book notes length mismatch")
	testutils.AssertEqual(t, note1.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "existing note 1 uuid mismatch")
	testutils.AssertEqual(t, note1.Content, "Booleans have toString()", "existing note 1 content mismatch")
	testutils.AssertEqual(t, note2.UUID, "06896551-8a06-4996-89cc-0d866308b0f6", "new note uuid mismatch")
	testutils.AssertEqual(t, note2.Content, "new content", "new note content mismatch")
	testutils.AssertEqual(t, note2.AddedOn, int64(1515199944), "new note added_on mismatch")
	testutils.AssertEqual(t, note3.UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "existing note 2 uuid mismatch")
	testutils.AssertEqual(t, note3.Content, "Date object implements mathematical comparisons", "existing note 2 content mismatch")
}

func TestReduceRemoveNote(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("../tmp")

	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)
	testutils.WriteFile(ctx, "../testutils/fixtures/dnote3.json", "dnote")

	// Execute
	b, err := json.Marshal(&RemoveNoteData{
		BookName: "js",
		NoteUUID: "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
	})
	action := Action{
		Type:      ActionRemoveNote,
		Data:      b,
		Timestamp: 1517629805,
	}
	if err = Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	targetBook := dnote["js"]
	otherBook := dnote["linux"]

	testutils.AssertEqual(t, len(dnote), 2, "number of books mismatch")
	testutils.AssertEqual(t, len(targetBook.Notes), 1, "target book notes length mismatch")
	testutils.AssertEqual(t, targetBook.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	testutils.AssertEqual(t, targetBook.Notes[0].Content, "Booleans have toString()", "remaining note content mismatch")
	testutils.AssertEqual(t, len(otherBook.Notes), 1, "other book notes length mismatch")
	testutils.AssertEqual(t, otherBook.Notes[0].UUID, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", "other book remaining note uuid mismatch")
	testutils.AssertEqual(t, otherBook.Notes[0].Content, "wc -l to count words", "other book remaining note content mismatch")
}

func TestReduceEditNote(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("../tmp")

	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)
	testutils.WriteFile(ctx, "../testutils/fixtures/dnote3.json", "dnote")

	// Execute
	b, err := json.Marshal(&EditNoteData{
		FromBook: "js",
		NoteUUID: "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
		Content:  "updated content",
	})
	action := Action{
		Type:      ActionEditNote,
		Data:      b,
		Timestamp: 1517629805,
	}
	err = Reduce(ctx, action)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	targetBook := dnote["js"]
	otherBook := dnote["linux"]

	testutils.AssertEqual(t, len(dnote), 2, "number of books mismatch")
	testutils.AssertEqual(t, len(targetBook.Notes), 2, "target book notes length mismatch")
	testutils.AssertEqual(t, targetBook.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	testutils.AssertEqual(t, targetBook.Notes[0].Content, "Booleans have toString()", "remaining note content mismatch")
	testutils.AssertEqual(t, targetBook.Notes[1].UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "edited note uuid mismatch")
	testutils.AssertEqual(t, targetBook.Notes[1].Content, "updated content", "edited note content mismatch")
	testutils.AssertEqual(t, targetBook.Notes[1].EditedOn, int64(1517629805), "edited note edited_on mismatch")
	testutils.AssertEqual(t, len(otherBook.Notes), 1, "other book notes length mismatch")
	testutils.AssertEqual(t, otherBook.Notes[0].UUID, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", "other book remaining note uuid mismatch")
	testutils.AssertEqual(t, otherBook.Notes[0].Content, "wc -l to count words", "other book remaining note content mismatch")
}

func TestReduceEditNote_changeBook(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("../tmp")

	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)
	testutils.WriteFile(ctx, "../testutils/fixtures/dnote3.json", "dnote")

	// Execute
	b, err := json.Marshal(&EditNoteData{
		FromBook: "js",
		NoteUUID: "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
		Content:  "updated content",
	})
	action := Action{
		Type:      ActionEditNote,
		Data:      b,
		Timestamp: 1517629805,
	}
	err = Reduce(ctx, action)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	targetBook := dnote["js"]
	otherBook := dnote["linux"]

	testutils.AssertEqual(t, len(dnote), 2, "number of books mismatch")
	testutils.AssertEqual(t, len(targetBook.Notes), 2, "target book notes length mismatch")
	testutils.AssertEqual(t, targetBook.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	testutils.AssertEqual(t, targetBook.Notes[0].Content, "Booleans have toString()", "remaining note content mismatch")
	testutils.AssertEqual(t, targetBook.Notes[1].UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "edited note uuid mismatch")
	testutils.AssertEqual(t, targetBook.Notes[1].Content, "updated content", "edited note content mismatch")
	testutils.AssertEqual(t, targetBook.Notes[1].EditedOn, int64(1517629805), "edited note edited_on mismatch")
	testutils.AssertEqual(t, len(otherBook.Notes), 1, "other book notes length mismatch")
	testutils.AssertEqual(t, otherBook.Notes[0].UUID, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", "other book remaining note uuid mismatch")
	testutils.AssertEqual(t, otherBook.Notes[0].Content, "wc -l to count words", "other book remaining note content mismatch")
}

func TestReduceAddBook(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("../tmp")

	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)
	testutils.WriteFile(ctx, "../testutils/fixtures/dnote4.json", "dnote")

	// Execute
	b, err := json.Marshal(&AddBookData{BookName: "new_book"})
	action := Action{
		Type:      ActionAddBook,
		Data:      b,
		Timestamp: 1517629805,
	}
	if err = Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	newBook := dnote["new_book"]

	testutils.AssertEqual(t, len(dnote), 3, "number of books mismatch")
	testutils.AssertEqual(t, newBook.Name, "new_book", "new book name mismatch")
	testutils.AssertEqual(t, len(newBook.Notes), 0, "new book number of notes mismatch")
}

func TestReduceRemoveBook(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("../tmp")

	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)
	testutils.WriteFile(ctx, "../testutils/fixtures/dnote3.json", "dnote")

	// Execute
	b, err := json.Marshal(&RemoveBookData{BookName: "linux"})
	action := Action{
		Type:      ActionRemoveBook,
		Data:      b,
		Timestamp: 1517629805,
	}
	if err = Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	remainingBook := dnote["js"]

	testutils.AssertEqual(t, len(dnote), 1, "number of books mismatch")
	testutils.AssertEqual(t, remainingBook.Name, "js", "remaining book name mismatch")
	testutils.AssertEqual(t, len(remainingBook.Notes), 2, "remaining book number of notes mismatch")
	testutils.AssertEqual(t, remainingBook.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	testutils.AssertEqual(t, remainingBook.Notes[0].Content, "Booleans have toString()", "remaining note content mismatch")
	testutils.AssertEqual(t, remainingBook.Notes[1].UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "edited note uuid mismatch")
	testutils.AssertEqual(t, remainingBook.Notes[1].Content, "Date object implements mathematical comparisons", "edited note content mismatch")
}
