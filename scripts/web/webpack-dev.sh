#!/usr/bin/env bash
set -eux

dir=$(dirname "${BASH_SOURCE[0]}")
basePath="$dir/../.."
appPath="$basePath/web"

(
  cd "$appPath" &&
  PUBLIC_PATH=$PUBLIC_PATH \
  COMPILED_PATH=$COMPILED_PATH \
  ASSET_BASE_URL=$ASSET_BASE_URL \
    "$dir/setup.sh" &&

  BUNDLE_BASE_URL=$BUNDLE_BASE_URL
  ASSET_BASE_URL=$ASSET_BASE_URL \
  COMPILED_PATH=$COMPILED_PATH \
  PUBLIC_PATH=$PUBLIC_PATH \
    node "$dir/placeholder.js" &&

  ROOT_URL=$ROOT_URL \
  VERSION="$VERSION" \
  "$appPath"/node_modules/.bin/webpack-dev-server \
    --env.standalone="$STANDALONE" \
    --host "$WEBPACK_HOST" \
    --config "$appPath"/webpack/dev.config.js
)
