#!/usr/bin/env bash
# shellcheck disable=SC1090
# test-local.sh runs api tests using local setting
set -eux

dir=$(dirname "${BASH_SOURCE[0]}")

set -a
source "$dir/../../pkg/server/.env.test"
set +a

"$dir/test.sh"
