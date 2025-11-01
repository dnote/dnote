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

package views

import (
	"net/http"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
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
	// CSRF  template.HTML
	User  *database.User
	Yield map[string]interface{}
}

func getErrMessage(err error) string {
	if pErr, ok := err.(PublicError); ok {
		return pErr.Public()
	}

	return AlertMsgGeneric
}

// PutAlert puts an alert in the given data.
func (d *Data) PutAlert(alert Alert, alertInYield bool) {
	if alertInYield {
		if d.Yield == nil {
			d.Yield = map[string]interface{}{}
		}
		d.Yield["Alert"] = &alert
	} else {
		d.Alert = &alert
	}
}

// SetAlert sets alert in the given data for given error.
func (d *Data) SetAlert(err error, alertInYield bool) {
	errC := errors.Cause(err)

	var alert Alert
	if pErr, ok := errC.(PublicError); ok {
		alert = Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
		alert = Alert{
			Level:   AlertLvlError,
			Message: AlertMsgGeneric,
		}
	}

	d.PutAlert(alert, alertInYield)
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
		Path:     "/",
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		Path:     "/",
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
