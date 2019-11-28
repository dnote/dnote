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

// Package mailer provides a functionality to send emails
package mailer

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/aymerick/douceur/inliner"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

var (
	// EmailTypeResetPassword represents a reset password email
	EmailTypeResetPassword = "reset_password"
	// EmailTypeWeeklyDigest represents a weekly digest email
	EmailTypeWeeklyDigest = "digest"
	// EmailTypeEmailVerification represents an email verification email
	EmailTypeEmailVerification = "email_verification"
)

// Templates holds the parsed email templates
type Templates map[string]*template.Template

// NewTemplates initializes templates
func NewTemplates(srcDir *string) Templates {
	var box *packr.Box

	if srcDir != nil {
		box = packr.Folder(*srcDir)
	} else {
		box = packr.New("emailTemplates", "./templates/src")
	}

	weeklyDigestTmpl, err := initTemplate(box, EmailTypeWeeklyDigest)
	if err != nil {
		panic(errors.Wrap(err, "initializing weekly digest template"))
	}
	emailVerificationTmpl, err := initTemplate(box, EmailTypeEmailVerification)
	if err != nil {
		panic(errors.Wrap(err, "initializing email verification template"))
	}
	passwowrdResetTmpl, err := initTemplate(box, EmailTypeResetPassword)
	if err != nil {
		panic(errors.Wrap(err, "initializing password reset template"))
	}

	T := map[string]*template.Template{}
	T[EmailTypeWeeklyDigest] = weeklyDigestTmpl
	T[EmailTypeEmailVerification] = emailVerificationTmpl
	T[EmailTypeResetPassword] = passwowrdResetTmpl

	return T
}

// initTemplate returns a template instance by parsing the template with the
// given name along with partials
func initTemplate(box *packr.Box, templateName string) (*template.Template, error) {
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

	t := template.New(templateName)
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

// Execute executes the template with the given name with the givn data, and inlines
// CSS rules.
func (tmpl Templates) Execute(templateName string, data interface{}) (string, error) {
	t := tmpl[templateName]
	if t == nil {
		return "", errors.Errorf("unsupported template '%s'", templateName)
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", errors.Wrap(err, "executing the template")
	}

	html, err := inliner.Inline(buf.String())
	if err != nil {
		return "", errors.Wrap(err, "inlining the css rules")
	}

	return html, nil
}
