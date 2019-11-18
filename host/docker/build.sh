#!/usr/bin/env bash
set -eu

version=$1

# copy over the build artifact to the Docker build context
dir=$(dirname "${BASH_SOURCE[0]}")
projectDir="$dir/../.."
cp "$projectDir/build/server/dnote_server_${version}_linux_amd64.tar.gz" "$dir"

docker build -t dnote/dnote:"$version" --build-arg version="$version" .
