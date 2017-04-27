package delete

import (
	"fmt"
	"strconv"
	"os"

	"github.com/dnote-io/cli/utils"
)

func Delete() error {
	if len(os.Args) == 3 {
		current_book, err := utils.GetCurrentBook()
		if err != nil {
			return err
		}

		note_index, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return err
		}

		deleteNote(note_index, current_book)
	} else if len(os.Args) == 5 {
		if os.Args[2] == "-b" {
			note_index, err := strconv.Atoi(os.Args[4])
			if err != nil {
				return err
			}

			deleteNote(note_index, os.Args[3])
		}
	} else if len(os.Args) == 4 {		
		if os.Args[2] == "--book" {
			deleteBook(os.Args[3])
		}
	} else {
		fmt.Println("Invalid arguments passed to Delete.")
	}

	return nil
}

func deleteNote(index int, book string) error {
	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	for i, _ := range dnote[book] {
		if i == index{
			dnote[book] = append(dnote[book][:i], dnote[book][i+1:]...)
			err = utils.WriteDnote(dnote)
			if err != nil {
				return err
			}

			fmt.Printf("[-] Deleted : %d", index)
			return nil
		}
	}

	fmt.Println("Error : The note with that index is not found.")
	return nil
}

func deleteBook(bookName string) error {
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
