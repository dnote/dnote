package upgrade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

func TestFetchLatestStableTag(t *testing.T) {
	tagCLI0_1_0 := "cli-v0.1.0"
	tagCLI0_1_1 := "cli-v0.1.1"
	tagCLI0_1_2Beta := "cli-v0.1.2-beta"
	tagCLI0_1_3 := "cli-v0.1.3"
	tagServer0_1_0 := "server-v0.1.0"

	prereleaseTrue := true

	testCases := []struct {
		releases []*github.RepositoryRelease
		expected string
	}{
		{
			releases: []*github.RepositoryRelease{{TagName: &tagCLI0_1_0}},
			expected: tagCLI0_1_0,
		},
		{
			releases: []*github.RepositoryRelease{
				{TagName: &tagCLI0_1_1},
				{TagName: &tagServer0_1_0},
				{TagName: &tagCLI0_1_0},
			},
			expected: tagCLI0_1_1,
		},
		{
			releases: []*github.RepositoryRelease{
				{TagName: &tagServer0_1_0},
				{TagName: &tagCLI0_1_1},
				{TagName: &tagCLI0_1_0},
			},
			expected: tagCLI0_1_1,
		},
		{
			releases: []*github.RepositoryRelease{
				{TagName: &tagCLI0_1_2Beta, Prerelease: &prereleaseTrue},
				{TagName: &tagServer0_1_0},
				{TagName: &tagCLI0_1_1},
				{TagName: &tagCLI0_1_0},
			},
			expected: tagCLI0_1_1,
		},
		{
			releases: []*github.RepositoryRelease{
				{TagName: &tagCLI0_1_3},
				{TagName: &tagCLI0_1_2Beta, Prerelease: &prereleaseTrue},
				{TagName: &tagCLI0_1_1},
				{TagName: &tagCLI0_1_0},
			},
			expected: tagCLI0_1_3,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			// setup
			apiHandler := http.NewServeMux()
			apiHandler.HandleFunc("/repos/dnote/dnote/releases", func(w http.ResponseWriter, r *http.Request) {
				if err := json.NewEncoder(w).Encode(tc.releases); err != nil {
					t.Fatal(errors.Wrap(err, "responding with mock releases"))
				}
			})

			server := httptest.NewServer(apiHandler)
			url, err := url.Parse(server.URL + "/")
			if err != nil {
				t.Fatal(errors.Wrap(err, "parsing mock server url"))
			}

			client := github.NewClient(nil)
			client.BaseURL = url
			client.UploadURL = url

			// execute
			got, err := fetchLatestStableTag(client, 0)
			if err != nil {
				t.Fatal(errors.Wrap(err, "performing"))
			}

			// test
			assert.Equal(t, got, tc.expected, "result mismatch")
		})
	}

}
