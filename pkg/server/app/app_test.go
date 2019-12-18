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
				Config: Config{
					WebURL: "http://mock.url",
				},
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
				Config: Config{
					WebURL: "http://mock.url",
				},
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
				Config: Config{
					WebURL: "http://mock.url",
				},
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
				Config: Config{
					WebURL: "http://mock.url",
				},
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
				Config: Config{
					WebURL: "http://mock.url",
				},
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
				Config: Config{
					WebURL: "",
				},
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
