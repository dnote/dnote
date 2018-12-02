package core

import (
	"context"
	"fmt"
	"time"

	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/dnote/cli/utils"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// upgradeInterval is 3 weeks
var upgradeInterval int64 = 86400 * 7 * 3

// shouldCheckUpdate checks if update should be checked
func shouldCheckUpdate(ctx infra.DnoteCtx) (bool, error) {
	db := ctx.DB

	var lastUpgrade int64
	err := db.QueryRow("SELECT value FROM system WHERE key = ?", infra.SystemLastUpgrade).Scan(&lastUpgrade)
	if err != nil {
		return false, errors.Wrap(err, "getting last_udpate")
	}

	now := time.Now().Unix()

	return now-lastUpgrade > upgradeInterval, nil
}

func touchLastUpgrade(ctx infra.DnoteCtx) error {
	db := ctx.DB

	now := time.Now().Unix()
	_, err := db.Exec("UPDATE system SET value = ? WHERE key = ?", now, infra.SystemLastUpgrade)
	if err != nil {
		return errors.Wrap(err, "updating last_upgrade")
	}

	return nil
}

func checkVersion(ctx infra.DnoteCtx) error {
	log.Infof("current version is %s\n", ctx.Version)

	// Fetch the latest version
	gh := github.NewClient(nil)
	releases, _, err := gh.Repositories.ListReleases(context.Background(), "dnote", "cli", nil)
	if err != nil {
		return errors.Wrap(err, "fetching releases")
	}

	latest := releases[0]
	latestVersion := (*latest.TagName)[1:]

	log.Infof("latest version is %s\n", latestVersion)

	if latestVersion == ctx.Version {
		log.Success("you are up-to-date\n\n")
	} else {
		log.Infof("to upgrade, see https://github.com/dnote/cli/blob/master/README.md\n")
	}

	return nil
}

// CheckUpdate triggers update if needed
func CheckUpdate(ctx infra.DnoteCtx) error {
	shouldCheck, err := shouldCheckUpdate(ctx)
	if err != nil {
		return errors.Wrap(err, "checking if dnote should check update")
	}
	if !shouldCheck {
		return nil
	}

	err = touchLastUpgrade(ctx)
	if err != nil {
		return errors.Wrap(err, "updating the last upgrade timestamp")
	}

	fmt.Printf("\n")
	willCheck, err := utils.AskConfirmation("check for upgrade?", true)
	if err != nil {
		return errors.Wrap(err, "getting user confirmation")
	}
	if !willCheck {
		return nil
	}

	err = checkVersion(ctx)
	if err != nil {
		return errors.Wrap(err, "checking version")
	}

	return nil
}
