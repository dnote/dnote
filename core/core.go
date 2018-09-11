package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/dnote/actions"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	// TimestampFilename is the name of the file containing upgrade info
	TimestampFilename  = "timestamps"
	ConfigFilename     = "dnoterc"
	DnoteFilename      = "dnote"
	ActionFilename     = "actions"
	TmpContentFilename = "DNOTE_TMPCONTENT"
)

type RunEFunc func(*cobra.Command, []string) error

// GetConfigPath returns the path to the dnote config file
func GetConfigPath(ctx infra.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, ConfigFilename)
}

// GetDnotePath returns the path to the dnote file
func GetDnotePath(ctx infra.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, DnoteFilename)
}

// GetTimestampPath returns the path to the file containing dnote upgrade
// information
func GetTimestampPath(ctx infra.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, TimestampFilename)
}

// GetActionPath returns the path to the file containing user actions
func GetActionPath(ctx infra.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, ActionFilename)
}

// GetDnoteTmpContentPath returns the path to the temporary file containing
// content being added or edited
func GetDnoteTmpContentPath(ctx infra.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, TmpContentFilename)
}

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

func getEditorCommand() string {
	editor := os.Getenv("EDITOR")

	switch editor {
	case "atom":
		return "atom -w"
	case "subl":
		return "subl -n -w"
	case "mate":
		return "mate -w"
	case "vim":
		return "vim"
	case "nvim":
		return "nvim"
	case "nano":
		return "nano"
	case "emacs":
		return "emacs"
	default:
		return "vi"
	}
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
		return errors.Wrap(err, "Failed to marshal config into YAML")
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return errors.Wrapf(err, "Failed to write the config file at '%s'", path)
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

// TODO: delete
// ReadTimestamp gets the content of the timestamp file
func ReadTimestamp(ctx infra.DnoteCtx) (infra.Timestamp, error) {
	var ret infra.Timestamp

	path := GetTimestampPath(ctx)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ret, err
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to unmarshal timestamp content")
	}

	return ret, nil
}

func WriteTimestamp(ctx infra.DnoteCtx, timestamp infra.Timestamp) error {
	d, err := yaml.Marshal(timestamp)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal timestamp into YAML")
	}

	path := GetTimestampPath(ctx)
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to write timestamp to the file")
	}

	return nil
}

// ReadNoteContent reads the content of dnote
func ReadNoteContent(ctx infra.DnoteCtx) ([]byte, error) {
	notePath := GetDnotePath(ctx)

	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetDnote reads and parses the dnote
func GetDnote(ctx infra.DnoteCtx) (infra.Dnote, error) {
	ret := infra.Dnote{}

	b, err := ReadNoteContent(ctx)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to read note content")
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to unmarshal note content")
	}

	return ret, nil
}

// WriteDnote persists the state of Dnote into the dnote file
func WriteDnote(ctx infra.DnoteCtx, dnote infra.Dnote) error {
	d, err := json.MarshalIndent(dnote, "", "  ")
	if err != nil {
		return err
	}

	notePath := GetDnotePath(ctx)

	err = ioutil.WriteFile(notePath, d, 0644)
	if err != nil {
		errors.Wrap(err, "Failed to write to the dnote file")
	}

	return nil
}

func WriteConfig(ctx infra.DnoteCtx, config infra.Config) error {
	d, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	configPath := GetConfigPath(ctx)

	err = ioutil.WriteFile(configPath, d, 0644)
	if err != nil {
		errors.Wrap(err, "Failed to write to the config file")
	}

	return nil
}

// LogAction logs action
func LogAction(tx *sql.Tx, schema int, actionType, data string, timestamp int64) error {
	uuid := uuid.NewV4().String()

	_, err := tx.Exec(`INSERT INTO actions (uuid, schema, type, data, timestamp)
	VALUES (?, ?, ?, ?, ?)`, uuid, schema, actionType, data, timestamp)
	if err != nil {
		return errors.Wrap(err, "inserting an action")
	}

	return nil
}

func WriteActionLog(ctx infra.DnoteCtx, ats []actions.Action) error {
	path := GetActionPath(ctx)

	d, err := json.Marshal(ats)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal newly generated actions to JSON")
	}

	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to write to the actions file")
	}

	return nil
}

func ClearActionLog(ctx infra.DnoteCtx) error {
	var content []actions.Action

	if err := WriteActionLog(ctx, content); err != nil {
		return errors.Wrap(err, "Failed to write action log")
	}

	return nil
}

func ReadActionLogContent(ctx infra.DnoteCtx) ([]byte, error) {
	path := GetActionPath(ctx)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, errors.Wrap(err, "Failed to read the action file")
	}

	return b, nil
}

// ReadActionLog returns the action log content
func ReadActionLog(ctx infra.DnoteCtx) ([]actions.Action, error) {
	var ret []actions.Action

	b, err := ReadActionLogContent(ctx)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to read the action log content")
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to unmarshal action log JSON")
	}

	return ret, nil
}

func ReadConfig(ctx infra.DnoteCtx) (infra.Config, error) {
	var ret infra.Config

	configPath := GetConfigPath(ctx)
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return ret, err
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to unmarshal config YAML")
	}

	return ret, nil
}

func UpdateLastActionTimestamp(ctx infra.DnoteCtx, val int64) error {
	ts, err := ReadTimestamp(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to read the timestamp file")
	}

	ts.LastAction = val

	err = WriteTimestamp(ctx, ts)
	if err != nil {
		return errors.Wrap(err, "Failed to write the timestamp to the file")
	}

	return nil
}

// NewNote returns a note
func NewNote(content string, ts int64) infra.Note {
	return infra.Note{
		UUID:    utils.GenerateUUID(),
		Content: content,
		AddedOn: ts,
	}
}

// NewBook returns a book
func NewBook(name string) infra.Book {
	return infra.Book{
		Name:  name,
		Notes: make([]infra.Note, 0),
	}
}

func GetUpdatedBook(book infra.Book, notes []infra.Note) infra.Book {
	b := NewBook(book.Name)

	b.Notes = notes

	return b
}

func FilterNotes(notes []infra.Note, testFunc func(infra.Note) bool) []infra.Note {
	var ret []infra.Note

	for _, note := range notes {
		if testFunc(note) {
			ret = append(ret, note)
		}
	}

	return ret
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

func getEditorCmd(ctx infra.DnoteCtx, fpath string) (*exec.Cmd, error) {
	config, err := ReadConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read the config")
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
			return errors.Wrap(err, "Failed to create a temporary file for content")
		}
		err = f.Close()
		if err != nil {
			return errors.Wrap(err, "Failed to close the temporary file for content")
		}
	}

	cmd, err := getEditorCmd(ctx, fpath)
	if err != nil {
		return errors.Wrap(err, "Failed to create the editor command")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "Failed to launch the editor")
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Wrap(err, "Failed to wait for the editor")
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return errors.Wrap(err, "Failed to read the file")
	}

	err = os.Remove(fpath)
	if err != nil {
		return errors.Wrap(err, "Failed to remove the temporary content file")
	}

	raw := string(b)
	c := SanitizeContent(raw)

	*content = c

	return nil
}
