#!/bin/bash
# dev.sh builds and starts development environment for standalone app
set -eux -o pipefail

basePath="$GOPATH/src/github.com/dnote/dnote"
appPath="$basePath"/web

# run webpack-dev-server for js
(
  cd "$appPath" &&

  BASE_URL=http://localhost:8080 \
  ASSET_BASE_URL=http://localhost:3000 \
  COMPILED_PATH="$appPath"/compiled \
  PUBLIC_PATH="$appPath"/public \
  STANDALONE=true \
  COMPILED_PATH="$basePath/web/compiled" \
    "$appPath"/scripts/webpack-dev.sh
) &

# run server
(cd "$appPath" && PORT=3000 go run main.go)
