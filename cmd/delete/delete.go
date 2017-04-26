package delete

import (
	"fmt"
	"strconv"
	"os"

	"github.com/dnote-io/cli/utils"
)

func Delete() error {
	if os.Args[2] == "-b" {
		book(os.Args[3])
	}else if os.Args[2] == "-n" {
		note()
	}else{
		fmt.Println("Invalid command.")
	}

	return nil
}

// Note deletes the book
func note() error {
	book, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	var target_book string
	var index int

	if len(os.Args) == 4 {
		target_book = book
		index, err = strconv.Atoi(os.Args[3])
		if err != nil {
			return nil
		}
	}else if len(os.Args) == 5 {
		target_book = os.Args[3]
		index, err = strconv.Atoi(os.Args[4])
		if err != nil {
			return err
		}
	}

	for i, _ := range dnote[target_book] {
		if i == index{
			dnote[target_book] = append(dnote[target_book][:i], dnote[target_book][i+1:]...)
			err = utils.WriteDnote(dnote)
			if err != nil {
				return err
			}

			fmt.Printf("[-] Deleted %d \n", index)
			return nil
		}
	}

	fmt.Println("The note with that index is not found.")
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

			fmt.Printf("[-] Deleted the book %s", bookName)
			return nil
		}
	}

	fmt.Println("[+] The book with that name is not found.")
	return nil
}
