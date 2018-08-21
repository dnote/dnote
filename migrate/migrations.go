package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/utils"
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
		var notes = make([]migrateToV2PostNote, 0, len(book))
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

	var actions []migrateToV3Action

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

	switch editor {
	case "atom":
		return "atom -w"
	case "subl":
		return "subl -n -w"
	case "mate":
		return "mate -w"
	case "vim":
		return "vim"
	case "nano":
		return "nano"
	case "emacs":
		return "emacs"
	default:
		return "vi"
	}
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

// migrateToV5 migrates actions
func migrateToV5(ctx infra.DnoteCtx) error {
	actionsPath := fmt.Sprintf("%s/actions", ctx.DnoteDir)

	b, err := ioutil.ReadFile(actionsPath)
	if err != nil {
		return errors.Wrap(err, "reading the actions file")
	}

	var actions []migrateToV5PreAction
	err = json.Unmarshal(b, &actions)
	if err != nil {
		return errors.Wrap(err, "unmarshalling actions to JSON")
	}

	result := []migrateToV5PostAction{}

	for _, action := range actions {
		var data json.RawMessage

		switch action.Type {
		case migrateToV5ActionEditNote:
			var oldData migrateToV5PreEditNoteData
			if err = json.Unmarshal(action.Data, &oldData); err != nil {
				return errors.Wrapf(err, "unmarshalling old data of an edit note action %s", action.ID)
			}

			migratedData := migrateToV5PostEditNoteData{
				NoteUUID: oldData.NoteUUID,
				FromBook: oldData.BookName,
				Content:  oldData.Content,
			}
			b, err = json.Marshal(migratedData)
			if err != nil {
				return errors.Wrap(err, "marshalling data")
			}

			data = b
		default:
			data = action.Data
		}

		migrated := migrateToV5PostAction{
			UUID:      uuid.NewV4().String(),
			Schema:    1,
			Type:      action.Type,
			Data:      data,
			Timestamp: action.Timestamp,
		}

		result = append(result, migrated)
	}

	a, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "marshalling result into JSON")
	}
	err = ioutil.WriteFile(actionsPath, a, 0644)
	if err != nil {
		return errors.Wrap(err, "writing the result into a file")
	}

	return nil
}

// migrateToV6 adds a 'public' field to notes
func migrateToV6(ctx infra.DnoteCtx) error {
	notePath := fmt.Sprintf("%s/dnote", ctx.DnoteDir)

	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return errors.Wrap(err, "Failed to read the note file")
	}

	var preDnote migrateToV6PreDnote
	postDnote := migrateToV6PostDnote{}

	err = json.Unmarshal(b, &preDnote)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal existing dnote into JSON")
	}

	for bookName, book := range preDnote {
		var notes = make([]migrateToV6PostNote, 0, len(book.Notes))
		public := false
		for _, note := range book.Notes {
			newNote := migrateToV6PostNote{
				UUID:     note.UUID,
				Content:  note.Content,
				AddedOn:  note.AddedOn,
				EditedOn: note.EditedOn,
				Public:   &public,
			}

			notes = append(notes, newNote)
		}

		b := migrateToV6PostBook{
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
