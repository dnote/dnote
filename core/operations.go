package core

import (
	"database/sql"
	"encoding/base64"

	"github.com/dnote/cli/infra"
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

// TODO: make DB struct that abstracts sql.Tx and sql.DB
// Use the DB strcut in GetSystem instead of *sql.Tx.
// From sync.go, simply pass DB to GetCipherKey and GetValidSession after calling StartTx
// Remove SetupCtx
type DB struct {
	IsTx bool
}

func StartTx(ctx infra.DnoteCtx) DB {

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

// GetCipherKey retrieves the cipher key and decode the base64 into bytes.
func GetCipherKey(ctx infra.DnoteCtx) ([]byte, error) {
	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "beginning transaction")
	}

	var cipherKeyB64 string
	err = GetSystem(tx, infra.SystemCipherKey, &cipherKeyB64)
	if err != nil {
		return []byte{}, errors.Wrap(err, "getting enc key")
	}

	cipherKey, err := base64.StdEncoding.DecodeString(cipherKeyB64)
	if err != nil {
		return nil, errors.Wrap(err, "decoding cipherKey from base64")
	}

	return cipherKey, nil
}
