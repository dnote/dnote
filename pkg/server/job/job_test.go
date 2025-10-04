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

package job

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
	"gorm.io/gorm"
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

			c := config.Load()
			c.WebURL = tc.webURL

			_, err := NewRunner(tc.db, tc.clock, tc.emailTmpl, tc.emailBackend, c)

			assert.Equal(t, errors.Cause(err), tc.expectedErr, "error mismatch")
		})
	}
}
