package migrate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnote/actions"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/testutils"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
)

func TestExecute_bump_schema(t *testing.T) {
	testCases := []struct {
		schemaKey string
	}{
		{
			schemaKey: infra.SystemSchema,
		},
		{
			schemaKey: infra.SystemRemoteSchema,
		},
	}

	for _, tc := range testCases {
		func() {
			// set up
			ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", false)
			defer testutils.TeardownEnv(ctx)

			db := ctx.DB
			testutils.MustExec(t, "inserting a schema", db, "INSERT INTO system (key, value) VALUES (?, ?)", tc.schemaKey, 8)

			m1 := migration{
				name: "noop",
				run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
					return nil
				},
			}
			m2 := migration{
				name: "noop",
				run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
					return nil
				},
			}

			// execute
			err := execute(ctx, m1, tc.schemaKey)
			if err != nil {
				t.Fatal(errors.Wrap(err, "failed to execute"))
			}
			err = execute(ctx, m2, tc.schemaKey)
			if err != nil {
				t.Fatal(errors.Wrap(err, "failed to execute"))
			}

			// test
			var schema int
			testutils.MustScan(t, "getting schema", db.QueryRow("SELECT value FROM system WHERE key = ?", tc.schemaKey), &schema)
			testutils.AssertEqual(t, schema, 10, "schema was not incremented properly")
		}()
	}
}

func TestRun_nonfresh(t *testing.T) {
	testCases := []struct {
		mode      int
		schemaKey string
	}{
		{
			mode:      LocalMode,
			schemaKey: infra.SystemSchema,
		},
		{
			mode:      RemoteMode,
			schemaKey: infra.SystemRemoteSchema,
		},
	}

	for _, tc := range testCases {
		func() {
			// set up
			ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", false)
			defer testutils.TeardownEnv(ctx)

			db := ctx.DB
			testutils.MustExec(t, "inserting a schema", db, "INSERT INTO system (key, value) VALUES (?, ?)", tc.schemaKey, 2)
			testutils.MustExec(t, "creating a temporary table for testing", db,
				"CREATE TABLE migrate_run_test ( name string )")

			sequence := []migration{
				migration{
					name: "v1",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v1 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v1")
						return nil
					},
				},
				migration{
					name: "v2",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v2 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v2")
						return nil
					},
				},
				migration{
					name: "v3",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v3 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v3")
						return nil
					},
				},
				migration{
					name: "v4",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v4 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v4")
						return nil
					},
				},
			}

			// execute
			err := Run(ctx, sequence, tc.mode)
			if err != nil {
				t.Fatal(errors.Wrap(err, "failed to run"))
			}

			// test
			var schema int
			testutils.MustScan(t, fmt.Sprintf("getting schema for %s", tc.schemaKey), db.QueryRow("SELECT value FROM system WHERE key = ?", tc.schemaKey), &schema)
			testutils.AssertEqual(t, schema, 4, fmt.Sprintf("schema was not updated for %s", tc.schemaKey))

			var testRunCount int
			testutils.MustScan(t, "counting test runs", db.QueryRow("SELECT count(*) FROM migrate_run_test"), &testRunCount)
			testutils.AssertEqual(t, testRunCount, 2, "test run count mismatch")

			var testRun1, testRun2 string
			testutils.MustScan(t, "finding test run 1", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v3"), &testRun1)
			testutils.MustScan(t, "finding test run 2", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v4"), &testRun2)
		}()
	}

}

func TestRun_fresh(t *testing.T) {
	testCases := []struct {
		mode      int
		schemaKey string
	}{
		{
			mode:      LocalMode,
			schemaKey: infra.SystemSchema,
		},
		{
			mode:      RemoteMode,
			schemaKey: infra.SystemRemoteSchema,
		},
	}

	for _, tc := range testCases {
		func() {
			// set up
			ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", false)
			defer testutils.TeardownEnv(ctx)

			db := ctx.DB
			testutils.MustExec(t, "creating a temporary table for testing", db,
				"CREATE TABLE migrate_run_test ( name string )")

			sequence := []migration{
				migration{
					name: "v1",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v1 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v1")
						return nil
					},
				},
				migration{
					name: "v2",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v2 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v2")
						return nil
					},
				},
				migration{
					name: "v3",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v3 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v3")
						return nil
					},
				},
			}

			// execute
			err := Run(ctx, sequence, tc.mode)
			if err != nil {
				t.Fatal(errors.Wrap(err, "failed to run"))
			}

			// test
			var schema int
			testutils.MustScan(t, "getting schema", db.QueryRow("SELECT value FROM system WHERE key = ?", tc.schemaKey), &schema)
			testutils.AssertEqual(t, schema, 3, "schema was not updated")

			var testRunCount int
			testutils.MustScan(t, "counting test runs", db.QueryRow("SELECT count(*) FROM migrate_run_test"), &testRunCount)
			testutils.AssertEqual(t, testRunCount, 3, "test run count mismatch")

			var testRun1, testRun2, testRun3 string
			testutils.MustScan(t, "finding test run 1", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v1"), &testRun1)
			testutils.MustScan(t, "finding test run 2", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v2"), &testRun2)
			testutils.MustScan(t, "finding test run 2", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v3"), &testRun3)
		}()
	}
}

