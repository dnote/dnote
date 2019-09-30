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

package add

import (
	"database/sql"
	"time"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/output"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/dnote/dnote/pkg/cli/upgrade"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/dnote/dnote/pkg/cli/validate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var contentFlag string

var example = `
 * Open an editor to write content
 dnote add git

 * Skip the editor by providing content directly
 dnote add git -c "time is a part of the commit hash"`

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

// NewCmd returns a new add command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <book>",
		Short:   "Add a new note",
		Aliases: []string{"a", "n", "new"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&contentFlag, "content", "c", "", "The new content for the note")

	return cmd
}

func getContent(ctx context.DnoteCtx) (string, error) {
	if contentFlag != "" {
		return contentFlag, nil
	}

	fpath, err := ui.GetTmpContentPath(ctx)
	if err != nil {
		return "", errors.Wrap(err, "getting temporarily content file path")
	}

	c, err := ui.GetEditorInput(ctx, fpath)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get editor input")
	}

	return c, nil
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		bookName := args[0]
		if err := validate.BookName(bookName); err != nil {
			return errors.Wrap(err, "invalid book name")
		}

		content, err := getContent(ctx)
		if err != nil {
			return errors.Wrap(err, "getting content")
		}
		if content == "" {
			return errors.New("Empty content")
		}

		ts := time.Now().UnixNano()
		noteRowID, err := writeNote(ctx, bookName, content, ts)
		if err != nil {
			return errors.Wrap(err, "Failed to write note")
		}

		log.Successf("added to %s\n", bookName)

		db := ctx.DB
		info, err := database.GetNoteInfo(db, noteRowID)
		if err != nil {
			return err
		}

		output.NoteInfo(info)

		if err := upgrade.Check(ctx); err != nil {
			log.Error(errors.Wrap(err, "automatically checking updates").Error())
		}

		return nil
	}
}

func writeNote(ctx context.DnoteCtx, bookLabel string, content string, ts int64) (int, error) {
	tx, err := ctx.DB.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "beginning a transaction")
	}

	var bookUUID string
	err = tx.QueryRow("SELECT uuid FROM books WHERE label = ?", bookLabel).Scan(&bookUUID)
	if err == sql.ErrNoRows {
		bookUUID = utils.GenerateUUID()

		b := database.NewBook(bookUUID, bookLabel, 0, false, true)
		err = b.Insert(tx)
		if err != nil {
			tx.Rollback()
			return 0, errors.Wrap(err, "creating the book")
		}
	} else if err != nil {
		return 0, errors.Wrap(err, "finding the book")
	}

	noteUUID := utils.GenerateUUID()
	n := database.NewNote(noteUUID, bookUUID, content, ts, 0, 0, false, false, true)

	err = n.Insert(tx)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "creating the note")
	}

	var noteRowID int
	err = tx.QueryRow(`SELECT notes.rowid
			FROM notes
			WHERE notes.uuid = ?`, noteUUID).
		Scan(&noteRowID)
	if err != nil {
		tx.Rollback()
		return noteRowID, errors.Wrap(err, "getting the note rowid")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return noteRowID, errors.Wrap(err, "committing a transaction")
	}

	return noteRowID, nil
}
