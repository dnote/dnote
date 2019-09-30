#!/usr/bin/env bash
set -eux

basePath="$GOPATH/src/github.com/dnote/dnote"
appPath="$basePath"/web
rootUrl=$ROOT_URL

echo "here is $rootUrl"

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
  "$appPath"/node_modules/.bin/webpack-dev-server\
    --env.isTest="$IS_TEST"\
    --config "$appPath"/webpack/dev.config.js
)
