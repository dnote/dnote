package controllers

import (
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/jinzhu/gorm"
)

// Controllers is a group of controllers
type Controllers struct {
	Users *Users
}

// New returns a new group of controllers
func New(cfg config.Config, db *gorm.DB, cl clock.Clock) *Controllers {
	c := Controllers{}

	c.Users = NewUsers(cfg, db)

	return &c
}
