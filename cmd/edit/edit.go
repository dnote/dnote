package edit

import (
	"fmt"

	"io/ioutil"
    "strings"

	"../../utils"
)

func replaceLine(replace_data string, target_line string) error {
	path, err := utils.GetDnotePath()

	if err != nil {
		return err
	}

	input, err := ioutil.ReadFile(path)

	if err != nil {
		return nil
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, target_line) {
			lines[i] = replace_data
		}
	}

	return nil
}

func Edit(notename string, newcontent string) error {
	book, err := utils.GetCurrentBook()

	if err != nil {
		return err
	}

	json_data, err := utils.GetDnote()

	if err != nil {
		return err
	}

	for _, note := range json_data[book] {
		if note.Name == notename {
			note.Content = newcontent
			replaceLine(note.Content, "Content")
		}
	}

	return nil 
}