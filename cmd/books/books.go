package books

import (
	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/log"
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
		books, err := core.GetBookNames(ctx)
		if err != nil {
			return err
		}

		for _, book := range books {
			log.Printf("  %s%s\n", "  ", book)
		}

		return nil
	}
}