func TestRun_up_to_date(t *testing.T) {
	testCases := []struct {
		mode      int
		schemaKey string
	}{
		{
			mode:      LocalMode,
			schemaKey: infra.SystemSchema,
		},
		{
			mode:      RemoteMode,
			schemaKey: infra.SystemRemoteSchema,
		},
	}

	for _, tc := range testCases {
		func() {
			// set up
			ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", false)
			defer testutils.TeardownEnv(ctx)

			db := ctx.DB
			testutils.MustExec(t, "creating a temporary table for testing", db,
				"CREATE TABLE migrate_run_test ( name string )")

			testutils.MustExec(t, "inserting a schema", db, "INSERT INTO system (key, value) VALUES (?, ?)", tc.schemaKey, 3)

			sequence := []migration{
				migration{
					name: "v1",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v1 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v1")
						return nil
					},
				},
				migration{
					name: "v2",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v2 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v2")
						return nil
					},
				},
				migration{
					name: "v3",
					run: func(ctx infra.DnoteCtx, tx *sql.Tx) error {
						testutils.MustExec(t, "marking v3 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v3")
						return nil
					},
				},
			}

			// execute
			err := Run(ctx, sequence, tc.mode)
			if err != nil {
				t.Fatal(errors.Wrap(err, "failed to run"))
			}

			// test
			var schema int
			testutils.MustScan(t, "getting schema", db.QueryRow("SELECT value FROM system WHERE key = ?", tc.schemaKey), &schema)
			testutils.AssertEqual(t, schema, 3, "schema was not updated")

			var testRunCount int
			testutils.MustScan(t, "counting test runs", db.QueryRow("SELECT count(*) FROM migrate_run_test"), &testRunCount)
			testutils.AssertEqual(t, testRunCount, 0, "test run count mismatch")
		}()
	}
}

