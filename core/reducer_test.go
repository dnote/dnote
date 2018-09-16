package core

import (
	"encoding/json"
	"testing"

	"github.com/dnote/actions"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/testutils"
	"github.com/pkg/errors"
)

func TestReduceAddNote(t *testing.T) {
	// Setup
	ctx := testutils.InitEnv("../tmp", "../testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	testutils.Setup1(t, ctx)

	// Execute
	b, err := json.Marshal(&actions.AddNoteDataV1{
		Content:  "new content",
		BookName: "js",
		NoteUUID: "06896551-8a06-4996-89cc-0d866308b0f6",
	})
	action := actions.Action{
		Type:      actions.ActionAddNote,
		Data:      b,
		Timestamp: 1517629805,
	}

	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		panic(errors.Wrap(err, "beginning a transaction"))
	}
	if err = Reduce(ctx, tx, action); err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "processing action"))
	}
	tx.Commit()

	// Test
	var noteCount, jsNoteCount, linuxNoteCount int
	testutils.MustScan(t, "counting note", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	testutils.MustScan(t, "counting js note", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "js-book-uuid"), &jsNoteCount)
	testutils.MustScan(t, "counting linux note", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "linux-book-uuid"), &linuxNoteCount)
	testutils.AssertEqual(t, noteCount, 2, "notes length mismatch")
	testutils.AssertEqual(t, jsNoteCount, 2, "js notes length mismatch")
	testutils.AssertEqual(t, linuxNoteCount, 0, "linux notes length mismatch")

	var existingNote, newNote infra.Note
	testutils.MustScan(t, "scanning existing note", db.QueryRow("SELECT uuid, content FROM notes WHERE uuid = ?", "43827b9a-c2b0-4c06-a290-97991c896653"), &existingNote.UUID, &existingNote.Content)
	testutils.MustScan(t, "scanning new note", db.QueryRow("SELECT uuid, content, added_on FROM notes WHERE uuid = ?", "06896551-8a06-4996-89cc-0d866308b0f6"), &newNote.UUID, &newNote.Content, &newNote.AddedOn)

	testutils.AssertEqual(t, existingNote.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "existing note uuid mismatch")
	testutils.AssertEqual(t, existingNote.Content, "Booleans have toString()", "existing note content mismatch")
	testutils.AssertEqual(t, newNote.UUID, "06896551-8a06-4996-89cc-0d866308b0f6", "new note uuid mismatch")
	testutils.AssertEqual(t, newNote.Content, "new content", "new note content mismatch")
	testutils.AssertEqual(t, newNote.AddedOn, int64(1517629805), "new note added_on mismatch")
}

