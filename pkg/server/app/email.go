package app

import (
	"fmt"
	"net/url"

	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/pkg/errors"
)

var defaultSender = "sung@getdnote.com"

func (a *App) getSenderEmail(want string) (string, error) {
	if !a.SelfHosted {
		return want, nil
	}

	addr, err := a.getNoreplySender()
	if err != nil {
		return "", errors.Wrap(err, "getting sender email address")
	}

	return addr, nil
}

func (a *App) getNoreplySender() (string, error) {
	u, err := url.Parse(a.WebURL)
	if err != nil {
		return "", errors.Wrap(err, "parsing web url")
	}

	hostname := u.Hostname()
	addr := fmt.Sprintf("noreply@%s", hostname)
	return addr, nil
}

// SendVerificationEmail sends verification email
func (a *App) SendVerificationEmail(email, tokenValue string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeEmailVerification, mailer.EmailKindText, mailer.EmailVerificationTmplData{
		Token:  tokenValue,
		WebURL: a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}

	from, err := a.getSenderEmail(defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Verify your Dnote email address", from, []string{email}, "text/plain", body); err != nil {
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

	from, err := a.getSenderEmail(defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Welcome to Dnote!", from, []string{email}, "text/plain", body); err != nil {
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

	from, err := a.getSenderEmail(defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Reset your password", from, []string{email}, "text/plain", body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}
