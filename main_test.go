package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/testutils"
	"github.com/dnote-io/cli/utils"
)

var binaryName = "test-dnote"

func TestMain(m *testing.M) {
	if err := exec.Command("go", "build", "-o", binaryName).Run(); err != nil {
		log.Printf(errors.Wrap(err, "Failed to build a binary").Error())
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
	// Setup
	ctx := testutils.InitCtx("./tmp")
	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)

	// Execute
	runDnoteCmd(ctx)

	// Test
	if !utils.FileExists(fmt.Sprintf("%s", ctx.DnoteDir)) {
		t.Errorf("dnote directory was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, core.DnoteFilename)) {
		t.Errorf("dnote file was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, core.ConfigFilename)) {
		t.Errorf("config file was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, core.TimestampFilename)) {
		t.Errorf("timestamp file was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, core.ActionFilename)) {
		t.Errorf("action file was not initialized")
	}
}

func TestAdd_NewBook(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("./tmp")
	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)

	// Execute
	runDnoteCmd(ctx, "add", "js", "foo")

	// Test
	dnote, err := core.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := core.ReadActionLog(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read actions"))
	}

	if len(actions) != 2 {
		t.Fatalf("action log length mismatch. got %d", len(actions))
	}

	book := dnote["js"]
	note := book.Notes[0]
	bookAction := actions[0]
	noteAction := actions[1]

	var noteActionData core.AddNoteData
	var bookActionData core.AddBookData
	err = json.Unmarshal(bookAction.Data, &bookActionData)
	if err != nil {
		log.Fatalln("Failed to unmarshal the action data: %s", err)
	}
	err = json.Unmarshal(noteAction.Data, &noteActionData)
	if err != nil {
		log.Fatalln("Failed to unmarshal the action data: %s", err)
	}

	testutils.AssertEqual(t, bookAction.Type, core.ActionAddBook, "bookAction type mismatch")
	testutils.AssertNotEqual(t, bookActionData.BookName, "", "bookAction data note_uuid mismatch")
	testutils.AssertNotEqual(t, bookAction.Timestamp, 0, "bookAction timestamp mismatch")
	testutils.AssertEqual(t, noteAction.Type, core.ActionAddNote, "noteAction type mismatch")
	testutils.AssertEqual(t, noteActionData.Content, "foo", "noteAction data name mismatch")
	testutils.AssertNotEqual(t, noteActionData.NoteUUID, nil, "noteAction data note_uuid mismatch")
	testutils.AssertNotEqual(t, noteActionData.BookName, "", "noteAction data note_uuid mismatch")
	testutils.AssertNotEqual(t, noteAction.Timestamp, 0, "noteAction timestamp mismatch")
	testutils.AssertEqual(t, len(book.Notes), 1, "Book should have one note")
	testutils.AssertNotEqual(t, note.UUID, "", "Note should have UUID")
	testutils.AssertEqual(t, note.Content, "foo", "Note content mismatch")
	testutils.AssertNotEqual(t, note.AddedOn, int64(0), "Note added_on mismatch")
}

func TestAdd_ExistingBook(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("./tmp")
	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)

	// init files by running root command
	runDnoteCmd(ctx)
	testutils.WriteFile(ctx, "./testutils/fixtures/dnote1.json", "dnote")

	// Execute
	runDnoteCmd(ctx, "add", "js", "foo")

	// Test
	dnote, err := core.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := core.ReadActionLog(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read actions"))
	}

	book := dnote["js"]
	action := actions[0]

	var actionData core.AddNoteData
	err = json.Unmarshal(action.Data, &actionData)
	if err != nil {
		log.Fatalln("Failed to unmarshal the action data: %s", err)
	}

	testutils.AssertEqual(t, len(actions), 1, "There should be 1 action")
	testutils.AssertEqual(t, action.Type, core.ActionAddNote, "action type mismatch")
	testutils.AssertEqual(t, actionData.Content, "foo", "action data name mismatch")
	testutils.AssertNotEqual(t, actionData.NoteUUID, "", "action data note_uuid mismatch")
	testutils.AssertEqual(t, actionData.BookName, "js", "action data book_name mismatch")
	testutils.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	testutils.AssertEqual(t, len(book.Notes), 2, "Book should have one note")
	testutils.AssertNotEqual(t, book.Notes[0].UUID, "", "Note should have UUID")
	testutils.AssertEqual(t, book.Notes[0].Content, "Booleans have toString()", "Note content mismatch")
	testutils.AssertNotEqual(t, book.Notes[1].UUID, "", "Note should have UUID")
	testutils.AssertEqual(t, book.Notes[1].Content, "foo", "Note content mismatch")
}

func TestEdit(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("./tmp")
	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)

	// init files by running root command
	runDnoteCmd(ctx)
	testutils.WriteFile(ctx, "./testutils/fixtures/dnote2.json", "dnote")

	// Execute
	runDnoteCmd(ctx, "edit", "js", "1", "foo bar")

	// Test
	dnote, err := core.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := core.ReadActionLog(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read actions"))
	}

	book := dnote["js"]
	action := actions[0]

	var actionData core.EditNoteData
	err = json.Unmarshal(action.Data, &actionData)
	if err != nil {
		log.Fatalln("Failed to unmarshal the action data: %s", err)
	}

	testutils.AssertEqual(t, len(actions), 1, "There should be 1 action")
	testutils.AssertEqual(t, action.Type, core.ActionEditNote, "action type mismatch")
	testutils.AssertEqual(t, actionData.Content, "foo bar", "action data name mismatch")
	testutils.AssertEqual(t, actionData.BookName, "js", "action data book_name mismatch")
	testutils.AssertEqual(t, actionData.NoteUUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "action data note_uuis mismatch")
	testutils.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	testutils.AssertEqual(t, len(book.Notes), 2, "Book should have one note")
	testutils.AssertEqual(t, book.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "Note should have UUID")
	testutils.AssertEqual(t, book.Notes[0].Content, "Booleans have toString()", "Note content mismatch")
	testutils.AssertEqual(t, book.Notes[1].UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "Note should have UUID")
	testutils.AssertEqual(t, book.Notes[1].Content, "foo bar", "Note content mismatch")
	testutils.AssertNotEqual(t, book.Notes[1].EditedOn, int64(0), "Note edited_on mismatch")
}

func TestRemoveNote(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("./tmp")
	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)

	// init files by running root command
	runDnoteCmd(ctx)
	testutils.WriteFile(ctx, "./testutils/fixtures/dnote3.json", "dnote")

	// Execute
	cmd, stderr, err := newDnoteCmd(ctx, "remove", "js", "1")
	if err != nil {
		panic(errors.Wrap(err, "Failed to get command"))
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(errors.Wrap(err, "Failed to get stdin %s"))
	}
	defer stdin.Close()

	// Start the program
	err = cmd.Start()
	if err != nil {
		panic(errors.Wrap(err, "Failed to start command"))
	}

	// confirm
	_, err = io.WriteString(stdin, "y\n")
	if err != nil {
		panic(errors.Wrap(err, "Failed to write to stdin"))
	}

	err = cmd.Wait()
	if err != nil {
		panic(errors.Wrapf(err, "Failed to run command %s", stderr.String()))
	}

	// Test
	dnote, err := core.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := core.ReadActionLog(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read actions"))
	}

	if len(actions) != 1 {
		t.Fatalf("action log length mismatch. got %d", len(actions))
	}

	book := dnote["js"]
	otherBook := dnote["linux"]
	action := actions[0]

	var actionData core.RemoveNoteData
	err = json.Unmarshal(action.Data, &actionData)
	if err != nil {
		log.Fatalln("Failed to unmarshal the action data: %s", err)
	}

	testutils.AssertEqual(t, len(actions), 1, "There should be 1 action")
	testutils.AssertEqual(t, action.Type, core.ActionRemoveNote, "action type mismatch")
	testutils.AssertEqual(t, actionData.NoteUUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "action data note_uuid mismatch")
	testutils.AssertEqual(t, actionData.BookName, "js", "action data book_name mismatch")
	testutils.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	testutils.AssertEqual(t, len(book.Notes), 1, "Book should have one note")
	testutils.AssertEqual(t, len(otherBook.Notes), 1, "Other book should have one note")
	testutils.AssertEqual(t, book.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "Note should have UUID")
	testutils.AssertEqual(t, book.Notes[0].Content, "Booleans have toString()", "Note content mismatch")
}

