package migrate

import (
	"encoding/json"
	"github.com/dnote-io/cli/test"
	"github.com/dnote-io/cli/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestMigrateAll(t *testing.T) {
	ctx := test.InitCtx("../tmp")

	// set up
	test.SetupTmp(ctx)
	test.WriteFile(ctx, "./fixtures/2-pre-dnote.json", "dnote")
	if err := InitSchemaFile(ctx, false); err != nil {
		panic(errors.Wrap(err, "Failed to initialize schema file"))
	}
	defer test.ClearTmp(ctx)

	// Execute
	if err := Migrate(ctx); err != nil {
		t.Fatalf("Failed to migrate %s", err.Error())
	}

	// Test
	schema, err := readSchema(ctx)
	if err != nil {
		panic(errors.Wrap(err, "Failed to read the schema"))
	}

	b := test.ReadFile(ctx, "dnote")
	var dnote migrateToV3PostDnote
	if err := json.Unmarshal(b, &dnote); err != nil {
		t.Error(errors.Wrap(err, "Failed to unmarshal result into dnote").Error())
	}

	test.AssertEqual(t, schema.CurrentVersion, len(migrationSequence), "current schema version mismatch")
}

func TestMigrateToV1(t *testing.T) {
	ctx := test.InitCtx("../tmp")

	t.Run("yaml exists", func(t *testing.T) {
		// set up
		test.SetupTmp(ctx)
		defer test.ClearTmp(ctx)

		yamlPath, err := filepath.Abs(filepath.Join(ctx.HomeDir, ".dnote-yaml-archived"))
		if err != nil {
			panic(errors.Wrap(err, "Failed to get absolute YAML path").Error())
		}
		ioutil.WriteFile(yamlPath, []byte{}, 0644)

		// execute
		if err := deleteDnoteYAMLArchive(ctx); err != nil {
			t.Fatal(errors.Wrapf(err, "Failed to migrate").Error())
		}

		// test
		if utils.FileExists(yamlPath) {
			t.Fatal("YAML archive file has not been deleted")
		}
	})

	t.Run("yaml does not exist", func(t *testing.T) {
		// set up
		test.SetupTmp(ctx)
		defer test.ClearTmp(ctx)

		yamlPath, err := filepath.Abs(filepath.Join(ctx.HomeDir, ".dnote-yaml-archived"))
		if err != nil {
			panic(errors.Wrap(err, "Failed to get absolute YAML path").Error())
		}

		// execute
		if err := deleteDnoteYAMLArchive(ctx); err != nil {
			t.Fatal(errors.Wrapf(err, "Failed to migrate").Error())
		}

		// test
		if utils.FileExists(yamlPath) {
			t.Fatal("YAML archive file must not exist")
		}
	})
}

func TestMigrateToV2(t *testing.T) {
	ctx := test.InitCtx("../tmp")

	// set up
	test.SetupTmp(ctx)
	test.WriteFile(ctx, "./fixtures/2-pre-dnote.json", "dnote")
	defer test.ClearTmp(ctx)

	// execute
	if err := migrateToV2(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := test.ReadFile(ctx, "dnote")

	var postDnote migrateToV2PostDnote
	if err := json.Unmarshal(b, &postDnote); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	for _, book := range postDnote {
		test.AssertNotEqual(t, len(book.UUID), 8, "UUID was not generated")

		for _, note := range book.Notes {
			if len(note.UUID) == 8 {
				t.Errorf("Note UUID was not migrated. It has length of %d", len(note.UUID))
			}
		}
	}
}

func TestMigrateToV3(t *testing.T) {
	ctx := test.InitCtx("../tmp")

	// set up
	test.SetupTmp(ctx)
	test.WriteFile(ctx, "./fixtures/3-pre-dnote.json", "dnote")
	defer test.ClearTmp(ctx)

	// execute
	if err := migrateToV3(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := test.ReadFile(ctx, "dnote")
	var postDnote migrateToV3PostDnote
	if err := json.Unmarshal(b, &postDnote); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	b = test.ReadFile(ctx, "actions")
	var actions []migrateToV3Action
	if err := json.Unmarshal(b, &actions); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the actions").Error())
	}

	test.AssertEqual(t, len(actions), 6, "actions length mismatch")
}
