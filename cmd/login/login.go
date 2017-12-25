package login

import (
	"fmt"

	"github.com/dnote-io/cli/infra"
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

func newRun(ctx infra.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Print("Please enter your APIKey: ")

		var apiKey string
		fmt.Scanln(&apiKey)

		config, err := infra.ReadConfig(ctx)
		if err != nil {
			return err
		}

		config.APIKey = apiKey
		err = infra.WriteConfig(ctx, config)
		if err != nil {
			return err
		}

		return nil
	}

}
