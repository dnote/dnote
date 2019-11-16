#!/usr/bin/env bash
set -eu

version=$1

docker login
docker push dnote/dnote:"$version"
