package add

import (
	"time"

	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var content string

var example = `
 * Write a note in the current book
 dnote new "time is a part of the commit hash"

 * Specify the book name
 dnote new git "time is a part of the commit hash"`

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Incorrect number of argument")
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

	f := cmd.Flags()
	f.StringVarP(&content, "content", "c", "", "The new content for the note")

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		bookName := args[0]

		if content == "" {
			fpath := core.GetDnoteTmpContentPath(ctx)
			input, err := core.GetEditorInput(ctx, fpath)
			if err != nil {
				return errors.Wrap(err, "Failed to get editor input")
			}

			content = input
		}

		if content == "" {
			return errors.New("Empty content")
		}

		ts := time.Now().Unix()
		note := core.NewNote(content, ts)
		err := writeNote(ctx, bookName, note, ts)
		if err != nil {
			return errors.Wrap(err, "Failed to write note")
		}

		log.Infof("content: \"%s\"\n", content)
		log.Infof("added to %s\n", bookName)
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
