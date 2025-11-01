/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package infra provides operations and definitions for the
// local infrastructure for Dnote
package infra

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dnote/dnote/pkg/cli/client"
	"github.com/dnote/dnote/pkg/cli/config"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/migrate"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/dirs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	// DefaultAPIEndpoint is the default API endpoint used when none is configured
	DefaultAPIEndpoint = "http://localhost:3001/api"
)

// RunEFunc is a function type of dnote commands
type RunEFunc func(*cobra.Command, []string) error

func checkLegacyDBPath() (string, bool) {
	legacyDnoteDir := getLegacyDnotePath(dirs.Home)
	ok, err := utils.FileExists(legacyDnoteDir)
	if ok {
		return legacyDnoteDir, true
	}

	if err != nil {
		log.Error(errors.Wrapf(err, "checking legacy dnote directory at %s", legacyDnoteDir).Error())
	}

	return "", false
}

func getDBPath(paths context.Paths, customPath string) string {
	// If custom path is provided, use it
	if customPath != "" {
		return customPath
	}

	legacyDnoteDir, ok := checkLegacyDBPath()
	if ok {
		return fmt.Sprintf("%s/%s", legacyDnoteDir, consts.DnoteDBFileName)
	}

	return fmt.Sprintf("%s/%s/%s", paths.Data, consts.DnoteDirName, consts.DnoteDBFileName)
}

// newBaseCtx creates a minimal context with paths and database connection.
// This base context is used for file and database initialization before
// being enriched with config values by setupCtx.
func newBaseCtx(versionTag, customDBPath string) (context.DnoteCtx, error) {
	dnoteDir := getLegacyDnotePath(dirs.Home)
	paths := context.Paths{
		Home:        dirs.Home,
		Config:      dirs.ConfigHome,
		Data:        dirs.DataHome,
		Cache:       dirs.CacheHome,
		LegacyDnote: dnoteDir,
	}

	dbPath := getDBPath(paths, customDBPath)

	db, err := database.Open(dbPath)
	if err != nil {
		return context.DnoteCtx{}, errors.Wrap(err, "conntecting to db")
	}

	ctx := context.DnoteCtx{
		Paths:   paths,
		Version: versionTag,
		DB:      db,
	}

	return ctx, nil
}

// Init initializes the Dnote environment and returns a new dnote context
// apiEndpoint is used when creating a new config file (e.g., from ldflags during tests)
func Init(versionTag, apiEndpoint, dbPath string) (*context.DnoteCtx, error) {
	ctx, err := newBaseCtx(versionTag, dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "initializing a context")
	}

	if err := initFiles(ctx, apiEndpoint); err != nil {
		return nil, errors.Wrap(err, "initializing files")
	}

	if err := InitDB(ctx); err != nil {
		return nil, errors.Wrap(err, "initializing database")
	}
	if err := InitSystem(ctx); err != nil {
		return nil, errors.Wrap(err, "initializing system data")
	}

	if err := migrate.Legacy(ctx); err != nil {
		return nil, errors.Wrap(err, "running legacy migration")
	}
	if err := migrate.Run(ctx, migrate.LocalSequence, migrate.LocalMode); err != nil {
		return nil, errors.Wrap(err, "running migration")
	}

	ctx, err = setupCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "setting up the context")
	}

	log.Debug("context: %+v\n", context.Redact(ctx))

	return &ctx, nil
}

// setupCtx enriches the base context with values from config file and database.
// This is called after files and database have been initialized.
func setupCtx(ctx context.DnoteCtx) (context.DnoteCtx, error) {
	db := ctx.DB

	var sessionKey string
	var sessionKeyExpiry int64

	err := db.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemSessionKey).Scan(&sessionKey)
	if err != nil && err != sql.ErrNoRows {
		return ctx, errors.Wrap(err, "finding sesison key")
	}
	err = db.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemSessionKeyExpiry).Scan(&sessionKeyExpiry)
	if err != nil && err != sql.ErrNoRows {
		return ctx, errors.Wrap(err, "finding sesison key expiry")
	}

	cf, err := config.Read(ctx)
	if err != nil {
		return ctx, errors.Wrap(err, "reading config")
	}

	ret := context.DnoteCtx{
		Paths:              ctx.Paths,
		Version:            ctx.Version,
		DB:                 ctx.DB,
		SessionKey:         sessionKey,
		SessionKeyExpiry:   sessionKeyExpiry,
		APIEndpoint:        cf.APIEndpoint,
		Editor:             cf.Editor,
		Clock:              clock.New(),
		EnableUpgradeCheck: cf.EnableUpgradeCheck,
		HTTPClient:         client.NewRateLimitedHTTPClient(),
	}

	return ret, nil
}

