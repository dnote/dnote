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
	"fmt"
	"strings"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
)

func TestGetNewlineIdx(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{
			input:    "hello\nworld",
			expected: 5,
		},
		{
			input:    "hello\r\nworld",
			expected: 5,
		},
		{
			input:    "no newline here",
			expected: -1,
		},
		{
			input:    "",
			expected: -1,
		},
		{
			input:    "\n",
			expected: 0,
		},
		{
			input:    "\r\n",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input: %q", tc.input), func(t *testing.T) {
			got := getNewlineIdx(tc.input)
			assert.Equal(t, got, tc.expected, "newline index mismatch")
		})
	}
}

func TestFormatBody(t *testing.T) {
	testCases := []struct {
		input           string
		expectedBody    string
		expectedExcerpt bool
	}{
		{
			input:           "single line",
			expectedBody:    "single line",
			expectedExcerpt: false,
		},
		{
			input:           "first line\nsecond line",
			expectedBody:    "first line",
			expectedExcerpt: true,
		},
		{
			input:           "first line\r\nsecond line",
			expectedBody:    "first line",
			expectedExcerpt: true,
		},
		{
			input:           "  spaced line  ",
			expectedBody:    "spaced line",
			expectedExcerpt: false,
		},
		{
			input:           "  first line  \nsecond line",
			expectedBody:    "first line",
			expectedExcerpt: true,
		},
		{
			input:           "",
			expectedBody:    "",
			expectedExcerpt: false,
		},
		{
			input:           "line with trailing newline\n",
			expectedBody:    "line with trailing newline",
			expectedExcerpt: false,
		},
		{
			input:           "line with trailing newlines\n\n",
			expectedBody:    "line with trailing newlines",
			expectedExcerpt: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input: %q", tc.input), func(t *testing.T) {
			gotBody, gotExcerpt := formatBody(tc.input)
			assert.Equal(t, gotBody, tc.expectedBody, "formatted body mismatch")
			assert.Equal(t, gotExcerpt, tc.expectedExcerpt, "excerpt flag mismatch")
		})
	}
}

func TestListNotes(t *testing.T) {
	// Setup
	db := database.InitTestMemoryDB(t)
	defer db.Close()

	bookUUID := "js-book-uuid"
	database.MustExec(t, "inserting book", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", bookUUID, "javascript")
	database.MustExec(t, "inserting note 1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on) VALUES (?, ?, ?, ?)", "note-1", bookUUID, "first note", 1515199943)
	database.MustExec(t, "inserting note 2", db, "INSERT INTO notes (uuid, book_uuid, body, added_on) VALUES (?, ?, ?, ?)", "note-2", bookUUID, "multiline note\nwith second line", 1515199945)

	ctx := context.DnoteCtx{DB: db}
	var buf bytes.Buffer

	// Execute
	err := listNotes(ctx, &buf, "javascript")
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	// Verify output
	assert.Equal(t, strings.Contains(got, "on book javascript"), true, "should show book name")
	assert.Equal(t, strings.Contains(got, "first note"), true, "should contain first note")
	assert.Equal(t, strings.Contains(got, "multiline note"), true, "should show first line of multiline note")
	assert.Equal(t, strings.Contains(got, "[---More---]"), true, "should show more indicator for multiline note")
	assert.Equal(t, strings.Contains(got, "with second line"), false, "should not show second line of multiline note")
}

func TestListBooks(t *testing.T) {
	// Setup
	db := database.InitTestMemoryDB(t)
	defer db.Close()

	b1UUID := "js-book-uuid"
	b2UUID := "linux-book-uuid"

	database.MustExec(t, "inserting book 1", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b1UUID, "javascript")
	database.MustExec(t, "inserting book 2", db, "INSERT INTO books (uuid, label) VALUES (?, ?)", b2UUID, "linux")

	// Add notes to test count
	database.MustExec(t, "inserting note 1", db, "INSERT INTO notes (uuid, book_uuid, body, added_on) VALUES (?, ?, ?, ?)", "note-1", b1UUID, "note body 1", 1515199943)
	database.MustExec(t, "inserting note 2", db, "INSERT INTO notes (uuid, book_uuid, body, added_on) VALUES (?, ?, ?, ?)", "note-2", b1UUID, "note body 2", 1515199944)

	ctx := context.DnoteCtx{DB: db}
	var buf bytes.Buffer

	// Execute
	err := listBooks(ctx, &buf, false)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	// Verify output
	assert.Equal(t, strings.Contains(got, "javascript"), true, "should contain javascript book")
	assert.Equal(t, strings.Contains(got, "linux"), true, "should contain linux book")
	assert.Equal(t, strings.Contains(got, "(2)"), true, "should show 2 notes for javascript")
}
