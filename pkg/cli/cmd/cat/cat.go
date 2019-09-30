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

package cat

import (
	"strconv"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/output"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
 * See the notes with index 2 from a book 'javascript'
 dnote cat javascript 2
 `

var deprecationWarning = `and "view" will replace it in the future version.

 Run "dnote view --help" for more information.
`

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Incorrect number of arguments")
	}

	return nil
}

// NewCmd returns a new cat command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:        "cat <book name> <note index>",
		Aliases:    []string{"c"},
		Short:      "See a note",
		Example:    example,
		RunE:       NewRun(ctx),
		PreRunE:    preRun,
		Deprecated: deprecationWarning,
	}

	return cmd
}

// NewRun returns a new run function
func NewRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
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

		db := ctx.DB
		info, err := database.GetNoteInfo(db, noteRowID)
		if err != nil {
			return err
		}

		output.NoteInfo(info)

		return nil
	}
}
