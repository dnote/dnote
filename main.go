package main

import (
	"os"

	"github.com/dnote-io/cli/cmd/root"

	// commands
	_ "github.com/dnote-io/cli/cmd/books"
	_ "github.com/dnote-io/cli/cmd/delete"
	_ "github.com/dnote-io/cli/cmd/edit"
	_ "github.com/dnote-io/cli/cmd/login"
	_ "github.com/dnote-io/cli/cmd/new"
	_ "github.com/dnote-io/cli/cmd/notes"
	_ "github.com/dnote-io/cli/cmd/sync"
	_ "github.com/dnote-io/cli/cmd/use"
	_ "github.com/dnote-io/cli/cmd/version"
)

func init() {
	err := root.Prepare()
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
