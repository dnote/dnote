/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package find

import (
	"database/sql"
	"fmt"
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
		return errors.New("Incorrect number of argument")
	}

	return nil
}

// NewCmd returns a new remove command
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

// noteInfo is an information about the note to be printed on screen
type noteInfo struct {
	RowID     int
	BookLabel string
	Body      string
}

// formatFTSSnippet turns the matched snippet from a full text search
// into a format suitable for CLI output
func formatFTSSnippet(s string) (string, error) {
	// first, strip all new lines
	body := newLineReg.ReplaceAllString(s, " ")

	var format, buf strings.Builder
	var args []interface{}

	toks := tokenize(body)

	for _, tok := range toks {
		if tok.Kind == tokenKindHLBegin || tok.Kind == tokenKindEOL {
			format.WriteString("%s")
			args = append(args, buf.String())

			buf.Reset()
		} else if tok.Kind == tokenKindHLEnd {
			format.WriteString("%s")
			str := log.ColorYellow.Sprintf("%s", buf.String())
			args = append(args, str)

			buf.Reset()
		} else {
			if err := buf.WriteByte(tok.Value); err != nil {
				return "", errors.Wrap(err, "building string")
			}
		}
	}

	return fmt.Sprintf(format.String(), args...), nil
}

// escapePhrase escapes the user-supplied FTS keywords by wrapping each term around
// double quotations so that they are treated as 'strings' as defined by SQLite FTS5.
func escapePhrase(s string) (string, error) {
	var b strings.Builder

	terms := strings.Fields(s)

	for idx, term := range terms {
		if _, err := b.WriteString(fmt.Sprintf("\"%s\"", term)); err != nil {
			return "", errors.Wrap(err, "writing string to builder")
		}

		if idx != len(term)-1 {
			if err := b.WriteByte(' '); err != nil {
				return "", errors.Wrap(err, "writing space to builder")
			}
		}
	}

	return b.String(), nil
}

func doQuery(ctx context.DnoteCtx, query, bookName string) (*sql.Rows, error) {
	db := ctx.DB

	sql := `SELECT
		notes.rowid,
		books.label AS book_label,
		snippet(note_fts, 0, '<dnotehl>', '</dnotehl>', '...', 28)
	FROM note_fts
	INNER JOIN notes ON notes.rowid = note_fts.rowid
	INNER JOIN books ON notes.book_uuid = books.uuid
	WHERE note_fts MATCH ?`
	args := []interface{}{query}

	if bookName != "" {
		sql = fmt.Sprintf("%s AND books.label = ?", sql)
		args = append(args, bookName)
	}

	rows, err := db.Query(sql, args...)

	return rows, err
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		phrase, err := escapePhrase(args[0])
		if err != nil {
			return errors.Wrap(err, "escaping phrase")
		}

		rows, err := doQuery(ctx, phrase, bookName)
		if err != nil {
			return errors.Wrap(err, "querying notes")
		}
		defer rows.Close()

		infos := []noteInfo{}
		for rows.Next() {
			var info noteInfo

			var body string
			err = rows.Scan(&info.RowID, &info.BookLabel, &body)
			if err != nil {
				return errors.Wrap(err, "scanning a row")
			}

			body, err := formatFTSSnippet(body)
			if err != nil {
				return errors.Wrap(err, "formatting a body")
			}

			info.Body = body

			infos = append(infos, info)
		}

		for _, info := range infos {
			bookLabel := log.ColorYellow.Sprintf("(%s)", info.BookLabel)
			rowid := log.ColorYellow.Sprintf("(%d)", info.RowID)

			log.Plainf("%s %s %s\n", bookLabel, rowid, info.Body)
		}

		return nil
	}
}
