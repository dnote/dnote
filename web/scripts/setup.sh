#!/bin/bash
# setup.sh prepares the directory structure and copies static files
set -eux -o pipefail

basePath="$GOPATH/src/github.com/dnote/dnote"
appStaticDir=$PUBLIC_PATH
compiledPath=$COMPILED_PATH

# prepare directories
rm -rf "$compiledPath"
rm -rf "$appStaticDir"
mkdir -p "$compiledPath"
mkdir -p  "$appStaticDir/dist"

# copy the assets and artifacts
cp -r "$basePath"/web/static/* "$appStaticDir/dist"
cp "$basePath/web/index.html" "$appStaticDir"
