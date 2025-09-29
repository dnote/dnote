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

package controllers

import (
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/views"
)

// Controllers is a group of controllers
type Controllers struct {
	Users  *Users
	Notes  *Notes
	Books  *Books
	Sync   *Sync
	Static *Static
	Health *Health
}

// New returns a new group of controllers
func New(app *app.App) *Controllers {
	c := Controllers{}

	viewEngine := views.NewDefaultEngine()

	c.Users = NewUsers(app, viewEngine)
	c.Notes = NewNotes(app)
	c.Books = NewBooks(app)
	c.Sync = NewSync(app)
	c.Static = NewStatic(app, viewEngine)
	c.Health = NewHealth(app)

	return &c
}
