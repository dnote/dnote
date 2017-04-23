package upgrade

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/dnote-io/cli/utils"
)

func Migrate() error {
	err := migrateYAMLToJSON()
	if err != nil {
		return err
	}

	err = migrateNoNameToName()
	if err != nil {
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
	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	books, err := utils.GetBooks()
	for _, book := range books {
		for i, note := range dnote[book] {
			if note.Name == "" {
				noteName, err := utils.GenerateNoteName()
				if err != nil {
					return err
				}

				note.Name = noteName
				dnote[book][i] = note
				err := utils.WriteDnote(dnote)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
