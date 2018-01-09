package ls

import (
	"fmt"
	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
 * List notes in the current book
 dnote notes
 dnote ls

 * List notes in a certain book
 dnote ls javascript
 `

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls <book name?>",
		Aliases: []string{"notes"},
		Short:   "List all notes",
		Example: example,
		RunE:    newRun(ctx),
		PreRunE: preRun,
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		bookName := args[0]

		log.Infof("on book %s\n", bookName)

		dnote, err := core.GetDnote(ctx)
		if err != nil {
			return err
		}

		book := dnote[bookName]

		for i, note := range book.Notes {
			fmt.Printf("  \033[%dm(%d)\033[0m %s\n", log.ColorYellow, i, note.Content)
		}

		return nil
	}

}
