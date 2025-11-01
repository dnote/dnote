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

package add

import (
	"database/sql"
	"time"
	"os"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/output"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/dnote/dnote/pkg/cli/upgrade"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/dnote/dnote/pkg/cli/validate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var contentFlag string

var example = `
 * Open an editor to write content
 dnote add git

 * Skip the editor by providing content directly
 dnote add git -c "time is a part of the commit hash"

 * Send stdin content to a note
 echo "a branch is just a pointer to a commit" | dnote add git
 # or
 dnote add git << EOF
 pull is fetch with a merge
 EOF`

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

// NewCmd returns a new add command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <book>",
		Short:   "Add a new note",
		Aliases: []string{"a", "n", "new"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&contentFlag, "content", "c", "", "The new content for the note")

	return cmd
}

func getContent(ctx context.DnoteCtx) (string, error) {
	if contentFlag != "" {
		return contentFlag, nil
	}

	// check for piped content
	fInfo, _ := os.Stdin.Stat()
	if fInfo.Mode() & os.ModeCharDevice == 0 {
		c, err := ui.ReadStdInput()
		if err != nil {
			return "", errors.Wrap(err, "Failed to get piped input")
		}
		return c, nil
	}

	fpath, err := ui.GetTmpContentPath(ctx)
	if err != nil {
		return "", errors.Wrap(err, "getting temporarily content file path")
	}

	c, err := ui.GetEditorInput(ctx, fpath)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get editor input")
	}

	return c, nil
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		bookName := args[0]
		if err := validate.BookName(bookName); err != nil {
			return errors.Wrap(err, "invalid book name")
		}

		content, err := getContent(ctx)
		if err != nil {
			return errors.Wrap(err, "getting content")
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

		db := ctx.DB
		info, err := database.GetNoteInfo(db, noteRowID)
		if err != nil {
			return err
		}

		output.NoteInfo(info)

		if err := upgrade.Check(ctx); err != nil {
			log.Error(errors.Wrap(err, "automatically checking updates").Error())
		}

		return nil
	}
}

func writeNote(ctx context.DnoteCtx, bookLabel string, content string, ts int64) (int, error) {
	tx, err := ctx.DB.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "beginning a transaction")
	}

	var bookUUID string
	err = tx.QueryRow("SELECT uuid FROM books WHERE label = ?", bookLabel).Scan(&bookUUID)
	if err == sql.ErrNoRows {
		bookUUID, err = utils.GenerateUUID()
		if err != nil {
			return 0, errors.Wrap(err, "generating uuid")
		}

		b := database.NewBook(bookUUID, bookLabel, 0, false, true)
		err = b.Insert(tx)
		if err != nil {
			tx.Rollback()
			return 0, errors.Wrap(err, "creating the book")
		}
	} else if err != nil {
		return 0, errors.Wrap(err, "finding the book")
	}

	noteUUID, err := utils.GenerateUUID()
	if err != nil {
		return 0, errors.Wrap(err, "generating uuid")
	}

	n := database.NewNote(noteUUID, bookUUID, content, ts, 0, 0, false, true)

	err = n.Insert(tx)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "creating the note")
	}

	var noteRowID int
	err = tx.QueryRow(`SELECT notes.rowid
			FROM notes
			WHERE notes.uuid = ?`, noteUUID).
		Scan(&noteRowID)
	if err != nil {
		tx.Rollback()
		return noteRowID, errors.Wrap(err, "getting the note rowid")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return noteRowID, errors.Wrap(err, "committing a transaction")
	}

	return noteRowID, nil
}
