package edit

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/dnote-io/cli/utils"
	"github.com/spf13/cobra"
)

var example = `
  * Edit the note by index in the current book
  dnote edit 3 "new content"

  * Edit the note by index in a certain book
  dnote edit JS 3 "new content"`

var cmd = &cobra.Command{
	Use:     "edit",
	Short:   "Edit a note or a book",
	Aliases: []string{"e"},
	Example: example,
	PreRunE: preRun,
	RunE:    run,
}

func init() {
	root.Register(cmd)
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Missing argument")
	}

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	var targetBook string
	var index int
	var content string

	if len(args) == 2 {
		targetBook, err = utils.GetCurrentBook()
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

	for i, note := range dnote[targetBook] {
		if i == index {
			note.Content = content
			dnote[targetBook][i] = note

			err := utils.WriteDnote(dnote)
			fmt.Printf("Edited Note : %d \n", index)
			return err
		}
	}

	// If loop finishes without returning, note did not exist
	fmt.Println("Error : The note with that index is not found.")
	return nil
}
