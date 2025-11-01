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
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/buildinfo"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/token"
	"github.com/dnote/dnote/pkg/server/views"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	pkgErrors "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var commonHelpers = map[string]interface{}{
	"getPathWithReferrer": func(base string, referrer string) string {
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
func NewUsers(app *app.App, viewEngine *views.Engine) *Users {
	return &Users{
		NewView: viewEngine.NewView(app,
			views.Config{Title: "Join", Layout: "base", HelperFuncs: commonHelpers, AlertInBody: true},
			"users/new",
		),
		LoginView: viewEngine.NewView(app,
			views.Config{Title: "Sign In", Layout: "base", HelperFuncs: commonHelpers, AlertInBody: true},
			"users/login",
		),
		PasswordResetView: viewEngine.NewView(app,
			views.Config{Title: "Reset Password", Layout: "base", HelperFuncs: commonHelpers, AlertInBody: true},
			"users/password_reset",
		),
		PasswordResetConfirmView: viewEngine.NewView(app,
			views.Config{Title: "Reset Password", Layout: "base", HelperFuncs: commonHelpers, AlertInBody: true},
			"users/password_reset_confirm",
		),
		SettingView: viewEngine.NewView(app,
			views.Config{Layout: "base", HelperFuncs: commonHelpers, HeaderTemplate: "navbar"},
			"users/settings",
		),
		AboutView: viewEngine.NewView(app,
			views.Config{Title: "About", Layout: "base", HelperFuncs: commonHelpers, HeaderTemplate: "navbar"},
			"users/settings_about",
		),
		app: app,
	}
}

// Users is a user controller.
type Users struct {
	NewView                  *views.View
	LoginView                *views.View
	SettingView              *views.View
	AboutView                *views.View
	PasswordResetView        *views.View
	PasswordResetConfirmView *views.View
	app                      *app.App
}

// New renders user registration page
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	vd := getDataWithReferrer(r)
	u.NewView.Render(w, r, &vd, http.StatusOK)
}

// RegistrationForm is the form data for registering
type RegistrationForm struct {
	Email                string `schema:"email"`
	Password             string `schema:"password"`
	PasswordConfirmation string `schema:"password_confirmation"`
}

// Create handles register
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	vd := getDataWithReferrer(r)

	var form RegistrationForm
	if err := parseForm(r, &form); err != nil {
		handleHTMLError(w, r, err, "parsing form", u.NewView, vd)
		return
	}

	vd.Yield["Email"] = form.Email

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

func (u *Users) login(form LoginForm) (*database.Session, error) {
	if form.Email == "" {
		return nil, app.ErrEmailRequired
	}
	if form.Password == "" {
		return nil, app.ErrPasswordRequired
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

func getDataWithReferrer(r *http.Request) views.Data {
	vd := views.Data{}

	vd.Yield = map[string]interface{}{
		"Referrer": r.URL.Query().Get("referrer"),
	}

	return vd
}

// NewLogin renders user login page
func (u *Users) NewLogin(w http.ResponseWriter, r *http.Request) {
	vd := getDataWithReferrer(r)
	u.LoginView.Render(w, r, &vd, http.StatusOK)
}

// Login handles login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := getDataWithReferrer(r)

	var form LoginForm
	if err := parseRequestData(r, &form); err != nil {
		handleHTMLError(w, r, err, "parsing payload", u.LoginView, vd)
		return
	}

	session, err := u.login(form)
	if err != nil {
		vd.Yield["Email"] = form.Email
		handleHTMLError(w, r, err, "logging in user", u.LoginView, vd)
		return
	}

	setSessionCookie(w, session.Key, session.ExpiresAt)

	dest := getPathOrReferrer("/", r)
	http.Redirect(w, r, dest, http.StatusFound)
}

// V3Login handles login
func (u *Users) V3Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := parseRequestData(r, &form); err != nil {
		handleJSONError(w, err, "parsing payload")
		return
	}

	session, err := u.login(form)
	if err != nil {
		handleJSONError(w, err, "logging in user")
		return
	}

	respondWithSession(w, http.StatusOK, session)
}

