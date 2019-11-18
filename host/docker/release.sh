#!/usr/bin/env bash
set -eu

version=$1

docker login
docker tag dnote-test dnote/test:"$version"
docker push dnote/test:"$version"
