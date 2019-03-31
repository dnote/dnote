package core

import (
	"encoding/base64"

	"github.com/dnote/cli/infra"
	"github.com/pkg/errors"
)

// InsertSystem inserets a system configuration
func InsertSystem(db *infra.DB, key, val string) error {
	if _, err := db.Exec("INSERT INTO system (key, value) VALUES (? , ?);", key, val); err != nil {
		return errors.Wrap(err, "saving system config")
	}

	return nil
}

// UpsertSystem inserts or updates a system configuration
func UpsertSystem(db *infra.DB, key, val string) error {
	var count int
	if err := db.QueryRow("SELECT count(*) FROM system WHERE key = ?", key).Scan(&count); err != nil {
		return errors.Wrap(err, "counting system record")
	}

	if count == 0 {
		if _, err := db.Exec("INSERT INTO system (key, value) VALUES (? , ?);", key, val); err != nil {
			return errors.Wrap(err, "saving system config")
		}
	} else {
		if _, err := db.Exec("UPDATE system SET value = ? WHERE key = ?", val, key); err != nil {
			return errors.Wrap(err, "updating system config")
		}
	}

	return nil
}

// UpdateSystem updates a system configuration
func UpdateSystem(db *infra.DB, key, val interface{}) error {
	if _, err := db.Exec("UPDATE system SET value = ? WHERE key = ?", val, key); err != nil {
		return errors.Wrap(err, "updating system config")
	}

	return nil
}

// GetSystem scans the given system configuration record onto the destination
func GetSystem(db *infra.DB, key string, dest interface{}) error {
	if err := db.QueryRow("SELECT value FROM system WHERE key = ?", key).Scan(dest); err != nil {
		return errors.Wrap(err, "finding system configuration record")
	}

	return nil
}

// DeleteSystem delets the given system record
func DeleteSystem(db *infra.DB, key string) error {
	if _, err := db.Exec("DELETE FROM system WHERE key = ?", key); err != nil {
		return errors.Wrap(err, "deleting system config")
	}

	return nil
}

// GetCipherKey retrieves the cipher key and decode the base64 into bytes.
func GetCipherKey(ctx infra.DnoteCtx) ([]byte, error) {
	db, err := ctx.DB.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "beginning transaction")
	}

	var cipherKeyB64 string
	err = GetSystem(db, infra.SystemCipherKey, &cipherKeyB64)
	if err != nil {
		return []byte{}, errors.Wrap(err, "getting enc key")
	}

	cipherKey, err := base64.StdEncoding.DecodeString(cipherKeyB64)
	if err != nil {
		return nil, errors.Wrap(err, "decoding cipherKey from base64")
	}

	return cipherKey, nil
}
