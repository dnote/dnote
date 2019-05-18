/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package add

import (
	"database/sql"
	"time"

	"github.com/dnote/dnote/cli/core"
	"github.com/dnote/dnote/cli/infra"
	"github.com/dnote/dnote/cli/log"
	"github.com/dnote/dnote/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var reservedBookNames = []string{"trash", "conflicts"}

var content string

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
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <book>",
		Short:   "Add a new note",
		Aliases: []string{"a", "n", "new"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&content, "content", "c", "", "The new content for the note")

	return cmd
}

func isReservedName(name string) bool {
	for _, n := range reservedBookNames {
		if name == n {
			return true
		}
	}

	return false
}

// ErrBookNameReserved is an error incidating that the specified book name is reserved
var ErrBookNameReserved = errors.New("The book name is reserved")

// ErrNumericBookName is an error for book names that only contain numbers
var ErrNumericBookName = errors.New("The book name cannot contain only numbers")

func validateBookName(name string) error {
	if isReservedName(name) {
		return ErrBookNameReserved
	}

	if utils.IsNumber(name) {
		return ErrNumericBookName
	}

	return nil
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		bookName := args[0]

		if err := validateBookName(bookName); err != nil {
			return errors.Wrap(err, "invalid book name")
		}

		if content == "" {
			fpath := core.GetDnoteTmpContentPath(ctx)
			err := core.GetEditorInput(ctx, fpath, &content)
			if err != nil {
				return errors.Wrap(err, "Failed to get editor input")
			}
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

		info, err := core.GetNoteInfo(ctx, noteRowID)
		if err != nil {
			return err
		}

		core.PrintNoteInfo(info)

		if err := core.CheckUpdate(ctx); err != nil {
			log.Error(errors.Wrap(err, "automatically checking updates").Error())
		}

		return nil
	}
}

func writeNote(ctx infra.DnoteCtx, bookLabel string, content string, ts int64) (string, error) {
	tx, err := ctx.DB.Begin()
	if err != nil {
		return "", errors.Wrap(err, "beginning a transaction")
	}

	var bookUUID string
	err = tx.QueryRow("SELECT uuid FROM books WHERE label = ?", bookLabel).Scan(&bookUUID)
	if err == sql.ErrNoRows {
		bookUUID = utils.GenerateUUID()

		b := core.NewBook(bookUUID, bookLabel, 0, false, true)
		err = b.Insert(tx)
		if err != nil {
			tx.Rollback()
			return "", errors.Wrap(err, "creating the book")
		}
	} else if err != nil {
		return "", errors.Wrap(err, "finding the book")
	}

	noteUUID := utils.GenerateUUID()
	n := core.NewNote(noteUUID, bookUUID, content, ts, 0, 0, false, false, true)

	err = n.Insert(tx)
	if err != nil {
		tx.Rollback()
		return "", errors.Wrap(err, "creating the note")
	}

	var noteRowID string
	err = tx.QueryRow(`SELECT notes.rowid
			FROM notes
			WHERE notes.uuid = ?`, noteUUID).
		Scan(&noteRowID)
	if err != nil {
		return noteRowID, errors.Wrap(err, "getting the note rowid")
	}

	tx.Commit()

	return noteRowID, nil
}
