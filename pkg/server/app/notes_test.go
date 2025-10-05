/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package app

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestCreateNote(t *testing.T) {
	serverTime := time.Date(2017, time.March, 14, 21, 15, 0, 0, time.UTC)
	mockClock := clock.NewMock()
	mockClock.SetNow(serverTime)

	ts1 := time.Date(2018, time.November, 12, 10, 11, 0, 0, time.UTC).UnixNano()
	ts2 := time.Date(2018, time.November, 15, 0, 1, 10, 0, time.UTC).UnixNano()

	testCases := []struct {
		userUSN          int
		addedOn          *int64
		editedOn         *int64
		expectedUSN      int
		expectedAddedOn  int64
		expectedEditedOn int64
	}{
		{
			userUSN:          8,
			addedOn:          nil,
			editedOn:         nil,
			expectedUSN:      9,
			expectedAddedOn:  serverTime.UnixNano(),
			expectedEditedOn: 0,
		},
		{
			userUSN:          102229,
			addedOn:          &ts1,
			editedOn:         nil,
			expectedUSN:      102230,
			expectedAddedOn:  ts1,
			expectedEditedOn: 0,
		},
		{
			userUSN:          8099,
			addedOn:          &ts1,
			editedOn:         &ts2,
			expectedUSN:      8100,
			expectedAddedOn:  ts1,
			expectedEditedOn: ts2,
		},
	}

	for idx, tc := range testCases {
		func() {
			db := testutils.InitMemoryDB(t)

			user := testutils.SetupUserData(db)
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))
			fmt.Println(user)

			anotherUser := testutils.SetupUserData(db)
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			b1 := database.Book{UserID: user.ID, Label: "js", Deleted: false}
			testutils.MustExec(t, db.Save(&b1), fmt.Sprintf("preparing b1 for test case %d", idx))

			a := NewTest()
			a.DB = db
			a.Clock = mockClock

			if _, err := a.CreateNote(user, b1.UUID, "note content", tc.addedOn, tc.editedOn, false, ""); err != nil {
				t.Fatal(errors.Wrapf(err, "creating note for test case %d", idx))
			}

			var bookCount, noteCount int64
			var noteRecord database.Note
			var userRecord database.User

			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), fmt.Sprintf("counting book for test case %d", idx))
			testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), fmt.Sprintf("counting notes for test case %d", idx))
			testutils.MustExec(t, db.First(&noteRecord), fmt.Sprintf("finding note for test case %d", idx))
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), fmt.Sprintf("finding user for test case %d", idx))

			assert.Equal(t, bookCount, int64(1), "book count mismatch")
			assert.Equal(t, noteCount, int64(1), "note count mismatch")
			assert.NotEqual(t, noteRecord.UUID, "", "note UUID should have been generated")
			assert.Equal(t, noteRecord.UserID, user.ID, "note UserID mismatch")
			assert.Equal(t, noteRecord.Body, "note content", "note Body mismatch")
			assert.Equal(t, noteRecord.Deleted, false, "note Deleted mismatch")
			assert.Equal(t, noteRecord.USN, tc.expectedUSN, "note Label mismatch")
			assert.Equal(t, noteRecord.AddedOn, tc.expectedAddedOn, "note AddedOn mismatch")
			assert.Equal(t, noteRecord.EditedOn, tc.expectedEditedOn, "note EditedOn mismatch")

			assert.Equal(t, userRecord.MaxUSN, tc.expectedUSN, "user max_usn mismatch")

			// Assert FTS table is updated
			var ftsBody string
			testutils.MustExec(t, db.Raw("SELECT body FROM notes_fts WHERE rowid = ?", noteRecord.ID).Scan(&ftsBody), fmt.Sprintf("querying notes_fts for test case %d", idx))
			assert.Equal(t, ftsBody, "note content", "FTS body mismatch")
			var searchCount int64
			testutils.MustExec(t, db.Raw("SELECT COUNT(*) FROM notes_fts WHERE notes_fts MATCH ?", "content").Scan(&searchCount), "searching notes_fts")
			assert.Equal(t, searchCount, int64(1), "Note should still be searchable")
		}()
	}
}

