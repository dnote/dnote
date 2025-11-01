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
	"os"
	"strconv"

	"github.com/dnote/dnote/pkg/server/log"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

// ErrSMTPNotConfigured is an error indicating that SMTP is not configured
var ErrSMTPNotConfigured = errors.New("SMTP is not configured")

// Backend is an interface for sending emails.
type Backend interface {
	SendEmail(templateType, from string, to []string, data interface{}) error
}

// EmailDialer is an interface for sending email messages
type EmailDialer interface {
	DialAndSend(m ...*gomail.Message) error
}

// gomailDialer wraps gomail.Dialer to implement EmailDialer interface
type gomailDialer struct {
	*gomail.Dialer
}

// DefaultBackend is an implementation of the Backend
// that sends an email without queueing.
// This backend is always enabled and will send emails via SMTP.
type DefaultBackend struct {
	Dialer    EmailDialer
	Templates Templates
}

type dialerParams struct {
	Host     string
	Port     int
	Username string
	Password string
}

func getSMTPParams() (*dialerParams, error) {
	portEnv := os.Getenv("SmtpPort")
	hostEnv := os.Getenv("SmtpHost")
	usernameEnv := os.Getenv("SmtpUsername")
	passwordEnv := os.Getenv("SmtpPassword")

	if portEnv == "" || hostEnv == "" || usernameEnv == "" || passwordEnv == "" {
		return nil, ErrSMTPNotConfigured
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		return nil, errors.Wrap(err, "parsing SMTP port")
	}

	p := &dialerParams{
		Host:     hostEnv,
		Port:     port,
		Username: usernameEnv,
		Password: passwordEnv,
	}

	return p, nil
}

// NewDefaultBackend creates a default backend
func NewDefaultBackend() (*DefaultBackend, error) {
	p, err := getSMTPParams()
	if err != nil {
		return nil, err
	}

	d := gomail.NewDialer(p.Host, p.Port, p.Username, p.Password)

	return &DefaultBackend{
		Dialer:    &gomailDialer{Dialer: d},
		Templates: NewTemplates(),
	}, nil
}

// SendEmail is an implementation of Backend.SendEmail.
// It renders the template and sends the email immediately via SMTP.
func (b *DefaultBackend) SendEmail(templateType, from string, to []string, data interface{}) error {
	subject, body, err := b.Templates.Execute(templateType, EmailKindText, data)
	if err != nil {
		return errors.Wrap(err, "executing template")
	}

	return b.queue(subject, from, to, EmailKindText, body)
}

// queue sends the email immediately via SMTP.
func (b *DefaultBackend) queue(subject, from string, to []string, contentType, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, body)

	if err := b.Dialer.DialAndSend(m); err != nil {
		return errors.Wrap(err, "dialing and sending email")
	}

	return nil
}

// StdoutBackend is an implementation of the Backend
// that prints emails to stdout instead of sending them.
// This is useful for development and testing.
type StdoutBackend struct {
	Templates Templates
}

// NewStdoutBackend creates a stdout backend
func NewStdoutBackend() *StdoutBackend {
	return &StdoutBackend{
		Templates: NewTemplates(),
	}
}

// SendEmail is an implementation of Backend.SendEmail.
// It renders the template and logs the email to stdout instead of sending it.
func (b *StdoutBackend) SendEmail(templateType, from string, to []string, data interface{}) error {
	subject, body, err := b.Templates.Execute(templateType, EmailKindText, data)
	if err != nil {
		return errors.Wrap(err, "executing template")
	}

	log.WithFields(log.Fields{
		"subject": subject,
		"to":      to,
		"from":    from,
		"body":    body,
	}).Info("Email (not sent, using StdoutBackend)")
	return nil
}
