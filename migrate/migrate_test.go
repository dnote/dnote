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

func TestDeleteDnoteYAMLArchive(t *testing.T) {
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

func TestGenerateBookMetadata(t *testing.T) {
	ctx := test.InitCtx("../tmp")

	// set up
	test.SetupTmp(ctx)
	test.WriteFile(ctx, "./fixtures/2-pre-dnote.json", "dnote")
	defer test.ClearTmp(ctx)

	// execute
	if err := generateBookMetadata(ctx); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to migrate").Error())
	}

	// test
	b := test.ReadFile(ctx, "dnote")

	var postDnote generateBookMetadataPostDnote
	if err := json.Unmarshal(b, &postDnote); err != nil {
		t.Fatal(errors.Wrap(err, "Failed to unmarshal the result into Dnote").Error())
	}

	for _, book := range postDnote {
		if book.UID == "" {
			t.Error("UID has not been generated")
		}
	}
}
