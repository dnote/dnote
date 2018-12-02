package remove

import (
	"database/sql"
	"fmt"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var targetBookName string

var example = `
  * Delete a note by its index from a book
  dnote delete js 2

  * Delete a book
  dnote delete -b js`

// NewCmd returns a new remove command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a note or a book",
		Aliases: []string{"rm", "d", "delete"},
		Example: example,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&targetBookName, "book", "b", "", "The book name to delete")

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		if targetBookName != "" {
			if err := removeBook(ctx, targetBookName); err != nil {
				return errors.Wrap(err, "removing the book")
			}

			return nil
		}

		if len(args) < 2 {
			return errors.New("Missing argument")
		}

		targetBook := args[0]
		noteRowID := args[1]

		if err := removeNote(ctx, noteRowID, targetBook); err != nil {
			return errors.Wrap(err, "removing the note")
		}

		return nil
	}
}

func removeNote(ctx infra.DnoteCtx, noteRowID, bookLabel string) error {
	db := ctx.DB

	bookUUID, err := core.GetBookUUID(ctx, bookLabel)
	if err != nil {
		return errors.Wrap(err, "finding book uuid")
	}

	var noteUUID, noteContent string
	err = db.QueryRow("SELECT uuid, body FROM notes WHERE rowid = ? AND book_uuid = ?", noteRowID, bookUUID).Scan(&noteUUID, &noteContent)
	if err == sql.ErrNoRows {
		return errors.Errorf("note %s not found in the book '%s'", noteRowID, bookLabel)
	} else if err != nil {
		return errors.Wrap(err, "querying the book")
	}

	// todo: multiline
	log.Printf("body: \"%s\"\n", noteContent)

	ok, err := utils.AskConfirmation("remove this note?", false)
	if err != nil {
		return errors.Wrap(err, "getting confirmation")
	}
	if !ok {
		log.Warnf("aborted by user\n")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	if _, err = tx.Exec("UPDATE notes SET deleted = ?, dirty = ?, body = ? WHERE uuid = ? AND book_uuid = ?", true, true, "", noteUUID, bookUUID); err != nil {
		return errors.Wrap(err, "removing the note")
	}
	tx.Commit()

	log.Successf("removed from %s\n", bookLabel)

	return nil
}

func removeBook(ctx infra.DnoteCtx, bookLabel string) error {
	db := ctx.DB

	bookUUID, err := core.GetBookUUID(ctx, bookLabel)
	if err != nil {
		return errors.Wrap(err, "finding book uuid")
	}

	ok, err := utils.AskConfirmation(fmt.Sprintf("delete book '%s' and all its notes?", bookLabel), false)
	if err != nil {
		return errors.Wrap(err, "getting confirmation")
	}
	if !ok {
		log.Warnf("aborted by user\n")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	if _, err = tx.Exec("UPDATE notes SET deleted = ?, dirty = ?, body = ? WHERE book_uuid = ?", true, true, "", bookUUID); err != nil {
		return errors.Wrap(err, "removing notes in the book")
	}

	// override the label with a random string
	uniqLabel := utils.GenerateUUID()
	if _, err = tx.Exec("UPDATE books SET deleted = ?, dirty = ?, label = ? WHERE uuid = ?", true, true, uniqLabel, bookUUID); err != nil {
		return errors.Wrap(err, "removing the book")
	}

	tx.Commit()

	log.Success("removed book\n")

	return nil
}
