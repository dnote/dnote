package delete

import (
	"fmt"
	"strconv"
	"os"

	"github.com/dnote-io/cli/utils"
)

// Bind the rest to one function for easier maintainance.
func Delete() error {
	if os.Args[2] == "-n" && len(os.Args) == 4{
		note_index, err := strconv.Atoi(os.Args[3])
		if err != nil {
			return err
		}

		target_book, err := utils.GetCurrentBook()
		if err != nil {
			return err
		}

		note(note_index, target_book)
	
	}else if os.Args[2] == "-n" && len(os.Args) == 5{
		note_index, err:= strconv.Atoi(os.Args[4])
		if err != nil {
			return err
		}

		target_book := os.Args[3]
		note(note_index, target_book)
	}else if os.Args[2] == "-b" {
		book(os.Args[3])
	}else{
		fmt.Println("Error : Invalid argument passed to delete.")
	}

	return nil
}

// Note deletes the note in a certain index.
func note(index int, book string) error {
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

			fmt.Printf("[-] Deleted : %d | Content : %s\n", index, dnote[book][index].Content)
			return nil
		}
	}

	fmt.Println("Error : The note with that index is not found.")
	return nil
}

// Book deletes a book with the given name
func book(bookName string) error {
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
