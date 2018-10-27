package version

import (
	"fmt"

	"github.com/dnote/cli/infra"
	"github.com/spf13/cobra"
)

// NewCmd returns a new version command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Dnote",
		Long:  "Print the version number of Dnote",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dnote %s\n", ctx.Version)
		},
	}

	return cmd
}
