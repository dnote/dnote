#!/usr/bin/env bash
# test.sh runs test files sequentially
# https://stackoverflow.com/questions/23715302/go-how-to-run-tests-for-multiple-packages
set -eux

dir=$(dirname "${BASH_SOURCE[0]}")
pushd "$dir/../../pkg/cli"
# clear tmp dir in case not properly torn down
rm -rf "./tmp"

go test -a ./... \
  -p 1\
  --tags "fts5"
popd
