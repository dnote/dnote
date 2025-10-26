/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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

package sync

import (
	"database/sql"
	"fmt"

	"github.com/dnote/dnote/pkg/cli/client"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/migrate"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/dnote/dnote/pkg/cli/upgrade"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	modeInsert = iota
	modeUpdate
)

var example = `
  dnote sync`

var isFullSync bool
var apiEndpointFlag string

// NewCmd returns a new sync command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Aliases: []string{"s"},
		Short:   "Sync data with the server",
		Example: example,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.BoolVarP(&isFullSync, "full", "f", false, "perform a full sync instead of incrementally syncing only the changed data.")
	f.StringVar(&apiEndpointFlag, "apiEndpoint", "", "API endpoint to connect to (defaults to value in config)")

	return cmd
}

func getLastSyncAt(tx *database.DB) (int, error) {
	var ret int

	if err := database.GetSystem(tx, consts.SystemLastSyncAt, &ret); err != nil {
		return ret, errors.Wrap(err, "querying last sync time")
	}

	return ret, nil
}

func getLastMaxUSN(tx *database.DB) (int, error) {
	var ret int

	if err := database.GetSystem(tx, consts.SystemLastMaxUSN, &ret); err != nil {
		return ret, errors.Wrap(err, "querying last user max_usn")
	}

	return ret, nil
}

// syncList is an aggregation of resources represented in the sync fragments
type syncList struct {
	Notes          map[string]client.SyncFragNote
	Books          map[string]client.SyncFragBook
	ExpungedNotes  map[string]bool
	ExpungedBooks  map[string]bool
	MaxUSN         int
	UserMaxUSN     int // Server's actual max USN (for distinguishing empty fragment vs empty server)
	MaxCurrentTime int64
}

func (l syncList) getLength() int {
	return len(l.Notes) + len(l.Books) + len(l.ExpungedNotes) + len(l.ExpungedBooks)
}

// processFragments categorizes items in sync fragments into a sync list.
func processFragments(fragments []client.SyncFragment) (syncList, error) {
	notes := map[string]client.SyncFragNote{}
	books := map[string]client.SyncFragBook{}
	expungedNotes := map[string]bool{}
	expungedBooks := map[string]bool{}
	var maxUSN int
	var userMaxUSN int
	var maxCurrentTime int64

	for _, fragment := range fragments {
		for _, note := range fragment.Notes {
			notes[note.UUID] = note
		}
		for _, book := range fragment.Books {
			books[book.UUID] = book
		}
		for _, uuid := range fragment.ExpungedBooks {
			expungedBooks[uuid] = true
		}
		for _, uuid := range fragment.ExpungedNotes {
			expungedNotes[uuid] = true
		}

		if fragment.FragMaxUSN > maxUSN {
			maxUSN = fragment.FragMaxUSN
		}
		if fragment.UserMaxUSN > userMaxUSN {
			userMaxUSN = fragment.UserMaxUSN
		}
		if fragment.CurrentTime > maxCurrentTime {
			maxCurrentTime = fragment.CurrentTime
		}
	}

	sl := syncList{
		Notes:          notes,
		Books:          books,
		ExpungedNotes:  expungedNotes,
		ExpungedBooks:  expungedBooks,
		MaxUSN:         maxUSN,
		UserMaxUSN:     userMaxUSN,
		MaxCurrentTime: maxCurrentTime,
	}

	return sl, nil
}

// getSyncList gets a list of all sync fragments after the specified usn
// and aggregates them into a syncList data structure
func getSyncList(ctx context.DnoteCtx, afterUSN int) (syncList, error) {
	fragments, err := getSyncFragments(ctx, afterUSN)
	if err != nil {
		return syncList{}, errors.Wrap(err, "getting sync fragments")
	}

	ret, err := processFragments(fragments)
	if err != nil {
		return syncList{}, errors.Wrap(err, "making sync list")
	}

	return ret, nil
}

// getSyncFragments repeatedly gets all sync fragments after the specified usn until there is no more new data
// remaining and returns the buffered list
func getSyncFragments(ctx context.DnoteCtx, afterUSN int) ([]client.SyncFragment, error) {
	var buf []client.SyncFragment

	nextAfterUSN := afterUSN

	for {
		resp, err := client.GetSyncFragment(ctx, nextAfterUSN)
		if err != nil {
			return buf, errors.Wrap(err, "getting sync fragment")
		}

		frag := resp.Fragment
		buf = append(buf, frag)

		nextAfterUSN = frag.FragMaxUSN

		// if there is no more data, break
		if nextAfterUSN == 0 {
			break
		}
	}

	log.Debug("received sync fragments: %+v\n", redactSyncFragments(buf))

	return buf, nil
}

