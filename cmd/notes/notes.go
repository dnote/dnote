package notes

import (
	"fmt"

	"github.com/dnote-io/cli/cmd/root"
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

var cmd = &cobra.Command{
	Use:     "notes <book name?>",
	Aliases: []string{"ls"},
	Short:   "List all notes",
	Example: example,
	RunE:    run,
}

func init() {
	root.Register(cmd)
}

func run(cmd *cobra.Command, args []string) error {
	var bookName string

	if len(args) == 1 {
		bookName = args[0]
	} else {
		var err error
		bookName, err = infra.GetCurrentBook()
		if err != nil {
			return err
		}
	}

	fmt.Printf("On note %s\n", bookName)

	dnote, err := infra.GetDnote()
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
