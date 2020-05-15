#!/usr/bin/env bash
# build.sh builds a bundle
set -ex

dir=$(dirname "${BASH_SOURCE[0]}")
basePath="$dir/../.."
standalone=${STANDALONE:-true}

set -u
rm -rf "$basePath/web/public"
mkdir -p "$basePath/web/public/static"

pushd "$basePath/web"
  PUBLIC_PATH="$PUBLIC_PATH" \
  COMPILED_PATH="$COMPILED_PATH" \
  ASSET_BASE_URL="$ASSET_BASE_URL" \
    "$dir/setup.sh"

  OUTPUT_PATH="$COMPILED_PATH" \
  ROOT_URL="$ROOT_URL" \
  VERSION="$VERSION" \
    "$basePath"/web/node_modules/.bin/webpack\
      --colors\
      --display-error-details\
      --env.standalone="$standalone"\
      --config "$(realpath "$basePath/web/webpack/prod.config.js")"

  NODE_ENV=PRODUCTION \
  BUNDLE_BASE_URL=$BUNDLE_BASE_URL \
  ASSET_BASE_URL=$ASSET_BASE_URL \
  PUBLIC_PATH=$PUBLIC_PATH \
  COMPILED_PATH=$COMPILED_PATH \
    node "$dir/placeholder.js"

  cp "$COMPILED_PATH"/*.js "$COMPILED_PATH"/*.css "$PUBLIC_PATH"/static

  # clean up compiled
  rm -rf "$basePath"/web/compiled/*
popd
