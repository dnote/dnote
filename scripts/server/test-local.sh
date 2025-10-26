#!/usr/bin/env bash
# shellcheck disable=SC1090
# test-local.sh runs api tests using local setting
set -ex

dir=$(dirname "${BASH_SOURCE[0]}")

"$dir/test.sh" "$1"
