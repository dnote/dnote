package core

import (
	"encoding/json"
	"testing"

	"github.com/dnote/actions"
	"github.com/dnote/cli/testutils"
	"github.com/pkg/errors"
)

func TestLogActionEditNote(t *testing.T) {
	// Setup
	ctx := testutils.InitEnv("../tmp", "../testutils/fixtures/schema.sql")
	defer testutils.TeardownEnv(ctx)

	// Execute
	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		panic(errors.Wrap(err, "beginning a transaction"))
	}

	if err := LogActionEditNote(tx, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "js", "updated content", 1536168581); err != nil {
		t.Fatalf("Failed to perform %s", err.Error())
	}

	tx.Commit()

	// Test
	var actionCount int
	if err := db.QueryRow("SELECT count(*) FROM actions;").Scan(&actionCount); err != nil {
		panic(errors.Wrap(err, "counting actions"))
	}
	var action actions.Action
	if err := db.QueryRow("SELECT uuid, schema, type, timestamp, data FROM actions").
		Scan(&action.UUID, &action.Schema, &action.Type, &action.Timestamp, &action.Data); err != nil {
		panic(errors.Wrap(err, "querying action"))
	}
	var actionData actions.EditNoteDataV2
	if err := json.Unmarshal(action.Data, &actionData); err != nil {
		panic(errors.Wrap(err, "unmarshalling action data"))
	}

	if actionCount != 1 {
		t.Fatalf("action count mismatch. got %d", actionCount)
	}
	testutils.AssertNotEqual(t, action.UUID, "", "action uuid mismatch")
	testutils.AssertEqual(t, action.Schema, 2, "action schema mismatch")
	testutils.AssertEqual(t, action.Type, actions.ActionEditNote, "action type mismatch")
	testutils.AssertNotEqual(t, action.Timestamp, 0, "action timestamp mismatch")
	testutils.AssertEqual(t, actionData.NoteUUID, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "action data note_uuid mismatch")
	testutils.AssertEqual(t, actionData.FromBook, "js", "action data from_book mismatch")
	testutils.AssertEqual(t, *actionData.Content, "updated content", "action data content mismatch")
	if actionData.ToBook != nil {
		t.Errorf("action data to_book mismatch. Expected %+v. Got %+v", nil, actionData.ToBook)
	}
	if actionData.Public != nil {
		t.Errorf("action data public mismatch. Expected %+v. Got %+v", nil, actionData.ToBook)
	}
}
