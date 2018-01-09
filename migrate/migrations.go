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
	"gopkg.in/yaml.v2"
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
		// Find the minimum added_on timestamp from the notes that belong to the book
		// to give timstamp to the add_book action.
		// Logically add_book must have happened no later than the first add_note
		// to the book in order for sync to work.
		minTs := time.Now().Unix()
		for _, note := range book.Notes {
			if note.AddedOn < minTs {
				minTs = note.AddedOn
			}
		}

		action := migrateToV3Action{
			Type: migrateToV3ActionAddBook,
			Data: map[string]interface{}{
				"book_name": bookName,
			},
			Timestamp: minTs,
		}
		actions = append(actions, action)

		for _, note := range book.Notes {
			action := migrateToV3Action{
				Type: migrateToV3ActionAddNote,
				Data: map[string]interface{}{
					"note_uuid": note.UUID,
					"book_name": book.Name,
					"content":   note.Content,
				},
				Timestamp: note.AddedOn,
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

func getEditorCommand() string {
	editor := os.Getenv("EDITOR")

	if editor == "atom" {
		return "atom -w"
	} else if editor == "subl" {
		return "subl -n -w"
	} else if editor == "mate" {
		return "mate -w"
	}

	return "vim"
}

func migrateToV4(ctx infra.DnoteCtx) error {
	configPath := fmt.Sprintf("%s/dnoterc", ctx.DnoteDir)

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.Wrap(err, "Failed to read the config file")
	}

	var preConfig migrateToV4PreConfig
	err = yaml.Unmarshal(b, &preConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal existing config into JSON")
	}

	postConfig := migrateToV4PostConfig{
		APIKey: preConfig.APIKey,
		Editor: getEditorCommand(),
	}

	data, err := yaml.Marshal(postConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal config into JSON")
	}

	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to write the config into a file")
	}

	return nil
}
