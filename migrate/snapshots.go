package migrate

import "encoding/json"

// v2
type migrateToV2PreNote struct {
	UID     string
	Content string
	AddedOn int64
}
type migrateToV2PostNote struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"editd_on"`
}
type migrateToV2PreBook []migrateToV2PreNote
type migrateToV2PostBook struct {
	Name  string                `json:"name"`
	Notes []migrateToV2PostNote `json:"notes"`
}
type migrateToV2PreDnote map[string]migrateToV2PreBook
type migrateToV2PostDnote map[string]migrateToV2PostBook

//v3
var (
	migrateToV3ActionAddNote = "add_note"
	migrateToV3ActionAddBook = "add_book"
)

type migrateToV3Note struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
}
type migrateToV3Book struct {
	UUID  string            `json:"uuid"`
	Name  string            `json:"name"`
	Notes []migrateToV3Note `json:"notes"`
}
type migrateToV3Dnote map[string]migrateToV3Book
type migrateToV3Action struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

// v4
type migrateToV4PreConfig struct {
	Book   string
	APIKey string
}
type migrateToV4PostConfig struct {
	Editor string
	APIKey string
}

// v5
type migrateToV5AddNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookName string `json:"book_name"`
	Content  string `json:"content"`
}
type migrateToV5RemoveNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookName string `json:"book_name"`
}
type migrateToV5AddBookData struct {
	BookName string `json:"book_name"`
}
type migrateToV5RemoveBookData struct {
	BookName string `json:"book_name"`
}
type migrateToV5PreEditNoteData struct {
	NoteUUID string `json:"note_uuid"`
	BookName string `json:"book_name"`
	Content  string `json:"content"`
}
type migrateToV5PostEditNoteData struct {
	NoteUUID string `json:"note_uuid"`
	FromBook string `json:"from_book"`
	ToBook   string `json:"to_book"`
	Content  string `json:"content"`
}
type migrateToV5PreAction struct {
	ID        int             `json:"id"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}
type migrateToV5PostAction struct {
	UUID      string          `json:"uuid"`
	Schema    int             `json:"schema"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

var (
	migrateToV5ActionAddNote    = "add_note"
	migrateToV5ActionRemoveNote = "remove_note"
	migrateToV5ActionEditNote   = "edit_note"
	migrateToV5ActionAddBook    = "add_book"
	migrateToV5ActionRemoveBook = "remove_book"
)

// v6
type migrateToV6PreNote struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
}
type migrateToV6PostNote struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
	// Make a pointer to test absent values
	Public *bool `json:"public"`
}
type migrateToV6PreBook struct {
	Name  string               `json:"name"`
	Notes []migrateToV6PreNote `json:"notes"`
}
type migrateToV6PostBook struct {
	Name  string                `json:"name"`
	Notes []migrateToV6PostNote `json:"notes"`
}
type migrateToV6PreDnote map[string]migrateToV6PreBook
type migrateToV6PostDnote map[string]migrateToV6PostBook
