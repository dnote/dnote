package new

import (
	"fmt"
	"time"

	"github.com/dnote-io/cli/utils"
)

func Run(bookName string, content string) error {
	note := makeNote(content)
	err := writeNote(bookName, note)
	if err != nil {
		return err
	}

	fmt.Printf("[+] Added to %s\n", bookName)
	return nil
}

func makeNote(content string) utils.Note {
	return utils.Note{
		UID:     utils.GenerateNoteID(),
		Content: content,
		AddedOn: time.Now().Unix(),
	}
}

func writeNote(bookName string, note utils.Note) error {
	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	if _, ok := dnote[bookName]; ok {
		dnote[bookName] = append(dnote[bookName], note)
	} else {
		dnote[bookName] = []utils.Note{note}
	}

	err = utils.WriteDnote(dnote)
	if err != nil {
		return err
	}

	return nil
}
