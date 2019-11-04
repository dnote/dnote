#!/usr/bin/env bash
set -eux

basePath="$GOPATH/src/github.com/dnote/dnote"
appPath="$basePath"/web

(
  cd "$appPath" &&
  PUBLIC_PATH=$PUBLIC_PATH \
  COMPILED_PATH=$COMPILED_PATH \
  ASSET_BASE_URL=$ASSET_BASE_URL \
    "$appPath"/scripts/setup.sh &&

  BUNDLE_BASE_URL=$BUNDLE_BASE_URL
  ASSET_BASE_URL=$ASSET_BASE_URL \
  COMPILED_PATH=$COMPILED_PATH \
  PUBLIC_PATH=$PUBLIC_PATH \
  IS_TEST=true \
    node "$appPath"/scripts/placeholder.js &&

  ROOT_URL=$ROOT_URL \
  VERSION="$VERSION" \
  "$appPath"/node_modules/.bin/webpack-dev-server\
    --env.isTest="$IS_TEST" \
    --host 0.0.0.0 \
    --config "$appPath"/webpack/dev.config.js
)
