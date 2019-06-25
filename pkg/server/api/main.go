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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dnote/dnote/pkg/server/api/clock"
	"github.com/dnote/dnote/pkg/server/api/handlers"
	"github.com/dnote/dnote/pkg/server/api/logger"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/mailer"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/pkg/errors"
)

var (
	emailTemplateDir = flag.String("emailTemplateDir", "../mailer/templates/src", "the path to the template directory")
)

func getOauthCallbackURL(provider string) string {
	if os.Getenv("GO_ENV") == "PRODUCTION" {
		return fmt.Sprintf("%s/api/auth/%s/callback", os.Getenv("WebHost"), provider)
	}

	return fmt.Sprintf("%s:%s/api/auth/%s/callback", os.Getenv("Host"), os.Getenv("WebPort"), provider)
}

func init() {
	// Set up Oauth
	gothic.Store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
	goth.UseProviders(
		github.New(
			os.Getenv("GithubClientID"),
			os.Getenv("GithubClientSecret"),
			getOauthCallbackURL("github"),
		),
		gplus.New(
			os.Getenv("GoogleClientID"),
			os.Getenv("GoogleClientSecret"),
			getOauthCallbackURL("gplus"),
			"https://www.googleapis.com/auth/plus.me",
		),
	)

	gothic.GetProviderName = func(r *http.Request) (name string, err error) {
		vars := mux.Vars(r)
		name = vars["provider"]
		return
	}
}

func main() {
	flag.Parse()

	mailer.InitTemplates(*emailTemplateDir)

	database.InitDB()
	database.InitSchema()
	defer database.CloseDB()

	if err := logger.Init(); err != nil {
		log.Println(errors.Wrap(err, "initializing logger"))
	}

	app := handlers.App{
		Clock:            clock.New(),
		StripeAPIBackend: nil,
	}
	r := handlers.NewRouter(&app)

	port := os.Getenv("PORT")
	logger.Notice("API listening on port %s", port)

	log.Println(http.ListenAndServe(":"+port, r))
}
