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

package migrate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/dnote/actions"
	"github.com/dnote/dnote/pkg/cli/client"
	"github.com/dnote/dnote/pkg/cli/config"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/pkg/errors"
)

type migration struct {
	name string
	run  func(ctx context.DnoteCtx, tx *database.DB) error
}

var lm1 = migration{
	name: "upgrade-edit-note-from-v1-to-v3",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		rows, err := tx.Query("SELECT uuid, data FROM actions WHERE type = ? AND schema = ?", "edit_note", 1)
		if err != nil {
			return errors.Wrap(err, "querying rows")
		}
		defer rows.Close()

		f := false

		for rows.Next() {
			var uuid, dat string

			err = rows.Scan(&uuid, &dat)
			if err != nil {
				return errors.Wrap(err, "scanning a row")
			}

			var oldData actions.EditNoteDataV1
			err = json.Unmarshal([]byte(dat), &oldData)
			if err != nil {
				return errors.Wrap(err, "unmarshalling existing data")
			}

			newData := actions.EditNoteDataV3{
				NoteUUID: oldData.NoteUUID,
				Content:  &oldData.Content,
				// With edit_note v1, CLI did not support changing books or public
				BookName: nil,
				Public:   &f,
			}

			b, err := json.Marshal(newData)
			if err != nil {
				return errors.Wrap(err, "marshalling new data")
			}

			_, err = tx.Exec("UPDATE actions SET data = ?, schema = ? WHERE uuid = ?", string(b), 3, uuid)
			if err != nil {
				return errors.Wrap(err, "updating a row")
			}
		}

		return nil
	},
}

var lm2 = migration{
	name: "upgrade-edit-note-from-v2-to-v3",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		rows, err := tx.Query("SELECT uuid, data FROM actions WHERE type = ? AND schema = ?", "edit_note", 2)
		if err != nil {
			return errors.Wrap(err, "querying rows")
		}
		defer rows.Close()

		for rows.Next() {
			var uuid, dat string

			err = rows.Scan(&uuid, &dat)
			if err != nil {
				return errors.Wrap(err, "scanning a row")
			}

			var oldData actions.EditNoteDataV2
			err = json.Unmarshal([]byte(dat), &oldData)
			if err != nil {
				return errors.Wrap(err, "unmarshalling existing data")
			}

			newData := actions.EditNoteDataV3{
				NoteUUID: oldData.NoteUUID,
				BookName: oldData.ToBook,
				Content:  oldData.Content,
				Public:   oldData.Public,
			}

			b, err := json.Marshal(newData)
			if err != nil {
				return errors.Wrap(err, "marshalling new data")
			}

			_, err = tx.Exec("UPDATE actions SET data = ?, schema = ? WHERE uuid = ?", string(b), 3, uuid)
			if err != nil {
				return errors.Wrap(err, "updating a row")
			}
		}

		return nil
	},
}

var lm3 = migration{
	name: "upgrade-remove-note-from-v1-to-v2",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		rows, err := tx.Query("SELECT uuid, data FROM actions WHERE type = ? AND schema = ?", "remove_note", 1)
		if err != nil {
			return errors.Wrap(err, "querying rows")
		}
		defer rows.Close()

		for rows.Next() {
			var uuid, dat string

			err = rows.Scan(&uuid, &dat)
			if err != nil {
				return errors.Wrap(err, "scanning a row")
			}

			var oldData actions.RemoveNoteDataV1
			err = json.Unmarshal([]byte(dat), &oldData)
			if err != nil {
				return errors.Wrap(err, "unmarshalling existing data")
			}

			newData := actions.RemoveNoteDataV2{
				NoteUUID: oldData.NoteUUID,
			}

			b, err := json.Marshal(newData)
			if err != nil {
				return errors.Wrap(err, "marshalling new data")
			}

			_, err = tx.Exec("UPDATE actions SET data = ?, schema = ? WHERE uuid = ?", string(b), 2, uuid)
			if err != nil {
				return errors.Wrap(err, "updating a row")
			}
		}

		return nil
	},
}

