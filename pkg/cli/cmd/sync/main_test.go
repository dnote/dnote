package sync

import (
	"github.com/dnote/dnote/pkg/cli/context"
	"path/filepath"
)

var paths context.Paths = context.Paths{
	Home:   "../../tmp",
	Cache:  "../../tmp",
	Config: "../../tmp",
	Data:   "../../tmp",
}

var dbPath = filepath.Join(paths.Data, ".dnote.db")
