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
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/mailer"
	"gorm.io/gorm"
	"github.com/pkg/errors"
)

var (
	// ErrEmptyDB is an error for missing database connection in the app configuration
	ErrEmptyDB = errors.New("No database connection was provided")
	// ErrEmptyClock is an error for missing clock in the app configuration
	ErrEmptyClock = errors.New("No clock was provided")
	// ErrEmptyWebURL is an error for missing WebURL content in the app configuration
	ErrEmptyWebURL = errors.New("No WebURL was provided")
	// ErrEmptyEmailTemplates is an error for missing EmailTemplates content in the app configuration
	ErrEmptyEmailTemplates = errors.New("No EmailTemplate store was provided")
	// ErrEmptyEmailBackend is an error for missing EmailBackend content in the app configuration
	ErrEmptyEmailBackend = errors.New("No EmailBackend was provided")
	// ErrEmptyHTTP500Page is an error for missing HTTP 500 page content
	ErrEmptyHTTP500Page = errors.New("No HTTP 500 error page was set")
)

// App is an application context
type App struct {
	DB                  *gorm.DB
	Clock               clock.Clock
	EmailTemplates      mailer.Templates
	EmailBackend        mailer.Backend
	Files               map[string][]byte
	HTTP500Page         []byte
	AppEnv              string
	WebURL              string
	DisableRegistration bool
	Port                string
	DBPath              string
	AssetBaseURL        string
}

// Validate validates the app configuration
func (a *App) Validate() error {
	if a.WebURL == "" {
		return ErrEmptyWebURL
	}
	if a.Clock == nil {
		return ErrEmptyClock
	}
	if a.EmailTemplates == nil {
		return ErrEmptyEmailTemplates
	}
	if a.EmailBackend == nil {
		return ErrEmptyEmailBackend
	}
	if a.DB == nil {
		return ErrEmptyDB
	}
	if a.HTTP500Page == nil {
		return ErrEmptyHTTP500Page
	}

	return nil
}
