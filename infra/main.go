// Package infra defines dnote structure
package infra

// DnoteCtx is a context holding the information of the current runtime
type DnoteCtx struct {
	HomeDir     string
	DnoteDir    string
	APIEndpoint string
}

// Config holds dnote configuration
type Config struct {
	Book   string
	Editor string
	APIKey string
}

// Dnote holds the whole dnote data
type Dnote map[string]Book

// Book holds a metadata and its notes
type Book struct {
	Name  string `json:"name"`
	Notes []Note `json:"notes"`
}

// Note represents a single microlesson
type Note struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
}

// Timestamp holds time information
type Timestamp struct {
	LastUpgrade int64 `yaml:"last_upgrade"`
	// id of the most recent action synced from the server
	Bookmark int `yaml:"bookmark"`
	// timestamp of the most recent action performed by the cli
	LastAction int64 `yaml:"last_action"`
}