func TestLocalMigration1(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-1-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	data := testutils.MustMarshalJSON(t, actions.AddBookDataV1{BookName: "js"})
	a1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a1UUID, 1, "add_book", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV1{NoteUUID: "note-1-uuid", FromBook: "js", ToBook: "", Content: "note 1"})
	a2UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a2UUID, 1, "edit_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV1{NoteUUID: "note-2-uuid", FromBook: "js", ToBook: "", Content: "note 2"})
	a3UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a3UUID, 1, "edit_note", string(data), 1537829463)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm1.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var actionCount int
	testutils.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	testutils.AssertEqual(t, actionCount, 3, "action count mismatch")

	var a1, a2, a3 actions.Action
	testutils.MustScan(t, "getting action 1", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a1UUID),
		&a1.Schema, &a1.Type, &a1.Data, &a1.Timestamp)
	testutils.MustScan(t, "getting action 2", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a2UUID),
		&a2.Schema, &a2.Type, &a2.Data, &a2.Timestamp)
	testutils.MustScan(t, "getting action 3", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a3UUID),
		&a3.Schema, &a3.Type, &a3.Data, &a3.Timestamp)

	var a1Data actions.AddBookDataV1
	var a2Data, a3Data actions.EditNoteDataV3
	testutils.MustUnmarshalJSON(t, a1.Data, &a1Data)
	testutils.MustUnmarshalJSON(t, a2.Data, &a2Data)
	testutils.MustUnmarshalJSON(t, a3.Data, &a3Data)

	testutils.AssertEqual(t, a1.Schema, 1, "a1 schema mismatch")
	testutils.AssertEqual(t, a1.Type, "add_book", "a1 type mismatch")
	testutils.AssertEqual(t, a1.Timestamp, int64(1537829463), "a1 timestamp mismatch")
	testutils.AssertEqual(t, a1Data.BookName, "js", "a1 data book_name mismatch")

	testutils.AssertEqual(t, a2.Schema, 3, "a2 schema mismatch")
	testutils.AssertEqual(t, a2.Type, "edit_note", "a2 type mismatch")
	testutils.AssertEqual(t, a2.Timestamp, int64(1537829463), "a2 timestamp mismatch")
	testutils.AssertEqual(t, a2Data.NoteUUID, "note-1-uuid", "a2 data note_uuid mismatch")
	testutils.AssertEqual(t, a2Data.BookName, (*string)(nil), "a2 data book_name mismatch")
	testutils.AssertEqual(t, *a2Data.Content, "note 1", "a2 data content mismatch")
	testutils.AssertEqual(t, *a2Data.Public, false, "a2 data public mismatch")

	testutils.AssertEqual(t, a3.Schema, 3, "a3 schema mismatch")
	testutils.AssertEqual(t, a3.Type, "edit_note", "a3 type mismatch")
	testutils.AssertEqual(t, a3.Timestamp, int64(1537829463), "a3 timestamp mismatch")
	testutils.AssertEqual(t, a3Data.NoteUUID, "note-2-uuid", "a3 data note_uuid mismatch")
	testutils.AssertEqual(t, a3Data.BookName, (*string)(nil), "a3 data book_name mismatch")
	testutils.AssertEqual(t, *a3Data.Content, "note 2", "a3 data content mismatch")
	testutils.AssertEqual(t, *a3Data.Public, false, "a3 data public mismatch")
}

func TestLocalMigration2(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-1-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	c1 := "note 1 - v1"
	c2 := "note 1 - v2"
	css := "css"

	b1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting css book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "css")

	data := testutils.MustMarshalJSON(t, actions.AddNoteDataV2{NoteUUID: "note-1-uuid", BookName: "js", Content: "note 1", Public: false})
	a1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a1UUID, 2, "add_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV2{NoteUUID: "note-1-uuid", FromBook: "js", ToBook: nil, Content: &c1, Public: nil})
	a2UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a2UUID, 2, "edit_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV2{NoteUUID: "note-1-uuid", FromBook: "js", ToBook: &css, Content: &c2, Public: nil})
	a3UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a3UUID, 2, "edit_note", string(data), 1537829463)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm2.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var actionCount int
	testutils.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	testutils.AssertEqual(t, actionCount, 3, "action count mismatch")

	var a1, a2, a3 actions.Action
	testutils.MustScan(t, "getting action 1", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a1UUID),
		&a1.Schema, &a1.Type, &a1.Data, &a1.Timestamp)
	testutils.MustScan(t, "getting action 2", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a2UUID),
		&a2.Schema, &a2.Type, &a2.Data, &a2.Timestamp)
	testutils.MustScan(t, "getting action 3", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a3UUID),
		&a3.Schema, &a3.Type, &a3.Data, &a3.Timestamp)

	var a1Data actions.AddNoteDataV2
	var a2Data, a3Data actions.EditNoteDataV3
	testutils.MustUnmarshalJSON(t, a1.Data, &a1Data)
	testutils.MustUnmarshalJSON(t, a2.Data, &a2Data)
	testutils.MustUnmarshalJSON(t, a3.Data, &a3Data)

	testutils.AssertEqual(t, a1.Schema, 2, "a1 schema mismatch")
	testutils.AssertEqual(t, a1.Type, "add_note", "a1 type mismatch")
	testutils.AssertEqual(t, a1.Timestamp, int64(1537829463), "a1 timestamp mismatch")
	testutils.AssertEqual(t, a1Data.NoteUUID, "note-1-uuid", "a1 data note_uuid mismatch")
	testutils.AssertEqual(t, a1Data.BookName, "js", "a1 data book_name mismatch")
	testutils.AssertEqual(t, a1Data.Public, false, "a1 data public mismatch")

	testutils.AssertEqual(t, a2.Schema, 3, "a2 schema mismatch")
	testutils.AssertEqual(t, a2.Type, "edit_note", "a2 type mismatch")
	testutils.AssertEqual(t, a2.Timestamp, int64(1537829463), "a2 timestamp mismatch")
	testutils.AssertEqual(t, a2Data.NoteUUID, "note-1-uuid", "a2 data note_uuid mismatch")
	testutils.AssertEqual(t, a2Data.BookName, (*string)(nil), "a2 data book_name mismatch")
	testutils.AssertEqual(t, *a2Data.Content, c1, "a2 data content mismatch")
	testutils.AssertEqual(t, a2Data.Public, (*bool)(nil), "a2 data public mismatch")

	testutils.AssertEqual(t, a3.Schema, 3, "a3 schema mismatch")
	testutils.AssertEqual(t, a3.Type, "edit_note", "a3 type mismatch")
	testutils.AssertEqual(t, a3.Timestamp, int64(1537829463), "a3 timestamp mismatch")
	testutils.AssertEqual(t, a3Data.NoteUUID, "note-1-uuid", "a3 data note_uuid mismatch")
	testutils.AssertEqual(t, *a3Data.BookName, "css", "a3 data book_name mismatch")
	testutils.AssertEqual(t, *a3Data.Content, c2, "a3 data content mismatch")
	testutils.AssertEqual(t, a3Data.Public, (*bool)(nil), "a3 data public mismatch")
}

