package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/user"
	"sort"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Book   string
	APIKey string
}

// Deprecated. See upgrade/migrate.go
type YAMLDnote map[string][]string

type Dnote map[string]Book
type Book []Note
type Note struct {
	UID     string
	Content string
	AddedOn int64
}

const configFilename = ".dnoterc"
const DnoteUpdateFilename = ".dnote-upgrade"
const dnoteFilename = ".dnote"
const Version = "0.0.3"

const letterRunes = "abcdefghipqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateNoteID() string {
	result := make([]byte, 8)
	for i := range result {
		result[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(result)
}

func GetConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, configFilename), nil
}

func GetDnotePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, dnoteFilename), nil
}

func GetYAMLDnoteArchivePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, ".dnote-yaml-archived"), nil
}

func GenerateConfigFile() error {
	content := []byte("book: general\n")
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, content, 0644)
	return err
}

func TouchDnoteFile() error {
	dnotePath, err := GetDnotePath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dnotePath, []byte{}, 0644)
	return err
}

func TouchDnoteUpgradeFile() error {
	dnoteUpdatePath, err := GetDnoteUpdatePath()
	if err != nil {
		return err
	}

	epoch := strconv.FormatInt(time.Now().Unix(), 10)
	content := []byte(fmt.Sprintf("LAST_UPGRADE_EPOCH: %s\n", epoch))

	err = ioutil.WriteFile(dnoteUpdatePath, content, 0644)
	return err
}

func GetDnoteUpdatePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, DnoteUpdateFilename), nil
}

func AskConfirmation(question string) (bool, error) {
	fmt.Printf("%s [Y/n]: ", question)

	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	ok := res == "y\n" || res == "Y\n" || res == "\n"

	return ok, nil
}

// ReadNoteContent reads the content of dnote
func ReadNoteContent() ([]byte, error) {
	notePath, err := GetDnotePath()
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
func GetDnote() (Dnote, error) {
	ret := Dnote{}

	b, err := ReadNoteContent()
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// Deprecated. See upgrade/upgrade.go
func GetNote() (YAMLDnote, error) {
	ret := YAMLDnote{}

	b, err := ReadNoteContent()
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
func WriteConfig(config Config) error {
	d, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfig() (Config, error) {
	var ret Config

	configPath, err := GetConfigPath()
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

func GetCurrentBook() (string, error) {
	config, err := ReadConfig()
	if err != nil {
		return "", err
	}

	return config.Book, nil
}

func GetBooks() ([]string, error) {
	dnote, err := GetDnote()
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
