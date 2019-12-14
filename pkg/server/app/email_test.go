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

package app

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestSendVerificationEmail(t *testing.T) {
	testCases := []struct {
		onPremise     bool
		expectedSender string
	}{
		{
			onPremise:     false,
			expectedSender: "sung@getdnote.com",
		},
		{
			onPremise:     true,
			expectedSender: "noreply@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("self hosted %t", tc.onPremise), func(t *testing.T) {
			emailBackend := testutils.MockEmailbackendImplementation{}
			a := NewTest(&App{
				OnPremise:   tc.onPremise,
				WebURL:       "http://example.com",
				EmailBackend: &emailBackend,
			})

			if err := a.SendVerificationEmail("alice@example.com", "mockTokenValue"); err != nil {
				t.Fatal(err, "failed to perform")
			}

			assert.Equalf(t, len(emailBackend.Emails), 1, "email queue count mismatch")
			assert.Equal(t, emailBackend.Emails[0].From, tc.expectedSender, "email sender mismatch")
			assert.DeepEqual(t, emailBackend.Emails[0].To, []string{"alice@example.com"}, "email sender mismatch")
		})
	}
}

func TestSendWelcomeEmail(t *testing.T) {
	testCases := []struct {
		onPremise     bool
		expectedSender string
	}{
		{
			onPremise:     false,
			expectedSender: "sung@getdnote.com",
		},
		{
			onPremise:     true,
			expectedSender: "noreply@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("self hosted %t", tc.onPremise), func(t *testing.T) {
			emailBackend := testutils.MockEmailbackendImplementation{}
			a := NewTest(&App{
				OnPremise:   tc.onPremise,
				WebURL:       "http://example.com",
				EmailBackend: &emailBackend,
			})

			if err := a.SendWelcomeEmail("alice@example.com"); err != nil {
				t.Fatal(err, "failed to perform")
			}

			assert.Equalf(t, len(emailBackend.Emails), 1, "email queue count mismatch")
			assert.Equal(t, emailBackend.Emails[0].From, tc.expectedSender, "email sender mismatch")
			assert.DeepEqual(t, emailBackend.Emails[0].To, []string{"alice@example.com"}, "email sender mismatch")
		})
	}
}

func TestSendPasswordResetEmail(t *testing.T) {
	testCases := []struct {
		onPremise     bool
		expectedSender string
	}{
		{
			onPremise:     false,
			expectedSender: "sung@getdnote.com",
		},
		{
			onPremise:     true,
			expectedSender: "noreply@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("self hosted %t", tc.onPremise), func(t *testing.T) {
			emailBackend := testutils.MockEmailbackendImplementation{}
			a := NewTest(&App{
				OnPremise:   tc.onPremise,
				WebURL:       "http://example.com",
				EmailBackend: &emailBackend,
			})

			if err := a.SendPasswordResetEmail("alice@example.com", "mockTokenValue"); err != nil {
				t.Fatal(err, "failed to perform")
			}

			assert.Equalf(t, len(emailBackend.Emails), 1, "email queue count mismatch")
			assert.Equal(t, emailBackend.Emails[0].From, tc.expectedSender, "email sender mismatch")
			assert.DeepEqual(t, emailBackend.Emails[0].To, []string{"alice@example.com"}, "email sender mismatch")
		})
	}
}
