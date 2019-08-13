#!/usr/bin/env bash
set -eux

basePath="$GOPATH/src/github.com/dnote/dnote/pkg/server/api"

cd "$basePath"
GOOS=linux GOARCH=amd64 go build -o "$basePath/build/api" main.go
