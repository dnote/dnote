package sync

import (
	"database/sql"
	"fmt"

	"github.com/dnote/cli/client"
	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/dnote/cli/migrate"
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

// NewCmd returns a new sync command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Aliases: []string{"s"},
		Short:   "Sync data with the server",
		Example: example,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.BoolVarP(&isFullSync, "full", "f", false, "perform a full sync instead of incrementally syncing only the changed data.")

	return cmd
}

func getLastSyncAt(tx *sql.Tx) (int, error) {
	var ret int

	err := tx.QueryRow("SELECT value FROM system WHERE key = ?", infra.SystemLastSyncAt).Scan(&ret)
	if err != nil {
		return ret, errors.Wrap(err, "querying last sync time")
	}

	return ret, nil
}

func getLastMaxUSN(tx *sql.Tx) (int, error) {
	var ret int

	err := tx.QueryRow("SELECT value FROM system WHERE key = ?", infra.SystemLastMaxUSN).Scan(&ret)
	if err != nil {
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
	MaxCurrentTime int64
}

func (l syncList) getLength() int {
	return len(l.Notes) + len(l.Books) + len(l.ExpungedNotes) + len(l.ExpungedBooks)
}

func newSyncList(fragments []client.SyncFragment) syncList {
	notes := map[string]client.SyncFragNote{}
	books := map[string]client.SyncFragBook{}
	expungedNotes := map[string]bool{}
	expungedBooks := map[string]bool{}
	var maxUSN int
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
		if fragment.CurrentTime > maxCurrentTime {
			maxCurrentTime = fragment.CurrentTime
		}
	}

	return syncList{
		Notes:          notes,
		Books:          books,
		ExpungedNotes:  expungedNotes,
		ExpungedBooks:  expungedBooks,
		MaxUSN:         maxUSN,
		MaxCurrentTime: maxCurrentTime,
	}
}

// getSyncList gets a list of all sync fragments after the specified usn
// and aggregates them into a syncList data structure
func getSyncList(ctx infra.DnoteCtx, apiKey string, afterUSN int) (syncList, error) {
	fragments, err := getSyncFragments(ctx, apiKey, afterUSN)
	if err != nil {
		return syncList{}, errors.Wrap(err, "getting sync fragments")
	}

	ret := newSyncList(fragments)

	return ret, nil
}

