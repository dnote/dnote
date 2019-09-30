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

package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/testutils"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func setupEnv(t *testing.T, homeDir string) context.DnoteCtx {
	dnoteDir := fmt.Sprintf("%s/.dnote", homeDir)
	if err := os.MkdirAll(dnoteDir, 0755); err != nil {
		t.Fatal(errors.Wrap(err, "preparing dnote dir"))
	}

	return context.DnoteCtx{
		HomeDir:  homeDir,
		DnoteDir: dnoteDir,
	}
}

func teardownEnv(t *testing.T, ctx context.DnoteCtx) {
	if err := os.RemoveAll(ctx.DnoteDir); err != nil {
		t.Fatal(errors.Wrap(err, "tearing down the dnote dir"))
	}
}

func TestMigrateToV1(t *testing.T) {
	t.Run("yaml exists", func(t *testing.T) {
		// set up
		ctx := setupEnv(t, "../tmp")
		defer teardownEnv(t, ctx)

		yamlPath, err := filepath.Abs(filepath.Join(ctx.HomeDir, ".dnote-yaml-archived"))
		if err != nil {
			panic(errors.Wrap(err, "Failed to get absolute YAML path").Error())
		}
		ioutil.WriteFile(yamlPath, []byte{}, 0644)

		// execute
		if err := migrateToV1(ctx); err != nil {
			t.Fatal(errors.Wrapf(err, "Failed to migrate").Error())
		}

		// test
		ok, err := utils.FileExists(yamlPath)
		if err != nil {
			t.Fatal(errors.Wrap(err, "checking if yaml file exists"))
		}
		if ok {
			t.Fatal("YAML archive file has not been deleted")
		}
	})

	t.Run("yaml does not exist", func(t *testing.T) {
		// set up
		ctx := setupEnv(t, "../tmp")
		defer teardownEnv(t, ctx)

		yamlPath, err := filepath.Abs(filepath.Join(ctx.HomeDir, ".dnote-yaml-archived"))
		if err != nil {
			panic(errors.Wrap(err, "Failed to get absolute YAML path").Error())
		}

		// execute
		if err := migrateToV1(ctx); err != nil {
			t.Fatal(errors.Wrapf(err, "Failed to migrate").Error())
		}

		// test
		ok, err := utils.FileExists(yamlPath)
		if err != nil {
			t.Fatal(errors.Wrap(err, "checking if yaml file exists"))
		}
		if ok {
			t.Fatal("YAML archive file must not exist")
		}
	})
}

