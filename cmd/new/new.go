package new

import (
	"fmt"
	"time"

	"github.com/dnote-io/cli/utils"
)

// Run makes a new note
func Run(noteName string, content string) error {
	currentBook, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	note, err := makeNote(noteName, content)
	if err != nil {
		return err
	}

	err = writeNote(note)
	if err != nil {
		return err
	}

	fmt.Printf("[+] Added to %s\n", currentBook)
	return nil
}

func makeNote(noteName string, content string) (utils.Note, error) {
	var note utils.Note
	if noteName == "" {
		newName, err := utils.GenerateNoteName()
		if err != nil {
			return note, err
		}

		noteName = newName
	}

	note = utils.Note{
		UID:     utils.GenerateNoteID(),
		Name:    noteName,
		Content: content,
		AddedOn: time.Now().Unix(),
	}

	return note, nil
}

func writeNote(note utils.Note) error {
	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	book, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	if _, ok := dnote[book]; ok {
		dnote[book] = append(dnote[book], note)
	} else {
		dnote[book] = []utils.Note{note}
	}

	err = utils.WriteDnote(dnote)
	if err != nil {
		return err
	}

	return nil
}
