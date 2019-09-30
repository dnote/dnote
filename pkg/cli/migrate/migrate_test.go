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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnote/actions"
	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/testutils"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
)

func TestExecute_bump_schema(t *testing.T) {
	testCases := []struct {
		schemaKey string
	}{
		{
			schemaKey: consts.SystemSchema,
		},
		{
			schemaKey: consts.SystemRemoteSchema,
		},
	}

	for _, tc := range testCases {
		func() {
			// set up
			opts := database.TestDBOptions{SkipMigration: true}
			ctx := context.InitTestCtx(t, "../tmp", &opts)
			defer context.TeardownTestCtx(t, ctx)

			db := ctx.DB

			database.MustExec(t, "inserting a schema", db, "INSERT INTO system (key, value) VALUES (?, ?)", tc.schemaKey, 8)

			m1 := migration{
				name: "noop",
				run: func(ctx context.DnoteCtx, db *database.DB) error {
					return nil
				},
			}
			m2 := migration{
				name: "noop",
				run: func(ctx context.DnoteCtx, db *database.DB) error {
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
			database.MustScan(t, "getting schema", db.QueryRow("SELECT value FROM system WHERE key = ?", tc.schemaKey), &schema)
			assert.Equal(t, schema, 10, "schema was not incremented properly")
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
			schemaKey: consts.SystemSchema,
		},
		{
			mode:      RemoteMode,
			schemaKey: consts.SystemRemoteSchema,
		},
	}

	for _, tc := range testCases {
		func() {
			// set up
			opts := database.TestDBOptions{SkipMigration: true}
			ctx := context.InitTestCtx(t, "../tmp", &opts)
			defer context.TeardownTestCtx(t, ctx)

			db := ctx.DB
			database.MustExec(t, "inserting a schema", db, "INSERT INTO system (key, value) VALUES (?, ?)", tc.schemaKey, 2)
			database.MustExec(t, "creating a temporary table for testing", db,
				"CREATE TABLE migrate_run_test ( name string )")

			sequence := []migration{
				{
					name: "v1",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v1 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v1")
						return nil
					},
				},
				{
					name: "v2",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v2 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v2")
						return nil
					},
				},
				{
					name: "v3",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v3 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v3")
						return nil
					},
				},
				{
					name: "v4",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v4 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v4")
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
			database.MustScan(t, fmt.Sprintf("getting schema for %s", tc.schemaKey), db.QueryRow("SELECT value FROM system WHERE key = ?", tc.schemaKey), &schema)
			assert.Equal(t, schema, 4, fmt.Sprintf("schema was not updated for %s", tc.schemaKey))

			var testRunCount int
			database.MustScan(t, "counting test runs", db.QueryRow("SELECT count(*) FROM migrate_run_test"), &testRunCount)
			assert.Equal(t, testRunCount, 2, "test run count mismatch")

			var testRun1, testRun2 string
			database.MustScan(t, "finding test run 1", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v3"), &testRun1)
			database.MustScan(t, "finding test run 2", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v4"), &testRun2)
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
			schemaKey: consts.SystemSchema,
		},
		{
			mode:      RemoteMode,
			schemaKey: consts.SystemRemoteSchema,
		},
	}

	for _, tc := range testCases {
		func() {
			// set up
			opts := database.TestDBOptions{SkipMigration: true}
			ctx := context.InitTestCtx(t, "../tmp", &opts)
			defer context.TeardownTestCtx(t, ctx)

			db := ctx.DB

			database.MustExec(t, "creating a temporary table for testing", db,
				"CREATE TABLE migrate_run_test ( name string )")

			sequence := []migration{
				{
					name: "v1",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v1 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v1")
						return nil
					},
				},
				{
					name: "v2",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v2 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v2")
						return nil
					},
				},
				{
					name: "v3",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v3 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v3")
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
			database.MustScan(t, "getting schema", db.QueryRow("SELECT value FROM system WHERE key = ?", tc.schemaKey), &schema)
			assert.Equal(t, schema, 3, "schema was not updated")

			var testRunCount int
			database.MustScan(t, "counting test runs", db.QueryRow("SELECT count(*) FROM migrate_run_test"), &testRunCount)
			assert.Equal(t, testRunCount, 3, "test run count mismatch")

			var testRun1, testRun2, testRun3 string
			database.MustScan(t, "finding test run 1", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v1"), &testRun1)
			database.MustScan(t, "finding test run 2", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v2"), &testRun2)
			database.MustScan(t, "finding test run 2", db.QueryRow("SELECT name FROM migrate_run_test WHERE name = ?", "v3"), &testRun3)
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
			schemaKey: consts.SystemSchema,
		},
		{
			mode:      RemoteMode,
			schemaKey: consts.SystemRemoteSchema,
		},
	}

	for _, tc := range testCases {
		func() {
			// set up
			opts := database.TestDBOptions{SkipMigration: true}
			ctx := context.InitTestCtx(t, "../tmp", &opts)
			defer context.TeardownTestCtx(t, ctx)

			db := ctx.DB

			database.MustExec(t, "creating a temporary table for testing", db,
				"CREATE TABLE migrate_run_test ( name string )")

			database.MustExec(t, "inserting a schema", db, "INSERT INTO system (key, value) VALUES (?, ?)", tc.schemaKey, 3)

			sequence := []migration{
				{
					name: "v1",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v1 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v1")
						return nil
					},
				},
				{
					name: "v2",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v2 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v2")
						return nil
					},
				},
				{
					name: "v3",
					run: func(ctx context.DnoteCtx, db *database.DB) error {
						database.MustExec(t, "marking v3 completed", db, "INSERT INTO migrate_run_test (name) VALUES (?)", "v3")
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
			database.MustScan(t, "getting schema", db.QueryRow("SELECT value FROM system WHERE key = ?", tc.schemaKey), &schema)
			assert.Equal(t, schema, 3, "schema was not updated")

			var testRunCount int
			database.MustScan(t, "counting test runs", db.QueryRow("SELECT count(*) FROM migrate_run_test"), &testRunCount)
			assert.Equal(t, testRunCount, 0, "test run count mismatch")
		}()
	}
}

func TestLocalMigration1(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-1-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB
	data := testutils.MustMarshalJSON(t, actions.AddBookDataV1{BookName: "js"})
	a1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a1UUID, 1, "add_book", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV1{NoteUUID: "note-1-uuid", FromBook: "js", ToBook: "", Content: "note 1"})
	a2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a2UUID, 1, "edit_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV1{NoteUUID: "note-2-uuid", FromBook: "js", ToBook: "", Content: "note 2"})
	a3UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
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
	database.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	assert.Equal(t, actionCount, 3, "action count mismatch")

	var a1, a2, a3 actions.Action
	database.MustScan(t, "getting action 1", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a1UUID),
		&a1.Schema, &a1.Type, &a1.Data, &a1.Timestamp)
	database.MustScan(t, "getting action 2", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a2UUID),
		&a2.Schema, &a2.Type, &a2.Data, &a2.Timestamp)
	database.MustScan(t, "getting action 3", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a3UUID),
		&a3.Schema, &a3.Type, &a3.Data, &a3.Timestamp)

	var a1Data actions.AddBookDataV1
	var a2Data, a3Data actions.EditNoteDataV3
	testutils.MustUnmarshalJSON(t, a1.Data, &a1Data)
	testutils.MustUnmarshalJSON(t, a2.Data, &a2Data)
	testutils.MustUnmarshalJSON(t, a3.Data, &a3Data)

	assert.Equal(t, a1.Schema, 1, "a1 schema mismatch")
	assert.Equal(t, a1.Type, "add_book", "a1 type mismatch")
	assert.Equal(t, a1.Timestamp, int64(1537829463), "a1 timestamp mismatch")
	assert.Equal(t, a1Data.BookName, "js", "a1 data book_name mismatch")

	assert.Equal(t, a2.Schema, 3, "a2 schema mismatch")
	assert.Equal(t, a2.Type, "edit_note", "a2 type mismatch")
	assert.Equal(t, a2.Timestamp, int64(1537829463), "a2 timestamp mismatch")
	assert.Equal(t, a2Data.NoteUUID, "note-1-uuid", "a2 data note_uuid mismatch")
	assert.Equal(t, a2Data.BookName, (*string)(nil), "a2 data book_name mismatch")
	assert.Equal(t, *a2Data.Content, "note 1", "a2 data content mismatch")
	assert.Equal(t, *a2Data.Public, false, "a2 data public mismatch")

	assert.Equal(t, a3.Schema, 3, "a3 schema mismatch")
	assert.Equal(t, a3.Type, "edit_note", "a3 type mismatch")
	assert.Equal(t, a3.Timestamp, int64(1537829463), "a3 timestamp mismatch")
	assert.Equal(t, a3Data.NoteUUID, "note-2-uuid", "a3 data note_uuid mismatch")
	assert.Equal(t, a3Data.BookName, (*string)(nil), "a3 data book_name mismatch")
	assert.Equal(t, *a3Data.Content, "note 2", "a3 data content mismatch")
	assert.Equal(t, *a3Data.Public, false, "a3 data public mismatch")
}

func TestLocalMigration2(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-1-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB
	c1 := "note 1 - v1"
	c2 := "note 1 - v2"
	css := "css"

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting css book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "css")

	data := testutils.MustMarshalJSON(t, actions.AddNoteDataV2{NoteUUID: "note-1-uuid", BookName: "js", Content: "note 1", Public: false})
	a1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a1UUID, 2, "add_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV2{NoteUUID: "note-1-uuid", FromBook: "js", ToBook: nil, Content: &c1, Public: nil})
	a2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a2UUID, 2, "edit_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV2{NoteUUID: "note-1-uuid", FromBook: "js", ToBook: &css, Content: &c2, Public: nil})
	a3UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
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
	database.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	assert.Equal(t, actionCount, 3, "action count mismatch")

	var a1, a2, a3 actions.Action
	database.MustScan(t, "getting action 1", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a1UUID),
		&a1.Schema, &a1.Type, &a1.Data, &a1.Timestamp)
	database.MustScan(t, "getting action 2", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a2UUID),
		&a2.Schema, &a2.Type, &a2.Data, &a2.Timestamp)
	database.MustScan(t, "getting action 3", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a3UUID),
		&a3.Schema, &a3.Type, &a3.Data, &a3.Timestamp)

	var a1Data actions.AddNoteDataV2
	var a2Data, a3Data actions.EditNoteDataV3
	testutils.MustUnmarshalJSON(t, a1.Data, &a1Data)
	testutils.MustUnmarshalJSON(t, a2.Data, &a2Data)
	testutils.MustUnmarshalJSON(t, a3.Data, &a3Data)

	assert.Equal(t, a1.Schema, 2, "a1 schema mismatch")
	assert.Equal(t, a1.Type, "add_note", "a1 type mismatch")
	assert.Equal(t, a1.Timestamp, int64(1537829463), "a1 timestamp mismatch")
	assert.Equal(t, a1Data.NoteUUID, "note-1-uuid", "a1 data note_uuid mismatch")
	assert.Equal(t, a1Data.BookName, "js", "a1 data book_name mismatch")
	assert.Equal(t, a1Data.Public, false, "a1 data public mismatch")

	assert.Equal(t, a2.Schema, 3, "a2 schema mismatch")
	assert.Equal(t, a2.Type, "edit_note", "a2 type mismatch")
	assert.Equal(t, a2.Timestamp, int64(1537829463), "a2 timestamp mismatch")
	assert.Equal(t, a2Data.NoteUUID, "note-1-uuid", "a2 data note_uuid mismatch")
	assert.Equal(t, a2Data.BookName, (*string)(nil), "a2 data book_name mismatch")
	assert.Equal(t, *a2Data.Content, c1, "a2 data content mismatch")
	assert.Equal(t, a2Data.Public, (*bool)(nil), "a2 data public mismatch")

	assert.Equal(t, a3.Schema, 3, "a3 schema mismatch")
	assert.Equal(t, a3.Type, "edit_note", "a3 type mismatch")
	assert.Equal(t, a3.Timestamp, int64(1537829463), "a3 timestamp mismatch")
	assert.Equal(t, a3Data.NoteUUID, "note-1-uuid", "a3 data note_uuid mismatch")
	assert.Equal(t, *a3Data.BookName, "css", "a3 data book_name mismatch")
	assert.Equal(t, *a3Data.Content, c2, "a3 data content mismatch")
	assert.Equal(t, a3Data.Public, (*bool)(nil), "a3 data public mismatch")
}

