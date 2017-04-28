package delete

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dnote-io/cli/utils"
)


// Delete is a facade for deleting either note or book
func Delete() error {
	if os.Args[2] == "-b" {
		book(os.Args[3])
	} else if len(os.Args) == 4 {
		targetBook := os.Args[2]
		noteIndex, err := strconv.Atoi(os.Args[3])
		if err != nil {
			return err
		}

		note(noteIndex, targetBook)
	} else {
		fmt.Println("Error : Invalid argument passed to delete.")
	}

	return nil
}

// note deletes the note in a certain index.
func note(index int, book string) error {
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

	if len(dnote[book])-1 < index {
		fmt.Println("Error : The note with that index is not found.")
		return nil
	}

	content := dnote[book][index].Content
	dnote[book] = append(dnote[book][:index], dnote[book][index+1:]...)
	err = utils.WriteDnote(dnote)
	if err != nil {
		return err
	}

	fmt.Printf("[-] Deleted : %d | Content : %s\n", index, content)
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