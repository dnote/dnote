package controllers

import (
	"net/http"
	"net/url"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/views"
	"github.com/pkg/errors"
)

var commonHelpers = map[string]interface{}{
	"getPathWithReferrer": func(base, referrer string) string {
		if referrer == "" {
			return base
		}

		query := url.Values{}
		query.Set("referrer", referrer)

		return helpers.GetPath(base, &query)
	},
}

// NewUsers creates a new Users controller.
// It panics if the necessary templates are not parsed.
func NewUsers(app *app.App) *Users {
	return &Users{
		NewView: views.NewView(
			app.Config.PageTemplateDir,
			views.Config{Title: "Join", Layout: "base", HelperFuncs: commonHelpers},
			"users/new",
		),
		LoginView: views.NewView(
			app.Config.PageTemplateDir,
			views.Config{Title: "Sign In", Layout: "base", HelperFuncs: commonHelpers},
			"users/login",
		),
		app: app,
	}
}

// Users is a user controller.
type Users struct {
	NewView   *views.View
	LoginView *views.View
	app       *app.App
}

// NewLogin renders user login page
func (u *Users) NewLogin(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	vd.Yield = struct {
		Referrer string
	}{
		Referrer: r.URL.Query().Get("referrer"),
	}

	u.LoginView.Render(w, r, vd)
}

// New renders user registration page
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	vd.Yield = struct {
		Email    string
		Referrer string
	}{
		Email:    "",
		Referrer: r.URL.Query().Get("referrer"),
	}

	u.NewView.Render(w, r, vd)
}

// RegistrationForm is the form data for registering
type RegistrationForm struct {
	Email                string `schema:"email"`
	Password             string `schema:"password"`
	PasswordConfirmation string `schema:"password_confirmation"`
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

	user, err := u.app.CreateUser(form.Email, form.Password, form.PasswordConfirmation)
	if err != nil {
		handleHTMLError(w, r, err, "creating user", u.NewView, vd)
		return
	}

	session, err := u.app.SignIn(&user)
	if err != nil {
		handleHTMLError(w, r, err, "signing in a user", u.LoginView, vd)
		return
	}

	if err := u.app.SendWelcomeEmail(form.Email); err != nil {
		log.ErrorWrap(err, "sending welcome email")
	}

	setSessionCookie(w, session.Key, session.ExpiresAt)

	dest := getPathOrReferrer("/", r)
	http.Redirect(w, r, dest, http.StatusFound)
}

// LoginForm is the form data for log in
type LoginForm struct {
	Email    string `schema:"email" json:"email"`
	Password string `schema:"password" json:"password"`
}

func (u *Users) login(r *http.Request) (*database.Session, error) {
	var form LoginForm
	if err := parseRequestData(r, &form); err != nil {
		return nil, err
	}

	user, err := u.app.Authenticate(form.Email, form.Password)
	if err != nil {
		// If the user is not found, treat it as invalid login
		if err == app.ErrNotFound {
			return nil, app.ErrLoginInvalid
		}

		return nil, err
	}

	s, err := u.app.SignIn(user)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func getPathOrReferrer(path string, r *http.Request) string {
	q := r.URL.Query()
	referrer := q.Get("referrer")

	if referrer == "" {
		return path
	}

	return referrer
}

// Login handles login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	session, err := u.login(r)
	if err != nil {
		handleHTMLError(w, r, err, "logging in user", u.LoginView, vd)
		return
	}

	dest := getPathOrReferrer("/", r)

	setSessionCookie(w, session.Key, session.ExpiresAt)
	http.Redirect(w, r, dest, http.StatusFound)
}

// V3Login handles login
func (u *Users) V3Login(w http.ResponseWriter, r *http.Request) {
	session, err := u.login(r)
	if err != nil {
		handleJSONError(w, err, "logging in user")
		return
	}

	respondWithSession(w, http.StatusOK, session)
}

func (u *Users) logout(r *http.Request) (bool, error) {
	key, err := GetCredential(r)
	if err != nil {
		return false, errors.Wrap(err, "getting credentials")
	}

	if key == "" {
		return false, nil
	}

	if err = u.app.DeleteSession(key); err != nil {
		return false, errors.Wrap(err, "deleting session")
	}

	return true, nil
}

// Logout handles logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	ok, err := u.logout(r)
	if err != nil {
		handleHTMLError(w, r, err, "logging out", u.LoginView, vd)
		return
	}

	if ok {
		unsetSessionCookie(w)
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

// V3Logout handles logout via API
func (u *Users) V3Logout(w http.ResponseWriter, r *http.Request) {
	ok, err := u.logout(r)
	if err != nil {
		handleJSONError(w, err, "logging out")
		return
	}

	if ok {
		unsetSessionCookie(w)
	}

	w.WriteHeader(http.StatusNoContent)
}