func TestLocalMigration3(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-1-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB
	data := testutils.MustMarshalJSON(t, actions.AddNoteDataV2{NoteUUID: "note-1-uuid", BookName: "js", Content: "note 1", Public: false})
	a1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a1UUID, 2, "add_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.RemoveNoteDataV1{NoteUUID: "note-1-uuid", BookName: "js"})
	a2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", a2UUID, 1, "remove_note", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.RemoveNoteDataV1{NoteUUID: "note-2-uuid", BookName: "js"})
	a3UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
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
	database.MustScan(t, "counting actions", db.QueryRow("SELECT count(*) FROM actions"), &actionCount)
	assert.Equal(t, actionCount, 3, "action count mismatch")

	var a1, a2, a3 actions.Action
	database.MustScan(t, "getting action 1", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a1UUID),
		&a1.Schema, &a1.Type, &a1.Data, &a1.Timestamp)
	database.MustScan(t, "getting action 2", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a2UUID),
		&a2.Schema, &a2.Type, &a2.Data, &a2.Timestamp)
	database.MustScan(t, "getting action 3", db.QueryRow("SELECT schema, type, data, timestamp FROM actions WHERE uuid = ?", a3UUID),
		&a3.Schema, &a3.Type, &a3.Data, &a3.Timestamp)

	var a1Data actions.AddNoteDataV2
	var a2Data, a3Data actions.RemoveNoteDataV2
	testutils.MustUnmarshalJSON(t, a1.Data, &a1Data)
	testutils.MustUnmarshalJSON(t, a2.Data, &a2Data)
	testutils.MustUnmarshalJSON(t, a3.Data, &a3Data)

	assert.Equal(t, a1.Schema, 2, "a1 schema mismatch")
	assert.Equal(t, a1.Type, "add_note", "a1 type mismatch")
	assert.Equal(t, a1.Timestamp, int64(1537829463), "a1 timestamp mismatch")
	assert.Equal(t, a1Data.NoteUUID, "note-1-uuid", "a1 data note_uuid mismatch")
	assert.Equal(t, a1Data.BookName, "js", "a1 data book_name mismatch")
	assert.Equal(t, a1Data.Content, "note 1", "a1 data content mismatch")
	assert.Equal(t, a1Data.Public, false, "a1 data public mismatch")

	assert.Equal(t, a2.Schema, 2, "a2 schema mismatch")
	assert.Equal(t, a2.Type, "remove_note", "a2 type mismatch")
	assert.Equal(t, a2.Timestamp, int64(1537829463), "a2 timestamp mismatch")
	assert.Equal(t, a2Data.NoteUUID, "note-1-uuid", "a2 data note_uuid mismatch")

	assert.Equal(t, a3.Schema, 2, "a3 schema mismatch")
	assert.Equal(t, a3.Type, "remove_note", "a3 type mismatch")
	assert.Equal(t, a3.Timestamp, int64(1537829463), "a3 timestamp mismatch")
	assert.Equal(t, a3Data.NoteUUID, "note-2-uuid", "a3 data note_uuid mismatch")
}

