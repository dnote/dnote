package edit

import (
	"fmt"

	"github.com/dnote-io/cli/utils"
)

// Edit edits dnote with the given name or uid
func Edit(nameOrUID string, content string) error {
	book, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	for i, note := range dnote[book] {
		if note.Name == nameOrUID || note.UID == nameOrUID {
			note.Content = content
			dnote[book][i] = note

			err := utils.WriteDnote(dnote)
			return err
		}
	}

	// If loop finishes without returning, note did not exist
	fmt.Println("[+] The note with that name / UID is not found.")
	return nil
}
