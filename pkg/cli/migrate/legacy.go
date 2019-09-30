/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package migrate provides migration logic for both sqlite and
// legacy JSON-based notes used until v0.4.x releases
package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

var (
	schemaFilename = "schema"
	backupDirName  = ".dnote-bak"
)

// migration IDs
const (
	_ = iota
	legacyMigrationV1
	legacyMigrationV2
	legacyMigrationV3
	legacyMigrationV4
	legacyMigrationV5
	legacyMigrationV6
	legacyMigrationV7
	legacyMigrationV8
)

var migrationSequence = []int{
	legacyMigrationV1,
	legacyMigrationV2,
	legacyMigrationV3,
	legacyMigrationV4,
	legacyMigrationV5,
	legacyMigrationV6,
	legacyMigrationV7,
	legacyMigrationV8,
}

type schema struct {
	CurrentVersion int `yaml:"current_version"`
}

func makeSchema(complete bool) schema {
	s := schema{}

	var CurrentVersion int
	if complete {
		CurrentVersion = len(migrationSequence)
	}

	s.CurrentVersion = CurrentVersion

	return s
}

// Legacy performs migration on JSON-based dnote if necessary
func Legacy(ctx context.DnoteCtx) error {
	// If schema does not exist, no need run a legacy migration
	schemaPath := getSchemaPath(ctx)
	ok, err := utils.FileExists(schemaPath)
	if err != nil {
		return errors.Wrap(err, "checking if schema exists")
	}
	if !ok {
		return nil
	}

	unrunMigrations, err := getUnrunMigrations(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get unrun migrations")
	}

	for _, mig := range unrunMigrations {
		log.Debug("running legacy migration %d\n", mig)
		if err := performMigration(ctx, mig); err != nil {
			return errors.Wrapf(err, "running migration #%d", mig)
		}
	}

	return nil
}

// performMigration backs up current .dnote data, performs migration, and
// restores or cleans backups depending on if there is an error
func performMigration(ctx context.DnoteCtx, migrationID int) error {
	// legacyMigrationV8 is the final migration of the legacy JSON Dnote migration
	// migrate to sqlite and return
	if migrationID == legacyMigrationV8 {
		if err := migrateToV8(ctx); err != nil {
			return errors.Wrap(err, "migrating to sqlite")
		}

		return nil
	}

	if err := backupDnoteDir(ctx); err != nil {
		return errors.Wrap(err, "Failed to back up dnote directory")
	}

	var migrationError error

	switch migrationID {
	case legacyMigrationV1:
		migrationError = migrateToV1(ctx)
	case legacyMigrationV2:
		migrationError = migrateToV2(ctx)
	case legacyMigrationV3:
		migrationError = migrateToV3(ctx)
	case legacyMigrationV4:
		migrationError = migrateToV4(ctx)
	case legacyMigrationV5:
		migrationError = migrateToV5(ctx)
	case legacyMigrationV6:
		migrationError = migrateToV6(ctx)
	case legacyMigrationV7:
		migrationError = migrateToV7(ctx)
	default:
		return errors.Errorf("Unrecognized migration id %d", migrationID)
	}

	if migrationError != nil {
		if err := restoreBackup(ctx); err != nil {
			panic(errors.Wrap(err, "Failed to restore backup for a failed migration"))
		}

		return errors.Wrapf(migrationError, "Failed to perform migration #%d", migrationID)
	}

	if err := clearBackup(ctx); err != nil {
		return errors.Wrap(err, "Failed to clear backup")
	}

	if err := updateSchemaVersion(ctx, migrationID); err != nil {
		return errors.Wrap(err, "Failed to update schema version")
	}

	return nil
}

// backupDnoteDir backs up the dnote directory to a temporary backup directory
func backupDnoteDir(ctx context.DnoteCtx) error {
	srcPath := fmt.Sprintf("%s/.dnote", ctx.HomeDir)
	tmpPath := fmt.Sprintf("%s/%s", ctx.HomeDir, backupDirName)

	if err := utils.CopyDir(srcPath, tmpPath); err != nil {
		return errors.Wrap(err, "Failed to copy the .dnote directory")
	}

	return nil
}

func restoreBackup(ctx context.DnoteCtx) error {
	var err error

	defer func() {
		if err != nil {
			log.Printf(`Failed to restore backup for a failed migration.
	Don't worry. Your data is still intact in the backup directory.
	Get help on https://github.com/dnote/dnote/pkg/cli/issues`)
		}
	}()

	srcPath := fmt.Sprintf("%s/.dnote", ctx.HomeDir)
	backupPath := fmt.Sprintf("%s/%s", ctx.HomeDir, backupDirName)

	if err = os.RemoveAll(srcPath); err != nil {
		return errors.Wrapf(err, "Failed to clear current dnote data at %s", backupPath)
	}

	if err = os.Rename(backupPath, srcPath); err != nil {
		return errors.Wrap(err, `Failed to copy backup data to the original directory.`)
	}

	return nil
}