func (u *Users) logout(r *http.Request) (bool, error) {
	key, err := GetCredential(r)
	if err != nil {
		return false, pkgErrors.Wrap(err, "getting credentials")
	}

	if key == "" {
		return false, nil
	}

	if err = u.app.DeleteSession(key); err != nil {
		return false, pkgErrors.Wrap(err, "deleting session")
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

type createResetTokenPayload struct {
	Email string `schema:"email" json:"email"`
}

func (u *Users) CreateResetToken(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	var form createResetTokenPayload
	if err := parseForm(r, &form); err != nil {
		handleHTMLError(w, r, err, "parsing form", u.PasswordResetView, vd)
		return
	}

	if form.Email == "" {
		handleHTMLError(w, r, app.ErrEmailRequired, "email is not provided", u.PasswordResetView, vd)
		return
	}

	var user database.User
	err := u.app.DB.Where("email = ?", form.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if err != nil {
		handleHTMLError(w, r, err, "finding user", u.PasswordResetView, vd)
		return
	}

	resetToken, err := token.Create(u.app.DB, user.ID, database.TokenTypeResetPassword)
	if err != nil {
		handleHTMLError(w, r, err, "generating token", u.PasswordResetView, vd)
		return
	}

	if err := u.app.SendPasswordResetEmail(user.Email.String, resetToken.Value); err != nil {
		handleHTMLError(w, r, err, "sending password reset email", u.PasswordResetView, vd)
		return
	}

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Check your email for a link to reset your password.",
	}
	views.RedirectAlert(w, r, "/password-reset", http.StatusFound, alert)
}

// PasswordResetConfirm renders password reset view
func (u *Users) PasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	vars := mux.Vars(r)
	token := vars["token"]

	vd.Yield = map[string]interface{}{
		"Token": token,
	}

	u.PasswordResetConfirmView.Render(w, r, &vd, http.StatusOK)
}

type resetPasswordPayload struct {
	Password             string `schema:"password" json:"password"`
	PasswordConfirmation string `schema:"password_confirmation" json:"password_confirmation"`
	Token                string `schema:"token" json:"token"`
}

// PasswordReset renders password reset view
func (u *Users) PasswordReset(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	var params resetPasswordPayload
	if err := parseForm(r, &params); err != nil {
		handleHTMLError(w, r, err, "parsing params", u.NewView, vd)
		return
	}

	vd.Yield = map[string]interface{}{
		"Token": params.Token,
	}

	if params.Password != params.PasswordConfirmation {
		handleHTMLError(w, r, app.ErrPasswordConfirmationMismatch, "password mismatch", u.PasswordResetConfirmView, vd)
		return
	}

	var token database.Token
	err := u.app.DB.Where("value = ? AND type =? AND used_at IS NULL", params.Token, database.TokenTypeResetPassword).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handleHTMLError(w, r, app.ErrInvalidToken, "invalid token", u.PasswordResetConfirmView, vd)
		return
	}
	if err != nil {
		handleHTMLError(w, r, err, "finding token", u.PasswordResetConfirmView, vd)
		return
	}

	if token.UsedAt != nil {
		handleHTMLError(w, r, app.ErrInvalidToken, "invalid token", u.PasswordResetConfirmView, vd)
		return
	}

	// Expire after 10 minutes
	if time.Since(token.CreatedAt).Minutes() > 10 {
		handleHTMLError(w, r, app.ErrPasswordResetTokenExpired, "expired token", u.PasswordResetConfirmView, vd)
		return
	}

	var user database.User
	if err := u.app.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		handleHTMLError(w, r, err, "finding user", u.PasswordResetConfirmView, vd)
		return
	}

	tx := u.app.DB.Begin()

	// Update the password
	if err := app.UpdateUserPassword(tx, &user, params.Password); err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "updating password", u.PasswordResetConfirmView, vd)
		return
	}

	if err := tx.Model(&token).Update("used_at", time.Now()).Error; err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "updating password reset token", u.PasswordResetConfirmView, vd)
		return
	}

	if err := u.app.DeleteUserSessions(tx, user.ID); err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "deleting user sessions", u.PasswordResetConfirmView, vd)
		return
	}

	tx.Commit()

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Password reset successful",
	}
	views.RedirectAlert(w, r, "/login", http.StatusFound, alert)

	if err := u.app.SendPasswordResetAlertEmail(user.Email.String); err != nil {
		log.ErrorWrap(err, "sending password reset email")
	}
}

