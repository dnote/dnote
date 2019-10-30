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
	"os"
	"path"
	"strconv"

	"github.com/aymerick/douceur/inliner"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

// Email represents email to be sent out
type Email struct {
	from    string
	to      []string
	subject string
	Body    string
}

var (
	// T is a map of templates
	T = map[string]*template.Template{}
	// EmailTypeResetPassword represents a reset password email
	EmailTypeResetPassword = "reset_password"
	// EmailTypeWeeklyDigest represents a weekly digest email
	EmailTypeWeeklyDigest = "digest"
	// EmailTypeEmailVerification represents an email verification email
	EmailTypeEmailVerification = "email_verification"
)

func getTemplatePath(templateDirPath, filename string) string {
	return path.Join(templateDirPath, fmt.Sprintf("%s.html", filename))
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

// InitTemplates initializes templates
func InitTemplates(srcDir *string) {
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

	T[EmailTypeWeeklyDigest] = weeklyDigestTmpl
	T[EmailTypeEmailVerification] = emailVerificationTmpl
	T[EmailTypeResetPassword] = passwowrdResetTmpl
}

// NewEmail returns a pointer to an Email struct with the given data
func NewEmail(from string, to []string, subject string) *Email {
	return &Email{
		from:    from,
		to:      to,
		subject: subject,
	}
}

type dialerParams struct {
	Host     string
	Port     int
	Username string
	Password string
}

func getSMTPParams() (*dialerParams, error) {
	portStr := os.Getenv("SmtpPort")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.Wrap(err, "parsing SMTP port")
	}

	p := &dialerParams{
		Host:     os.Getenv("SmtpHost"),
		Port:     port,
		Username: os.Getenv("SmtpUsername"),
		Password: os.Getenv("SmtpPassword"),
	}

	return p, nil
}

// Send sends the email
func (e *Email) Send() error {
	// If not production, never actually send an email
	if os.Getenv("GO_ENV") != "PRODUCTION" {
		fmt.Println("Not sending email because not production")
		fmt.Println(e.subject, e.to, e.from)
		fmt.Println("Body", e.Body)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", e.to...)
	m.SetHeader("Subject", e.subject)
	m.SetBody("text/html", e.Body)

	p, err := getSMTPParams()
	if err != nil {
		return errors.Wrap(err, "getting dialer params")
	}

	d := gomail.NewPlainDialer(p.Host, p.Port, p.Username, p.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// ParseTemplate sets the email body by parsing the file at the given path,
// evaluating all partials and inlining CSS rules
func (e *Email) ParseTemplate(templateName string, data interface{}) error {
	t := T[templateName]
	if t == nil {
		return errors.Errorf("unsupported template '%s'", templateName)
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return errors.Wrap(err, "executing the template")
	}

	html, err := inliner.Inline(buf.String())
	if err != nil {
		return errors.Wrap(err, "inlining the css rules")
	}

	e.Body = html
	return nil
}
