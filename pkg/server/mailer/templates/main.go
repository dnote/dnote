/* Copyright (C) 2019 Monomax Software Pty Ltd
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
	"os"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/job/repetition"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func weeklyDigestHandler(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	q := r.URL.Query()
	digestUUID := q.Get("digest_uuid")
	if digestUUID == "" {
		http.Error(w, errors.New("Please provide digest_uuid query param").Error(), http.StatusBadRequest)
		return
	}

	var user database.User
	if err := db.First(&user).Error; err != nil {
		http.Error(w, errors.Wrap(err, "Failed to find user").Error(), http.StatusInternalServerError)
		return
	}

	var digest database.Digest
	if err := db.Where("uuid = ?", digestUUID).First(&digest).Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding digest").Error(), http.StatusInternalServerError)
		return
	}

	var rule database.RepetitionRule
	if err := db.Where("id = ?", digest.RuleID).First(&rule).Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding digest").Error(), http.StatusInternalServerError)
		return
	}

	now := time.Now()
	email, err := repetition.BuildEmail(now, user, "sung@getdnote.com", digest, rule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := email.Body
	w.Write([]byte(body))
}

func emailVerificationHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Subject string
		Token   string
	}{
		"Verify your email",
		"testToken",
	}
	email := mailer.NewEmail("noreply@getdnote.com", []string{"sung@getdnote.com"}, "Reset your password")
	err := email.ParseTemplate(mailer.EmailTypeEmailVerification, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := email.Body
	w.Write([]byte(body))
}

func init() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		panic(err)
	}
}

func main() {
	c := database.Config{
		Host:     os.Getenv("DBHost"),
		Port:     os.Getenv("DBPort"),
		Name:     os.Getenv("DBName"),
		User:     os.Getenv("DBUser"),
		Password: os.Getenv("DBPassword"),
	}
	database.Open(c)
	defer database.Close()

	mailer.InitTemplates(nil)

	log.Println("Email template development server running on http://127.0.0.1:2300")

	http.HandleFunc("/weekly-digest", weeklyDigestHandler)
	http.HandleFunc("/email-verification", emailVerificationHandler)
	log.Fatal(http.ListenAndServe(":2300", nil))
}