func (u *Users) logoutOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Version")
}

func (u *Users) Settings(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	u.SettingView.Render(w, r, &vd, http.StatusOK)
}

func (u *Users) About(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	vd.Yield = map[string]interface{}{
		"Version": buildinfo.Version,
	}

	u.AboutView.Render(w, r, &vd, http.StatusOK)
}

type updatePasswordForm struct {
	OldPassword             string `schema:"old_password"`
	NewPassword             string `schema:"new_password"`
	NewPasswordConfirmation string `schema:"new_password_confirmation"`
}

func (u *Users) PasswordUpdate(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	user := context.User(r.Context())
	if user == nil {
		handleHTMLError(w, r, app.ErrLoginRequired, "No authenticated user found", u.SettingView, vd)
		return
	}

	var form updatePasswordForm
	if err := parseRequestData(r, &form); err != nil {
		handleHTMLError(w, r, err, "parsing payload", u.LoginView, vd)
		return
	}

	if form.OldPassword == "" || form.NewPassword == "" {
		handleHTMLError(w, r, app.ErrInvalidPasswordChangeInput, "invalid params", u.SettingView, vd)
		return
	}
	if form.NewPassword != form.NewPasswordConfirmation {
		handleHTMLError(w, r, app.ErrPasswordConfirmationMismatch, "passwords do not match", u.SettingView, vd)
		return
	}

	password := []byte(form.OldPassword)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), password); err != nil {
		log.WithFields(log.Fields{
			"user_id": user.ID,
		}).Warn("invalid password update attempt")
		handleHTMLError(w, r, app.ErrInvalidPassword, "invalid password", u.SettingView, vd)
		return
	}

	if err := app.UpdateUserPassword(u.app.DB, user, form.NewPassword); err != nil {
		handleHTMLError(w, r, err, "updating password", u.SettingView, vd)
		return
	}

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Password change successful",
	}
	views.RedirectAlert(w, r, "/", http.StatusFound, alert)
}

type updateProfileForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Users) ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	user := context.User(r.Context())
	if user == nil {
		handleHTMLError(w, r, app.ErrLoginRequired, "No authenticated user found", u.SettingView, vd)
		return
	}

	var form updateProfileForm
	if err := parseRequestData(r, &form); err != nil {
		handleHTMLError(w, r, err, "parsing payload", u.SettingView, vd)
		return
	}

	password := []byte(form.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), password); err != nil {
		log.WithFields(log.Fields{
			"user_id": user.ID,
		}).Warn("invalid email update attempt")
		handleHTMLError(w, r, app.ErrInvalidPassword, "Wrong password", u.SettingView, vd)
		return
	}

	// Validate
	if len(form.Email) > 60 {
		handleHTMLError(w, r, app.ErrEmailTooLong, "Email is too long", u.SettingView, vd)
		return
	}

	user.Email.String = form.Email

	if err := u.app.DB.Save(&user).Error; err != nil {
		handleHTMLError(w, r, err, "saving user", u.SettingView, vd)
		return
	}

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Email change successful",
	}
	views.RedirectAlert(w, r, "/", http.StatusFound, alert)
}

