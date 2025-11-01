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

func getEmailBackend() mailer.Backend {
	defaultBackend, err := mailer.NewDefaultBackend()
	if err != nil {
		log.Debug("SMTP not configured, using StdoutBackend for emails")
		return mailer.NewStdoutBackend()
	}

	log.Debug("Email backend configured")
	return defaultBackend
}

func initApp(cfg config.Config) app.App {
	db := initDB(cfg.DBPath)
	emailBackend := getEmailBackend()

	return app.App{
		DB:                  db,
		Clock:               clock.New(),
		EmailBackend:        emailBackend,
		HTTP500Page:         cfg.HTTP500Page,
		BaseURL:             cfg.BaseURL,
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
