package mailer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func init() {
	testutils.InitTestDB()

	templatePath := fmt.Sprintf("%s/mailer/templates/src", testutils.ServerPath)
	InitTemplates(&templatePath)
}

func TestEmailVerificationEmail(t *testing.T) {
	testCases := []struct {
		token  string
		webURL string
	}{
		{
			token:  "someRandomToken1",
			webURL: "http://localhost:3000",
		},
		{
			token:  "someRandomToken2",
			webURL: "http://localhost:3001",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with WebURL %s", tc.webURL), func(t *testing.T) {
			m := NewEmail("alice@example.com", []string{"bob@example.com"}, "Test email")

			dat := EmailVerificationTmplData{
				Subject: "Test email verification email",
				Token:   tc.token,
				WebURL:  tc.webURL,
			}
			err := m.ParseTemplate(EmailTypeEmailVerification, dat)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			if ok := strings.Contains(m.Body, tc.webURL); !ok {
				t.Errorf("email body did not contain %s", tc.webURL)
			}
			if ok := strings.Contains(m.Body, tc.token); !ok {
				t.Errorf("email body did not contain %s", tc.token)
			}
		})
	}
}

func TestResetPasswordEmail(t *testing.T) {
	testCases := []struct {
		token  string
		webURL string
	}{
		{
			token:  "someRandomToken1",
			webURL: "http://localhost:3000",
		},
		{
			token:  "someRandomToken2",
			webURL: "http://localhost:3001",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with WebURL %s", tc.webURL), func(t *testing.T) {
			m := NewEmail("alice@example.com", []string{"bob@example.com"}, "Test email")

			dat := EmailVerificationTmplData{
				Subject: "Test reset passowrd email",
				Token:   tc.token,
				WebURL:  tc.webURL,
			}
			err := m.ParseTemplate(EmailTypeResetPassword, dat)
			if err != nil {
				t.Fatal(errors.Wrap(err, "executing"))
			}

			if ok := strings.Contains(m.Body, tc.webURL); !ok {
				t.Errorf("email body did not contain %s", tc.webURL)
			}
			if ok := strings.Contains(m.Body, tc.token); !ok {
				t.Errorf("email body did not contain %s", tc.token)
			}
		})
	}
}
