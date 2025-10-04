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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/buildinfo"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/controllers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/job"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var port = flag.String("port", "3000", "port to connect to")

func initDB(c config.Config) *gorm.DB {
	db := database.Open(c.DB.Path)
	database.InitSchema(db)
	database.Migrate(db)

	return db
}

func initApp(cfg config.Config) app.App {
	db := initDB(cfg)

	isProduction := os.Getenv("GO_ENV") == "PRODUCTION"
	emailBackend, err := mailer.NewDefaultBackend(isProduction)
	if err != nil {
		log.Printf("Email backend not fully configured: %v. Emails will only be logged.", err)
		emailBackend = &mailer.DefaultBackend{Enabled: false}
	}

	return app.App{
		DB:             db,
		Clock:          clock.New(),
		EmailTemplates: mailer.NewTemplates(),
		EmailBackend:   emailBackend,
		Config:         cfg,
		HTTP500Page:    cfg.HTTP500Page,
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
	cfg.SetAssetBaseURL("/static")

	app := initApp(cfg)
	defer func() {
		sqlDB, err := app.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	if err := database.Migrate(app.DB); err != nil {
		panic(errors.Wrap(err, "running migrations"))
	}
	if err := runJob(app); err != nil {
		panic(errors.Wrap(err, "running job"))
	}

	ctl := controllers.New(&app)
	rc := controllers.RouteConfig{
		WebRoutes:   controllers.NewWebRoutes(&app, ctl),
		APIRoutes:   controllers.NewAPIRoutes(&app, ctl),
		Controllers: ctl,
	}

	r, err := controllers.NewRouter(&app, rc)
	if err != nil {
		panic(errors.Wrap(err, "initializing router"))
	}

	log.Printf("Dnote version %s is running on port %s", buildinfo.Version, *port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", *port), r))
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