func TestLocalMigration4(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-1-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting css book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "css")
	n1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting css note", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", n1UUID, b1UUID, "n1 content", time.Now().UnixNano())

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
	database.MustScan(t, "scanning the newly added dirty flag of n1", db.QueryRow("SELECT dirty, deleted, usn FROM notes WHERE uuid = ?", n1UUID), &n1Dirty, &n1Deleted, &n1USN)
	database.MustScan(t, "scanning the newly added dirty flag of b1", db.QueryRow("SELECT dirty, deleted, usn FROM books WHERE uuid = ?", b1UUID), &b1Dirty, &b1Deleted, &b1USN)

	assert.Equal(t, n1Dirty, false, "n1 dirty flag should be false by default")
	assert.Equal(t, b1Dirty, false, "b1 dirty flag should be false by default")

	assert.Equal(t, n1Deleted, false, "n1 deleted flag should be false by default")
	assert.Equal(t, b1Deleted, false, "b1 deleted flag should be false by default")

	assert.Equal(t, n1USN, 0, "n1 usn flag should be 0 by default")
	assert.Equal(t, b1USN, 0, "b1 usn flag should be 0 by default")
}

func TestLocalMigration5(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-5-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting css book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "css")
	b2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting js book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "js")

	n1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting css note", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", n1UUID, b1UUID, "n1 content", time.Now().UnixNano())
	n2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting css note", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", n2UUID, b1UUID, "n2 content", time.Now().UnixNano())
	n3UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting css note", db, "INSERT INTO notes (uuid, book_uuid, content, added_on) VALUES (?, ?, ?, ?)", n3UUID, b1UUID, "n3 content", time.Now().UnixNano())

	data := testutils.MustMarshalJSON(t, actions.AddBookDataV1{BookName: "js"})
	database.MustExec(t, "inserting a1", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", "a1-uuid", 1, "add_book", string(data), 1537829463)

	data = testutils.MustMarshalJSON(t, actions.AddNoteDataV2{NoteUUID: n1UUID, BookName: "css", Content: "n1 content", Public: false})
	database.MustExec(t, "inserting a2", db,
		"INSERT INTO actions (uuid, schema, type, data, timestamp) VALUES (?, ?, ?, ?, ?)", "a2-uuid", 1, "add_note", string(data), 1537829463)

	updatedContent := "updated content"
	data = testutils.MustMarshalJSON(t, actions.EditNoteDataV3{NoteUUID: n2UUID, BookName: (*string)(nil), Content: &updatedContent, Public: (*bool)(nil)})
	database.MustExec(t, "inserting a3", db,
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
	database.MustScan(t, "scanning the newly added dirty flag of b1", db.QueryRow("SELECT dirty FROM books WHERE uuid = ?", b1UUID), &b1Dirty)
	database.MustScan(t, "scanning the newly added dirty flag of b2", db.QueryRow("SELECT dirty FROM books WHERE uuid = ?", b2UUID), &b2Dirty)
	database.MustScan(t, "scanning the newly added dirty flag of n1", db.QueryRow("SELECT dirty FROM notes WHERE uuid = ?", n1UUID), &n1Dirty)
	database.MustScan(t, "scanning the newly added dirty flag of n2", db.QueryRow("SELECT dirty FROM notes WHERE uuid = ?", n2UUID), &n2Dirty)
	database.MustScan(t, "scanning the newly added dirty flag of n3", db.QueryRow("SELECT dirty FROM notes WHERE uuid = ?", n3UUID), &n3Dirty)

	assert.Equal(t, b1Dirty, false, "b1 dirty flag should be false by default")
	assert.Equal(t, b2Dirty, true, "b2 dirty flag should be false by default")
	assert.Equal(t, n1Dirty, true, "n1 dirty flag should be false by default")
	assert.Equal(t, n2Dirty, true, "n2 dirty flag should be false by default")
	assert.Equal(t, n3Dirty, false, "n3 dirty flag should be false by default")
}

func TestLocalMigration6(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-5-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	data := testutils.MustMarshalJSON(t, actions.AddBookDataV1{BookName: "js"})
	a1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting action", db,
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
	assert.Equal(t, count, 0, "actions table should have been deleted")
}

func TestLocalMigration7_trash(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-7-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting trash book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "trash")

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
	database.MustScan(t, "scanning b1 label", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b1UUID), &b1Label, &b1Dirty)
	assert.Equal(t, b1Label, "trash (2)", "b1 label was not migrated")
	assert.Equal(t, b1Dirty, true, "b1 was not marked dirty")
}

