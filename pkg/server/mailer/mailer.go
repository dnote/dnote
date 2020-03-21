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

// Package mailer provides a functionality to send emails
package mailer

import (
	"bytes"
	"fmt"
	htemplate "html/template"
	"io"
	ttemplate "text/template"

	"github.com/aymerick/douceur/inliner"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

var (
	// EmailTypeResetPassword represents a reset password email
	EmailTypeResetPassword = "reset_password"
	// EmailTypeResetPasswordAlert represents a password change notification email
	EmailTypeResetPasswordAlert = "reset_password_alert"
	// EmailTypeEmailVerification represents an email verification email
	EmailTypeEmailVerification = "verify_email"
	// EmailTypeWelcome represents an welcome email
	EmailTypeWelcome = "welcome"
	// EmailTypeInactiveReminder represents an inactivity reminder email
	EmailTypeInactiveReminder = "inactive"
	// EmailTypeSubscriptionConfirmation represents an inactivity reminder email
	EmailTypeSubscriptionConfirmation = "subscription_confirmation"
)

var (
	// EmailKindHTML is the type of html email
	EmailKindHTML = "text/html"
	// EmailKindText is the type of text email
	EmailKindText = "text/plain"
)

// template is the common interface shared between Template from
// html/template and text/template
type template interface {
	Execute(wr io.Writer, data interface{}) error
}

// Templates holds the parsed email templates
type Templates map[string]template

func getTemplateKey(name, kind string) string {
	return fmt.Sprintf("%s.%s", name, kind)
}

func (tmpl Templates) get(name, kind string) (template, error) {
	key := getTemplateKey(name, kind)
	t := tmpl[key]
	if t == nil {
		return nil, errors.Errorf("unsupported template '%s' with type '%s'", name, kind)
	}

	return t, nil
}

func (tmpl Templates) set(name, kind string, t template) {
	key := getTemplateKey(name, kind)
	tmpl[key] = t
}

// NewTemplates initializes templates
func NewTemplates(srcDir *string) Templates {
	var box *packr.Box

	if srcDir != nil {
		box = packr.Folder(*srcDir)
	} else {
		box = packr.New("emailTemplates", "./templates/src")
	}

	welcomeText, err := initTextTmpl(box, EmailTypeWelcome)
	if err != nil {
		panic(errors.Wrap(err, "initializing welcome template"))
	}
	verifyEmailText, err := initTextTmpl(box, EmailTypeEmailVerification)
	if err != nil {
		panic(errors.Wrap(err, "initializing email verification template"))
	}
	passwordResetText, err := initTextTmpl(box, EmailTypeResetPassword)
	if err != nil {
		panic(errors.Wrap(err, "initializing password reset template"))
	}
	passwordResetAlertText, err := initTextTmpl(box, EmailTypeResetPasswordAlert)
	if err != nil {
		panic(errors.Wrap(err, "initializing password reset template"))
	}
	inactiveReminderText, err := initTextTmpl(box, EmailTypeInactiveReminder)
	if err != nil {
		panic(errors.Wrap(err, "initializing password reset template"))
	}
	subscriptionConfirmationText, err := initTextTmpl(box, EmailTypeSubscriptionConfirmation)
	if err != nil {
		panic(errors.Wrap(err, "initializing password reset template"))
	}

	T := Templates{}
	T.set(EmailTypeResetPassword, EmailKindText, passwordResetText)
	T.set(EmailTypeResetPasswordAlert, EmailKindText, passwordResetAlertText)
	T.set(EmailTypeEmailVerification, EmailKindText, verifyEmailText)
	T.set(EmailTypeWelcome, EmailKindText, welcomeText)
	T.set(EmailTypeInactiveReminder, EmailKindText, inactiveReminderText)
	T.set(EmailTypeSubscriptionConfirmation, EmailKindText, subscriptionConfirmationText)

	return T
}

// initHTMLTmpl returns a template instance by parsing the template with the
// given name along with partials
func initHTMLTmpl(box *packr.Box, templateName string) (template, error) {
	filename := fmt.Sprintf("%s.html", templateName)

	content, err := box.FindString(filename)
	if err != nil {
		return nil, errors.Wrap(err, "reading template")
	}
	headerContent, err := box.FindString("header.html")
	if err != nil {
		return nil, errors.Wrap(err, "reading header template")
	}
	footerContent, err := box.FindString("footer.html")
	if err != nil {
		return nil, errors.Wrap(err, "reading footer template")
	}

	t := htemplate.New(templateName)
	if _, err = t.Parse(content); err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}
	if _, err = t.Parse(headerContent); err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}
	if _, err = t.Parse(footerContent); err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}

	return t, nil
}

// initTextTmpl returns a template instance by parsing the template with the given name
func initTextTmpl(box *packr.Box, templateName string) (template, error) {
	filename := fmt.Sprintf("%s.txt", templateName)

	content, err := box.FindString(filename)
	if err != nil {
		return nil, errors.Wrap(err, "reading template")
	}

	t := ttemplate.New(templateName)
	if _, err = t.Parse(content); err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}

	return t, nil
}

// Execute executes the template with the given name with the givn data
func (tmpl Templates) Execute(name, kind string, data interface{}) (string, error) {
	t, err := tmpl.get(name, kind)
	if err != nil {
		return "", errors.Wrap(err, "getting template")
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", errors.Wrap(err, "executing the template")
	}

	// If HTML email, inline the CSS rules
	if kind == EmailKindHTML {
		html, err := inliner.Inline(buf.String())
		if err != nil {
			return "", errors.Wrap(err, "inlining the css rules")
		}

		return html, nil
	}

	return buf.String(), nil
}
