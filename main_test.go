package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/test"
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

	dnoteDirPath, err := filepath.Abs("./tmp/.dnote")
	if err != nil {
		return &exec.Cmd{}, &stderr, errors.Wrap(err, "Failed to get the absolute path to the dnote dir")
	}

	cmd := exec.Command(binaryPath, arg...)
	cmd.Env = []string{fmt.Sprintf("DNOTE_DIR=%s", dnoteDirPath)}
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

	book := dnote["js"]
	note := book.Notes[0]

	test.AssertNotEqual(t, book.UID, "", "Book should have UID")
	test.AssertEqual(t, len(book.Notes), 1, "Book should have one note")
	test.AssertNotEqual(t, note.UID, "", "Note should have UID")
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

	book := dnote["js"]

	test.AssertNotEqual(t, book.UID, "", "Book should have UID")
	test.AssertEqual(t, len(book.Notes), 2, "Book should have one note")
	test.AssertNotEqual(t, book.Notes[0].UID, "", "Note should have UID")
	test.AssertEqual(t, book.Notes[0].Content, "Booleans have toString()", "Note content mismatch")
	test.AssertNotEqual(t, book.Notes[1].UID, "", "Note should have UID")
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

	book := dnote["js"]

	test.AssertNotEqual(t, book.UID, "", "Book should have UID")
	test.AssertEqual(t, len(book.Notes), 2, "Book should have one note")
	test.AssertEqual(t, book.Notes[0].UID, "hy07v63d", "Note should have UID")
	test.AssertEqual(t, book.Notes[0].Content, "Booleans have toString()", "Note content mismatch")
	test.AssertEqual(t, book.Notes[1].UID, "pzuz03c4", "Note should have UID")
	test.AssertEqual(t, book.Notes[1].Content, "foo bar", "Note content mismatch")
}
