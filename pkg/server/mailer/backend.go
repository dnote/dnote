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
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

// ErrSMTPNotConfigured is an error indicating that SMTP is not configured
var ErrSMTPNotConfigured = errors.New("SMTP is not configured")

// Backend is an interface for sending emails.
type Backend interface {
	Queue(subject, from string, to []string, contentType, body string) error
}

// SimpleBackendImplementation is an implementation of the Backend
// that sends an email without queueing.
type SimpleBackendImplementation struct {
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

// Queue is an implementation of Backend.Queue.
func (b *SimpleBackendImplementation) Queue(subject, from string, to []string, contentType, body string) error {
	// If not production, never actually send an email
	if os.Getenv("GO_ENV") != "PRODUCTION" {
		log.Println("Not sending email because Dnote is not running in a production environment.")
		log.Printf("Subject: %s, to: %s, from: %s", subject, to, from)
		fmt.Println(body)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, body)

	p, err := getSMTPParams()
	if err != nil {
		return errors.Wrap(err, "getting dialer params")
	}

	d := gomail.NewPlainDialer(p.Host, p.Port, p.Username, p.Password)
	if err := d.DialAndSend(m); err != nil {
		return errors.Wrap(err, "dialing and sending email")
	}

	return nil
}
