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
			if len(json_data[book]) == 1 {
				json_data[book] = json_data[book][:0]
			}else if len(json_data[book]) > 1{
				json_data[book][i] = json_data[book][len(json_data[book]) - 1]
				json_data[book] = json_data[book][:len(json_data)]
			}
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
	fmt.Println("[+] Deleted note.")

	return nil
}

func DeleteBook(book_name string) error {	// Delete any entry of the book, .dnote config and .dnote file.
	json_data, err := utils.GetDnote()
	if err != nil {
		return err
	}

	book_data, err := utils.GetBooks()
	if err != nil {
		return err
	}

	for _, book := range book_data {
		if book == book_name {
			delete(json_data, book_name)
		}
	}

	fmt.Println("[+] Deleted book.")
	return nil
}