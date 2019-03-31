package main

import (
	"os"

	"github.com/dnote/cli/cmd/root"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	// commands
	"github.com/dnote/cli/cmd/add"
	"github.com/dnote/cli/cmd/cat"
	"github.com/dnote/cli/cmd/edit"
	"github.com/dnote/cli/cmd/find"
	"github.com/dnote/cli/cmd/login"
	"github.com/dnote/cli/cmd/logout"
	"github.com/dnote/cli/cmd/ls"
	"github.com/dnote/cli/cmd/remove"
	"github.com/dnote/cli/cmd/sync"
	"github.com/dnote/cli/cmd/version"
	"github.com/dnote/cli/cmd/view"
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
