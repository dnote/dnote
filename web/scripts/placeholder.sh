#!/bin/bash
set -eux

basePath="$GOPATH/src/github.com/dnote/dnote"

if [ -z "$BASE_URL" ]; then
  echo "BASE_URL environment variable is not set"
  exit 1
fi
if [ -z "$ASSET_BASE_URL" ]; then
  echo "ASSET_BASE_URL environment variable is not set"
  exit 1
fi
if [ -z "$PUBLIC_PATH" ]; then
  echo "PUBLIC_PATH environment variable is not set"
  exit 1
fi
if [ -z "$COMPILED_PATH" ]; then
  echo "COMPILED_PATH environment variable is not set"
  exit 1
fi

BASE_URL=$BASE_URL \
ASSET_BASE_URL=$ASSET_BASE_URL \
PUBLIC_PATH=$PUBLIC_PATH \
COMPILED_PATH=$COMPILED_PATH \
  node "$basePath/web/scripts/placeholder.js"
