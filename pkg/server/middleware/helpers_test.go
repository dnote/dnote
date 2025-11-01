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

package middleware

import (
	"net/http"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/pkg/errors"
)

func TestGetSessionKeyFromCookie(t *testing.T) {
	testCases := []struct {
		cookie   *http.Cookie
		expected string
	}{
		{
			cookie: &http.Cookie{
				Name:     "id",
				Value:    "foo",
				HttpOnly: true,
			},
			expected: "foo",
		},
		{
			cookie:   nil,
			expected: "",
		},
		{
			cookie: &http.Cookie{
				Name:     "foo",
				Value:    "bar",
				HttpOnly: true,
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		// set up
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(errors.Wrap(err, "constructing request"))
		}

		if tc.cookie != nil {
			r.AddCookie(tc.cookie)
		}

		// execute
		got, err := getSessionKeyFromCookie(r)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}

func TestGetSessionKeyFromAuth(t *testing.T) {
	testCases := []struct {
		authHeaderStr string
		expected      string
	}{
		{
			authHeaderStr: "Bearer foo",
			expected:      "foo",
		},
	}

	for _, tc := range testCases {
		// set up
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(errors.Wrap(err, "constructing request"))
		}

		r.Header.Set("Authorization", tc.authHeaderStr)

		// execute
		got, err := getSessionKeyFromAuth(r)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}

func mustMakeRequest(t *testing.T) *http.Request {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(errors.Wrap(err, "constructing request"))
	}

	return r
}

func TestGetCredential(t *testing.T) {
	r1 := mustMakeRequest(t)
	r2 := mustMakeRequest(t)
	r2.Header.Set("Authorization", "Bearer foo")
	r3 := mustMakeRequest(t)
	r3.Header.Set("Authorization", "Bearer bar")

	r4 := mustMakeRequest(t)
	c4 := http.Cookie{
		Name:     "id",
		Value:    "foo",
		HttpOnly: true,
	}
	r4.AddCookie(&c4)

	r5 := mustMakeRequest(t)
	c5 := http.Cookie{
		Name:     "id",
		Value:    "foo",
		HttpOnly: true,
	}
	r5.AddCookie(&c5)
	r5.Header.Set("Authorization", "Bearer foo")

	testCases := []struct {
		request  *http.Request
		expected string
	}{
		{
			request:  r1,
			expected: "",
		},
		{
			request:  r2,
			expected: "foo",
		},
		{
			request:  r3,
			expected: "bar",
		},
		{
			request:  r4,
			expected: "foo",
		},
		{
			request:  r5,
			expected: "foo",
		},
	}

	for _, tc := range testCases {
		// execute
		got, err := GetCredential(tc.request)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		assert.Equal(t, got, tc.expected, "result mismatch")
	}
}

