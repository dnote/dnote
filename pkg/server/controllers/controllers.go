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
