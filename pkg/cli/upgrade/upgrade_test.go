/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

func setupGithubClient(t *testing.T) (*github.Client, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	url, err := url.Parse(server.URL + "/")
	if err != nil {
		t.Fatal(errors.Wrap(err, "parsing mock server url"))
	}

	client := github.NewClient(nil)
	client.BaseURL = url
	client.UploadURL = url

	return client, mux
}

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
			gh, mux := setupGithubClient(t)
			mux.HandleFunc("/repos/dnote/dnote/releases", func(w http.ResponseWriter, r *http.Request) {
				if err := json.NewEncoder(w).Encode(tc.releases); err != nil {
					t.Fatal(errors.Wrap(err, "responding with mock releases"))
				}
			})

			// execute
			got, err := fetchLatestStableTag(gh, 0)
			if err != nil {
				t.Fatal(errors.Wrap(err, "performing"))
			}

			// test
			assert.Equal(t, got, tc.expected, "result mismatch")
		})
	}
}

func TestFetchLatestStableTag_paginated(t *testing.T) {
	tagServer0_1_0 := "server-v0.1.0"
	tagCLI0_1_2Beta := "cli-v0.1.2-beta"
	tagCLI0_1_1 := "cli-v0.1.1"
	prereleaseTrue := true

	// set up
	gh, mux := setupGithubClient(t)
	path := "/repos/dnote/dnote/releases"
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		page := r.FormValue("page")

		releasesPage1 := []*github.RepositoryRelease{
			{TagName: &tagServer0_1_0},
		}
		releasesPage2 := []*github.RepositoryRelease{
			{TagName: &tagCLI0_1_2Beta, Prerelease: &prereleaseTrue},
			{TagName: &tagCLI0_1_1},
		}

		baseURL := gh.BaseURL.String()

		switch page {
		case "", "1":
			linkHeader := fmt.Sprintf("<%s%s?page=2>; rel=\"next\" <%s%s?page=2>; rel=\"last\"", baseURL, path, baseURL, path)
			w.Header().Set("Link", linkHeader)

			if err := json.NewEncoder(w).Encode(releasesPage1); err != nil {
				t.Fatal(errors.Wrap(err, "responding with mock releases"))
			}
		case "2":
			linkHeader := fmt.Sprintf("<%s%s?page=1>; rel=\"prev\" <%s%s?page=1>; rel=\"first\"", baseURL, path, baseURL, path)
			w.Header().Set("Link", linkHeader)

			if err := json.NewEncoder(w).Encode(releasesPage2); err != nil {
				t.Fatal(errors.Wrap(err, "responding with mock releases"))
			}
		default:
			t.Fatal("Should have stopped walking")
		}
	})

	// execute
	got, err := fetchLatestStableTag(gh, 0)
	if err != nil {
		t.Fatal(errors.Wrap(err, "performing"))
	}

	// test
	assert.Equal(t, got, tagCLI0_1_1, "result mismatch")
}
