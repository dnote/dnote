package new

import (
	"fmt"
	"time"

	"github.com/dnote-io/cli/utils"
)

func Run(content string) error {
	currentBook, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	note := makeNote(content)
	err = writeNote(note)
	if err != nil {
		return err
	}

	fmt.Printf("[+] Added to %s\n", currentBook)
	return nil
}

func makeNote(content string) utils.Note {
	return utils.Note{
		UID:     utils.GenerateNoteID(),
		Content: content,
		AddedOn: time.Now().Unix(),
	}
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

	err := utils.WriteDnote(dnote)
	if err != nil {
		return err
	}

	return nil
}
