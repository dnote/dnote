package add

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var content string

var example = `
 * Open an editor to write content
 dnote add git

 * Skip the editor by providing content directly
 dnote add git -c "time is a part of the commit hash"`

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

// NewCmd returns a new add command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <content>",
		Short:   "Add a note",
		Aliases: []string{"a", "n", "new"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&content, "content", "c", "", "The new content for the note")

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		bookName := args[0]

		if content == "" {
			fpath := core.GetDnoteTmpContentPath(ctx)
			err := core.GetEditorInput(ctx, fpath, &content)
			if err != nil {
				return errors.Wrap(err, "Failed to get editor input")
			}
		}

		if content == "" {
			return errors.New("Empty content")
		}

		ts := time.Now().Unix()
		err := writeNote(ctx, bookName, content, ts)
		if err != nil {
			return errors.Wrap(err, "Failed to write note")
		}

		log.Successf("added to %s\n", bookName)
		fmt.Printf("\n------------------------content------------------------\n")
		fmt.Printf("%s", content)
		fmt.Printf("\n-------------------------------------------------------\n")

		if err := core.CheckUpdate(ctx); err != nil {
			log.Error(errors.Wrap(err, "automatically checking updates").Error())
		}

		return nil
	}
}

func writeNote(ctx infra.DnoteCtx, bookLabel string, content string, ts int64) error {
	tx, err := ctx.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	var bookUUID string
	err = tx.QueryRow("SELECT uuid FROM books WHERE label = ?", bookLabel).Scan(&bookUUID)
	if err == sql.ErrNoRows {
		bookUUID = utils.GenerateUUID()
		_, err = tx.Exec("INSERT INTO books (uuid, label) VALUES (?, ?)", bookUUID, bookLabel)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "creating the book")
		}

		err = core.LogActionAddBook(tx, bookLabel)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "logging action")
		}
	} else if err != nil {
		return errors.Wrap(err, "finding the book")
	}

	noteUUID := utils.GenerateUUID()
	_, err = tx.Exec(`INSERT INTO notes (uuid, book_uuid, content, added_on, public)
		VALUES (?, ?, ?, ?, ?);`, noteUUID, bookUUID, content, ts, false)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "creating the note")
	}
	err = core.LogActionAddNote(tx, noteUUID, bookLabel, content, ts)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "logging action")
	}

	tx.Commit()

	return nil
}
