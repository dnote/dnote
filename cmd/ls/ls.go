package ls

import (
	"fmt"
	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
 * List all books
 dnote ls

 * List notes in a book
 dnote ls javascript
 `

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls <book name?>",
		Aliases: []string{"l", "notes"},
		Short:   "List all notes",
		Example: example,
		RunE:    newRun(ctx),
		PreRunE: preRun,
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dnote, err := core.GetDnote(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to read dnote")
		}

		if len(args) == 0 {
			if err := printBooks(dnote); err != nil {
				return errors.Wrap(err, "Failed to print books")
			}

			return nil
		}

		bookName := args[0]
		if err := printNotes(dnote, bookName); err != nil {
			return errors.Wrapf(err, "Failed to print notes for the book %s", bookName)
		}

		return nil
	}
}

func printBooks(dnote infra.Dnote) error {
	for bookName, book := range dnote {
		log.Printf("%s \033[%dm(%d)\033[0m\n", bookName, log.ColorYellow, len(book.Notes))
	}

	return nil
}

func printNotes(dnote infra.Dnote, bookName string) error {
	log.Infof("on book %s\n", bookName)

	book := dnote[bookName]

	for i, note := range book.Notes {
		fmt.Printf("  \033[%dm(%d)\033[0m %s\n", log.ColorYellow, i, note.Content)
	}

	return nil
}
