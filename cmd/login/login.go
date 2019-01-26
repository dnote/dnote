package login

import (
	"fmt"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
  dnote login`

// NewCmd returns a new login command
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
		log.Plain("\n")
		log.Plain("   _(  )_( )_\n")
		log.Plain("  (_   _    _)\n")
		log.Plain("    (_) (__)\n\n")
		log.Plain("Welcome to Dnote Cloud :)\n\n")
		log.Plain("A home for your engineering microlessons\n")
		log.Plain("You can register at https://dnote.io/cloud\n\n")
		log.Printf("API key: ")

		var apiKey string
		fmt.Scanln(&apiKey)

		if apiKey == "" {
			return errors.New("Empty API key")
		}

		config, err := core.ReadConfig(ctx)
		if err != nil {
			return err
		}

		config.APIKey = apiKey
		err = core.WriteConfig(ctx, config)
		if err != nil {
			return errors.Wrap(err, "Failed to write to config file")
		}

		log.Success("configured\n")

		return nil
	}

}
