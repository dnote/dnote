package core

import (
	"encoding/json"
	"testing"

	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/test"
	"github.com/pkg/errors"
)

func TestReduceAddNote(t *testing.T) {
	// Setup
	ctx := test.InitCtx("../tmp")

	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)
	test.WriteFile(ctx, "../fixtures/dnote4.json", "dnote")

	ts := infra.Timestamp{Bookmark: 1517629804}
	if err := WriteTimestamp(ctx, ts); err != nil {
		panic(errors.Wrap(err, "Failed to prepare timestamp"))
	}

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
	if err := Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	ts, err = ReadTimestamp(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read timestamp"))
	}

	book := dnote["js"]
	otherBook := dnote["linux"]

	test.AssertEqual(t, len(book.Notes), 2, "notes length mismatch")
	test.AssertEqual(t, len(otherBook.Notes), 0, "other book notes length mismatch")
	test.AssertEqual(t, book.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "existing note uuid mismatch")
	test.AssertEqual(t, book.Notes[0].Content, "Booleans have toString()", "existing note content mismatch")
	test.AssertEqual(t, book.Notes[1].UUID, "06896551-8a06-4996-89cc-0d866308b0f6", "new note uuid mismatch")
	test.AssertEqual(t, book.Notes[1].Content, "new content", "new note content mismatch")
	test.AssertEqual(t, ts.Bookmark, int64(1517629805), "bookmark was not updated")
}

func TestReduceRemoveNote(t *testing.T) {
	// Setup
	ctx := test.InitCtx("../tmp")

	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)
	test.WriteFile(ctx, "../fixtures/dnote3.json", "dnote")

	ts := infra.Timestamp{Bookmark: 1517629806}
	if err := WriteTimestamp(ctx, ts); err != nil {
		panic(errors.Wrap(err, "Failed to prepare timestamp"))
	}

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
	if err := Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	ts, err = ReadTimestamp(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read timestamp"))
	}

	targetBook := dnote["js"]
	otherBook := dnote["linux"]

	test.AssertEqual(t, len(dnote), 2, "number of books mismatch")
	test.AssertEqual(t, len(targetBook.Notes), 1, "target book notes length mismatch")
	test.AssertEqual(t, targetBook.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	test.AssertEqual(t, targetBook.Notes[0].Content, "Booleans have toString()", "remaining note content mismatch")
	test.AssertEqual(t, len(otherBook.Notes), 1, "other book notes length mismatch")
	test.AssertEqual(t, otherBook.Notes[0].UUID, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", "other book remaining note uuid mismatch")
	test.AssertEqual(t, otherBook.Notes[0].Content, "wc -l to count words", "other book remaining note content mismatch")
	test.AssertEqual(t, ts.Bookmark, int64(1517629806), "bookmark was updated")
}

func TestReduceEditNote(t *testing.T) {
	// Setup
	ctx := test.InitCtx("../tmp")

	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)
	test.WriteFile(ctx, "../fixtures/dnote3.json", "dnote")

	// Execute
	b, err := json.Marshal(&EditNoteData{
		BookName: "js",
		NoteUUID: "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
		Content:  "updated content",
	})
	action := Action{
		Type:      ActionEditNote,
		Data:      b,
		Timestamp: 1517629805,
	}
	if err := Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	targetBook := dnote["js"]
	otherBook := dnote["linux"]

	test.AssertEqual(t, len(dnote), 2, "number of books mismatch")
	test.AssertEqual(t, len(targetBook.Notes), 2, "target book notes length mismatch")
	test.AssertEqual(t, targetBook.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	test.AssertEqual(t, targetBook.Notes[0].Content, "Booleans have toString()", "remaining note content mismatch")
	test.AssertEqual(t, targetBook.Notes[1].UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "edited note uuid mismatch")
	test.AssertEqual(t, targetBook.Notes[1].Content, "updated content", "edited note content mismatch")
	test.AssertEqual(t, len(otherBook.Notes), 1, "other book notes length mismatch")
	test.AssertEqual(t, otherBook.Notes[0].UUID, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", "other book remaining note uuid mismatch")
	test.AssertEqual(t, otherBook.Notes[0].Content, "wc -l to count words", "other book remaining note content mismatch")
}

func TestReduceAddBook(t *testing.T) {
	// Setup
	ctx := test.InitCtx("../tmp")

	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)
	test.WriteFile(ctx, "../fixtures/dnote4.json", "dnote")

	// Execute
	b, err := json.Marshal(&AddBookData{Name: "new_book"})
	action := Action{
		Type:      ActionAddBook,
		Data:      b,
		Timestamp: 1517629805,
	}
	if err := Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	newBook := dnote["new_book"]

	test.AssertEqual(t, len(dnote), 3, "number of books mismatch")
	test.AssertEqual(t, newBook.Name, "new_book", "new book name mismatch")
	test.AssertEqual(t, len(newBook.Notes), 0, "new book number of notes mismatch")
}

func TestReduceRemoveBook(t *testing.T) {
	// Setup
	ctx := test.InitCtx("../tmp")

	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)
	test.WriteFile(ctx, "../fixtures/dnote3.json", "dnote")

	// Execute
	b, err := json.Marshal(&RemoveBookData{Name: "linux"})
	action := Action{
		Type:      ActionRemoveBook,
		Data:      b,
		Timestamp: 1517629805,
	}
	if err := Reduce(ctx, action); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}

	// Test
	dnote, err := GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}

	remainingBook := dnote["js"]

	test.AssertEqual(t, len(dnote), 1, "number of books mismatch")
	test.AssertEqual(t, remainingBook.Name, "js", "remaining book name mismatch")
	test.AssertEqual(t, len(remainingBook.Notes), 2, "remaining book number of notes mismatch")
	test.AssertEqual(t, remainingBook.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	test.AssertEqual(t, remainingBook.Notes[0].Content, "Booleans have toString()", "remaining note content mismatch")
	test.AssertEqual(t, remainingBook.Notes[1].UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "edited note uuid mismatch")
	test.AssertEqual(t, remainingBook.Notes[1].Content, "Date object implements mathematical comparisons", "edited note content mismatch")
}