// getLegacyDnotePath returns a legacy dnote directory path placed under
// the user's home directory
func getLegacyDnotePath(homeDir string) string {
	return fmt.Sprintf("%s/%s", homeDir, consts.LegacyDnoteDirName)
}

// InitDB initializes the database.
// Ideally this process must be a part of migration sequence. But it is performed
// seaprately because it is a prerequisite for legacy migration.
func InitDB(ctx context.DnoteCtx) error {
	log.Debug("initializing the database\n")

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

func initSystemKV(db *database.DB, key string, val string) error {
	var count int
	if err := db.QueryRow("SELECT count(*) FROM system WHERE key = ?", key).Scan(&count); err != nil {
		return errors.Wrapf(err, "counting %s", key)
	}

	if count > 0 {
		return nil
	}

	if _, err := db.Exec("INSERT INTO system (key, value) VALUES (?, ?)", key, val); err != nil {
		db.Rollback()
		return errors.Wrapf(err, "inserting %s %s", key, val)
	}

	return nil
}

// InitSystem inserts system data if missing
func InitSystem(ctx context.DnoteCtx) error {
	log.Debug("initializing the system\n")

	db := ctx.DB

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	nowStr := strconv.FormatInt(time.Now().Unix(), 10)
	if err := initSystemKV(tx, consts.SystemLastUpgrade, nowStr); err != nil {
		return errors.Wrapf(err, "initializing system config for %s", consts.SystemLastUpgrade)
	}
	if err := initSystemKV(tx, consts.SystemLastMaxUSN, "0"); err != nil {
		return errors.Wrapf(err, "initializing system config for %s", consts.SystemLastMaxUSN)
	}
	if err := initSystemKV(tx, consts.SystemLastSyncAt, "0"); err != nil {
		return errors.Wrapf(err, "initializing system config for %s", consts.SystemLastSyncAt)
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "committing transaction")
	}

	return nil
}

// getEditorCommand returns the system's editor command with appropriate flags,
// if necessary, to make the command wait until editor is close to exit.
func getEditorCommand() string {
	editor := os.Getenv("EDITOR")

	var ret string

	switch editor {
	case "atom":
		ret = "atom -w"
	case "subl":
		ret = "subl -n -w"
	case "code":
		ret = "code -n -w"
	case "mate":
		ret = "mate -w"
	case "vim":
		ret = "vim"
	case "nano":
		ret = "nano"
	case "emacs":
		ret = "emacs"
	case "nvim":
		ret = "nvim"
	default:
		ret = "vi"
	}

	return ret
}


// initConfigFile populates a new config file if it does not exist yet
func initConfigFile(ctx context.DnoteCtx, apiEndpoint string) error {
	path := config.GetPath(ctx)
	ok, err := utils.FileExists(path)
	if err != nil {
		return errors.Wrap(err, "checking if config exists")
	}
	if ok {
		return nil
	}

	editor := getEditorCommand()

	// Use default API endpoint if none provided
	endpoint := apiEndpoint
	if endpoint == "" {
		endpoint = DefaultAPIEndpoint
	}

	cf := config.Config{
		Editor:             editor,
		APIEndpoint:        endpoint,
		EnableUpgradeCheck: true,
	}

	if err := config.Write(ctx, cf); err != nil {
		return errors.Wrap(err, "writing config")
	}

	return nil
}

// initFiles creates, if necessary, the dnote directory and files inside
func initFiles(ctx context.DnoteCtx, apiEndpoint string) error {
	if err := context.InitDnoteDirs(ctx.Paths); err != nil {
		return errors.Wrap(err, "creating the dnote dir")
	}
	if err := initConfigFile(ctx, apiEndpoint); err != nil {
		return errors.Wrap(err, "generating the config file")
	}

	return nil
}
