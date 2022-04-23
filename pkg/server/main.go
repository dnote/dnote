/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/buildinfo"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/controllers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/job"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/views"
	"github.com/jinzhu/gorm"

	"github.com/pkg/errors"
)

var pageDir = flag.String("pageDir", "views", "the path to a directory containing page templates")
var staticDir = flag.String("staticDir", "./static/", "the path to the static directory ")

func initDB(c config.Config) *gorm.DB {
	db, err := gorm.Open("postgres", c.DB.GetConnectionStr())
	if err != nil {
		panic(errors.Wrap(err, "opening database connection"))
	}
	database.InitSchema(db)

	return db
}

func mustReadFile(path string) []byte {
	ret, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrap(err, "reading file"))
	}

	return ret
}

func initApp(cfg config.Config) app.App {
	db := initDB(cfg)

	files := map[string][]byte{}
	files[views.ServerErrorPageFileKey] = mustReadFile(fmt.Sprintf("%s/500.html", cfg.StaticDir))

	return app.App{
		DB:             db,
		Clock:          clock.New(),
		EmailTemplates: mailer.NewTemplates(nil),
		EmailBackend:   &mailer.SimpleBackendImplementation{},
		Config:         cfg,
		Files:          files,
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
	cfg := config.Load()
	cfg.SetPageTemplateDir(*pageDir)
	cfg.SetStaticDir(*staticDir)
	cfg.SetAssetBaseURL("/static")

	app := initApp(cfg)
	defer app.DB.Close()

	if err := database.Migrate(app.DB); err != nil {
		panic(errors.Wrap(err, "running migrations"))
	}
	if err := runJob(app); err != nil {
		panic(errors.Wrap(err, "running job"))
	}

	ctl := controllers.New(&app, *pageDir)
	rc := controllers.RouteConfig{
		WebRoutes:   controllers.NewWebRoutes(&app, ctl),
		APIRoutes:   controllers.NewAPIRoutes(&app, ctl),
		Controllers: ctl,
	}

	r, err := controllers.NewRouter(&app, rc)
	if err != nil {
		panic(errors.Wrap(err, "initializing router"))
	}

	log.Printf("Dnote version %s is running on port %s", buildinfo.Version, cfg.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r))
}

func versionCmd() {
	fmt.Printf("dnote-server-%s\n", buildinfo.Version)
}

func rootCmd() {
	fmt.Printf(`Dnote server - a simple personal knowledge base

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
