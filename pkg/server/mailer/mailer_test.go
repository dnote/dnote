/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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

package mailer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestEmailVerificationEmail(t *testing.T) {
	testCases := []struct {
		token  string
		webURL string
	}{
		{
			token:  "someRandomToken1",
			webURL: "http://localhost:3000",
		},
		{
			token:  "someRandomToken2",
			webURL: "http://localhost:3001",
		},
	}

	tmpl := NewTemplates()

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with WebURL %s", tc.webURL), func(t *testing.T) {
			dat := EmailVerificationTmplData{
				Token:  tc.token,
				WebURL: tc.webURL,
			}
			body, err := tmpl.Execute(EmailTypeEmailVerification, EmailKindText, dat)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			if ok := strings.Contains(body, tc.webURL); !ok {
				t.Errorf("email body did not contain %s", tc.webURL)
			}
			if ok := strings.Contains(body, tc.token); !ok {
				t.Errorf("email body did not contain %s", tc.token)
			}
		})
	}
}

func TestResetPasswordEmail(t *testing.T) {
	testCases := []struct {
		token  string
		webURL string
	}{
		{
			token:  "someRandomToken1",
			webURL: "http://localhost:3000",
		},
		{
			token:  "someRandomToken2",
			webURL: "http://localhost:3001",
		},
	}

	tmpl := NewTemplates()

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with WebURL %s", tc.webURL), func(t *testing.T) {
			dat := EmailResetPasswordTmplData{
				Token:  tc.token,
				WebURL: tc.webURL,
			}
			body, err := tmpl.Execute(EmailTypeResetPassword, EmailKindText, dat)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			if ok := strings.Contains(body, tc.webURL); !ok {
				t.Errorf("email body did not contain %s", tc.webURL)
			}
			if ok := strings.Contains(body, tc.token); !ok {
				t.Errorf("email body did not contain %s", tc.token)
			}
		})
	}
}