func TestLocalMigration7_conflicts(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-7-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "conflicts")

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
	database.MustScan(t, "scanning b1 label", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b1UUID), &b1Label, &b1Dirty)
	assert.Equal(t, b1Label, "conflicts (2)", "b1 label was not migrated")
	assert.Equal(t, b1Dirty, true, "b1 was not marked dirty")
}

func TestLocalMigration7_conflicts_dup(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-7-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "conflicts")
	b2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "conflicts (2)")

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
	database.MustScan(t, "scanning b1 label", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b1UUID), &b1Label, &b1Dirty)
	database.MustScan(t, "scanning b2 label", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b2UUID), &b2Label, &b2Dirty)
	assert.Equal(t, b1Label, "conflicts (3)", "b1 label was not migrated")
	assert.Equal(t, b2Label, "conflicts (2)", "b1 label was not migrated")
	assert.Equal(t, b1Dirty, true, "b1 was not marked dirty")
	assert.Equal(t, b2Dirty, false, "b2 should not have been marked dirty")
}

func TestLocalMigration8(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-8-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1")

	n1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting n1", db, `INSERT INTO notes
		(id, uuid, book_uuid, content, added_on, edited_on, public, dirty, usn, deleted) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 1, n1UUID, b1UUID, "n1 Body", 1, 2, true, true, 20, false)
	n2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting n2", db, `INSERT INTO notes
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
	database.MustScan(t, "scanning n1", db.QueryRow("SELECT book_uuid, body, added_on, edited_on, usn,  public, dirty, deleted FROM notes WHERE uuid = ?", n1UUID), &n1BookUUID, &n1Body, &n1AddedOn, &n1EditedOn, &n1USN, &n1Public, &n1Dirty, &n1Deleted)

	var n2BookUUID, n2Body string
	var n2AddedOn, n2EditedOn int64
	var n2USN int
	var n2Public, n2Dirty, n2Deleted bool
	database.MustScan(t, "scanning n2", db.QueryRow("SELECT book_uuid, body, added_on, edited_on, usn,  public, dirty, deleted FROM notes WHERE uuid = ?", n2UUID), &n2BookUUID, &n2Body, &n2AddedOn, &n2EditedOn, &n2USN, &n2Public, &n2Dirty, &n2Deleted)

	assert.Equal(t, n1BookUUID, b1UUID, "n1 BookUUID mismatch")
	assert.Equal(t, n1Body, "n1 Body", "n1 Body mismatch")
	assert.Equal(t, n1AddedOn, int64(1), "n1 AddedOn mismatch")
	assert.Equal(t, n1EditedOn, int64(2), "n1 EditedOn mismatch")
	assert.Equal(t, n1USN, 20, "n1 USN mismatch")
	assert.Equal(t, n1Public, true, "n1 Public mismatch")
	assert.Equal(t, n1Dirty, true, "n1 Dirty mismatch")
	assert.Equal(t, n1Deleted, false, "n1 Deleted mismatch")

	assert.Equal(t, n2BookUUID, b1UUID, "n2 BookUUID mismatch")
	assert.Equal(t, n2Body, "", "n2 Body mismatch")
	assert.Equal(t, n2AddedOn, int64(3), "n2 AddedOn mismatch")
	assert.Equal(t, n2EditedOn, int64(4), "n2 EditedOn mismatch")
	assert.Equal(t, n2USN, 21, "n2 USN mismatch")
	assert.Equal(t, n2Public, false, "n2 Public mismatch")
	assert.Equal(t, n2Dirty, true, "n2 Dirty mismatch")
	assert.Equal(t, n2Deleted, true, "n2 Deleted mismatch")
}

