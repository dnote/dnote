package infra

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/dnote-io/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	// Version is the current version of dnote
	Version = "0.2.0"

	// TimestampFilename is the name of the file containing upgrade info
	TimestampFilename = "timestamps"
	// DnoteDirName is the name of the directory containing dnote files
	DnoteDirName   = ".dnote"
	ConfigFilename = "dnoterc"
	DnoteFilename  = "dnote"
	ActionFilename = "actions"
)

// Config holds dnote configuration
type Config struct {
	Book   string
	APIKey string
}

// Dnote holds the whole dnote data
type Dnote map[string]Book

// Book holds a metadata and its notes
type Book struct {
	UUID  string
	Notes []Note
}

// Note represents a single microlesson
type Note struct {
	UUID    string
	Content string
}

// DnoteCtx is a context holding the information of the current runtime
type DnoteCtx struct {
	HomeDir  string
	DnoteDir string
}

type RunEFunc func(*cobra.Command, []string) error

// GetConfigPath returns the path to the dnote config file
func GetConfigPath(ctx DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, ConfigFilename)
}

// GetDnotePath returns the path to the dnote file
func GetDnotePath(ctx DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, DnoteFilename)
}

// GetTimestampPath returns the path to the file containing dnote upgrade
// information
func GetTimestampPath(ctx DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, TimestampFilename)
}

// GetActionPath returns the path to the file containing user actions
func GetActionPath(ctx DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, ActionFilename)
}

// InitActionFile populates action file if it does not exist
func InitActionFile(ctx DnoteCtx) error {
	path := GetActionPath(ctx)

	if utils.FileExists(path) {
		return nil
	}

	b, err := json.Marshal(&[]Action{})
	if err != nil {
		return errors.Wrap(err, "Failed to get initial action content")
	}

	err = ioutil.WriteFile(path, b, 0644)
	return err
}

// InitConfigFile populates a new config file if it does not exist yet
func InitConfigFile(ctx DnoteCtx) error {
	content := []byte("book: general\n")
	path := GetConfigPath(ctx)

	if utils.FileExists(path) {
		return nil
	}

	err := ioutil.WriteFile(path, content, 0644)
	return errors.Wrapf(err, "Failed to write the config file at '%s'", path)
}

// InitDnoteDir initializes dnote directory if it does not exist yet
func InitDnoteDir(ctx DnoteCtx) error {
	path := ctx.DnoteDir

	if utils.FileExists(path) {
		return nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return errors.Wrap(err, "Failed to create dnote directory")
	}

	return nil
}

// InitDnoteFile creates an empty dnote file
func InitDnoteFile(ctx DnoteCtx) error {
	path := GetDnotePath(ctx)

	if utils.FileExists(path) {
		return nil
	}

	b, err := json.Marshal(&Dnote{})
	if err != nil {
		return errors.Wrap(err, "Failed to get initial dnote content")
	}

	err = ioutil.WriteFile(path, b, 0644)
	return err
}

// InitTimestampFile creates an empty dnote upgrade file
func InitTimestampFile(ctx DnoteCtx) error {
	path := GetTimestampPath(ctx)

	if utils.FileExists(path) {
		return nil
	}

	epoch := strconv.FormatInt(time.Now().Unix(), 10)
	content := []byte(fmt.Sprintf("LAST_UPGRADE_EPOCH: %s\n", epoch))

	err := ioutil.WriteFile(path, content, 0644)
	return err
}

