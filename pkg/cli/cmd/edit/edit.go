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
