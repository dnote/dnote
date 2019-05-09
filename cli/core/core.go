/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package core

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/dnote/dnote/cli/infra"
	"github.com/dnote/dnote/cli/log"
	"github.com/dnote/dnote/cli/utils"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	// ConfigFilename is the name of the config file
	ConfigFilename = "dnoterc"
	// TmpContentFilename is the name of the temporary file that holds editor input
	TmpContentFilename = "DNOTE_TMPCONTENT"
)

// RunEFunc is a function type of dnote commands
type RunEFunc func(*cobra.Command, []string) error

// GetConfigPath returns the path to the dnote config file
func GetConfigPath(ctx infra.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, ConfigFilename)
}

// GetDnoteTmpContentPath returns the path to the temporary file containing
// content being added or edited
func GetDnoteTmpContentPath(ctx infra.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, TmpContentFilename)
}

// GetBookUUID returns a uuid of a book given a label
func GetBookUUID(ctx infra.DnoteCtx, label string) (string, error) {
	db := ctx.DB

	var ret string
	err := db.QueryRow("SELECT uuid FROM books WHERE label = ?", label).Scan(&ret)
	if err == sql.ErrNoRows {
		return ret, errors.Errorf("book '%s' not found", label)
	} else if err != nil {
		return ret, errors.Wrap(err, "querying the book")
	}

	return ret, nil
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

// InitFiles creates, if necessary, the dnote directory and files inside
func InitFiles(ctx infra.DnoteCtx) error {
	if err := initDnoteDir(ctx); err != nil {
		return errors.Wrap(err, "creating the dnote dir")
	}
	if err := initConfigFile(ctx); err != nil {
		return errors.Wrap(err, "generating the config file")
	}

	return nil
}

// initConfigFile populates a new config file if it does not exist yet
func initConfigFile(ctx infra.DnoteCtx) error {
	path := GetConfigPath(ctx)

	if utils.FileExists(path) {
		return nil
	}

	editor := getEditorCommand()

	config := infra.Config{
		Editor: editor,
	}

	b, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "marshalling config into YAML")
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return errors.Wrap(err, "writing the config file")
	}

	return nil
}

// initDnoteDir initializes dnote directory if it does not exist yet
func initDnoteDir(ctx infra.DnoteCtx) error {
	path := ctx.DnoteDir

	if utils.FileExists(path) {
		return nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return errors.Wrap(err, "Failed to create dnote directory")
	}

	return nil
}

// WriteConfig writes the config to the config file
func WriteConfig(ctx infra.DnoteCtx, config infra.Config) error {
	d, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "marhsalling config")
	}

	configPath := GetConfigPath(ctx)

	err = ioutil.WriteFile(configPath, d, 0644)
	if err != nil {
		errors.Wrap(err, "writing the config file")
	}

	return nil
}

// LogAction logs action and updates the last_action
func LogAction(tx *sql.Tx, schema int, actionType, data string, timestamp int64) error {
	uuid := uuid.NewV4().String()

	_, err := tx.Exec(`INSERT INTO actions (uuid, schema, type, data, timestamp)
	VALUES (?, ?, ?, ?, ?)`, uuid, schema, actionType, data, timestamp)
	if err != nil {
		return errors.Wrap(err, "inserting an action")
	}

	_, err = tx.Exec("UPDATE system SET value = ? WHERE key = ?", timestamp, "last_action")
	if err != nil {
		return errors.Wrap(err, "updating last_action")
	}

	return nil
}

// ReadConfig reads the config file
func ReadConfig(ctx infra.DnoteCtx) (infra.Config, error) {
	var ret infra.Config

	configPath := GetConfigPath(ctx)
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return ret, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, errors.Wrap(err, "unmarshalling config")
	}

	return ret, nil
}

// SanitizeContent sanitizes note content
func SanitizeContent(s string) string {
	var ret string

	ret = strings.Trim(s, " ")

	// Remove newline at the end of the file because POSIX defines a line as
	// characters followed by a newline
	ret = strings.TrimSuffix(ret, "\n")
	ret = strings.TrimSuffix(ret, "\r\n")

	return ret
}

func newEditorCmd(ctx infra.DnoteCtx, fpath string) (*exec.Cmd, error) {
	config, err := ReadConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	args := strings.Fields(config.Editor)
	args = append(args, fpath)

	return exec.Command(args[0], args[1:]...), nil
}

// GetEditorInput gets the user input by launching a text editor and waiting for
// it to exit
func GetEditorInput(ctx infra.DnoteCtx, fpath string, content *string) error {
	if !utils.FileExists(fpath) {
		f, err := os.Create(fpath)
		if err != nil {
			return errors.Wrap(err, "creating a temporary content file")
		}
		err = f.Close()
		if err != nil {
			return errors.Wrap(err, "closing the temporary content file")
		}
	}

	cmd, err := newEditorCmd(ctx, fpath)
	if err != nil {
		return errors.Wrap(err, "creating an editor command")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "launching an editor")
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Wrap(err, "waiting for the editor")
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return errors.Wrap(err, "reading the temporary content file")
	}

	err = os.Remove(fpath)
	if err != nil {
		return errors.Wrap(err, "removing the temporary content file")
	}

	raw := string(b)
	c := SanitizeContent(raw)

	*content = c

	return nil
}

func initSystemKV(db *infra.DB, key string, val string) error {
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
func InitSystem(ctx infra.DnoteCtx) error {
	db := ctx.DB

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	nowStr := strconv.FormatInt(time.Now().Unix(), 10)
	if err := initSystemKV(tx, infra.SystemLastUpgrade, nowStr); err != nil {
		return errors.Wrapf(err, "initializing system config for %s", infra.SystemLastUpgrade)
	}
	if err := initSystemKV(tx, infra.SystemLastMaxUSN, "0"); err != nil {
		return errors.Wrapf(err, "initializing system config for %s", infra.SystemLastMaxUSN)
	}
	if err := initSystemKV(tx, infra.SystemLastSyncAt, "0"); err != nil {
		return errors.Wrapf(err, "initializing system config for %s", infra.SystemLastSyncAt)
	}

	tx.Commit()

	return nil
}

// GetValidSession returns a session key from the local storage if one exists and is not expired
// If one does not exist or is expired, it prints out an instruction and returns false
func GetValidSession(ctx infra.DnoteCtx) (string, bool, error) {
	db := ctx.DB

	var sessionKey string
	var sessionKeyExpires int64

	if err := GetSystem(db, infra.SystemSessionKey, &sessionKey); err != nil {
		return "", false, errors.Wrap(err, "getting session key")
	}
	if err := GetSystem(db, infra.SystemSessionKeyExpiry, &sessionKeyExpires); err != nil {
		return "", false, errors.Wrap(err, "getting session key expiry")
	}

	if sessionKey == "" {
		log.Error("login required. please run `dnote login`\n")
		return "", false, nil
	}
	if sessionKeyExpires < time.Now().Unix() {
		log.Error("sesison expired. please run `dnote login`\n")
		return "", false, nil
	}

	return sessionKey, true, nil
}
