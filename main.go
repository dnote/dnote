package main

import (
	"fmt"
	"os"

	// For testing purposes.
	"./cmd/books"
	"./cmd/login"
	"./cmd/new"
	"./cmd/edit"
	"./cmd/delete"
	"./cmd/notes"
	"./cmd/sync"
	"./upgrade"
	"./utils"

	// For GitHub.
	/*
	"github.com/dnote-io/cli/cmd/books"
	"github.com/dnote-io/cli/cmd/delete"
	"github.com/dnote-io/cli/cmd/edit"
	"github.com/dnote-io/cli/cmd/login"
	"github.com/dnote-io/cli/cmd/new"
	"github.com/dnote-io/cli/cmd/notes"
	"github.com/dnote-io/cli/cmd/sync"
	"github.com/dnote-io/cli/upgrade"
	"github.com/dnote-io/cli/utils"
	*/
)

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

	err = upgrade.Migrate()
	if err != nil {
		return err
	}

	return nil

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// changeBook replaces the book name in the dnote config file
func changeBook(bookName string) error {
	config, err := utils.ReadConfig()
	if err != nil {
		return err
	}

	config.Book = bookName

	err = utils.WriteConfig(config)
	if err != nil {
		return err
	}

	// Now add this book to the .dnote file, for issue #2
	dnote, err := utils.GetDnote()
	if err != nil {
		return err
	}

	_, exists := dnote[bookName]
	if exists == false {
		dnote[bookName] = make([]utils.Note, 0)
		err := utils.WriteDnote(dnote)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Now using %s\n", bookName)

	return nil
}

func checkFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func init() {
	err := initDnote()
	check(err)
}

func main() {
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
		var notename string
		var note string
		var err error

		if os.Args[2] != "-t" {
			notename, err = utils.GenerateNoteName()
			note = os.Args[2]
			check(err)
		} else if os.Args[2] == "-t" {
			notename = os.Args[3]
			note = os.Args[4]
		}

		err = new.Run(notename, note)
		check(err)
	case "edit", "e":
		notename := os.Args[2]
		newcontent := os.Args[3]
		err := edit.Edit(notename, newcontent)
		check(err)
	case "delete", "d":
		if os.Args[2] == "-b" {
			err := delete.DeleteBook(os.Args[3])
			check(err)
		} else if os.Args[2] == "-n" {
			err := delete.DeleteNote(os.Args[3])
			check(err)
		}
	case "books", "b":
		err := books.Run()
		check(err)
	case "upgrade":
		err := upgrade.Upgrade()
		check(err)
	case "--version":
		fmt.Println(utils.Version)
	case "notes":
		err := notes.Run()
		check(err)
	case "sync":
		err := sync.Sync()
		check(err)
	case "login":
		err := login.Run()
		check(err)
	default:
		break
	}

	err := upgrade.AutoUpgrade()
	if err != nil {
		fmt.Println("Warning - Failed to check for update:", err)
	}
}
