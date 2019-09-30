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

package edit

import (
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var contentFlag string
var bookFlag string
var nameFlag string

var example = `
  * Edit a note by id
  dnote edit 3

  * Edit a note without launching an editor
  dnote edit 3 -c "new content"

  * Move a note to another book
  dnote edit 3 -b javascript

  * Rename a book
  dnote edit javascript

  * Rename a book without launching an editor
  dnote edit javascript -n js
`

// NewCmd returns a new edit command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit <note id|book name>",
		Short:   "Edit a note or a book",
		Aliases: []string{"e"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&contentFlag, "content", "c", "", "a new content for the note")
	f.StringVarP(&bookFlag, "book", "b", "", "the name of the book to move the note to")
	f.StringVarP(&nameFlag, "name", "n", "", "a new name for a book")

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
		// DEPRECATED: Remove in 1.0.0
		if len(args) == 2 {
			log.Plain(log.ColorYellow.Sprintf("DEPRECATED: you no longer need to pass book name to the view command. e.g. `dnote view 123`.\n\n"))

			target := args[1]

			if err := runNote(ctx, target); err != nil {
				return errors.Wrap(err, "editing note")
			}

			return nil
		}

		target := args[0]

		if utils.IsNumber(target) {
			if err := runNote(ctx, target); err != nil {
				return errors.Wrap(err, "editing note")
			}
		} else {
			if err := runBook(ctx, target); err != nil {
				return errors.Wrap(err, "editing book")
			}
		}

		return nil
	}
}
