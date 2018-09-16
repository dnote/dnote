package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"

	"github.com/dnote/actions"
	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/testutils"
	"github.com/dnote/cli/utils"
)

var binaryName = "test-dnote"

func TestMain(m *testing.M) {
	if err := exec.Command("go", "build", "-o", binaryName).Run(); err != nil {
		log.Print(errors.Wrap(err, "Failed to build a binary").Error())
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func newDnoteCmd(ctx infra.DnoteCtx, arg ...string) (*exec.Cmd, *bytes.Buffer, error) {
	var stderr bytes.Buffer

	binaryPath, err := filepath.Abs(binaryName)
	if err != nil {
		return &exec.Cmd{}, &stderr, errors.Wrap(err, "Failed to get the absolute path to the test binary")
	}

	cmd := exec.Command(binaryPath, arg...)
	cmd.Env = []string{fmt.Sprintf("DNOTE_DIR=%s", ctx.DnoteDir), fmt.Sprintf("DNOTE_HOME_DIR=%s", ctx.HomeDir)}
	cmd.Stderr = &stderr

	return cmd, &stderr, nil
}

func runDnoteCmd(ctx infra.DnoteCtx, arg ...string) {
	cmd, stderr, err := newDnoteCmd(ctx, arg...)
	if err != nil {
		panic(errors.Wrap(err, "Failed to get command").Error())
	}

	if err := cmd.Run(); err != nil {
		panic(errors.Wrapf(err, "Failed to run command %s", stderr.String()))
	}
}

func TestInit(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv("../tmp", "./testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	// Execute
	runDnoteCmd(ctx)

	// Test
	if !utils.FileExists(ctx.DnoteDir) {
		t.Errorf("dnote directory was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, core.ConfigFilename)) {
		t.Errorf("config file was not initialized")
	}

	db := ctx.DB

	var notesTableCount, booksTableCount, actionsTableCount, systemTableCount int
	testutils.MustScan(t, "counting notes",
		db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = ? AND name = ?", "table", "notes"), &notesTableCount)
	testutils.MustScan(t, "counting books",
		db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = ? AND name = ?", "table", "books"), &booksTableCount)
	testutils.MustScan(t, "counting actions",
		db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = ? AND name = ?", "table", "actions"), &actionsTableCount)
	testutils.MustScan(t, "counting system",
		db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = ? AND name = ?", "table", "system"), &systemTableCount)

	testutils.AssertEqual(t, notesTableCount, 1, "notes table count mismatch")
	testutils.AssertEqual(t, booksTableCount, 1, "books table count mismatch")
	testutils.AssertEqual(t, actionsTableCount, 1, "actions table count mismatch")
	testutils.AssertEqual(t, systemTableCount, 1, "system table count mismatch")
}

func TestAddNote_NewBook_ContentFlag(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv("../tmp", "./testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	// Execute
	runDnoteCmd(ctx, "add", "js", "-c", "foo")

	// Test
	db := ctx.DB

	var actionCount, noteCount, bookCount int
	testutils.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)

	testutils.AssertEqualf(t, actionCount, 2, "action count mismatch")
	testutils.AssertEqualf(t, bookCount, 1, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 1, "note count mismatch")

	var jsBookUUID string
	testutils.MustScan(t, "getting js book uuid", db.QueryRow("SELECT uuid FROM books where label = ?", "js"), &jsBookUUID)
	var note infra.Note
	testutils.MustScan(t, "getting note",
		db.QueryRow("SELECT uuid, content, added_on FROM notes where book_uuid = ?", jsBookUUID), &note.UUID, &note.Content, &note.AddedOn)
	var bookAction, noteAction actions.Action
	testutils.MustScan(t, "getting book action",
		db.QueryRow("SELECT data, timestamp FROM actions where type = ?", actions.ActionAddBook), &bookAction.Data, &bookAction.Timestamp)
	testutils.MustScan(t, "getting note action",
		db.QueryRow("SELECT data, timestamp FROM actions where type = ?", actions.ActionAddNote), &noteAction.Data, &noteAction.Timestamp)

	var noteActionData actions.AddNoteDataV1
	var bookActionData actions.AddBookDataV1
	if err := json.Unmarshal(bookAction.Data, &bookActionData); err != nil {
		log.Fatalf("unmarshalling the action data: %s", err)
	}
	if err := json.Unmarshal(noteAction.Data, &noteActionData); err != nil {
		log.Fatalf("unmarshalling the action data: %s", err)
	}

	testutils.AssertNotEqual(t, bookActionData.BookName, "", "bookAction data note_uuid mismatch")
	testutils.AssertNotEqual(t, bookAction.Timestamp, 0, "bookAction timestamp mismatch")
	testutils.AssertEqual(t, noteActionData.Content, "foo", "noteAction data name mismatch")
	testutils.AssertNotEqual(t, noteActionData.NoteUUID, nil, "noteAction data note_uuid mismatch")
	testutils.AssertNotEqual(t, noteActionData.BookName, "", "noteAction data note_uuid mismatch")
	testutils.AssertNotEqual(t, noteAction.Timestamp, 0, "noteAction timestamp mismatch")
	testutils.AssertNotEqual(t, note.UUID, "", "Note should have UUID")
	testutils.AssertEqual(t, note.Content, "foo", "Note content mismatch")
	testutils.AssertNotEqual(t, note.AddedOn, int64(0), "Note added_on mismatch")
}

func TestAddNote_ExistingBook_ContentFlag(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv("../tmp", "./testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	testutils.Setup3(t, ctx)

	// Execute
	runDnoteCmd(ctx, "add", "js", "-c", "foo")

	// Test
	db := ctx.DB

	var actionCount, noteCount, bookCount int
	testutils.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)

	testutils.AssertEqualf(t, actionCount, 1, "action count mismatch")
	testutils.AssertEqualf(t, bookCount, 1, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 2, "note count mismatch")

	var n1, n2 infra.Note
	testutils.MustScan(t, "getting n1",
		db.QueryRow("SELECT uuid, content, added_on FROM notes WHERE book_uuid = ? AND uuid = ?", "js-book-uuid", "43827b9a-c2b0-4c06-a290-97991c896653"), &n1.UUID, &n1.Content, &n1.AddedOn)
	testutils.MustScan(t, "getting n2",
		db.QueryRow("SELECT uuid, content, added_on FROM notes WHERE book_uuid = ? AND content = ?", "js-book-uuid", "foo"), &n2.UUID, &n2.Content, &n2.AddedOn)
	var noteAction actions.Action
	testutils.MustScan(t, "getting note action",
		db.QueryRow("SELECT data, timestamp FROM actions WHERE type = ?", actions.ActionAddNote), &noteAction.Data, &noteAction.Timestamp)

	var noteActionData actions.AddNoteDataV1
	if err := json.Unmarshal(noteAction.Data, &noteActionData); err != nil {
		log.Fatalf("unmarshalling the action data: %s", err)
	}

	testutils.AssertEqual(t, noteActionData.Content, "foo", "action data name mismatch")
	testutils.AssertNotEqual(t, noteActionData.NoteUUID, "", "action data note_uuid mismatch")
	testutils.AssertEqual(t, noteActionData.BookName, "js", "action data book_name mismatch")
	testutils.AssertNotEqual(t, noteAction.Timestamp, 0, "action timestamp mismatch")
	testutils.AssertNotEqual(t, n1.UUID, "", "Note should have UUID")
	testutils.AssertEqual(t, n1.Content, "Booleans have toString()", "Note content mismatch")
	testutils.AssertEqual(t, n1.AddedOn, int64(1515199943), "Note added_on mismatch")
	testutils.AssertNotEqual(t, n2.UUID, "", "Note should have UUID")
	testutils.AssertEqual(t, n2.Content, "foo", "Note content mismatch")
}

func TestEditNote_ContentFlag(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv("../tmp", "./testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	testutils.Setup4(t, ctx)

	// Execute
	runDnoteCmd(ctx, "edit", "js", "2", "-c", "foo bar")

	// Test
	db := ctx.DB

	var actionCount, noteCount, bookCount int
	testutils.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)

	testutils.AssertEqualf(t, actionCount, 1, "action count mismatch")
	testutils.AssertEqualf(t, bookCount, 1, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 2, "note count mismatch")

	var n1, n2 infra.Note
	testutils.MustScan(t, "getting n1",
		db.QueryRow("SELECT uuid, content, added_on FROM notes where book_uuid = ? AND uuid = ?", "js-book-uuid", "43827b9a-c2b0-4c06-a290-97991c896653"), &n1.UUID, &n1.Content, &n1.AddedOn)
	testutils.MustScan(t, "getting n2",
		db.QueryRow("SELECT uuid, content, added_on FROM notes where book_uuid = ? AND uuid = ?", "js-book-uuid", "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f"), &n2.UUID, &n2.Content, &n2.AddedOn)
	var noteAction actions.Action
	testutils.MustScan(t, "getting note action",
		db.QueryRow("SELECT data, type, schema FROM actions where type = ?", actions.ActionEditNote),
		&noteAction.Data, &noteAction.Type, &noteAction.Schema)

	var actionData actions.EditNoteDataV2
	if err := json.Unmarshal(noteAction.Data, &actionData); err != nil {
		log.Fatalf("Failed to unmarshal the action data: %s", err)
	}

	testutils.AssertEqual(t, noteAction.Type, actions.ActionEditNote, "action type mismatch")
	testutils.AssertEqual(t, noteAction.Schema, 2, "action schema mismatch")
	testutils.AssertEqual(t, *actionData.Content, "foo bar", "action data name mismatch")
	testutils.AssertEqual(t, actionData.FromBook, "js", "action data from_book mismatch")
	if actionData.ToBook != nil {
		t.Errorf("action data to_book mismatch. Expected %+v. Got %+v", nil, actionData.ToBook)
	}
	testutils.AssertEqual(t, actionData.NoteUUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "action data note_uuis mismatch")
	testutils.AssertNotEqual(t, noteAction.Timestamp, 0, "action timestamp mismatch")
	testutils.AssertEqual(t, n1.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "Note should have UUID")
	testutils.AssertEqual(t, n1.Content, "Booleans have toString()", "Note content mismatch")
	testutils.AssertEqual(t, n2.UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "Note should have UUID")
	testutils.AssertEqual(t, n2.Content, "foo bar", "Note content mismatch")
	testutils.AssertNotEqual(t, n2.EditedOn, 0, "Note edited_on mismatch")
}

func TestRemoveNote(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv("../tmp", "./testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	testutils.Setup2(t, ctx)

	// Execute
	cmd, stderr, err := newDnoteCmd(ctx, "remove", "js", "1")
	if err != nil {
		panic(errors.Wrap(err, "getting command"))
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(errors.Wrap(err, "getting stdin %s"))
	}
	defer stdin.Close()

	// Start the program
	err = cmd.Start()
	if err != nil {
		panic(errors.Wrap(err, "starting command"))
	}

	// confirm
	_, err = io.WriteString(stdin, "y\n")
	if err != nil {
		panic(errors.Wrap(err, "writing to stdin"))
	}

	err = cmd.Wait()
	if err != nil {
		panic(errors.Wrapf(err, "running command %s", stderr.String()))
	}

	// Test
	db := ctx.DB

	var actionCount, noteCount, bookCount, jsNoteCount, linuxNoteCount int
	testutils.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	testutils.MustScan(t, "counting js notes", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "js-book-uuid"), &jsNoteCount)
	testutils.MustScan(t, "counting linux notes", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "linux-book-uuid"), &linuxNoteCount)

	testutils.AssertEqualf(t, actionCount, 1, "action count mismatch")
	testutils.AssertEqualf(t, bookCount, 2, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 2, "note count mismatch")
	testutils.AssertEqual(t, jsNoteCount, 1, "Book should have one note")
	testutils.AssertEqual(t, linuxNoteCount, 1, "Other book should have one note")

	var b1, b2 infra.Book
	var n1 infra.Note
	testutils.MustScan(t, "getting b1",
		db.QueryRow("SELECT label FROM books WHERE uuid = ?", "js-book-uuid"),
		&b1.Name)
	testutils.MustScan(t, "getting b2",
		db.QueryRow("SELECT label FROM books WHERE uuid = ?", "linux-book-uuid"),
		&b2.Name)
	testutils.MustScan(t, "getting n1",
		db.QueryRow("SELECT uuid, content, added_on FROM notes WHERE book_uuid = ? AND id = ?", "js-book-uuid", 2),
		&n1.UUID, &n1.Content, &n1.AddedOn)

	var noteAction actions.Action
	testutils.MustScan(t, "getting note action",
		db.QueryRow("SELECT type, schema, data FROM actions WHERE type = ?", actions.ActionRemoveNote), &noteAction.Type, &noteAction.Schema, &noteAction.Data)

	var actionData actions.RemoveNoteDataV1
	if err := json.Unmarshal(noteAction.Data, &actionData); err != nil {
		log.Fatalf("unmarshalling the action data: %s", err)
	}

	testutils.AssertEqual(t, b1.Name, "js", "b1 label mismatch")
	testutils.AssertEqual(t, b2.Name, "linux", "b2 label mismatch")
	testutils.AssertEqual(t, noteAction.Schema, 1, "action schema mismatch")
	testutils.AssertEqual(t, noteAction.Type, actions.ActionRemoveNote, "action type mismatch")
	testutils.AssertEqual(t, actionData.NoteUUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "action data note_uuid mismatch")
	testutils.AssertEqual(t, actionData.BookName, "js", "action data book_name mismatch")
	testutils.AssertNotEqual(t, noteAction.Timestamp, 0, "action timestamp mismatch")
	testutils.AssertEqual(t, n1.UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "Note should have UUID")
	testutils.AssertEqual(t, n1.Content, "Booleans have toString()", "Note content mismatch")
}

func TestRemoveBook(t *testing.T) {
	// Set up
	ctx := testutils.InitEnv("../tmp", "./testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	testutils.Setup2(t, ctx)

	// Execute
	cmd, stderr, err := newDnoteCmd(ctx, "remove", "-b", "js")
	if err != nil {
		panic(errors.Wrap(err, "getting command"))
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(errors.Wrap(err, "getting stdin %s"))
	}
	defer stdin.Close()

	// Start the program
	err = cmd.Start()
	if err != nil {
		panic(errors.Wrap(err, "starting command"))
	}

	// confirm
	_, err = io.WriteString(stdin, "y\n")
	if err != nil {
		panic(errors.Wrap(err, "writing to stdin"))
	}

	err = cmd.Wait()
	if err != nil {
		panic(errors.Wrapf(err, "running command %s", stderr.String()))
	}

	// Test
	db := ctx.DB

	var actionCount, noteCount, bookCount, jsNoteCount, linuxNoteCount int
	testutils.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	testutils.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	testutils.MustScan(t, "counting notes", db.QueryRow("SELECT count(*) FROM notes"), &noteCount)
	testutils.MustScan(t, "counting js notes", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "js-book-uuid"), &jsNoteCount)
	testutils.MustScan(t, "counting linux notes", db.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ?", "linux-book-uuid"), &linuxNoteCount)

	testutils.AssertEqualf(t, actionCount, 1, "action count mismatch")
	testutils.AssertEqualf(t, bookCount, 1, "book count mismatch")
	testutils.AssertEqualf(t, noteCount, 1, "note count mismatch")
	testutils.AssertEqual(t, jsNoteCount, 0, "some notes in book were not deleted")
	testutils.AssertEqual(t, linuxNoteCount, 1, "Other book should have one note")

	var b1 infra.Book
	testutils.MustScan(t, "getting b1",
		db.QueryRow("SELECT label FROM books WHERE uuid = ?", "linux-book-uuid"),
		&b1.Name)

	var action actions.Action
	testutils.MustScan(t, "getting an action",
		db.QueryRow("SELECT type, schema, data FROM actions WHERE type = ?", actions.ActionRemoveBook), &action.Type, &action.Schema, &action.Data)

	var actionData actions.RemoveBookDataV1
	if err = json.Unmarshal(action.Data, &actionData); err != nil {
		log.Fatalf("unmarshalling the action data: %s", err)
	}

	testutils.AssertEqual(t, action.Type, actions.ActionRemoveBook, "action type mismatch")
	testutils.AssertEqual(t, actionData.BookName, "js", "action data name mismatch")
	testutils.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	testutils.AssertEqual(t, b1.Name, "linux", "Remaining book name mismatch")
}
