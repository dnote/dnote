#!/usr/bin/env bash
# test.sh runs tests for CLI packages
set -eux

dir=$(dirname "${BASH_SOURCE[0]}")
pushd "$dir/../../pkg/cli"
# clear tmp dir in case not properly torn down
rm -rf "./tmp"

go test ./... --tags "fts5"
popd
