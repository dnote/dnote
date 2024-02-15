package find

import (
	"database/sql"
	"strings"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
    # find notes by a keyword
    dnote find rpoplpush

    # find notes by multiple keywords
    dnote find "building a heap"

    # find notes within a book
    dnote find "merge sort" -b algorithm
`

var bookName string

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Incorrect number of arguments")
	}

	return nil
}

// NewCmd returns a new find command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "find",
		Short:   "Find notes by keywords",
		Aliases: []string{"f"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&bookName, "book", "b", "", "book name to find notes in")

	return cmd
}

func escapePhrase(s string) string {
	escaped := strings.ReplaceAll(s, "%", "\\%")
	escaped = strings.ReplaceAll(escaped, "_", "\\_")

	return "%" + escaped + "%"
}

func doQuery(ctx context.DnoteCtx, query, bookName string) (*sql.Rows, error) {
	db := ctx.DB

	sqlQuery := `SELECT
            notes.rowid,
            books.label AS book_label,
            notes.body
        FROM notes
        JOIN books ON notes.book_uuid = books.uuid
        WHERE notes.body LIKE ?`
	args := []interface{}{query}

	if bookName != "" {
		sqlQuery += " AND books.label = ?"
		args = append(args, bookName)
	}

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func highlightMatchesInLines(text, searchPhrase string) []string {
	highlightedLines := []string{}
	lines := strings.Split(text, "\n")

	// Convertimos la frase de búsqueda a minúsculas para la comparación insensible a mayúsculas
	searchPhraseLower := strings.ToLower(searchPhrase)

	for _, line := range lines {
		lineLower := strings.ToLower(line)
		if strings.Contains(lineLower, searchPhraseLower) {
			// Encuentra todas las coincidencias de la frase de búsqueda en la línea actual, insensible a mayúsculas
			var startIndex int
			var highlightedLine strings.Builder

			for startIndex < len(line) {
				matchIndex := strings.Index(lineLower[startIndex:], searchPhraseLower)
				if matchIndex == -1 {
					// No hay más coincidencias, agregar el resto de la línea
					highlightedLine.WriteString(line[startIndex:])
					break
				}

				// Agregar texto antes de la coincidencia
				highlightedLine.WriteString(line[startIndex : startIndex+matchIndex])

				// Agregar coincidencia resaltada
				match := line[startIndex+matchIndex : startIndex+matchIndex+len(searchPhraseLower)]
				highlightedLine.WriteString(log.ColorRed.Sprintf("%s", match))

				// Actualizar el índice de inicio para buscar la próxima coincidencia
				startIndex += matchIndex + len(searchPhraseLower)
			}

			highlightedLines = append(highlightedLines, highlightedLine.String())
		}
	}

	return highlightedLines
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		phrase := escapePhrase(args[0])

		rows, err := doQuery(ctx, phrase, bookName)
		if err != nil {
			return errors.Wrap(err, "querying notes")
		}
		defer rows.Close()

		for rows.Next() {
			var rowID int
			var bookLabel, body string

			err = rows.Scan(&rowID, &bookLabel, &body)
			if err != nil {
				return errors.Wrap(err, "scanning a row")
			}

			highlightedLines := highlightMatchesInLines(body, args[0])

			bookLabelStr := log.ColorYellow.Sprintf("(%s)", bookLabel)
			rowIDStr := log.ColorYellow.Sprintf("(%d)", rowID)

			for _, line := range highlightedLines {
				log.Plainf("%s %s %s\n", bookLabelStr, rowIDStr, line)
			}
		}

		return nil
	}
}