func TestMigrateToV2(t *testing.T) {
	ctx := setupEnv(t, "../tmp")
	defer teardownEnv(t, ctx)

	testutils.CopyFixture(t, ctx, "./fixtures/legacy-2-pre-dnote.json", "dnote")

	// execute
	if err := migrateToV2(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := testutils.ReadFile(ctx, "dnote")

	var postDnote migrateToV2PostDnote
	if err := json.Unmarshal(b, &postDnote); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	for _, book := range postDnote {
		assert.NotEqual(t, book.Name, "", "Book name was not populated")

		for _, note := range book.Notes {
			if len(note.UUID) == 8 {
				t.Errorf("Note UUID was not migrated. It has length of %d", len(note.UUID))
			}

			assert.NotEqual(t, note.AddedOn, int64(0), "AddedOn was not carried over")
			assert.Equal(t, note.EditedOn, int64(0), "EditedOn was not created properly")
		}
	}
}

func TestMigrateToV3(t *testing.T) {
	// set up
	ctx := setupEnv(t, "../tmp")
	defer teardownEnv(t, ctx)

	testutils.CopyFixture(t, ctx, "./fixtures/legacy-3-pre-dnote.json", "dnote")

	// execute
	if err := migrateToV3(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := testutils.ReadFile(ctx, "dnote")
	var postDnote migrateToV3Dnote
	if err := json.Unmarshal(b, &postDnote); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	b = testutils.ReadFile(ctx, "actions")
	var actions []migrateToV3Action
	if err := json.Unmarshal(b, &actions); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the actions").Error())
	}

	assert.Equal(t, len(actions), 6, "actions length mismatch")

	for _, book := range postDnote {
		for _, note := range book.Notes {
			assert.NotEqual(t, note.AddedOn, int64(0), "AddedOn was not carried over")
		}
	}
}

func TestMigrateToV4(t *testing.T) {
	// set up
	ctx := setupEnv(t, "../tmp")
	defer teardownEnv(t, ctx)
	defer os.Setenv("EDITOR", "")

	testutils.CopyFixture(t, ctx, "./fixtures/legacy-4-pre-dnoterc.yaml", "dnoterc")

	// execute
	os.Setenv("EDITOR", "vim")
	if err := migrateToV4(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := testutils.ReadFile(ctx, "dnoterc")
	var config migrateToV4PostConfig
	if err := yaml.Unmarshal(b, &config); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	assert.Equal(t, config.APIKey, "Oev6e1082ORasdf9rjkfjkasdfjhgei", "api key mismatch")
	assert.Equal(t, config.Editor, "vim", "editor mismatch")
}

func TestMigrateToV5(t *testing.T) {
	// set up
	ctx := setupEnv(t, "../tmp")
	defer teardownEnv(t, ctx)

	testutils.CopyFixture(t, ctx, "./fixtures/legacy-5-pre-actions.json", "actions")

	// execute
	if err := migrateToV5(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "migrating").Error())
	}

	// test
	var oldActions []migrateToV5PreAction
	testutils.ReadJSON("./fixtures/legacy-5-pre-actions.json", &oldActions)

	b := testutils.ReadFile(ctx, "actions")
	var migratedActions []migrateToV5PostAction
	if err := json.Unmarshal(b, &migratedActions); err != nil {
		t.Fatal(errors.Wrap(err, "unmarshalling migrated actions").Error())
	}

	if len(oldActions) != len(migratedActions) {
		t.Fatalf("There were %d actions but after migration there were %d", len(oldActions), len(migratedActions))
	}

	for idx := range migratedActions {
		migrated := migratedActions[idx]
		old := oldActions[idx]

		assert.NotEqual(t, migrated.UUID, "", fmt.Sprintf("uuid mismatch for migrated item with index %d", idx))
		assert.Equal(t, migrated.Schema, 1, fmt.Sprintf("schema mismatch for migrated item with index %d", idx))
		assert.Equal(t, migrated.Timestamp, old.Timestamp, fmt.Sprintf("timestamp mismatch for migrated item with index %d", idx))
		assert.Equal(t, migrated.Type, old.Type, fmt.Sprintf("timestamp mismatch for migrated item with index %d", idx))

		switch migrated.Type {
		case migrateToV5ActionAddNote:
			var oldData, migratedData migrateToV5AddNoteData
			if err := json.Unmarshal(old.Data, &oldData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling old data").Error())
			}
			if err := json.Unmarshal(migrated.Data, &migratedData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling new data").Error())
			}

			assert.Equal(t, oldData.BookName, migratedData.BookName, fmt.Sprintf("data book_name mismatch for item idx %d", idx))
			assert.Equal(t, oldData.Content, migratedData.Content, fmt.Sprintf("data content mismatch for item idx %d", idx))
			assert.Equal(t, oldData.NoteUUID, migratedData.NoteUUID, fmt.Sprintf("data note_uuid mismatch for item idx %d", idx))
		case migrateToV5ActionRemoveNote:
			var oldData, migratedData migrateToV5RemoveNoteData
			if err := json.Unmarshal(old.Data, &oldData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling old data").Error())
			}
			if err := json.Unmarshal(migrated.Data, &migratedData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling new data").Error())
			}

			assert.Equal(t, oldData.BookName, migratedData.BookName, fmt.Sprintf("data book_name mismatch for item idx %d", idx))
			assert.Equal(t, oldData.NoteUUID, migratedData.NoteUUID, fmt.Sprintf("data note_uuid mismatch for item idx %d", idx))
		case migrateToV5ActionAddBook:
			var oldData, migratedData migrateToV5AddBookData
			if err := json.Unmarshal(old.Data, &oldData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling old data").Error())
			}
			if err := json.Unmarshal(migrated.Data, &migratedData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling new data").Error())
			}

			assert.Equal(t, oldData.BookName, migratedData.BookName, fmt.Sprintf("data book_name mismatch for item idx %d", idx))
		case migrateToV5ActionRemoveBook:
			var oldData, migratedData migrateToV5RemoveBookData
			if err := json.Unmarshal(old.Data, &oldData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling old data").Error())
			}
			if err := json.Unmarshal(migrated.Data, &migratedData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling new data").Error())
			}

			assert.Equal(t, oldData.BookName, migratedData.BookName, fmt.Sprintf("data book_name mismatch for item idx %d", idx))
		case migrateToV5ActionEditNote:
			var oldData migrateToV5PreEditNoteData
			var migratedData migrateToV5PostEditNoteData
			if err := json.Unmarshal(old.Data, &oldData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling old data").Error())
			}
			if err := json.Unmarshal(migrated.Data, &migratedData); err != nil {
				t.Fatal(errors.Wrap(err, "unmarshalling new data").Error())
			}

			assert.Equal(t, oldData.NoteUUID, migratedData.NoteUUID, fmt.Sprintf("data note_uuid mismatch for item idx %d", idx))
			assert.Equal(t, oldData.Content, migratedData.Content, fmt.Sprintf("data content mismatch for item idx %d", idx))
			assert.Equal(t, oldData.BookName, migratedData.FromBook, "book_name should have been renamed to from_book")
			assert.Equal(t, migratedData.ToBook, "", "to_book should be empty")
		}
	}
}

