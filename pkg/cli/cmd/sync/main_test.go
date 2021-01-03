package sync

import (
	"github.com/dnote/dnote/pkg/cli/context"
	"path/filepath"
)

var testDir = "../../tmp"

var paths context.Paths = context.Paths{
	Home:   testDir,
	Cache:  testDir,
	Config: testDir,
	Data:   testDir,
}

var dbPath = filepath.Join(testDir, "test.db")
