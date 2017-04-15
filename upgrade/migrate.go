package upgrade

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	// For testing purposes.
	//"../utils"

	// For GitHub.
	
	"github.com/dnote-io/cli/utils"
	
)

func Migrate() error {
	err := migrateYAMLToJSON()
	if err != nil {
		return err
	}

	err = migrateNoNameToName()
	if err != nil{
		return err
	}

	return nil
}

// v0.0.4
// file format of .dnote changed to JSON from YAML
// Remove YAML support in v0.1.1, and archive this method
func migrateYAMLToJSON() error {
	usingYAML, err := isDnoteUsingYAML()
	if err != nil {
		return err
	}
	if !usingYAML {
		return nil
	}

	dnote, err := utils.GetNote()
	if err != nil {
		return err
	}

	jsonDnote := utils.Dnote{}

	for bookName, notes := range dnote {
		book := []utils.Note{}

		for _, note := range notes {
			note := utils.Note{
				UID:     utils.GenerateNoteID(),
				Content: note,
				AddedOn: time.Now().Unix(),
			}

			book = append(book, note)
		}

		jsonDnote[bookName] = book
	}

	migratedContent, err := json.MarshalIndent(jsonDnote, "", "  ")
	if err != nil {
		return err
	}

	dnotePath, err := utils.GetDnotePath()
	if err != nil {
		return err
	}

	archivePath, err := utils.GetYAMLDnoteArchivePath()
	if err != nil {
		return err
	}

	if err := os.Rename(dnotePath, archivePath); err != nil {
		return err
	}

	err = ioutil.WriteFile(dnotePath, migratedContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func migrateNoNameToName() error {
	json_data, err := utils.GetDnote()
	if err != nil {
		return err
	}

	new_note := utils.Dnote{}
	books, err := utils.GetBooks()

	for _, book := range books {
		new_note = json_data
		
		for i, note := range new_note[book] {
			if note.Name == "" {
				note_name, err := utils.GenerateNoteName()
				if err != nil {
					return err
				}

				note.Name = note_name
				json_data[book][i] = note
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

	return nil 
}