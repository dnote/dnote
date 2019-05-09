#!/bin/bash
# build.sh builds a production bundle
set -eux

basePath="$GOPATH/src/github.com/dnote/dnote"
publicPath="$basePath/web/public"
compiledPath="$basePath/web/compiled"


baseUrl=https://dnote.io
assetBaseUrl=https://dnote.io

# baseUrl=http://localhost:3000
# assetBaseUrl=http://localhost:3000

STANDALONE=true \
BASE_URL=$baseUrl \
ASSET_BASE_URL=$assetBaseUrl \
PUBLIC_PATH=$publicPath \
COMPILED_PATH=$compiledPath \
"$basePath"/web/scripts/build.sh