func TestReduceRemoveNote(t *testing.T) {
	// Setup
	ctx := testutils.InitEnv("../tmp", "../testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	testutils.Setup2(t, ctx)

	// Execute
	b, err := json.Marshal(&actions.RemoveNoteDataV1{
		BookName: "js",
		NoteUUID: "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
	})
	action := actions.Action{
		Type:      actions.ActionRemoveNote,
		Data:      b,
		Timestamp: 1517629805,
	}

	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		panic(errors.Wrap(err, "beginning a transaction"))
	}
	if err = Reduce(ctx, tx, action); err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "processing action"))
	}
	tx.Commit()

	// Test
	var bookCount, noteCount, jsNoteCount, linuxNoteCount int
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting note", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	testutils.MustScan(t, "counting js note", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "js-book-uuid"), &jsNoteCount)
	testutils.MustScan(t, "counting linux note", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "linux-book-uuid"), &linuxNoteCount)

	var n1, n2 infra.Note
	testutils.MustScan(t, "scanning note 1", db.QueryRow("SELECT uuid, content FROM notes WHERE uuid = ?", "43827b9a-c2b0-4c06-a290-97991c896653"), &n1.UUID, &n1.Content)
	testutils.MustScan(t, "scanning note 2", db.QueryRow("SELECT uuid, content FROM notes WHERE uuid = ?", "3e065d55-6d47-42f2-a6bf-f5844130b2d2"), &n2.UUID, &n2.Content)

	testutils.AssertEqual(t, bookCount, 2, "number of books mismatch")
	testutils.AssertEqual(t, jsNoteCount, 1, "target book notes length mismatch")
	testutils.AssertEqual(t, linuxNoteCount, 1, "other book notes length mismatch")
	testutils.AssertEqual(t, n1.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	testutils.AssertEqual(t, n1.Content, "Booleans have toString()", "remaining note content mismatch")
	testutils.AssertEqual(t, n2.UUID, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", "other book remaining note uuid mismatch")
	testutils.AssertEqual(t, n2.Content, "wc -l to count words", "other book remaining note content mismatch")
}

func TestReduceEditNote(t *testing.T) {
	testCases := []struct {
		data                   string
		expectedNoteUUID       string
		expectedNoteBookUUID   string
		expectedNoteContent    string
		expectedNoteAddedOn    int64
		expectedNoteEditedOn   int64
		expectedNotePublic     bool
		expectedJsNoteCount    int
		expectedLinuxNoteCount int
	}{
		{
			data:                   `{"note_uuid": "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "from_book": "js", "content": "updated content"}`,
			expectedNoteUUID:       "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
			expectedNoteBookUUID:   "js-book-uuid",
			expectedNoteContent:    "updated content",
			expectedNoteAddedOn:    int64(1515199951),
			expectedNoteEditedOn:   int64(1517629805),
			expectedNotePublic:     false,
			expectedJsNoteCount:    2,
			expectedLinuxNoteCount: 1,
		},
		{
			data:                   `{"note_uuid": "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "from_book": "js", "public": true}`,
			expectedNoteUUID:       "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
			expectedNoteBookUUID:   "js-book-uuid",
			expectedNoteContent:    "Date object implements mathematical comparisons",
			expectedNoteAddedOn:    int64(1515199951),
			expectedNoteEditedOn:   int64(1517629805),
			expectedNotePublic:     true,
			expectedJsNoteCount:    2,
			expectedLinuxNoteCount: 1,
		},
		{
			data:                   `{"note_uuid": "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "from_book": "js", "to_book": "linux", "content": "updated content"}`,
			expectedNoteUUID:       "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f",
			expectedNoteBookUUID:   "linux-book-uuid",
			expectedNoteContent:    "updated content",
			expectedNoteAddedOn:    int64(1515199951),
			expectedNoteEditedOn:   int64(1517629805),
			expectedNotePublic:     false,
			expectedJsNoteCount:    1,
			expectedLinuxNoteCount: 2,
		}}

	for _, tc := range testCases {
		// Setup
		func() {
			ctx := testutils.InitEnv("../tmp", "../testutils/fixtures/schema.sql")
			defer testutils.TeardownEnv(ctx)

			testutils.Setup2(t, ctx)
			db := ctx.DB

			// Execute
			action := actions.Action{
				Type:      actions.ActionEditNote,
				Data:      json.RawMessage(tc.data),
				Schema:    2,
				Timestamp: 1517629805,
			}

			tx, err := db.Begin()
			if err != nil {
				panic(errors.Wrap(err, "beginning a transaction"))
			}
			err = Reduce(ctx, tx, action)
			if err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "Failed to process action"))
			}

			tx.Commit()

			// Test
			var bookCount, noteCount, jsNoteCount, linuxNoteCount int
			testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
			testutils.MustScan(t, "counting note", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
			testutils.MustScan(t, "counting js note", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "js-book-uuid"), &jsNoteCount)
			testutils.MustScan(t, "counting linux note", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "linux-book-uuid"), &linuxNoteCount)

			var n1, n2, n3 infra.Note
			testutils.MustScan(t, "scanning note 1", db.QueryRow("SELECT uuid, content FROM notes WHERE uuid = ?", "43827b9a-c2b0-4c06-a290-97991c896653"), &n1.UUID, &n1.Content)
			testutils.MustScan(t, "scanning note 2", db.QueryRow("SELECT uuid, content FROM notes WHERE uuid = ?", "3e065d55-6d47-42f2-a6bf-f5844130b2d2"), &n2.UUID, &n2.Content)
			testutils.MustScan(t, "scanning note 2", db.QueryRow("SELECT uuid, content, added_on, edited_on, public FROM notes WHERE uuid = ?", "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f"), &n3.UUID, &n3.Content, &n3.AddedOn, &n3.EditedOn, &n3.Public)

			testutils.AssertEqual(t, bookCount, 2, "number of books mismatch")
			testutils.AssertEqual(t, noteCount, 3, "number of notes mismatch")
			testutils.AssertEqual(t, jsNoteCount, tc.expectedJsNoteCount, "js book notes length mismatch")
			testutils.AssertEqual(t, linuxNoteCount, tc.expectedLinuxNoteCount, "linux book notes length mismatch")

			testutils.AssertEqual(t, n1.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "n1 mismatch")
			testutils.AssertEqual(t, n1.Content, "Booleans have toString()", "n1 content mismatch")
			testutils.AssertEqual(t, n2.UUID, "3e065d55-6d47-42f2-a6bf-f5844130b2d2", "n2 uuid mismatch")
			testutils.AssertEqual(t, n2.Content, "wc -l to count words", "n2 content mismatch")
			testutils.AssertEqual(t, n3.UUID, tc.expectedNoteUUID, "edited note uuid mismatch")
			testutils.AssertEqual(t, n3.Content, tc.expectedNoteContent, "edited note content mismatch")
			testutils.AssertEqual(t, n3.AddedOn, tc.expectedNoteAddedOn, "edited note added_on mismatch")
			testutils.AssertEqual(t, n3.EditedOn, tc.expectedNoteEditedOn, "edited note edited_on mismatch")
			testutils.AssertEqual(t, n3.Public, tc.expectedNotePublic, "edited note public mismatch")
		}()
	}
}

