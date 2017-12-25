package books

import (
	"fmt"

	"github.com/dnote-io/cli/infra"
	"github.com/spf13/cobra"
)

var example = `
 dnote books`

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "books",
		Short:   "List books",
		Aliases: []string{"b"},
		Example: example,
		RunE:    newRun(ctx),
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		currentBook, err := infra.GetCurrentBook(ctx)
		if err != nil {
			return err
		}

		books, err := infra.GetBooks(ctx)
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

}
