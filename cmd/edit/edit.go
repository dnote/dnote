package edit

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dnote-io/cli/infra"
	"github.com/spf13/cobra"
)

var example = `
  * Edit the note by index in the current book
  dnote edit 3 "new content"

  * Edit the note by index in a certain book
  dnote edit JS 3 "new content"`

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Edit a note or a book",
		Aliases: []string{"e"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Missing argument")
	}

	return nil
}

func newRun(ctx infra.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dnote, err := infra.GetDnote(ctx)
		if err != nil {
			return err
		}

		var targetBook string
		var index int
		var content string

		if len(args) == 2 {
			targetBook, err = infra.GetCurrentBook(ctx)
			if err != nil {
				return err
			}
			index, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			content = args[1]
		} else if len(args) == 3 {
			targetBook = args[0]
			index, err = strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			content = args[2]
		}

		for i, note := range dnote[targetBook].Notes {
			if i == index {
				note.Content = content
				dnote[targetBook].Notes[i] = note

				err := infra.WriteDnote(ctx, dnote)
				fmt.Printf("Edited Note : %d \n", index)
				return err
			}
		}

		// If loop finishes without returning, note did not exist
		fmt.Println("Error : The note with that index is not found.")
		return nil
	}
}
