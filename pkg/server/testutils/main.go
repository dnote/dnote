/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package testutils provides utilities used in tests
package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/pkg/errors"
)

// HTTPDo makes an HTTP request and returns a response
func HTTPDo(t *testing.T, req *http.Request) *http.Response {
	hc := http.Client{
		// Do not follow redirects.
		// e.g. /logout redirects to a page but we'd like to test the redirect
		// itself, not what happens after the redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := hc.Do(req)
	if err != nil {
		t.Fatal(errors.Wrap(err, "performing http request"))
	}

	return res
}

// MakeReq makes an HTTP request and returns a response
func MakeReq(endpoint string, method, path, data string) *http.Request {
	u := fmt.Sprintf("%s%s", endpoint, path)

	req, err := http.NewRequest(method, u, strings.NewReader(data))
	if err != nil {
		panic(errors.Wrap(err, "constructing http request"))
	}

	return req
}

// GetCookieByName returns a cookie with the given name
func GetCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	var ret *http.Cookie

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == name {
			ret = cookies[i]
			break
		}
	}

	return ret
}

// MustRespondJSON responds with the JSON-encoding of the given interface. If the encoding
// fails, the test fails. It is used by test servers.
func MustRespondJSON(t *testing.T, w http.ResponseWriter, i interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(i); err != nil {
		t.Fatal(message)
	}
}

// MockEmail is a mock email data
type MockEmail struct {
	Subject string
	From    string
	To      []string
	Body    string
}

// MockEmailbackendImplementation is an email backend that simply discards the emails
type MockEmailbackendImplementation struct {
	mu     sync.RWMutex
	Emails []MockEmail
}

// Clear clears the mock email queue
func (b *MockEmailbackendImplementation) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Emails = []MockEmail{}
}

// Queue is an implementation of Backend.Queue.
func (b *MockEmailbackendImplementation) Queue(subject, from string, to []string, contentType, body string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Emails = append(b.Emails, MockEmail{
		Subject: subject,
		From:    from,
		To:      to,
		Body:    body,
	})

	return nil
}
