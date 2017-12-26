package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/test"
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
	ctx := test.InitCtx("./tmp")
	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)

	// Execute
	runDnoteCmd(ctx)

	// Test
	if !utils.FileExists(fmt.Sprintf("%s", ctx.DnoteDir)) {
		t.Errorf("dnote directory was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, infra.DnoteFilename)) {
		t.Errorf("dnote file was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, infra.ConfigFilename)) {
		t.Errorf("config file was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, infra.TimestampFilename)) {
		t.Errorf("timestamp file was not initialized")
	}
	if !utils.FileExists(fmt.Sprintf("%s/%s", ctx.DnoteDir, infra.ActionFilename)) {
		t.Errorf("action file was not initialized")
	}
}

func TestAdd_NewBook(t *testing.T) {
	// Setup
	ctx := test.InitCtx("./tmp")
	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)

	// Execute
	runDnoteCmd(ctx, "add", "js", "foo")

	// Test
	dnote, err := infra.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := infra.ReadActionLog(ctx)
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

	test.AssertEqual(t, bookAction.Type, infra.ActionAddBook, "bookAction type mismatch")
	test.AssertNotEqual(t, bookAction.Data["UUID"], "", "bookAction data UUID mismatch")
	test.AssertNotEqual(t, bookAction.Timestamp, 0, "bookAction timestamp mismatch")
	test.AssertEqual(t, noteAction.Type, infra.ActionAddNote, "noteAction type mismatch")
	test.AssertEqual(t, noteAction.Data["Content"], "foo", "noteAction data name mismatch")
	test.AssertNotEqual(t, noteAction.Data["UUID"], "", "noteAction data UUID mismatch")
	test.AssertNotEqual(t, noteAction.Timestamp, 0, "noteAction timestamp mismatch")
	test.AssertNotEqual(t, book.UUID, "", "Book should have UUID")
	test.AssertEqual(t, len(book.Notes), 1, "Book should have one note")
	test.AssertNotEqual(t, note.UUID, "", "Note should have UUID")
	test.AssertEqual(t, note.Content, "foo", "Note content mismatch")
}

func TestAdd_ExistingBook(t *testing.T) {
	// Setup
	ctx := test.InitCtx("./tmp")
	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)

	// init files by running root command
	runDnoteCmd(ctx)
	test.WriteFile(ctx, "./fixtures/dnote1.json", "dnote")

	// Execute
	runDnoteCmd(ctx, "add", "js", "foo")

	// Test
	dnote, err := infra.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := infra.ReadActionLog(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read actions"))
	}

	book := dnote["js"]
	action := actions[0]

	test.AssertEqual(t, len(actions), 1, "There should be 1 action")
	test.AssertEqual(t, action.Type, infra.ActionAddNote, "action type mismatch")
	test.AssertEqual(t, action.Data["Content"], "foo", "action data name mismatch")
	test.AssertNotEqual(t, action.Data["UUID"], "", "action data UUID mismatch")
	test.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	test.AssertNotEqual(t, book.UUID, "", "Book should have UUID")
	test.AssertEqual(t, len(book.Notes), 2, "Book should have one note")
	test.AssertNotEqual(t, book.Notes[0].UUID, "", "Note should have UUID")
	test.AssertEqual(t, book.Notes[0].Content, "Booleans have toString()", "Note content mismatch")
	test.AssertNotEqual(t, book.Notes[1].UUID, "", "Note should have UUID")
	test.AssertEqual(t, book.Notes[1].Content, "foo", "Note content mismatch")
}

