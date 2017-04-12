package new

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dnote-io/cli/utils"
)

func Run(notename string, content string) error {
	currentBook, err := utils.GetCurrentBook()
	if err != nil {
		return err
	}

	note := makeNote(notename, content)
	err = writeNote(note)
	if err != nil {
		return err
	}

	fmt.Printf("[+] Added to %s\n", currentBook)
	return nil
}

func makeNote(notename string, content string) utils.Note {
	return utils.Note {
		UID: utils.GenerateNoteID(),
		Name: notename,
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

	d, err := json.MarshalIndent(dnote, "", "  ")
	if err != nil {
		return err
	}

	notePath, err := utils.GetDnotePath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(notePath, d, 0644)
	if err != nil {
		return err
	}

	return nil
}
