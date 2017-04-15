package delete

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"../../utils"
)

func DeleteNote(note_nu string) error {
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
            json_data[book][i] = json_data[book][len(json_data[book]) - 1]
            json_data[book] = json_data[book][:len(json_data)] // Currently this function deletes all the notes in a book.
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
	fmt.Printf("[+] Deleted note.")

	return nil
}