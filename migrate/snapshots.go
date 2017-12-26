package migrate

// v2
type migrateToV2PreNote struct {
	UID     string
	Content string
	AddedOn int64
}
type migrateToV2PostNote struct {
	UUID    string
	Content string
	AddedOn int64
}
type migrateToV2PreBook []migrateToV2PreNote
type migrateToV2PostBook struct {
	UUID  string
	Notes []migrateToV2PostNote
}
type migrateToV2PreDnote map[string]migrateToV2PreBook
type migrateToV2PostDnote map[string]migrateToV2PostBook

//v3
var (
	migrateToV3ActionAddNote = "add_note"
	migrateToV3ActionAddBook = "add_book"
)

type migrateToV3PreNote struct {
	UUID    string
	Content string
	AddedOn int64
}
type migrateToV3PostNote struct {
	UUID    string
	Content string
}
type migrateToV3PreBook struct {
	UUID  string
	Notes []migrateToV3PreNote
}
type migrateToV3PostBook struct {
	UUID  string
	Notes []migrateToV3PostNote
}
type migrateToV3PreDnote map[string]migrateToV3PreBook
type migrateToV3PostDnote map[string]migrateToV3PostBook
type migrateToV3Action struct {
	Type      string
	Data      map[string]interface{}
	Timestamp int64
}
