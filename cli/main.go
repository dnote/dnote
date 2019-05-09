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

package main

import (
	"os"

	"github.com/dnote/dnote/cli/cmd/root"
	"github.com/dnote/dnote/cli/infra"
	"github.com/dnote/dnote/cli/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	// commands
	"github.com/dnote/dnote/cli/cmd/add"
	"github.com/dnote/dnote/cli/cmd/cat"
	"github.com/dnote/dnote/cli/cmd/edit"
	"github.com/dnote/dnote/cli/cmd/find"
	"github.com/dnote/dnote/cli/cmd/login"
	"github.com/dnote/dnote/cli/cmd/logout"
	"github.com/dnote/dnote/cli/cmd/ls"
	"github.com/dnote/dnote/cli/cmd/remove"
	"github.com/dnote/dnote/cli/cmd/sync"
	"github.com/dnote/dnote/cli/cmd/version"
	"github.com/dnote/dnote/cli/cmd/view"
)

// apiEndpoint and versionTag are populated during link time
var apiEndpoint string
var versionTag = "master"

func main() {
	ctx, err := infra.NewCtx(apiEndpoint, versionTag)
	if err != nil {
		panic(errors.Wrap(err, "initializing context"))
	}
	defer ctx.DB.Close()

	if err := root.Prepare(ctx); err != nil {
		panic(errors.Wrap(err, "preparing dnote run"))
	}

	ctx, err = infra.SetupCtx(ctx)
	if err != nil {
		panic(errors.Wrap(err, "setting up context"))
	}

	root.Register(remove.NewCmd(ctx))
	root.Register(edit.NewCmd(ctx))
	root.Register(login.NewCmd(ctx))
	root.Register(logout.NewCmd(ctx))
	root.Register(add.NewCmd(ctx))
	root.Register(ls.NewCmd(ctx))
	root.Register(sync.NewCmd(ctx))
	root.Register(version.NewCmd(ctx))
	root.Register(cat.NewCmd(ctx))
	root.Register(view.NewCmd(ctx))
	root.Register(find.NewCmd(ctx))

	if err := root.Execute(); err != nil {
		log.Errorf("%s\n", err.Error())
		os.Exit(1)
	}
}
