package routes

import (
	"net/http"
	"time"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/controllers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/pkg/errors"
)

// AuthParams is the params for the authentication middleware
type AuthParams struct {
	ProOnly               bool
	RedirectGuestsToLogin bool
}

// AuthWithSession performs user authentication with session
func AuthWithSession(app *app.App, r *http.Request, p *AuthParams) (database.User, bool, error) {
	var user database.User

	sessionKey, err := controllers.GetCredential(r)
	if err != nil {
		return user, false, errors.Wrap(err, "getting credential")
	}
	if sessionKey == "" {
		return user, false, nil
	}

	var session database.Session
	conn := app.DB.Where("key = ?", sessionKey).First(&session)

	if conn.RecordNotFound() {
		return user, false, nil
	} else if err := conn.Error; err != nil {
		return user, false, errors.Wrap(err, "finding session")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return user, false, nil
	}

	conn = app.DB.Where("id = ?", session.UserID).First(&user)

	if conn.RecordNotFound() {
		return user, false, nil
	} else if err := conn.Error; err != nil {
		return user, false, errors.Wrap(err, "finding user from token")
	}

	return user, true, nil
}