func clearBackup(ctx context.DnoteCtx) error {
	backupPath := fmt.Sprintf("%s/%s", ctx.HomeDir, backupDirName)

	if err := os.RemoveAll(backupPath); err != nil {
		return errors.Wrapf(err, "Failed to remove backup at %s", backupPath)
	}

	return nil
}

// getSchemaPath returns the path to the file containing schema info
func getSchemaPath(ctx context.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, schemaFilename)
}

func readSchema(ctx context.DnoteCtx) (schema, error) {
	var ret schema

	path := getSchemaPath(ctx)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to read schema file")
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to unmarshal the schema JSON")
	}

	return ret, nil
}

func writeSchema(ctx context.DnoteCtx, s schema) error {
	path := getSchemaPath(ctx)
	d, err := yaml.Marshal(&s)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal schema into yaml")
	}

	if err := ioutil.WriteFile(path, d, 0644); err != nil {
		return errors.Wrap(err, "Failed to write schema file")
	}

	return nil
}

func getUnrunMigrations(ctx context.DnoteCtx) ([]int, error) {
	var ret []int

	schema, err := readSchema(ctx)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to read schema")
	}

	log.Debug("current legacy schema: %d\n", schema.CurrentVersion)

	if schema.CurrentVersion == len(migrationSequence) {
		return ret, nil
	}

	nextVersion := schema.CurrentVersion
	ret = migrationSequence[nextVersion:]

	return ret, nil
}

func updateSchemaVersion(ctx context.DnoteCtx, mID int) error {
	s, err := readSchema(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to read schema")
	}

	s.CurrentVersion = mID

	err = writeSchema(ctx, s)
	if err != nil {
		return errors.Wrap(err, "Failed to write schema")
	}

	return nil
}

/***** snapshots **/

// v2
type migrateToV2PreNote struct {
	UID     string
	Content string
	AddedOn int64
}
type migrateToV2PostNote struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"editd_on"`
}
type migrateToV2PreBook []migrateToV2PreNote
type migrateToV2PostBook struct {
	Name  string                `json:"name"`
	Notes []migrateToV2PostNote `json:"notes"`
}
type migrateToV2PreDnote map[string]migrateToV2PreBook
type migrateToV2PostDnote map[string]migrateToV2PostBook

//v3
var (
	migrateToV3ActionAddNote = "add_note"
	migrateToV3ActionAddBook = "add_book"
)

type migrateToV3Note struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
}
type migrateToV3Book struct {
	UUID  string            `json:"uuid"`
	Name  string            `json:"name"`
	Notes []migrateToV3Note `json:"notes"`
}
type migrateToV3Dnote map[string]migrateToV3Book
type migrateToV3Action struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

// v4
type migrateToV4PreConfig struct {
	Book   string
	APIKey string
}
type migrateToV4PostConfig struct {
	Editor string
	APIKey string
}

// v5
type migrateToV5AddNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookName string `json:"book_name"`
	Content  string `json:"content"`
}
type migrateToV5RemoveNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookName string `json:"book_name"`
}
type migrateToV5AddBookData struct {
	BookName string `json:"book_name"`
}
type migrateToV5RemoveBookData struct {
	BookName string `json:"book_name"`
}
type migrateToV5PreEditNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookName string `json:"book_name"`
	Content  string `json:"content"`
}
type migrateToV5PostEditNoteData struct {
	NoteUUID string `json:"note_uuid"`
	FromBook string `json:"from_book"`
	ToBook   string `json:"to_book"`
	Content  string `json:"content"`
}
type migrateToV5PreAction struct {
	ID        int             `json:"id"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}
type migrateToV5PostAction struct {
	UUID      string          `json:"uuid"`
	Schema    int             `json:"schema"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

var (
	migrateToV5ActionAddNote    = "add_note"
	migrateToV5ActionRemoveNote = "remove_note"
	migrateToV5ActionEditNote   = "edit_note"
	migrateToV5ActionAddBook    = "add_book"
	migrateToV5ActionRemoveBook = "remove_book"
)

// v6
type migrateToV6PreNote struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
}
type migrateToV6PostNote struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
	// Make a pointer to test absent values
	Public *bool `json:"public"`
}
type migrateToV6PreBook struct {
	Name  string               `json:"name"`
	Notes []migrateToV6PreNote `json:"notes"`
}
type migrateToV6PostBook struct {
	Name  string                `json:"name"`
	Notes []migrateToV6PostNote `json:"notes"`
}
type migrateToV6PreDnote map[string]migrateToV6PreBook
type migrateToV6PostDnote map[string]migrateToV6PostBook

// v7
var migrateToV7ActionTypeEditNote = "edit_note"

type migrateToV7Action struct {
	UUID      string          `json:"uuid"`
	Schema    int             `json:"schema"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}
