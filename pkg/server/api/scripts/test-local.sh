#!/usr/bin/env bash
# shellcheck disable=SC1090
# test-local.sh runs api tests using local setting
set -eux

basePath=$GOPATH/src/github.com/dnote/dnote/pkg/server

set -a
source "$basePath/.env.test"
set +a

"$basePath/api/scripts/test.sh"
