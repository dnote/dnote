package cat

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
 * See the notes with index 2 from a book 'javascript'
 dnote cat javascript 2
 `

var deprecationWarning = `and "view" will replace it in v0.5.0.

 Run "dnote view --help" for more information.
`

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Incorrect number of arguments")
	}

	return nil
}

// NewCmd returns a new cat command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
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

type noteInfo struct {
	BookLabel string
	UUID      string
	Content   string
	AddedOn   int64
	EditedOn  int64
}

func NewRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		db := ctx.DB
		bookLabel := args[0]
		noteID := args[1]

		var bookUUID string
		err := db.QueryRow("SELECT uuid FROM books WHERE label = ?", bookLabel).Scan(&bookUUID)
		if err == sql.ErrNoRows {
			return errors.Errorf("book '%s' not found", bookLabel)
		} else if err != nil {
			return errors.Wrap(err, "querying the book")
		}

		var info noteInfo
		err = db.QueryRow(`SELECT books.label, notes.uuid, notes.content, notes.added_on, notes.edited_on
			FROM notes
			INNER JOIN books ON books.uuid = notes.book_uuid
			WHERE notes.id = ? AND books.uuid = ?`, noteID, bookUUID).
			Scan(&info.BookLabel, &info.UUID, &info.Content, &info.AddedOn, &info.EditedOn)
		if err == sql.ErrNoRows {
			return errors.Errorf("note %s not found in the book '%s'", noteID, bookLabel)
		} else if err != nil {
			return errors.Wrap(err, "querying the note")
		}

		log.Infof("book name: %s\n", info.BookLabel)
		log.Infof("note uuid: %s\n", info.UUID)
		log.Infof("created at: %s\n", time.Unix(info.AddedOn, 0).Format("Jan 2, 2006 3:04pm (MST)"))
		if info.EditedOn != 0 {
			log.Infof("updated at: %s\n", time.Unix(info.EditedOn, 0).Format("Jan 2, 2006 3:04pm (MST)"))
		}
		fmt.Printf("\n------------------------content------------------------\n")
		fmt.Printf("%s", info.Content)
		fmt.Printf("\n-------------------------------------------------------\n")

		return nil
	}
}
