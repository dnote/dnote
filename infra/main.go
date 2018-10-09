// Package infra defines dnote structure
package infra

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"time"

	// use sqlite
	_ "github.com/mattn/go-sqlite3"

	"github.com/pkg/errors"
)

var (
	// DnoteDirName is the name of the directory containing dnote files
	DnoteDirName = ".dnote"

	// SystemSchema is the key for schema in the system table
	SystemSchema = "schema"
)

// DnoteCtx is a context holding the information of the current runtime
type DnoteCtx struct {
	HomeDir     string
	DnoteDir    string
	APIEndpoint string
	Version     string
	DB          *sql.DB
}

// Config holds dnote configuration
type Config struct {
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
	Public   bool   `json:"public"`
}

// Timestamp holds time information
type Timestamp struct {
	LastUpgrade int64 `yaml:"last_upgrade"`
	// id of the most recent action synced from the server
	Bookmark int `yaml:"bookmark"`
	// timestamp of the most recent action performed by the cli
	LastAction int64 `yaml:"last_action"`
}

// NewCtx returns a new dnote context
func NewCtx(apiEndpoint, versionTag string) (DnoteCtx, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return DnoteCtx{}, errors.Wrap(err, "Failed to get home dir")
	}
	dnoteDir := getDnoteDir(homeDir)

	dnoteDBPath := fmt.Sprintf("%s/dnote.db", dnoteDir)
	db, err := sql.Open("sqlite3", dnoteDBPath)
	if err != nil {
		return DnoteCtx{}, errors.Wrap(err, "conntecting to db")
	}

	ret := DnoteCtx{
		HomeDir:     homeDir,
		DnoteDir:    dnoteDir,
		APIEndpoint: apiEndpoint,
		Version:     versionTag,
		DB:          db,
	}

	return ret, nil
}

func getDnoteDir(homeDir string) string {
	var ret string

	dnoteDirEnv := os.Getenv("DNOTE_DIR")
	if dnoteDirEnv == "" {
		ret = fmt.Sprintf("%s/%s", homeDir, DnoteDirName)
	} else {
		ret = dnoteDirEnv
	}

	return ret
}

func getHomeDir() (string, error) {
	homeDirEnv := os.Getenv("DNOTE_HOME_DIR")
	if homeDirEnv != "" {
		return homeDirEnv, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get current user")
	}

	return usr.HomeDir, nil
}

// InitDB initializes the database.
// Ideally this process must be a part of migration sequence. But it is performed
// seaprately because it is a prerequisite for legacy migration.
func InitDB(ctx DnoteCtx) error {
	db := ctx.DB

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS notes
		(
			id integer PRIMARY KEY AUTOINCREMENT,
			uuid text NOT NULL,
			book_uuid text NOT NULL,
			content text NOT NULL,
			added_on integer NOT NULL,
			edited_on integer DEFAULT 0,
			public bool DEFAULT false
		)`)
	if err != nil {
		return errors.Wrap(err, "creating notes table")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books
		(
			uuid text PRIMARY KEY,
			label text NOT NULL
		)`)
	if err != nil {
		return errors.Wrap(err, "creating books table")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS actions
		(
			uuid text PRIMARY KEY,
			schema integer NOT NULL,
			type text NOT NULL,
			data text NOT NULL,
			timestamp integer NOT NULL
		)`)
	if err != nil {
		return errors.Wrap(err, "creating actions table")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS system
		(
			key string NOT NULL,
			value text NOT NULL
		)`)
	if err != nil {
		return errors.Wrap(err, "creating system table")
	}

	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_books_label ON books(label);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_notes_uuid ON notes(uuid);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_books_uuid ON books(uuid);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_notes_id ON notes(id);
		CREATE INDEX IF NOT EXISTS idx_notes_book_uuid ON notes(book_uuid);`)
	if err != nil {
		return errors.Wrap(err, "creating indices")
	}

	return nil
}

// InitSystem inserts system data if missing
func InitSystem(ctx DnoteCtx) error {
	db := ctx.DB

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	var bookmarkCount, lastUpgradeCount int
	if err := db.QueryRow("SELECT count(*) FROM system WHERE key = ?", "bookmark").
		Scan(&bookmarkCount); err != nil {
		return errors.Wrap(err, "counting bookmarks")
	}
	if bookmarkCount == 0 {
		_, err := tx.Exec("INSERT INTO system (key, value) VALUES (?, ?)", "bookmark", 0)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "inserting bookmark")
		}
	}

	if err := db.QueryRow("SELECT count(*) FROM system WHERE key = ?", "last_upgrade").
		Scan(&lastUpgradeCount); err != nil {
		return errors.Wrap(err, "counting last_upgrade")
	}
	if lastUpgradeCount == 0 {
		now := time.Now().Unix()
		_, err := tx.Exec("INSERT INTO system (key, value) VALUES (?, ?)", "last_upgrade", now)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "inserting bookmark")
		}
	}

	tx.Commit()

	return nil
}
