package add

import (
	"fmt"
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

func newRun(ctx infra.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var bookName string
		var content string

		if len(args) == 1 {
			var err error
			bookName, err = infra.GetCurrentBook(ctx)
			if err != nil {
				return err
			}

			content = args[0]
		} else {
			bookName = args[0]
			content = args[1]
		}

		note := infra.MakeNote(content)
		err := writeNote(ctx, bookName, note)
		if err != nil {
			return errors.Wrap(err, "Failed to write note")
		}

		fmt.Printf("[+] Added to %s\n", bookName)
		return nil
	}
}

func writeNote(ctx infra.DnoteCtx, bookName string, note infra.Note) error {
	dnote, err := infra.GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}

	if _, ok := dnote[bookName]; ok {
		notes := append(dnote[bookName].Notes, note)
		dnote[bookName] = infra.GetUpdatedBook(dnote[bookName], notes)
	} else {
		book := infra.MakeBook()
		book.Notes = []infra.Note{note}
		dnote[bookName] = book
	}

	err = infra.WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write to dnote file")
	}

	return nil
}
