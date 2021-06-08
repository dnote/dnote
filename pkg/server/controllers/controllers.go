package controllers

import (
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/log"
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
	log.Info(app.Config.PageTemplateDir)

	c := Controllers{}

	c.Users = NewUsers(app)
	c.Notes = NewNotes(app)
	c.Books = NewBooks(app)
	c.Sync = NewSync(app)
	c.Static = NewStatic(app)
	c.Health = NewHealth(app)

	return &c
}
