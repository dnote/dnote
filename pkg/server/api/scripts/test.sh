#!/usr/bin/env bash
# test.sh runs api tests. It is to be invoked by other scripts that set
# appropriate env vars.
set -eux

pushd "$GOPATH"/src/github.com/dnote/dnote/pkg/server/api
go test ./... -cover -p 1
popd
