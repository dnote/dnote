package notes

import (
	"fmt"
	"os"

	"github.com/dnote-io/cli/utils"
)

func Run() error {
	defaultBookName, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	var bookName string

	if len(os.Args) == 2 {
		bookName = defaultBookName
	} else if len(os.Args) == 4 && os.Args[2] == "-b" {
		bookName = os.Args[3]
	} else {
		fmt.Println("Invalid argument passed to notes")
		os.Exit(1)
	}

	fmt.Printf("On note %s\n", bookName)

	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	for k, v := range dnote {
		if k == bookName {
			for _, note := range v {
				fmt.Printf("* %s\n", note.Content)
			}
		}
	}

	//sort.Strings(notes)

	return nil
}
