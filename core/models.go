package core

import (
	"github.com/dnote/cli/infra"
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
	Body     string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
	USN      int    `json:"usn"`
	Public   bool   `json:"public"`
	Deleted  bool   `json:"deleted"`
	Dirty    bool   `json:"dirty"`
}

// NewNote constructs a note with the given data
func NewNote(uuid, bookUUID, body string, addedOn, editedOn int64, usn int, public, deleted, dirty bool) Note {
	return Note{
		UUID:     uuid,
		BookUUID: bookUUID,
		Body:     body,
		AddedOn:  addedOn,
		EditedOn: editedOn,
		USN:      usn,
		Public:   public,
		Deleted:  deleted,
		Dirty:    dirty,
	}
}

// Insert inserts a new note
func (n Note) Insert(db *infra.DB) error {
	_, err := db.Exec("INSERT INTO notes (uuid, book_uuid, body, added_on, edited_on, usn, public, deleted, dirty) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		n.UUID, n.BookUUID, n.Body, n.AddedOn, n.EditedOn, n.USN, n.Public, n.Deleted, n.Dirty)

	if err != nil {
		return errors.Wrapf(err, "inserting note with uuid %s", n.UUID)
	}

	return nil
}

// Update updates the note with the given data
func (n Note) Update(db *infra.DB) error {
	_, err := db.Exec("UPDATE notes SET book_uuid = ?, body = ?, added_on = ?, edited_on = ?, usn = ?, public = ?, deleted = ?, dirty = ? WHERE uuid = ?",
		n.BookUUID, n.Body, n.AddedOn, n.EditedOn, n.USN, n.Public, n.Deleted, n.Dirty, n.UUID)

	if err != nil {
		return errors.Wrapf(err, "updating the note with uuid %s", n.UUID)
	}

	return nil
}

// UpdateUUID updates the uuid of a book
func (n *Note) UpdateUUID(db *infra.DB, newUUID string) error {
	_, err := db.Exec("UPDATE notes SET uuid = ? WHERE uuid = ?", newUUID, n.UUID)

	if err != nil {
		return errors.Wrapf(err, "updating note uuid from '%s' to '%s'", n.UUID, newUUID)
	}

	n.UUID = newUUID

	return nil
}

// Expunge hard-deletes the note from the database
func (n Note) Expunge(db *infra.DB) error {
	_, err := db.Exec("DELETE FROM notes WHERE uuid = ?", n.UUID)
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
func (b Book) Insert(db *infra.DB) error {
	_, err := db.Exec("INSERT INTO books (uuid, label, usn, dirty, deleted) VALUES (?, ?, ?, ?, ?)",
		b.UUID, b.Label, b.USN, b.Dirty, b.Deleted)

	if err != nil {
		return errors.Wrapf(err, "inserting book with uuid %s", b.UUID)
	}

	return nil
}

// Update updates the book with the given data
func (b Book) Update(db *infra.DB) error {
	_, err := db.Exec("UPDATE books SET label = ?, usn = ?, dirty = ?, deleted = ? WHERE uuid = ?",
		b.Label, b.USN, b.Dirty, b.Deleted, b.UUID)

	if err != nil {
		return errors.Wrapf(err, "updating the book with uuid %s", b.UUID)
	}

	return nil
}

// UpdateUUID updates the uuid of a book
func (b *Book) UpdateUUID(db *infra.DB, newUUID string) error {
	_, err := db.Exec("UPDATE books SET uuid = ? WHERE uuid = ?", newUUID, b.UUID)

	if err != nil {
		return errors.Wrapf(err, "updating book uuid from '%s' to '%s'", b.UUID, newUUID)
	}

	b.UUID = newUUID

	return nil
}

// Expunge hard-deletes the book from the database
func (b Book) Expunge(db *infra.DB) error {
	_, err := db.Exec("DELETE FROM books WHERE uuid = ?", b.UUID)
	if err != nil {
		return errors.Wrap(err, "expunging a book locally")
	}

	return nil
}
