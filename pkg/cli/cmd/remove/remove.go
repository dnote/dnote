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

package remove

import (
	"fmt"
	"strconv"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/output"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var bookFlag string
var yesFlag bool

var example = `
  * Delete a note by id
  dnote delete 2

  * Delete a book by name
  dnote delete js
`

// NewCmd returns a new remove command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <note id|book name>",
		Short:   "Remove a note or a book",
		Aliases: []string{"rm", "d", "delete"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&bookFlag, "book", "b", "", "The book name to delete")
	f.BoolVarP(&yesFlag, "yes", "y", false, "Assume yes to the prompts and run in non-interactive mode")

	f.MarkDeprecated("book", "Pass the book name as an argument. e.g. `dnote rm book_name`")

	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 && len(args) != 2 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

func maybeConfirm(message string, defaultValue bool) (bool, error) {
	if yesFlag {
		return true, nil
	}

	return ui.Confirm(message, defaultValue)
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		// DEPRECATED: Remove in 1.0.0
		if bookFlag != "" {
			if err := runBook(ctx, bookFlag); err != nil {
				return errors.Wrap(err, "removing the book")
			}

			return nil
		}

		// DEPRECATED: Remove in 1.0.0
		if len(args) == 2 {
			log.Plain(log.ColorYellow.Sprintf("DEPRECATED: you no longer need to pass book name to the remove command. e.g. `dnote remove 123`.\n\n"))

			target := args[1]
			if err := runNote(ctx, target); err != nil {
				return errors.Wrap(err, "removing the note")
			}

			return nil
		}

		target := args[0]

		if utils.IsNumber(target) {
			if err := runNote(ctx, target); err != nil {
				return errors.Wrap(err, "removing the note")
			}
		} else {
			if err := runBook(ctx, target); err != nil {
				return errors.Wrap(err, "removing the book")
			}
		}

		return nil
	}
}

func runNote(ctx context.DnoteCtx, rowIDArg string) error {
	db := ctx.DB

	noteRowID, err := strconv.Atoi(rowIDArg)
	if err != nil {
		return errors.Wrap(err, "invalid rowid")
	}

	noteInfo, err := database.GetNoteInfo(db, noteRowID)
	if err != nil {
		return err
	}

	output.NoteInfo(noteInfo)

	ok, err := maybeConfirm("remove this note?", false)
	if err != nil {
		return errors.Wrap(err, "getting confirmation")
	}
	if !ok {
		log.Warnf("aborted by user\n")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	if _, err = tx.Exec("UPDATE notes SET deleted = ?, dirty = ?, body = ? WHERE uuid = ?", true, true, "", noteInfo.UUID); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "removing the note")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "comitting transaction")
	}

	log.Successf("removed from %s\n", noteInfo.BookLabel)

	return nil
}

func runBook(ctx context.DnoteCtx, bookLabel string) error {
	db := ctx.DB

	bookUUID, err := database.GetBookUUID(db, bookLabel)
	if err != nil {
		return errors.Wrap(err, "finding book uuid")
	}

	ok, err := maybeConfirm(fmt.Sprintf("delete book '%s' and all its notes?", bookLabel), false)
	if err != nil {
		return errors.Wrap(err, "getting confirmation")
	}
	if !ok {
		log.Warnf("aborted by user\n")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	if _, err = tx.Exec("UPDATE notes SET deleted = ?, dirty = ?, body = ? WHERE book_uuid = ?", true, true, "", bookUUID); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "removing notes in the book")
	}

	// override the label with a random string
	uniqLabel := utils.GenerateUUID()
	if _, err = tx.Exec("UPDATE books SET deleted = ?, dirty = ?, label = ? WHERE uuid = ?", true, true, uniqLabel, bookUUID); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "removing the book")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "committing transaction")
	}

	log.Success("removed book\n")

	return nil
}
