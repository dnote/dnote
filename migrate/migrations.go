package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/utils"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
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

func migrateToV2(ctx infra.DnoteCtx) error {
	notePath := fmt.Sprintf("%s/dnote", ctx.DnoteDir)

	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return errors.Wrap(err, "Failed to read the note file")
	}

	var preDnote migrateToV2PreDnote
	postDnote := migrateToV2PostDnote{}

	err = json.Unmarshal(b, &preDnote)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal existing dnote into JSON")
	}

	for bookName, book := range preDnote {
		notes := []migrateToV2PostNote{}
		for _, note := range book {
			newNote := migrateToV2PostNote{
				UUID:    uuid.NewV4().String(),
				Content: note.Content,
				AddedOn: note.AddedOn,
			}

			notes = append(notes, newNote)
		}

		b := migrateToV2PostBook{
			UUID:  uuid.NewV4().String(),
			Notes: notes,
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

func migrateToV3(ctx infra.DnoteCtx) error {
	notePath := fmt.Sprintf("%s/dnote", ctx.DnoteDir)
	actionsPath := fmt.Sprintf("%s/actions", ctx.DnoteDir)

	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return errors.Wrap(err, "Failed to read the note file")
	}

	var preDnote migrateToV3PreDnote
	postDnote := migrateToV3PostDnote{}

	err = json.Unmarshal(b, &preDnote)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal existing dnote into JSON")
	}

	actions := []migrateToV3Action{}

	for bookName, book := range preDnote {
		action := migrateToV3Action{
			Type: migrateToV3ActionAddBook,
			Data: map[string]interface{}{
				"UUID": book.UUID,
				"Name": bookName,
			},
			Timestamp: time.Now().Unix(),
		}
		actions = append(actions, action)

		notes := []migrateToV3PostNote{}
		for _, note := range book.Notes {
			newNote := migrateToV3PostNote{
				UUID:    note.UUID,
				Content: note.Content,
			}
			action := migrateToV3Action{
				Type: migrateToV3ActionAddNote,
				Data: map[string]interface{}{
					"UUID":    note.UUID,
					"Content": note.Content,
				},
				Timestamp: time.Now().Unix(),
			}
			actions = append(actions, action)

			notes = append(notes, newNote)
		}

		b := migrateToV3PostBook{
			UUID:  book.UUID,
			Notes: notes,
		}

		postDnote[bookName] = b
	}

	d, err := json.MarshalIndent(postDnote, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Failed to marshal new dnote into JSON")
	}

	a, err := json.Marshal(actions)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal actions into JSON")
	}

	err = ioutil.WriteFile(notePath, d, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to write the new dnote into the file")
	}
	err = ioutil.WriteFile(actionsPath, a, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to write the actions into a file")
	}

	return nil
}
