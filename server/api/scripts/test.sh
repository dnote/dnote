#!/bin/bash
# test.sh runs api tests. It is to be invoked by other scripts that set
# appropriate env vars.
set -eux

pushd "$GOPATH"/src/github.com/dnote/dnote/server/api
go test ./handlers/... ./operations/...  -cover -p 1
popd
