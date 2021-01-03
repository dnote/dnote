package controllers

import (
	"net/http"

	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/views"
	"github.com/jinzhu/gorm"
)

// NewUsers creates a new Users controller.
// It panics if the necessary templates are not parsed.
func NewUsers(cfg config.Config, db *gorm.DB) *Users {
	return &Users{
		NewView:   views.NewView(cfg.PageTemplateDir, views.Config{Title: "Join", Layout: "base"}, "users/new"),
		onPremise: cfg.OnPremise,
		db:        db,
	}
}

// Users is a user controller.
type Users struct {
	NewView   *views.View
	db        *gorm.DB
	onPremise bool
}

// New handles GET /register
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	var form RegistrationForm
	parseURLParams(r, &form)
	u.NewView.Render(w, r, form)
}

// RegistrationForm is the form data for registering
type RegistrationForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// LoginForm is the form data for log in
type LoginForm struct {
	Email    string `schema:"email" json:"email"`
	Password string `schema:"password" json:"password"`
}
