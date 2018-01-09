package books

import (
	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/log"
	"github.com/pkg/errors"
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

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dnote, err := core.GetDnote(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to read dnote")
		}

		for bookName, book := range dnote {
			log.Printf("%s \033[%dm(%d)\033[0m\n", bookName, log.ColorYellow, len(book.Notes))
		}

		return nil
	}
}
