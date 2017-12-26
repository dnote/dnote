package upgrade

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/utils"
	"github.com/google/go-github/github"
)

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

// getLastUpdateEpoch reads and parses the last update epoch
func getLastUpdateEpoch(ctx infra.DnoteCtx) (int64, error) {
	updatePath := infra.GetTimestampPath(ctx)

	b, err := ioutil.ReadFile(updatePath)
	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`LAST_UPGRADE_EPOCH: (\d+)`)
	match := re.FindStringSubmatch(string(b))

	if len(match) != 2 {
		msg := fmt.Sprintf("Error parsing %s: %s", infra.TimestampFilename, string(b))
		return 0, errors.New(msg)
	}

	lastEpoch, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return lastEpoch, nil
}

// shouldCheckUpdate checks if update should be checked
func shouldCheckUpdate(ctx infra.DnoteCtx) (bool, error) {
	var updatePeriod int64 = 86400 * 7

	now := time.Now().Unix()
	lastEpoch, err := getLastUpdateEpoch(ctx)
	if err != nil {
		return false, err
	}

	return now-lastEpoch > updatePeriod, nil
}

// AutoUpgrade triggers update if needed
func AutoUpgrade(ctx infra.DnoteCtx) error {
	shouldCheck, err := shouldCheckUpdate(ctx)
	if err != nil {
		return err
	}

	if shouldCheck {
		willCheck, err := utils.AskConfirmation("Would you like to check for an update?")
		infra.InitTimestampFile(ctx)
		if err != nil {
			return err
		}

		if willCheck {
			err := Upgrade(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Upgrade(ctx infra.DnoteCtx) error {
	// Fetch the latest version
	gh := github.NewClient(nil)
	releases, _, err := gh.Repositories.ListReleases(context.Background(), "dnote-io", "cli", nil)

	if err != nil {
		return err
	}

	latest := releases[0]
	latestVersion := (*latest.TagName)[1:]

	if err != nil {
		return err
	}

	// Check if up to date
	if latestVersion == infra.Version {
		fmt.Printf("Up-to-date: %s\n", infra.Version)
		infra.InitTimestampFile(ctx)
		return nil
	}

	asset := getAsset(latest)
	if asset == nil {
		infra.InitTimestampFile(ctx)
		fmt.Printf("Could not find the release for %s %s", runtime.GOOS, runtime.GOARCH)
		return nil
	}

	// Download temporary file
	fmt.Printf("Downloading: %s\n", latestVersion)
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

	infra.InitTimestampFile(ctx)

	fmt.Printf("Updated: v%s -> v%s\n", infra.Version, latestVersion)
	fmt.Println("Changelog: https://github.com/dnote-io/cli/releases")
	return nil
}
