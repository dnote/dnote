package new

import (
	"fmt"
	"time"

	// For testing purposes.
	//"../../utils"

	// For GitHub.
	"github.com/dnote-io/cli/utils"
	
)

func Run(notename string, content string) error {
	currentBook, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	note, err:= makeNote(notename, content)
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

func makeNote(notename string, content string) (utils.Note, error) {
	var note utils.Note
    if notename == "" {
        auto_gen_name, err := utils.GenerateNoteName()
        if err != nil {
            return note, err
        }

        note = utils.Note {
            UID: utils.GenerateNoteID(),
            Name: auto_gen_name,
            Content: content,
            AddedOn: time.Now().Unix(),
        }
    } else {
        note = utils.Note {
            UID: utils.GenerateNoteID(),
            Name: notename,
            Content: content,
            AddedOn: time.Now().Unix(),
        }
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

	err := utils.WriteDnote(dnote)
	if err != nil {
		return err
	}

	return nil
}