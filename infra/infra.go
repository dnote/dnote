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
	Version = "1.0.0"

	// TimestampFilename is the name of the file containing upgrade info
	TimestampFilename = "timestamps"
	// TimestampFilename is the name of the directory containing dnote files
	DnoteDirName   = ".dnote"
	configFilename = "dnoterc"
	dnoteFilename  = "dnote"
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
	UID   string
	Notes []Note
}

// Note represents a single microlesson
type Note struct {
	UID     string
	Content string
	Dirty   bool
	AddedOn int64
}

// DnoteCtx is a context holding the information of the current runtime
type DnoteCtx struct {
	HomeDir  string
	DnoteDir string
}

type RunEFunc func(*cobra.Command, []string) error

// GetConfigPath returns the path to the dnote config file
func GetConfigPath(ctx DnoteCtx) (string, error) {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, configFilename), nil
}

// GetDnotePath returns the path to the dnote file
func GetDnotePath(ctx DnoteCtx) (string, error) {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, dnoteFilename), nil
}

// GetTimestampPath returns the path to the file containing dnote upgrade
// information
func GetTimestampPath(ctx DnoteCtx) (string, error) {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, TimestampFilename), nil
}

func InitConfigFile(ctx DnoteCtx) error {
	content := []byte("book: general\n")
	path, err := GetConfigPath(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get config path")
	}

	if utils.FileExists(path) {
		return nil
	}

	err = ioutil.WriteFile(path, content, 0644)
	return errors.Wrapf(err, "Failed to write the config file at '%s'", path)
}

// InitDnoteDir initializes dnote directory
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
	path, err := GetDnotePath(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get dnote path")
	}

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
	path, err := GetTimestampPath(ctx)
	if err != nil {
		return err
	}

	if utils.FileExists(path) {
		return nil
	}

	epoch := strconv.FormatInt(time.Now().Unix(), 10)
	content := []byte(fmt.Sprintf("LAST_UPGRADE_EPOCH: %s\n", epoch))

	err = ioutil.WriteFile(path, content, 0644)
	return err
}

// ReadNoteContent reads the content of dnote
func ReadNoteContent(ctx DnoteCtx) ([]byte, error) {
	notePath, err := GetDnotePath(ctx)
	if err != nil {
		return nil, err
	}

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

	notePath, err := GetDnotePath(ctx)
	if err != nil {
		return err
	}

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

	configPath, err := GetConfigPath(ctx)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfig(ctx DnoteCtx) (Config, error) {
	var ret Config

	configPath, err := GetConfigPath(ctx)
	if err != nil {
		return ret, err
	}

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
		UID:     utils.GenerateUID(),
		Content: content,
		AddedOn: time.Now().Unix(),
	}
}

// MakeBook returns a book
func MakeBook() Book {
	return Book{
		UID:   utils.GenerateUID(),
		Notes: make([]Note, 0),
	}
}

func GetUpdatedBook(book Book, notes []Note) Book {
	b := MakeBook()

	b.UID = book.UID
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
