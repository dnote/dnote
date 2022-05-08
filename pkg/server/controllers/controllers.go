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
