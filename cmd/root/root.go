package root

import (
	"github.com/dnote-io/cli/upgrade"
	"github.com/dnote-io/cli/utils"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "dnote",
	Short: "Dnote - Instantly capture what you learn while coding",
}

// Register adds a new command
func Register(cmd *cobra.Command) {
	root.AddCommand(cmd)
}

// Execute runs the main command
func Execute() error {
	return root.Execute()
}

// Prepare initializes necessary files
func Prepare() error {
	configPath, err := utils.GetConfigPath()
	if err != nil {
		return err
	}
	dnotePath, err := utils.GetDnotePath()
	if err != nil {
		return err
	}
	dnoteUpdatePath, err := utils.GetDnoteUpdatePath()
	if err != nil {
		return err
	}

	if !utils.CheckFileExists(configPath) {
		err := utils.GenerateConfigFile()
		if err != nil {
			return err
		}
	}
	if !utils.CheckFileExists(dnotePath) {
		err := utils.TouchDnoteFile()
		if err != nil {
			return err
		}
	}
	if !utils.CheckFileExists(dnoteUpdatePath) {
		err := utils.TouchDnoteUpgradeFile()
		if err != nil {
			return err
		}
	}

	err = upgrade.Migrate()
	if err != nil {
		return err
	}

	return nil
}
