#!/usr/bin/env bash
# build.sh builds styles
set -ex

dir=$(dirname "${BASH_SOURCE[0]}")
basePath="$dir/../../.."
serverDir="$dir/../.."
outputDir="$serverDir/static"
inputDir="$dir/src"

task="cp $inputDir/*.js $outputDir"

(
  cd "$basePath/watcher" && \
  go run main.go \
  --task="$task" \
  --context="$inputDir" \
  "$inputDir"
)
