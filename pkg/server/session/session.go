package session

import (
	"github.com/dnote/dnote/pkg/server/database"
)

// Session represents user session
type Session struct {
	UUID          string `json:"uuid"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Pro           bool   `json:"pro"`
}

// New returns a new session for the given user
func New(user database.User, account database.Account) Session {
	return Session{
		UUID:          user.UUID,
		Pro:           user.Cloud,
		Email:         account.Email.String,
		EmailVerified: account.EmailVerified,
	}
}
