package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dnote/cli/cmd/root"
	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/pkg/errors"

	// commands
	"github.com/dnote/cli/cmd/add"
	"github.com/dnote/cli/cmd/cat"
	"github.com/dnote/cli/cmd/edit"
	"github.com/dnote/cli/cmd/login"
	"github.com/dnote/cli/cmd/ls"
	"github.com/dnote/cli/cmd/remove"
	"github.com/dnote/cli/cmd/sync"
	"github.com/dnote/cli/cmd/upgrade"
	"github.com/dnote/cli/cmd/version"
	"github.com/dnote/cli/cmd/view"
)

// apiEndpoint is populated during link time
var apiEndpoint string

func main() {
	ctx, err := newCtx()
	if err != nil {
		panic(errors.Wrap(err, "Failed to initialize the dnote context"))
	}

	err = root.Prepare(ctx)
	if err != nil {
		panic(errors.Wrap(err, "Failed to prepare dnote run"))
	}

	root.Register(remove.NewCmd(ctx))
	root.Register(edit.NewCmd(ctx))
	root.Register(login.NewCmd(ctx))
	root.Register(add.NewCmd(ctx))
	root.Register(ls.NewCmd(ctx))
	root.Register(sync.NewCmd(ctx))
	root.Register(version.NewCmd(ctx))
	root.Register(upgrade.NewCmd(ctx))
	root.Register(cat.NewCmd(ctx))
	root.Register(view.NewCmd(ctx))

	if err := root.Execute(); err != nil {
		log.Errorf("%s\n", err.Error())
		os.Exit(1)
	}
}

func newCtx() (infra.DnoteCtx, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return infra.DnoteCtx{}, errors.Wrap(err, "Failed to get home dir")
	}
	dnoteDir := getDnoteDir(homeDir)

	ret := infra.DnoteCtx{
		HomeDir:     homeDir,
		DnoteDir:    dnoteDir,
		APIEndpoint: apiEndpoint,
	}

	return ret, nil
}

func getDnoteDir(homeDir string) string {
	var ret string

	dnoteDirEnv := os.Getenv("DNOTE_DIR")
	if dnoteDirEnv == "" {
		ret = fmt.Sprintf("%s/%s", homeDir, core.DnoteDirName)
	} else {
		ret = dnoteDirEnv
	}

	return ret
}

func getHomeDir() (string, error) {
	homeDirEnv := os.Getenv("DNOTE_HOME_DIR")
	if homeDirEnv != "" {
		return homeDirEnv, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get current user")
	}

	return usr.HomeDir, nil
}
