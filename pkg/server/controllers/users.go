package controllers

import (
	"net/http"
	"net/url"
	"time"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/buildinfo"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/token"
	"github.com/dnote/dnote/pkg/server/views"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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
func NewUsers(app *app.App, baseDir string) *Users {
	return &Users{
		NewView: views.NewView(baseDir, app,
			views.Config{Title: "Join", Layout: "base", HelperFuncs: commonHelpers, AlertInBody: true},
			"users/new",
		),
		LoginView: views.NewView(baseDir, app,
			views.Config{Title: "Sign In", Layout: "base", HelperFuncs: commonHelpers, AlertInBody: true},
			"users/login",
		),
		PasswordResetView: views.NewView(baseDir, app,
			views.Config{Title: "Reset Password", Layout: "base", HelperFuncs: commonHelpers, AlertInBody: true},
			"users/password_reset",
		),
		PasswordResetConfirmView: views.NewView(baseDir, app,
			views.Config{Title: "Reset Password", Layout: "base", HelperFuncs: commonHelpers, AlertInBody: true},
			"users/password_reset_confirm",
		),
		SettingView: views.NewView(baseDir, app,
			views.Config{Layout: "base", HelperFuncs: commonHelpers, HeaderTemplate: "navbar"},
			"users/settings",
		),
		AboutView: views.NewView(baseDir, app,
			views.Config{Title: "About", Layout: "base", HelperFuncs: commonHelpers, HeaderTemplate: "navbar"},
			"users/settings_about",
		),
		EmailVerificationView: views.NewView(baseDir, app,
			views.Config{Layout: "base", HelperFuncs: commonHelpers, HeaderTemplate: "navbar"},
			"users/email_verification",
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
	EmailVerificationView    *views.View
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

	var account database.Account
	conn := u.app.DB.Where("email = ?", form.Email).First(&account)
	if conn.RecordNotFound() {
		return
	}
	if err := conn.Error; err != nil {
		handleHTMLError(w, r, err, "finding account", u.PasswordResetView, vd)
		return
	}

	resetToken, err := token.Create(u.app.DB, account.UserID, database.TokenTypeResetPassword)
	if err != nil {
		handleHTMLError(w, r, err, "generating token", u.PasswordResetView, vd)
		return
	}

	if err := u.app.SendPasswordResetEmail(account.Email.String, resetToken.Value); err != nil {
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
	conn := u.app.DB.Where("value = ? AND type =? AND used_at IS NULL", params.Token, database.TokenTypeResetPassword).First(&token)
	if conn.RecordNotFound() {
		handleHTMLError(w, r, app.ErrInvalidToken, "invalid token", u.PasswordResetConfirmView, vd)
		return
	}
	if err := conn.Error; err != nil {
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

	tx := u.app.DB.Begin()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "hashing password", u.PasswordResetConfirmView, vd)
		return
	}

	var account database.Account
	if err := u.app.DB.Where("user_id = ?", token.UserID).First(&account).Error; err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "finding user", u.PasswordResetConfirmView, vd)
		return
	}

	if err := tx.Model(&account).Update("password", string(hashedPassword)).Error; err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "updating password", u.PasswordResetConfirmView, vd)
		return
	}
	if err := tx.Model(&token).Update("used_at", time.Now()).Error; err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "updating password reset token", u.PasswordResetConfirmView, vd)
		return
	}

	if err := u.app.DeleteUserSessions(tx, account.UserID); err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "deleting user sessions", u.PasswordResetConfirmView, vd)
		return
	}

	tx.Commit()

	var user database.User
	if err := u.app.DB.Where("id = ?", account.UserID).First(&user).Error; err != nil {
		handleHTMLError(w, r, err, "finding user", u.PasswordResetConfirmView, vd)
		return
	}

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Password reset successful",
	}
	views.RedirectAlert(w, r, "/login", http.StatusFound, alert)

	if err := u.app.SendPasswordResetAlertEmail(account.Email.String); err != nil {
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

	var account database.Account
	if err := u.app.DB.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		handleHTMLError(w, r, err, "getting account", u.SettingView, vd)
		return
	}

	password := []byte(form.OldPassword)
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password.String), password); err != nil {
		log.WithFields(log.Fields{
			"user_id": user.ID,
		}).Warn("invalid password update attempt")
		handleHTMLError(w, r, app.ErrInvalidPassword, "invalid password", u.SettingView, vd)
		return
	}

	if err := validatePassword(form.NewPassword); err != nil {
		handleHTMLError(w, r, err, "invalid password", u.SettingView, vd)
		return
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(form.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		handleHTMLError(w, r, err, "hashing password", u.SettingView, vd)
		return
	}

	if err := u.app.DB.Model(&account).Update("password", string(hashedNewPassword)).Error; err != nil {
		handleHTMLError(w, r, err, "updating password", u.SettingView, vd)
		return
	}

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Password change successful",
	}
	views.RedirectAlert(w, r, "/", http.StatusFound, alert)
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return app.ErrPasswordTooShort
	}

	return nil
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

	var account database.Account
	if err := u.app.DB.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		handleHTMLError(w, r, err, "getting account", u.SettingView, vd)
		return
	}

	var form updateProfileForm
	if err := parseRequestData(r, &form); err != nil {
		handleHTMLError(w, r, err, "parsing payload", u.SettingView, vd)
		return
	}

	password := []byte(form.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password.String), password); err != nil {
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

	tx := u.app.DB.Begin()
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "saving user", u.SettingView, vd)
		return
	}

	// check if email was changed
	if form.Email != account.Email.String {
		account.EmailVerified = false
	}
	account.Email.String = form.Email

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "saving account", u.SettingView, vd)
		return
	}

	tx.Commit()

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Email change successful",
	}
	views.RedirectAlert(w, r, "/", http.StatusFound, alert)
}

