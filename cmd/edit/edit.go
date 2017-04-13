package edit

import (
	"fmt"

	"encoding/json"
	"../../utils"
)

func Edit(notename string, newcontent string) error {
	// Get the current book.
	book, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	json_data, err := utils.GetDnote()
	if err != nil {
		return err
	}

	for _, note := range json_data[book] {

		if note.Name == notename {
			note.Content = newcontent
			out, err := json.Marshal(line)
			if err != nil {
				return err
			}
		
			note_data := "[" + string(out) + "]"
			json_data[book] = note_data // Assigns a new data to json map, does nto work (error: cannot use note_data (type string) as type utils.Book in assignment)
			fmt.Println(json_data[book])
		}
	}

	return nil 
}