func TestReduceAddBook(t *testing.T) {
	// Setup
	ctx := testutils.InitEnv("../tmp", "../testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	testutils.Setup1(t, ctx)

	// Execute
	b, err := json.Marshal(&actions.AddBookDataV1{BookName: "new_book"})
	action := actions.Action{
		Type:      actions.ActionAddBook,
		Data:      b,
		Timestamp: 1517629805,
	}
	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		panic(errors.Wrap(err, "beginning a transaction"))
	}
	if err = Reduce(ctx, tx, action); err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}
	tx.Commit()

	// Test
	var bookCount, newBookNoteCount int
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting note in the new book", db.QueryRow("SELECT count(*) FROM notes INNER JOIN books ON books.uuid = notes.book_uuid WHERE books.label = ?", "new_book"), &newBookNoteCount)

	testutils.AssertEqual(t, bookCount, 3, "number of books mismatch")
	testutils.AssertEqual(t, newBookNoteCount, 0, "new book number of notes mismatch")
}

func TestReduceRemoveBook(t *testing.T) {
	// Setup
	ctx := testutils.InitEnv("../tmp", "../testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	testutils.Setup2(t, ctx)

	// Execute
	b, err := json.Marshal(&actions.RemoveBookDataV1{BookName: "linux"})
	action := actions.Action{
		Type:      actions.ActionRemoveBook,
		Data:      b,
		Timestamp: 1517629805,
	}

	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		panic(errors.Wrap(err, "beginning a transaction"))
	}
	if err = Reduce(ctx, tx, action); err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "Failed to process action"))
	}
	tx.Commit()

	// Test
	var bookCount, noteCount, jsNoteCount, linuxNoteCount int
	var jsBookLabel string
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting note", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	testutils.MustScan(t, "counting js note", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "js-book-uuid"), &jsNoteCount)
	testutils.MustScan(t, "counting linux note", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "linux-book-uuid"), &linuxNoteCount)
	testutils.MustScan(t, "scanning book", db.QueryRow("SELECT label FROM books WHERE uuid = ?", "js-book-uuid"), &jsBookLabel)

	var n1, n2 infra.Note
	testutils.MustScan(t, "scanning note 1", db.QueryRow("SELECT uuid, content FROM notes WHERE uuid = ?", "43827b9a-c2b0-4c06-a290-97991c896653"), &n1.UUID, &n1.Content)
	testutils.MustScan(t, "scanning note 2", db.QueryRow("SELECT uuid, content FROM notes WHERE uuid = ?", "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f"), &n2.UUID, &n2.Content)

	testutils.AssertEqual(t, bookCount, 1, "number of books mismatch")
	testutils.AssertEqual(t, noteCount, 2, "number of notes mismatch")
	testutils.AssertEqual(t, jsNoteCount, 2, "js note count mismatch")
	testutils.AssertEqual(t, linuxNoteCount, 0, "linux note count mismatch")
	testutils.AssertEqual(t, jsBookLabel, "js", "remaining book name mismatch")

	testutils.AssertEqual(t, n1.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "remaining note uuid mismatch")
	testutils.AssertEqual(t, n1.Content, "Booleans have toString()", "remaining note content mismatch")
	testutils.AssertEqual(t, n2.UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "edited note uuid mismatch")
	testutils.AssertEqual(t, n2.Content, "Date object implements mathematical comparisons", "edited note content mismatch")
}
