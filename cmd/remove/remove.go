package remove

import (
	"fmt"
	"strconv"

	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var targetBookName string

var example = `
  * Delete a note by its index from a book
  dnote delete js 2

  * Delete a book
  dnote delete -b js`

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a note or a book",
		Aliases: []string{"rm", "d", "delete"},
		Example: example,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&targetBookName, "book", "b", "", "The book name to delete")

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		if targetBookName != "" {
			book(ctx, targetBookName)
		} else {
			if len(args) < 2 {
				return errors.New("Missing argument")
			}

			targetBook := args[0]
			noteIndex, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			err = note(ctx, noteIndex, targetBook)
			if err != nil {
				return errors.Wrap(err, "Failed to delete the note")
			}
		}

		return nil
	}
}

// note deletes the note in a certain index.
func note(ctx infra.DnoteCtx, index int, bookName string) error {
	dnote, err := core.GetDnote(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote")
	}

	book, exists := dnote[bookName]
	if !exists {
		return errors.Errorf("Book with the name '%s' does not exist", bookName)
	}
	notes := book.Notes

	if len(notes)-1 < index {
		fmt.Println("Error : The note with that index is not found.")
		return nil
	}

	content := notes[index].Content
	fmt.Printf("Deleting note: %s\n", content)

	ok, err := utils.AskConfirmation("Are you sure?")
	if err != nil {
		return errors.Wrap(err, "Failed to get confirmation")
	}
	if !ok {
		return nil
	}

	dnote[bookName] = core.GetUpdatedBook(dnote[bookName], append(notes[:index], notes[index+1:]...))

	note := notes[index]
	err = core.LogActionRemoveNote(ctx, note.UUID, book.Name)
	if err != nil {
		return errors.Wrap(err, "Failed to log action")
	}

	err = core.WriteDnote(ctx, dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to write dnote")
	}

	fmt.Printf("Deleted!\n")
	return nil
}

// book deletes a book with the given name
func book(ctx infra.DnoteCtx, bookName string) error {
	ok, err := utils.AskConfirmation("Are you sure?")
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	dnote, err := core.GetDnote(ctx)
	if err != nil {
		return err
	}

	for n, book := range dnote {
		if n == bookName {
			delete(dnote, n)

			err = core.LogActionRemoveBook(ctx, book.Name)
			if err != nil {
				return errors.Wrap(err, "Failed to log action")
			}
			err := core.WriteDnote(ctx, dnote)
			if err != nil {
				return err
			}

			fmt.Printf("[-] Deleted book : %s \n", bookName)
			return nil
		}
	}

	fmt.Println("Error : The book with that name is not found.")
	return nil
}
