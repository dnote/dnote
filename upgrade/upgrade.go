package upgrade

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"time"

	"github.com/google/go-github/github"
)

const dnoteUpdateFilename = ".dnote-update"
const version = "0.0.3"

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
func getDnoteUpdatePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, dnoteUpdateFilename), nil
}

// shouldCheckUpdate checks if update should be checked
func shouldCheckUpdate() (bool, error) {
	updatePath, err := getDnoteUpdatePath()
	if err != nil {
		return false, err
	}

	b, err := ioutil.ReadFile(updatePath)
	if err != nil {
		return false, err
	}

	buf := bytes.NewBuffer(b)
	lastEpoch, err := binary.ReadVarint(buf)
	if err != nil {
		return false, err
	}

	now := time.Now().Unix()
	var epochTarget int64 = 86400 * 7 // 7 days

	return now-lastEpoch > epochTarget, nil
}

// AutoUpdate triggers update if needed
func AutoUpdate() error {
	shouldCheck, err := shouldCheckUpdate()
	if err != nil {
		return err
	}

	if shouldCheck {
		Update()
	}

	return nil
}

func Update() error {
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
	fmt.Println("Change note: https://github.com/dnote-io/cli/releases")
	return nil
}
