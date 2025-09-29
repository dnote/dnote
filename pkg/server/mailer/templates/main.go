/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"log"
	"net/http"

	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func (c Context) passwordResetHandler(w http.ResponseWriter, r *http.Request) {
	data := mailer.EmailResetPasswordTmplData{
		AccountEmail: "alice@example.com",
		Token:        "testToken",
		WebURL:       "http://localhost:3000",
	}
	body, err := c.Tmpl.Execute(mailer.EmailTypeResetPassword, mailer.EmailKindText, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(body))
}

func (c Context) passwordResetAlertHandler(w http.ResponseWriter, r *http.Request) {
	data := mailer.EmailResetPasswordAlertTmplData{
		AccountEmail: "alice@example.com",
		WebURL:       "http://localhost:3000",
	}
	body, err := c.Tmpl.Execute(mailer.EmailTypeResetPasswordAlert, mailer.EmailKindText, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(body))
}

func (c Context) emailVerificationHandler(w http.ResponseWriter, r *http.Request) {
	data := mailer.EmailVerificationTmplData{
		Token:  "testToken",
		WebURL: "http://localhost:3000",
	}
	body, err := c.Tmpl.Execute(mailer.EmailTypeEmailVerification, mailer.EmailKindText, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(body))
}

func (c Context) welcomeHandler(w http.ResponseWriter, r *http.Request) {
	data := mailer.WelcomeTmplData{
		AccountEmail: "alice@example.com",
		WebURL:       "http://localhost:3000",
	}
	body, err := c.Tmpl.Execute(mailer.EmailTypeWelcome, mailer.EmailKindText, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(body))
}

func (c Context) inactiveHandler(w http.ResponseWriter, r *http.Request) {
	data := mailer.InactiveReminderTmplData{
		SampleNoteUUID: "some-uuid",
		WebURL:         "http://localhost:3000",
		Token:          "some-random-token",
	}
	body, err := c.Tmpl.Execute(mailer.EmailTypeInactiveReminder, mailer.EmailKindText, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(body))
}

func (c Context) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Email development server is running."))
}

func init() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		panic(err)
	}
}

// Context is a context holding global information
type Context struct {
	DB   *gorm.DB
	Tmpl mailer.Templates
}

func main() {
	c := config.Load()
	db := database.Open(c)
	defer db.Close()

	log.Println("Email template development server running on http://127.0.0.1:2300")

	tmpl := mailer.NewTemplates()
	ctx := Context{DB: db, Tmpl: tmpl}

	http.HandleFunc("/", ctx.homeHandler)
	http.HandleFunc("/email-verification", ctx.emailVerificationHandler)
	http.HandleFunc("/password-reset", ctx.passwordResetHandler)
	http.HandleFunc("/password-reset-alert", ctx.passwordResetAlertHandler)
	http.HandleFunc("/welcome", ctx.welcomeHandler)
	http.HandleFunc("/inactive-reminder", ctx.inactiveHandler)
	log.Fatal(http.ListenAndServe(":2300", nil))
}
