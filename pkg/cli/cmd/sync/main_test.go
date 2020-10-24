package sync

import (
	"github.com/dnote/dnote/pkg/cli/context"
	"path/filepath"
)

var dirs context.Dirs = context.Dirs{
	Home:   "../../tmp",
	Cache:  "../../tmp",
	Config: "../../tmp",
	Data:   "../../tmp",
}

var dbPath = filepath.Join(dirs.Data, ".dnote.db")
