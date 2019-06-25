#!/bin/bash
# test-local.sh runs api tests using local setting
set -eux

basePath=$GOPATH/src/github.com/dnote/dnote/pkg/server/api

export $(cat "$basePath"/.env.test | xargs)
"$basePath"/scripts/test.sh
