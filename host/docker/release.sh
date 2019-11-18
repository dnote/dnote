#!/usr/bin/env bash
set -eux

version=$1

docker login

# tag the release
docker tag dnote/dnote:"$version" dnote/dnote:latest

# publish
docker push dnote/dnote:"$version"
docker push dnote/dnote:latest