type migrateToV7EditNoteDataV1 struct {
	NoteUUID string `json:"note_uuid"`
	FromBook string `json:"from_book"`
	ToBook   string `json:"to_book"`
	Content  string `json:"content"`
}
type migrateToV7EditNoteDataV2 struct {
	NoteUUID string  `json:"note_uuid"`
	FromBook string  `json:"from_book"`
	ToBook   *string `json:"to_book"`
	Content  *string `json:"content"`
	Public   *bool   `json:"public"`
}

// v8
type migrateToV8Action struct {
	UUID      string          `json:"uuid"`
	Schema    int             `json:"schema"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}
type migrateToV8Note struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
	// Make a pointer to test absent values
	Public *bool `json:"public"`
}
type migrateToV8Book struct {
	Name  string            `json:"name"`
	Notes []migrateToV8Note `json:"notes"`
}
type migrateToV8Dnote map[string]migrateToV8Book
type migrateToV8Timestamp struct {
	LastUpgrade int64 `yaml:"last_upgrade"`
	Bookmark    int   `yaml:"bookmark"`
	LastAction  int64 `yaml:"last_action"`
}

var migrateToV8SystemKeyLastUpgrade = "last_upgrade"
var migrateToV8SystemKeyLastAction = "last_action"
var migrateToV8SystemKeyBookMark = "bookmark"

/***** migrations **/

// migrateToV1 deletes YAML archive if exists
func migrateToV1(ctx context.DnoteCtx) error {
	yamlPath := fmt.Sprintf("%s/%s", ctx.HomeDir, ".dnote-yaml-archived")
	ok, err := utils.FileExists(yamlPath)
	if err != nil {
		return errors.Wrap(err, "checking if yaml file exists")
	}
	if !ok {
		return nil
	}

	if err := os.Remove(yamlPath); err != nil {
		return errors.Wrap(err, "Failed to delete .dnote archive")
	}

	return nil
}

func migrateToV2(ctx context.DnoteCtx) error {
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
func migrateToV3(ctx context.DnoteCtx) error {
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

func migrateToV4(ctx context.DnoteCtx) error {
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
func migrateToV5(ctx context.DnoteCtx) error {
	actionsPath := fmt.Sprintf("%s/actions", ctx.DnoteDir)

	b, err := ioutil.ReadFile(actionsPath)
	if err != nil {
		return errors.Wrap(err, "reading the actions file")
	}

	var actions []migrateToV5PreAction
	err = json.Unmarshal(b, &actions)
	if err != nil {
		return errors.Wrap(err, "unmarshalling actions from JSON")
	}

	result := []migrateToV5PostAction{}

	for _, action := range actions {
		var data json.RawMessage

		switch action.Type {
		case migrateToV5ActionEditNote:
			var oldData migrateToV5PreEditNoteData
			if err = json.Unmarshal(action.Data, &oldData); err != nil {
				return errors.Wrapf(err, "unmarshalling old data of an edit note action %d", action.ID)
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
func migrateToV6(ctx context.DnoteCtx) error {
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

// migrateToV7 migrates data of edit_note action to the proper version which is
// EditNoteDataV2. Due to a bug, edit logged actions with schema version '2'
// but with a data of EditNoteDataV1. https://github.com/dnote/dnote/pkg/cli/issues/107
func migrateToV7(ctx context.DnoteCtx) error {
	actionPath := fmt.Sprintf("%s/actions", ctx.DnoteDir)

	b, err := ioutil.ReadFile(actionPath)
	if err != nil {
		return errors.Wrap(err, "reading actions file")
	}

	var preActions []migrateToV7Action
	postActions := []migrateToV7Action{}
	err = json.Unmarshal(b, &preActions)
	if err != nil {
		return errors.Wrap(err, "unmarshalling existing actions")
	}

	for _, action := range preActions {
		var newAction migrateToV7Action

		if action.Type == migrateToV7ActionTypeEditNote {
			var oldData migrateToV7EditNoteDataV1
			if e := json.Unmarshal(action.Data, &oldData); e != nil {
				return errors.Wrapf(e, "unmarshalling data of action with uuid %s", action.Data)
			}

			newData := migrateToV7EditNoteDataV2{
				NoteUUID: oldData.NoteUUID,
				FromBook: oldData.FromBook,
				ToBook:   nil,
				Content:  &oldData.Content,
				Public:   nil,
			}
			d, e := json.Marshal(newData)
			if e != nil {
				return errors.Wrapf(e, "marshalling new data of action with uuid %s", action.Data)
			}

			newAction = migrateToV7Action{
				UUID:      action.UUID,
				Schema:    action.Schema,
				Type:      action.Type,
				Timestamp: action.Timestamp,
				Data:      d,
			}
		} else {
			newAction = action
		}

		postActions = append(postActions, newAction)
	}

	d, err := json.Marshal(postActions)
	if err != nil {
		return errors.Wrap(err, "marshalling new actions")
	}

	err = ioutil.WriteFile(actionPath, d, 0644)
	if err != nil {
		return errors.Wrap(err, "writing new actions to a file")
	}

	return nil
}

// migrateToV8 migrates dnote data to sqlite database
func migrateToV8(ctx context.DnoteCtx) error {
	tx, err := ctx.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	// 1. Migrate the the dnote file
	dnoteFilePath := fmt.Sprintf("%s/dnote", ctx.DnoteDir)
	b, err := ioutil.ReadFile(dnoteFilePath)
	if err != nil {
		return errors.Wrap(err, "reading the notes")
	}

	var dnote migrateToV8Dnote
	err = json.Unmarshal(b, &dnote)
	if err != nil {
		return errors.Wrap(err, "unmarshalling notes to JSON")
	}

	for bookName, book := range dnote {
		bookUUID := uuid.NewV4().String()
		_, err = tx.Exec(`INSERT INTO books (uuid, label) VALUES (?, ?)`, bookUUID, bookName)
		if err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "inserting book %s", book.Name)
		}

		for _, note := range book.Notes {
			_, err = tx.Exec(`INSERT INTO notes
			(uuid, book_uuid, content, added_on, edited_on, public)
			VALUES (?, ?, ?, ?, ?, ?)
			`, note.UUID, bookUUID, note.Content, note.AddedOn, note.EditedOn, note.Public)

			if err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "inserting the note %s", note.UUID)
			}
		}
	}

	// 2. Migrate the actions file
	actionsPath := fmt.Sprintf("%s/actions", ctx.DnoteDir)
	b, err = ioutil.ReadFile(actionsPath)
	if err != nil {
		return errors.Wrap(err, "reading the actions")
	}

	var actions []migrateToV8Action
	err = json.Unmarshal(b, &actions)
	if err != nil {
		return errors.Wrap(err, "unmarshalling actions from JSON")
	}

	for _, action := range actions {
		_, err = tx.Exec(`INSERT INTO actions
			(uuid, schema, type, data, timestamp)
			VALUES (?, ?, ?, ?, ?)
			`, action.UUID, action.Schema, action.Type, action.Data, action.Timestamp)

		if err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "inserting the action %s", action.UUID)
		}
	}

	// 3. Migrate the timestamps file
	timestampsPath := fmt.Sprintf("%s/timestamps", ctx.DnoteDir)
	b, err = ioutil.ReadFile(timestampsPath)
	if err != nil {
		return errors.Wrap(err, "reading the timestamps")
	}

	var timestamp migrateToV8Timestamp
	err = yaml.Unmarshal(b, &timestamp)
	if err != nil {
		return errors.Wrap(err, "unmarshalling timestamps from YAML")
	}

	_, err = tx.Exec(`INSERT INTO system (key, value) VALUES (?, ?)`,
		migrateToV8SystemKeyLastUpgrade, timestamp.LastUpgrade)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "inserting the last_upgrade value")
	}
	_, err = tx.Exec(`INSERT INTO system (key, value) VALUES (?, ?)`,
		migrateToV8SystemKeyLastAction, timestamp.LastAction)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "inserting the last_action value")
	}
	_, err = tx.Exec(`INSERT INTO system (key, value) VALUES (?, ?)`,
		migrateToV8SystemKeyBookMark, timestamp.Bookmark)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "inserting the bookmark value")
	}

	tx.Commit()

	if err := os.RemoveAll(dnoteFilePath); err != nil {
		return errors.Wrap(err, "removing the old dnote file")
	}
	if err := os.RemoveAll(actionsPath); err != nil {
		return errors.Wrap(err, "removing the actions file")
	}
	if err := os.RemoveAll(timestampsPath); err != nil {
		return errors.Wrap(err, "removing the timestamps file")
	}
	schemaPath := fmt.Sprintf("%s/schema", ctx.DnoteDir)
	if err := os.RemoveAll(schemaPath); err != nil {
		return errors.Wrap(err, "removing the schema file")
	}

	return nil
}
