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

package app

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/pkg/errors"
)

var defaultSender = "admin@getdnote.com"

// GetSenderEmail returns the sender email
func GetSenderEmail(webURL, want string) (string, error) {
	addr, err := getNoreplySender(webURL)
	if err != nil {
		return "", errors.Wrap(err, "getting sender email address")
	}

	return addr, nil
}

func getDomainFromURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.Wrap(err, "parsing url")
	}

	host := u.Hostname()
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return host, nil
	}
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

	return domain, nil
}

func getNoreplySender(webURL string) (string, error) {
	domain, err := getDomainFromURL(webURL)
	if err != nil {
		return "", errors.Wrap(err, "parsing web url")
	}

	addr := fmt.Sprintf("noreply@%s", domain)
	return addr, nil
}

// SendWelcomeEmail sends welcome email
func (a *App) SendWelcomeEmail(email string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeWelcome, mailer.EmailKindText, mailer.WelcomeTmplData{
		AccountEmail: email,
		WebURL:       a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}

	from, err := GetSenderEmail(a.WebURL, defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Welcome to Dnote!", from, []string{email}, mailer.EmailKindText, body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}

// SendPasswordResetEmail sends password reset email
func (a *App) SendPasswordResetEmail(email, tokenValue string) error {
	if email == "" {
		return ErrEmailRequired
	}

	body, err := a.EmailTemplates.Execute(mailer.EmailTypeResetPassword, mailer.EmailKindText, mailer.EmailResetPasswordTmplData{
		AccountEmail: email,
		Token:        tokenValue,
		WebURL:       a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset password template for %s", email)
	}

	from, err := GetSenderEmail(a.WebURL, defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Reset your password", from, []string{email}, mailer.EmailKindText, body); err != nil {
		if errors.Cause(err) == mailer.ErrSMTPNotConfigured {
			return ErrInvalidSMTPConfig
		}

		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}

// SendPasswordResetAlertEmail sends email that notifies users of a password change
func (a *App) SendPasswordResetAlertEmail(email string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeResetPasswordAlert, mailer.EmailKindText, mailer.EmailResetPasswordAlertTmplData{
		AccountEmail: email,
		WebURL:       a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset password alert template for %s", email)
	}

	from, err := GetSenderEmail(a.WebURL, defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Dnote password changed", from, []string{email}, mailer.EmailKindText, body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}
