/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"os"

	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	// commands
	"github.com/dnote/dnote/pkg/cli/cmd/add"
	"github.com/dnote/dnote/pkg/cli/cmd/cat"
	"github.com/dnote/dnote/pkg/cli/cmd/edit"
	"github.com/dnote/dnote/pkg/cli/cmd/find"
	"github.com/dnote/dnote/pkg/cli/cmd/login"
	"github.com/dnote/dnote/pkg/cli/cmd/logout"
	"github.com/dnote/dnote/pkg/cli/cmd/ls"
	"github.com/dnote/dnote/pkg/cli/cmd/remove"
	"github.com/dnote/dnote/pkg/cli/cmd/root"
	"github.com/dnote/dnote/pkg/cli/cmd/sync"
	"github.com/dnote/dnote/pkg/cli/cmd/version"
	"github.com/dnote/dnote/pkg/cli/cmd/view"
)

// apiEndpoint and versionTag are populated during link time
var apiEndpoint string
var versionTag = "master"

func main() {
	ctx, err := infra.Init(apiEndpoint, versionTag)
	if err != nil {
		panic(errors.Wrap(err, "initializing context"))
	}
	defer ctx.DB.Close()

	root.Register(remove.NewCmd(*ctx))
	root.Register(edit.NewCmd(*ctx))
	root.Register(login.NewCmd(*ctx))
	root.Register(logout.NewCmd(*ctx))
	root.Register(add.NewCmd(*ctx))
	root.Register(ls.NewCmd(*ctx))
	root.Register(sync.NewCmd(*ctx))
	root.Register(version.NewCmd(*ctx))
	root.Register(cat.NewCmd(*ctx))
	root.Register(view.NewCmd(*ctx))
	root.Register(find.NewCmd(*ctx))

	if err := root.Execute(); err != nil {
		log.Errorf("%s\n", err.Error())
		os.Exit(1)
	}
}
