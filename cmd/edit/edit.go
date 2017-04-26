package edit

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dnote-io/cli/utils"
)

func Edit() error {
	current_book, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	var target_book string
	var index int
	var content string

	if len(os.Args) == 4 {
		target_book = current_book
		index, err = strconv.Atoi(os.Args[2])
		if err != nil {
			return nil
		}
		content = os.Args[3]
	}else if len(os.Args) == 5 {
		target_book = os.Args[2]
		index, err = strconv.Atoi(os.Args[3])
		if err != nil {
			return err
		}
		content = os.Args[4]
	}

	for i, note := range dnote[target_book] {
		if i == index {
			note.Content = content
			dnote[target_book][i] = note

			err := utils.WriteDnote(dnote)
			fmt.Printf("[+] Edited Note : %d \n", index)
			return err
		}
	}

	// If loop finishes without returning, note did not exist
	fmt.Println("Error : The note with that index is not found.")
	return nil
}
