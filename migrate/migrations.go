package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/dnote-io/cli/utils"
	"github.com/pkg/errors"
)

func deleteDnoteYAMLArchive() error {
	yamlPath, err := getYAMLDnoteArchivePath()
	if err != nil {
		return errors.Wrap(err, "Failed to get YAML path")
	}

	if !utils.FileExists(yamlPath) {
		return nil
	}

	if err := os.Remove(yamlPath); err != nil {
		return errors.Wrap(err, "Failed to delete .dnote archive")
	}

	return nil
}

func generateBookMetadata() error {
	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "Failed to get current os user")
	}

	notePath := fmt.Sprintf("%s/.dnote/dnote", usr.HomeDir)
	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return errors.Wrap(err, "Failed to read the note file")
	}

	type Note struct {
		UID     string
		Content string
		AddedOn int64
	}
	type PreBook []Note
	type PostBook struct {
		UID   string
		Notes []Note
	}

	type PreDnote map[string]PreBook
	type PostDnote map[string]PostBook

	var preDnote PreDnote
	postDnote := PostDnote{}

	err = json.Unmarshal(b, &preDnote)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal existing dnote into JSON")
	}

	for bookName, book := range preDnote {
		b := PostBook{
			UID:   utils.GenerateUID(),
			Notes: book,
		}

		postDnote[bookName] = b
	}

	d, err := json.MarshalIndent(postDnote, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Failed to marshal new dnote into JSON")
	}

	err = ioutil.WriteFile(notePath, d, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to write the new dnote into the file")
	}

	return nil
}
