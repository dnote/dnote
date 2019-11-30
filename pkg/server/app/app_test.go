package app

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

func TestValidate(t *testing.T) {
	testCases := []struct {
		app         App
		expectedErr error
	}{
		{
			app: App{
				DB:               &gorm.DB{},
				Clock:            clock.NewMock(),
				StripeAPIBackend: nil,
				EmailTemplates:   mailer.Templates{},
				EmailBackend:     &testutils.MockEmailbackendImplementation{},
				WebURL:           "http://mock.url",
			},
			expectedErr: nil,
		},
		{
			app: App{
				DB:               nil,
				Clock:            clock.NewMock(),
				StripeAPIBackend: nil,
				EmailTemplates:   mailer.Templates{},
				EmailBackend:     &testutils.MockEmailbackendImplementation{},
				WebURL:           "http://mock.url",
			},
			expectedErr: ErrEmptyDB,
		},
		{
			app: App{
				DB:               &gorm.DB{},
				Clock:            nil,
				StripeAPIBackend: nil,
				EmailTemplates:   mailer.Templates{},
				EmailBackend:     &testutils.MockEmailbackendImplementation{},
				WebURL:           "http://mock.url",
			},
			expectedErr: ErrEmptyClock,
		},
		{
			app: App{
				DB:               &gorm.DB{},
				Clock:            clock.NewMock(),
				StripeAPIBackend: nil,
				EmailTemplates:   nil,
				EmailBackend:     &testutils.MockEmailbackendImplementation{},
				WebURL:           "http://mock.url",
			},
			expectedErr: ErrEmptyEmailTemplates,
		},
		{
			app: App{
				DB:               &gorm.DB{},
				Clock:            clock.NewMock(),
				StripeAPIBackend: nil,
				EmailTemplates:   mailer.Templates{},
				EmailBackend:     nil,
				WebURL:           "http://mock.url",
			},
			expectedErr: ErrEmptyEmailBackend,
		},
		{
			app: App{
				DB:               &gorm.DB{},
				Clock:            clock.NewMock(),
				StripeAPIBackend: nil,
				EmailTemplates:   mailer.Templates{},
				EmailBackend:     &testutils.MockEmailbackendImplementation{},
				WebURL:           "",
			},
			expectedErr: ErrEmptyWebURL,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", idx), func(t *testing.T) {
			err := tc.app.Validate()

			assert.Equal(t, errors.Cause(err), tc.expectedErr, "error mismatch")
		})
	}
}