// redactSyncFragments returns a deep copy of sync fragments with sensitive fields (note body, book label) removed for safe logging
func redactSyncFragments(fragments []client.SyncFragment) []client.SyncFragment {
	redacted := make([]client.SyncFragment, len(fragments))
	for i, frag := range fragments {
		// Create new notes with redacted bodies
		notes := make([]client.SyncFragNote, len(frag.Notes))
		for j, note := range frag.Notes {
			notes[j] = client.SyncFragNote{
				UUID:      note.UUID,
				BookUUID:  note.BookUUID,
				USN:       note.USN,
				CreatedAt: note.CreatedAt,
				UpdatedAt: note.UpdatedAt,
				AddedOn:   note.AddedOn,
				EditedOn:  note.EditedOn,
				Body: func() string {
					if note.Body != "" {
						return "<redacted>"
					}
					return ""
				}(),
				Deleted: note.Deleted,
			}
		}

		// Create new books with redacted labels
		books := make([]client.SyncFragBook, len(frag.Books))
		for j, book := range frag.Books {
			books[j] = client.SyncFragBook{
				UUID:      book.UUID,
				USN:       book.USN,
				CreatedAt: book.CreatedAt,
				UpdatedAt: book.UpdatedAt,
				AddedOn:   book.AddedOn,
				Label: func() string {
					if book.Label != "" {
						return "<redacted>"
					}
					return ""
				}(),
				Deleted: book.Deleted,
			}
		}

		redacted[i] = client.SyncFragment{
			FragMaxUSN:    frag.FragMaxUSN,
			UserMaxUSN:    frag.UserMaxUSN,
			CurrentTime:   frag.CurrentTime,
			Notes:         notes,
			Books:         books,
			ExpungedNotes: frag.ExpungedNotes,
			ExpungedBooks: frag.ExpungedBooks,
		}
	}
	return redacted
}

// resolveLabel resolves a book label conflict by repeatedly appending an increasing integer
// to the label until it finds a unique label. It returns the first non-conflicting label.
func resolveLabel(tx *database.DB, label string) (string, error) {
	var ret string

	for i := 2; ; i++ {
		ret = fmt.Sprintf("%s_%d", label, i)

		var cnt int
		if err := tx.QueryRow("SELECT count(*) FROM books WHERE label = ?", ret).Scan(&cnt); err != nil {
			return "", errors.Wrapf(err, "checking availability of label %s", ret)
		}

		if cnt == 0 {
			break
		}
	}

	return ret, nil
}

// mergeBook inserts or updates the given book in the local database.
// If a book with a duplicate label exists locally, it renames the duplicate by appending a number.
func mergeBook(tx *database.DB, b client.SyncFragBook, mode int) error {
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM books WHERE label = ?", b.Label).Scan(&count); err != nil {
		return errors.Wrapf(err, "checking for books with a duplicate label %s", b.Label)
	}

	// if duplicate exists locally, rename it and mark it dirty
	if count > 0 {
		newLabel, err := resolveLabel(tx, b.Label)
		if err != nil {
			return errors.Wrap(err, "getting a new book label for conflict resolution")
		}

		if _, err := tx.Exec("UPDATE books SET label = ?, dirty = ? WHERE label = ? AND uuid != ?", newLabel, true, b.Label, b.UUID); err != nil {
			return errors.Wrap(err, "resolving duplicate book label")
		}
	}

	if mode == modeInsert {
		book := database.NewBook(b.UUID, b.Label, b.USN, false, false)
		if err := book.Insert(tx); err != nil {
			return errors.Wrapf(err, "inserting note with uuid %s", b.UUID)
		}
	} else if mode == modeUpdate {
		// The state from the server overwrites the local state. In other words, the server change always wins.
		if _, err := tx.Exec("UPDATE books SET usn = ?, uuid = ?, label = ?, deleted = ? WHERE uuid = ?",
			b.USN, b.UUID, b.Label, b.Deleted, b.UUID); err != nil {
			return errors.Wrapf(err, "updating local book %s", b.UUID)
		}
	}

	return nil
}

func stepSyncBook(tx *database.DB, b client.SyncFragBook) error {
	var localUSN int
	var dirty bool
	err := tx.QueryRow("SELECT usn, dirty FROM books WHERE uuid = ?", b.UUID).Scan(&localUSN, &dirty)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrapf(err, "getting local book %s", b.UUID)
	}

	// if book exists in the server and does not exist in the client
	if err == sql.ErrNoRows {
		if e := mergeBook(tx, b, modeInsert); e != nil {
			return errors.Wrapf(e, "resolving book")
		}

		return nil
	}

	if e := mergeBook(tx, b, modeUpdate); e != nil {
		return errors.Wrapf(e, "resolving book")
	}

	return nil
}

