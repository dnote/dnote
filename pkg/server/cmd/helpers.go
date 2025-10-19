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
