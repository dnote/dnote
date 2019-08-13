#!/usr/bin/env bash
# build.sh builds a bundle
set -ex

basePath="$GOPATH/src/github.com/dnote/dnote"

standalone=${STANDALONE:-false}
isTest=${IS_TEST:-false}
baseUrl=$BASE_URL
assetBaseUrl=$ASSET_BASE_URL

set -u
rm -rf "$basePath/web/public"
mkdir -p "$basePath/web/public/dist"

pushd "$basePath/web"
  PUBLIC_PATH="$PUBLIC_PATH" \
  COMPILED_PATH="$COMPILED_PATH" \
    "$basePath"/web/scripts/setup.sh

  OUTPUT_PATH="$COMPILED_PATH" \
    "$basePath"/web/node_modules/.bin/webpack\
      --colors\
      --display-error-details\
      --env.standalone="$standalone"\
      --env.isTest="$isTest"\
      --config "$basePath"/web/webpack/prod.config.js

  NODE_ENV=PRODUCTION \
  BASE_URL=$baseUrl \
  ASSET_BASE_URL=$assetBaseUrl \
  PUBLIC_PATH=$PUBLIC_PATH \
  COMPILED_PATH=$COMPILED_PATH \
    node "$basePath"/web/scripts/placeholder.js

  cp "$COMPILED_PATH"/*.js "$COMPILED_PATH"/*.css "$PUBLIC_PATH"/dist

  # clean up compiled
  rm -rf "$basePath"/web/compiled/*
popd
