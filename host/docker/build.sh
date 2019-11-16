#!/usr/bin/env bash
set -eu

version=$1

dir=$(dirname "${BASH_SOURCE[0]}")
projectDir="$dir/../.."
cp "$projectDir/build/server/dnote_server_${version}_linux_amd64.tar.gz" "$dir"

docker build --build-arg version="$version" .
