package controllers

import (
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
		app:       app,
	}
}

// Users is a user controller.
type Users struct {
	NewView   *views.View
	LoginView *views.View
	app       *app.App
}

// New renders user registration page
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

// Create handles register
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	var form RegistrationForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	user, err := u.app.CreateUser(form.Email, form.Password)
	if err != nil {
		handleHTMLError(w, r, err, "creating user", u.NewView, &vd)
		return
	}

	session, err := u.app.SignIn(&user)
	if err != nil {
		handleHTMLError(w, r, err, "signing in a user", u.LoginView, &vd)
		return
	}

	setSessionCookie(w, session.Key, session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusFound)

	if err := u.app.SendWelcomeEmail(form.Email); err != nil {
		log.ErrorWrap(err, "sending welcome email")
	}
}

// LoginForm is the form data for log in
type LoginForm struct {
	Email    string `schema:"email" json:"email"`
	Password string `schema:"password" json:"password"`
}

// Login handles login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	var form LoginForm
	if err := parseRequestData(r, &form); err != nil {
		handleHTMLError(w, r, err, "parsing request data", u.LoginView, &vd)
		return
	}

	user, err := u.app.Authenticate(form.Email, form.Password)
	if err != nil {
		handleHTMLError(w, r, err, "authenticating user", u.LoginView, &vd)
		return
	}

	session, err := u.app.SignIn(user)
	if err != nil {
		handleHTMLError(w, r, err, "signing in a user", u.LoginView, &vd)
		return
	}

	setSessionCookie(w, session.Key, session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Logout handles logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	key, err := GetCredential(r)
	if err != nil {
		handleHTMLError(w, r, err, "getting credentials", u.LoginView, &vd)
		u.LoginView.Render(w, r, vd)
		return
	}

	if err = u.app.DeleteSession(key); err != nil {
		handleHTMLError(w, r, err, "deleting session", u.LoginView, &vd)
		return
	}

	unsetSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}