func TestLocalMigration3(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-1-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	data := testutils.MustMarshalJSON(t, actions.AddNoteDataV2{NoteUUID: "note-1-uuid", BookName: "js", Content: "note 1", Public: false})
	a1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a1UUID, 2, "add_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.RemoveNoteDataV1{NoteUUID: "note-1-uuid", BookName: "js"})
	a2UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a2UUID, 1, "remove_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.RemoveNoteDataV1{NoteUUID: "note-2-uuid", BookName: "js"})
	a3UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a3UUID, 1, "remove_note", string(data), 1537829463)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm3.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var actionCount int
	testutils.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	testutils.AssertEqual(t, actionCount, 3, "action count mismatch")

	var a1, a2, a3 actions.Action
	testutils.MustScan(t, "getting action 1", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a1UUID),
		&a1.Schema, &a1.Type, &a1.Data, &a1.Timestamp)
	testutils.MustScan(t, "getting action 2", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a2UUID),
		&a2.Schema, &a2.Type, &a2.Data, &a2.Timestamp)
	testutils.MustScan(t, "getting action 3", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a3UUID),
		&a3.Schema, &a3.Type, &a3.Data, &a3.Timestamp)

	var a1Data actions.AddNoteDataV2
	var a2Data, a3Data actions.RemoveNoteDataV2
	testutils.MustUnmarshalJSON(t, a1.Data, &a1Data)
	testutils.MustUnmarshalJSON(t, a2.Data, &a2Data)
	testutils.MustUnmarshalJSON(t, a3.Data, &a3Data)

	testutils.AssertEqual(t, a1.Schema, 2, "a1 schema mismatch")
	testutils.AssertEqual(t, a1.Type, "add_note", "a1 type mismatch")
	testutils.AssertEqual(t, a1.Timestamp, int64(1537829463), "a1 timestamp mismatch")
	testutils.AssertEqual(t, a1Data.NoteUUID, "note-1-uuid", "a1 data note_uuid mismatch")
	testutils.AssertEqual(t, a1Data.BookName, "js", "a1 data book_name mismatch")
	testutils.AssertEqual(t, a1Data.Content, "note 1", "a1 data content mismatch")
	testutils.AssertEqual(t, a1Data.Public, false, "a1 data public mismatch")

	testutils.AssertEqual(t, a2.Schema, 2, "a2 schema mismatch")
	testutils.AssertEqual(t, a2.Type, "remove_note", "a2 type mismatch")
	testutils.AssertEqual(t, a2.Timestamp, int64(1537829463), "a2 timestamp mismatch")
	testutils.AssertEqual(t, a2Data.NoteUUID, "note-1-uuid", "a2 data note_uuid mismatch")

	testutils.AssertEqual(t, a3.Schema, 2, "a3 schema mismatch")
	testutils.AssertEqual(t, a3.Type, "remove_note", "a3 type mismatch")
	testutils.AssertEqual(t, a3.Timestamp, int64(1537829463), "a3 timestamp mismatch")
	testutils.AssertEqual(t, a3Data.NoteUUID, "note-2-uuid", "a3 data note_uuid mismatch")
}

