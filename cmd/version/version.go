package cmd

import (
	"fmt"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/dnote-io/cli/infra"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dnote",
	Long:  "Print the version number of Dnote",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dnote v%s\n", infra.Version)
	},
}

func init() {
	root.Register(cmd)
}
