package add

import (
	"fmt"
	"time"

	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
 * Write a note in the current book
 dnote new "time is a part of the commit hash"

 * Specify the book name
 dnote new git "time is a part of the commit hash"`

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing argument")
	}

	return nil
}

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <content>",
		Short:   "Add a add note",
		Aliases: []string{"a", "n", "new"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	return cmd
}

func parseArgs(ctx infra.DnoteCtx, args []string) (bookName string, content string, err error) {
	if len(args) == 1 {
		bookName, err = core.GetCurrentBook(ctx)
		if err != nil {
			return
		}

		content = args[0]
	} else {
		bookName = args[0]
		content = args[1]
	}

	return
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		bookName, content, err := parseArgs(ctx, args)
		if err != nil {
			return errors.Wrap(err, "Failed to parse args")
		}

		ts := time.Now().Unix()

		note := core.NewNote(content, ts)
		err = writeNote(ctx, bookName, note, ts)
		if err != nil {
			return errors.Wrap(err, "Failed to write note")
		}

		fmt.Printf("[+] Added to %s\n", bookName)
		return nil
	}
}

func writeNote(ctx infra.DnoteCtx, bookName string, note infra.Note, ts int64) error {
	dnote, err := core.GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}

	var book infra.Book

	book, ok := dnote[bookName]
	if ok {
		notes := append(dnote[bookName].Notes, note)
		dnote[bookName] = core.GetUpdatedBook(dnote[bookName], notes)
	} else {
		book = core.NewBook(bookName)
		book.Notes = []infra.Note{note}
		dnote[bookName] = book

		err = core.LogActionAddBook(ctx, bookName)
		if err != nil {
			return errors.Wrap(err, "Failed to log action")
		}
	}

	err = core.LogActionAddNote(ctx, note.UUID, book.Name, note.Content, ts)
	if err != nil {
		return errors.Wrap(err, "Failed to log action")
	}

	err = core.WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write to dnote file")
	}

	return nil
}
