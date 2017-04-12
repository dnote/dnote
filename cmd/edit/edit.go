package edit

import (
	"fmt"
	"github.com/dnote-io/cli/utils"
)

func Edit(notename string, newcontent string) error{
	book, book_err := utils.GetCurrentBook()
	json_data, json_err := utils.GetDnote()

	if book_err != nil {return book_err}
	if json_err != nil {return json_err}

	fmt.Println(book)
	fmt.Println(json_data)

	return nil 
}