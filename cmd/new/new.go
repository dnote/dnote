package new

import (
	"errors"
	"fmt"
	"time"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/dnote-io/cli/utils"
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
		bookName, err = utils.GetCurrentBook()
		if err != nil {
			return err
		}

		content = args[0]
	} else {
		bookName = args[0]
		content = args[1]
	}

	note := makeNote(content)
	err := writeNote(bookName, note)
	if err != nil {
		return err
	}

	fmt.Printf("[+] Added to %s\n", bookName)
	return nil
}

func makeNote(content string) utils.Note {
	return utils.Note{
		UID:     utils.GenerateNoteID(),
		Content: content,
		AddedOn: time.Now().Unix(),
	}
}

func writeNote(bookName string, note utils.Note) error {
	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	if _, ok := dnote[bookName]; ok {
		dnote[bookName] = append(dnote[bookName], note)
	} else {
		dnote[bookName] = []utils.Note{note}
	}

	err = utils.WriteDnote(dnote)
	if err != nil {
		return err
	}

	return nil
}
