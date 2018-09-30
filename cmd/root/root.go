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
	if err := core.InitFiles(ctx); err != nil {
		return errors.Wrap(err, "initializing files")
	}

	if err := infra.InitDB(ctx); err != nil {
		return errors.Wrap(err, "initializing database")
	}
	if err := infra.InitSystem(ctx); err != nil {
		return errors.Wrap(err, "initializing system data")
	}

	if err := migrate.Legacy(ctx); err != nil {
		return errors.Wrap(err, "running legacy migration")
	}
	if err := migrate.Run(ctx, migrate.LocalSequence); err != nil {
		return errors.Wrap(err, "running migration")
	}

	return nil
}
