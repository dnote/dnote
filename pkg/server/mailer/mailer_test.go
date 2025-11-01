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

package mailer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestAllTemplatesInitialized(t *testing.T) {
	tmpl := NewTemplates()

	emailTypes := []string{
		EmailTypeResetPassword,
		EmailTypeResetPasswordAlert,
		EmailTypeWelcome,
	}

	for _, emailType := range emailTypes {
		t.Run(emailType, func(t *testing.T) {
			_, err := tmpl.get(emailType, EmailKindText)
			if err != nil {
				t.Errorf("template %s not initialized: %v", emailType, err)
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
			subject, body, err := tmpl.Execute(EmailTypeResetPassword, EmailKindText, dat)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			if subject != "Reset your Dnote password" {
				t.Errorf("expected subject 'Reset your Dnote password', got '%s'", subject)
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

func TestWelcomeEmail(t *testing.T) {
	testCases := []struct {
		accountEmail string
		webURL       string
	}{
		{
			accountEmail: "test@example.com",
			webURL:       "http://localhost:3000",
		},
		{
			accountEmail: "user@example.org",
			webURL:       "http://localhost:3001",
		},
	}

	tmpl := NewTemplates()

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with WebURL %s and email %s", tc.webURL, tc.accountEmail), func(t *testing.T) {
			dat := WelcomeTmplData{
				AccountEmail: tc.accountEmail,
				WebURL:       tc.webURL,
			}
			subject, body, err := tmpl.Execute(EmailTypeWelcome, EmailKindText, dat)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			if subject != "Welcome to Dnote!" {
				t.Errorf("expected subject 'Welcome to Dnote!', got '%s'", subject)
			}
			if ok := strings.Contains(body, tc.webURL); !ok {
				t.Errorf("email body did not contain %s", tc.webURL)
			}
			if ok := strings.Contains(body, tc.accountEmail); !ok {
				t.Errorf("email body did not contain %s", tc.accountEmail)
			}
		})
	}
}

func TestResetPasswordAlertEmail(t *testing.T) {
	testCases := []struct {
		accountEmail string
		webURL       string
	}{
		{
			accountEmail: "test@example.com",
			webURL:       "http://localhost:3000",
		},
		{
			accountEmail: "user@example.org",
			webURL:       "http://localhost:3001",
		},
	}

	tmpl := NewTemplates()

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with WebURL %s and email %s", tc.webURL, tc.accountEmail), func(t *testing.T) {
			dat := EmailResetPasswordAlertTmplData{
				AccountEmail: tc.accountEmail,
				WebURL:       tc.webURL,
			}
			subject, body, err := tmpl.Execute(EmailTypeResetPasswordAlert, EmailKindText, dat)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			if subject != "Your Dnote password was changed" {
				t.Errorf("expected subject 'Your Dnote password was changed', got '%s'", subject)
			}
			if ok := strings.Contains(body, tc.webURL); !ok {
				t.Errorf("email body did not contain %s", tc.webURL)
			}
			if ok := strings.Contains(body, tc.accountEmail); !ok {
				t.Errorf("email body did not contain %s", tc.accountEmail)
			}
		})
	}
}
