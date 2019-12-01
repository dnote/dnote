package app

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestSendVerificationEmail(t *testing.T) {
	testCases := []struct {
		selfHosted     bool
		expectedSender string
	}{
		{
			selfHosted:     false,
			expectedSender: "sung@getdnote.com",
		},
		{
			selfHosted:     true,
			expectedSender: "noreply@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("self hosted %t", tc.selfHosted), func(t *testing.T) {
			emailBackend := testutils.MockEmailbackendImplementation{}
			a := NewTest(&App{
				SelfHosted:   tc.selfHosted,
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
		selfHosted     bool
		expectedSender string
	}{
		{
			selfHosted:     false,
			expectedSender: "sung@getdnote.com",
		},
		{
			selfHosted:     true,
			expectedSender: "noreply@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("self hosted %t", tc.selfHosted), func(t *testing.T) {
			emailBackend := testutils.MockEmailbackendImplementation{}
			a := NewTest(&App{
				SelfHosted:   tc.selfHosted,
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
		selfHosted     bool
		expectedSender string
	}{
		{
			selfHosted:     false,
			expectedSender: "sung@getdnote.com",
		},
		{
			selfHosted:     true,
			expectedSender: "noreply@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("self hosted %t", tc.selfHosted), func(t *testing.T) {
			emailBackend := testutils.MockEmailbackendImplementation{}
			a := NewTest(&App{
				SelfHosted:   tc.selfHosted,
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
