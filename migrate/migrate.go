package migrate

import (
	"database/sql"

	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
)

type migration struct {
	name string
	sql  string
}

var migrations = []migration{}

func initSchema(db *sql.DB) (int, error) {
	schemaVersion := 0

	_, err := db.Exec("INSERT INTO system (key, value) VALUES (? ,?)", "schema", schemaVersion)
	if err != nil {
		return schemaVersion, errors.Wrap(err, "inserting schema")
	}

	return schemaVersion, nil
}

func getSchema(db *sql.DB) (int, error) {
	var ret int

	err := db.QueryRow("SELECT value FROM system where key = ?", "schema").Scan(&ret)
	if err == sql.ErrNoRows {
		ret, err = initSchema(db)

		if err != nil {
			return ret, errors.Wrap(err, "initializing schema")
		}
	} else if err != nil {
		return ret, errors.Wrap(err, "querying schema")
	}

	return ret, nil
}

func execute(ctx infra.DnoteCtx, nextSchema int, m migration) error {
	log.Debug("running migration %s\n", m.name)

	tx, err := ctx.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	_, err = tx.Exec(m.sql)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "running sql")
	}

	_, err = tx.Exec("UPDATE system SET value = ? WHERE key = ?", nextSchema, "schema")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "incrementing schema")
	}

	tx.Commit()

	return nil
}

// Run performs unrun migrations
func Run(ctx infra.DnoteCtx) error {
	db := ctx.DB

	schema, err := getSchema(db)
	if err != nil {
		return errors.Wrap(err, "getting the current schema")
	}

	log.Debug("current schema %d\n", schema)

	if schema == len(migrations) {
		return nil
	}

	toRun := migrations[schema:]

	for idx, m := range toRun {
		nextSchema := schema + idx + 1
		if err := execute(ctx, nextSchema, m); err != nil {
			return errors.Wrapf(err, "running migration %s", m.name)
		}
	}

	return nil
}
