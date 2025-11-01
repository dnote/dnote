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
