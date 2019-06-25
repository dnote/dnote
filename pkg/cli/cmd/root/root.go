/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package root

import (
	"github.com/dnote/dnote/pkg/cli/context"
	// "github.com/dnote/dnote/pkg/cli/core"
	// "github.com/dnote/dnote/pkg/cli/infra"
	// "github.com/dnote/dnote/pkg/cli/migrate"
	// "github.com/pkg/errors"
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
func Prepare(ctx context.DnoteCtx) error {
	// 	if err := core.InitFiles(ctx); err != nil {
	// 		return errors.Wrap(err, "initializing files")
	// 	}
	//
	// 	if err := infra.InitDB(ctx); err != nil {
	// 		return errors.Wrap(err, "initializing database")
	// 	}
	// 	if err := core.InitSystem(ctx); err != nil {
	// 		return errors.Wrap(err, "initializing system data")
	// 	}
	//
	// 	if err := migrate.Legacy(ctx); err != nil {
	// 		return errors.Wrap(err, "running legacy migration")
	// 	}
	// 	if err := migrate.Run(ctx, migrate.LocalSequence, migrate.LocalMode); err != nil {
	// 		return errors.Wrap(err, "running migration")
	// 	}

	return nil
}