func TestCreateNote_EmptyBody(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db)
	b1 := database.Book{UserID: user.ID, Label: "testBook"}
	testutils.MustExec(t, db.Save(&b1), "preparing book")

	a := NewTest()
	a.DB = db
	a.Clock = clock.NewMock()

	// Create note with empty body
	note, err := a.CreateNote(user, b1.UUID, "", nil, nil, false, "")
	if err != nil {
		t.Fatal(errors.Wrap(err, "creating note with empty body"))
	}

	// Assert FTS entry exists with empty body
	var ftsBody string
	testutils.MustExec(t, db.Raw("SELECT body FROM notes_fts WHERE rowid = ?", note.ID).Scan(&ftsBody), "querying notes_fts for empty body note")
	assert.Equal(t, ftsBody, "", "FTS body should be empty for note created with empty body")
}

func TestUpdateNote(t *testing.T) {
	testCases := []struct {
		userUSN int
	}{
		{
			userUSN: 8,
		},
		{
			userUSN: 102229,
		},
		{
			userUSN: 8099,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			db := testutils.InitMemoryDB(t)

			user := testutils.SetupUserData(db)
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.userUSN), "preparing user max_usn for test case")

			anotherUser := testutils.SetupUserData(db)
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 55), "preparing user max_usn for test case")

			b1 := database.Book{UserID: user.ID, Label: "js", Deleted: false}
			testutils.MustExec(t, db.Save(&b1), "preparing b1 for test case")

			note := database.Note{UserID: user.ID, Deleted: false, Body: "test content", BookUUID: b1.UUID}
			testutils.MustExec(t, db.Save(&note), "preparing note for test case")

			// Assert FTS table has original content
			var ftsBodyBefore string
			testutils.MustExec(t, db.Raw("SELECT body FROM notes_fts WHERE rowid = ?", note.ID).Scan(&ftsBodyBefore), "querying notes_fts before update")
			assert.Equal(t, ftsBodyBefore, "test content", "FTS body mismatch before update")

			c := clock.NewMock()
			content := "updated test content"
			public := true

			a := NewTest()
			a.DB = db
			a.Clock = c

			tx := db.Begin()
			if _, err := a.UpdateNote(tx, user, note, &UpdateNoteParams{
				Content: &content,
				Public:  &public,
			}); err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "updating note"))
			}
			tx.Commit()

			var bookCount, noteCount int64
			var noteRecord database.Note
			var userRecord database.User

			testutils.MustExec(t, db.Model(&database.Book{}).Count(&bookCount), "counting book for test case")
			testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), "counting notes for test case")
			testutils.MustExec(t, db.First(&noteRecord), "finding note for test case")
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), "finding user for test case")

			expectedUSN := tc.userUSN + 1
			assert.Equal(t, bookCount, int64(1), "book count mismatch")
			assert.Equal(t, noteCount, int64(1), "note count mismatch")
			assert.Equal(t, noteRecord.UserID, user.ID, "note UserID mismatch")
			assert.Equal(t, noteRecord.Body, content, "note Body mismatch")
			assert.Equal(t, noteRecord.Public, public, "note Public mismatch")
			assert.Equal(t, noteRecord.Deleted, false, "note Deleted mismatch")
			assert.Equal(t, noteRecord.USN, expectedUSN, "note USN mismatch")
			assert.Equal(t, userRecord.MaxUSN, expectedUSN, "user MaxUSN mismatch")

			// Assert FTS table is updated with new content
			var ftsBodyAfter string
			testutils.MustExec(t, db.Raw("SELECT body FROM notes_fts WHERE rowid = ?", noteRecord.ID).Scan(&ftsBodyAfter), "querying notes_fts after update")
			assert.Equal(t, ftsBodyAfter, content, "FTS body mismatch after update")
			var searchCount int64
			testutils.MustExec(t, db.Raw("SELECT COUNT(*) FROM notes_fts WHERE notes_fts MATCH ?", "updated").Scan(&searchCount), "searching notes_fts")
			assert.Equal(t, searchCount, int64(1), "Note should still be searchable")
		})
	}
}

