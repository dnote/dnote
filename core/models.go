package core

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Book holds a metadata and its notes
type Book struct {
	UUID    string `json:"uuid"`
	Label   string `json:"label"`
	USN     int    `json:"usn"`
	Notes   []Note `json:"notes"`
	Deleted bool   `json:"deleted"`
	Dirty   bool   `json:"dirty"`
}

// Note represents a note
type Note struct {
	UUID     string `json:"uuid"`
	BookUUID string `json:"book_uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
	USN      int    `json:"usn"`
	Public   bool   `json:"public"`
	Deleted  bool   `json:"deleted"`
	Dirty    bool   `json:"dirty"`
}

// NewNote constructs a note with the given data
func NewNote(uuid, bookUUID, content string, addedOn, editedOn int64, usn int, public, deleted, dirty bool) Note {
	return Note{
		UUID:     uuid,
		BookUUID: bookUUID,
		Content:  content,
		AddedOn:  addedOn,
		EditedOn: editedOn,
		USN:      usn,
		Public:   public,
		Deleted:  deleted,
		Dirty:    dirty,
	}
}

// Insert inserts a new note
func (n Note) Insert(tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO notes (uuid, book_uuid, content, added_on, edited_on, usn, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		n.UUID, n.BookUUID, n.Content, n.AddedOn, n.EditedOn, n.USN, n.Public, n.Deleted, n.Dirty)

	if err != nil {
		return errors.Wrapf(err, "inserting note with uuid %s", n.UUID)
	}

	return nil
}

// Update updates the note with the given data
func (n Note) Update(tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE notes SET book_uuid = ?, content = ?, added_on = ?, edited_on = ?, usn = ?, public = ?, deleted = ?, dirty = ? WHERE uuid = ?",
		n.BookUUID, n.Content, n.AddedOn, n.EditedOn, n.USN, n.Public, n.Deleted, n.Dirty, n.UUID)

	if err != nil {
		return errors.Wrapf(err, "updating the note with uuid %s", n.UUID)
	}

	return nil
}

// Expunge hard-deletes the note from the database
func (n Note) Expunge(tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM notes WHERE uuid = ?", n.UUID)
	if err != nil {
		return errors.Wrap(err, "expunging a note locally")
	}

	return nil
}

// NewBook constructs a book with the given data
func NewBook(uuid, label string, usn int, deleted, dirty bool) Book {
	return Book{
		UUID:    uuid,
		Label:   label,
		USN:     usn,
		Deleted: deleted,
		Dirty:   dirty,
	}
}

// Insert inserts a new book
func (b Book) Insert(tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO books (uuid, label, usn, dirty, deleted) VALUES (?, ?, ?, ?, ?)",
		b.UUID, b.Label, b.USN, b.Dirty, b.Deleted)

	if err != nil {
		return errors.Wrapf(err, "inserting book with uuid %s", b.UUID)
	}

	return nil
}

// Update updates the book with the given data
func (b Book) Update(tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE books SET label = ?, usn = ?, dirty = ?, deleted = ? WHERE uuid = ?",
		b.Label, b.USN, b.Dirty, b.Deleted, b.UUID)

	if err != nil {
		return errors.Wrapf(err, "updating the book with uuid %s", b.UUID)
	}

	return nil
}

// Expunge hard-deletes the book from the database
func (b Book) Expunge(tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM books WHERE uuid = ?", b.UUID)
	if err != nil {
		return errors.Wrap(err, "expunging a book locally")
	}

	return nil
}
