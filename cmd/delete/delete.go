package delete

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/utils"
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
		Use:     "delete",
		Short:   "Delete a note or a book",
		Aliases: []string{"d"},
		Example: example,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&targetBookName, "book", "b", "", "The book name to delete")

	return cmd
}

func newRun(ctx infra.DnoteCtx) infra.RunEFunc {
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

			note(ctx, noteIndex, targetBook)
		}

		return nil
	}
}

// note deletes the note in a certain index.
func note(ctx infra.DnoteCtx, index int, book string) error {
	dnote, err := infra.GetDnote(ctx)
	if err != nil {
		return err
	}
	notes := dnote[book].Notes

	if len(notes)-1 < index {
		fmt.Println("Error : The note with that index is not found.")
		return nil
	}

	content := notes[index].Content
	fmt.Printf("Deleting note: %s\n", content)

	ok, err := utils.AskConfirmation("Are you sure?")
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	dnote[book] = infra.GetUpdatedBook(dnote[book], append(notes[:index], notes[index+1:]...))
	err = infra.WriteDnote(ctx, dnote)
	if err != nil {
		return err
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

	dnote, err := infra.GetDnote(ctx)
	if err != nil {
		return err
	}

	books, err := infra.GetBooks(ctx)
	if err != nil {
		return err
	}

	for _, book := range books {
		if book == bookName {
			delete(dnote, bookName)
			err := infra.WriteDnote(ctx, dnote)
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
