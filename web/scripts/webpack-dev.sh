#!/bin/bash
set -eux

basePath="$GOPATH/src/github.com/dnote/dnote"
appPath="$basePath"/web

(
  cd "$appPath" &&
  PUBLIC_PATH=$PUBLIC_PATH \
  COMPILED_PATH=$COMPILED_PATH \
    "$appPath"/scripts/setup.sh &&

  BASE_URL=$BASE_URL \
  ASSET_BASE_URL=$ASSET_BASE_URL \
  COMPILED_PATH=$COMPILED_PATH \
  PUBLIC_PATH=$PUBLIC_PATH \
    "$appPath"/scripts/placeholder.sh &&

  "$appPath"/node_modules/.bin/webpack-dev-server\
    --env.standalone="$STANDALONE"\
    --config "$appPath"/webpack/dev.config.js
)
