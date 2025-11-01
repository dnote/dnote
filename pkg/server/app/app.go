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
	// ErrEmptyBaseURL is an error for missing BaseURL content in the app configuration
	ErrEmptyBaseURL = errors.New("No BaseURL was provided")
	// ErrEmptyEmailBackend is an error for missing EmailBackend content in the app configuration
	ErrEmptyEmailBackend = errors.New("No EmailBackend was provided")
	// ErrEmptyHTTP500Page is an error for missing HTTP 500 page content
	ErrEmptyHTTP500Page = errors.New("No HTTP 500 error page was set")
)

// App is an application context
type App struct {
	DB                  *gorm.DB
	Clock               clock.Clock
	EmailBackend        mailer.Backend
	Files               map[string][]byte
	HTTP500Page         []byte
	BaseURL             string
	DisableRegistration bool
	Port                string
	DBPath              string
	AssetBaseURL        string
}

// Validate validates the app configuration
func (a *App) Validate() error {
	if a.BaseURL == "" {
		return ErrEmptyBaseURL
	}
	if a.Clock == nil {
		return ErrEmptyClock
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