func (u *Users) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	vars := mux.Vars(r)
	tokenValue := vars["token"]

	if tokenValue == "" {
		handleHTMLError(w, r, app.ErrMissingToken, "Missing email verification token", u.EmailVerificationView, vd)
		return
	}

	var token database.Token
	if err := u.app.DB.
		Where("value = ? AND type = ?", tokenValue, database.TokenTypeEmailVerification).
		First(&token).Error; err != nil {
		handleHTMLError(w, r, app.ErrInvalidToken, "Finding token", u.EmailVerificationView, vd)
		return
	}

	if token.UsedAt != nil {
		handleHTMLError(w, r, app.ErrInvalidToken, "Token has already been used.", u.EmailVerificationView, vd)
		return
	}

	// Expire after ttl
	if time.Since(token.CreatedAt).Minutes() > 30 {
		handleHTMLError(w, r, app.ErrExpiredToken, "Token has expired.", u.EmailVerificationView, vd)
		return
	}

	var account database.Account
	if err := u.app.DB.Where("user_id = ?", token.UserID).First(&account).Error; err != nil {
		handleHTMLError(w, r, err, "finding account", u.EmailVerificationView, vd)
		return
	}
	if account.EmailVerified {
		handleHTMLError(w, r, app.ErrEmailAlreadyVerified, "Already verified", u.EmailVerificationView, vd)
		return
	}

	tx := u.app.DB.Begin()
	account.EmailVerified = true
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "updating email_verified", u.EmailVerificationView, vd)
		return
	}
	if err := tx.Model(&token).Update("used_at", time.Now()).Error; err != nil {
		tx.Rollback()
		handleHTMLError(w, r, err, "updating reset token", u.EmailVerificationView, vd)
		return
	}
	tx.Commit()

	var user database.User
	if err := u.app.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		handleHTMLError(w, r, err, "finding user", u.EmailVerificationView, vd)
		return
	}

	session, err := u.app.SignIn(&user)
	if err != nil {
		handleHTMLError(w, r, err, "Creating session", u.EmailVerificationView, vd)
	}

	setSessionCookie(w, session.Key, session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u *Users) CreateEmailVerificationToken(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}

	user := context.User(r.Context())
	if user == nil {
		handleHTMLError(w, r, app.ErrLoginRequired, "No authenticated user found", u.SettingView, vd)
		return
	}

	var account database.Account
	err := u.app.DB.Where("user_id = ?", user.ID).First(&account).Error
	if err != nil {
		handleHTMLError(w, r, err, "finding account", u.SettingView, vd)
		return
	}

	if account.EmailVerified {
		handleHTMLError(w, r, app.ErrEmailAlreadyVerified, "email is already verified.", u.SettingView, vd)
		return
	}
	if account.Email.String == "" {
		handleHTMLError(w, r, app.ErrEmailRequired, "email is empty.", u.SettingView, vd)
		return
	}

	tok, err := token.Create(u.app.DB, account.UserID, database.TokenTypeEmailVerification)
	if err != nil {
		handleHTMLError(w, r, err, "saving token", u.SettingView, vd)
		return
	}

	if err := u.app.SendVerificationEmail(account.Email.String, tok.Value); err != nil {
		if errors.Cause(err) == mailer.ErrSMTPNotConfigured {
			handleHTMLError(w, r, app.ErrInvalidSMTPConfig, "SMTP config is not configured correctly.", u.SettingView, vd)
		} else {
			handleHTMLError(w, r, err, "sending verification email", u.SettingView, vd)
		}

		return
	}

	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Please check your email for the verification",
	}
	views.RedirectAlert(w, r, "/", http.StatusFound, alert)
}
