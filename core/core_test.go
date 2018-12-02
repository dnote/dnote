package core

import (
	"testing"

	"github.com/dnote/cli/testutils"
	"github.com/pkg/errors"
)

func TestInitSystemKV(t *testing.T) {
	// Setup
	ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", true)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB

	var originalCount int
	testutils.MustScan(t, "counting system configs", db.QueryRow("SELECT count(*) FROM system"), &originalCount)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	if err := initSystemKV(tx, "testKey", "testVal"); err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "executing"))
	}

	tx.Commit()

	// Test
	var count int
	testutils.MustScan(t, "counting system configs", db.QueryRow("SELECT count(*) FROM system"), &count)
	testutils.AssertEqual(t, count, originalCount+1, "system count mismatch")

	var val string
	testutils.MustScan(t, "getting system value",
		db.QueryRow("SELECT value FROM system WHERE key = ?", "testKey"), &val)
	testutils.AssertEqual(t, val, "testVal", "system value mismatch")
}

func TestInitSystemKV_existing(t *testing.T) {
	// Setup
	ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", true)
	defer testutils.TeardownEnv(ctx)

	db := ctx.DB
	testutils.MustExec(t, "inserting a system config", db, "INSERT INTO system (key, value) VALUES (?, ?)", "testKey", "testVal")

	var originalCount int
	testutils.MustScan(t, "counting system configs", db.QueryRow("SELECT count(*) FROM system"), &originalCount)

	// Execute
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(errors.Wrap(err, "beginning a transaction"))
	}

	if err := initSystemKV(tx, "testKey", "newTestVal"); err != nil {
		tx.Rollback()
		t.Fatal(errors.Wrap(err, "executing"))
	}

	tx.Commit()

	// Test
	var count int
	testutils.MustScan(t, "counting system configs", db.QueryRow("SELECT count(*) FROM system"), &count)
	testutils.AssertEqual(t, count, originalCount, "system count mismatch")

	var val string
	testutils.MustScan(t, "getting system value",
		db.QueryRow("SELECT value FROM system WHERE key = ?", "testKey"), &val)
	testutils.AssertEqual(t, val, "testVal", "system value should not have been updated")
}
