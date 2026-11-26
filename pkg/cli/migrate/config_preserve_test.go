package migrate

import (
	"fmt"
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/pkg/errors"
)

// A new database runs the migrations again against the existing config file.
func TestConfigMigrationsKeepExistingSettings(t *testing.T) {
	db := database.InitTestMemoryDBRaw(t, "./fixtures/local-12-pre-schema.sql")
	ctx := context.InitTestCtxWithDB(t, db)

	data := "editor: vim\napiEndpoint: https://test.com/api\nenableUpgradeCheck: false\n"
	path := fmt.Sprintf("%s/%s/dnoterc", ctx.Paths.Config, consts.DnoteDirName)
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(errors.Wrap(err, "writing config"))
	}

	for _, m := range []migration{lm12, lm13} {
		if err := m.run(ctx, nil); err != nil {
			t.Fatal(errors.Wrap(err, "running migration"))
		}
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(errors.Wrap(err, "reading config"))
	}
	assert.Equal(t, string(b), data, "config file changed")
}
