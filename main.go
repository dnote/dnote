package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"sort"

	"github.com/dnote-io/cli/upgrade"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Book string
}

type Note map[string][]string

const configFilename = ".dnoterc"
const dnoteFilename = ".dnote"

func getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, configFilename), nil
}

func getDnotePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, dnoteFilename), nil
}

func generateConfigFile() error {
	content := []byte("book: general\n")
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, content, 0644)
	return err
}

func touchDnoteFile() error {
	dnotePath, err := getDnotePath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dnotePath, []byte{}, 0644)
	return err
}

// initDnote creates a config file if one does not exist
func initDnote() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	dnotePath, err := getDnotePath()
	if err != nil {
		return err
	}

	if !checkFileExists(configPath) {
		err := generateConfigFile()
		if err != nil {
			return err
		}
	}
	if !checkFileExists(dnotePath) {
		err := touchDnoteFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readConfig() (Config, error) {
	var ret Config

	configPath, err := getConfigPath()
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

func getCurrentBook() (string, error) {
	config, err := readConfig()
	if err != nil {
		return "", err
	}

	return config.Book, nil
}

func writeConfig(config Config) error {
	d, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

// changeBook replaces the book name in the dnote config file
func changeBook(bookName string) error {
	config, err := readConfig()
	if err != nil {
		return err
	}

	config.Book = bookName

	err = writeConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func readNote() (Note, error) {
	ret := Note{}

	notePath, err := getDnotePath()
	if err != nil {
		return ret, err
	}

	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return ret, nil
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func writeNote(content string) error {
	note, err := readNote()
	if err != nil {
		return err
	}

	book, err := getCurrentBook()
	if err != nil {
		return err
	}

	if _, ok := note[book]; ok {
		note[book] = append(note[book], content)
	} else {
		note[book] = []string{content}
	}

	d, err := yaml.Marshal(note)
	if err != nil {
		return err
	}

	notePath, err := getDnotePath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(notePath, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getBooks() ([]string, error) {
	note, err := readNote()
	if err != nil {
		return nil, err
	}

	books := make([]string, 0, len(note))
	for k := range note {
		books = append(books, k)
	}

	sort.Strings(books)

	return books, nil
}

func checkFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func main() {
	err := initDnote()
	check(err)
	err = upgrade.AutoUpdate()
	if err != nil {
		fmt.Println("Warning: Failed to check for update", err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Dnote - Spontaneously capture new engineering lessons\n")
		fmt.Println("Commands:")
		fmt.Println("  use [u] - choose the book")
		fmt.Println("  new [n] - write a new note")
		fmt.Println("  books [b] - show books")
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "use", "u":
		book := os.Args[2]
		err := changeBook(book)
		check(err)
	case "new", "n":
		note := os.Args[2]
		fmt.Println(note)
		err := writeNote(note)
		check(err)
	case "books", "b":
		currentBook, err := getCurrentBook()
		check(err)
		books, err := getBooks()
		check(err)

		for _, book := range books {
			if book == currentBook {
				fmt.Printf("* %v\n", book)
			} else {
				fmt.Printf("  %v\n", book)
			}
		}
	default:
		break
	}

}
