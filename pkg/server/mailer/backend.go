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
	Queue(subject, from string, to []string, contentType, body string) error
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
type DefaultBackend struct {
	Dialer  EmailDialer
	Enabled bool
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
func NewDefaultBackend(enabled bool) (*DefaultBackend, error) {
	p, err := getSMTPParams()
	if err != nil {
		return nil, err
	}

	d := gomail.NewDialer(p.Host, p.Port, p.Username, p.Password)

	return &DefaultBackend{
		Dialer:  &gomailDialer{Dialer: d},
		Enabled: enabled,
	}, nil
}

// Queue is an implementation of Backend.Queue.
func (b *DefaultBackend) Queue(subject, from string, to []string, contentType, body string) error {
	// If not enabled, just log the email
	if !b.Enabled {
		log.WithFields(log.Fields{
			"subject": subject,
			"to":      to,
			"from":    from,
			"body":    body,
		}).Info("Not sending email because email backend is not configured.")
		return nil
	}

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
