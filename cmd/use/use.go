package use

import (
	"fmt"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/dnote-io/cli/infra"
	"github.com/spf13/cobra"
)

var example = `
  dnote use JS`

var cmd = &cobra.Command{
	Use:     "use",
	Short:   "Change the current book",
	Aliases: []string{"u"},
	Example: example,
	RunE:    run,
}

func init() {
	root.Register(cmd)
}

func run(cmd *cobra.Command, args []string) error {
	targetBookName := args[0]

	err := infra.ChangeBook(targetBookName)
	if err != nil {
		return err
	}

	fmt.Printf("Now using %s\n", targetBookName)
	return nil
}
