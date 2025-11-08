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

package edit

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/output"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/pkg/errors"
)

func validateRunNoteFlags() error {
	if nameFlag != "" {
		return errors.New("--name is invalid for editing a book")
	}

	return nil
}

func waitEditorNoteContent(ctx context.DnoteCtx, note database.Note) (string, error) {
	fpath, err := ui.GetTmpContentPath(ctx)
	if err != nil {
		return "", errors.Wrap(err, "getting temporarily content file path")
	}

	if err := os.WriteFile(fpath, []byte(note.Body), 0644); err != nil {
		return "", errors.Wrap(err, "preparing tmp content file")
	}

	c, err := ui.GetEditorInput(ctx, fpath)
	if err != nil {
		return "", errors.Wrap(err, "getting editor input")
	}

	return c, nil
}

func getContent(ctx context.DnoteCtx, note database.Note) (string, error) {
	if contentFlag != "" {
		return contentFlag, nil
	}

	c, err := waitEditorNoteContent(ctx, note)
	if err != nil {
		return "", errors.Wrap(err, "getting content from editor")
	}

	return c, nil
}

func changeContent(ctx context.DnoteCtx, tx *database.DB, note database.Note, content string) error {
	if note.Body == content {
		return errors.New("Nothing changed")
	}

	if err := database.UpdateNoteContent(tx, ctx.Clock, note.RowID, content); err != nil {
		return errors.Wrap(err, "updating the note")
	}

	return nil
}

func moveBook(ctx context.DnoteCtx, tx *database.DB, note database.Note, bookName string) error {
	targetBookUUID, err := database.GetBookUUID(tx, bookName)
	if err != nil {
		return errors.Wrap(err, "finding book uuid")
	}

	if note.BookUUID == targetBookUUID {
		return errors.New("book has not changed")
	}

	if err := database.UpdateNoteBook(tx, ctx.Clock, note.RowID, targetBookUUID); err != nil {
		return errors.Wrap(err, "moving book")
	}

	return nil
}

func updateNote(ctx context.DnoteCtx, tx *database.DB, note database.Note, bookName, content string) error {
	if bookName != "" {
		if err := moveBook(ctx, tx, note, bookName); err != nil {
			return errors.Wrap(err, "moving book")
		}
	}
	if content != "" {
		if err := changeContent(ctx, tx, note, content); err != nil {
			return errors.Wrap(err, "changing content")
		}
	}

	return nil
}

func runNote(ctx context.DnoteCtx, rowIDArg string) error {
	err := validateRunNoteFlags()
	if err != nil {
		return errors.Wrap(err, "validating flags.")
	}

	rowID, err := strconv.Atoi(rowIDArg)
	if err != nil {
		return errors.Wrap(err, "invalid rowid")
	}

	db := ctx.DB
	note, err := database.GetActiveNote(db, rowID)
	if err == sql.ErrNoRows {
		return errors.Errorf("note %d not found", rowID)
	} else if err != nil {
		return errors.Wrap(err, "querying the book")
	}

	content := contentFlag

	// If no flag was provided, launch an editor to get the content
	if bookFlag == "" && contentFlag == "" {
		c, err := getContent(ctx, note)
		if err != nil {
			return errors.Wrap(err, "getting content from editor")
		}

		content = c
	}

	tx, err := ctx.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	err = updateNote(ctx, tx, note, bookFlag, content)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "updating note fields")
	}

	noteInfo, err := database.GetNoteInfo(tx, rowID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "getting note info")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "committing a transaction")
	}

	log.Success("edited the note\n")
	output.NoteInfo(os.Stdout, noteInfo)

	return nil
}
