package ls

import (
	"fmt"

	"github.com/dnote-io/cli/infra"
	"github.com/spf13/cobra"
)

var example = `
 * List notes in the current book
 dnote notes
 dnote ls

 * List notes in a certain book
 dnote ls javascript
 `

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls <book name?>",
		Aliases: []string{"notes"},
		Short:   "List all notes",
		Example: example,
		RunE:    newRun(ctx),
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var bookName string

		if len(args) == 1 {
			bookName = args[0]
		} else {
			var err error
			bookName, err = infra.GetCurrentBook(ctx)
			if err != nil {
				return err
			}
		}

		fmt.Printf("On note %s\n", bookName)

		dnote, err := infra.GetDnote(ctx)
		if err != nil {
			return err
		}

		for k, v := range dnote {
			if k == bookName {
				for i, note := range v.Notes {
					fmt.Printf("* [%d] - %s\n", i, note.Content)
				}
			}
		}

		return nil
	}

}
