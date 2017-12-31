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
	APIKey string
}

// Dnote holds the whole dnote data
type Dnote map[string]Book

// Book holds a metadata and its notes
type Book struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Notes []Note `json:"notes"`
}

// Note represents a single microlesson
type Note struct {
	UUID    string `json:"uuid"`
	Content string `json:"content"`
}
