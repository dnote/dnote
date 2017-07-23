package delete

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/dnote-io/cli/utils"
	"github.com/spf13/cobra"
)

var targetBookName string

var example = `
  * Delete a note by its index from a book
  dnote delete js 2

  * Delete a book
  dnote delete -b js`

var cmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a note or a book",
	Aliases: []string{"d"},
	Example: example,
	RunE:    run,
}

func init() {
	root.Register(cmd)

	f := cmd.Flags()
	f.StringVarP(&targetBookName, "book", "b", "", "The book name to delete")
}

func run(cmd *cobra.Command, args []string) error {
	if targetBookName != "" {
		book(targetBookName)
	} else {
		if len(args) < 2 {
			return errors.New("Missing argument")
		}

		targetBook := args[0]
		noteIndex, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		note(noteIndex, targetBook)
	}

	return nil
}

// note deletes the note in a certain index.
func note(index int, book string) error {
	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	if len(dnote[book])-1 < index {
		fmt.Println("Error : The note with that index is not found.")
		return nil
	}
	
	content := dnote[book][index].Content
	fmt.Printf("Deleting note: %s\n", content)
	
	ok, err := utils.AskConfirmation("Are you sure?")
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	dnote[book] = append(dnote[book][:index], dnote[book][index+1:]...)
	err = utils.WriteDnote(dnote)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted!\n")
	return nil
}

// book deletes a book with the given name
func book(bookName string) error {
	ok, err := utils.AskConfirmation("Are you sure?")
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	books, err := utils.GetBooks()
	if err != nil {
		return err
	}

	for _, book := range books {
		if book == bookName {
			delete(dnote, bookName)
			err := utils.WriteDnote(dnote)
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
