package root

import (
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/migrate"
	"github.com/dnote-io/cli/upgrade"
	"github.com/pkg/errors"
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
func Prepare(ctx infra.DnoteCtx) error {
	err := infra.MigrateToDnoteDir(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize dnote dir")
	}

	fresh, err := infra.IsFreshInstall(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to check if fresh install")
	}

	err = infra.InitDnoteDir(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create dnote dir")
	}
	err = infra.InitConfigFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to generate config file")
	}
	err = infra.InitDnoteFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create dnote file")
	}
	err = infra.InitTimestampFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create dnote upgrade file")
	}
	err = migrate.InitSchemaFile(ctx, fresh)
	if err != nil {
		return errors.Wrap(err, "Failed to create migration file")
	}

	err = migrate.Migrate(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to perform migration")
	}

	err = upgrade.AutoUpgrade(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to auto upgrade")
	}

	return nil
}
