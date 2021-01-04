package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/views"
)

// NewUsers creates a new Users controller.
// It panics if the necessary templates are not parsed.
func NewUsers(cfg config.Config, app *app.App) *Users {
	return &Users{
		NewView:   views.NewView(cfg.PageTemplateDir, views.Config{Title: "Join", Layout: "base"}, "users/new"),
		LoginView: views.NewView(cfg.PageTemplateDir, views.Config{Title: "Login", Layout: "base"}, "users/login"),
		onPremise: cfg.OnPremise,
		app:       app,
	}
}

// Users is a user controller.
type Users struct {
	NewView   *views.View
	LoginView *views.View
	app       *app.App
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

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := parseRequestData(r, &form); err != nil {
		log.Error(err.Error())
		w.WriteHeader(500)
		return
	}

	user, err := u.app.Authenticate(form.Email, form.Password)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(500)
		return
	}

	session, err := u.app.SignIn(user)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(500)
		return
	}

	setSessionCookie(w, session.Key, session.ExpiresAt)

	response := SessionResponse{
		Key:       session.Key,
		ExpiresAt: session.ExpiresAt.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err.Error())
		return
	}

}