func TestLocalMigration4(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-1-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting css book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "css")
	n1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting css note", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", n1UUID, b1UUID, "n1 content", time.Now().UnixNano())

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm4.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var n1Dirty, b1Dirty bool
	var n1Deleted, b1Deleted bool
	var n1USN, b1USN int
	testutils.MustScan(t, "scanning the newly added dirty flag of n1", db.QueryRow("SELECT dirty, deleted, usn FROM notes WHERE uuid = ?", n1UUID), &n1Dirty, &n1Deleted, &n1USN)
	testutils.MustScan(t, "scanning the newly added dirty flag of b1", db.QueryRow("SELECT dirty, deleted, usn FROM books WHERE uuid = ?", b1UUID), &b1Dirty, &b1Deleted, &b1USN)

	testutils.AssertEqual(t, n1Dirty, false, "n1 dirty flag should be false by default")
	testutils.AssertEqual(t, b1Dirty, false, "b1 dirty flag should be false by default")

	testutils.AssertEqual(t, n1Deleted, false, "n1 deleted flag should be false by default")
	testutils.AssertEqual(t, b1Deleted, false, "b1 deleted flag should be false by default")

	testutils.AssertEqual(t, n1USN, 0, "n1 usn flag should be 0 by default")
	testutils.AssertEqual(t, b1USN, 0, "b1 usn flag should be 0 by default")
}

func TestLocalMigration5(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-5-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting css book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "css")
	b2UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting js book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "js")

	n1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting css note", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", n1UUID, b1UUID, "n1 content", time.Now().UnixNano())
	n2UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting css note", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", n2UUID, b1UUID, "n2 content", time.Now().UnixNano())
	n3UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting css note", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", n3UUID, b1UUID, "n3 content", time.Now().UnixNano())

	data := testutils.MustMarshalJSON(t, actions.AddBookDataV1{BookName: "js"})
	testutils.MustExec(t, "inserting a1", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", "a1-uuid", 1, "add_book", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.AddNoteDataV2{NoteUUID: n1UUID, BookName: "css", Content: "n1 content", Public: false})
	testutils.MustExec(t, "inserting a2", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", "a2-uuid", 1, "add_note", string(data), 1537829463)

	updatedContent := "updated content"
	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV3{NoteUUID: n2UUID, BookName: (*string)(nil), Content: &updatedContent, Public: (*bool)(nil)})
	testutils.MustExec(t, "inserting a3", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", "a3-uuid", 1, "edit_note", string(data), 1537829463)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm5.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var b1Dirty, b2Dirty, n1Dirty, n2Dirty, n3Dirty bool
	testutils.MustScan(t, "scanning the newly added dirty flag of b1", db.QueryRow("SELECT dirty FROM books WHERE uuid = ?", b1UUID), &b1Dirty)
	testutils.MustScan(t, "scanning the newly added dirty flag of b2", db.QueryRow("SELECT dirty FROM books WHERE uuid = ?", b2UUID), &b2Dirty)
	testutils.MustScan(t, "scanning the newly added dirty flag of n1", db.QueryRow("SELECT dirty FROM notes WHERE uuid = ?", n1UUID), &n1Dirty)
	testutils.MustScan(t, "scanning the newly added dirty flag of n2", db.QueryRow("SELECT dirty FROM notes WHERE uuid = ?", n2UUID), &n2Dirty)
	testutils.MustScan(t, "scanning the newly added dirty flag of n3", db.QueryRow("SELECT dirty FROM notes WHERE uuid = ?", n3UUID), &n3Dirty)

	testutils.AssertEqual(t, b1Dirty, false, "b1 dirty flag should be false by default")
	testutils.AssertEqual(t, b2Dirty, true, "b2 dirty flag should be false by default")
	testutils.AssertEqual(t, n1Dirty, true, "n1 dirty flag should be false by default")
	testutils.AssertEqual(t, n2Dirty, true, "n2 dirty flag should be false by default")
	testutils.AssertEqual(t, n3Dirty, false, "n3 dirty flag should be false by default")
}

func TestLocalMigration6(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-5-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	data := testutils.MustMarshalJSON(t, actions.AddBookDataV1{BookName: "js"})
	a1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a1UUID, 1, "add_book", string(data), 1537829463)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm5.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var count int
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name = ?;", "actions").Scan(&count)
	testutils.AssertEqual(t, count, 0, "actions table should have been deleted")
}

