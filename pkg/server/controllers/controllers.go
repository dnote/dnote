package controllers

import (
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
)

// Controllers is a group of controllers
type Controllers struct {
	Users  *Users
	Static *Static
}

// New returns a new group of controllers
func New(cfg config.Config, app *app.App) *Controllers {
	c := Controllers{}

	c.Users = NewUsers(cfg, app)
	c.Static = NewStatic(cfg)

	return &c
}