func mergeNote(tx *database.DB, serverNote client.SyncFragNote, localNote database.Note) error {
	var bookDeleted bool
	err := tx.QueryRow("SELECT deleted FROM books WHERE uuid = ?", localNote.BookUUID).Scan(&bookDeleted)
	if err != nil {
		return errors.Wrapf(err, "checking if local book %s is deleted", localNote.BookUUID)
	}

	// if the book is deleted, noop
	if bookDeleted {
		return nil
	}

	// if the local copy is deleted, and it was edited on the server, override with server values and mark it not dirty.
	if localNote.Deleted {
		if _, err := tx.Exec("UPDATE notes SET usn = ?, book_uuid = ?, body = ?, edited_on = ?, deleted = ?, dirty = ? WHERE uuid = ?",
			serverNote.USN, serverNote.BookUUID, serverNote.Body, serverNote.EditedOn, serverNote.Deleted, false, serverNote.UUID); err != nil {
			return errors.Wrapf(err, "updating local note %s", serverNote.UUID)
		}

		return nil
	}

	mr, err := mergeNoteFields(tx, localNote, serverNote)
	if err != nil {
		return errors.Wrapf(err, "reporting note conflict for note %s", localNote.UUID)
	}

	if _, err := tx.Exec("UPDATE notes SET usn = ?, book_uuid = ?, body = ?, edited_on = ?, deleted = ?  WHERE uuid = ?",
		serverNote.USN, mr.bookUUID, mr.body, mr.editedOn, serverNote.Deleted, serverNote.UUID); err != nil {
		return errors.Wrapf(err, "updating local note %s", serverNote.UUID)
	}

	return nil
}

func stepSyncNote(tx *database.DB, n client.SyncFragNote) error {
	var localNote database.Note
	err := tx.QueryRow("SELECT body, usn, book_uuid, dirty, deleted FROM notes WHERE uuid = ?", n.UUID).
		Scan(&localNote.Body, &localNote.USN, &localNote.BookUUID, &localNote.Dirty, &localNote.Deleted)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrapf(err, "getting local note %s", n.UUID)
	}

	// if note exists in the server and does not exist in the client, insert the note.
	if err == sql.ErrNoRows {
		note := database.NewNote(n.UUID, n.BookUUID, n.Body, n.AddedOn, n.EditedOn, n.USN, n.Deleted, false)

		if err := note.Insert(tx); err != nil {
			return errors.Wrapf(err, "inserting note with uuid %s", n.UUID)
		}
	} else {
		if err := mergeNote(tx, n, localNote); err != nil {
			return errors.Wrap(err, "merging local note")
		}
	}

	return nil
}

func fullSyncNote(tx *database.DB, n client.SyncFragNote) error {
	var localNote database.Note
	err := tx.QueryRow("SELECT body, usn, book_uuid, dirty, deleted FROM notes WHERE uuid = ?", n.UUID).
		Scan(&localNote.Body, &localNote.USN, &localNote.BookUUID, &localNote.Dirty, &localNote.Deleted)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrapf(err, "getting local note %s", n.UUID)
	}

	// if note exists in the server and does not exist in the client, insert the note.
	if err == sql.ErrNoRows {
		note := database.NewNote(n.UUID, n.BookUUID, n.Body, n.AddedOn, n.EditedOn, n.USN, n.Deleted, false)

		if err := note.Insert(tx); err != nil {
			return errors.Wrapf(err, "inserting note with uuid %s", n.UUID)
		}
	} else if n.USN > localNote.USN {
		if err := mergeNote(tx, n, localNote); err != nil {
			return errors.Wrap(err, "merging local note")
		}
	}

	return nil
}

func syncDeleteNote(tx *database.DB, noteUUID string) error {
	var localUSN int
	var dirty bool
	err := tx.QueryRow("SELECT usn, dirty FROM notes WHERE uuid = ?", noteUUID).Scan(&localUSN, &dirty)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrapf(err, "getting local note %s", noteUUID)
	}

	// if note does not exist on client, noop
	if err == sql.ErrNoRows {
		return nil
	}

	// if local copy is not dirty, delete
	if !dirty {
		_, err = tx.Exec("DELETE FROM notes WHERE uuid = ?", noteUUID)
		if err != nil {
			return errors.Wrapf(err, "deleting local note %s", noteUUID)
		}
	}

	return nil
}

// checkNotesPristine checks that none of the notes in the given book are dirty
func checkNotesPristine(tx *database.DB, bookUUID string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ? AND dirty = ?", bookUUID, true).Scan(&count); err != nil {
		return false, errors.Wrapf(err, "counting notes that are dirty in book %s", bookUUID)
	}

	if count > 0 {
		return false, nil
	}

	return true, nil
}

