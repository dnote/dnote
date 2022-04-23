/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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
	"os"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
)

// NewTest returns an app for a testing environment
func NewTest(appParams *App) App {
	emailTmplDir := os.Getenv("DNOTE_TEST_EMAIL_TEMPLATE_DIR")
	c := config.Load()
	c.SetOnPremise(false)

	a := App{
		DB:             testutils.DB,
		Clock:          clock.NewMock(),
		EmailTemplates: mailer.NewTemplates(&emailTmplDir),
		EmailBackend:   &testutils.MockEmailbackendImplementation{},
		Config:         c,
	}

	// Allow to override with appParams
	if appParams != nil && appParams.EmailBackend != nil {
		a.EmailBackend = appParams.EmailBackend
	}
	if appParams != nil && appParams.Clock != nil {
		a.Clock = appParams.Clock
	}
	if appParams != nil && appParams.EmailTemplates != nil {
		a.EmailTemplates = appParams.EmailTemplates
	}
	if appParams != nil && appParams.Config.OnPremise {
		a.Config.OnPremise = appParams.Config.OnPremise
	}
	if appParams != nil && appParams.Config.WebURL != "" {
		a.Config.WebURL = appParams.Config.WebURL
	}
	if appParams != nil && appParams.Config.DisableRegistration {
		a.Config.DisableRegistration = appParams.Config.DisableRegistration
	}
	if appParams != nil && appParams.Config.PageTemplateDir != "" {
		a.Config.PageTemplateDir = appParams.Config.PageTemplateDir
	}

	fmt.Printf("%+v\n", appParams)
	fmt.Printf("%+v\n", a)

	return a
}