func TestEdit(t *testing.T) {
	// Setup
	ctx := test.InitCtx("./tmp")
	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)

	// init files by running root command
	runDnoteCmd(ctx)
	test.WriteFile(ctx, "./fixtures/dnote2.json", "dnote")

	// Execute
	runDnoteCmd(ctx, "edit", "js", "1", "foo bar")

	// Test
	dnote, err := infra.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := infra.ReadActionLog(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read actions"))
	}

	book := dnote["js"]
	action := actions[0]

	test.AssertEqual(t, len(actions), 1, "There should be 1 action")
	test.AssertEqual(t, action.Type, infra.ActionEditNote, "action type mismatch")
	test.AssertEqual(t, action.Data["Content"], "foo bar", "action data name mismatch")
	test.AssertEqual(t, action.Data["UUID"], "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "action data UUID mismatch")
	test.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	test.AssertNotEqual(t, book.UUID, "", "Book should have UUID")
	test.AssertEqual(t, len(book.Notes), 2, "Book should have one note")
	test.AssertEqual(t, book.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "Note should have UUID")
	test.AssertEqual(t, book.Notes[0].Content, "Booleans have toString()", "Note content mismatch")
	test.AssertEqual(t, book.Notes[1].UUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "Note should have UUID")
	test.AssertEqual(t, book.Notes[1].Content, "foo bar", "Note content mismatch")
}

func TestRemoveNote(t *testing.T) {
	// Setup
	ctx := test.InitCtx("./tmp")
	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)

	// init files by running root command
	runDnoteCmd(ctx)
	test.WriteFile(ctx, "./fixtures/dnote3.json", "dnote")

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

	// Hit return to confirm
	_, err = io.WriteString(stdin, "\n")
	if err != nil {
		panic(errors.Wrap(err, "Failed to write to stdin"))
	}

	err = cmd.Wait()
	if err != nil {
		panic(errors.Wrapf(err, "Failed to run command %s", stderr.String()))
	}

	// Test
	dnote, err := infra.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := infra.ReadActionLog(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read actions"))
	}

	if len(actions) != 1 {
		t.Fatalf("action log length mismatch. got %d", len(actions))
	}

	book := dnote["js"]
	otherBook := dnote["linux"]
	action := actions[0]

	test.AssertEqual(t, len(actions), 1, "There should be 1 action")
	test.AssertEqual(t, action.Type, infra.ActionRemoveNote, "action type mismatch")
	test.AssertEqual(t, action.Data["UUID"], "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "action data UUID mismatch")
	test.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	test.AssertNotEqual(t, book.UUID, "", "Book should have UUID")
	test.AssertEqual(t, len(book.Notes), 1, "Book should have one note")
	test.AssertEqual(t, len(otherBook.Notes), 1, "Other book should have one note")
	test.AssertEqual(t, book.Notes[0].UUID, "43827b9a-c2b0-4c06-a290-97991c896653", "Note should have UUID")
	test.AssertEqual(t, book.Notes[0].Content, "Booleans have toString()", "Note content mismatch")
}

func TestRemoveBook(t *testing.T) {
	// Setup
	ctx := test.InitCtx("./tmp")
	test.SetupTmp(ctx)
	defer test.ClearTmp(ctx)

	// init files by running root command
	runDnoteCmd(ctx)
	test.WriteFile(ctx, "./fixtures/dnote3.json", "dnote")

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

	// Hit return to confirm
	_, err = io.WriteString(stdin, "\n")
	if err != nil {
		panic(errors.Wrap(err, "Failed to write to stdin"))
	}

	err = cmd.Wait()
	if err != nil {
		panic(errors.Wrapf(err, "Failed to run command %s", stderr.String()))
	}

	// Test
	dnote, err := infra.GetDnote(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to get dnote"))
	}
	actions, err := infra.ReadActionLog(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to read actions"))
	}

	if len(actions) != 1 {
		t.Fatalf("action log length mismatch. got %d", len(actions))
	}

	book := dnote["linux"]
	action := actions[0]

	test.AssertEqual(t, len(actions), 1, "There should be 1 action")
	test.AssertEqual(t, action.Type, infra.ActionRemoveBook, "action type mismatch")
	test.AssertEqual(t, action.Data["UUID"], "3e6c9401-833b-485f-bcda-c2525a5dc389", "action data UUID mismatch")
	test.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	test.AssertEqual(t, len(dnote), 1, "There should be 1 book")
	test.AssertEqual(t, book.UUID, "94b829e6-fec8-4e65-95db-7ad2ab0d3a39", "Remaining book uid mismatch")
	test.AssertEqual(t, len(book.Notes), 1, "Remaining book should have one note")
}