func syncDeleteBook(tx *database.DB, bookUUID string) error {
	var localUSN int
	var dirty bool
	err := tx.QueryRow("SELECT usn, dirty FROM books WHERE uuid = ?", bookUUID).Scan(&localUSN, &dirty)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrapf(err, "getting local book %s", bookUUID)
	}

	// if book does not exist on client, noop
	if err == sql.ErrNoRows {
		return nil
	}

	// if local copy is dirty, noop. it will be uploaded to the server later
	if dirty {
		return nil
	}

	ok, err := checkNotesPristine(tx, bookUUID)
	if err != nil {
		return errors.Wrap(err, "checking if any notes are dirty in book")
	}
	// if the local book is not pristine, do not delete but mark it as dirty
	// so that it can be uploaded to the server later and become un-deleted
	if !ok {
		_, err = tx.Exec("UPDATE books SET dirty = ? WHERE uuid = ?", true, bookUUID)
		if err != nil {
			return errors.Wrapf(err, "marking a book dirty with uuid %s", bookUUID)
		}

		return nil
	}

	_, err = tx.Exec("DELETE FROM notes WHERE book_uuid = ?", bookUUID)
	if err != nil {
		return errors.Wrapf(err, "deleting local notes of the book %s", bookUUID)
	}

	_, err = tx.Exec("DELETE FROM books WHERE uuid = ?", bookUUID)
	if err != nil {
		return errors.Wrapf(err, "deleting local book %s", bookUUID)
	}

	return nil
}

func fullSyncBook(tx *database.DB, b client.SyncFragBook) error {
	var localUSN int
	var dirty bool
	err := tx.QueryRow("SELECT usn, dirty FROM books WHERE uuid = ?", b.UUID).Scan(&localUSN, &dirty)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrapf(err, "getting local book %s", b.UUID)
	}

	// if book exists in the server and does not exist in the client
	if err == sql.ErrNoRows {
		if e := mergeBook(tx, b, modeInsert); e != nil {
			return errors.Wrapf(e, "resolving book")
		}
	} else if b.USN > localUSN {
		if e := mergeBook(tx, b, modeUpdate); e != nil {
			return errors.Wrapf(e, "resolving book")
		}
	}

	return nil
}

// checkNoteInList checks if the given syncList contains the note with the given uuid
func checkNoteInList(uuid string, list *syncList) bool {
	if _, ok := list.Notes[uuid]; ok {
		return true
	}

	if _, ok := list.ExpungedNotes[uuid]; ok {
		return true
	}

	return false
}

// checkBookInList checks if the given syncList contains the book with the given uuid
func checkBookInList(uuid string, list *syncList) bool {
	if _, ok := list.Books[uuid]; ok {
		return true
	}

	if _, ok := list.ExpungedBooks[uuid]; ok {
		return true
	}

	return false
}

// cleanLocalNotes deletes from the local database any notes that are in invalid state
// judging by the full list of resources in the server. Concretely, the only acceptable
// situation in which a local note is not present in the server is if it is new and has not been
// uploaded (i.e. dirty and usn is 0). Otherwise, it is a result of some kind of error and should be cleaned.
func cleanLocalNotes(tx *database.DB, fullList *syncList) error {
	rows, err := tx.Query("SELECT uuid, usn, dirty FROM notes")
	if err != nil {
		return errors.Wrap(err, "getting local notes")
	}
	defer rows.Close()

	for rows.Next() {
		var note database.Note
		if err := rows.Scan(&note.UUID, &note.USN, &note.Dirty); err != nil {
			return errors.Wrap(err, "scanning a row for local note")
		}

		ok := checkNoteInList(note.UUID, fullList)
		if !ok && (!note.Dirty || note.USN != 0) {
			err = note.Expunge(tx)
			if err != nil {
				return errors.Wrap(err, "expunging a note")
			}
		}
	}

	return nil
}

// cleanLocalBooks deletes from the local database any books that are in invalid state
func cleanLocalBooks(tx *database.DB, fullList *syncList) error {
	rows, err := tx.Query("SELECT uuid, usn, dirty FROM books")
	if err != nil {
		return errors.Wrap(err, "getting local books")
	}
	defer rows.Close()

	for rows.Next() {
		var book database.Book
		if err := rows.Scan(&book.UUID, &book.USN, &book.Dirty); err != nil {
			return errors.Wrap(err, "scanning a row for local book")
		}

		ok := checkBookInList(book.UUID, fullList)
		if !ok && (!book.Dirty || book.USN != 0) {
			err = book.Expunge(tx)
			if err != nil {
				return errors.Wrap(err, "expunging a book")
			}
		}
	}

	return nil
}

