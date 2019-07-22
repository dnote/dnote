package database

import (
	"log"

	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

// Migrate runs the migrations
func Migrate() error {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "../database/migrations/"),
	}

	migrate.SetTable(MigrationTableName)

	db := DBConn.DB()
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "running migrations")
	}

	log.Printf("Performed %d migrations", n)

	return nil
}
