package books

import (
	"fmt"

	"github.com/dnote-io/cli/utils"
)

func Run() error {
	currentBook, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	books, err := utils.GetBooks()
	if err != nil {
		return err
	}

	for _, book := range books {
		if book == currentBook {
			fmt.Printf("* %v\n", book)
		} else {
			fmt.Printf("  %v\n", book)
		}
	}

	return nil
}
