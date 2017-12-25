package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/utils"
	"github.com/pkg/errors"
)

func deleteDnoteYAMLArchive(ctx infra.DnoteCtx) error {
	yamlPath, err := getYAMLDnoteArchivePath(ctx)
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

type generateBookMetadataNote struct {
	UID     string
	Content string
	AddedOn int64
}
type generateBookMetadataPreBook []generateBookMetadataNote
type generateBookMetadataPostBook struct {
	UID   string
	Notes []generateBookMetadataNote
}
type generateBookMetadataPreDnote map[string]generateBookMetadataPreBook
type generateBookMetadataPostDnote map[string]generateBookMetadataPostBook

func generateBookMetadata(ctx infra.DnoteCtx) error {
	notePath := fmt.Sprintf("%s/dnote", ctx.DnoteDir)
	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return errors.Wrap(err, "Failed to read the note file")
	}

	var preDnote generateBookMetadataPreDnote
	postDnote := generateBookMetadataPostDnote{}

	err = json.Unmarshal(b, &preDnote)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal existing dnote into JSON")
	}

	for bookName, book := range preDnote {
		b := generateBookMetadataPostBook{
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
