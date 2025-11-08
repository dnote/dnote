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
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/pkg/errors"
)

// bookInfo is an information about the book to be printed on screen
type bookInfo struct {
	BookLabel string
	NoteCount int
}

// noteInfo is an information about the note to be printed on screen
type noteInfo struct {
	RowID int
	Body  string
}

// getNewlineIdx returns the index of newline character in a string
func getNewlineIdx(str string) int {
	// Check for \r\n first
	if idx := strings.Index(str, "\r\n"); idx != -1 {
		return idx
	}

	// Then check for \n
	return strings.Index(str, "\n")
}

// formatBody returns an excerpt of the given raw note content and a boolean
// indicating if the returned string has been excertped
func formatBody(noteBody string) (string, bool) {
	trimmed := strings.TrimRight(noteBody, "\r\n")
	newlineIdx := getNewlineIdx(trimmed)

	if newlineIdx > -1 {
		ret := strings.Trim(trimmed[0:newlineIdx], " ")

		return ret, true
	}

	return strings.Trim(trimmed, " "), false
}

func printBookLine(w io.Writer, info bookInfo, nameOnly bool) {
	if nameOnly {
		fmt.Fprintln(w, info.BookLabel)
	} else {
		fmt.Fprintf(w, "%s %s\n", info.BookLabel, log.ColorYellow.Sprintf("(%d)", info.NoteCount))
	}
}

func listBooks(ctx context.DnoteCtx, w io.Writer, nameOnly bool) error {
	db := ctx.DB

	rows, err := db.Query(`SELECT books.label, count(notes.uuid) note_count
	FROM books
	LEFT JOIN notes ON notes.book_uuid = books.uuid AND notes.deleted = false
	WHERE books.deleted = false
	GROUP BY books.uuid
	ORDER BY books.label ASC;`)
	if err != nil {
		return errors.Wrap(err, "querying books")
	}
	defer rows.Close()

	infos := []bookInfo{}
	for rows.Next() {
		var info bookInfo
		err = rows.Scan(&info.BookLabel, &info.NoteCount)
		if err != nil {
			return errors.Wrap(err, "scanning a row")
		}

		infos = append(infos, info)
	}

	for _, info := range infos {
		printBookLine(w, info, nameOnly)
	}

	return nil
}

func listNotes(ctx context.DnoteCtx, w io.Writer, bookName string) error {
	db := ctx.DB

	var bookUUID string
	err := db.QueryRow("SELECT uuid FROM books WHERE label = ?", bookName).Scan(&bookUUID)
	if err == sql.ErrNoRows {
		return errors.New("book not found")
	} else if err != nil {
		return errors.Wrap(err, "querying the book")
	}

	rows, err := db.Query(`SELECT rowid, body FROM notes WHERE book_uuid = ? AND deleted = ? ORDER BY added_on ASC;`, bookUUID, false)
	if err != nil {
		return errors.Wrap(err, "querying notes")
	}
	defer rows.Close()

	infos := []noteInfo{}
	for rows.Next() {
		var info noteInfo
		err = rows.Scan(&info.RowID, &info.Body)
		if err != nil {
			return errors.Wrap(err, "scanning a row")
		}

		infos = append(infos, info)
	}

	fmt.Fprintf(w, "on book %s\n", bookName)

	for _, info := range infos {
		body, isExcerpt := formatBody(info.Body)

		rowid := log.ColorYellow.Sprintf("(%d)", info.RowID)
		if isExcerpt {
			body = fmt.Sprintf("%s %s", body, log.ColorYellow.Sprintf("[---More---]"))
		}

		fmt.Fprintf(w, "%s %s\n", rowid, body)
	}

	return nil
}
