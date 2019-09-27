#!/usr/bin/env bash
# build.sh builds a bundle
set -ex

basePath="$GOPATH/src/github.com/dnote/dnote"

isTest=${IS_TEST:-false}
bundleBaseUrl=$BUNDLE_BASE_URL
assetBaseUrl=$ASSET_BASE_URL

set -u
rm -rf "$basePath/web/public"
mkdir -p "$basePath/web/public/static"

pushd "$basePath/web"
  PUBLIC_PATH="$PUBLIC_PATH" \
  COMPILED_PATH="$COMPILED_PATH" \
    "$basePath"/web/scripts/setup.sh

  OUTPUT_PATH="$COMPILED_PATH" \
  ROOT_URL="$ROOT_URL" \
    "$basePath"/web/node_modules/.bin/webpack\
      --colors\
      --display-error-details\
      --env.isTest="$isTest"\
      --config "$basePath"/web/webpack/prod.config.js

  NODE_ENV=PRODUCTION \
  BUNDLE_BASE_URL=$bundleBaseUrl \
  ASSET_BASE_URL=$assetBaseUrl \
  PUBLIC_PATH=$PUBLIC_PATH \
  COMPILED_PATH=$COMPILED_PATH \
    node "$basePath"/web/scripts/placeholder.js

  cp "$COMPILED_PATH"/*.js "$COMPILED_PATH"/*.css "$PUBLIC_PATH"/static

  # clean up compiled
  rm -rf "$basePath"/web/compiled/*
popd
