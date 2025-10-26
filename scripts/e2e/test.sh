#!/bin/bash
set -eux

dir=$(dirname "${BASH_SOURCE[0]}")
basePath=$(realpath "$dir/../../")

pushd "$basePath"/pkg/e2e
go test --tags "fts5" ./... -v -timeout 5m
popd
