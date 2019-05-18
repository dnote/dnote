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

package view

import (
	"github.com/dnote/dnote/cli/core"
	"github.com/dnote/dnote/cli/infra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/dnote/dnote/cli/cmd/cat"
	"github.com/dnote/dnote/cli/cmd/ls"
	"github.com/dnote/dnote/cli/utils"
)

var example = `
 * View all books
 dnote view

 * List notes in a book
 dnote view javascript

 * View a particular note in a book
 dnote view javascript 0
 `

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) > 2 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

// NewCmd returns a new view command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <book name?> <note index?>",
		Aliases: []string{"v"},
		Short:   "List books, notes or view a content",
		Example: example,
		RunE:    newRun(ctx),
		PreRunE: preRun,
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var run core.RunEFunc

		if len(args) == 0 {
			run = ls.NewRun(ctx)
		} else if len(args) == 1 {
			if utils.IsNumber(args[0]) {
				run = cat.NewRun(ctx)
			} else {
				run = ls.NewRun(ctx)
			}
		} else if len(args) == 2 {
			// DEPRECATED: passing book name to view command is deprecated
			run = cat.NewRun(ctx)
		} else {
			return errors.New("Incorrect number of arguments")
		}

		return run(cmd, args)
	}
}
