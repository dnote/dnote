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

package core

import (
	"context"
	"fmt"
	"time"

	"github.com/dnote/dnote/cli/infra"
	"github.com/dnote/dnote/cli/log"
	"github.com/dnote/dnote/cli/utils"
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
		log.Infof("to upgrade, see https://github.com/dnote/dnote/cli/blob/master/README.md\n")
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
