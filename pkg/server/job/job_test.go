package job

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func TestNewRunner(t *testing.T) {
	testCases := []struct {
		db           *gorm.DB
		clock        clock.Clock
		emailTmpl    mailer.Templates
		emailBackend mailer.Backend
		webURL       string
		expectedErr  error
	}{
		{
			db:           &gorm.DB{},
			clock:        clock.NewMock(),
			emailTmpl:    mailer.Templates{},
			emailBackend: &testutils.MockEmailbackendImplementation{},
			webURL:       "http://mock.url",
			expectedErr:  nil,
		},
		{
			db:           nil,
			clock:        clock.NewMock(),
			emailTmpl:    mailer.Templates{},
			emailBackend: &testutils.MockEmailbackendImplementation{},
			webURL:       "http://mock.url",
			expectedErr:  ErrEmptyDB,
		},
		{
			db:           &gorm.DB{},
			clock:        nil,
			emailTmpl:    mailer.Templates{},
			emailBackend: &testutils.MockEmailbackendImplementation{},
			webURL:       "http://mock.url",
			expectedErr:  ErrEmptyClock,
		},
		{
			db:           &gorm.DB{},
			clock:        clock.NewMock(),
			emailTmpl:    nil,
			emailBackend: &testutils.MockEmailbackendImplementation{},
			webURL:       "http://mock.url",
			expectedErr:  ErrEmptyEmailTemplates,
		},
		{
			db:           &gorm.DB{},
			clock:        clock.NewMock(),
			emailTmpl:    mailer.Templates{},
			emailBackend: nil,
			webURL:       "http://mock.url",
			expectedErr:  ErrEmptyEmailBackend,
		},
		{
			db:           &gorm.DB{},
			clock:        clock.NewMock(),
			emailTmpl:    mailer.Templates{},
			emailBackend: &testutils.MockEmailbackendImplementation{},
			webURL:       "",
			expectedErr:  ErrEmptyWebURL,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			_, err := NewRunner(tc.db, tc.clock, tc.emailTmpl, tc.emailBackend, tc.webURL)

			assert.Equal(t, errors.Cause(err), tc.expectedErr, "error mismatch")
		})
	}
}
