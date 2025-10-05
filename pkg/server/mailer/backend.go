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