func fullSync(ctx context.DnoteCtx, tx *database.DB) error {
	log.Debug("performing a full sync\n")
	log.Info("resolving delta.")

	log.DebugNewline()

	list, err := getSyncList(ctx, 0)
	if err != nil {
		return errors.Wrap(err, "getting sync list")
	}

	fmt.Printf(" (total %d).", list.getLength())

	log.DebugNewline()

	// clean resources that are in erroneous states
	if err := cleanLocalNotes(tx, &list); err != nil {
		return errors.Wrap(err, "cleaning up local notes")
	}
	if err := cleanLocalBooks(tx, &list); err != nil {
		return errors.Wrap(err, "cleaning up local books")
	}

	for _, note := range list.Notes {
		if err := fullSyncNote(tx, note); err != nil {
			return errors.Wrap(err, "merging note")
		}
	}
	for _, book := range list.Books {
		if err := fullSyncBook(tx, book); err != nil {
			return errors.Wrap(err, "merging book")
		}
	}

	for noteUUID := range list.ExpungedNotes {
		if err := syncDeleteNote(tx, noteUUID); err != nil {
			return errors.Wrap(err, "deleting note")
		}
	}
	for bookUUID := range list.ExpungedBooks {
		if err := syncDeleteBook(tx, bookUUID); err != nil {
			return errors.Wrap(err, "deleting book")
		}
	}

	err = saveSyncState(tx, list.MaxCurrentTime, list.MaxUSN, list.UserMaxUSN)
	if err != nil {
		return errors.Wrap(err, "saving sync state")
	}

	fmt.Println(" done.")

	return nil
}

func stepSync(ctx context.DnoteCtx, tx *database.DB, afterUSN int) error {
	log.Debug("performing a step sync\n")

	log.Info("resolving delta.")

	log.DebugNewline()

	list, err := getSyncList(ctx, afterUSN)
	if err != nil {
		return errors.Wrap(err, "getting sync list")
	}

	fmt.Printf(" (total %d).", list.getLength())

	for _, note := range list.Notes {
		if err := stepSyncNote(tx, note); err != nil {
			return errors.Wrap(err, "merging note")
		}
	}
	for _, book := range list.Books {
		if err := stepSyncBook(tx, book); err != nil {
			return errors.Wrap(err, "merging book")
		}
	}

	for noteUUID := range list.ExpungedNotes {
		if err := syncDeleteNote(tx, noteUUID); err != nil {
			return errors.Wrap(err, "deleting note")
		}
	}
	for bookUUID := range list.ExpungedBooks {
		if err := syncDeleteBook(tx, bookUUID); err != nil {
			return errors.Wrap(err, "deleting book")
		}
	}

	err = saveSyncState(tx, list.MaxCurrentTime, list.MaxUSN, list.UserMaxUSN)
	if err != nil {
		return errors.Wrap(err, "saving sync state")
	}

	fmt.Println(" done.")

	return nil
}

// isConflictError checks if an error is a 409 Conflict error from the server
func isConflictError(err error) bool {
	if err == nil {
		return false
	}

	var httpErr *client.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.IsConflict()
	}

	return false
}

func sendBooks(ctx context.DnoteCtx, tx *database.DB) (bool, error) {
	isBehind := false

	rows, err := tx.Query("SELECT uuid, label, usn, deleted FROM books WHERE dirty")
	if err != nil {
		return isBehind, errors.Wrap(err, "getting syncable books")
	}
	defer rows.Close()

	for rows.Next() {
		var book database.Book

		if err = rows.Scan(&book.UUID, &book.Label, &book.USN, &book.Deleted); err != nil {
			return isBehind, errors.Wrap(err, "scanning a syncable book")
		}

		log.Debug("sending book %s\n", book.UUID)

		var respUSN int

		// if new, create it in the server, or else, update.
		if book.USN == 0 {
			if book.Deleted {
				err = book.Expunge(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "expunging a book locally")
				}

				continue
			} else {
				resp, err := client.CreateBook(ctx, book.Label)
				if err != nil {
					log.Debug("error creating book (will retry after stepSync): %v\n", err)
					isBehind = true
					continue
				}

				_, err = tx.Exec("UPDATE notes SET book_uuid = ? WHERE book_uuid = ?", resp.Book.UUID, book.UUID)
				if err != nil {
					return isBehind, errors.Wrap(err, "updating book_uuids of notes")
				}

				book.Dirty = false
				book.USN = resp.Book.USN
				err = book.Update(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "marking book dirty")
				}

				err = book.UpdateUUID(tx, resp.Book.UUID)
				if err != nil {
					return isBehind, errors.Wrap(err, "updating book uuid")
				}

				respUSN = resp.Book.USN
			}
		} else {
			if book.Deleted {
				resp, err := client.DeleteBook(ctx, book.UUID)
				if err != nil {
					return isBehind, errors.Wrap(err, "deleting a book")
				}

				err = book.Expunge(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "expunging a book locally")
				}

				respUSN = resp.Book.USN
			} else {
				resp, err := client.UpdateBook(ctx, book.Label, book.UUID)
				if err != nil {
					return isBehind, errors.Wrap(err, "updating a book")
				}

				book.Dirty = false
				book.USN = resp.Book.USN
				err = book.Update(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "marking book dirty")
				}

				respUSN = resp.Book.USN
			}
		}

		lastMaxUSN, err := getLastMaxUSN(tx)
		if err != nil {
			return isBehind, errors.Wrap(err, "getting last max usn")
		}

		log.Debug("sent book %s. response USN %d. last max usn: %d\n", book.UUID, respUSN, lastMaxUSN)

		if respUSN == lastMaxUSN+1 {
			err = updateLastMaxUSN(tx, lastMaxUSN+1)
			if err != nil {
				return isBehind, errors.Wrap(err, "updating last max usn")
			}
		} else {
			isBehind = true
		}
	}

	return isBehind, nil
}

