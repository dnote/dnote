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

package app

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/pkg/errors"
)

var defaultSender = "admin@getdnote.com"

// GetSenderEmail returns the sender email
func GetSenderEmail(c config.Config, want string) (string, error) {
	if !c.OnPremises {
		return want, nil
	}

	addr, err := getNoreplySender(c)
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

func getNoreplySender(c config.Config) (string, error) {
	domain, err := getDomainFromURL(c.WebURL)
	if err != nil {
		return "", errors.Wrap(err, "parsing web url")
	}

	addr := fmt.Sprintf("noreply@%s", domain)
	return addr, nil
}

// SendVerificationEmail sends verification email
func (a *App) SendVerificationEmail(email, tokenValue string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeEmailVerification, mailer.EmailKindText, mailer.EmailVerificationTmplData{
		Token:  tokenValue,
		WebURL: a.Config.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}

	from, err := GetSenderEmail(a.Config, defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Verify your Dnote email address", from, []string{email}, mailer.EmailKindText, body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}

// SendWelcomeEmail sends welcome email
func (a *App) SendWelcomeEmail(email string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeWelcome, mailer.EmailKindText, mailer.WelcomeTmplData{
		AccountEmail: email,
		WebURL:       a.Config.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}

	from, err := GetSenderEmail(a.Config, defaultSender)
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
		WebURL:       a.Config.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset password template for %s", email)
	}

	from, err := GetSenderEmail(a.Config, defaultSender)
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
		WebURL:       a.Config.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset password alert template for %s", email)
	}

	from, err := GetSenderEmail(a.Config, defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Dnote password changed", from, []string{email}, mailer.EmailKindText, body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}
