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

package sync

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dnote/dnote/pkg/cli/client"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/dnote/dnote/pkg/cli/utils/diff"
	"github.com/pkg/errors"
)

const (
	modeNormal = iota
	modeRemote
	modeLocal
)

const (
	conflictLabelLocal  = "<<<<<<< Local\n"
	conflictLabelServer = ">>>>>>> Server\n"
	conflictLabelDivide = "=======\n"
)

func sanitize(s string) string {
	var textBuilder strings.Builder
	textBuilder.WriteString(s)
	if !strings.HasSuffix(s, "\n") {
		textBuilder.WriteString("\n")
	}

	return textBuilder.String()
}

// reportBodyConflict returns a conflict report of the local and the remote version
// of a body
func reportBodyConflict(localBody, remoteBody string) string {
	diffs := diff.Do(localBody, remoteBody)

	var ret strings.Builder
	mode := modeNormal
	maxIdx := len(diffs) - 1

	for idx, d := range diffs {
		if d.Type == diff.DiffEqual {
			if mode != modeNormal {
				mode = modeNormal
				ret.WriteString(conflictLabelServer)
			}

			ret.WriteString(d.Text)
		}

		// within the conflict area, append a linebreak to the text if it is missing one
		// to make sure conflict labels are separated by new lines
		sanitized := sanitize(d.Text)

		if d.Type == diff.DiffDelete {
			if mode == modeNormal {
				mode = modeLocal
				ret.WriteString(conflictLabelLocal)
			}

			ret.WriteString(sanitized)
		}

		if d.Type == diff.DiffInsert {
			if mode == modeLocal {
				mode = modeRemote
				ret.WriteString(conflictLabelDivide)
			}

			ret.WriteString(sanitized)

			if idx == maxIdx {
				ret.WriteString(conflictLabelServer)
			}
		}
	}

	return ret.String()
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

func reportBookConflict(tx *database.DB, body, localBookUUID, serverBookUUID string) (string, error) {
	var builder strings.Builder

	var localBookName, serverBookName string
	if err := tx.QueryRow("SELECT label FROM books WHERE uuid = ?", localBookUUID).Scan(&localBookName); err != nil {
		return "", errors.Wrapf(err, "getting book label for %s", localBookUUID)
	}
	if err := tx.QueryRow("SELECT label FROM books WHERE uuid = ?", serverBookUUID).Scan(&serverBookName); err != nil {
		return "", errors.Wrapf(err, "getting book label for %s", serverBookUUID)
	}

	builder.WriteString(conflictLabelLocal)
	builder.WriteString(fmt.Sprintf("Moved to the book %s\n", localBookName))
	builder.WriteString(conflictLabelDivide)
	builder.WriteString(fmt.Sprintf("Moved to the book %s\n", serverBookName))
	builder.WriteString(conflictLabelServer)
	builder.WriteString("\n")
	builder.WriteString(body)

	return builder.String(), nil
}

func getConflictsBookUUID(tx *database.DB) (string, error) {
	var ret string

	err := tx.QueryRow("SELECT uuid FROM books WHERE label = ?", "conflicts").Scan(&ret)
	if err == sql.ErrNoRows {
		// Create a conflicts book
		ret = utils.GenerateUUID()
		b := database.NewBook(ret, "conflicts", 0, false, true)
		err = b.Insert(tx)
		if err != nil {
			tx.Rollback()
			return "", errors.Wrap(err, "creating the conflicts book")
		}
	} else if err != nil {
		return "", errors.Wrap(err, "getting uuid for conflicts book")
	}

	return ret, nil
}

// noteMergeReport holds the result of a field-by-field merge of two copies of notes
type noteMergeReport struct {
	body     string
	bookUUID string
	editedOn int64
}

// mergeNoteFields  performs a field-by-field merge between the local and the server copy. It returns a merge report
// between the local and the server copy of the note.
func mergeNoteFields(tx *database.DB, localNote database.Note, serverNote client.SyncFragNote) (*noteMergeReport, error) {
	if !localNote.Dirty {
		return &noteMergeReport{
			body:     serverNote.Body,
			bookUUID: serverNote.BookUUID,
			editedOn: serverNote.EditedOn,
		}, nil
	}

	body := reportBodyConflict(localNote.Body, serverNote.Body)

	var bookUUID string
	if serverNote.BookUUID != localNote.BookUUID {
		b, err := reportBookConflict(tx, body, localNote.BookUUID, serverNote.BookUUID)
		if err != nil {
			return nil, errors.Wrapf(err, "reporting book conflict for note %s", localNote.UUID)
		}

		body = b

		conflictsBookUUID, err := getConflictsBookUUID(tx)
		if err != nil {
			return nil, errors.Wrap(err, "getting the conflicts book uuid")
		}

		bookUUID = conflictsBookUUID
	} else {
		bookUUID = serverNote.BookUUID
	}

	ret := noteMergeReport{
		body:     body,
		bookUUID: bookUUID,
		editedOn: maxInt64(localNote.EditedOn, serverNote.EditedOn),
	}

	return &ret, nil
}
