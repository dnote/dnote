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

package view

import (
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/dnote/dnote/pkg/cli/cmd/cat"
	"github.com/dnote/dnote/pkg/cli/cmd/ls"
	"github.com/dnote/dnote/pkg/cli/utils"
)

var example = `
 * View all books
 dnote view

 * List notes in a book
 dnote view javascript

 * View a particular note in a book
 dnote view javascript 0
 `

var nameOnly bool
var contentOnly bool

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) > 2 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

// NewCmd returns a new view command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <book name?> <note index?>",
		Aliases: []string{"v"},
		Short:   "List books, notes or view a content",
		Example: example,
		RunE:    newRun(ctx),
		PreRunE: preRun,
	}

	f := cmd.Flags()
	f.BoolVarP(&nameOnly, "name-only", "", false, "print book names only")
	f.BoolVarP(&contentOnly, "content-only", "", false, "print the note content only")

	return cmd
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var run infra.RunEFunc

		if len(args) == 0 {
			run = ls.NewRun(ctx, nameOnly)
		} else if len(args) == 1 {
			if nameOnly {
				return errors.New("--name-only flag is only valid when viewing books")
			}

			if utils.IsNumber(args[0]) {
				run = cat.NewRun(ctx, contentOnly)
			} else {
				run = ls.NewRun(ctx, false)
			}
		} else if len(args) == 2 {
			// DEPRECATED: passing book name to view command is deprecated
			run = cat.NewRun(ctx, false)
		} else {
			return errors.New("Incorrect number of arguments")
		}

		return run(cmd, args)
	}
}
