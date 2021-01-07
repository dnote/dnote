package views

import (
	"html/template"
	"net/http"
	"time"

	"github.com/dnote/dnote/pkg/server/models"
	"github.com/pkg/errors"
)

const (
	// AlertLvlError is an alert level for error
	AlertLvlError = "danger"
	// AlertLvlWarning is an alert level for warning
	AlertLvlWarning = "warning"
	// AlertLvlInfo is an alert level for info
	AlertLvlInfo = "info"
	// AlertLvlSuccess is an alert level for success
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is a generic message for a server error
	AlertMsgGeneric = "Something went wrong. Please try again."
)

// Alert is used to render Bootstrap Alert messages in templates
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data to come in.
type Data struct {
	Alert *Alert
	CSRF  template.HTML
	User  *models.User
	Yield interface{}
}

func getErrMessage(err error) string {
	if pErr, ok := err.(PublicError); ok {
		return pErr.Public()
	}

	return AlertMsgGeneric
}

// SetAlert sets alert in the given data for given error.
func (d *Data) SetAlert(err error) {
	errC := errors.Cause(err)

	if pErr, ok := errC.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: AlertMsgGeneric,
		}
	}
}

// AlertError returns a new error alert using the given message.
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func persistAlert(w http.ResponseWriter, alert Alert) {
	expiresAt := time.Now().Add(5 * time.Minute)
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func clearAlert(w http.ResponseWriter) {
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func getAlert(r *http.Request) *Alert {
	lvl, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}
	msg, err := r.Cookie("alert_message")
	if err != nil {
		return nil
	}
	alert := Alert{
		Level:   lvl.Value,
		Message: msg.Value,
	}
	return &alert
}

// RedirectAlert redirects to a URL after persisting the provided alert data
// into a cookie so that it can be displayed when the page is rendered.
func RedirectAlert(w http.ResponseWriter, r *http.Request, urlStr string, code int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, r, urlStr, code)
}

// PublicError is an error meant to be displayed to the public
type PublicError interface {
	error
	Public() string
}
