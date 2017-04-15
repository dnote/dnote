package delete

import (
	"fmt"
	"io/ioutil"
	"encoding/json"

	// For testing purposes.
	//"../../utils"

	// For GitHub.
	"github.com/dnote-io/cli/utils"
	
)

func DeleteNote(note_name_uid string) error {
	confirmation_message := "Delete note " + note_name_uid + " ?"
	isConfirmed, err := utils.AskConfirmation(confirmation_message)
	if err != nil {
		return err
	}
	var isDeleted bool

	if isConfirmed {
		book, err := utils.GetCurrentBook()
		if err != nil {
			return err
		}

		json_data, err := utils.GetDnote()
		if err != nil {
			return err
		}

		for i, note := range json_data[book] {
			if note.Name == note_name_uid || note.UID == note_name_uid {
				if len(json_data[book]) == 1 {
					json_data[book] = json_data[book][:0]
				}else if len(json_data[book]) > 1{				
					json_data[book][i] = json_data[book][len(json_data[book]) - 1]
					json_data[book] = json_data[book][:len(json_data)]
					isDeleted = true
				}
			}else if note.Name != note_name_uid || note.UID != note_name_uid {
				isDeleted = false
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

		if isDeleted {
			fmt.Printf("[+] Deleted %s \n", note_name_uid)
		}else {
			fmt.Println("[+] The note with that name is not found.")
		}
	}else{
		fmt.Println("[+] Deletion cancelled.")
	}

	return nil
}

func DeleteBook(book_name string) error {
	confirmation_message := "Delete book " + book_name + " ?"
	isConfirmed, err := utils.AskConfirmation(confirmation_message)
	if err != nil {
		return err
	}
	var isDnoteDeleted bool

	if isConfirmed {
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
				isDnoteDeleted = true
			}else if book != book_name {
				isDnoteDeleted = false
			}
		}

		dnote_path, err := utils.GetDnotePath()
		if err !=  nil {
			return err
		}

		new_data, err := json.MarshalIndent(json_data, "", "	")
		if err != nil {
			return err
		}

		ioutil.WriteFile(dnote_path, new_data, 0644)

		if isDnoteDeleted {
			fmt.Printf("[+] Deleted %s \n", book_name)
		}else {
			fmt.Println("[+] The book with that name is not found.")
		}
	}else{
		fmt.Println("[+] Deletion cancelled.")
	}
	
	return nil
}