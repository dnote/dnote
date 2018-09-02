package root

import (
	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/migrate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:           "dnote",
	Short:         "Dnote - Instantly capture what you learn while coding",
	SilenceErrors: true,
	SilenceUsage:  true,
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
	err := core.MigrateToDnoteDir(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize dnote dir")
	}

	fresh, err := core.IsFreshInstall(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to check if fresh install")
	}

	err = core.InitDnoteDir(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create dnote dir")
	}
	err = core.InitConfigFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to generate config file")
	}
	err = core.InitDnoteFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create dnote file")
	}
	err = core.InitTimestampFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create dnote upgrade file")
	}
	err = core.InitActionFile(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to create action file")
	}
	err = migrate.InitSchemaFile(ctx, fresh)
	if err != nil {
		return errors.Wrap(err, "Failed to create migration file")
	}

	err = migrate.Migrate(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to perform migration")
	}

	return nil
}