var lm4 = migration{
	name: "add-dirty-usn-deleted-to-notes-and-books",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		_, err := tx.Exec("ALTER TABLE books ADD COLUMN dirty bool DEFAULT false")
		if err != nil {
			return errors.Wrap(err, "adding dirty column to books")
		}

		_, err = tx.Exec("ALTER TABLE books ADD COLUMN usn int DEFAULT 0 NOT NULL")
		if err != nil {
			return errors.Wrap(err, "adding usn column to books")
		}

		_, err = tx.Exec("ALTER TABLE books ADD COLUMN deleted bool DEFAULT false")
		if err != nil {
			return errors.Wrap(err, "adding deleted column to books")
		}

		_, err = tx.Exec("ALTER TABLE notes ADD COLUMN dirty bool DEFAULT false")
		if err != nil {
			return errors.Wrap(err, "adding dirty column to notes")
		}

		_, err = tx.Exec("ALTER TABLE notes ADD COLUMN usn int DEFAULT 0 NOT NULL")
		if err != nil {
			return errors.Wrap(err, "adding usn column to notes")
		}

		_, err = tx.Exec("ALTER TABLE notes ADD COLUMN deleted bool DEFAULT false")
		if err != nil {
			return errors.Wrap(err, "adding deleted column to notes")
		}

		return nil
	},
}

var lm5 = migration{
	name: "mark-action-targets-dirty",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		rows, err := tx.Query("SELECT uuid, data, type FROM actions")
		if err != nil {
			return errors.Wrap(err, "querying rows")
		}
		defer rows.Close()

		for rows.Next() {
			var uuid, dat, actionType string

			err = rows.Scan(&uuid, &dat, &actionType)
			if err != nil {
				return errors.Wrap(err, "scanning a row")
			}

			// removed notes and removed books cannot be reliably derived retrospectively
			// because books did not use to have uuid. Users will find locally deleted
			// notes and books coming back to existence if they have not synced the change.
			// But there will be no data loss.
			switch actionType {
			case "add_note":
				var data actions.AddNoteDataV2
				err = json.Unmarshal([]byte(dat), &data)
				if err != nil {
					return errors.Wrap(err, "unmarshalling existing data")
				}

				_, err := tx.Exec("UPDATE notes SET dirty = true WHERE uuid = ?", data.NoteUUID)
				if err != nil {
					return errors.Wrapf(err, "markig note dirty '%s'", data.NoteUUID)
				}
			case "edit_note":
				var data actions.EditNoteDataV3
				err = json.Unmarshal([]byte(dat), &data)
				if err != nil {
					return errors.Wrap(err, "unmarshalling existing data")
				}

				_, err := tx.Exec("UPDATE notes SET dirty = true WHERE uuid = ?", data.NoteUUID)
				if err != nil {
					return errors.Wrapf(err, "markig note dirty '%s'", data.NoteUUID)
				}
			case "add_book":
				var data actions.AddBookDataV1
				err = json.Unmarshal([]byte(dat), &data)
				if err != nil {
					return errors.Wrap(err, "unmarshalling existing data")
				}

				_, err := tx.Exec("UPDATE books SET dirty = true WHERE label = ?", data.BookName)
				if err != nil {
					return errors.Wrapf(err, "markig note dirty '%s'", data.BookName)
				}
			}
		}

		return nil
	},
}

var lm6 = migration{
	name: "drop-actions",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		_, err := tx.Exec("DROP TABLE actions;")
		if err != nil {
			return errors.Wrap(err, "dropping the actions table")
		}

		return nil
	},
}

var lm7 = migration{
	name: "resolve-conflicts-with-reserved-book-names",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		migrateBook := func(name string) error {
			var uuid string

			err := tx.QueryRow("SELECT uuid FROM books WHERE label = ?", name).Scan(&uuid)
			if err == sql.ErrNoRows {
				// if not found, noop
				return nil
			} else if err != nil {
				return errors.Wrap(err, "finding trash book")
			}

			for i := 2; ; i++ {
				candidate := fmt.Sprintf("%s (%d)", name, i)

				var count int
				err := tx.QueryRow("SELECT count(*) FROM books WHERE label = ?", candidate).Scan(&count)
				if err != nil {
					return errors.Wrap(err, "counting candidate")
				}

				if count == 0 {
					_, err := tx.Exec("UPDATE books SET label = ?, dirty = ? WHERE uuid = ?", candidate, true, uuid)
					if err != nil {
						return errors.Wrapf(err, "updating book '%s'", name)
					}

					break
				}
			}

			return nil
		}

		if err := migrateBook("trash"); err != nil {
			return errors.Wrap(err, "migrating trash book")
		}
		if err := migrateBook("conflicts"); err != nil {
			return errors.Wrap(err, "migrating conflicts book")
		}

		return nil
	},
}

