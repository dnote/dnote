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
		RunE:       NewRun(ctx, false),
		PreRunE:    preRun,
		Deprecated: deprecationWarning,
	}

	return cmd
}

// NewRun returns a new run function
func NewRun(ctx context.DnoteCtx, contentOnly bool) infra.RunEFunc {
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

		if contentOnly {
			output.NoteContent(info)
		} else {
			output.NoteInfo(info)
		}

		return nil
	}
}
