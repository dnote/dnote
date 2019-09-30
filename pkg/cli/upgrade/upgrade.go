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

package upgrade

import (
	stdCtx "context"
	"fmt"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// upgradeInterval is 3 weeks
var upgradeInterval int64 = 86400 * 7 * 3

// shouldCheckUpdate checks if update should be checked
func shouldCheckUpdate(ctx context.DnoteCtx) (bool, error) {
	db := ctx.DB

	var lastUpgrade int64
	err := db.QueryRow("SELECT value FROM system WHERE key = ?", consts.SystemLastUpgrade).Scan(&lastUpgrade)
	if err != nil {
		return false, errors.Wrap(err, "getting last_udpate")
	}

	now := time.Now().Unix()

	return now-lastUpgrade > upgradeInterval, nil
}

func touchLastUpgrade(ctx context.DnoteCtx) error {
	db := ctx.DB

	now := time.Now().Unix()
	_, err := db.Exec("UPDATE system SET value = ? WHERE key = ?", now, consts.SystemLastUpgrade)
	if err != nil {
		return errors.Wrap(err, "updating last_upgrade")
	}

	return nil
}

func fetchLatestStableTag(gh *github.Client, page int) (string, error) {
	params := github.ListOptions{
		Page: page,
	}
	releases, resp, err := gh.Repositories.ListReleases(stdCtx.Background(), "dnote", "dnote", &params)
	if err != nil {
		return "", errors.Wrapf(err, "fetching releases page %d", page)
	}

	for _, release := range releases {
		tag := release.GetTagName()
		isStable := !release.GetPrerelease()

		if strings.HasPrefix(tag, "cli-") && isStable {
			return tag, nil
		}
	}

	if page == resp.LastPage {
		return "", errors.New("No CLI release was found")
	}

	return fetchLatestStableTag(gh, page+1)
}

func checkVersion(ctx context.DnoteCtx) error {
	log.Infof("current version is %s\n", ctx.Version)

	// Fetch the latest version
	gh := github.NewClient(nil)
	latestTag, err := fetchLatestStableTag(gh, 1)
	if err != nil {
		return errors.Wrap(err, "fetching the latest stable release")
	}

	// releases are tagged in a form of cli-v1.0.0
	latestVersion := latestTag[5:]
	log.Infof("latest version is %s\n", latestVersion)

	if latestVersion == ctx.Version {
		log.Success("you are up-to-date\n\n")
	} else {
		log.Infof("to upgrade, see https://github.com/dnote/dnote/pkg/cli/blob/master/README.md\n")
	}

	return nil
}

// Check triggers update if needed
func Check(ctx context.DnoteCtx) error {
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
	willCheck, err := ui.Confirm("check for upgrade?", true)
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
