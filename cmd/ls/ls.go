package ls

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
 * List all books
 dnote ls

 * List notes in a book
 dnote ls javascript
 `

var deprecationWarning = `and "view" will replace it in v1.0.0.

Run "dnote view --help" for more information.
`

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

// NewCmd returns a new ls command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:        "ls <book name?>",
		Aliases:    []string{"l", "notes"},
		Short:      "List all notes",
		Example:    example,
		RunE:       NewRun(ctx),
		PreRunE:    preRun,
		Deprecated: deprecationWarning,
	}

	return cmd
}

// NewRun returns a new run function for ls
func NewRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if err := printBooks(ctx); err != nil {
				return errors.Wrap(err, "viewing books")
			}

			return nil
		}

		bookName := args[0]
		if err := printNotes(ctx, bookName); err != nil {
			return errors.Wrapf(err, "viewing book '%s'", bookName)
		}

		return nil
	}
}

// bookInfo is an information about the book to be printed on screen
type bookInfo struct {
	BookLabel string
	NoteCount int
}

// noteInfo is an information about the note to be printed on screen
type noteInfo struct {
	RowID int
	Body  string
}

// getNewlineIdx returns the index of newline character in a string
func getNewlineIdx(str string) int {
	var ret int

	ret = strings.Index(str, "\n")

	if ret == -1 {
		ret = strings.Index(str, "\r\n")
	}

	return ret
}

// formatBody returns an excerpt of the given raw note content and a boolean
// indicating if the returned string has been excertped
func formatBody(noteBody string) (string, bool) {
	newlineIdx := getNewlineIdx(noteBody)

	if newlineIdx > -1 {
		ret := strings.Trim(noteBody[0:newlineIdx], " ")

		return ret, true
	}

	return strings.Trim(noteBody, " "), false
}

func printBooks(ctx infra.DnoteCtx) error {
	db := ctx.DB

	rows, err := db.Query(`SELECT books.label, count(notes.uuid) note_count
	FROM books
	INNER JOIN notes ON notes.book_uuid = books.uuid
	GROUP BY books.uuid
	ORDER BY books.label ASC;`)
	if err != nil {
		return errors.Wrap(err, "querying books")
	}
	defer rows.Close()

	infos := []bookInfo{}
	for rows.Next() {
		var info bookInfo
		err = rows.Scan(&info.BookLabel, &info.NoteCount)
		if err != nil {
			return errors.Wrap(err, "scanning a row")
		}

		infos = append(infos, info)
	}

	for _, info := range infos {
		log.Printf("%s %s\n", info.BookName, log.ColorYellow.Sprintf("(%d)", info.NoteCount))
	}

	return nil
}

func printNotes(ctx infra.DnoteCtx, bookName string) error {
	db := ctx.DB

	var bookUUID string
	err := db.QueryRow("SELECT uuid FROM books WHERE label = ?", bookName).Scan(&bookUUID)
	if err == sql.ErrNoRows {
		return errors.New("book not found")
	} else if err != nil {
		return errors.Wrap(err, "querying the book")
	}

	rows, err := db.Query(`SELECT rowid, body FROM notes WHERE book_uuid = ? ORDER BY added_on ASC;`, bookUUID)
	if err != nil {
		return errors.Wrap(err, "querying notes")
	}
	defer rows.Close()

	infos := []noteInfo{}
	for rows.Next() {
		var info noteInfo
		err = rows.Scan(&info.RowID, &info.Body)
		if err != nil {
			return errors.Wrap(err, "scanning a row")
		}

		infos = append(infos, info)
	}

	log.Infof("on book %s\n", bookName)

	for _, info := range infos {
		body, isExcerpt := formatBody(info.Body)

		rowid := log.SprintfYellow("(%d)", info.RowID)
		if isExcerpt {
			body = fmt.Sprintf("%s %s", body, log.SprintfYellow("[---More---]"))
		}

		log.Plainf("%s %s\n", rowid, body)
	}

	return nil
}
