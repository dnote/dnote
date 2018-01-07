package use

import (
	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/log"
	"github.com/spf13/cobra"
)

var example = `
  dnote use JS`

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "use",
		Short:   "Change the current book",
		Aliases: []string{"u"},
		Example: example,
		RunE:    newRun(ctx),
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		targetBookName := args[0]

		err := core.ChangeBook(ctx, targetBookName)
		if err != nil {
			return err
		}

		log.Infof("now on book '%s'\n", targetBookName)
		return nil
	}
}
