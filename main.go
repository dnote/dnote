package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/dnote-io/cli/upgrade"
	"github.com/dnote-io/cli/utils"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Book string
}

type Note map[string][]string

// initDnote creates a config file if one does not exist
func initDnote() error {
	configPath, err := utils.GetConfigPath()
	if err != nil {
		return err
	}
	dnotePath, err := utils.GetDnotePath()
	if err != nil {
		return err
	}
	dnoteUpdatePath, err := utils.GetDnoteUpdatePath()
	if err != nil {
		return err
	}

	if !checkFileExists(configPath) {
		err := utils.GenerateConfigFile()
		if err != nil {
			return err
		}
	}
	if !checkFileExists(dnotePath) {
		err := utils.TouchDnoteFile()
		if err != nil {
			return err
		}
	}
	if !checkFileExists(dnoteUpdatePath) {
		err := utils.TouchDnoteUpgradeFile()
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

	configPath, err := utils.GetConfigPath()
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

	configPath, err := utils.GetConfigPath()
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

	notePath, err := utils.GetDnotePath()
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

	notePath, err := utils.GetDnotePath()
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

func getNotesInBook(bookName string) ([]string, error) {
	note, err := readNote()
	if err != nil {
		return nil, err
	}

	notes := make([]string, 0, len(note))
	for k, v := range note {
		if k == bookName {
			for _, noteContent := range v {
				notes = append(notes, noteContent)
			}
		}
	}

	sort.Strings(notes)

	return notes, nil
}

func getNotesInCurrentBook() ([]string, error) {
	currentBook, err := getCurrentBook()
	if err != nil {
		return nil, err
	}

	return getNotesInBook(currentBook)
}

func checkFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func main() {
	err := initDnote()
	check(err)

	if len(os.Args) < 2 {
		fmt.Println("Dnote - Spontaneously capture new engineering lessons\n")
		fmt.Println("Main commands:")
		fmt.Println("  use [u] - choose the book")
		fmt.Println("  new [n] - write a new note")
		fmt.Println("  books [b] - show books")
		fmt.Println("  notes - show notes for book")
		fmt.Println("")
		fmt.Println("Other commands:")
		fmt.Println("  upgrade - upgrade dnote")
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
		currentBook, err := getCurrentBook()
		check(err)
		fmt.Printf("[+] Added to: %s\n", currentBook)
		err = writeNote(note)
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
	case "upgrade":
		err := upgrade.Upgrade()
		check(err)
	case "--version":
		fmt.Println(utils.Version)
	case "notes":
		defaultBookName, err := getCurrentBook()

		check(err)

		var bookName string

		if len(os.Args) == 2 {
			bookName = defaultBookName
		} else if len(os.Args) == 4 && os.Args[2] == "-b" {
			bookName = os.Args[3]
		} else {
			fmt.Println("Invalid argument passed to notes")
			os.Exit(1)
		}

		notes, err := getNotesInBook(bookName)
		check(err)

		fmt.Printf("Notes in book %s:\n", bookName)

		for _, note := range notes {
			fmt.Printf("%s\n", note)
		}
	default:
		break
	}

	err = upgrade.AutoUpgrade()
	if err != nil {
		fmt.Println("Warning - Failed to check for update:", err)
	}
}
