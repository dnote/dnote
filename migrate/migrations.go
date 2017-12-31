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

// migrateToV1 deletes YAML archive if exists
func migrateToV1(ctx infra.DnoteCtx) error {
	yamlPath := fmt.Sprintf("%s/%s", ctx.HomeDir, ".dnote-yaml-archived")
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
				UUID:     uuid.NewV4().String(),
				Content:  note.Content,
				AddedOn:  note.AddedOn,
				EditedOn: 0,
			}

			notes = append(notes, newNote)
		}

		b := migrateToV2PostBook{
			UUID:  uuid.NewV4().String(),
			Name:  bookName,
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

// migrateToV3 generates actions for existing dnote
func migrateToV3(ctx infra.DnoteCtx) error {
	notePath := fmt.Sprintf("%s/dnote", ctx.DnoteDir)
	actionsPath := fmt.Sprintf("%s/actions", ctx.DnoteDir)

	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return errors.Wrap(err, "Failed to read the note file")
	}

	var dnote migrateToV3Dnote

	err = json.Unmarshal(b, &dnote)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal existing dnote into JSON")
	}

	actions := []migrateToV3Action{}

	for bookName, book := range dnote {
		action := migrateToV3Action{
			Type: migrateToV3ActionAddBook,
			Data: map[string]interface{}{
				"uuid": book.UUID,
				"name": bookName,
			},
			Timestamp: time.Now().Unix(),
		}
		actions = append(actions, action)

		for _, note := range book.Notes {
			action := migrateToV3Action{
				Type: migrateToV3ActionAddNote,
				Data: map[string]interface{}{
					"note_uuid": note.UUID,
					"book_uuid": book.UUID,
					"content":   note.Content,
				},
				Timestamp: time.Now().Unix(),
			}
			actions = append(actions, action)

		}
	}

	a, err := json.Marshal(actions)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal actions into JSON")
	}

	err = ioutil.WriteFile(actionsPath, a, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to write the actions into a file")
	}

	return nil
}
