package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dnote-io/cli/cmd/root"
	"github.com/dnote-io/cli/infra"
	"github.com/pkg/errors"

	// commands
	"github.com/dnote-io/cli/cmd/add"
	"github.com/dnote-io/cli/cmd/books"
	"github.com/dnote-io/cli/cmd/edit"
	"github.com/dnote-io/cli/cmd/login"
	"github.com/dnote-io/cli/cmd/ls"
	"github.com/dnote-io/cli/cmd/remove"
	"github.com/dnote-io/cli/cmd/sync"
	"github.com/dnote-io/cli/cmd/upgrade"
	"github.com/dnote-io/cli/cmd/use"
	"github.com/dnote-io/cli/cmd/version"
)

var ctx infra.DnoteCtx

// initialize the dnote context and prepare dir structure
func init() {
	var err error
	ctx, err = newCtx()
	if err != nil {
		panic(errors.Wrap(err, "Failed to initialize the dnote context"))
	}

	err = root.Prepare(ctx)
	if err != nil {
		panic(errors.Wrap(err, "Failed to prepare dnote run"))
	}
}

func main() {
	root.Register(books.NewCmd(ctx))
	root.Register(remove.NewCmd(ctx))
	root.Register(edit.NewCmd(ctx))
	root.Register(login.NewCmd(ctx))
	root.Register(add.NewCmd(ctx))
	root.Register(ls.NewCmd(ctx))
	root.Register(sync.NewCmd(ctx))
	root.Register(use.NewCmd(ctx))
	root.Register(version.NewCmd(ctx))
	root.Register(upgrade.NewCmd(ctx))

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newCtx() (infra.DnoteCtx, error) {
	usr, err := user.Current()
	if err != nil {
		return infra.DnoteCtx{}, errors.Wrap(err, "Failed to get current user")
	}

	homeDir := usr.HomeDir
	dnoteDir := getDnoteDir(homeDir)

	ret := infra.DnoteCtx{
		HomeDir:  homeDir,
		DnoteDir: dnoteDir,
	}

	return ret, nil
}

func getDnoteDir(homeDir string) string {
	var ret string

	dnoteDirEnv := os.Getenv("DNOTE_DIR")
	if dnoteDirEnv == "" {
		ret = fmt.Sprintf("%s/%s", homeDir, infra.DnoteDirName)
	} else {
		ret = dnoteDirEnv
	}

	return ret
}
