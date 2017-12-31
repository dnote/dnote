package migrate

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
	UUID  string                `json:"uuid"`
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