// findOrphanedNotes returns a list of all orphaned notes
func findOrphanedNotes(db *database.DB) (int, []struct{ noteUUID, bookUUID string }, error) {
	var orphanCount int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM notes n
		WHERE NOT EXISTS (
			SELECT 1 FROM books b
			WHERE b.uuid = n.book_uuid
			AND NOT b.deleted
		)
	`).Scan(&orphanCount)
	if err != nil {
		return 0, nil, err
	}

	if orphanCount == 0 {
		return 0, nil, nil
	}

	rows, err := db.Query(`
		SELECT n.uuid, n.book_uuid
		FROM notes n
		WHERE NOT EXISTS (
			SELECT 1 FROM books b
			WHERE b.uuid = n.book_uuid
			AND NOT b.deleted
		)
	`)
	if err != nil {
		return orphanCount, nil, err
	}
	defer rows.Close()

	var orphans []struct{ noteUUID, bookUUID string }
	for rows.Next() {
		var noteUUID, bookUUID string
		if err := rows.Scan(&noteUUID, &bookUUID); err != nil {
			continue
		}
		orphans = append(orphans, struct{ noteUUID, bookUUID string }{noteUUID, bookUUID})
	}

	return orphanCount, orphans, nil
}

func warnOrphanedNotes(tx *database.DB) {
	count, orphans, err := findOrphanedNotes(tx)
	if err != nil {
		log.Debug("error checking orphaned notes: %v\n", err)
		return
	}

	if count == 0 {
		return
	}

	log.Debug("Found %d orphaned notes (book doesn't exist locally):\n", count)
	for _, o := range orphans {
		log.Debug("note %s (book %s)\n", o.noteUUID, o.bookUUID)
	}
}

// checkPostSyncIntegrity checks for data integrity issues after sync and warns the user
func checkPostSyncIntegrity(db *database.DB) {
	count, orphans, err := findOrphanedNotes(db)
	if err != nil {
		log.Debug("error checking orphaned notes: %v\n", err)
		return
	}

	if count == 0 {
		return
	}

	log.Warnf("Found %d orphaned notes (referencing non-existent or deleted books):\n", count)
	for _, o := range orphans {
		log.Plainf("  - note %s (missing book: %s)\n", o.noteUUID, o.bookUUID)
	}
}

func sendNotes(ctx context.DnoteCtx, tx *database.DB) (bool, error) {
	isBehind := false

	warnOrphanedNotes(tx)

	rows, err := tx.Query("SELECT uuid, book_uuid, body, deleted, usn, added_on FROM notes WHERE dirty")
	if err != nil {
		return isBehind, errors.Wrap(err, "getting syncable notes")
	}
	defer rows.Close()

	for rows.Next() {
		var note database.Note

		if err = rows.Scan(&note.UUID, &note.BookUUID, &note.Body, &note.Deleted, &note.USN, &note.AddedOn); err != nil {
			return isBehind, errors.Wrap(err, "scanning a syncable note")
		}

		log.Debug("sending note %s (book: %s)\n", note.UUID, note.BookUUID)

		var respUSN int

		// if new, create it in the server, or else, update.
		if note.USN == 0 {
			if note.Deleted {
				// if a note was added and deleted locally, simply expunge
				err = note.Expunge(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "expunging a note locally")
				}

				continue
			} else {
				resp, err := client.CreateNote(ctx, note.BookUUID, note.Body)
				if err != nil {
					log.Debug("failed to create note %s (book: %s): %v\n", note.UUID, note.BookUUID, err)
					isBehind = true
					continue
				}

				note.Dirty = false
				note.USN = resp.Result.USN
				err = note.Update(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "marking note dirty")
				}

				err = note.UpdateUUID(tx, resp.Result.UUID)
				if err != nil {
					return isBehind, errors.Wrap(err, "updating note uuid")
				}

				respUSN = resp.Result.USN
			}
		} else {
			if note.Deleted {
				resp, err := client.DeleteNote(ctx, note.UUID)
				if err != nil {
					return isBehind, errors.Wrap(err, "deleting a note")
				}

				err = note.Expunge(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "expunging a note locally")
				}

				respUSN = resp.Result.USN
			} else {
				resp, err := client.UpdateNote(ctx, note.UUID, note.BookUUID, note.Body)
				if err != nil {
					return isBehind, errors.Wrap(err, "updating a note")
				}

				note.Dirty = false
				note.USN = resp.Result.USN
				err = note.Update(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "marking note dirty")
				}

				respUSN = resp.Result.USN
			}
		}

		lastMaxUSN, err := getLastMaxUSN(tx)
		if err != nil {
			return isBehind, errors.Wrap(err, "getting last max usn")
		}

		log.Debug("sent note %s. response USN %d. last max usn: %d\n", note.UUID, respUSN, lastMaxUSN)

		if respUSN == lastMaxUSN+1 {
			err = updateLastMaxUSN(tx, lastMaxUSN+1)
			if err != nil {
				return isBehind, errors.Wrap(err, "updating last max usn")
			}
		} else {
			isBehind = true
		}
	}

	return isBehind, nil
}

func sendChanges(ctx context.DnoteCtx, tx *database.DB) (bool, error) {
	log.Info("sending changes.")

	var delta int
	err := tx.QueryRow("SELECT (SELECT count(*) FROM notes WHERE dirty) + (SELECT count(*) FROM books WHERE dirty)").Scan(&delta)

	fmt.Printf(" (total %d).", delta)

	log.DebugNewline()

	behind1, err := sendBooks(ctx, tx)
	if err != nil {
		return behind1, errors.Wrap(err, "sending books")
	}

	behind2, err := sendNotes(ctx, tx)
	if err != nil {
		return behind2, errors.Wrap(err, "sending notes")
	}

	fmt.Println(" done.")

	isBehind := behind1 || behind2

	return isBehind, nil
}

func updateLastMaxUSN(tx *database.DB, val int) error {
	if err := database.UpdateSystem(tx, consts.SystemLastMaxUSN, val); err != nil {
		return errors.Wrapf(err, "updating %s", consts.SystemLastMaxUSN)
	}

	return nil
}

func updateLastSyncAt(tx *database.DB, val int64) error {
	if err := database.UpdateSystem(tx, consts.SystemLastSyncAt, val); err != nil {
		return errors.Wrapf(err, "updating %s", consts.SystemLastSyncAt)
	}

	return nil
}

func saveSyncState(tx *database.DB, serverTime int64, serverMaxUSN int, userMaxUSN int) error {
	// Handle last_max_usn update based on server state:
	// - If serverMaxUSN > 0: we got data, update to serverMaxUSN
	// - If serverMaxUSN == 0 && userMaxUSN > 0: empty fragment (caught up), preserve existing
	// - If serverMaxUSN == 0 && userMaxUSN == 0: empty server, reset to 0
	if serverMaxUSN > 0 {
		if err := updateLastMaxUSN(tx, serverMaxUSN); err != nil {
			return errors.Wrap(err, "updating last max usn")
		}
	} else if userMaxUSN == 0 {
		// Server is empty, reset to 0
		if err := updateLastMaxUSN(tx, 0); err != nil {
			return errors.Wrap(err, "updating last max usn")
		}
	}
	// else: empty fragment but server has data, preserve existing last_max_usn

	// Always update last_sync_at (we did communicate with server)
	if err := updateLastSyncAt(tx, serverTime); err != nil {
		return errors.Wrap(err, "updating last sync at")
	}

	return nil
}

// prepareEmptyServerSync marks all local books and notes as dirty when syncing to an empty server.
// This is typically used when switching to a new empty server but wanting to upload existing local data.
// Returns true if preparation was done, false otherwise.
func prepareEmptyServerSync(tx *database.DB) error {
	// Mark all books and notes as dirty and reset USN to 0
	if _, err := tx.Exec("UPDATE books SET usn = 0, dirty = 1 WHERE deleted = 0"); err != nil {
		return errors.Wrap(err, "marking books as dirty")
	}
	if _, err := tx.Exec("UPDATE notes SET usn = 0, dirty = 1 WHERE deleted = 0"); err != nil {
		return errors.Wrap(err, "marking notes as dirty")
	}

	// Reset lastMaxUSN to 0 to match the server
	if err := updateLastMaxUSN(tx, 0); err != nil {
		return errors.Wrap(err, "resetting last max usn")
	}

	return nil
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		// Override APIEndpoint if flag was provided
		if apiEndpointFlag != "" {
			ctx.APIEndpoint = apiEndpointFlag
		}

		if ctx.SessionKey == "" {
			return errors.New("not logged in")
		}

		if err := migrate.Run(ctx, migrate.RemoteSequence, migrate.RemoteMode); err != nil {
			return errors.Wrap(err, "running remote migrations")
		}

		tx, err := ctx.DB.Begin()
		if err != nil {
			return errors.Wrap(err, "beginning a transaction")
		}

		syncState, err := client.GetSyncState(ctx)
		if err != nil {
			return errors.Wrap(err, "getting the sync state from the server")
		}
		lastSyncAt, err := getLastSyncAt(tx)
		if err != nil {
			return errors.Wrap(err, "getting the last sync time")
		}
		lastMaxUSN, err := getLastMaxUSN(tx)
		if err != nil {
			return errors.Wrap(err, "getting the last max_usn")
		}

		log.Debug("lastSyncAt: %d, lastMaxUSN: %d, syncState: %+v\n", lastSyncAt, lastMaxUSN, syncState)

		// Handle a case where server has MaxUSN=0 but local has data (server switch)
		var bookCount, noteCount int
		if err := tx.QueryRow("SELECT count(*) FROM books WHERE deleted = 0").Scan(&bookCount); err != nil {
			return errors.Wrap(err, "counting local books")
		}
		if err := tx.QueryRow("SELECT count(*) FROM notes WHERE deleted = 0").Scan(&noteCount); err != nil {
			return errors.Wrap(err, "counting local notes")
		}

		// If a client has previously synced (lastMaxUSN > 0) but the server was never synced to (MaxUSN = 0),
		// and the client has undeleted books or notes, allow to upload all data to the server.
		// The client might have switched servers or the server might need to be restored for any reasons.
		if syncState.MaxUSN == 0 && lastMaxUSN > 0 && (bookCount > 0 || noteCount > 0) {
			log.Debug("empty server detected: server.MaxUSN=%d, local.MaxUSN=%d, books=%d, notes=%d\n",
				syncState.MaxUSN, lastMaxUSN, bookCount, noteCount)

			log.Warnf("The server is empty but you have local data. Maybe you switched servers?\n")
			log.Debug("server state: MaxUSN = 0 (empty)\n")
			log.Debug("local state: %d books, %d notes (MaxUSN = %d)\n", bookCount, noteCount, lastMaxUSN)

			confirmed, err := ui.Confirm(fmt.Sprintf("Upload %d books and %d notes to the server?", bookCount, noteCount), false)
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, "getting user confirmation")
			}

			if !confirmed {
				tx.Rollback()
				return errors.New("sync cancelled by user")
			}

			fmt.Println() // Add newline after confirmation.

			if err := prepareEmptyServerSync(tx); err != nil {
				return errors.Wrap(err, "preparing for empty server sync")
			}

			// Re-fetch lastMaxUSN after prepareEmptyServerSync
			lastMaxUSN, err = getLastMaxUSN(tx)
			if err != nil {
				return errors.Wrap(err, "getting the last max_usn after prepare")
			}

			log.Debug("prepared empty server sync: marked %d books and %d notes as dirty\n", bookCount, noteCount)
		}

		var syncErr error
		if isFullSync || lastSyncAt < syncState.FullSyncBefore {
			syncErr = fullSync(ctx, tx)
		} else if lastMaxUSN != syncState.MaxUSN {
			syncErr = stepSync(ctx, tx, lastMaxUSN)
		} else {
			// if no need to sync from the server, simply update the last sync timestamp and proceed to send changes
			err = updateLastSyncAt(tx, syncState.CurrentTime)
			if err != nil {
				return errors.Wrap(err, "updating last sync at")
			}
		}
		if syncErr != nil {
			tx.Rollback()
			return errors.Wrap(syncErr, "syncing changes from the server")
		}

		isBehind, err := sendChanges(ctx, tx)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "sending changes")
		}

		// if server state gets ahead of that of client during the sync, do an additional step sync
		if isBehind {
			log.Debug("performing another step sync because client is behind\n")

			updatedLastMaxUSN, err := getLastMaxUSN(tx)
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, "getting the new last max_usn")
			}

			err = stepSync(ctx, tx, updatedLastMaxUSN)
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, "performing the follow-up step sync")
			}

			// After syncing server changes (which resolves conflicts), send local changes again
			// This uploads books/notes that were skipped due to 409 conflicts
			_, err = sendChanges(ctx, tx)
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, "sending changes after conflict resolution")
			}
		}

		if err := tx.Commit(); err != nil {
			return errors.Wrap(err, "committing transaction")
		}

		log.Success("success\n")

		checkPostSyncIntegrity(ctx.DB)

		if err := upgrade.Check(ctx); err != nil {
			log.Error(errors.Wrap(err, "automatically checking updates").Error())
		}

		return nil
	}
}
