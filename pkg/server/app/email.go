package app

import (
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/pkg/errors"
)

// SendVerificationEmail sends verification email
func (a *App) SendVerificationEmail(email, tokenValue string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeEmailVerification, mailer.EmailKindText, mailer.EmailVerificationTmplData{
		Token:  tokenValue,
		WebURL: a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}
	if err := a.EmailBackend.Queue("Verify your Dnote email address", "sung@getdnote.com", []string{email}, "text/plain", body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}

// SendWelcomeEmail sends welcome email
func (a *App) SendWelcomeEmail(email string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeWelcome, mailer.EmailKindText, mailer.WelcomeTmplData{
		AccountEmail: email,
		WebURL:       a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}
	if err := a.EmailBackend.Queue("Welcome to Dnote!", "sung@getdnote.com", []string{email}, "text/plain", body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}

// SendPasswordResetEmail sends verification email
func (a *App) SendPasswordResetEmail(email, tokenValue string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeResetPassword, mailer.EmailKindText, mailer.EmailResetPasswordTmplData{
		AccountEmail: email,
		Token:        tokenValue,
		WebURL:       a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}

	if err := a.EmailBackend.Queue("Reset your password", "sung@getdnote.com", []string{email}, "text/plain", body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}