func TestMigrateToV6(t *testing.T) {
	// set up
	ctx := setupEnv(t, "../tmp")
	defer teardownEnv(t, ctx)

	testutils.CopyFixture(t, ctx, "./fixtures/legacy-6-pre-dnote.json", "dnote")

	// execute
	if err := migrateToV6(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := testutils.ReadFile(ctx, "dnote")
	var got migrateToV6PostDnote
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	b = utils.ReadFileAbs("./fixtures/legacy-6-post-dnote.json")
	var expected migrateToV6PostDnote
	if err := json.Unmarshal(b, &expected); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	if ok := reflect.DeepEqual(expected, got); !ok {
		t.Errorf("Payload does not match.\nActual:   %+v\nExpected: %+v", got, expected)
	}
}

func TestMigrateToV7(t *testing.T) {
	// set up
	ctx := setupEnv(t, "../tmp")
	defer teardownEnv(t, ctx)

	testutils.CopyFixture(t, ctx, "./fixtures/legacy-7-pre-actions.json", "actions")

	// execute
	if err := migrateToV7(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "migrating").Error())
	}

	// test
	b := testutils.ReadFile(ctx, "actions")
	var got []migrateToV7Action
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(errors.Wrap(err, "unmarshalling the result").Error())
	}

	b2 := utils.ReadFileAbs("./fixtures/legacy-7-post-actions.json")
	var expected []migrateToV7Action
	if err := json.Unmarshal(b, &expected); err != nil {
		t.Fatal(errors.Wrap(err, "unmarshalling the result into Dnote").Error())
	}

	assert.EqualJSON(t, string(b), string(b2), "Result does not match")
}

