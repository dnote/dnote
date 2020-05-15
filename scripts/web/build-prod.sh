#!/usr/bin/env bash
# build.sh builds a production bundle
set -eux

dir=$(dirname "${BASH_SOURCE[0]}")

basePath="$dir/../.."
publicPath="$basePath/web/public"
compiledPath="$basePath/web/compiled"

bundleBaseUrl="/static"
assetBaseUrl="/static"
rootUrl=""

BUNDLE_BASE_URL="$bundleBaseUrl" \
ASSET_BASE_URL="$assetBaseUrl" \
ROOT_URL="$rootUrl" \
PUBLIC_PATH="$publicPath" \
COMPILED_PATH="$compiledPath" \
STANDALONE=true \
VERSION="$VERSION" \
  "$dir/build.sh"