func TestUpdateNote_SameContent(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db)
	b1 := database.Book{UserID: user.ID, Label: "testBook"}
	testutils.MustExec(t, db.Save(&b1), "preparing book")

	note := database.Note{UserID: user.ID, Deleted: false, Body: "test content", BookUUID: b1.UUID}
	testutils.MustExec(t, db.Save(&note), "preparing note")

	a := NewTest()
	a.DB = db
	a.Clock = clock.NewMock()

	// Update note with same content
	sameContent := "test content"
	tx := db.Begin()
	_, err := a.UpdateNote(tx, user, note, &UpdateNoteParams{
		Content: &sameContent,
	})
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "updating note with same content"))
	}
	tx.Commit()

	// Assert FTS still has the same content
	var ftsBody string
	testutils.MustExec(t, db.Raw("SELECT body FROM notes_fts WHERE rowid = ?", note.ID).Scan(&ftsBody), "querying notes_fts after update")
	assert.Equal(t, ftsBody, "test content", "FTS body should still be 'test content'")

	// Assert it's still searchable
	var searchCount int64
	testutils.MustExec(t, db.Raw("SELECT COUNT(*) FROM notes_fts WHERE notes_fts MATCH ?", "test").Scan(&searchCount), "searching notes_fts")
	assert.Equal(t, searchCount, int64(1), "Note should still be searchable")
}