// getSyncFragments repeatedly gets all sync fragments after the specified usn until there is no more new data
// remaining and returns the buffered list
func getSyncFragments(ctx infra.DnoteCtx, apiKey string, afterUSN int) ([]client.SyncFragment, error) {
	var buf []client.SyncFragment

	nextAfterUSN := afterUSN

	for {
		resp, err := client.GetSyncFragment(ctx, apiKey, nextAfterUSN)
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

	log.Debug("received sync fragments: %+v\n", buf)

	return buf, nil
}

// resolveLabel resolves a book label conflict by repeatedly appending an increasing integer
// to the label until it finds a unique label. It returns the first non-conflicting label.
func resolveLabel(tx *sql.Tx, label string) (string, error) {
	var ret string

	for i := 2; ; i++ {
		ret = fmt.Sprintf("%s (%d)", label, i)

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
func mergeBook(tx *sql.Tx, b client.SyncFragBook, mode int) error {
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

		if _, err := tx.Exec("UPDATE books SET label = ?, dirty = ? WHERE label = ?", newLabel, true, b.Label); err != nil {
			return errors.Wrap(err, "resolving duplicate book label")
		}
	}

	if mode == modeInsert {
		book := core.NewBook(b.UUID, b.Label, b.USN, false, false)
		if err := book.Insert(tx); err != nil {
			return errors.Wrapf(err, "inserting note with uuid %s", b.UUID)
		}
	} else if mode == modeUpdate {
		// TODO: if the client copy is dirty, perform field-by-field merge and report conflict instead of overwriting
		if _, err := tx.Exec("UPDATE books SET usn = ?, uuid = ?, label = ?, deleted = ? WHERE uuid = ?",
			b.USN, b.UUID, b.Label, b.Deleted, b.UUID); err != nil {
			return errors.Wrapf(err, "updating local book %s", b.UUID)
		}
	}

	return nil
}

func stepSyncBook(tx *sql.Tx, b client.SyncFragBook) error {
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

func mergeNote(tx *sql.Tx, serverNote client.SyncFragNote, localNote core.Note) error {
	var bookDeleted bool
	err := tx.QueryRow("SELECT deleted FROM books WHERE uuid = ?", localNote.BookUUID).Scan(&bookDeleted)
	if err != nil {
		return errors.Wrapf(err, "checking if local book %s is deleted", localNote.BookUUID)
	}

	// if the book is deleted, noop
	if bookDeleted {
		return nil
	}

	// if the local copy is deleted, and the it was edited on the server, override with server values and mark it not dirty.
	if localNote.Deleted {
		if _, err := tx.Exec("UPDATE notes SET usn = ?, book_uuid = ?, content = ?, edited_on = ?, deleted = ?, public = ?, dirty = ? WHERE uuid = ?",
			serverNote.USN, serverNote.BookUUID, serverNote.Content, serverNote.EditedOn, serverNote.Deleted, serverNote.Public, false, serverNote.UUID); err != nil {
			return errors.Wrapf(err, "updating local note %s", serverNote.UUID)
		}

		return nil
	}

	// TODO: if the client copy is dirty, perform field-by-field merge and report conflict instead of overwriting
	if _, err := tx.Exec("UPDATE notes SET usn = ?, book_uuid = ?, content = ?, edited_on = ?, deleted = ?, public = ?  WHERE uuid = ?",
		serverNote.USN, serverNote.BookUUID, serverNote.Content, serverNote.EditedOn, serverNote.Deleted, serverNote.Public, serverNote.UUID); err != nil {
		return errors.Wrapf(err, "updating local note %s", serverNote.UUID)
	}

	return nil
}

func stepSyncNote(tx *sql.Tx, n client.SyncFragNote) error {
	var localNote core.Note
	err := tx.QueryRow("SELECT usn, book_uuid, dirty, deleted FROM notes WHERE uuid = ?", n.UUID).
		Scan(&localNote.USN, &localNote.BookUUID, &localNote.Dirty, &localNote.Deleted)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrapf(err, "getting local note %s", n.UUID)
	}

	// if note exists in the server and does not exist in the client, insert the note.
	if err == sql.ErrNoRows {
		note := core.NewNote(n.UUID, n.BookUUID, n.Content, n.AddedOn, n.EditedOn, n.USN, n.Public, n.Deleted, false)

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

func fullSyncNote(tx *sql.Tx, n client.SyncFragNote) error {
	var localNote core.Note
	err := tx.QueryRow("SELECT usn,book_uuid, dirty, deleted FROM notes WHERE uuid = ?", n.UUID).
		Scan(&localNote.USN, &localNote.BookUUID, &localNote.Dirty, &localNote.Deleted)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrapf(err, "getting local note %s", n.UUID)
	}

	// if note exists in the server and does not exist in the client, insert the note.
	if err == sql.ErrNoRows {
		note := core.NewNote(n.UUID, n.BookUUID, n.Content, n.AddedOn, n.EditedOn, n.USN, n.Public, n.Deleted, false)

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

func syncDeleteNote(tx *sql.Tx, noteUUID string) error {
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
func checkNotesPristine(tx *sql.Tx, bookUUID string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM notes WHERE book_uuid = ? AND dirty = ?", bookUUID, true).Scan(&count); err != nil {
		return false, errors.Wrapf(err, "counting notes that are dirty in book %s", bookUUID)
	}

	if count > 0 {
		return false, nil
	}

	return true, nil
}

func syncDeleteBook(tx *sql.Tx, bookUUID string) error {
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

func fullSyncBook(tx *sql.Tx, b client.SyncFragBook) error {
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
func cleanLocalNotes(tx *sql.Tx, fullList *syncList) error {
	rows, err := tx.Query("SELECT uuid, usn, dirty FROM notes")
	if err != nil {
		return errors.Wrap(err, "getting local notes")
	}
	defer rows.Close()

	for rows.Next() {
		var note core.Note
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
func cleanLocalBooks(tx *sql.Tx, fullList *syncList) error {
	rows, err := tx.Query("SELECT uuid, usn, dirty FROM books")
	if err != nil {
		return errors.Wrap(err, "getting local books")
	}
	defer rows.Close()

	for rows.Next() {
		var book core.Book
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

func fullSync(ctx infra.DnoteCtx, tx *sql.Tx, apiKey string) error {
	log.Debug("performing a full sync\n")

	list, err := getSyncList(ctx, apiKey, 0)
	if err != nil {
		return errors.Wrap(err, "getting sync list")
	}

	log.Infof("resolving delta (total %d).", list.getLength())

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

	err = saveSyncState(tx, list.MaxCurrentTime, list.MaxUSN)
	if err != nil {
		return errors.Wrap(err, "saving sync state")
	}

	fmt.Println(" done.")

	return nil
}

func stepSync(ctx infra.DnoteCtx, tx *sql.Tx, apiKey string, afterUSN int) error {
	log.Debug("performing a step sync\n")

	list, err := getSyncList(ctx, apiKey, afterUSN)
	if err != nil {
		return errors.Wrap(err, "getting sync list")
	}

	log.Infof("resolving delta (total %d).", list.getLength())

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

	err = saveSyncState(tx, list.MaxCurrentTime, list.MaxUSN)
	if err != nil {
		return errors.Wrap(err, "saving sync state")
	}

	fmt.Println(" done.")

	return nil
}

func sendBooks(ctx infra.DnoteCtx, tx *sql.Tx, apiKey string) (bool, error) {
	isBehind := false

	rows, err := tx.Query("SELECT uuid, label, usn, deleted FROM books WHERE dirty")
	if err != nil {
		return isBehind, errors.Wrap(err, "getting syncable books")
	}
	defer rows.Close()

	for rows.Next() {
		var book core.Book

		if err = rows.Scan(&book.UUID, &book.Label, &book.USN, &book.Deleted); err != nil {
			return isBehind, errors.Wrap(err, "scanning a syncable book")
		}

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
				resp, err := client.CreateBook(ctx, apiKey, book.Label, book.UUID)
				if err != nil {
					return isBehind, errors.Wrap(err, "creating a book")
				}

				book.Dirty = false
				book.USN = resp.Book.USN
				err = book.Update(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "marking book dirty")
				}

				respUSN = resp.Book.USN
			}
		} else {
			if book.Deleted {
				resp, err := client.DeleteBook(ctx, apiKey, book.UUID)
				if err != nil {
					return isBehind, errors.Wrap(err, "deleting a book")
				}

				err = book.Expunge(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "expunging a book locally")
				}

				respUSN = resp.Book.USN
			} else {
				resp, err := client.UpdateBook(ctx, apiKey, book.Label, book.UUID)
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

func sendNotes(ctx infra.DnoteCtx, tx *sql.Tx, apiKey string) (bool, error) {
	isBehind := false

	rows, err := tx.Query("SELECT uuid, book_uuid, content, public, deleted, usn FROM notes WHERE dirty")
	if err != nil {
		return isBehind, errors.Wrap(err, "getting syncable notes")
	}
	defer rows.Close()

	for rows.Next() {
		var note core.Note

		if err = rows.Scan(&note.UUID, &note.BookUUID, &note.Content, &note.Public, &note.Deleted, &note.USN); err != nil {
			return isBehind, errors.Wrap(err, "scanning a syncable note")
		}

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
				resp, err := client.CreateNote(ctx, apiKey, note.UUID, note.BookUUID, note.Content)
				if err != nil {
					return isBehind, errors.Wrap(err, "creating a note")
				}

				note.Dirty = false
				note.USN = resp.Result.USN
				err = note.Update(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "marking note dirty")
				}

				respUSN = resp.Result.USN
			}
		} else {
			if note.Deleted {
				resp, err := client.DeleteNote(ctx, apiKey, note.UUID)
				if err != nil {
					return isBehind, errors.Wrap(err, "deleting a note")
				}

				// TODO: server-side, DELETE endpoint should not return error in case note does not exist
				// rather rpely with 202 or something and handle here accordingly. reason is if the program
				// fails after sending DELETE http call, the note will not be expunged locally and cli will try to
				// delete again
				err = note.Expunge(tx)
				if err != nil {
					return isBehind, errors.Wrap(err, "expunging a note locally")
				}

				respUSN = resp.Result.USN
			} else {
				resp, err := client.UpdateNote(ctx, apiKey, note.UUID, note.BookUUID, note.Content, note.Public)
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

func sendChanges(ctx infra.DnoteCtx, tx *sql.Tx, apiKey string) (bool, error) {
	var delta int
	err := tx.QueryRow("SELECT (SELECT count(*) FROM notes WHERE dirty) + (SELECT count(*) FROM books WHERE dirty)").Scan(&delta)

	log.Infof("sending changes (total %d).", delta)

	isBehind, err := sendBooks(ctx, tx, apiKey)
	if err != nil {
		return isBehind, errors.Wrap(err, "sending books")
	}

	isBehind, err = sendNotes(ctx, tx, apiKey)
	if err != nil {
		return isBehind, errors.Wrap(err, "sending notes")
	}

	fmt.Println(" done.")

	return isBehind, nil
}

func updateLastMaxUSN(tx *sql.Tx, val int) error {
	_, err := tx.Exec("UPDATE system SET value = ? WHERE key = ?", val, infra.SystemLastMaxUSN)
	if err != nil {
		return errors.Wrapf(err, "updating %s", infra.SystemLastMaxUSN)
	}

	return nil
}

func updateLastSyncAt(tx *sql.Tx, val int64) error {
	_, err := tx.Exec("UPDATE system SET value = ? WHERE key = ?", val, infra.SystemLastSyncAt)
	if err != nil {
		return errors.Wrapf(err, "updating %s", infra.SystemLastSyncAt)
	}

	return nil
}

func saveSyncState(tx *sql.Tx, serverTime int64, serverMaxUSN int) error {
	if err := updateLastMaxUSN(tx, serverMaxUSN); err != nil {
		return errors.Wrap(err, "updating last max usn")
	}
	if err := updateLastSyncAt(tx, serverTime); err != nil {
		return errors.Wrap(err, "updating last sync at")
	}

	return nil
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		config, err := core.ReadConfig(ctx)
		if err != nil {
			return errors.Wrap(err, "reading the config")
		}
		if config.APIKey == "" {
			log.Error("login required. please run `dnote login`\n")
			return nil
		}

		if err := migrate.Run(ctx, migrate.RemoteSequence, migrate.RemoteMode); err != nil {
			return errors.Wrap(err, "running remote migrations")
		}

		db := ctx.DB
		tx, err := db.Begin()
		if err != nil {
			return errors.Wrap(err, "beginning a transaction")
		}

		syncState, err := client.GetSyncState(config.APIKey, ctx)
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

		var syncErr error
		if isFullSync || lastSyncAt < syncState.FullSyncBefore {
			syncErr = fullSync(ctx, tx, config.APIKey)
		} else if lastMaxUSN != syncState.MaxUSN {
			syncErr = stepSync(ctx, tx, config.APIKey, lastMaxUSN)
		} else {
			// if no need to sync from the server, simply update the last sync timestamp and proceed to send changes
			err = updateLastSyncAt(tx, syncState.CurrentTime)
			if err != nil {
				return errors.Wrap(err, "updating last sync at")
			}
		}
		if syncErr != nil {
			tx.Rollback()
			return errors.Wrap(err, "syncing changes from the server")
		}

		isBehind, err := sendChanges(ctx, tx, config.APIKey)
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

			err = stepSync(ctx, tx, config.APIKey, updatedLastMaxUSN)
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, "performing the follow-up step sync")
			}
		}

		tx.Commit()

		log.Success("success\n")

		if err := core.CheckUpdate(ctx); err != nil {
			log.Error(errors.Wrap(err, "automatically checking updates").Error())
		}

		return nil
	}
}
