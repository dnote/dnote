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
		return errors.Wrap(err, "initializing dnote dir")
	}

	err = core.InitFiles(ctx)
	if err != nil {
		return errors.Wrap(err, "initiating files")
	}

	err = migrate.Migrate(ctx)
	if err != nil {
		return errors.Wrap(err, "running migration")
	}

	return nil
}
