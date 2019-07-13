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
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var newContent string

var example = `
  * Edit the note by its id
  dnote edit 3

	* Skip the prompt by providing new content directly
	dnote edit 3 -c "new content"`

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

	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 && len(args) != 2 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		db := ctx.DB

		var noteRowID string

		if len(args) == 2 {
			log.Plain(log.ColorYellow.Sprintf("DEPRECATED: you no longer need to pass book name to the view command. e.g. `dnote view 123`.\n\n"))

			noteRowID = args[1]
		} else {
			noteRowID = args[0]
		}

		var noteUUID, oldContent string
		err := db.QueryRow("SELECT uuid, body FROM notes WHERE rowid = ? AND deleted = false", noteRowID).Scan(&noteUUID, &oldContent)
		if err == sql.ErrNoRows {
			return errors.Errorf("note %s not found", noteRowID)
		} else if err != nil {
			return errors.Wrap(err, "querying the book")
		}

		if newContent == "" {
			fpath, err := ui.GetTmpContentPath(ctx)
			if err != nil {
				return errors.Wrap(err, "getting temporarily content file path")
			}

			e := ioutil.WriteFile(fpath, []byte(oldContent), 0644)
			if e != nil {
				return errors.Wrap(e, "preparing tmp content file")
			}

			e = ui.GetEditorInput(ctx, fpath, &newContent)
			if e != nil {
				return errors.Wrap(err, "getting editor input")
			}
		}

		if oldContent == newContent {
			return errors.New("Nothing changed")
		}

		ts := time.Now().UnixNano()
		newContent = ui.SanitizeContent(newContent)

		tx, err := db.Begin()
		if err != nil {
			return errors.Wrap(err, "beginning a transaction")
		}

		_, err = tx.Exec(`UPDATE notes
			SET body = ?, edited_on = ?, dirty = ?
			WHERE rowid = ?`, newContent, ts, true, noteRowID)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "updating the note")
		}

		tx.Commit()

		log.Success("edited the note\n")
		fmt.Printf("\n------------------------content------------------------\n")
		fmt.Printf("%s", newContent)
		fmt.Printf("\n-------------------------------------------------------\n")

		return nil
	}
}
