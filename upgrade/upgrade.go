package upgrade

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/dnote/cli/utils"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// upgradeInterval is 3 weeks
var upgradeInterval int64 = 86400 * 7 * 3

// getAsset finds the asset to download from the liast of assets in a release
func getAsset(release *github.RepositoryRelease) *github.ReleaseAsset {
	filename := fmt.Sprintf("dnote-%s-%s", runtime.GOOS, runtime.GOARCH)

	for _, asset := range release.Assets {
		if *asset.Name == filename {
			return &asset
		}
	}

	return nil
}

// shouldCheckUpdate checks if update should be checked
func shouldCheckUpdate(ctx infra.DnoteCtx) (bool, error) {
	timestamp, err := core.ReadTimestamp(ctx)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get timestamp content")
	}

	now := time.Now().Unix()

	return now-timestamp.LastUpgrade > upgradeInterval, nil
}

func touchLastUpgrade(ctx infra.DnoteCtx) error {
	timestamp, err := core.ReadTimestamp(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get timestamp content")
	}

	now := time.Now().Unix()
	timestamp.LastUpgrade = now

	if err := core.WriteTimestamp(ctx, timestamp); err != nil {
		return errors.Wrap(err, "Failed to write the updated timestamp to the file")
	}

	return nil
}

// AutoUpgrade triggers update if needed
func AutoUpgrade(ctx infra.DnoteCtx) error {
	shouldCheck, err := shouldCheckUpdate(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to check if dnote should check update")
	}

	if shouldCheck {
		willCheck, err := utils.AskConfirmation("check for upgrade?", true)
		if err != nil {
			return errors.Wrap(err, "Failed to get user confirmation for checking upgrade")
		}

		err = touchLastUpgrade(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to update last upgrade timestamp")
		}

		if willCheck {
			if err := Upgrade(ctx); err != nil {
				return errors.Wrap(err, "Failed to upgrade")
			}
		}
	}

	return nil
}

func Upgrade(ctx infra.DnoteCtx) error {
	log.Infof("current version is %s\n", core.Version)

	// Fetch the latest version
	gh := github.NewClient(nil)
	releases, _, err := gh.Repositories.ListReleases(context.Background(), "dnote", "cli", nil)

	if err != nil {
		return err
	}

	latest := releases[0]
	latestVersion := (*latest.TagName)[1:]

	if err != nil {
		return err
	}

	log.Infof("latest version is %s\n", latestVersion)

	// Check if up to date
	if latestVersion == core.Version {
		log.Success("you are up-to-date\n\n")
		err = touchLastUpgrade(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to update the upgrade timestamp")
		}

		return nil
	}

	asset := getAsset(latest)
	if asset == nil {
		err = touchLastUpgrade(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to update the upgrade timestamp")
		}

		return errors.Errorf("Could not find the release for %s %s", runtime.GOOS, runtime.GOARCH)
	}

	// Download temporary file
	log.Infof("downloading: %s\n", latestVersion)
	tmpPath := path.Join(os.TempDir(), "dnote_update")

	out, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(*asset.BrowserDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Override the binary
	cmdPath, err := exec.LookPath("dnote")
	if err != nil {
		return err
	}

	err = os.Rename(tmpPath, cmdPath)
	if err != nil {
		return err
	}

	// Make it executable
	err = os.Chmod(cmdPath, 0755)
	if err != nil {
		return err
	}

	err = touchLastUpgrade(ctx)
	if err != nil {
		return errors.Wrap(err, "Upgrade is done, but failed to update the last_upgrade timestamp.")
	}

	log.Successf("updated: v%s -> v%s\n", core.Version, latestVersion)
	log.Plain("changelog: https://github.com/dnote/cli/releases\n\n")
	return nil
}
