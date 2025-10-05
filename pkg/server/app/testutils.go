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
	"github.com/dnote/dnote/pkg/server/assets"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
)

// NewTest returns an app for a testing environment
func NewTest() App {
	return App{
		Clock:               clock.NewMock(),
		EmailTemplates:      mailer.NewTemplates(),
		EmailBackend:        &testutils.MockEmailbackendImplementation{},
		HTTP500Page:         assets.MustGetHTTP500ErrorPage(),
		AppEnv:              "TEST",
		WebURL:              "http://127.0.0.0.1",
		Port:                "3000",
		DisableRegistration: false,
		DBPath:              ":memory:",
		AssetBaseURL:        "",
	}
}
