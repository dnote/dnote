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

// Package mailer provides a functionality to send emails
package mailer

import (
	"bytes"
	"fmt"
	"io"
	ttemplate "text/template"

	"github.com/dnote/dnote/pkg/server/mailer/templates"
	"github.com/pkg/errors"
)

var (
	// EmailTypeResetPassword represents a reset password email
	EmailTypeResetPassword = "reset_password"
	// EmailTypeResetPasswordAlert represents a password change notification email
	EmailTypeResetPasswordAlert = "reset_password_alert"
	// EmailTypeWelcome represents an welcome email
	EmailTypeWelcome = "welcome"
)

var (
	// EmailKindText is the type of text email
	EmailKindText = "text/plain"
)

// tmpl is the common interface shared between Template from
// html/template and text/template
type tmpl interface {
	Execute(wr io.Writer, data interface{}) error
}

// template wraps a template with its subject line
type template struct {
	tmpl    tmpl
	subject string
}

// Templates holds the parsed email templates with their subjects
type Templates map[string]template

func getTemplateKey(name, kind string) string {
	return fmt.Sprintf("%s.%s", name, kind)
}

func (tmpl Templates) get(name, kind string) (template, error) {
	key := getTemplateKey(name, kind)
	t := tmpl[key]
	if t.tmpl == nil {
		return template{}, errors.Errorf("unsupported template '%s' with type '%s'", name, kind)
	}

	return t, nil
}

func (tmpl Templates) set(name, kind string, t tmpl, subject string) {
	key := getTemplateKey(name, kind)
	tmpl[key] = template{
		tmpl:    t,
		subject: subject,
	}
}

// NewTemplates initializes templates
func NewTemplates() Templates {
	welcomeText, err := initTextTmpl(EmailTypeWelcome)
	if err != nil {
		panic(errors.Wrap(err, "initializing welcome template"))
	}
	passwordResetText, err := initTextTmpl(EmailTypeResetPassword)
	if err != nil {
		panic(errors.Wrap(err, "initializing password reset template"))
	}
	passwordResetAlertText, err := initTextTmpl(EmailTypeResetPasswordAlert)
	if err != nil {
		panic(errors.Wrap(err, "initializing password reset template"))
	}

	T := Templates{}
	T.set(EmailTypeResetPassword, EmailKindText, passwordResetText, "Reset your Dnote password")
	T.set(EmailTypeResetPasswordAlert, EmailKindText, passwordResetAlertText, "Your Dnote password was changed")
	T.set(EmailTypeWelcome, EmailKindText, welcomeText, "Welcome to Dnote!")

	return T
}

// initTextTmpl returns a template instance by parsing the template with the given name
func initTextTmpl(templateName string) (tmpl, error) {
	filename := fmt.Sprintf("%s.txt", templateName)

	content, err := templates.Files.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "reading template")
	}

	t := ttemplate.New(templateName)
	if _, err = t.Parse(string(content)); err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}

	return t, nil
}

// Execute executes the template and returns the subject, body, and any error
func (tmpl Templates) Execute(name, kind string, data any) (subject, body string, err error) {
	t, err := tmpl.get(name, kind)
	if err != nil {
		return "", "", errors.Wrap(err, "getting template")
	}

	buf := new(bytes.Buffer)
	if err := t.tmpl.Execute(buf, data); err != nil {
		return "", "", errors.Wrap(err, "executing the template")
	}

	return t.subject, buf.String(), nil
}
