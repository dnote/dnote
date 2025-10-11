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

package app

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestSendWelcomeEmail(t *testing.T) {
	emailBackend := testutils.MockEmailbackendImplementation{}
	a := NewTest()
	a.EmailBackend = &emailBackend
	a.WebURL = "http://example.com"

	if err := a.SendWelcomeEmail("alice@example.com"); err != nil {
		t.Fatal(err, "failed to perform")
	}

	assert.Equalf(t, len(emailBackend.Emails), 1, "email queue count mismatch")
	assert.Equal(t, emailBackend.Emails[0].From, "noreply@example.com", "email sender mismatch")
	assert.DeepEqual(t, emailBackend.Emails[0].To, []string{"alice@example.com"}, "email sender mismatch")

}

func TestSendPasswordResetEmail(t *testing.T) {
	emailBackend := testutils.MockEmailbackendImplementation{}
	a := NewTest()
	a.EmailBackend = &emailBackend
	a.WebURL = "http://example.com"

	if err := a.SendPasswordResetEmail("alice@example.com", "mockTokenValue"); err != nil {
		t.Fatal(err, "failed to perform")
	}

	assert.Equalf(t, len(emailBackend.Emails), 1, "email queue count mismatch")
	assert.Equal(t, emailBackend.Emails[0].From, "noreply@example.com", "email sender mismatch")
	assert.DeepEqual(t, emailBackend.Emails[0].To, []string{"alice@example.com"}, "email sender mismatch")

}

func TestGetSenderEmail(t *testing.T) {
	testCases := []struct {
		webURL         string
		expectedSender string
	}{
		{
			webURL:         "https://www.example.com",
			expectedSender: "noreply@example.com",
		},
		{
			webURL:         "https://www.example2.com",
			expectedSender: "alice@example2.com",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("web url %s", tc.webURL), func(t *testing.T) {
		})
	}
}
