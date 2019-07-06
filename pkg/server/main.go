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
	"strings"

	"github.com/dnote/dnote/pkg/server/api/clock"
	"github.com/dnote/dnote/pkg/server/api/handlers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/job"
	"github.com/dnote/dnote/pkg/server/mailer"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var versionTag = "master"
var port = flag.String("port", "8080", "port to connect to")

func init() {
}

func getAppHandler() http.HandlerFunc {
	box := packr.New("web", "../../web/public")

	fs := http.FileServer(box)
	appShell, err := box.Find("index.html")
	if err != nil {
		panic(errors.Wrap(err, "getting index.html content"))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 2 && parts[1] == "dist" {
			fs.ServeHTTP(w, r)
			return
		}

		// All other requests should serve the index.html file
		w.Write(appShell)
	}
}

func initServer() *mux.Router {
	srv := mux.NewRouter()

	apiRouter := handlers.NewRouter(&handlers.App{
		Clock:            clock.New(),
		StripeAPIBackend: nil,
	})

	srv.PathPrefix("/api").Handler(http.StripPrefix("/api", apiRouter))
	srv.PathPrefix("/").HandlerFunc(getAppHandler())

	return srv
}

func startCmd() {
	mailer.InitTemplates()
	database.InitDB()
	database.InitSchema()
	defer database.CloseDB()

	// Perform database migration
	if err := database.Migrate(); err != nil {
		panic(errors.Wrap(err, "running migrations"))
	}

	// Run job in the background
	go job.Run()

	srv := initServer()

	log.Printf("Dnote version %s is running on port %s", versionTag, *port)
	addr := fmt.Sprintf(":%s", *port)
	log.Println(http.ListenAndServe(addr, srv))
}

func versionCmd() {
	fmt.Printf("dnote-server-%s\n", versionTag)
}

func main() {
	flag.Parse()
	cmd := flag.Arg(0)

	switch cmd {
	case "":
		fmt.Printf(`Dnote Server - A simple notebook for developers

Usage:
  dnote-server [command]

Available commands:
  start: Start the server
  version: Print the version
`)
	case "start":
		startCmd()
	case "version":
		versionCmd()
	default:
		fmt.Printf("Unknown command %s", cmd)
	}
}
