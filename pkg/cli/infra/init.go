/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package infra provides operations and definitions for the
// local infrastructure for Dnote
package infra

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/dnote/dnote/pkg/cli/config"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/migrate"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// RunEFunc is a function type of dnote commands
type RunEFunc func(*cobra.Command, []string) error

func newCtx(versionTag string) (context.DnoteCtx, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return context.DnoteCtx{}, errors.Wrap(err, "Failed to get home dir")
	}
	dnoteDir := getDnoteDir(homeDir)

	dnoteDBPath := fmt.Sprintf("%s/%s", dnoteDir, consts.DnoteDBFileName)
	db, err := database.Open(dnoteDBPath)
	if err != nil {
		return context.DnoteCtx{}, errors.Wrap(err, "conntecting to db")
	}

	ctx := context.DnoteCtx{
		HomeDir:  homeDir,
		DnoteDir: dnoteDir,
		Version:  versionTag,
		DB:       db,
	}

	return ctx, nil
}

// Init initializes the Dnote environment and returns a new dnote context
func Init(apiEndpoint, versionTag string) (*context.DnoteCtx, error) {
	ctx, err := newCtx(versionTag)
	if err != nil {
		return nil, errors.Wrap(err, "initializing a context")
	}

	if err := InitFiles(ctx, apiEndpoint); err != nil {
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

	ctx, err = SetupCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "setting up the context")
	}

	log.Debug("Running with Dnote context: %+v\n", context.Redact(ctx))

	return &ctx, nil
}

// SetupCtx populates the context and returns a new context
func SetupCtx(ctx context.DnoteCtx) (context.DnoteCtx, error) {
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
		HomeDir:          ctx.HomeDir,
		DnoteDir:         ctx.DnoteDir,
		Version:          ctx.Version,
		DB:               ctx.DB,
		SessionKey:       sessionKey,
		SessionKeyExpiry: sessionKeyExpiry,
		APIEndpoint:      cf.APIEndpoint,
		Editor:           cf.Editor,
		Clock:            clock.New(),
	}

	return ret, nil
}

func getDnoteDir(homeDir string) string {
	var ret string

	dnoteDirEnv := os.Getenv("DNOTE_DIR")
	if dnoteDirEnv == "" {
		ret = fmt.Sprintf("%s/%s", homeDir, consts.DnoteDirName)
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

	tx.Commit()

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

// initDnoteDir initializes dnote directory if it does not exist yet
func initDnoteDir(ctx context.DnoteCtx) error {
	path := ctx.DnoteDir

	ok, err := utils.FileExists(path)
	if err != nil {
		return errors.Wrap(err, "checking if dnote dir exists")
	}
	if ok {
		return nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return errors.Wrap(err, "Failed to create dnote directory")
	}

	return nil
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

	cf := config.Config{
		Editor:      editor,
		APIEndpoint: apiEndpoint,
	}

	if err := config.Write(ctx, cf); err != nil {
		return errors.Wrap(err, "writing config")
	}

	return nil
}

// InitFiles creates, if necessary, the dnote directory and files inside
func InitFiles(ctx context.DnoteCtx, apiEndpoint string) error {
	if err := initDnoteDir(ctx); err != nil {
		return errors.Wrap(err, "creating the dnote dir")
	}
	if err := initConfigFile(ctx, apiEndpoint); err != nil {
		return errors.Wrap(err, "generating the config file")
	}

	return nil
}
