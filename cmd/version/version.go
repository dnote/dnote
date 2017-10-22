package cmd

import (
	"fmt"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/spf13/cobra"
)

// Version is the current version of dnote
const Version = "0.1.3"

var cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dnote",
	Long:  "Print the version number of Dnote",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dnote v%s\n", Version)
	},
}

func init() {
	root.Register(cmd)
}
