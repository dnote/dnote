#!/usr/bin/env bash
# setup.sh prepares the directory structure and copies static files
set -eux -o pipefail

basePath="$GOPATH/src/github.com/dnote/dnote"
publicPath=$PUBLIC_PATH
compiledPath=$COMPILED_PATH
assetBaseUrl=$ASSET_BASE_URL

# prepare directories
rm -rf "$compiledPath"
rm -rf "$publicPath"
mkdir -p "$compiledPath"
mkdir -p  "$publicPath"

# copy the assets and artifacts
cp -r "$basePath"/web/assets/* "$publicPath"

# populate placeholders
assetBaseUrlEscaped=$(echo "$assetBaseUrl" | sed -e 's/[\/&]/\\&/g')
sed -i -e "s/ASSET_BASE_PLACEHOLDER/$assetBaseUrlEscaped/g" "$publicPath"/static/manifest.json
