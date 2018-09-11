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
	ctx := testutils.InitCtx("../tmp")

	testutils.SetupTmp(ctx)
	testutils.SetupDB(ctx)
	defer testutils.ClearTmp(ctx)
	testutils.WriteFile(ctx, "../testutils/fixtures/dnote3.json", "dnote")
	InitFiles(ctx)

	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
	}

	if err := LogActionEditNote(tx, "f0d0fbb7-31ff-45ae-9f0f-4e429c0c797f", "js", "updated content", 1536168581); err != nil {
		t.Fatalf("Failed to perform %s", err.Error())
	}

	tx.Commit()

	b := testutils.ReadFile(ctx, "actions")
	var got []actions.Action

	if err := json.Unmarshal(b, &got); err != nil {
		panic(errors.Wrap(err, "unmarshalling actions"))
	}

	var actionData actions.EditNoteDataV2
	if err := json.Unmarshal(got[0].Data, &actionData); err != nil {
		panic(errors.Wrap(err, "unmarshalling action data"))
	}

	testutils.AssertEqual(t, len(got), 1, "action length mismatch")
	testutils.AssertNotEqual(t, got[0].UUID, "", "action uuid mismatch")
	testutils.AssertEqual(t, got[0].Schema, 2, "action schema mismatch")
	testutils.AssertEqual(t, got[0].Type, actions.ActionEditNote, "action type mismatch")
	testutils.AssertNotEqual(t, got[0].Timestamp, 0, "action timestamp mismatch")
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
