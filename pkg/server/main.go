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
	"net/http"
	"os"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/buildinfo"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/controllers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func initDB(dbPath string) *gorm.DB {
	db := database.Open(dbPath)
	database.InitSchema(db)
	database.Migrate(db)

	return db
}

func initApp(cfg config.Config) app.App {
	db := initDB(cfg.DBPath)

	emailBackend, err := mailer.NewDefaultBackend(cfg.IsProd())
	if err != nil {
		emailBackend = &mailer.DefaultBackend{Enabled: false}
	} else {
		log.Info("Email backend configured")
	}

	return app.App{
		DB:                  db,
		Clock:               clock.New(),
		EmailTemplates:      mailer.NewTemplates(),
		EmailBackend:        emailBackend,
		HTTP500Page:         cfg.HTTP500Page,
		AppEnv:              cfg.AppEnv,
		WebURL:              cfg.WebURL,
		DisableRegistration: cfg.DisableRegistration,
		Port:                cfg.Port,
		DBPath:              cfg.DBPath,
		AssetBaseURL:        cfg.AssetBaseURL,
	}
}

func startCmd(args []string) {
	startFlags := flag.NewFlagSet("start", flag.ExitOnError)
	startFlags.Usage = func() {
		fmt.Printf(`Usage:
  dnote-server start [flags]

Flags:
`)
		startFlags.PrintDefaults()
	}

	appEnv := startFlags.String("appEnv", "", "Application environment (env: APP_ENV, default: PRODUCTION)")
	port := startFlags.String("port", "", "Server port (env: PORT, default: 3000)")
	webURL := startFlags.String("webUrl", "", "Full URL to server without trailing slash (env: WebURL, example: https://example.com)")
	dbPath := startFlags.String("dbPath", "", "Path to SQLite database file (env: DBPath, default: $XDG_DATA_HOME/dnote/server.db)")
	disableRegistration := startFlags.Bool("disableRegistration", false, "Disable user registration (env: DisableRegistration, default: false)")
	logLevel := startFlags.String("logLevel", "", "Log level: debug, info, warn, or error (env: LOG_LEVEL, default: info)")

	startFlags.Parse(args)

	cfg, err := config.New(config.Params{
		AppEnv:              *appEnv,
		Port:                *port,
		WebURL:              *webURL,
		DBPath:              *dbPath,
		DisableRegistration: *disableRegistration,
		LogLevel:            *logLevel,
	})
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		startFlags.Usage()
		os.Exit(1)
	}

	// Set log level
	log.SetLevel(cfg.LogLevel)

	app := initApp(cfg)
	defer func() {
		sqlDB, err := app.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

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

	log.WithFields(log.Fields{
		"version": buildinfo.Version,
		"port":    cfg.Port,
	}).Info("Dnote server starting")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r); err != nil {
		log.ErrorWrap(err, "server failed")
		os.Exit(1)
	}
}

func versionCmd() {
	fmt.Printf("dnote-server-%s\n", buildinfo.Version)
}

func rootCmd() {
	fmt.Printf(`Dnote server - a simple command line notebook

Usage:
  dnote-server [command] [flags]

Available commands:
  start: Start the server (use 'dnote-server start --help' for flags)
  version: Print the version
`)
}

func main() {
	if len(os.Args) < 2 {
		rootCmd()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "start":
		startCmd(os.Args[2:])
	case "version":
		versionCmd()
	default:
		fmt.Printf("Unknown command %s\n", cmd)
		rootCmd()
		os.Exit(1)
	}
}
