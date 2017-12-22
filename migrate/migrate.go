package migrate

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"

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
	migrationDeleteYAMLArchive
	migrationAddBookMetadata
)

var migrationSequence = []int{
	migrationDeleteYAMLArchive,
	migrationAddBookMetadata,
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
func Migrate() error {
	unrunMigrations, err := getUnrunMigrations()
	if err != nil {
		return errors.Wrap(err, "Failed to get unrun migrations")
	}

	for _, mid := range unrunMigrations {
		if err := performMigration(mid); err != nil {
			return errors.Wrapf(err, "Failed to run migration #%d", mid)
		}
	}

	return nil
}

// performMigration backs up current .dnote data, performs migration, and
// restores or cleans backups depending on if there is an error
func performMigration(migrationID int) error {
	if err := backupDnoteDir(); err != nil {
		return errors.Wrap(err, "Failed to back up dnote directory")
	}

	var migrationError error

	switch migrationID {
	case migrationDeleteYAMLArchive:
		migrationError = deleteDnoteYAMLArchive()
	case migrationAddBookMetadata:
		migrationError = generateBookMetadata()
	default:
		return errors.Errorf("Unrecognized migration id %d", migrationID)
	}

	if migrationError != nil {
		if err := restoreBackup(); err != nil {
			panic(errors.Wrap(err, "Failed to restore backup for a failed migration"))
		}

		return errors.Wrapf(migrationError, "Failed to perform migration #%d", migrationID)
	}

	if err := clearBackup(); err != nil {
		return errors.Wrap(err, "Failed to clear backup")
	}

	if err := updateSchemaVersion(migrationID); err != nil {
		return errors.Wrap(err, "Failed to update schema version")
	}

	return nil
}

// backupDnoteDir backs up the dnote directory to a temporary backup directory
func backupDnoteDir() error {
	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "Failed to get current os user")
	}

	srcPath := fmt.Sprintf("%s/.dnote", usr.HomeDir)
	tmpPath := fmt.Sprintf("%s/%s", usr.HomeDir, backupDirName)

	if err := utils.CopyDir(srcPath, tmpPath); err != nil {
		return errors.Wrap(err, "Failed to copy the .dnote directory")
	}

	return nil
}

func restoreBackup() error {
	var err error

	defer func() {
		if err != nil {
			log.Printf(`Failed to restore backup for a failed migration.
	Don't worry. Your data is still intact in the backup directory.
	Get help on https://github.com/dnote-io/cli/issues`)
		}
	}()

	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "Failed to get current os user")
	}

	srcPath := fmt.Sprintf("%s/.dnote", usr.HomeDir)
	backupPath := fmt.Sprintf("%s/%s", usr.HomeDir, backupDirName)

	if err = os.RemoveAll(srcPath); err != nil {
		return errors.Wrapf(err, "Failed to clear current dnote data at %s", backupPath)
	}

	if err = os.Rename(backupPath, srcPath); err != nil {
		return errors.Wrap(err, `Failed to copy backup data to the original directory.`)
	}

	return nil
}

func clearBackup() error {
	usr, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "Failed to get current os user")
	}

	backupPath := fmt.Sprintf("%s/%s", usr.HomeDir, backupDirName)

	if err := os.RemoveAll(backupPath); err != nil {
		return errors.Wrapf(err, "Failed to remove backup at %s", backupPath)
	}

	return nil
}

// getSchemaPath returns the path to the file containing schema info
func getSchemaPath() (string, error) {
	dnoteDirPath, err := infra.GetDnoteDirPath()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get dnote dir path")
	}

	return fmt.Sprintf("%s/%s", dnoteDirPath, schemaFilename), nil
}

// InitSchemaFile creates a migration file
func InitSchemaFile(pristine bool) error {
	path, err := getSchemaPath()
	if err != nil {
		return errors.Wrap(err, "Failed to get migration file path")
	}

	if utils.FileExists(path) {
		return nil
	}

	s := makeSchema(pristine)
	err = writeSchema(s)
	if err != nil {
		return errors.Wrap(err, "Failed to write schema")
	}

	return nil
}

func readSchema() (schema, error) {
	var ret schema

	path, err := getSchemaPath()
	if err != nil {
		return ret, errors.Wrap(err, "Failed to get schema file path")
	}

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

func writeSchema(s schema) error {
	path, err := getSchemaPath()
	if err != nil {
		return errors.Wrap(err, "Failed to get migration file path")
	}

	d, err := yaml.Marshal(&s)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal schema into yaml")
	}

	if err := ioutil.WriteFile(path, d, 0644); err != nil {
		return errors.Wrap(err, "Failed to write schema file")
	}

	return nil
}

func getUnrunMigrations() ([]int, error) {
	var ret []int

	schema, err := readSchema()
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

func updateSchemaVersion(mID int) error {
	s, err := readSchema()
	if err != nil {
		return errors.Wrap(err, "Failed to read schema")
	}

	s.CurrentVersion = mID

	err = writeSchema(s)
	if err != nil {
		return errors.Wrap(err, "Failed to write schema")
	}

	return nil
}

func getYAMLDnoteArchivePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", usr.HomeDir, ".dnote-yaml-archived"), nil
}
