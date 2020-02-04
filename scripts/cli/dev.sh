#!/usr/bin/env bash
# dev.sh builds a new binary and replaces the old one in the PATH with it
set -eux

dir=$(dirname "${BASH_SOURCE[0]}")
sudo rm -rf "$(which dnote)" "$GOPATH/bin/cli"

# change tags to darwin if on macos
go install -ldflags "-X main.apiEndpoint=http://127.0.0.1:3000/api" --tags "linux fts5" "$dir/../../pkg/cli"

sudo ln -s "$GOPATH/bin/cli" /usr/local/bin/dnote
