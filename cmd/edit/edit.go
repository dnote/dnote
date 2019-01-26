package edit

import (
	"database/sql"
	"io/ioutil"
	"time"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var newContent string

var example = `
  * Edit the note by index in a book
  dnote edit js 3

	* Skip the prompt by providing new content directly
	dntoe edit js 3 -c "new content"`

// NewCmd returns a new edit command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Edit a note or a book",
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
	if len(args) != 2 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		db := ctx.DB
		bookLabel := args[0]
		noteRowID := args[1]

		bookUUID, err := core.GetBookUUID(ctx, bookLabel)
		if err != nil {
			return errors.Wrap(err, "finding book uuid")
		}

		var noteUUID, oldContent string
		err = db.QueryRow("SELECT uuid, body FROM notes WHERE rowid = ? AND book_uuid = ?", noteRowID, bookUUID).Scan(&noteUUID, &oldContent)
		if err == sql.ErrNoRows {
			return errors.Errorf("note %s not found in the book '%s'", noteRowID, bookLabel)
		} else if err != nil {
			return errors.Wrap(err, "querying the book")
		}

		if newContent == "" {
			fpath := core.GetDnoteTmpContentPath(ctx)

			err := ioutil.WriteFile(fpath, []byte(oldContent), 0644)
			if err != nil {
				return errors.Wrap(e, "preparing tmp content file")
			}

			e = core.GetEditorInput(ctx, fpath, &newContent)
			if e != nil {
				return errors.Wrap(err, "getting editor input")
			}
		}

		if oldContent == newContent {
			return errors.New("Nothing changed")
		}

		ts := time.Now().UnixNano()
		newContent = core.SanitizeContent(newContent)

		tx, err := db.Begin()
		if err != nil {
			return errors.Wrap(err, "beginning a transaction")
		}

		_, err = tx.Exec(`UPDATE notes
			SET body = ?, edited_on = ?, dirty = ?
			WHERE rowid = ? AND book_uuid = ?`, newContent, ts, true, noteRowID, bookUUID)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "updating the note")
		}

		tx.Commit()

		log.Success("edited the note\n")
		log.PrintContent(content)

		return nil
	}
}
