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

package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/mailer"
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

// printFlags prints flags with -- prefix for consistency with CLI
func printFlags(fs *flag.FlagSet) {
	fs.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  --%s", f.Name)

		// Print type hint for non-boolean flags
		name, usage := flag.UnquoteUsage(f)
		if name != "" {
			fmt.Printf(" %s", name)
		}
		fmt.Println()

		// Print usage description with indentation
		if usage != "" {
			fmt.Printf("    \t%s", usage)
			if f.DefValue != "" && f.DefValue != "false" {
				fmt.Printf(" (default: %s)", f.DefValue)
			}
			fmt.Println()
		}
	})
}

// setupFlagSet creates a FlagSet with standard usage format
func setupFlagSet(name, usageCmd string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf(`Usage:
  %s [flags]

Flags:
`, usageCmd)
		printFlags(fs)
	}
	return fs
}

// requireString validates that a required string flag is not empty
func requireString(fs *flag.FlagSet, value, fieldName string) {
	if value == "" {
		fmt.Printf("Error: %s is required\n", fieldName)
		fs.Usage()
		os.Exit(1)
	}
}

// setupAppWithDB creates config, initializes app, and returns cleanup function
func setupAppWithDB(fs *flag.FlagSet, dbPath string) (*app.App, func()) {
	cfg, err := config.New(config.Params{
		DBPath: dbPath,
	})
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		fs.Usage()
		os.Exit(1)
	}

	a := initApp(cfg)
	cleanup := func() {
		sqlDB, err := a.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	return &a, cleanup
}
