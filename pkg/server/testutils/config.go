package testutils

import (
	"os"
	"path/filepath"
)

// ProjectPath is the path of the proprietary test suite relative to the "GOPATH"
var ProjectPath string

// CLIPath is the path to the CLI project
var CLIPath string

// ServerPath is the path to the Dnote server project
var ServerPath string

func init() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		panic("GOPATH is not set up")
	}

	// Populate paths
	ProjectPath = filepath.Join(goPath, "src/gitlab.com/monomax/dnote-infra")
	CLIPath = filepath.Join(goPath, "src/github.com/dnote/dnote/pkg/cli")
	ServerPath = filepath.Join(goPath, "src/github.com/dnote/dnote/pkg/server")
}
