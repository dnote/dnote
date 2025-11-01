/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dnote/dnote/pkg/server/buildinfo"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/controllers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/pkg/errors"
)

func startCmd(args []string) {
	fs := setupFlagSet("start", "dnote-server start")

	port := fs.String("port", "", "Server port (env: PORT, default: 3001)")
	webURL := fs.String("webUrl", "", "Full URL to server without trailing slash (env: WebURL, default: http://localhost:3001)")
	dbPath := fs.String("dbPath", "", "Path to SQLite database file (env: DBPath, default: $XDG_DATA_HOME/dnote/server.db)")
	disableRegistration := fs.Bool("disableRegistration", false, "Disable user registration (env: DisableRegistration, default: false)")
	logLevel := fs.String("logLevel", "", "Log level: debug, info, warn, or error (env: LOG_LEVEL, default: info)")

	fs.Parse(args)

	cfg, err := config.New(config.Params{
		Port:                *port,
		WebURL:              *webURL,
		DBPath:              *dbPath,
		DisableRegistration: *disableRegistration,
		LogLevel:            *logLevel,
	})
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		fs.Usage()
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

	// Start WAL checkpointing to prevent WAL file from growing unbounded.
	database.StartWALCheckpointing(app.DB, 5*time.Minute)

	// Start periodic VACUUM to reclaim space and defragment database.
	database.StartPeriodicVacuum(app.DB, 24*time.Hour)

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
