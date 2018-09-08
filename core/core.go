package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dnote/actions"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/migrate"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	// TimestampFilename is the name of the file containing upgrade info
	TimestampFilename = "timestamps"
	// DnoteDirName is the name of the directory containing dnote files
	DnoteDirName       = ".dnote"
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

// initActionFile populates action file if it does not exist
func initActionFile(ctx infra.DnoteCtx) error {
	path := GetActionPath(ctx)

	if utils.FileExists(path) {
		return nil
	}

	b, err := json.Marshal(&[]actions.Action{})
	if err != nil {
		return errors.Wrap(err, "Failed to get initial action content")
	}

	err = ioutil.WriteFile(path, b, 0644)
	return err
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
	fresh, err := isFreshInstall(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to check if fresh install")
	}

	err = initDnoteDir(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create dnote dir")
	}
	err = initConfigFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to generate config file")
	}
	err = initDnoteFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create dnote file")
	}
	err = initTimestampFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create dnote upgrade file")
	}
	err = initActionFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create action file")
	}
	err = migrate.InitSchemaFile(ctx, fresh)
	if err != nil {
		return errors.Wrap(err, "Failed to create migration file")
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

// initDnoteFile creates an empty dnote file
func initDnoteFile(ctx infra.DnoteCtx) error {
	path := GetDnotePath(ctx)

	if utils.FileExists(path) {
		return nil
	}

	b, err := json.Marshal(&infra.Dnote{})
	if err != nil {
		return errors.Wrap(err, "Failed to get initial dnote content")
	}

	err = ioutil.WriteFile(path, b, 0644)
	return err
}

// initTimestampFile creates an empty dnote upgrade file
func initTimestampFile(ctx infra.DnoteCtx) error {
	path := GetTimestampPath(ctx)

	if utils.FileExists(path) {
		return nil
	}

	now := time.Now().Unix()
	ts := infra.Timestamp{
		LastUpgrade: now,
	}

	b, err := yaml.Marshal(&ts)
	if err != nil {
		return errors.Wrap(err, "Failed to get initial timestamp content")
	}

	err = ioutil.WriteFile(path, b, 0644)
	return err
}

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

// LogAction appends the action to the action log and updates the last_action
// timestamp
func LogAction(ctx infra.DnoteCtx, action actions.Action) error {
	actions, err := ReadActionLog(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to read the action log")
	}

	actions = append(actions, action)

	err = WriteActionLog(ctx, actions)
	if err != nil {
		return errors.Wrap(err, "Failed to write action log")
	}

	err = UpdateLastActionTimestamp(ctx, action.Timestamp)
	if err != nil {
		return errors.Wrap(err, "Failed to update the last_action timestamp")
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
		UUID:    utils.GenerateUID(),
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

// MigrateToDnoteDir creates dnote directory if artifacts from the previous version
// of dnote are present, and moves the artifacts to the directory.
func MigrateToDnoteDir(ctx infra.DnoteCtx) error {
	homeDir := ctx.HomeDir

	temporaryDirPath := fmt.Sprintf("%s/.dnote-tmp", homeDir)
	oldDnotePath := fmt.Sprintf("%s/.dnote", homeDir)
	oldDnotercPath := fmt.Sprintf("%s/.dnoterc", homeDir)
	oldDnoteUpgradePath := fmt.Sprintf("%s/.dnote-upgrade", homeDir)

	// Check if a dnote file exists. Return early if it does not exist,
	// or exists but already a directory.
	fi, err := os.Stat(oldDnotePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return errors.Wrap(err, "Failed to look up old dnote path")
	}
	if fi.IsDir() {
		return nil
	}

	if err := os.Mkdir(temporaryDirPath, 0755); err != nil {
		return errors.Wrap(err, "Failed to make temporary .dnote directory")
	}

	// In the beta release for v0.2, backup user's .dnote
	if err := utils.CopyFile(oldDnotePath, fmt.Sprintf("%s/dnote-bak-5cdde2e83", homeDir)); err != nil {
		return errors.Wrap(err, "Failed to back up the old .dnote file")
	}

	if err := os.Rename(oldDnotePath, fmt.Sprintf("%s/dnote", temporaryDirPath)); err != nil {
		return errors.Wrap(err, "Failed to move .dnote file")
	}
	if err := os.Rename(oldDnotercPath, fmt.Sprintf("%s/dnoterc", temporaryDirPath)); err != nil {
		return errors.Wrap(err, "Failed to move .dnoterc file")
	}
	if err := os.Remove(oldDnoteUpgradePath); err != nil {
		return errors.Wrap(err, "Failed to delete the old upgrade file")
	}

	// Now that all files are moved to the temporary dir, rename the dir to .dnote
	if err := os.Rename(temporaryDirPath, fmt.Sprintf("%s/.dnote", homeDir)); err != nil {
		return errors.Wrap(err, "Failed to rename temporary dir to .dnote")
	}

	return nil
}

// isFreshInstall checks if the dnote files have been initialized
func isFreshInstall(ctx infra.DnoteCtx) (bool, error) {
	path := ctx.DnoteDir

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "Failed to get file info for dnote directory")
	}

	return false, nil
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
