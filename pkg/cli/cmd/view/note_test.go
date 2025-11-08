/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package view

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
)

func TestViewNote(t *testing.T) {
	db := database.InitTestMemoryDB(t)
	defer db.Close()

	bookUUID := "test-book-uuid"
	database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", bookUUID, "golang")
	database.MustExec(t, "inserting note", db, "INSERT INTO notes (uuid, book_uuid, body, added_on) VALUES (?, ?, ?, ?)",
		"note-uuid", bookUUID, "test note content", 1515199943000000000)

	ctx := context.DnoteCtx{DB: db}
	var buf bytes.Buffer

	err := viewNote(ctx, &buf, "1", false)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	assert.Equal(t, strings.Contains(got, "test note content"), true, "should contain note content")
}

func TestViewNoteContentOnly(t *testing.T) {
	db := database.InitTestMemoryDB(t)
	defer db.Close()

	bookUUID := "test-book-uuid"
	database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", bookUUID, "golang")
	database.MustExec(t, "inserting note", db, "INSERT INTO notes (uuid, book_uuid, body, added_on) VALUES (?, ?, ?, ?)",
		"note-uuid", bookUUID, "test note content", 1515199943000000000)

	ctx := context.DnoteCtx{DB: db}
	var buf bytes.Buffer

	err := viewNote(ctx, &buf, "1", true)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	assert.Equal(t, got, "test note content", "should contain only note content")
}

func TestViewNoteInvalidRowID(t *testing.T) {
	db := database.InitTestMemoryDB(t)
	defer db.Close()

	ctx := context.DnoteCtx{DB: db}
	var buf bytes.Buffer

	err := viewNote(ctx, &buf, "not-a-number", false)
	assert.NotEqual(t, err, nil, "should return error for invalid rowid")
}

func TestViewNoteNotFound(t *testing.T) {
	db := database.InitTestMemoryDB(t)
	defer db.Close()

	ctx := context.DnoteCtx{DB: db}
	var buf bytes.Buffer

	err := viewNote(ctx, &buf, "999", false)
	assert.NotEqual(t, err, nil, "should return error for non-existent note")
}
