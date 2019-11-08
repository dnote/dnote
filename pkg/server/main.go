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
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/handlers"
	"github.com/dnote/dnote/pkg/server/job"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/web"

	"github.com/gobuffalo/packr/v2"
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

func initContext() web.Context {
	staticBox := packr.New("static", "../../web/public/static")

	return web.Context{
		IndexHTML:        mustFind(rootBox, "index.html"),
		RobotsTxt:        mustFind(rootBox, "robots.txt"),
		ServiceWorkerJs:  mustFind(rootBox, "service-worker.js"),
		StaticFileSystem: staticBox,
	}
}

func initServer() (*http.ServeMux, error) {
	apiRouter, err := handlers.NewRouter(&handlers.App{
		Clock:            clock.New(),
		StripeAPIBackend: nil,
		WebURL:           os.Getenv("WebURL"),
	})
	if err != nil {
		return nil, errors.Wrap(err, "initializing router")
	}

	ctx := initContext()

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", apiRouter))
	mux.Handle("/static/", web.GetStaticHandler(ctx.StaticFileSystem))
	mux.HandleFunc("/service-worker.js", web.GetSWHandler(ctx.ServiceWorkerJs))
	mux.HandleFunc("/robots.txt", web.GetRobotsHandler(ctx.RobotsTxt))
	mux.HandleFunc("/", web.GetRootHandler(ctx.IndexHTML))

	return mux, nil
}

func startCmd() {
	mailer.InitTemplates(nil)

	database.Open(database.Config{
		Host:     os.Getenv("DBHost"),
		Port:     os.Getenv("DBPort"),
		Name:     os.Getenv("DBName"),
		User:     os.Getenv("DBUser"),
		Password: os.Getenv("DBPassword"),
	})
	database.InitSchema()
	defer database.Close()

	if err := database.Migrate(); err != nil {
		panic(errors.Wrap(err, "running migrations"))
	}
	if err := job.Run(); err != nil {
		panic(errors.Wrap(err, "running job"))
	}

	srv, err := initServer()
	if err != nil {
		panic(errors.Wrap(err, "initializing server"))
	}

	log.Printf("Dnote version %s is running on port %s", versionTag, *port)
	log.Fatalln(http.ListenAndServe(":"+*port, srv))
}

func versionCmd() {
	fmt.Printf("dnote-server-%s\n", versionTag)
}

func rootCmd() {
	fmt.Printf(`Dnote Server - A simple notebook for developers

Usage:
  dnote-server [command]

Available commands:
  start: Start the server
  version: Print the version
`)
}

func main() {
	flag.Parse()
	cmd := flag.Arg(0)

	switch cmd {
	case "":
		rootCmd()
	case "start":
		startCmd()
	case "version":
		versionCmd()
	default:
		fmt.Printf("Unknown command %s", cmd)
	}
}
