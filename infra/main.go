// Package infra defines dnote structure
package infra

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"os/user"

	// use sqlite
	_ "github.com/mattn/go-sqlite3"

	"github.com/pkg/errors"
)

var (
	// DnoteDirName is the name of the directory containing dnote files
	DnoteDirName = ".dnote"

	// SystemSchema is the key for schema in the system table
	SystemSchema = "schema"
	// SystemRemoteSchema is the key for remote schema in the system table
	SystemRemoteSchema = "remote_schema"
	// SystemLastSyncAt is the timestamp of the server at the last sync
	SystemLastSyncAt = "last_sync_time"
	// SystemLastMaxUSN is the user's max_usn from the server at the alst sync
	SystemLastMaxUSN = "last_max_usn"
	// SystemLastUpgrade is the timestamp at which the system more recently checked for an upgrade
	SystemLastUpgrade = "last_upgrade"
	// SystemCipherKey is the encryption key
	SystemCipherKey = "enc_key"
	// SystemSessionKey is the session key
	SystemSessionKey = "session_token"
	// SystemSessionKeyExpiry is the timestamp at which the session key will expire
	SystemSessionKeyExpiry = "session_token_expiry"
)

// DnoteCtx is a context holding the information of the current runtime
type DnoteCtx struct {
	HomeDir          string
	DnoteDir         string
	APIEndpoint      string
	Version          string
	DB               *DB
	SessionKey       string
	SessionKeyExpiry int64
	CipherKey        []byte
}

// Config holds dnote configuration
type Config struct {
	Editor string
}

// NewCtx returns a new dnote context
func NewCtx(apiEndpoint, versionTag string) (DnoteCtx, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return DnoteCtx{}, errors.Wrap(err, "Failed to get home dir")
	}
	dnoteDir := getDnoteDir(homeDir)

	dnoteDBPath := fmt.Sprintf("%s/dnote.db", dnoteDir)
	db, err := OpenDB(dnoteDBPath)
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

// SetupCtx populates context and returns a new context
func SetupCtx(ctx DnoteCtx) (DnoteCtx, error) {
	db := ctx.DB

	var sessionKey, cipherKeyB64 string
	var sessionKeyExpiry int64

	err := db.QueryRow("SELECT value FROM system WHERE key = ?", SystemSessionKey).Scan(&sessionKey)
	if err != nil && err != sql.ErrNoRows {
		return ctx, errors.Wrap(err, "finding sesison key")
	}
	err = db.QueryRow("SELECT value FROM system WHERE key = ?", SystemCipherKey).Scan(&cipherKeyB64)
	if err != nil && err != sql.ErrNoRows {
		return ctx, errors.Wrap(err, "finding sesison key")
	}
	err = db.QueryRow("SELECT value FROM system WHERE key = ?", SystemSessionKeyExpiry).Scan(&sessionKeyExpiry)
	if err != nil && err != sql.ErrNoRows {
		return ctx, errors.Wrap(err, "finding sesison key expiry")
	}

	cipherKey, err := base64.StdEncoding.DecodeString(cipherKeyB64)
	if err != nil {
		return ctx, errors.Wrap(err, "decoding cipherKey from base64")
	}

	ret := DnoteCtx{
		HomeDir:          ctx.HomeDir,
		DnoteDir:         ctx.DnoteDir,
		APIEndpoint:      ctx.APIEndpoint,
		Version:          ctx.Version,
		DB:               ctx.DB,
		SessionKey:       sessionKey,
		SessionKeyExpiry: sessionKeyExpiry,
		CipherKey:        cipherKey,
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS system
		(
			key string NOT NULL,
			value text NOT NULL
		)`)
	if err != nil {
		return errors.Wrap(err, "creating system table")
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

	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_books_label ON books(label);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_notes_uuid ON notes(uuid);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_books_uuid ON books(uuid);
		CREATE INDEX IF NOT EXISTS idx_notes_book_uuid ON notes(book_uuid);`)
	if err != nil {
		return errors.Wrap(err, "creating indices")
	}

	return nil
}
