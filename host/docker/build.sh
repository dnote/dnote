#!/usr/bin/env bash
set -eux

version=$1

dir=$(dirname "${BASH_SOURCE[0]}")
projectDir="$dir/../.."
tarballName="dnote_server_${version}_linux_amd64.tar.gz"

# copy over the build artifact to the Docker build context
cp "$projectDir/build/server/$tarballName" "$dir"

docker build --network=host -t dnote/dnote:"$version" --build-arg tarballName="$tarballName" .
