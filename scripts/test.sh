#!/bin/bash
# run_server_test.sh runs server test files sequentially
# https://stackoverflow.com/questions/23715302/go-how-to-run-tests-for-multiple-packages

set -e

# clear tmp dir in case not properly torn down
rm -rf $GOPATH/src/github.com/dnote/cli/tmp

# run test
pushd $GOPATH/src/github.com/dnote/cli

go test ./... \
  -p 1\
  --tags "fts5"

popd
