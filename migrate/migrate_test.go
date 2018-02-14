package migrate

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/dnote-io/cli/testutils"
	"github.com/dnote-io/cli/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func TestMigrateAll(t *testing.T) {
	ctx := testutils.InitCtx("../tmp")

	// set up
	testutils.SetupTmp(ctx)
	testutils.WriteFile(ctx, "./fixtures/2-pre-dnote.json", "dnote")
	testutils.WriteFile(ctx, "./fixtures/4-pre-dnoterc.yaml", "dnoterc")
	if err := InitSchemaFile(ctx, false); err != nil {
		panic(errors.Wrap(err, "Failed to initialize schema file"))
	}
	defer testutils.ClearTmp(ctx)

	// Execute
	if err := Migrate(ctx); err != nil {
		t.Fatalf("Failed to migrate %s", err.Error())
	}

	// Test
	schema, err := readSchema(ctx)
	if err != nil {
		panic(errors.Wrap(err, "Failed to read the schema"))
	}

	b := testutils.ReadFile(ctx, "dnote")
	var dnote migrateToV3Dnote
	if err := json.Unmarshal(b, &dnote); err != nil {
		t.Error(errors.Wrap(err, "Failed to unmarshal result into dnote").Error())
	}

	testutils.AssertEqual(t, schema.CurrentVersion, len(migrationSequence), "current schema version mismatch")

	note := dnote["algorithm"].Notes[0]
	testutils.AssertEqual(t, note.Content, "in-place means no extra space required. it mutates the input", "content was not carried over")
	testutils.AssertNotEqual(t, note.UUID, "", "note uuid was not generated")
	testutils.AssertNotEqual(t, note.AddedOn, int64(0), "note added_on was not generated")
	testutils.AssertEqual(t, note.EditedOn, int64(0), "note edited_on was not propertly generated")
}

func TestMigrateToV1(t *testing.T) {
	ctx := testutils.InitCtx("../tmp")

	t.Run("yaml exists", func(t *testing.T) {
		// set up
		testutils.SetupTmp(ctx)
		defer testutils.ClearTmp(ctx)

		yamlPath, err := filepath.Abs(filepath.Join(ctx.HomeDir, ".dnote-yaml-archived"))
		if err != nil {
			panic(errors.Wrap(err, "Failed to get absolute YAML path").Error())
		}
		ioutil.WriteFile(yamlPath, []byte{}, 0644)

		// execute
		if err := migrateToV1(ctx); err != nil {
			t.Fatal(errors.Wrapf(err, "Failed to migrate").Error())
		}

		// test
		if utils.FileExists(yamlPath) {
			t.Fatal("YAML archive file has not been deleted")
		}
	})

	t.Run("yaml does not exist", func(t *testing.T) {
		// set up
		testutils.SetupTmp(ctx)
		defer testutils.ClearTmp(ctx)

		yamlPath, err := filepath.Abs(filepath.Join(ctx.HomeDir, ".dnote-yaml-archived"))
		if err != nil {
			panic(errors.Wrap(err, "Failed to get absolute YAML path").Error())
		}

		// execute
		if err := migrateToV1(ctx); err != nil {
			t.Fatal(errors.Wrapf(err, "Failed to migrate").Error())
		}

		// test
		if utils.FileExists(yamlPath) {
			t.Fatal("YAML archive file must not exist")
		}
	})
}

func TestMigrateToV2(t *testing.T) {
	ctx := testutils.InitCtx("../tmp")

	// set up
	testutils.SetupTmp(ctx)
	testutils.WriteFile(ctx, "./fixtures/2-pre-dnote.json", "dnote")
	defer testutils.ClearTmp(ctx)

	// execute
	if err := migrateToV2(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := testutils.ReadFile(ctx, "dnote")

	var postDnote migrateToV2PostDnote
	if err := json.Unmarshal(b, &postDnote); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	for _, book := range postDnote {
		testutils.AssertNotEqual(t, book.Name, "", "Book name was not populated")

		for _, note := range book.Notes {
			if len(note.UUID) == 8 {
				t.Errorf("Note UUID was not migrated. It has length of %d", len(note.UUID))
			}

			testutils.AssertNotEqual(t, note.AddedOn, int64(0), "AddedOn was not carried over")
			testutils.AssertEqual(t, note.EditedOn, int64(0), "EditedOn was not created properly")
		}
	}
}

func TestMigrateToV3(t *testing.T) {
	ctx := testutils.InitCtx("../tmp")

	// set up
	testutils.SetupTmp(ctx)
	testutils.WriteFile(ctx, "./fixtures/3-pre-dnote.json", "dnote")
	defer testutils.ClearTmp(ctx)

	// execute
	if err := migrateToV3(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := testutils.ReadFile(ctx, "dnote")
	var postDnote migrateToV3Dnote
	if err := json.Unmarshal(b, &postDnote); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	b = testutils.ReadFile(ctx, "actions")
	var actions []migrateToV3Action
	if err := json.Unmarshal(b, &actions); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the actions").Error())
	}

	testutils.AssertEqual(t, len(actions), 6, "actions length mismatch")

	for _, book := range postDnote {
		for _, note := range book.Notes {
			testutils.AssertNotEqual(t, note.AddedOn, int64(0), "AddedOn was not carried over")
		}
	}
}

func TestMigrateToV4(t *testing.T) {
	ctx := testutils.InitCtx("../tmp")

	// set up
	testutils.SetupTmp(ctx)
	testutils.WriteFile(ctx, "./fixtures/4-pre-dnoterc.yaml", "dnoterc")
	defer testutils.ClearTmp(ctx)
	defer os.Setenv("EDITOR", "")

	// execute
	os.Setenv("EDITOR", "vim")
	if err := migrateToV4(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := testutils.ReadFile(ctx, "dnoterc")
	var config migrateToV4PostConfig
	if err := yaml.Unmarshal(b, &config); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	testutils.AssertEqual(t, config.APIKey, "Oev6e1082ORasdf9rjkfjkasdfjhgei", "api key mismatch")
	testutils.AssertEqual(t, config.Editor, "vim", "editor mismatch")
}
