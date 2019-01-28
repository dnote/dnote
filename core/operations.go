package core

import (
	"database/sql"

	"github.com/pkg/errors"
)

// InsertSystem inserets a system configuration
func InsertSystem(tx *sql.Tx, key, val string) error {
	if _, err := tx.Exec("INSERT INTO system (key, value) VALUES (? , ?);", key, val); err != nil {
		return errors.Wrap(err, "saving system config")
	}

	return nil
}

// UpsertSystem inserts or updates a system configuration
func UpsertSystem(tx *sql.Tx, key, val string) error {
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM system WHERE key = ?", key).Scan(&count); err != nil {
		return errors.Wrap(err, "counting system record")
	}

	if count == 0 {
		if _, err := tx.Exec("INSERT INTO system (key, value) VALUES (? , ?);", key, val); err != nil {
			return errors.Wrap(err, "saving system config")
		}
	} else {
		if _, err := tx.Exec("UPDATE system SET value = ? WHERE key = ?", val, key); err != nil {
			return errors.Wrap(err, "updating system config")
		}
	}

	return nil
}

// UpdateSystem updates a system configuration
func UpdateSystem(tx *sql.Tx, key, val interface{}) error {
	if _, err := tx.Exec("UPDATE system SET value = ? WHERE key = ?", val, key); err != nil {
		return errors.Wrap(err, "updating system config")
	}

	return nil
}

// GetSystem scans the given system configuration record onto the destination
func GetSystem(tx *sql.Tx, key string, dest interface{}) error {
	if err := tx.QueryRow("SELECT value FROM system WHERE key = ?", key).Scan(dest); err != nil {
		return errors.Wrap(err, "finding system configuration record")
	}

	return nil
}

// DeleteSystem delets the given system record
func DeleteSystem(tx *sql.Tx, key string) error {
	if _, err := tx.Exec("DELETE FROM system WHERE key = ?", key); err != nil {
		return errors.Wrap(err, "deleting system config")
	}

	return nil
}