func TestLocalMigration9(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-9-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "b1")

	n1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting n1", db, `INSERT INTO notes
		(uuid, book_uuid, body, added_on, edited_on, public, dirty, usn, deleted) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?)`, n1UUID, b1UUID, "n1 Body", 1, 2, true, true, 20, false)
	n2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting n2", db, `INSERT INTO notes
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
	database.MustScan(t, "counting note_fts", db.QueryRow("SELECT count(*) FROM note_fts;"), &noteFtsCount)
	assert.Equal(t, noteFtsCount, 2, "noteFtsCount mismatch")

	var resCount int
	database.MustScan(t, "counting result", db.QueryRow("SELECT count(*) FROM note_fts WHERE note_fts MATCH ?", "n1"), &resCount)
	assert.Equal(t, resCount, 1, "noteFtsCount mismatch")
}

func TestLocalMigration10(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-10-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "123")
	b2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 2", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "123 javascript")
	b3UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 3", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b3UUID, "foo")
	b4UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 4", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b4UUID, "+123")
	b5UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 5", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b5UUID, "0123")
	b6UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 6", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b6UUID, "javascript 123")
	b7UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 7", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b7UUID, "123 (1)")
	b8UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 8", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b8UUID, "5")

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm10.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test

	// assert that note_fts was populated with correct values
	var b1Label, b2Label, b3Label, b4Label, b5Label, b6Label, b7Label, b8Label string
	var b1Dirty, b2Dirty, b3Dirty, b4Dirty, b5Dirty, b6Dirty, b7Dirty, b8Dirty bool

	database.MustScan(t, "getting b1", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b1UUID), &b1Label, &b1Dirty)
	database.MustScan(t, "getting b2", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b2UUID), &b2Label, &b2Dirty)
	database.MustScan(t, "getting b3", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b3UUID), &b3Label, &b3Dirty)
	database.MustScan(t, "getting b4", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b4UUID), &b4Label, &b4Dirty)
	database.MustScan(t, "getting b5", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b5UUID), &b5Label, &b5Dirty)
	database.MustScan(t, "getting b6", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b6UUID), &b6Label, &b6Dirty)
	database.MustScan(t, "getting b7", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b7UUID), &b7Label, &b7Dirty)
	database.MustScan(t, "getting b8", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b8UUID), &b8Label, &b8Dirty)

	assert.Equal(t, b1Label, "123 (2)", "b1Label mismatch")
	assert.Equal(t, b1Dirty, true, "b1Dirty mismatch")
	assert.Equal(t, b2Label, "123 javascript", "b2Label mismatch")
	assert.Equal(t, b2Dirty, false, "b2Dirty mismatch")
	assert.Equal(t, b3Label, "foo", "b3Label mismatch")
	assert.Equal(t, b3Dirty, false, "b3Dirty mismatch")
	assert.Equal(t, b4Label, "+123", "b4Label mismatch")
	assert.Equal(t, b4Dirty, false, "b4Dirty mismatch")
	assert.Equal(t, b5Label, "0123 (1)", "b5Label mismatch")
	assert.Equal(t, b5Dirty, true, "b5Dirty mismatch")
	assert.Equal(t, b6Label, "javascript 123", "b6Label mismatch")
	assert.Equal(t, b6Dirty, false, "b6Dirty mismatch")
	assert.Equal(t, b7Label, "123 (1)", "b7Label mismatch")
	assert.Equal(t, b7Dirty, false, "b7Dirty mismatch")
	assert.Equal(t, b8Label, "5 (1)", "b8Label mismatch")
	assert.Equal(t, b8Dirty, true, "b8Dirty mismatch")
}

