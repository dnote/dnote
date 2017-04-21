package edit

import (
	"fmt"
	"io/ioutil"
	"encoding/json"

	// For testing purposes.
	//"../../utils"

	// For GitHub.
	"github.com/dnote-io/cli/utils"
	
)

func Edit(note_name_uid string, newcontent string) error {
	book, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	json_data, err := utils.GetDnote()
	if err != nil {
		return err
	}

	var noteFound bool
	for i, note := range json_data[book] {
		if note.Name == note_name_uid || note.UID == note_name_uid {
			note.Content = newcontent
			json_data[book][i] = note
			noteFound = true
			break
		}else{
			noteFound = false
		}
	}

	if noteFound != true{
		fmt.Println("[+] The note with that name / UID is not found.")
		return nil
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
