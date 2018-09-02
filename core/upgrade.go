package core

import (
	"context"
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
	timestamp, err := ReadTimestamp(ctx)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get timestamp content")
	}

	now := time.Now().Unix()

	return now-timestamp.LastUpgrade > upgradeInterval, nil
}

func touchLastUpgrade(ctx infra.DnoteCtx) error {
	timestamp, err := ReadTimestamp(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get timestamp content")
	}

	now := time.Now().Unix()
	timestamp.LastUpgrade = now

	if err := WriteTimestamp(ctx, timestamp); err != nil {
		return errors.Wrap(err, "Failed to write the updated timestamp to the file")
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
		log.Info("to upgrade, see https://github.com/dnote/cli/blob/master/README.md")
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
