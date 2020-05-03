/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/api"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/job"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/web"
	"github.com/jinzhu/gorm"

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

func initWebContext(db *gorm.DB) web.Context {
	staticBox := packr.New("static", "../../web/public/static")

	return web.Context{
		DB:               db,
		IndexHTML:        mustFind(rootBox, "index.html"),
		RobotsTxt:        mustFind(rootBox, "robots.txt"),
		ServiceWorkerJs:  mustFind(rootBox, "service-worker.js"),
		StaticFileSystem: staticBox,
	}
}

func initServer(a app.App) (*http.ServeMux, error) {
	apiRouter, err := api.NewRouter(&api.API{App: &a})
	if err != nil {
		return nil, errors.Wrap(err, "initializing router")
	}

	webCtx := initWebContext(a.DB)
	webHandlers, err := web.Init(webCtx)
	if err != nil {
		return nil, errors.Wrap(err, "initializing web handlers")
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", apiRouter))
	mux.Handle("/static/", webHandlers.GetStatic)
	mux.HandleFunc("/service-worker.js", webHandlers.GetServiceWorker)
	mux.HandleFunc("/robots.txt", webHandlers.GetRobots)
	mux.HandleFunc("/", webHandlers.GetRoot)

	return mux, nil
}

func initDB(c config.Config) *gorm.DB {
	db, err := gorm.Open("postgres", c.DB.GetConnectionStr())
	if err != nil {
		panic(errors.Wrap(err, "opening database connection"))
	}
	database.InitSchema(db)

	return db
}

func initApp(c config.Config) app.App {
	db := initDB(c)

	return app.App{
		DB:             db,
		Clock:          clock.New(),
		EmailTemplates: mailer.NewTemplates(nil),
		EmailBackend:   &mailer.SimpleBackendImplementation{},
		Config:         c,
	}
}

func runJob(a app.App) error {
	runner, err := job.NewRunner(a.DB, a.Clock, a.EmailTemplates, a.EmailBackend, a.Config)
	if err != nil {
		return errors.Wrap(err, "getting a job runner")
	}
	if err := runner.Do(); err != nil {
		return errors.Wrap(err, "running job")
	}

	return nil
}

func startCmd() {
	c := config.Load()

	app := initApp(c)
	defer app.DB.Close()

	if err := database.Migrate(app.DB); err != nil {
		panic(errors.Wrap(err, "running migrations"))
	}

	if err := runJob(app); err != nil {
		panic(errors.Wrap(err, "running job"))
	}

	srv, err := initServer(app)
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
	fmt.Printf(`Dnote Server - A simple personal knowledge base

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