func TestMigrateToV8(t *testing.T) {
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-1-pre-schema.sql", SkipMigration: true}
	db := database.InitTestDB(t, "../tmp/.dnote/dnote-test.db", &opts)
	defer database.CloseTestDB(t, db)

	ctx := context.DnoteCtx{HomeDir: "../tmp", DnoteDir: "../tmp/.dnote", DB: db}

	// set up
	testutils.CopyFixture(t, ctx, "./fixtures/legacy-8-actions.json", "actions")
	testutils.CopyFixture(t, ctx, "./fixtures/legacy-8-dnote.json", "dnote")
	testutils.CopyFixture(t, ctx, "./fixtures/legacy-8-dnoterc.yaml", "dnoterc")
	testutils.CopyFixture(t, ctx, "./fixtures/legacy-8-schema.yaml", "schema")
	testutils.CopyFixture(t, ctx, "./fixtures/legacy-8-timestamps.yaml", "timestamps")

	// execute
	if err := migrateToV8(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "migrating").Error())
	}

	// test

	// 1. test if files are migrated
	dnoteFilePath := fmt.Sprintf("%s/dnote", ctx.DnoteDir)
	dnotercPath := fmt.Sprintf("%s/dnoterc", ctx.DnoteDir)
	schemaFilePath := fmt.Sprintf("%s/schema", ctx.DnoteDir)
	timestampFilePath := fmt.Sprintf("%s/timestamps", ctx.DnoteDir)

	ok, err := utils.FileExists(dnoteFilePath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "checking if file exists"))
	}
	if ok {
		t.Errorf("%s still exists", dnoteFilePath)
	}

	ok, err = utils.FileExists(schemaFilePath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "checking if file exists"))
	}
	if ok {
		t.Errorf("%s still exists", schemaFilePath)
	}

	ok, err = utils.FileExists(timestampFilePath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "checking if file exists"))
	}
	if ok {
		t.Errorf("%s still exists", timestampFilePath)
	}

	ok, err = utils.FileExists(dnotercPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "checking if file exists"))
	}
	if !ok {
		t.Errorf("%s still exists", dnotercPath)
	}

	// 2. test if notes and books are migrated

	var bookCount, noteCount int
	err = db.QueryRow("SELECT count(*) FROM books").Scan(&bookCount)
	if err != nil {
		panic(errors.Wrap(err, "counting books"))
	}
	err = db.QueryRow("SELECT count(*) FROM notes").Scan(&noteCount)
	if err != nil {
		panic(errors.Wrap(err, "counting notes"))
	}
	assert.Equal(t, bookCount, 2, "book count mismatch")
	assert.Equal(t, noteCount, 3, "note count mismatch")

	type bookInfo struct {
		label string
		uuid  string
	}
	type noteInfo struct {
		id       int
		uuid     string
		bookUUID string
		content  string
		addedOn  int64
		editedOn int64
		public   bool
	}

	var b1, b2 bookInfo
	var n1, n2, n3 noteInfo
	err = db.QueryRow("SELECT label, uuid FROM books WHERE label = ?", "js").Scan(&b1.label, &b1.uuid)
	if err != nil {
		panic(errors.Wrap(err, "finding book 1"))
	}
	err = db.QueryRow("SELECT label, uuid FROM books WHERE label = ?", "css").Scan(&b2.label, &b2.uuid)
	if err != nil {
		panic(errors.Wrap(err, "finding book 2"))
	}
	err = db.QueryRow("SELECT id, uuid, book_uuid, content, added_on, edited_on, public FROM notes WHERE uuid = ?", "d69edb54-5b31-4cdd-a4a5-34f0a0bfa153").Scan(&n1.id, &n1.uuid, &n1.bookUUID, &n1.content, &n1.addedOn, &n1.editedOn, &n1.public)
	if err != nil {
		panic(errors.Wrap(err, "finding note 1"))
	}
	err = db.QueryRow("SELECT id, uuid, book_uuid, content, added_on, edited_on, public FROM notes WHERE uuid = ?", "35cbcab1-6a2a-4cc8-97e0-e73bbbd54626").Scan(&n2.id, &n2.uuid, &n2.bookUUID, &n2.content, &n2.addedOn, &n2.editedOn, &n2.public)
	if err != nil {
		panic(errors.Wrap(err, "finding note 2"))
	}
	err = db.QueryRow("SELECT id, uuid, book_uuid, content, added_on, edited_on, public FROM notes WHERE uuid = ?", "7c1fcfb2-de8b-4350-88f0-fb3cbaf6630a").Scan(&n3.id, &n3.uuid, &n3.bookUUID, &n3.content, &n3.addedOn, &n3.editedOn, &n3.public)
	if err != nil {
		panic(errors.Wrap(err, "finding note 3"))
	}

	assert.NotEqual(t, b1.uuid, "", "book 1 uuid should have been generated")
	assert.Equal(t, b1.label, "js", "book 1 label mismatch")
	assert.NotEqual(t, b2.uuid, "", "book 2 uuid should have been generated")
	assert.Equal(t, b2.label, "css", "book 2 label mismatch")

	assert.Equal(t, n1.uuid, "d69edb54-5b31-4cdd-a4a5-34f0a0bfa153", "note 1 uuid mismatch")
	assert.NotEqual(t, n1.id, 0, "note 1 id should have been generated")
	assert.Equal(t, n1.bookUUID, b2.uuid, "note 1 book_uuid mismatch")
	assert.Equal(t, n1.content, "css test 1", "note 1 content mismatch")
	assert.Equal(t, n1.addedOn, int64(1536977237), "note 1 added_on mismatch")
	assert.Equal(t, n1.editedOn, int64(1536977253), "note 1 edited_on mismatch")
	assert.Equal(t, n1.public, false, "note 1 public mismatch")

	assert.Equal(t, n2.uuid, "35cbcab1-6a2a-4cc8-97e0-e73bbbd54626", "note 2 uuid mismatch")
	assert.NotEqual(t, n2.id, 0, "note 2 id should have been generated")
	assert.Equal(t, n2.bookUUID, b1.uuid, "note 2 book_uuid mismatch")
	assert.Equal(t, n2.content, "js test 1", "note 2 content mismatch")
	assert.Equal(t, n2.addedOn, int64(1536977229), "note 2 added_on mismatch")
	assert.Equal(t, n2.editedOn, int64(0), "note 2 edited_on mismatch")
	assert.Equal(t, n2.public, false, "note 2 public mismatch")

	assert.Equal(t, n3.uuid, "7c1fcfb2-de8b-4350-88f0-fb3cbaf6630a", "note 3 uuid mismatch")
	assert.NotEqual(t, n3.id, 0, "note 3 id should have been generated")
	assert.Equal(t, n3.bookUUID, b1.uuid, "note 3 book_uuid mismatch")
	assert.Equal(t, n3.content, "js test 2", "note 3 content mismatch")
	assert.Equal(t, n3.addedOn, int64(1536977230), "note 3 added_on mismatch")
	assert.Equal(t, n3.editedOn, int64(0), "note 3 edited_on mismatch")
	assert.Equal(t, n3.public, false, "note 3 public mismatch")

	// 3. test if actions are migrated
	var actionCount int
	err = db.QueryRow("SELECT count(*) FROM actions").Scan(&actionCount)
	if err != nil {
		panic(errors.Wrap(err, "counting actions"))
	}

	assert.Equal(t, actionCount, 11, "action count mismatch")

	type actionInfo struct {
		uuid       string
		schema     int
		actionType string
		data       string
		timestamp  int
	}

	var a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11 actionInfo
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "6145c1b7-f286-4d9f-b0f6-00d274baefc6").Scan(&a1.uuid, &a1.schema, &a1.actionType, &a1.data, &a1.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a1"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "c048a56b-179c-4f31-9995-81e9b32b7dd6").Scan(&a2.uuid, &a2.schema, &a2.actionType, &a2.data, &a2.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a2"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "f557ef48-c304-47dc-adfb-46b7306e701f").Scan(&a3.uuid, &a3.schema, &a3.actionType, &a3.data, &a3.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a3"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "8d79db34-343d-4331-ae5b-24743f17ca7f").Scan(&a4.uuid, &a4.schema, &a4.actionType, &a4.data, &a4.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a4"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "b9c1ed4a-e6b3-41f2-983b-593ec7b8b7a1").Scan(&a5.uuid, &a5.schema, &a5.actionType, &a5.data, &a5.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a5"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "06ed7ef0-f171-4bd7-ae8e-97b5d06a4c49").Scan(&a6.uuid, &a6.schema, &a6.actionType, &a6.data, &a6.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a6"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "7f173cef-1688-4177-a373-145fcd822b2f").Scan(&a7.uuid, &a7.schema, &a7.actionType, &a7.data, &a7.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a7"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "64352e08-aa7a-45f4-b760-b3f38b5e11fa").Scan(&a8.uuid, &a8.schema, &a8.actionType, &a8.data, &a8.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a8"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "82e20a12-bda8-45f7-ac42-b453b6daa5ec").Scan(&a9.uuid, &a9.schema, &a9.actionType, &a9.data, &a9.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a9"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "a29055f4-ace4-44fd-8800-3396edbccaef").Scan(&a10.uuid, &a10.schema, &a10.actionType, &a10.data, &a10.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a10"))
	}
	err = db.QueryRow("SELECT uuid, schema, type, data, timestamp FROM actions WHERE uuid = ?", "871a5562-1bd0-43c1-b550-5bbb727ac7c4").Scan(&a11.uuid, &a11.schema, &a11.actionType, &a11.data, &a11.timestamp)
	if err != nil {
		panic(errors.Wrap(err, "finding a11"))
	}

	assert.Equal(t, a1.uuid, "6145c1b7-f286-4d9f-b0f6-00d274baefc6", "action 1 uuid mismatch")
	assert.Equal(t, a1.schema, 1, "action 1 schema mismatch")
	assert.Equal(t, a1.actionType, "add_book", "action 1 type mismatch")
	assert.Equal(t, a1.data, `{"book_name":"js"}`, "action 1 data mismatch")
	assert.Equal(t, a1.timestamp, 1536977229, "action 1 timestamp mismatch")

	assert.Equal(t, a2.uuid, "c048a56b-179c-4f31-9995-81e9b32b7dd6", "action 2 uuid mismatch")
	assert.Equal(t, a2.schema, 2, "action 2 schema mismatch")
	assert.Equal(t, a2.actionType, "add_note", "action 2 type mismatch")
	assert.Equal(t, a2.data, `{"note_uuid":"35cbcab1-6a2a-4cc8-97e0-e73bbbd54626","book_name":"js","content":"js test 1","public":false}`, "action 2 data mismatch")
	assert.Equal(t, a2.timestamp, 1536977229, "action 2 timestamp mismatch")

	assert.Equal(t, a3.uuid, "f557ef48-c304-47dc-adfb-46b7306e701f", "action 3 uuid mismatch")
	assert.Equal(t, a3.schema, 2, "action 3 schema mismatch")
	assert.Equal(t, a3.actionType, "add_note", "action 3 type mismatch")
	assert.Equal(t, a3.data, `{"note_uuid":"7c1fcfb2-de8b-4350-88f0-fb3cbaf6630a","book_name":"js","content":"js test 2","public":false}`, "action 3 data mismatch")
	assert.Equal(t, a3.timestamp, 1536977230, "action 3 timestamp mismatch")

	assert.Equal(t, a4.uuid, "8d79db34-343d-4331-ae5b-24743f17ca7f", "action 4 uuid mismatch")
	assert.Equal(t, a4.schema, 2, "action 4 schema mismatch")
	assert.Equal(t, a4.actionType, "add_note", "action 4 type mismatch")
	assert.Equal(t, a4.data, `{"note_uuid":"b23a88ba-b291-4294-9795-86b394db5dcf","book_name":"js","content":"js test 3","public":false}`, "action 4 data mismatch")
	assert.Equal(t, a4.timestamp, 1536977234, "action 4 timestamp mismatch")

	assert.Equal(t, a5.uuid, "b9c1ed4a-e6b3-41f2-983b-593ec7b8b7a1", "action 5 uuid mismatch")
	assert.Equal(t, a5.schema, 1, "action 5 schema mismatch")
	assert.Equal(t, a5.actionType, "add_book", "action 5 type mismatch")
	assert.Equal(t, a5.data, `{"book_name":"css"}`, "action 5 data mismatch")
	assert.Equal(t, a5.timestamp, 1536977237, "action 5 timestamp mismatch")

	assert.Equal(t, a6.uuid, "06ed7ef0-f171-4bd7-ae8e-97b5d06a4c49", "action 6 uuid mismatch")
	assert.Equal(t, a6.schema, 2, "action 6 schema mismatch")
	assert.Equal(t, a6.actionType, "add_note", "action 6 type mismatch")
	assert.Equal(t, a6.data, `{"note_uuid":"d69edb54-5b31-4cdd-a4a5-34f0a0bfa153","book_name":"css","content":"js test 3","public":false}`, "action 6 data mismatch")
	assert.Equal(t, a6.timestamp, 1536977237, "action 6 timestamp mismatch")

	assert.Equal(t, a7.uuid, "7f173cef-1688-4177-a373-145fcd822b2f", "action 7 uuid mismatch")
	assert.Equal(t, a7.schema, 2, "action 7 schema mismatch")
	assert.Equal(t, a7.actionType, "edit_note", "action 7 type mismatch")
	assert.Equal(t, a7.data, `{"note_uuid":"d69edb54-5b31-4cdd-a4a5-34f0a0bfa153","from_book":"css","to_book":null,"content":"css test 1","public":null}`, "action 7 data mismatch")
	assert.Equal(t, a7.timestamp, 1536977253, "action 7 timestamp mismatch")

	assert.Equal(t, a8.uuid, "64352e08-aa7a-45f4-b760-b3f38b5e11fa", "action 8 uuid mismatch")
	assert.Equal(t, a8.schema, 1, "action 8 schema mismatch")
	assert.Equal(t, a8.actionType, "add_book", "action 8 type mismatch")
	assert.Equal(t, a8.data, `{"book_name":"sql"}`, "action 8 data mismatch")
	assert.Equal(t, a8.timestamp, 1536977261, "action 8 timestamp mismatch")

	assert.Equal(t, a9.uuid, "82e20a12-bda8-45f7-ac42-b453b6daa5ec", "action 9 uuid mismatch")
	assert.Equal(t, a9.schema, 2, "action 9 schema mismatch")
	assert.Equal(t, a9.actionType, "add_note", "action 9 type mismatch")
	assert.Equal(t, a9.data, `{"note_uuid":"2f47d390-685b-4b84-89ac-704c6fb8d3fb","book_name":"sql","content":"blah","public":false}`, "action 9 data mismatch")
	assert.Equal(t, a9.timestamp, 1536977261, "action 9 timestamp mismatch")

	assert.Equal(t, a10.uuid, "a29055f4-ace4-44fd-8800-3396edbccaef", "action 10 uuid mismatch")
	assert.Equal(t, a10.schema, 1, "action 10 schema mismatch")
	assert.Equal(t, a10.actionType, "remove_book", "action 10 type mismatch")
	assert.Equal(t, a10.data, `{"book_name":"sql"}`, "action 10 data mismatch")
	assert.Equal(t, a10.timestamp, 1536977268, "action 10 timestamp mismatch")

	assert.Equal(t, a11.uuid, "871a5562-1bd0-43c1-b550-5bbb727ac7c4", "action 11 uuid mismatch")
	assert.Equal(t, a11.schema, 1, "action 11 schema mismatch")
	assert.Equal(t, a11.actionType, "remove_note", "action 11 type mismatch")
	assert.Equal(t, a11.data, `{"note_uuid":"b23a88ba-b291-4294-9795-86b394db5dcf","book_name":"js"}`, "action 11 data mismatch")
	assert.Equal(t, a11.timestamp, 1536977274, "action 11 timestamp mismatch")

	// 3. test if system is migrated
	var systemCount int
	err = db.QueryRow("SELECT count(*) FROM system").Scan(&systemCount)
	if err != nil {
		panic(errors.Wrap(err, "counting system"))
	}

	assert.Equal(t, systemCount, 3, "action count mismatch")

	var lastUpgrade, lastAction, bookmark int
	err = db.QueryRow("SELECT value FROM system WHERE key = ?", "last_upgrade").Scan(&lastUpgrade)
	if err != nil {
		panic(errors.Wrap(err, "finding last_upgrade"))
	}
	err = db.QueryRow("SELECT value FROM system WHERE key = ?", "last_action").Scan(&lastAction)
	if err != nil {
		panic(errors.Wrap(err, "finding last_action"))
	}
	err = db.QueryRow("SELECT value FROM system WHERE key = ?", "bookmark").Scan(&bookmark)
	if err != nil {
		panic(errors.Wrap(err, "finding bookmark"))
	}

	assert.Equal(t, lastUpgrade, 1536977220, "last_upgrade mismatch")
	assert.Equal(t, lastAction, 1536977274, "last_action mismatch")
	assert.Equal(t, bookmark, 9, "bookmark mismatch")
}