func TestLocalMigration11(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-11-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	db := ctx.DB

	b1UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "foo")
	b2UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 2", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "bar baz")
	b3UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 3", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b3UUID, "quz qux")
	b4UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 4", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b4UUID, "quz_qux")
	b5UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 5", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b5UUID, "foo bar baz quz 123")
	b6UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 6", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b6UUID, "foo_bar baz")
	b7UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 7", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b7UUID, "cool ideas")
	b8UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 8", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b8UUID, "cool_ideas")
	b9UUID := utils.GenerateUUID()
	database.MustExec(t, "inserting book 9", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b9UUID, "cool_ideas_2")

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	err = lm11.run(ctx, tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	tx.Commit()

	// Test
	var bookCount int
	database.MustScan(t, "counting books", db.QueryRow("SELECT count(*) FROM books"), &bookCount)
	assert.Equal(t, bookCount, 9, "bookCount mismatch")

	// assert that note_fts was populated with correct values
	var b1Label, b2Label, b3Label, b4Label, b5Label, b6Label, b7Label, b8Label, b9Label string
	var b1Dirty, b2Dirty, b3Dirty, b4Dirty, b5Dirty, b6Dirty, b7Dirty, b8Dirty, b9Dirty bool

	database.MustScan(t, "getting b1", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b1UUID), &b1Label, &b1Dirty)
	database.MustScan(t, "getting b2", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b2UUID), &b2Label, &b2Dirty)
	database.MustScan(t, "getting b3", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b3UUID), &b3Label, &b3Dirty)
	database.MustScan(t, "getting b4", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b4UUID), &b4Label, &b4Dirty)
	database.MustScan(t, "getting b5", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b5UUID), &b5Label, &b5Dirty)
	database.MustScan(t, "getting b6", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b6UUID), &b6Label, &b6Dirty)
	database.MustScan(t, "getting b7", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b7UUID), &b7Label, &b7Dirty)
	database.MustScan(t, "getting b8", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b8UUID), &b8Label, &b8Dirty)
	database.MustScan(t, "getting b9", db.QueryRow("SELECT label, dirty FROM books WHERE uuid = ?", b9UUID), &b9Label, &b9Dirty)

	assert.Equal(t, b1Label, "foo", "b1Label mismatch")
	assert.Equal(t, b1Dirty, false, "b1Dirty mismatch")
	assert.Equal(t, b2Label, "bar_baz", "b2Label mismatch")
	assert.Equal(t, b2Dirty, true, "b2Dirty mismatch")
	assert.Equal(t, b3Label, "quz_qux_2", "b3Label mismatch")
	assert.Equal(t, b3Dirty, true, "b3Dirty mismatch")
	assert.Equal(t, b4Label, "quz_qux", "b4Label mismatch")
	assert.Equal(t, b4Dirty, false, "b4Dirty mismatch")
	assert.Equal(t, b5Label, "foo_bar_baz_quz_123", "b5Label mismatch")
	assert.Equal(t, b5Dirty, true, "b5Dirty mismatch")
	assert.Equal(t, b6Label, "foo_bar_baz", "b6Label mismatch")
	assert.Equal(t, b6Dirty, true, "b6Dirty mismatch")
	assert.Equal(t, b7Label, "cool_ideas_3", "b7Label mismatch")
	assert.Equal(t, b7Dirty, true, "b7Dirty mismatch")
	assert.Equal(t, b8Label, "cool_ideas", "b8Label mismatch")
	assert.Equal(t, b8Dirty, false, "b8Dirty mismatch")
	assert.Equal(t, b9Label, "cool_ideas_2", "b9Label mismatch")
	assert.Equal(t, b9Dirty, false, "b9Dirty mismatch")
}

