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

package edit

import (
	"database/sql"
	"io/ioutil"
	"strconv"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/output"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var newContent string
var bookName string

var example = `
  * Edit the note by its id
  dnote edit 3

  * Skip the prompt by providing new content directly
  dnote edit 3 -c "new content"

  * Move a note to another book
  dnote edit 3 -b javascript
`

// NewCmd returns a new edit command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Edit a note",
		Aliases: []string{"e"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&newContent, "content", "c", "", "The new content for the note")
	f.StringVarP(&bookName, "book", "b", "", "The name of the book to move the note to")

	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 && len(args) != 2 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

func changeContent(ctx context.DnoteCtx, note database.Note) error {
	if newContent == "" {
		fpath, err := ui.GetTmpContentPath(ctx)
		if err != nil {
			return errors.Wrap(err, "getting temporarily content file path")
		}

		if err := ioutil.WriteFile(fpath, []byte(note.Body), 0644); err != nil {
			return errors.Wrap(err, "preparing tmp content file")
		}

		if err := ui.GetEditorInput(ctx, fpath, &newContent); err != nil {
			return errors.Wrap(err, "getting editor input")
		}
	}

	if note.Body == newContent {
		return errors.New("Nothing changed")
	}

	newContent = ui.SanitizeContent(newContent)

	if err := database.UpdateNoteContent(ctx.DB, ctx.Clock, note.RowID, newContent); err != nil {
		return errors.Wrap(err, "updating the note")
	}

	return nil
}

func moveBook(ctx context.DnoteCtx, note database.Note, bookName string) error {
	db := ctx.DB

	targetBookUUID, err := database.GetBookUUID(db, bookName)
	if err != nil {
		return errors.Wrap(err, "finding book uuid")
	}

	if note.BookUUID == targetBookUUID {
		return errors.New("book has not changed")
	}

	if err := database.UpdateNoteBook(db, ctx.Clock, note.RowID, targetBookUUID); err != nil {
		return errors.Wrap(err, "moving book")
	}

	return nil
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		db := ctx.DB

		var noteRowIDArg string

		if len(args) == 2 {
			log.Plain(log.ColorYellow.Sprintf("DEPRECATED: you no longer need to pass book name to the view command. e.g. `dnote view 123`.\n\n"))

			noteRowIDArg = args[1]
		} else {
			noteRowIDArg = args[0]
		}

		noteRowID, err := strconv.Atoi(noteRowIDArg)
		if err != nil {
			return errors.Wrap(err, "invalid rowid")
		}

		note, err := database.GetActiveNote(db, noteRowID)
		if err == sql.ErrNoRows {
			return errors.Errorf("note %d not found", noteRowID)
		} else if err != nil {
			return errors.Wrap(err, "querying the book")
		}

		if bookName != "" {
			if err := moveBook(ctx, note, bookName); err != nil {
				return errors.Wrap(err, "moving book")
			}
		} else {
			if err := changeContent(ctx, note); err != nil {
				return errors.Wrap(err, "changing content")
			}
		}

		noteInfo, err := database.GetNoteInfo(db, noteRowID)
		if err != nil {
			return errors.Wrap(err, "getting note info")
		}

		log.Success("edited the note\n")
		output.NoteInfo(noteInfo)

		return nil
	}
}
