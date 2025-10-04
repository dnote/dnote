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
	"testing"

	"gopkg.in/gomail.v2"
)

type mockDialer struct {
	sentMessages []*gomail.Message
	err          error
}

func (m *mockDialer) DialAndSend(msgs ...*gomail.Message) error {
	m.sentMessages = append(m.sentMessages, msgs...)
	return m.err
}

func TestDefaultBackendQueue(t *testing.T) {
	t.Run("enabled sends email", func(t *testing.T) {
		mock := &mockDialer{}
		backend := &DefaultBackend{
			Dialer:  mock,
			Enabled: true,
		}

		err := backend.Queue("Test Subject", "alice@example.com", []string{"bob@example.com"}, "text/plain", "Test body")
		if err != nil {
			t.Fatalf("Queue failed: %v", err)
		}

		if len(mock.sentMessages) != 1 {
			t.Errorf("expected 1 message sent, got %d", len(mock.sentMessages))
		}
	})

	t.Run("disabled does not send email", func(t *testing.T) {
		mock := &mockDialer{}
		backend := &DefaultBackend{
			Dialer:  mock,
			Enabled: false,
		}

		err := backend.Queue("Test Subject", "alice@example.com", []string{"bob@example.com"}, "text/plain", "Test body")
		if err != nil {
			t.Fatalf("Queue failed: %v", err)
		}

		if len(mock.sentMessages) != 0 {
			t.Errorf("expected 0 messages sent when disabled, got %d", len(mock.sentMessages))
		}
	})
}

func TestNewDefaultBackend(t *testing.T) {
	t.Run("with all env vars set", func(t *testing.T) {
		t.Setenv("SmtpHost", "smtp.example.com")
		t.Setenv("SmtpPort", "587")
		t.Setenv("SmtpUsername", "user@example.com")
		t.Setenv("SmtpPassword", "secret")

		backend, err := NewDefaultBackend(true)
		if err != nil {
			t.Fatalf("NewDefaultBackend failed: %v", err)
		}

		if backend.Enabled != true {
			t.Errorf("expected Enabled to be true, got %v", backend.Enabled)
		}
		if backend.Dialer == nil {
			t.Error("expected Dialer to be set")
		}
	})

	t.Run("missing SMTP config returns error", func(t *testing.T) {
		t.Setenv("SmtpHost", "")
		t.Setenv("SmtpPort", "")
		t.Setenv("SmtpUsername", "")
		t.Setenv("SmtpPassword", "")

		_, err := NewDefaultBackend(true)
		if err == nil {
			t.Error("expected error when SMTP not configured")
		}
		if err != ErrSMTPNotConfigured {
			t.Errorf("expected ErrSMTPNotConfigured, got %v", err)
		}
	})
}
