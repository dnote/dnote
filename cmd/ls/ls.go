package ls

import (
	"fmt"
	"sort"
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

var deprecationWarning = `and "view" will replace it in v0.5.0.

Run "dnote view --help" for more information.
`

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

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

func NewRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dnote, err := core.GetDnote(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to read dnote")
		}

		if len(args) == 0 {
			if err := printBooks(dnote); err != nil {
				return errors.Wrap(err, "Failed to print books")
			}

			return nil
		}

		bookName := args[0]
		if err := printNotes(dnote, bookName); err != nil {
			return errors.Wrapf(err, "Failed to print notes for the book %s", bookName)
		}

		return nil
	}
}

// bookInfo is an information about the book to be printed on screen
type bookInfo struct {
	BookName  string
	NoteCount int
}

func getBookInfos(dnote infra.Dnote) []bookInfo {
	var ret []bookInfo

	for bookName, book := range dnote {
		ret = append(ret, bookInfo{BookName: bookName, NoteCount: len(book.Notes)})
	}

	return ret
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

func printBooks(dnote infra.Dnote) error {
	infos := getBookInfos(dnote)

	// Show books with more notes first
	sort.SliceStable(infos, func(i, j int) bool {
		return infos[i].NoteCount > infos[j].NoteCount
	})

	for _, info := range infos {
		log.Printf("%s %s\n", info.BookName, log.SprintfYellow("(%d)", info.NoteCount))
	}

	return nil
}

func printNotes(dnote infra.Dnote, bookName string) error {
	log.Infof("on book %s\n", bookName)

	book := dnote[bookName]

	for i, note := range book.Notes {
		content, isExcerpt := formatContent(note.Content)

		index := log.SprintfYellow("(%d)", i)
		if isExcerpt {
			content = fmt.Sprintf("%s %s", content, log.SprintfYellow("[---More---]"))
		}

		log.Plainf("%s %s\n", index, content)
	}

	return nil
}