func TestRemoveBook(t *testing.T) {
	// Setup
	ctx := testutils.InitCtx("./tmp")
	testutils.SetupTmp(ctx)
	defer testutils.ClearTmp(ctx)

	// init files by running root command
	runDnoteCmd(ctx)
	testutils.WriteFile(ctx, "./testutils/fixtures/dnote3.json", "dnote")

	// Execute
	cmd, stderr, err := newDnoteCmd(ctx, "remove", "-b", "js")
	if err != nil {
		panic(errors.Wrap(err, "Failed to get command"))
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(errors.Wrap(err, "Failed to get stdin %s"))
	}
	defer stdin.Close()

	// Start the program
	err = cmd.Start()
	if err != nil {
		panic(errors.Wrap(err, "Failed to start command"))
	}

	// confirm
	_, err = io.WriteString(stdin, "y\n")
	if err != nil {
		panic(errors.Wrap(err, "Failed to write to stdin"))
	}

	err = cmd.Wait()
	if err != nil {
		panic(errors.Wrapf(err, "Failed to run command %s", stderr.String()))
	}

	// Test
	dnote, err := core.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := core.ReadActionLog(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read actions"))
	}

	if len(actions) != 1 {
		t.Fatalf("action log length mismatch. got %d", len(actions))
	}

	book := dnote["linux"]
	action := actions[0]

	var actionData core.RemoveBookData
	err = json.Unmarshal(action.Data, &actionData)
	if err != nil {
		log.Fatalln("Failed to unmarshal the action data: %s", err)
	}

	testutils.AssertEqual(t, len(actions), 1, "There should be 1 action")
	testutils.AssertEqual(t, action.Type, core.ActionRemoveBook, "action type mismatch")
	testutils.AssertEqual(t, actionData.BookName, "js", "action data name mismatch")
	testutils.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	testutils.AssertEqual(t, len(dnote), 1, "There should be 1 book")
	testutils.AssertEqual(t, book.Name, "linux", "Remaining book name mismatch")
	testutils.AssertEqual(t, len(book.Notes), 1, "Remaining book should have one note")
}