func TestLocalMigration12(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/local-12-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)

	data := []byte("editor: vim")
	path := fmt.Sprintf("%s/dnoterc", ctx.DnoteDir)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to write schema file"))
	}

	// execute
	err := lm12.run(ctx, nil)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to run"))
	}

	// test
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(errors.Wrap(err, "reading config"))
	}

	type config struct {
		APIEndpoint string `yaml:"apiEndpoint"`
	}

	var cf config
	err = yaml.Unmarshal(b, &cf)
	if err != nil {
		t.Fatal(errors.Wrap(err, "unmarshalling config"))
	}

	assert.NotEqual(t, cf.APIEndpoint, "", "apiEndpoint was not populated")
}

func TestRemoteMigration1(t *testing.T) {
	// set up
	opts := database.TestDBOptions{SchemaSQLPath: "./fixtures/remote-1-pre-schema.sql", SkipMigration: true}
	ctx := context.InitTestCtx(t, "../tmp", &opts)
	defer context.TeardownTestCtx(t, ctx)
	testutils.Login(t, &ctx)

	JSBookUUID := "existing-js-book-uuid"
	CSSBookUUID := "existing-css-book-uuid"
	linuxBookUUID := "existing-linux-book-uuid"
	newJSBookUUID := "new-js-book-uuid"
	newCSSBookUUID := "new-css-book-uuid"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/v3/books" {
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
	database.MustExec(t, "inserting js book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", JSBookUUID, "js")
	database.MustExec(t, "inserting css book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", CSSBookUUID, "css")
	database.MustExec(t, "inserting linux book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", linuxBookUUID, "linux")
	database.MustExec(t, "inserting sessionKey", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKey, "someSessionKey")
	database.MustExec(t, "inserting sessionKeyExpiry", db, "INSERT INTO system (key, value) VALUES (?, ?)", consts.SystemSessionKeyExpiry, time.Now().Add(24*time.Hour).Unix())

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
	database.MustScan(t, "getting js book uuid", db.QueryRow("SELECT uuid FROM books WHERE label = ?", "js"), &postJSBookUUID)
	database.MustScan(t, "getting css book uuid", db.QueryRow("SELECT uuid FROM books WHERE label = ?", "css"), &postCSSBookUUID)
	database.MustScan(t, "getting linux book uuid", db.QueryRow("SELECT uuid FROM books WHERE label = ?", "linux"), &postLinuxBookUUID)

	assert.Equal(t, postJSBookUUID, newJSBookUUID, "js book uuid was not updated correctly")
	assert.Equal(t, postCSSBookUUID, newCSSBookUUID, "css book uuid was not updated correctly")
	assert.Equal(t, postLinuxBookUUID, linuxBookUUID, "linux book uuid changed")
}
