package core

import (
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/test"
	"github.com/pkg/errors"
	"testing"
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
	action := Action{
		Type: ActionAddNote,
		Data: map[string]interface{}{
			"content":   "new content",
			"book_uuid": "3e6c9401-833b-485f-bcda-c2525a5dc389",
			"note_uuid": "06896551-8a06-4996-89cc-0d866308b0f6",
		},
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
	action := Action{
		Type: ActionRemoveNote,
		Data: map[string]interface{}{
			"book_uuid": "3e6c9401-833b-485f-bcda-c2525a5dc389",
			"note_uuid": "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
		},
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

	action := Action{
		Type: ActionEditNote,
		Data: map[string]interface{}{
			"book_uuid": "3e6c9401-833b-485f-bcda-c2525a5dc389",
			"note_uuid": "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
			"content":   "updated content",
		},
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
	action := Action{
		Type: ActionAddBook,
		Data: map[string]interface{}{
			"name": "new_book",
			"uuid": "4e6c9401-833b-485f-bcda-c2525aadc389",
		},
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
	test.AssertEqual(t, newBook.UUID, "4e6c9401-833b-485f-bcda-c2525aadc389", "new book uuid mismatch")
	test.AssertEqual(t, len(newBook.Notes), 0, "new book number of notes mismatch")
}

func TestReduceRemoveBook(t *testing.T) {
	// Setup
	ctx := test.InitCtx("../tmp")

	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)
	test.WriteFile(ctx, "../fixtures/dnote3.json", "dnote")

	// Execute
	action := Action{
		Type: ActionRemoveBook,
		Data: map[string]interface{}{
			"uuid": "94b829e6-fec8-4e65-95db-7ad2ab0d3a39",
		},
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
	test.AssertEqual(t, remainingBook.UUID, "3e6c9401-833b-485f-bcda-c2525a5dc389", "remaining book uuid mismatch")
	test.AssertEqual(t, len(remainingBook.Notes), 2, "remaining book number of notes mismatch")
	test.AssertEqual(t, remainingBook.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	test.AssertEqual(t, remainingBook.Notes[0].Content, "Booleans have toString()", "remaining note content mismatch")
	test.AssertEqual(t, remainingBook.Notes[1].UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "edited note uuid mismatch")
	test.AssertEqual(t, remainingBook.Notes[1].Content, "Date object implements mathematical comparisons", "edited note content mismatch")
}
