/* Copyright (C) 2019 Monomax Software Pty Ltd
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
	"os"
	"strings"
	"testing"

	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func init() {
	testutils.InitTestDB()

	templatePath := os.Getenv("DNOTE_TEST_EMAIL_TEMPLATE_DIR")
	InitTemplates(&templatePath)
}

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

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with WebURL %s", tc.webURL), func(t *testing.T) {
			m := NewEmail("alice@example.com", []string{"bob@example.com"}, "Test email")

			dat := EmailVerificationTmplData{
				Subject: "Test email verification email",
				Token:   tc.token,
				WebURL:  tc.webURL,
			}
			err := m.ParseTemplate(EmailTypeEmailVerification, dat)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			if ok := strings.Contains(m.Body, tc.webURL); !ok {
				t.Errorf("email body did not contain %s", tc.webURL)
			}
			if ok := strings.Contains(m.Body, tc.token); !ok {
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

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with WebURL %s", tc.webURL), func(t *testing.T) {
			m := NewEmail("alice@example.com", []string{"bob@example.com"}, "Test email")

			dat := EmailVerificationTmplData{
				Subject: "Test reset passowrd email",
				Token:   tc.token,
				WebURL:  tc.webURL,
			}
			err := m.ParseTemplate(EmailTypeResetPassword, dat)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			if ok := strings.Contains(m.Body, tc.webURL); !ok {
				t.Errorf("email body did not contain %s", tc.webURL)
			}
			if ok := strings.Contains(m.Body, tc.token); !ok {
				t.Errorf("email body did not contain %s", tc.token)
			}
		})
	}
}
