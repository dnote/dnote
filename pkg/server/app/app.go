package app

import (
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
)

var (
	// ErrEmptyDB is an error for missing database connection in the app configuration
	ErrEmptyDB = errors.New("No database connection was provided")
	// ErrEmptyClock is an error for missing clock in the app configuration
	ErrEmptyClock = errors.New("No clock was provided")
	// ErrEmptyWebURL is an error for missing WebURL content in the app configuration
	ErrEmptyWebURL = errors.New("No WebURL was provided")
	// ErrEmptyEmailTemplates is an error for missing EmailTemplates content in the app configuration
	ErrEmptyEmailTemplates = errors.New("No EmailTemplate store was provided")
	// ErrEmptyEmailBackend is an error for missing EmailBackend content in the app configuration
	ErrEmptyEmailBackend = errors.New("No EmailBackend was provided")
)

// App is an application configuration
type App struct {
	DB               *gorm.DB
	Clock            clock.Clock
	StripeAPIBackend stripe.Backend
	EmailTemplates   mailer.Templates
	EmailBackend     mailer.Backend
	WebURL           string
	SelfHosted       bool
}

// Validate validates the app configuration
func (a *App) Validate() error {
	if a.WebURL == "" {
		return ErrEmptyWebURL
	}
	if a.Clock == nil {
		return ErrEmptyClock
	}
	if a.EmailTemplates == nil {
		return ErrEmptyEmailTemplates
	}
	if a.EmailBackend == nil {
		return ErrEmptyEmailBackend
	}
	if a.DB == nil {
		return ErrEmptyDB
	}

	return nil
}
