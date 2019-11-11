#!/usr/bin/env bash
# shellcheck disable=SC1090
# dev.sh builds and starts development environment
set -eux -o pipefail

# clean up background processes
function cleanup {
  kill "$devServerPID"
}
trap cleanup EXIT

dir=$(dirname $"{BASH_SOURCE[0]}")
basePath="$dir/../.."
appPath="$basePath/web"
serverPath="$basePath/pkg/server"
serverPort=3000

# load env
set -a
dotenvPath="$serverPath/.env.dev"
source "$dotenvPath"
set +a

# run webpack-dev-server for js in the background
(
  cd "$appPath" &&

  BUNDLE_BASE_URL=http://localhost:8080 \
  ASSET_BASE_URL=http://localhost:3000/static \
  ROOT_URL=http://localhost:$serverPort \
  COMPILED_PATH="$appPath"/compiled \
  PUBLIC_PATH="$appPath"/public \
  COMPILED_PATH="$basePath/web/compiled" \
  IS_TEST=true \
  VERSION="$VERSION" \
  WEBPACK_HOST="0.0.0.0" \
    "$appPath"/scripts/webpack-dev.sh
) &
devServerPID=$!

# run server
(cd "$serverPath/watcher" && go run main.go)
