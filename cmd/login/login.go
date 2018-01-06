package login

import (
	"fmt"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/dnote-io/cli/utils"
	"github.com/spf13/cobra"
)

var example = `
  dnote login`

var cmd = &cobra.Command{
	Use:     "login",
	Short:   "Login to dnote server",
	Example: example,
	RunE:    run,
}

func init() {
	root.Register(cmd)
}

func run(cmd *cobra.Command, args []string) error {
	fmt.Print("Please enter your APIKey: ")

	var apiKey string
	fmt.Scanln(&apiKey)

	config, err := utils.ReadConfig()
	if err != nil {
		return err
	}

	config.APIKey = apiKey
	err = utils.WriteConfig(config)
	if err != nil {
		return err
	}

	return nil
}
