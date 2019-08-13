#!/usr/bin/env bash
# shellcheck disable=SC1090
# dev.sh builds and starts development environment for standalone app
set -eux -o pipefail

# clean up background processes
function cleanup {
  kill "$devServerPID"
}
trap cleanup EXIT

basePath="$GOPATH/src/github.com/dnote/dnote"
appPath="$basePath/web"
serverPath="$basePath/pkg/server"

# load env
set -a
dotenvPath="$serverPath/.env.dev"
source "$dotenvPath"
set +a

# run webpack-dev-server for js in the background
(
  cd "$appPath" &&

  BASE_URL=http://localhost:8080 \
  ASSET_BASE_URL=http://localhost:3000 \
  COMPILED_PATH="$appPath"/compiled \
  PUBLIC_PATH="$appPath"/public \
  STANDALONE=true \
  COMPILED_PATH="$basePath/web/compiled" \
  IS_TEST=true \
    "$appPath"/scripts/webpack-dev.sh
) &
devServerPID=$!

# run server
(cd "$serverPath" && CompileDaemon \
  -command="$serverPath/server start")
