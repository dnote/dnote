package migrate

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"

	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/utils"
)

var (
	schemaFilename = "schema"
	backupDirName  = ".dnote-bak"
)

// migration IDs
const (
	_ = iota
	migrationV1
	migrationV2
	migrationV3
)

var migrationSequence = []int{
	migrationV1,
	migrationV2,
	migrationV3,
}

type schema struct {
	CurrentVersion int `yaml:"current_version"`
}

func makeSchema(complete bool) schema {
	s := schema{}

	var currentVersion int
	if complete {
		currentVersion = len(migrationSequence)
	}

	s.CurrentVersion = currentVersion

	return s
}

// Migrate determines migrations to be run and performs them
func Migrate(ctx infra.DnoteCtx) error {
	unrunMigrations, err := getUnrunMigrations(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get unrun migrations")
	}

	for _, mid := range unrunMigrations {
		if err := performMigration(ctx, mid); err != nil {
			return errors.Wrapf(err, "Failed to run migration #%d", mid)
		}
	}

	return nil
}

// performMigration backs up current .dnote data, performs migration, and
// restores or cleans backups depending on if there is an error
func performMigration(ctx infra.DnoteCtx, migrationID int) error {
	if err := backupDnoteDir(ctx); err != nil {
		return errors.Wrap(err, "Failed to back up dnote directory")
	}

	var migrationError error

	switch migrationID {
	case migrationV1:
		migrationError = deleteDnoteYAMLArchive(ctx)
	case migrationV2:
		migrationError = migrateToV2(ctx)
	case migrationV3:
		migrationError = migrateToV3(ctx)
	default:
		return errors.Errorf("Unrecognized migration id %d", migrationID)
	}

	if migrationError != nil {
		if err := restoreBackup(ctx); err != nil {
			panic(errors.Wrap(err, "Failed to restore backup for a failed migration"))
		}

		return errors.Wrapf(migrationError, "Failed to perform migration #%d", migrationID)
	}

	if err := clearBackup(ctx); err != nil {
		return errors.Wrap(err, "Failed to clear backup")
	}

	if err := updateSchemaVersion(ctx, migrationID); err != nil {
		return errors.Wrap(err, "Failed to update schema version")
	}

	return nil
}

// backupDnoteDir backs up the dnote directory to a temporary backup directory
func backupDnoteDir(ctx infra.DnoteCtx) error {
	srcPath := fmt.Sprintf("%s/.dnote", ctx.HomeDir)
	tmpPath := fmt.Sprintf("%s/%s", ctx.HomeDir, backupDirName)

	if err := utils.CopyDir(srcPath, tmpPath); err != nil {
		return errors.Wrap(err, "Failed to copy the .dnote directory")
	}

	return nil
}

func restoreBackup(ctx infra.DnoteCtx) error {
	var err error

	defer func() {
		if err != nil {
			log.Printf(`Failed to restore backup for a failed migration.
	Don't worry. Your data is still intact in the backup directory.
	Get help on https://github.com/dnote-io/cli/issues`)
		}
	}()

	srcPath := fmt.Sprintf("%s/.dnote", ctx.HomeDir)
	backupPath := fmt.Sprintf("%s/%s", ctx.HomeDir, backupDirName)

	if err = os.RemoveAll(srcPath); err != nil {
		return errors.Wrapf(err, "Failed to clear current dnote data at %s", backupPath)
	}

	if err = os.Rename(backupPath, srcPath); err != nil {
		return errors.Wrap(err, `Failed to copy backup data to the original directory.`)
	}

	return nil
}

func clearBackup(ctx infra.DnoteCtx) error {
	backupPath := fmt.Sprintf("%s/%s", ctx.HomeDir, backupDirName)

	if err := os.RemoveAll(backupPath); err != nil {
		return errors.Wrapf(err, "Failed to remove backup at %s", backupPath)
	}

	return nil
}

// getSchemaPath returns the path to the file containing schema info
func getSchemaPath(ctx infra.DnoteCtx) string {
	return fmt.Sprintf("%s/%s", ctx.DnoteDir, schemaFilename)
}

// InitSchemaFile creates a migration file
func InitSchemaFile(ctx infra.DnoteCtx, pristine bool) error {
	path := getSchemaPath(ctx)
	if utils.FileExists(path) {
		return nil
	}

	s := makeSchema(pristine)
	err := writeSchema(ctx, s)
	if err != nil {
		return errors.Wrap(err, "Failed to write schema")
	}

	return nil
}

func readSchema(ctx infra.DnoteCtx) (schema, error) {
	var ret schema

	path := getSchemaPath(ctx)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to read schema file")
	}

	err = yaml.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func writeSchema(ctx infra.DnoteCtx, s schema) error {
	path := getSchemaPath(ctx)
	d, err := yaml.Marshal(&s)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal schema into yaml")
	}

	if err := ioutil.WriteFile(path, d, 0644); err != nil {
		return errors.Wrap(err, "Failed to write schema file")
	}

	return nil
}

func getUnrunMigrations(ctx infra.DnoteCtx) ([]int, error) {
	var ret []int

	schema, err := readSchema(ctx)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to read schema")
	}

	if schema.CurrentVersion == len(migrationSequence) {
		return ret, nil
	}

	nextVersion := schema.CurrentVersion
	ret = migrationSequence[nextVersion:]

	return ret, nil
}

func updateSchemaVersion(ctx infra.DnoteCtx, mID int) error {
	s, err := readSchema(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to read schema")
	}

	s.CurrentVersion = mID

	err = writeSchema(ctx, s)
	if err != nil {
		return errors.Wrap(err, "Failed to write schema")
	}

	return nil
}

func getYAMLDnoteArchivePath(ctx infra.DnoteCtx) (string, error) {
	return fmt.Sprintf("%s/%s", ctx.HomeDir, ".dnote-yaml-archived"), nil
}
