package login

import (
	"fmt"

	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/log"
	"github.com/spf13/cobra"
)

var example = `
  dnote login`

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to dnote server",
		Example: example,
		RunE:    newRun(ctx),
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		log.Plain("Welcome to Dnote Cloud :)\n\n")
		log.Plain("A home for your engineer microlessons\n")
		log.Plain("You can register at https://dnote.io\n\n")
		log.Printf("API key: ")

		var apiKey string
		fmt.Scanln(&apiKey)

		config, err := core.ReadConfig(ctx)
		if err != nil {
			return err
		}

		config.APIKey = apiKey
		err = core.WriteConfig(ctx, config)
		if err != nil {
			return err
		}

		log.Success("success\n")

		return nil
	}

}