func TestDeleteNote(t *testing.T) {
	testCases := []struct {
		userUSN     int
		expectedUSN int
	}{
		{
			userUSN:     3,
			expectedUSN: 4,
		},
		{
			userUSN:     9787,
			expectedUSN: 9788,
		},
		{
			userUSN:     787,
			expectedUSN: 788,
		},
	}

	for idx, tc := range testCases {
		func() {
			db := testutils.InitMemoryDB(t)

			user := testutils.SetupUserData(db)
			testutils.MustExec(t, db.Model(&user).Update("max_usn", tc.userUSN), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			anotherUser := testutils.SetupUserData(db)
			testutils.MustExec(t, db.Model(&anotherUser).Update("max_usn", 55), fmt.Sprintf("preparing user max_usn for test case %d", idx))

			b1 := database.Book{UserID: user.ID, Label: "testBook"}
			testutils.MustExec(t, db.Save(&b1), fmt.Sprintf("preparing b1 for test case %d", idx))

			note := database.Note{UserID: user.ID, Deleted: false, Body: "test content", BookUUID: b1.UUID}
			testutils.MustExec(t, db.Save(&note), fmt.Sprintf("preparing note for test case %d", idx))

			// Assert FTS table has content before delete
			var ftsCountBefore int64
			testutils.MustExec(t, db.Raw("SELECT COUNT(*) FROM notes_fts WHERE rowid = ?", note.ID).Scan(&ftsCountBefore), fmt.Sprintf("counting notes_fts before delete for test case %d", idx))
			assert.Equal(t, ftsCountBefore, int64(1), "FTS should have entry before delete")

			a := NewTest()
			a.DB = db

			tx := db.Begin()
			ret, err := a.DeleteNote(tx, user, note)
			if err != nil {
				tx.Rollback()
				t.Fatal(errors.Wrap(err, "deleting note"))
			}
			tx.Commit()

			var noteCount int64
			var noteRecord database.Note
			var userRecord database.User

			testutils.MustExec(t, db.Model(&database.Note{}).Count(&noteCount), fmt.Sprintf("counting notes for test case %d", idx))
			testutils.MustExec(t, db.First(&noteRecord), fmt.Sprintf("finding note for test case %d", idx))
			testutils.MustExec(t, db.Where("id = ?", user.ID).First(&userRecord), fmt.Sprintf("finding user for test case %d", idx))

			assert.Equal(t, noteCount, int64(1), "note count mismatch")

			assert.Equal(t, noteRecord.UserID, user.ID, "note user_id mismatch")
			assert.Equal(t, noteRecord.Body, "", "note content mismatch")
			assert.Equal(t, noteRecord.Deleted, true, "note deleted flag mismatch")
			assert.Equal(t, noteRecord.USN, tc.expectedUSN, "note label mismatch")
			assert.Equal(t, userRecord.MaxUSN, tc.expectedUSN, "user max_usn mismatch")

			assert.Equal(t, ret.UserID, user.ID, "note user_id mismatch")
			assert.Equal(t, ret.Body, "", "note content mismatch")
			assert.Equal(t, ret.Deleted, true, "note deleted flag mismatch")
			assert.Equal(t, ret.USN, tc.expectedUSN, "note label mismatch")

			// Assert FTS body is empty after delete (row still exists but content is cleared)
			var ftsBody string
			testutils.MustExec(t, db.Raw("SELECT body FROM notes_fts WHERE rowid = ?", noteRecord.ID).Scan(&ftsBody), fmt.Sprintf("querying notes_fts after delete for test case %d", idx))
			assert.Equal(t, ftsBody, "", "FTS body should be empty after delete")
		}()
	}
}

func TestGetNotes_FTSSearch(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db)
	b1 := database.Book{UserID: user.ID, Label: "testBook"}
	testutils.MustExec(t, db.Save(&b1), "preparing book")

	// Create notes with different content
	note1 := database.Note{UserID: user.ID, Deleted: false, Body: "foo bar baz bar", BookUUID: b1.UUID}
	testutils.MustExec(t, db.Save(&note1), "preparing note1")

	note2 := database.Note{UserID: user.ID, Deleted: false, Body: "hello run foo", BookUUID: b1.UUID}
	testutils.MustExec(t, db.Save(&note2), "preparing note2")

	note3 := database.Note{UserID: user.ID, Deleted: false, Body: "running quz succeeded", BookUUID: b1.UUID}
	testutils.MustExec(t, db.Save(&note3), "preparing note3")

	a := NewTest()
	a.DB = db
	a.Clock = clock.NewMock()

	// Search "baz"
	result, err := a.GetNotes(user.ID, GetNotesParams{
		Search:    "baz",
		Encrypted: false,
		Page:      1,
		PerPage:   30,
	})
	if err != nil {
		t.Fatal(errors.Wrap(err, "getting notes with FTS search"))
	}
	assert.Equal(t, result.Total, int64(1), "Should find 1 note with 'baz'")
	assert.Equal(t, len(result.Notes), 1, "Should return 1 note")
	for i, note := range result.Notes {
		assert.Equal(t, strings.Contains(note.Body, "<dnotehl>baz</dnotehl>"), true, fmt.Sprintf("Note %d should contain highlighted dnote", i))
	}

	// Search for "running" - should return 1 note
	result, err = a.GetNotes(user.ID, GetNotesParams{
		Search:    "running",
		Encrypted: false,
		Page:      1,
		PerPage:   30,
	})
	if err != nil {
		t.Fatal(errors.Wrap(err, "getting notes with FTS search for review"))
	}
	assert.Equal(t, result.Total, int64(2), "Should find 2 note with 'running'")
	assert.Equal(t, len(result.Notes), 2, "Should return 2 notes")
	assert.Equal(t, result.Notes[0].Body, "<dnotehl>running</dnotehl> quz succeeded", "Should return the review note with highlighting")
	assert.Equal(t, result.Notes[1].Body, "hello <dnotehl>run</dnotehl> foo", "Should return the review note with highlighting")

	// Search for non-existent term - should return 0 notes
	result, err = a.GetNotes(user.ID, GetNotesParams{
		Search:    "nonexistent",
		Encrypted: false,
		Page:      1,
		PerPage:   30,
	})
	if err != nil {
		t.Fatal(errors.Wrap(err, "getting notes with FTS search for nonexistent"))
	}

	assert.Equal(t, result.Total, int64(0), "Should find 0 notes with 'nonexistent'")
	assert.Equal(t, len(result.Notes), 0, "Should return 0 notes")
}

