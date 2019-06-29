package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

var (
	migrationDir = flag.String("migrationDir", "../migrations", "the path to the directory with migraiton files")
)

func init() {
	fmt.Println("Migrating Dnote database...")

	// Load env
	if os.Getenv("GO_ENV") != "PRODUCTION" {
		if err := godotenv.Load("../../api/.env.dev"); err != nil {
			panic(err)
		}
	}

	database.InitDB()
}

func main() {
	flag.Parse()

	db := database.DBConn

	migrations := &migrate.FileMigrationSource{
		Dir: *migrationDir,
	}

	n, err := migrate.Exec(db.DB(), "postgres", migrations, migrate.Up)
	if err != nil {
		panic(errors.Wrap(err, "executing migrations"))
	}

	fmt.Printf("Applied %d migrations\n", n)
}