// ReadNoteContent reads the content of dnote
func ReadNoteContent(ctx DnoteCtx) ([]byte, error) {
	notePath := GetDnotePath(ctx)

	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetDnote reads and parses the dnote
func GetDnote(ctx DnoteCtx) (Dnote, error) {
	ret := Dnote{}

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
func WriteDnote(ctx DnoteCtx, dnote Dnote) error {
	d, err := json.MarshalIndent(dnote, "", "  ")
	if err != nil {
		return err
	}

	notePath := GetDnotePath(ctx)

	err = ioutil.WriteFile(notePath, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func WriteConfig(ctx DnoteCtx, config Config) error {
	d, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	configPath := GetConfigPath(ctx)

	err = ioutil.WriteFile(configPath, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LogAction appends the action to the action log
func LogAction(ctx DnoteCtx, a Action) error {
	actions, err := ReadActionLog(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to read the action log")
	}

	actions = append(actions, a)
	d, err := json.Marshal(actions)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal newly generated actions to JSON")
	}

	path := GetActionPath(ctx)
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadActionLogContent(ctx DnoteCtx) ([]byte, error) {
	path := GetActionPath(ctx)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, errors.Wrap(err, "Failed to read the action file")
	}

	return b, nil
}

// ReadActionLog returns the action log content
func ReadActionLog(ctx DnoteCtx) ([]Action, error) {
	var ret []Action

	b, err := ReadActionLogContent(ctx)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to read the action log content")
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ReadConfig(ctx DnoteCtx) (Config, error) {
	var ret Config

	configPath := GetConfigPath(ctx)
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return ret, err
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func GetCurrentBook(ctx DnoteCtx) (string, error) {
	config, err := ReadConfig(ctx)
	if err != nil {
		return "", err
	}

	return config.Book, nil
}

func GetBooks(ctx DnoteCtx) ([]string, error) {
	dnote, err := GetDnote(ctx)
	if err != nil {
		return nil, err
	}

	books := make([]string, 0, len(dnote))
	for k := range dnote {
		books = append(books, k)
	}

	sort.Strings(books)

	return books, nil
}

// ChangeBook replaces the book name in the dnote config file
func ChangeBook(ctx DnoteCtx, bookName string) error {
	config, err := ReadConfig(ctx)
	if err != nil {
		return err
	}

	config.Book = bookName

	err = WriteConfig(ctx, config)
	if err != nil {
		return err
	}

	// Now add this book to the .dnote file, for issue #2
	dnote, err := GetDnote(ctx)
	if err != nil {
		return err
	}

	_, exists := dnote[bookName]
	if !exists {
		dnote[bookName] = MakeBook()
		err := WriteDnote(ctx, dnote)
		if err != nil {
			return err
		}
	}

	return nil
}

// MakeNote returns a note
func MakeNote(content string) Note {
	return Note{
		UUID:    utils.GenerateUID(),
		Content: content,
	}
}

// MakeBook returns a book
func MakeBook() Book {
	return Book{
		UUID:  utils.GenerateUID(),
		Notes: make([]Note, 0),
	}
}

func GetUpdatedBook(book Book, notes []Note) Book {
	b := MakeBook()

	b.UUID = book.UUID
	b.Notes = notes

	return b
}

// MigrateToDnoteDir creates dnote directory if artifacts from the previous version
// of dnote are present, and moves the artifacts to the directory.
func MigrateToDnoteDir(ctx DnoteCtx) error {
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

	if err := os.Rename(oldDnotePath, fmt.Sprintf("%s/dnote", temporaryDirPath)); err != nil {
		return errors.Wrap(err, "Failed to move .dnote file")
	}
	if err := os.Rename(oldDnotercPath, fmt.Sprintf("%s/dnoterc", temporaryDirPath)); err != nil {
		return errors.Wrap(err, "Failed to move .dnoterc file")
	}
	if err := os.Rename(oldDnoteUpgradePath, fmt.Sprintf("%s/timestamps", temporaryDirPath)); err != nil {
		return errors.Wrap(err, "Failed to move .dnote-upgrade file")
	}

	// Now that all files are moved to the temporary dir, rename the dir to .dnote
	if err := os.Rename(temporaryDirPath, fmt.Sprintf("%s/.dnote", homeDir)); err != nil {
		return errors.Wrap(err, "Failed to rename temporary dir to .dnote")
	}

	return nil
}

// IsFreshInstall checks if the dnote files have been initialized
func IsFreshInstall(ctx DnoteCtx) (bool, error) {
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