func TestGetNotes_FTSSearch_Snippet(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db)
	b1 := database.Book{UserID: user.ID, Label: "testBook"}
	testutils.MustExec(t, db.Save(&b1), "preparing book")

	// Create a long note to test snippet truncation with "..."
	// The snippet limit is 50 tokens, so we generate enough words to exceed it
	longBody := strings.Repeat("filler ", 100) + "the important keyword appears here"
	longNote := database.Note{UserID: user.ID, Deleted: false, Body: longBody, BookUUID: b1.UUID}
	testutils.MustExec(t, db.Save(&longNote), "preparing long note")

	a := NewTest()
	a.DB = db
	a.Clock = clock.NewMock()

	// Search for "keyword" in long note - should return snippet with "..."
	result, err := a.GetNotes(user.ID, GetNotesParams{
		Search:    "keyword",
		Encrypted: false,
		Page:      1,
		PerPage:   30,
	})
	if err != nil {
		t.Fatal(errors.Wrap(err, "getting notes with FTS search for keyword"))
	}

	assert.Equal(t, result.Total, int64(1), "Should find 1 note with 'keyword'")
	assert.Equal(t, len(result.Notes), 1, "Should return 1 note")
	// The snippet should contain "..." to indicate truncation and the highlighted keyword
	assert.Equal(t, strings.Contains(result.Notes[0].Body, "..."), true, "Snippet should contain '...' for truncation")
	assert.Equal(t, strings.Contains(result.Notes[0].Body, "<dnotehl>keyword</dnotehl>"), true, "Snippet should contain highlighted keyword")
}

func TestGetNotes_FTSSearch_ShortWord(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db)
	b1 := database.Book{UserID: user.ID, Label: "testBook"}
	testutils.MustExec(t, db.Save(&b1), "preparing book")

	// Create notes with short words
	note1 := database.Note{UserID: user.ID, Deleted: false, Body: "a b c", BookUUID: b1.UUID}
	testutils.MustExec(t, db.Save(&note1), "preparing note1")

	note2 := database.Note{UserID: user.ID, Deleted: false, Body: "d", BookUUID: b1.UUID}
	testutils.MustExec(t, db.Save(&note2), "preparing note2")

	a := NewTest()
	a.DB = db
	a.Clock = clock.NewMock()

	result, err := a.GetNotes(user.ID, GetNotesParams{
		Search:    "a",
		Encrypted: false,
		Page:      1,
		PerPage:   30,
	})
	if err != nil {
		t.Fatal(errors.Wrap(err, "getting notes with FTS search for 'a'"))
	}

	assert.Equal(t, result.Total, int64(1), "Should find 1 note")
	assert.Equal(t, len(result.Notes), 1, "Should return 1 note")
	assert.Equal(t, strings.Contains(result.Notes[0].Body, "<dnotehl>a</dnotehl>"), true, "Should contain highlighted 'a'")
}

func TestGetNotes_All(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	user := testutils.SetupUserData(db)
	b1 := database.Book{UserID: user.ID, Label: "testBook"}
	testutils.MustExec(t, db.Save(&b1), "preparing book")

	note1 := database.Note{UserID: user.ID, Deleted: false, Body: "a b c", BookUUID: b1.UUID}
	testutils.MustExec(t, db.Save(&note1), "preparing note1")

	note2 := database.Note{UserID: user.ID, Deleted: false, Body: "d", BookUUID: b1.UUID}
	testutils.MustExec(t, db.Save(&note2), "preparing note2")

	a := NewTest()
	a.DB = db
	a.Clock = clock.NewMock()

	result, err := a.GetNotes(user.ID, GetNotesParams{
		Search:    "",
		Encrypted: false,
		Page:      1,
		PerPage:   30,
	})
	if err != nil {
		t.Fatal(errors.Wrap(err, "getting notes with FTS search for 'a'"))
	}

	assert.Equal(t, result.Total, int64(2), "Should not find all notes")
	assert.Equal(t, len(result.Notes), 2, "Should not find all notes")

	for _, note := range result.Notes {
		assert.Equal(t, strings.Contains(note.Body, "<dnotehl>"), false, "There should be no keywords")
		assert.Equal(t, strings.Contains(note.Body, "</dnotehl>"), false, "There should be no keywords")
	}
	assert.Equal(t, result.Notes[0].Body, "d", "Full content should be returned")
	assert.Equal(t, result.Notes[1].Body, "a b c", "Full content should be returned")
}
