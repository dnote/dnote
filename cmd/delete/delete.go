package delete

import (
	"fmt"

	"github.com/dnote-io/cli/utils"
)

// Note deletes the book
func Note(nameOrUID string) error {
	confirmMsg := "Delete note " + nameOrUID + " ?"
	isConfirmed, err := utils.AskConfirmation(confirmMsg)
	if err != nil {
		return err
	}
	if !isConfirmed {
		fmt.Println("Deletion cancelled.")
		return nil
	}

	book, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	currentBook := dnote[book]
	for i, note := range currentBook {
		if note.Name == nameOrUID || note.UID == nameOrUID {
			currentBook = append(currentBook[:i], currentBook[i+1:]...)
			err = utils.WriteDnote(dnote)
			if err != nil {
				return err
			}

			fmt.Printf("[-] Deleted %s \n", nameOrUID)
			return nil
		}
	}

	fmt.Println("The note with that name is not found.")
	return nil
}

// Book deletes a book with the given name
func Book(bookName string) error {
	confirmMsg := "Delete book " + bookName + " ?"
	isConfirmed, err := utils.AskConfirmation(confirmMsg)
	if err != nil {
		return err
	}
	if !isConfirmed {
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

			fmt.Printf("[-] Deleted the book %s", bookName)
			return nil
		}
	}

	fmt.Println("[+] The book with that name is not found.")
	return nil
}
