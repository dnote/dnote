#!/usr/bin/env bash
# dev.sh builds a new binary and replaces the old one in the PATH with it
set -eux

sudo rm -rf "$(which dnote)" "$GOPATH/bin/cli"

# change tags to darwin if on macos
go install -ldflags "-X main.apiEndpoint=http://127.0.0.1:5000" --tags "linux fts5" "$GOPATH/src/github.com/dnote/dnote/pkg/cli/."

sudo ln -s "$GOPATH/bin/cli" /usr/local/bin/dnote
