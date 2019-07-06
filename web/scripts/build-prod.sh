#!/bin/bash
# build.sh builds a production bundle
set -eux

basePath="$GOPATH/src/github.com/dnote/dnote"
publicPath="$basePath/web/public"
compiledPath="$basePath/web/compiled"

baseUrl=""
assetBaseUrl=""

STANDALONE=true \
BASE_URL=$baseUrl \
ASSET_BASE_URL=$assetBaseUrl \
PUBLIC_PATH=$publicPath \
COMPILED_PATH=$compiledPath \
"$basePath"/web/scripts/build.sh
