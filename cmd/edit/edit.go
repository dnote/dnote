package edit

import (
	"fmt"
	"io/ioutil"
	"encoding/json"

	"../../utils"
)

func Edit(note_nu string, newcontent string) error {
	book, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	json_data, err := utils.GetDnote()
	if err != nil {
		return err
	}

	for i, note := range json_data[book] {
		if note.Name == note_nu || note.UID == note_nu {
			note.Content = newcontent
			json_data[book][i] = note
		}
	}

	dnote_path, err := utils.GetDnotePath()
	if err != nil {
		return err
	}

	new_data, err := json.MarshalIndent(json_data, "", "	")
	if err != nil {
		return err
	}

	ioutil.WriteFile(dnote_path, new_data, 0644)
	fmt.Printf("[+] Edited %s", book)

	return nil 
}
