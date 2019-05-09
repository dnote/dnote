#!/bin/bash
# build.sh builds a production bundle
set -eux

basePath="$GOPATH/src/github.com/dnote/dnote"
standalone=${STANDALONE:-false}

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
      --config "$basePath"/web/webpack/prod.config.js

  NODE_ENV=PRODUCTION \
  BASE_URL=$BASE_URL \
  ASSET_BASE_URL=$ASSET_BASE_URL \
  PUBLIC_PATH=$PUBLIC_PATH \
  COMPILED_PATH=$COMPILED_PATH \
    "$basePath"/web/scripts/placeholder.sh

  cp "$COMPILED_PATH"/*.js "$COMPILED_PATH"/*.css "$PUBLIC_PATH"/dist

  # clean up compiled
  rm -rf "$basePath"/web/compiled/*
popd
