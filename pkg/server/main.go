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

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/api/handlers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/job"
	"github.com/dnote/dnote/pkg/server/mailer"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var versionTag = "master"
var port = flag.String("port", "3000", "port to connect to")

var rootBox *packr.Box

func init() {
	rootBox = packr.New("root", "../../web/public")
}

func mustFind(box *packr.Box, path string) []byte {
	b, err := rootBox.Find(path)
	if err != nil {
		panic(errors.Wrapf(err, "getting file content for %s", path))
	}

	return b
}

func getStaticHandler() http.Handler {
	box := packr.New("static", "../../web/public/static")

	return http.StripPrefix("/static/", http.FileServer(box))
}

// getRootHandler returns an HTTP handler that serves the app shell
func getRootHandler() http.HandlerFunc {
	b := mustFind(rootBox, "index.html")

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(b)
	}
}

func getRobotsHandler() http.HandlerFunc {
	b := mustFind(rootBox, "robots.txt")

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(b)
	}
}

func getSWHandler() http.HandlerFunc {
	b := mustFind(rootBox, "service-worker.js")

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(b)
	}
}

func initServer() *mux.Router {
	srv := mux.NewRouter()

	apiRouter := handlers.NewRouter(&handlers.App{
		Clock:            clock.New(),
		StripeAPIBackend: nil,
	})

	srv.PathPrefix("/api").Handler(http.StripPrefix("/api", apiRouter))
	srv.PathPrefix("/static").Handler(getStaticHandler())
	srv.Handle("/service-worker.js", getSWHandler())
	srv.Handle("/robots.txt", getRobotsHandler())

	// For all other requests, serve the index.html file
	srv.PathPrefix("/").Handler(getRootHandler())

	return srv
}

func startCmd() {
	c := database.Config{
		Host:     os.Getenv("DBHost"),
		Port:     os.Getenv("DBPort"),
		Name:     os.Getenv("DBName"),
		User:     os.Getenv("DBUser"),
		Password: os.Getenv("DBPassword"),
	}
	database.Open(c)
	database.InitSchema()
	defer database.Close()

	mailer.InitTemplates(nil)

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
