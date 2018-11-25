package view

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <book name?> <note index?>",
		Aliases: []string{"v"},
		Short:   "List books, notes or view a note",
		Example: example,
		RunE:    newRun(ctx),
		PreRunE: preRun,
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var err error
		if len(args) == 0 {
			err = listBooks(ctx)
		} else if len(args) == 1 {
			err = listNotes(ctx, args[0])
		} else if len(args) == 2 {
			err = viewNote(ctx, args[0], args[1])
		} else {
			return errors.New("Incorrect number of arguments")
		}

		return err
	}
}

type noteInfo struct {
	BookLabel string
	ID        int
	UUID      string
	Content   string
	AddedOn   int64
	EditedOn  int64
}

func viewNote(ctx infra.DnoteCtx, bookLabel, noteID string) error {
	db := ctx.DB

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
	log.Infof("created at: %s\n", time.Unix(0, info.AddedOn).Format("Jan 2, 2006 3:04pm (MST)"))
	if info.EditedOn != 0 {
		log.Infof("updated at: %s\n", time.Unix(0, info.EditedOn).Format("Jan 2, 2006 3:04pm (MST)"))
	}
	fmt.Printf("\n------------------------content------------------------\n")
	fmt.Printf("%s", info.Content)
	fmt.Printf("\n-------------------------------------------------------\n")

	return nil
}

// bookInfo is an information about the book to be printed on screen
type bookInfo struct {
	BookLabel string
	NoteCount int
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

// formatContent returns an excerpt of the given raw note content and a boolean
// indicating if the returned string has been excertped
func formatContent(noteContent string) (string, bool) {
	newlineIdx := getNewlineIdx(noteContent)

	if newlineIdx > -1 {
		ret := strings.Trim(noteContent[0:newlineIdx], " ")

		return ret, true
	}

	return strings.Trim(noteContent, " "), false
}

func listBooks(ctx infra.DnoteCtx) error {
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
		log.Printf("%s %s\n", info.BookLabel, log.SprintfYellow("(%d)", info.NoteCount))
	}

	return nil
}

func listNotes(ctx infra.DnoteCtx, bookName string) error {
	db := ctx.DB

	var bookUUID string
	err := db.QueryRow("SELECT uuid FROM books WHERE label = ?", bookName).Scan(&bookUUID)
	if err == sql.ErrNoRows {
		return errors.New("book not found")
	} else if err != nil {
		return errors.Wrap(err, "querying the book")
	}

	rows, err := db.Query(`SELECT id, content FROM notes WHERE book_uuid = ? ORDER BY added_on ASC;`, bookUUID)
	if err != nil {
		return errors.Wrap(err, "querying notes")
	}
	defer rows.Close()

	infos := []noteInfo{}
	for rows.Next() {
		var info noteInfo
		err = rows.Scan(&info.ID, &info.Content)
		if err != nil {
			return errors.Wrap(err, "scanning a row")
		}

		infos = append(infos, info)
	}

	log.Infof("on book %s\n", bookName)

	for _, info := range infos {
		content, isExcerpt := formatContent(info.Content)

		index := log.SprintfYellow("(%d)", info.ID)
		if isExcerpt {
			content = fmt.Sprintf("%s %s", content, log.SprintfYellow("[---More---]"))
		}

		log.Plainf("%s %s\n", index, content)
	}

	return nil
}
