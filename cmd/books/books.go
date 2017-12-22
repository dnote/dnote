package books

import (
	"fmt"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/dnote-io/cli/infra"
	"github.com/spf13/cobra"
)

var example = `
 dnote books`

var cmd = &cobra.Command{
	Use:     "books",
	Short:   "List books",
	Aliases: []string{"b"},
	Example: example,
	RunE:    run,
}

func init() {
	root.Register(cmd)
}

func run(cmd *cobra.Command, args []string) error {
	currentBook, err := infra.GetCurrentBook()
	if err != nil {
		return err
	}

	books, err := infra.GetBooks()
	if err != nil {
		return err
	}

	for _, book := range books {
		if book == currentBook {
			fmt.Printf("* %v\n", book)
		} else {
			fmt.Printf("  %v\n", book)
		}
	}

	return nil
}