func TestLocalMigration7_trash(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-7-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting trash book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "trash")

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm7.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var b1Label string
	var b1Dirty bool
	testutils.MustScan(t, "scanning b1 label", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b1UUID), &b1Label, &b1Dirty)
	testutils.AssertEqual(t, b1Label, "trash (2)", "b1 label was not migrated")
	testutils.AssertEqual(t, b1Dirty, true, "b1 was not marked dirty")
}

func TestLocalMigration7_conflicts(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-7-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "conflicts")

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm7.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var b1Label string
	var b1Dirty bool
	testutils.MustScan(t, "scanning b1 label", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b1UUID), &b1Label, &b1Dirty)
	testutils.AssertEqual(t, b1Label, "conflicts (2)", "b1 label was not migrated")
	testutils.AssertEqual(t, b1Dirty, true, "b1 was not marked dirty")
}

func TestLocalMigration7_conflicts_dup(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-7-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "conflicts")
	b2UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "conflicts (2)")

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm7.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var b1Label, b2Label string
	var b1Dirty, b2Dirty bool
	testutils.MustScan(t, "scanning b1 label", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b1UUID), &b1Label, &b1Dirty)
	testutils.MustScan(t, "scanning b2 label", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b2UUID), &b2Label, &b2Dirty)
	testutils.AssertEqual(t, b1Label, "conflicts (3)", "b1 label was not migrated")
	testutils.AssertEqual(t, b2Label, "conflicts (2)", "b1 label was not migrated")
	testutils.AssertEqual(t, b1Dirty, true, "b1 was not marked dirty")
	testutils.AssertEqual(t, b2Dirty, false, "b2 should not have been marked dirty")
}

func TestLocalMigration8(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-8-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1")

	n1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting n1", db, `INSERT INTO notes
		(id, uuid, book_uuid, content, added_on, edited_on, public, dirty, usn, deleted) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 1, n1UUID, b1UUID, "n1 Body", 1, 2, true, true, 20, false)
	n2UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting n2", db, `INSERT INTO notes
		(id, uuid, book_uuid, content, added_on, edited_on, public, dirty, usn, deleted) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 2, n2UUID, b1UUID, "", 3, 4, false, true, 21, true)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm8.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var n1BookUUID, n1Body string
	var n1AddedOn, n1EditedOn int64
	var n1USN int
	var n1Public, n1Dirty, n1Deleted bool
	testutils.MustScan(t, "scanning n1", db.QueryRow("SELECT book_uuid, body, added_on, edited_on, usn,  public, dirty, deleted FROM notes WHERE uuid = ?", n1UUID), &n1BookUUID, &n1Body, &n1AddedOn, &n1EditedOn, &n1USN, &n1Public, &n1Dirty, &n1Deleted)

	var n2BookUUID, n2Body string
	var n2AddedOn, n2EditedOn int64
	var n2USN int
	var n2Public, n2Dirty, n2Deleted bool
	testutils.MustScan(t, "scanning n2", db.QueryRow("SELECT book_uuid, body, added_on, edited_on, usn,  public, dirty, deleted FROM notes WHERE uuid = ?", n2UUID), &n2BookUUID, &n2Body, &n2AddedOn, &n2EditedOn, &n2USN, &n2Public, &n2Dirty, &n2Deleted)

	testutils.AssertEqual(t, n1BookUUID, b1UUID, "n1 BookUUID mismatch")
	testutils.AssertEqual(t, n1Body, "n1 Body", "n1 Body mismatch")
	testutils.AssertEqual(t, n1AddedOn, int64(1), "n1 AddedOn mismatch")
	testutils.AssertEqual(t, n1EditedOn, int64(2), "n1 EditedOn mismatch")
	testutils.AssertEqual(t, n1USN, 20, "n1 USN mismatch")
	testutils.AssertEqual(t, n1Public, true, "n1 Public mismatch")
	testutils.AssertEqual(t, n1Dirty, true, "n1 Dirty mismatch")
	testutils.AssertEqual(t, n1Deleted, false, "n1 Deleted mismatch")

	testutils.AssertEqual(t, n2BookUUID, b1UUID, "n2 BookUUID mismatch")
	testutils.AssertEqual(t, n2Body, "", "n2 Body mismatch")
	testutils.AssertEqual(t, n2AddedOn, int64(3), "n2 AddedOn mismatch")
	testutils.AssertEqual(t, n2EditedOn, int64(4), "n2 EditedOn mismatch")
	testutils.AssertEqual(t, n2USN, 21, "n2 USN mismatch")
	testutils.AssertEqual(t, n2Public, false, "n2 Public mismatch")
	testutils.AssertEqual(t, n2Dirty, true, "n2 Dirty mismatch")
	testutils.AssertEqual(t, n2Deleted, true, "n2 Deleted mismatch")
}

