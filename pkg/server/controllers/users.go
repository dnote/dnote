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

	vd := views.Data{}

	user, err := u.app.Authenticate(form.Email, form.Password)
	if err != nil {
		handleHTMLError(w, err, "authenticating user", &vd)
		u.LoginView.Render(w, r, vd)
		return
	}

	session, err := u.app.SignIn(user)
	if err != nil {
		handleHTMLError(w, err, "signing in a user", &vd)
		u.LoginView.Render(w, r, vd)
		return
	}

	setSessionCookie(w, session.Key, session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	key, err := GetCredential(r)
	if err != nil {
		handleHTMLError(w, err, "signing in a user", &vd)
		u.LoginView.Render(w, r, vd)
		return
	}

	if err = u.app.DeleteSession(key); err != nil {
		handleHTMLError(w, err, "signing in a user", &vd)
		u.LoginView.Render(w, r, vd)
		return
	}

	unsetSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}