var lm8 = migration{
	name: "drop-note-id-and-rename-content-to-body",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		_, err := tx.Exec(`CREATE TABLE notes_tmp
		(
			uuid text NOT NULL,
			book_uuid text NOT NULL,
			body text NOT NULL,
			added_on integer NOT NULL,
			edited_on integer DEFAULT 0,
			public bool DEFAULT false,
			dirty bool DEFAULT false,
			usn int DEFAULT 0 NOT NULL,
			deleted bool DEFAULT false
		);`)
		if err != nil {
			return errors.Wrap(err, "creating temporary notes table for migration")
		}

		_, err = tx.Exec(`INSERT INTO notes_tmp
			SELECT uuid, book_uuid, content, added_on, edited_on, public, dirty, usn, deleted FROM notes;`)
		if err != nil {
			return errors.Wrap(err, "copying data to new table")
		}

		_, err = tx.Exec(`DROP TABLE notes;`)
		if err != nil {
			return errors.Wrap(err, "dropping the notes table")
		}

		_, err = tx.Exec(`ALTER TABLE notes_tmp RENAME to notes;`)
		if err != nil {
			return errors.Wrap(err, "renaming the temporary notes table")
		}

		return nil
	},
}

var lm9 = migration{
	name: "create-fts-index",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		_, err := tx.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS note_fts USING fts5(content=notes, body, tokenize="porter unicode61 categories 'L* N* Co Ps Pe'");`)
		if err != nil {
			return errors.Wrap(err, "creating note_fts")
		}

		// Create triggers to keep the indices in note_fts in sync with notes
		_, err = tx.Exec(`
			CREATE TRIGGER notes_after_insert AFTER INSERT ON notes BEGIN
				INSERT INTO note_fts(rowid, body) VALUES (new.rowid, new.body);
			END;
			CREATE TRIGGER notes_after_delete AFTER DELETE ON notes BEGIN
				INSERT INTO note_fts(note_fts, rowid, body) VALUES ('delete', old.rowid, old.body);
			END;
			CREATE TRIGGER notes_after_update AFTER UPDATE ON notes BEGIN
				INSERT INTO note_fts(note_fts, rowid, body) VALUES ('delete', old.rowid, old.body);
				INSERT INTO note_fts(rowid, body) VALUES (new.rowid, new.body);
			END;
		`)
		if err != nil {
			return errors.Wrap(err, "creating triggers for note_fts")
		}

		// populate fts indices
		_, err = tx.Exec(`INSERT INTO note_fts (rowid, body)
			SELECT rowid, body FROM notes;`)
		if err != nil {
			return errors.Wrap(err, "populating note_fts")
		}

		return nil
	},
}

var lm10 = migration{
	name: "rename-number-only-book",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		migrateBook := func(label string) error {
			var uuid string

			err := tx.QueryRow("SELECT uuid FROM books WHERE label = ?", label).Scan(&uuid)
			if err != nil {
				return errors.Wrap(err, "finding uuid")
			}

			for i := 1; ; i++ {
				candidate := fmt.Sprintf("%s (%d)", label, i)

				var count int
				err := tx.QueryRow("SELECT count(*) FROM books WHERE label = ?", candidate).Scan(&count)
				if err != nil {
					return errors.Wrap(err, "counting candidate")
				}

				if count == 0 {
					_, err := tx.Exec("UPDATE books SET label = ?, dirty = ? WHERE uuid = ?", candidate, true, uuid)
					if err != nil {
						return errors.Wrapf(err, "updating book '%s'", label)
					}

					break
				}
			}

			return nil
		}

		rows, err := tx.Query("SELECT label FROM books")
		defer rows.Close()
		if err != nil {
			return errors.Wrap(err, "getting labels")
		}

		var regexNumber = regexp.MustCompile(`^\d+$`)

		for rows.Next() {
			var label string
			err := rows.Scan(&label)
			if err != nil {
				return errors.Wrap(err, "scannign row")
			}

			if regexNumber.MatchString(label) {
				err = migrateBook(label)
				if err != nil {
					return errors.Wrapf(err, "migrating book %s", label)
				}
			}
		}

		return nil
	},
}

var lm11 = migration{
	name: "rename-book-labels-with-space",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		processLabel := func(label string) (string, error) {
			sanitized := strings.Replace(label, " ", "_", -1)

			var cnt int
			err := tx.QueryRow("SELECT count(*) FROM books WHERE label = ?", sanitized).Scan(&cnt)
			if err != nil {
				return "", errors.Wrap(err, "counting ret")
			}
			if cnt == 0 {
				return sanitized, nil
			}

			// if there is a collision, resolve by appending number
			var ret string
			for i := 2; ; i++ {
				ret = fmt.Sprintf("%s_%d", sanitized, i)

				var count int
				err := tx.QueryRow("SELECT count(*) FROM books WHERE label = ?", ret).Scan(&count)
				if err != nil {
					return "", errors.Wrap(err, "counting ret")
				}

				if count == 0 {
					break
				}
			}

			return ret, nil
		}

		rows, err := tx.Query("SELECT uuid, label FROM books")
		defer rows.Close()
		if err != nil {
			return errors.Wrap(err, "getting labels")
		}

		for rows.Next() {
			var uuid, label string
			err := rows.Scan(&uuid, &label)
			if err != nil {
				return errors.Wrap(err, "scanning row")
			}

			if strings.Contains(label, " ") {
				processed, err := processLabel(label)
				if err != nil {
					return errors.Wrapf(err, "resolving book name for %s", label)
				}

				_, err = tx.Exec("UPDATE books SET label = ?, dirty = ? WHERE uuid = ?", processed, true, uuid)
				if err != nil {
					return errors.Wrapf(err, "updating book '%s'", label)
				}
			}
		}

		return nil
	},
}

var lm12 = migration{
	name: "add apiEndpoint to the configuration file",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		cf, err := config.Read(ctx)
		if err != nil {
			return errors.Wrap(err, "reading config")
		}

		cf.APIEndpoint = "https://api.getdnote.com"

		err = config.Write(ctx, cf)
		if err != nil {
			return errors.Wrap(err, "writing config")
		}

		return nil
	},
}

var rm1 = migration{
	name: "sync-book-uuids-from-server",
	run: func(ctx context.DnoteCtx, tx *database.DB) error {
		sessionKey := ctx.SessionKey
		if sessionKey == "" {
			return errors.New("not logged in")
		}

		resp, err := client.GetBooks(ctx, sessionKey)
		if err != nil {
			return errors.Wrap(err, "getting books from the server")
		}
		log.Debug("book details from the server: %+v\n", resp)

		UUIDMap := map[string]string{}
		for _, book := range resp {
			// Build a map from uuid to label
			UUIDMap[book.Label] = book.UUID
		}

		for _, book := range resp {
			// update uuid in the books table
			log.Debug("Updating book %s\n", book.Label)

			//todo if does not exist, then continue loop
			var count int
			if err := tx.
				QueryRow("SELECT count(*) FROM books WHERE label = ?", book.Label).
				Scan(&count); err != nil {
				return errors.Wrapf(err, "checking if book exists: %s", book.Label)
			}

			if count == 0 {
				continue
			}

			var originalUUID string
			if err := tx.
				QueryRow("SELECT uuid FROM books WHERE label = ?", book.Label).
				Scan(&originalUUID); err != nil {
				return errors.Wrapf(err, "scanning the orignal uuid of the book %s", book.Label)
			}
			log.Debug("original uuid: %s. new_uuid %s\n", originalUUID, book.UUID)

			_, err := tx.Exec("UPDATE books SET uuid = ? WHERE label = ?", book.UUID, book.Label)
			if err != nil {
				return errors.Wrapf(err, "updating book '%s'", book.Label)
			}

			_, err = tx.Exec("UPDATE notes SET book_uuid = ? WHERE book_uuid = ?", book.UUID, originalUUID)
			if err != nil {
				return errors.Wrapf(err, "updating book_uuids of notes")
			}
		}

		return nil
	},
}
