#!/bin/bash
set -eux

dir=$(dirname "${BASH_SOURCE[0]}")
basePath=$(realpath "$dir/../../")

set -a
source "$basePath/pkg/server/.env.test"
set +a

pushd "$basePath"/pkg/e2e
go test --tags "fts5" ./... -p 1 -v -timeout 5m
popd
