package core

import (
	"github.com/dnote-io/cli/test"
	"github.com/dnote-io/cli/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateToDnoteDir(t *testing.T) {
	ctx := test.InitCtx("../tmp")

	t.Run("pre v1 files exist", func(t *testing.T) {
		// set up
		if err := os.MkdirAll(ctx.HomeDir, 0755); err != nil {
			panic(err)
		}
		defer func() {
			if err := os.RemoveAll(ctx.HomeDir); err != nil {
				panic(err)
			}
		}()

		dnotePath, err := filepath.Abs(filepath.Join(ctx.HomeDir, ".dnote"))
		if err != nil {
			panic(errors.Wrap(err, "Failed to get absolute .dnote path").Error())
		}
		dnotercPath, err := filepath.Abs(filepath.Join(ctx.HomeDir, ".dnoterc"))
		if err != nil {
			panic(errors.Wrap(err, "Failed to get absolute .dnote path").Error())
		}
		dnoteUpgradePath, err := filepath.Abs(filepath.Join(ctx.HomeDir, ".dnote-upgrade"))
		if err != nil {
			panic(errors.Wrap(err, "Failed to get absolute .dnote path").Error())
		}

		if err = ioutil.WriteFile(dnotePath, []byte{}, 0644); err != nil {
			panic(errors.Wrap(err, "Failed prepare .dnote").Error())
		}
		if err = ioutil.WriteFile(dnotercPath, []byte{}, 0644); err != nil {
			panic(errors.Wrap(err, "Failed prepare .dnoterc").Error())
		}
		if err = ioutil.WriteFile(dnoteUpgradePath, []byte{}, 0644); err != nil {
			panic(errors.Wrap(err, "Failed prepare .dnote-upgrade").Error())
		}

		// execute
		err = MigrateToDnoteDir(ctx)
		if err != nil {
			panic(errors.Wrap(err, "Failed to perform").Error())
		}

		// test
		newDnotePath, err := filepath.Abs(filepath.Join(ctx.DnoteDir, "dnote"))
		if err != nil {
			panic(errors.Wrap(err, "Failed get new dnote path").Error())
		}
		newDnotercPath, err := filepath.Abs(filepath.Join(ctx.DnoteDir, "dnoterc"))
		if err != nil {
			panic(errors.Wrap(err, "Failed get new dnoterc path").Error())
		}
		newTimestampPath, err := filepath.Abs(filepath.Join(ctx.DnoteDir, "timestamps"))
		if err != nil {
			panic(errors.Wrap(err, "Failed get new timestamp path").Error())
		}

		fi, err := os.Stat(dnotePath)
		if err != nil {
			panic(errors.Wrap(err, "Failed to look up file"))
		}
		if !fi.IsDir() {
			t.Fatal(".dnote must be a directory")
		}

		if utils.FileExists(dnotercPath) {
			t.Error(".dnoterc must not exist in the original location")
		}
		if utils.FileExists(dnoteUpgradePath) {
			t.Error(".dnote-upgrade must not exist in the original location")
		}
		if !utils.FileExists(newDnotePath) {
			t.Error("dnote must exist")
		}
		if !utils.FileExists(newDnotercPath) {
			t.Error("dnoterc must exist")
		}
		if !utils.FileExists(newTimestampPath) {
			t.Error("timestamp must exist")
		}
	})
}
