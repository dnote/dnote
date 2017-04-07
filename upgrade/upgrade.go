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
	"os/user"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/dnote-io/cli/utils"
	"github.com/google/go-github/github"
)

const version = "0.0.3"

func GetDnoteUpdatePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, utils.DnoteUpdateFilename), nil
}

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
func getLastUpdateEpoch() (int64, error) {
	updatePath, err := utils.GetDnoteUpdatePath()
	if err != nil {
		return 0, err
	}

	b, err := ioutil.ReadFile(updatePath)
	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`LAST_UPDATE_EPOCH: (\d+)\n`)
	match := re.FindStringSubmatch(string(b))

	if len(match) != 2 {
		msg := fmt.Sprintf("Error parsing %s", utils.DnoteUpdateFilename)
		return 0, errors.New(msg)
	}

	lastEpoch, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return lastEpoch, nil
}

// shouldCheckUpdate checks if update should be checked
func shouldCheckUpdate() (bool, error) {
	var updatePeriod int64 = 86400 * 7

	now := time.Now().Unix()
	lastEpoch, err := getLastUpdateEpoch()
	if err != nil {
		return false, err
	}

	return now-lastEpoch > updatePeriod, nil
}

// AutoUpdate triggers update if needed
func AutoUpdate() error {
	shouldCheck, err := shouldCheckUpdate()
	if err != nil {
		return err
	}

	if shouldCheck {
		tryUpgrade()
	}

	return nil
}

func tryUpgrade() error {
	defer utils.TouchDnoteUpgradeFile()

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
	if latestVersion == version {
		fmt.Printf("Up-to-date: %s", version)
		return nil
	}

	asset := getAsset(latest)
	if asset == nil {
		fmt.Printf("Could not find the release for %s %s", runtime.GOOS, runtime.GOARCH)
		return nil
	}

	// Download temporary file
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

	fmt.Printf("Updated: v%s -> v%s", version, latestVersion)
	fmt.Println("Changelog: https://github.com/dnote-io/cli/releases")
	return nil
}