func TestLocalMigration9(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/local-9-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1")

	n1UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting n1", db, `INSERT INTO notes
		(uuid, book_uuid, body, added_on, edited_on, public, dirty, usn, deleted) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?)`, n1UUID, b1UUID, "n1 Body", 1, 2, true, true, 20, false)
	n2UUID := utils.GenerateUUID()
	testutils.MustExec(t, "inserting n2", db, `INSERT INTO notes
		(uuid, book_uuid, body, added_on, edited_on, public, dirty, usn, deleted) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?)`, n2UUID, b1UUID, "n2 Body", 3, 4, false, true, 21, false)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm9.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test

	// assert that note_fts was populated with correct values
	var noteFtsCount int
	testutils.MustScan(t, "counting note_fts", db.QueryRow("SELECT count(*) FROM note_fts;"), &noteFtsCount)
	testutils.AssertEqual(t, noteFtsCount, 2, "noteFtsCount mismatch")

	var resCount int
	testutils.MustScan(t, "counting result", db.QueryRow("SELECT count(*) FROM note_fts WHERE note_fts MATCH ?", "n1"), &resCount)
	testutils.AssertEqual(t, resCount, 1, "noteFtsCount mismatch")
}

func TestRemoteMigration1(t *testing.T) {
	// set up
	ctx := testutils.InitEnv(t, "../tmp", "./fixtures/remote-1-pre-schema.sql", false)
	defer testutils.TeardownEnv(ctx)

	JSBookUUID := "existing-js-book-uuid"
	CSSBookUUID := "existing-css-book-uuid"
	linuxBookUUID := "existing-linux-book-uuid"
	newJSBookUUID := "new-js-book-uuid"
	newCSSBookUUID := "new-css-book-uuid"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/v1/books" {
			res := []struct {
				UUID  string `json:"uuid"`
				Label string `json:"label"`
			}{
				{
					UUID:  newJSBookUUID,
					Label: "js",
				},
				{
					UUID:  newCSSBookUUID,
					Label: "css",
				},
				// book that only exists on the server. client must ignore.
				{
					UUID:  "golang-book-uuid",
					Label: "golang",
				},
			}

			if err := json.NewEncoder(w).Encode(res); err != nil {
				t.Fatal(errors.Wrap(err, "encoding response"))
			}
		}
	}))
	defer server.Close()

	ctx.APIEndpoint = server.URL

	db := ctx.DB

	testutils.MustExec(t, "inserting js book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", JSBookUUID, "js")
	testutils.MustExec(t, "inserting css book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", CSSBookUUID, "css")
	testutils.MustExec(t, "inserting linux book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", linuxBookUUID, "linux")
	testutils.MustExec(t, "inserting sessionKey", db, "INSERT INTO system (key, value) VALUES (?, ?)", infra.SystemSessionKey, "someSessionKey")
	testutils.MustExec(t, "inserting sessionKeyExpiry", db, "INSERT INTO system (key, value) VALUES (?, ?)", infra.SystemSessionKeyExpiry, time.Now().Add(24*time.Hour).Unix())

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = rm1.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// test
	var postJSBookUUID, postCSSBookUUID, postLinuxBookUUID string
	testutils.MustScan(t, "getting js book uuid", db.QueryRow("SELECT uuid FROM books WHERE label = ?", "js"), &postJSBookUUID)
	testutils.MustScan(t, "getting css book uuid", db.QueryRow("SELECT uuid FROM books WHERE label = ?", "css"), &postCSSBookUUID)
	testutils.MustScan(t, "getting linux book uuid", db.QueryRow("SELECT uuid FROM books WHERE label = ?", "linux"), &postLinuxBookUUID)

	testutils.AssertEqual(t, postJSBookUUID, newJSBookUUID, "js book uuid was not updated correctly")
	testutils.AssertEqual(t, postCSSBookUUID, newCSSBookUUID, "css book uuid was not updated correctly")
	testutils.AssertEqual(t, postLinuxBookUUID, linuxBookUUID, "linux book uuid changed")
}
