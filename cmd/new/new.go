package new

import (
	"fmt"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/dnote-io/cli/infra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
 * Write a note in the current book
 dnote new "time is a part of the commit hash"

 * Specify the book name
 dnote new git "time is a part of the commit hash"`

var cmd = &cobra.Command{
	Use:     "new <content>",
	Short:   "Add a new note",
	Aliases: []string{"n", "add"},
	Example: example,
	PreRunE: preRun,
	RunE:    run,
}

func init() {
	root.Register(cmd)
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing argument")
	}

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	var bookName string
	var content string

	if len(args) == 1 {
		var err error
		bookName, err = infra.GetCurrentBook()
		if err != nil {
			return err
		}

		content = args[0]
	} else {
		bookName = args[0]
		content = args[1]
	}

	note := infra.MakeNote(content)
	err := writeNote(bookName, note)
	if err != nil {
		return errors.Wrap(err, "Failed to write note")
	}

	fmt.Printf("[+] Added to %s\n", bookName)
	return nil
}

func writeNote(bookName string, note infra.Note) error {
	dnote, err := infra.GetDnote()
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

	err = infra.WriteDnote(dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write to dnote file")
	}

	return nil
}
